# SQLite production/dev runbook

Tài liệu này đi cùng [ADR 0001](../adr/0001-sqlite-primary-database.md). Các command dưới đây là contract vận hành mục tiêu của issue chuyển SQLite; chỉ dùng sau khi implementation của issue đã được merge.

## Môi trường

`APP_ENV` luôn bắt buộc:

- `development`: database mặc định `data/dev.sqlite3` nếu không đặt `SQLITE_DB_PATH`;
- `production`: database mặc định `data/production.sqlite3` nếu không đặt `SQLITE_DB_PATH`.
- `SQLITE_DB_PATH` override cho mọi môi trường, dùng absolute path cho production deployment trên persistent volume.

**Lưu ý về đường dẫn vật lý:** Các giá trị mặc định trên là tên logical. Khi chạy từ `src/api/` (workflow được document), chúng resolve thành `../data/dev.sqlite3` và `../data/production.sqlite3` để đảm bảo database file nằm trong thư mục `data/` ở thư mục gốc repository.

Ví dụ xem `src/api/.env.example`. Database files và SQLite sidecars là runtime artifact, không commit vào repository.

## Development

```sh
cd src/api
# data/dev.sqlite3 được dùng mặc định (resolve thành ../data/dev.sqlite3);
# có thể ghi đè bằng SQLITE_DB_PATH.

# API sẽ migrate và seed khi database rỗng.
rtk go run .
```

Development seed không reset hoặc overwrite dữ liệu đã có. E2E/cmd/test phải dùng file SQLite tạm riêng, chạy migration và seed trong file đó.

## Production deploy

Production không auto-seed và không để API tự thay đổi schema. Quy trình tối thiểu:

1. Đảm bảo persistent volume và quyền đọc/ghi cho service user.
2. Tạo backup trước migration/import.
3. Chạy migration command explicit.
4. Kiểm tra migration hoàn tất và không còn pending.
5. Start API với `APP_ENV=production` (dùng `data/production.sqlite3` mặc định, resolve thành `../data/production.sqlite3` từ `src/api/`, hoặc `SQLITE_DB_PATH` override cho persistent volume).
6. Kiểm tra `/healthz` và `/readyz`.

API fail fast nếu database không mở được, migration/schema chưa sẵn sàng hoặc config không hợp lệ.

## Migration

Migration SQL được versioned, embedded bằng `go:embed` và chạy qua `pressly/goose`.

- Development/test: auto-migrate khi khởi tạo database.
- Production: chạy `cmd/migrate` explicit sau backup.
- Production chỉ forward-only; không tự động chạy `down`.

## Import catalog

Production bootstrap dùng `cmd/importcatalog` với JSON versioned contract dùng chung với dev seed.

- validate toàn bộ file trước khi ghi;
- import trong một transaction all-or-nothing;
- reject duplicate ID mặc định;
- không tự import fixture demo khi API start.

## Backup và restore

Backup chạy ngoài application, không thêm scheduler vào API:

- daily backup;
- backup bắt buộc trước migration/import;
- giữ tối thiểu 7 bản daily và 4 bản weekly;
- lưu backup ngoài database volume và mã hóa;
- định kỳ kiểm tra restore vào một database tạm trước khi coi backup là hợp lệ.

Khi backup database có WAL sidecar, quy trình backup phải dùng cơ chế nhất quán của SQLite và bao gồm/checkpoint các file WAL theo hướng dẫn vận hành của môi trường deploy.

## SQLite runtime baseline

Adapter bật:

- `foreign_keys=ON`;
- `journal_mode=WAL`;
- `busy_timeout`;
- `synchronous=NORMAL`.

MVP giới hạn database/sql pool ở một writer connection. Production chạy một API instance trên một persistent volume; không mount cùng SQLite file cho nhiều API writers/hosts.

## Health và shutdown

- `/healthz`: process còn sống;
- `/readyz`: database mở được, schema đã migrate và query kiểm tra thành công;
- khi nhận SIGINT/SIGTERM, API drain HTTP requests, đóng server rồi đóng SQLite connection.

## Khi nào cần database server khác?

Đánh giá PostgreSQL hoặc lựa chọn khác khi cần nhiều writers/replicas, shared multi-host writes, HA/read replicas, lock contention đáng kể hoặc transaction workflows vượt khả năng SQLite.
