# Sàn mua bán Sản phẩm số thiết kế và kỹ thuật

## Problem Statement

Người mua cá nhân tại Việt Nam đang khó tìm, đánh giá và tải các tài nguyên thiết kế, bản vẽ kỹ thuật và tài liệu liên quan trong một nơi thống nhất. Người bán cá nhân cũng thiếu một kênh để đăng Sản phẩm số, nhận doanh thu và quản lý quyền tải xuống.

Sản phẩm hiện cần phục vụ đồng thời việc khám phá, kiểm duyệt, mua/tải, quản lý Xu, báo cáo vi phạm và vận hành nội dung.

## Solution

Xây dựng một marketplace dành cho Việt Nam, nơi người bán cá nhân đăng **Sản phẩm số** chứa một hoặc nhiều **Tệp**, còn người mua cá nhân có thể tìm kiếm, xem chi tiết, tải miễn phí hoặc mua bằng **Xu**.

Sàn cung cấp:

- Danh mục kiến trúc, cơ khí, điện tử, đồ họa, đồ án và luận văn.
- Tìm kiếm, lọc và sắp xếp Sản phẩm số.
- Trang chi tiết với ảnh demo, mô tả, định dạng, dung lượng, giá, người bán và đánh giá.
- Quyền tải xuống sau khi tải miễn phí hoặc mua thành công.
- Ví Xu, nạp VND, thanh toán bằng Xu, lịch sử giao dịch và rút Xu về VND.
- Hồ sơ người bán, quy trình duyệt sản phẩm, doanh thu và yêu cầu rút.
- Yêu thích, bình luận, đánh giá, báo lỗi và Báo cáo vi phạm.
- Khu vực tin tức, hướng dẫn, sự kiện và chính sách.
- Khu vực quản trị cho người dùng, sản phẩm, giao dịch, báo cáo và nội dung.

## User Stories

### Khám phá và tìm kiếm

1. As a visitor, I want to xem danh sách Sản phẩm số công khai, so that I can biết sàn đang cung cấp những gì.
2. As a visitor, I want to xem sáu Danh mục, so that I can bắt đầu từ lĩnh vực phù hợp.
3. As a Người mua cá nhân, I want to tìm kiếm theo tên và mô tả, so that I can tìm được sản phẩm liên quan.
4. As a Người mua cá nhân, I want to lọc theo Danh mục, định dạng, khoảng giá và đánh giá, so that I can thu hẹp kết quả.
5. As a Người mua cá nhân, I want to sắp xếp theo mới nhất, phổ biến, giá và đánh giá, so that I can chọn cách xem phù hợp.
6. As a Người mua cá nhân, I want to xem ảnh demo và thông tin tóm tắt, so that I can đánh giá nhanh trước khi mở chi tiết.
7. As a Người mua cá nhân, I want to xem trang chi tiết, so that I can hiểu nội dung trước khi tải hoặc mua.
8. As a Người mua cá nhân, I want to xem định dạng và dung lượng của từng Tệp, so that I can biết sản phẩm có phù hợp với công cụ của mình không.
9. As a Người mua cá nhân, I want to xem lịch sử đánh giá và điểm đánh giá, so that I can giảm rủi ro khi chọn sản phẩm.
10. As a Người mua cá nhân, I want to xem thông tin Người bán cá nhân, so that I can đánh giá nguồn cung.
11. As a Người mua cá nhân, I want to xem sản phẩm liên quan, so that I can khám phá thêm tài nguyên cùng chủ đề.
12. As a Người mua cá nhân, I want to lưu Sản phẩm số yêu thích, so that I can quay lại sau.

### Tài khoản và quyền tải xuống

