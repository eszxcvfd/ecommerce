# Sàn sản phẩm số thiết kế và kỹ thuật

Bối cảnh này mô tả ngôn ngữ nghiệp vụ của một sàn mua bán tài nguyên thiết kế, bản vẽ kỹ thuật và tài liệu liên quan.

## Language

**Sản phẩm số**:
Một gói tài nguyên kỹ thuật số được đăng trên sàn, bao gồm một hoặc nhiều tệp và được định danh như một đơn vị để khám phá và giao dịch.
_Avoid_: File, Gói tài nguyên số khi nói về đơn vị được đăng bán.

**Tệp**:
Một đơn vị dữ liệu kỹ thuật số thuộc về một Sản phẩm số; có thể là bản vẽ, mô hình, hình ảnh, tài liệu hoặc tệp định dạng nguồn.

**Mô tả ngắn**: Nội dung tóm tắt giúp người mua nhanh chóng hiểu Sản phẩm số trong danh sách và phần đầu trang chi tiết.

**Mô tả chi tiết**: Nội dung đầy đủ giải thích phạm vi, đặc điểm và thông tin sử dụng của Sản phẩm số; khác với Mô tả ngắn.

**Giấy phép**: Thông tin công khai mô tả phạm vi quyền sử dụng đi kèm Sản phẩm số.

**Đánh giá**: Dữ liệu tổng hợp từ các đánh giá của người mua, gồm điểm trung bình và số lượng đánh giá; việc tạo đánh giá thuộc luồng riêng.

**Ngày đăng**: Thời điểm Sản phẩm số được duyệt và công khai trên sàn; khác với thời điểm Người bán cá nhân bắt đầu tạo sản phẩm.

**Sản phẩm đề xuất**: Các Sản phẩm số công khai khác được hiển thị để hỗ trợ khám phá, không bao gồm Sản phẩm số hiện tại.

**Người mua cá nhân**:
Một cá nhân sử dụng sàn để tìm kiếm, mua hoặc tải Sản phẩm số cho nhu cầu cá nhân, học tập hoặc công việc.
_Avoid_: Khách hàng doanh nghiệp khi nói về phạm vi người mua MVP.

**Người bán cá nhân**:
Một cá nhân tạo và đăng Sản phẩm số lên sàn để cung cấp miễn phí hoặc nhận doanh thu từ việc bán sản phẩm.
_Avoid_: Nhà cung cấp, doanh nghiệp khi nói về phạm vi người bán MVP.

**Quyền tải xuống**:
Quyền của người mua được tải lại một Sản phẩm số sau khi tải miễn phí hoặc hoàn tất mua hàng, trong phạm vi giấy phép của sản phẩm.
_Avoid_: Quyền sở hữu, Bản quyền.

**Xu**:
Vật trung gian thanh toán nội bộ của sàn, được người dùng nhận khi nạp VND, dùng để mua Sản phẩm số và được chuyển thành doanh thu cho Người bán cá nhân sau khi kết toán.
_Avoid_: Tiền mặt, VND khi nói về vật trung gian thanh toán.

**Ví xu**:
Nơi lưu giữ số dư Xu của người dùng, bao gồm Xu đã nạp, đã chi, đang chờ thanh toán hoặc đủ điều kiện rút.
_Avoid_: Tài khoản ngân hàng, Ví nội bộ khi nói về đơn vị lưu giữ Xu.

**Xu chưa sử dụng**:
Số Xu người mua đã nạp nhưng chưa dùng để mua Sản phẩm số; có thể được rút về VND theo tỷ lệ 1 Xu = 1 VND và ngưỡng rút của sàn sau thời gian giữ 1 ngày.
_Avoid_: Doanh thu người bán.

**Thị trường MVP**:
Thị trường Việt Nam, nơi người mua cá nhân và người bán cá nhân đầu tiên của sàn hoạt động.
_Avoid_: Thị trường quốc tế trong phạm vi MVP.

**Mục đích sử dụng**:
Mục tiêu mà người mua lựa chọn khi sử dụng Sản phẩm số sau khi có Quyền tải xuống; sàn không xác minh hoặc kiểm soát trực tiếp mục đích sử dụng cuối cùng.
_Avoid_: Cam kết mục đích sử dụng của người mua.

**Phạm vi trách nhiệm của sàn**:
Sàn không phán đoán mục đích sử dụng cuối cùng của người mua, nhưng tiếp nhận và xử lý các báo cáo hợp lệ về nội dung hoặc sản phẩm vi phạm quy định của sàn, pháp luật hoặc bản quyền.
_Avoid_: Miễn trừ toàn bộ trách nhiệm.

**Duyệt sản phẩm**:
Quy trình kiểm tra Sản phẩm số trước khi công khai, tập trung vào nội dung, định dạng, thông tin quyền sử dụng và quy định đăng bán.
_Avoid_: Kiểm soát mục đích sử dụng cuối cùng.

