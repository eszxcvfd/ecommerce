# Sàn sản phẩm số thiết kế và kỹ thuật

Bối cảnh này mô tả ngôn ngữ nghiệp vụ của một sàn mua bán tài nguyên thiết kế, bản vẽ kỹ thuật và tài liệu liên quan.

## Language

**Sản phẩm số**:
Một gói tài nguyên kỹ thuật số được đăng trên sàn, bao gồm một hoặc nhiều tệp và được định danh như một đơn vị để khám phá và giao dịch.
_Avoid_: File, Gói tài nguyên số khi nói về đơn vị được đăng bán.

**Tệp**:
Một đơn vị dữ liệu kỹ thuật số thuộc về một Sản phẩm số; có thể là bản vẽ, mô hình, hình ảnh, tài liệu hoặc tệp định dạng nguồn.

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
