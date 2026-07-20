# ADR 0002: Domain metadata cho trang chi tiết Sản phẩm số

- **Trạng thái**: Đã quyết định — 2026-07-20
- **Phạm vi**: module `catalog`, Issue #21 và Issue #5

Trang chi tiết phải giữ đầy đủ yêu cầu của Issue #5 và bổ sung metadata đã chốt: Mô tả ngắn và Mô tả chi tiết là hai trường riêng; Giấy phép là thông tin công khai; Đánh giá chỉ là điểm trung bình và số lượng; Người bán cá nhân chỉ hiện Tên hiển thị; và mỗi Tệp có Tên, Định dạng cùng Dung lượng tính bằng byte.

## Quyết định

- `Ngày đăng` là thời điểm Sản phẩm số được duyệt và công khai, tách khỏi thời điểm tạo (`NgayTao`). Dữ liệu legacy được backfill `Ngày đăng` từ `Ngày tạo` khi chưa có nguồn xuất bản riêng.
- `Sản phẩm đề xuất` là quan hệ dẫn xuất, không phải dữ liệu bắt buộc lưu trên Sản phẩm số. Phiên bản đầu chọn tối đa 4 Sản phẩm số khác cùng Danh mục, chỉ lấy sản phẩm approved/public, sắp xếp mới nhất trước và dùng ID làm tie-break ổn định.
- Chi tiết API trả một view gồm Sản phẩm số và danh sách Sản phẩm đề xuất; không mở rộng luồng mua, tải xuống hoặc tạo đánh giá.
- `CatalogRepository` tiếp tục là storage port. SQLite lưu metadata và Tệp qua adapter; SQL/schema không rò vào domain hoặc HTTP API.

## Hệ quả

- Domain/JSON/SQLite cần thêm mô tả chi tiết, giấy phép, Ngày đăng, hồ sơ Người bán công khai và danh sách Tệp có kích thước.
- Tệp không chứa URL tải xuống trong phạm vi này; quyền tải xuống thuộc các issue riêng.
- Có thể tích hợp hồ sơ Tài khoản/Người bán đầy đủ sau này mà không đưa dữ liệu riêng tư vào public detail.