**Báo cáo vi phạm**:
Phản ánh của người dùng hoặc bên liên quan về Sản phẩm số có dấu hiệu vi phạm quy định của sàn, pháp luật hoặc bản quyền, được sàn tiếp nhận để xem xét và xử lý.
_Avoid_: Báo lỗi kỹ thuật.

**Danh mục**:
Nhóm phân loại của Sản phẩm số trên sàn; MVP gồm kiến trúc, cơ khí, điện tử, đồ họa, đồ án và luận văn.
_Avoid_: Quy trình nghiệp vụ riêng theo từng nhóm.

**Tiền tệ MVP**:
VND là tiền pháp định dùng để nạp vào và rút ra; Xu là đơn vị dùng để niêm yết giá, thanh toán và ghi nhận doanh thu trong MVP. Tỷ lệ quy đổi cố định là 1 Xu = 1 VND.
_Avoid_: Ngoại tệ, thanh toán quốc tế trong phạm vi MVP.

**Hoa hồng sàn**:
Sàn giữ 25% số Xu của mỗi giao dịch trả phí; người bán được ghi nhận 75% số Xu còn lại vào doanh thu trước các phí rút tiền nếu có.
_Avoid_: Hoa hồng 10% trong phạm vi MVP.

**Doanh thu chờ**:
75% số Xu của Người bán cá nhân sau một giao dịch trả phí, được giữ trong 7 ngày trước khi kết toán vào tài khoản Xu của người bán.
_Avoid_: Doanh thu khả dụng trước thời hạn chờ.

**Kết toán giao dịch**:
Sau 7 ngày kể từ giao dịch trả phí, sàn hoàn tất việc phân bổ Xu: 25% giữ lại làm Hoa hồng sàn và 75% được cộng vào tài khoản Xu của Người bán cá nhân.
_Avoid_: Cộng doanh thu người bán trước thời hạn kết toán.

**Ngưỡng rút**:
Mỗi yêu cầu rút Xu của người dùng phải có giá trị tối thiểu 100.000 Xu.
_Avoid_: Rút dưới ngưỡng tối thiểu.

**Chính sách không hoàn tiền**:
Giao dịch mua Sản phẩm số không được hoàn tiền sau khi hoàn tất, theo chính sách nghiệp vụ của sàn.
_Avoid_: Hoàn tiền tự động, quyền hoàn tiền mặc định.

**Tài khoản**:
Một định danh người dùng duy nhất có thể thực hiện vai trò Người mua cá nhân và kích hoạt hồ sơ Người bán cá nhân.
_Avoid_: Tài khoản mua và tài khoản bán tách biệt.

**Hồ sơ người bán**:
Phần thông tin và trạng thái bán hàng được kích hoạt trên một Tài khoản khi người dùng muốn đăng Sản phẩm số và nhận doanh thu.
_Avoid_: Tài khoản người bán riêng.

## Architectural Decisions

### Modular monolith với Ports & Adapters thực dụng

**Trạng thái**: Đã quyết định — 2026-07-20

**Quyết định**: Giữ hệ thống ở dạng modular monolith với các module theo vertical slice, hiện tại là `catalog`. Sử dụng Ports & Adapters ở những seam có biến thể thật. Chưa chuyển toàn bộ sang Clean Architecture nhiều tầng và chưa tách thành microservices.

**Lý do**: Domain hiện còn nhỏ và chưa cần deploy/scale từng module độc lập. SQLite là storage chính thức phù hợp với topology một API instance, trong khi việc tạo sớm các lớp `entities`, `usecases`, `interfaces`, `infrastructure` vẫn dễ tạo module nông, boilerplate và mapping không cần thiết.

**Nguyên tắc áp dụng**:
- Giữ `catalog` làm module chính; `main.go` là composition root.
- SQLite adapter nằm trong module `catalog`; in-memory chỉ là adapter test.
- Giữ `CatalogRepository` làm storage port; không để SQLite-specific types/SQL rò vào domain/API.
- Chỉ thêm application/use-case layer khi xuất hiện workflow nghiệp vụ thực sự như duyệt sản phẩm, đơn hàng, thanh toán Xu hoặc quyền tải xuống.
- Không chọn microservices cho đến khi có nhu cầu deploy, scale hoặc ownership độc lập.

**Điều kiện xem xét lại**: Nhiều writers/replicas, shared multi-host writes, lock contention, HA/read replicas, hoặc workflow giao dịch vượt khả năng vận hành an toàn của SQLite.

### SQLite là database chính
- **Trạng thái**: Đã quyết định — 2026-07-20
- **ADR**: [`docs/adr/0001-sqlite-primary-database.md`](docs/adr/0001-sqlite-primary-database.md)
- **Runbook**: [`docs/operations/sqlite.md`](docs/operations/sqlite.md)
- **Tóm tắt**: SQLite là storage mặc định cho development và production, tách database theo môi trường. Production chạy một API instance trên persistent volume; dev/test có thể dùng seed và in-memory adapter theo contract tests.