13. As a visitor, I want to đăng ký và đăng nhập, so that I can mua, tải và quản lý hoạt động của mình.
14. As a Tài khoản, I want to dùng cùng một thông tin đăng nhập cho cả mua và bán, so that I do not need two accounts.
15. As a Tài khoản, I want to xem lịch sử các Sản phẩm số đã tải hoặc mua, so that I can tải lại khi cần.
16. As a Người mua cá nhân, I want to tải Sản phẩm số miễn phí, so that I can sử dụng tài nguyên không tính phí.
17. As a Người mua cá nhân, I want to mua Sản phẩm số trả phí bằng Xu, so that I can nhận Quyền tải xuống.
18. As a Người mua cá nhân, I want to tải lại sản phẩm đã có Quyền tải xuống, so that I do not lose access after the first download.
19. As a Người mua cá nhân, I want to xem giấy phép và cảnh báo sử dụng, so that I understand the permitted use of the Tệp.
20. As a Người mua cá nhân, I want to gửi bình luận và đánh giá sau khi có Quyền tải xuống, so that I can chia sẻ trải nghiệm thật.
21. As a Người mua cá nhân, I want to báo lỗi sản phẩm, so that I can notify the sàn about a broken or incorrect asset.
22. As a Người mua cá nhân, I want to gửi Báo cáo vi phạm, so that I can notify the sàn about copyright or legal concerns.
23. As a Người mua cá nhân, I want to nhận thông báo về trạng thái giao dịch và báo cáo của mình, so that I know what happened.

### Xu và giao dịch

24. As a người dùng, I want to nạp VND để nhận Xu theo tỷ lệ 1 Xu = 1 VND, so that I can mua Sản phẩm số trong sàn.
25. As a người dùng, I want to xem số dư Ví xu và lịch sử biến động, so that I can kiểm tra tiền của mình.
26. As a người dùng, I want to xem chi tiết mỗi lần nạp, chi, nhận hoặc rút Xu, so that I can đối soát.
27. As a Người mua cá nhân, I want to dùng Xu để thanh toán sản phẩm trả phí, so that I can nhận Quyền tải xuống.
28. As a Người mua cá nhân, I want to rút Xu chưa sử dụng về VND, so that I can lấy lại số dư chưa chi.
29. As a người dùng, I want to biết thời gian giữ Xu nạp trước khi rút, so that I understand when my balance becomes withdrawable.
30. As a người dùng, I want to gửi yêu cầu rút từ 100.000 Xu trở lên, so that I can rút số dư đủ ngưỡng.
31. As a người dùng, I want to xem trạng thái yêu cầu rút, so that I know whether it is pending, approved, rejected or completed.
32. As a Người mua cá nhân, I want to biết rằng giao dịch mua Sản phẩm số không được hoàn tiền, so that I can decide before confirming.
33. As a Người mua cá nhân, I want to biết rằng mục đích sử dụng cuối cùng là trách nhiệm của mình, so that I understand the boundary of the sàn.
34. As a Người bán cá nhân, I want to nhận 75% số Xu của giao dịch sau kết toán, so that I can earn from my work.
35. As a Người bán cá nhân, I want to xem Doanh thu chờ và doanh thu đã kết toán, so that I can plan withdrawals.
36. As a Người bán cá nhân, I want to biết sàn giữ 25% Xu cho mỗi giao dịch trả phí, so that I understand the revenue split.
37. As a Người bán cá nhân, I want to biết giao dịch được kết toán sau 7 ngày, so that I understand when revenue becomes available.
38. As a người dùng, I want to nhận cảnh báo khi Xu nạp chưa đủ thời gian giữ hoặc số dư rút dưới ngưỡng, so that I can avoid invalid withdrawal requests.

### Hồ sơ người bán và đăng sản phẩm

39. As a Tài khoản, I want to kích hoạt Hồ sơ người bán, so that I can đăng Sản phẩm số.
40. As a Người bán cá nhân, I want to tạo Sản phẩm số chứa một hoặc nhiều Tệp, so that I can bán một file hoặc một bộ tài nguyên.
41. As a Người bán cá nhân, I want to nhập tên, mô tả, Danh mục, định dạng, dung lượng, giá và giấy phép, so that buyers have enough information.
42. As a Người bán cá nhân, I want to tải ảnh demo, so that buyers can inspect the product before buying.
43. As a Người bán cá nhân, I want to chọn miễn phí hoặc đặt giá Xu, so that I can choose my offer model.
44. As a Người bán cá nhân, I want to hỗ trợ các định dạng như AutoCAD, SketchUp, Revit, 3ds Max, Photoshop và Illustrator, so that I can publish the assets I create.
45. As a Người bán cá nhân, I want to lưu bản nháp, so that I can hoàn thiện sản phẩm trước khi gửi duyệt.
46. As a Người bán cá nhân, I want to gửi sản phẩm để Duyệt sản phẩm, so that it can become public after review.
47. As a Người bán cá nhân, I want to xem lý do bị từ chối, so that I can correct and resubmit.
48. As a Người bán cá nhân, I want to cập nhật thông tin hoặc Tệp của sản phẩm, so that I can maintain accurate content.
49. As a Người bán cá nhân, I want to xem sản phẩm và doanh thu của mình, so that I can manage my catalog.
50. As a Người bán cá nhân, I want to gửi yêu cầu rút Xu, so that I can convert available revenue to VND.

