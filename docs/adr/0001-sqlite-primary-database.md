# ADR 0001: SQLite là database chính

- **Trạng thái**: Đã quyết định — 2026-07-20
- **Phạm vi đầu tiên**: module `catalog`
- **Supersedes**: giả định trước đây rằng hệ thống chưa có database adapter chính

## Bối cảnh

API hiện tại dùng `memoryRepo` và seed Go trực tiếp từ composition root. Điều này phù hợp cho prototype nhưng không tạo persistence chính thức cho development hoặc production.

Dự án cần một database chính tách biệt theo môi trường:

- **development**: database phục vụ phát triển, thử nghiệm và đánh giá;
- **production**: database chính thức khi deploy, không tự nạp dữ liệu demo.

## Quyết định

Dùng SQLite làm storage mặc định cho cả development và production. Giữ in-memory repository như adapter test cô lập, không dùng làm storage runtime mặc định.

SQLite adapter nằm trong module `catalog`, dùng `database/sql` và `modernc.org/sqlite`. SQLite-specific SQL/types không rò vào domain hoặc HTTP API; `CatalogRepository` tiếp tục là port và trả lỗi storage rõ ràng.

API catalog hiện tại được giữ nguyên contract. Chính sách approved/public được thực thi trong SQL adapter. Schema catalog chuẩn hóa gồm bảng sản phẩm và bảng liên quan cho các định dạng; tên bảng/cột dùng domain Vietnamese `snake_case`.

Search không dấu dùng các cột search đã chuẩn hóa và truy vấn `LIKE`. Tạo các index có mục tiêu cho status, category, price, date, popularity và format; chưa dùng FTS5.

## Môi trường và cấu hình

- `APP_ENV` bắt buộc, nhận `development` hoặc `production`.
- Development có default `SQLITE_DB_PATH=var/dev.sqlite3`.
- Production bắt buộc `SQLITE_DB_PATH` là absolute path ngoài repository, thường nằm trên persistent volume.
- Thiếu cấu hình bắt buộc hoặc database không mở được thì fail fast.
- Không commit database file vào repository.

Production chạy một API instance với một persistent volume và một writer connection. Bật `foreign_keys=ON`, `journal_mode=WAL`, `busy_timeout` và `synchronous=NORMAL`.

## Migration và dữ liệu

Dùng migration SQL versioned, embedded bằng `go:embed` và chạy qua `pressly/goose`.

- Development/test tự chạy migration khi khởi tạo database.
- Production chạy command migration explicit sau khi backup; API chỉ verify schema và fail nếu còn migration pending.
- Production migration forward-only; backup trước migration và dùng migration sửa tiếp thay vì tự động down.
- Dev seed và production import dùng chung JSON versioned contract.
- Development chỉ auto-seed khi database rỗng.
- Production không auto-seed; dùng command import JSON explicit, validate toàn bộ rồi import all-or-nothing, reject duplicate ID mặc định.

## Runtime và vận hành

- Repository methods trả `(result, error)`; lỗi storage trả HTTP 500 với error code generic `storage_unavailable`, chi tiết chỉ ghi server-side.
- Có `/healthz` cho process và `/readyz` kiểm tra database/schema.
- API graceful shutdown HTTP trước, sau đó đóng database connection.
- Backup production chạy ngoài application: daily và bắt buộc trước migration/import; giữ tối thiểu 7 bản daily và 4 bản weekly. Backup phải được mã hóa.
- File database được bảo vệ bằng filesystem/volume permissions; chưa dùng SQLCipher. Khi dữ liệu nhạy cảm xuất hiện, đánh giá lại at-rest encryption.

## Testing

- Contract tests chạy cho cả in-memory và SQLite adapters.
- HTTP tests dùng SQLite adapter làm mặc định cho runtime path.
- E2E dùng SQLite database tạm mỗi run, chạy migration và dev seed, không chạm `var/dev.sqlite3`.
- Test phải kiểm tra migration, public invariant, query/filter/sort, lỗi storage, import atomicity và health/readiness.

## Hệ quả

### Tích cực

- Có persistence chính thức ngay từ catalog.
- Dev/prod tách dữ liệu rõ ràng.
- Port giữ được khả năng thay adapter sau này.
- SQLite phù hợp deployment một instance và không cần CGO.

### Đánh đổi

- Cần migrations, backup/restore, config strict và operational runbook.
- SQLite không phù hợp nhiều API writers, shared multi-host writes, HA/read replicas hoặc lock contention cao.
- Khi các workflow users/orders/Xu xuất hiện, schema và transaction policy sẽ mở rộng trong cùng database.

## Điều kiện xem xét chuyển database server

Đánh giá PostgreSQL hoặc database server khác khi xuất hiện nhiều writers/replicas, cần shared multi-host writes, lock contention, HA/read replicas, hoặc workflow giao dịch vượt khả năng vận hành an toàn của SQLite.