### Quản trị và kiểm duyệt

51. As an admin, I want to xem hàng đợi sản phẩm chờ duyệt, so that I can kiểm tra trước khi công khai.
52. As an admin, I want to kiểm tra nội dung, định dạng, quyền sử dụng và quy định đăng bán, so that only compliant products become public.
53. As an admin, I want to duyệt hoặc từ chối sản phẩm kèm lý do, so that sellers can understand decisions.
54. As an admin, I want to ẩn hoặc gỡ sản phẩm sau Báo cáo vi phạm hợp lệ, so that the marketplace can respond to violations.
55. As an admin, I want to xem và xử lý báo cáo bản quyền, so that rights holders have a response path.
56. As an admin, I want to quản lý tài khoản và Hồ sơ người bán, so that I can suspend abusive accounts.
57. As an admin, I want to xem giao dịch nạp, mua, kết toán và rút, so that I can đối soát the Xu ledger.
58. As an admin, I want to xem và xử lý yêu cầu rút, so that I can control payout operations.
59. As an admin, I want to theo dõi các giao dịch sắp đến hạn kết toán, so that seller revenue is released at the right time.
60. As an admin, I want to xem báo cáo tổng hợp về doanh thu, hoa hồng, sản phẩm và người dùng, so that I can operate the marketplace.
61. As an admin, I want to quản lý Danh mục, so that the taxonomy can evolve without changing user accounts.
62. As an admin, I want to quản lý tin tức, hướng dẫn, sự kiện và chính sách, so that the sàn has a maintained information area.
63. As an admin, I want to xem audit trail cho các hành động nhạy cảm, so that financial and moderation decisions can be investigated.

### Nội dung và hỗ trợ

64. As a visitor, I want to đọc tin tức, hướng dẫn, sự kiện và chính sách, so that I can learn how to use the sàn.
65. As a user, I want to tìm kiếm nội dung hướng dẫn, so that I can solve common questions.
66. As a user, I want to xem chính sách Xu, tải xuống, bản quyền và không hoàn tiền, so that I can make informed decisions.
67. As a user, I want to nhận thông báo khi sản phẩm yêu thích thay đổi trạng thái, so that I can follow relevant activity.

## Implementation Decisions

- Backend dùng Go và cung cấp HTTP API cho các luồng nghiệp vụ.
- Frontend dùng Nuxt 3, Vue 3, TypeScript và SCSS/CSS Modules; không dùng React hoặc Tailwind.
- API là seam chính cho các hành vi nghiệp vụ; browser/page-level flow là seam cho luồng người dùng chính.
- Một Tài khoản có thể vừa là Người mua cá nhân vừa kích hoạt Hồ sơ người bán. Admin là vai trò phân quyền riêng.
- Sản phẩm số chứa một hoặc nhiều Tệp; Tệp có metadata về định dạng và dung lượng, cùng ảnh demo, mô tả, giá, giấy phép và Danh mục.
- MVP có sáu Danh mục: kiến trúc, cơ khí, điện tử, đồ họa, đồ án và luận văn.
- Sản phẩm phải qua Duyệt sản phẩm trước khi công khai; sau khi công khai, sàn tiếp nhận và xử lý Báo cáo vi phạm.
- Sàn không xác minh mục đích sử dụng cuối cùng của người mua. Sàn vẫn xử lý sản phẩm/nội dung có báo cáo hợp lệ về vi phạm quy định, pháp luật hoặc bản quyền.
- Xu là vật trung gian trong giao dịch. VND được nạp vào để nhận Xu theo tỷ lệ cố định 1 Xu = 1 VND.
- Hoa hồng sàn là 25% số Xu của giao dịch trả phí; 75% còn lại là Doanh thu chờ của người bán.
- Giao dịch trả phí được kết toán sau 7 ngày; 75% Xu được cộng vào tài khoản Xu của người bán sau kết toán.
- Xu chưa sử dụng của người mua được rút về VND sau thời gian giữ 1 ngày. Mỗi yêu cầu rút tối thiểu 100.000 Xu.
- Giao dịch mua Sản phẩm số không được hoàn tiền theo chính sách nghiệp vụ đã chốt.
- Tải miễn phí hoặc mua trả phí đều tạo Quyền tải xuống theo giấy phép của sản phẩm; Quyền tải xuống không đồng nghĩa với chuyển quyền sở hữu hoặc bản quyền.
- Chi tiết cổng nạp/rút, lưu trữ Tệp, xác minh danh tính và cơ chế chuyển VND thực tế là các quyết định cần research trước khi production; chúng phải nằm sau các interface nghiệp vụ và không khóa domain model vào một nhà cung cấp.

## Testing Decisions

- Test chỉ kiểm tra hành vi quan sát được qua HTTP API và browser/page-level flow, không kiểm tra implementation details.
- Go HTTP API là seam chính cho catalog, tìm kiếm, product detail, Quyền tải xuống, Ví xu, giao dịch, hoa hồng, kết toán và rút tiền.
- Nuxt page-level/browser flow kiểm tra luồng: duyệt Danh mục → xem chi tiết → nạp/dùng Xu hoặc tải miễn phí → nhận Quyền tải xuống.
- Các kiểm thử quan trọng phải bao phủ: tỷ lệ 1 Xu = 1 VND, hoa hồng 25%, kết toán sau 7 ngày, giữ Xu nạp 1 ngày, ngưỡng rút 100.000 Xu, không hoàn tiền, sản phẩm phải được duyệt trước khi công khai và báo cáo có thể làm sản phẩm bị ẩn.
- Test thời gian phải sử dụng một clock/time seam được kiểm soát để kiểm tra chuyển trạng thái sau 1 và 7 ngày.
- Repo chưa có prior art cho test; các test mới sẽ là chuẩn đầu tiên của backend Go và frontend Nuxt.
- Test được viết theo vertical slices: một test hành vi đỏ, implementation tối thiểu xanh, sau đó mới refactor.

## Out of Scope

- Thị trường ngoài Việt Nam, ngoại tệ và thanh toán quốc tế.
- Tài khoản doanh nghiệp hoặc tổ chức.
- Tự động chuyển khoản/rút tiền qua ngân hàng; provider cụ thể chưa được chọn.
- Hoàn tiền tự động hoặc cơ chế refund mặc định.
- Kiểm soát mục đích sử dụng cuối cùng của người mua.
- Đảm bảo người mua tuân thủ quy định học thuật; sàn chỉ hiển thị chính sách và xử lý báo cáo hợp lệ.
- Xác minh pháp lý đầy đủ cho mọi sản phẩm trước khi đăng; MVP dùng Duyệt sản phẩm và quy trình Báo cáo vi phạm.
- Hệ thống quảng cáo, đấu giá, thuê bao hoặc chia sẻ doanh thu nhiều tầng.
- Ứng dụng native mobile.
- Tích hợp nhà cung cấp lưu trữ, thanh toán, OCR, antivirus hoặc CDN cụ thể trước khi có research/decision riêng.

## Further Notes

- Mô hình Xu nạp/rút 1:1 có thể tạo nghĩa vụ pháp lý và yêu cầu tuân thủ thanh toán tại Việt Nam. Cần research pháp lý và provider trước khi dùng tiền thật; bản spec này không thay thế tư vấn pháp lý.
- Nội dung cảnh báo về mục đích sử dụng và chính sách không hoàn tiền cần được rà soát pháp lý trước khi phát hành.
- Các quyết định về Xu, kết toán và trách nhiệm của sàn đủ khó đảo ngược để được ghi thành ADR trước khi xây phần tài chính.
- Spec này mô tả phạm vi marketplace; triển khai nên được chia thành các ticket tracer-bullet theo dependency, bắt đầu từ catalog/public browsing trước khi đi vào Xu và payout.
