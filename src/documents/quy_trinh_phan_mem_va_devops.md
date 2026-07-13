## **TÀI LIỆU TỔNG QUAN QUY TRÌNH PHÁT TRIỂN PHẦN MỀM** 

_Phân tích các mô hình SDLC (Waterfall, Agile, Scrum) và Thiết lập Workflow DevOps_ 

## **Lời mở đầu:** 

Trong ngành công nghiệp phần mềm, việc lựa chọn và áp dụng một quy trình phát triển phần mềm (Software Development Life Cycle - SDLC) phù hợp đóng vai trò quyết định đến sự thành bại của dự án. Tài liệu này cung cấp cái nhìn toàn diện về các mô hình quản lý phổ biến bao gồm Waterfall, Agile, Scrum, đồng thời làm rõ cách thiết lập một quy trình làm việc (workflow) chuẩn hóa dưới góc nhìn của kỹ sư DevOps. 

## **1. Mô hình Thác Nước (Waterfall Model)** 

Mô hình Thác nước là phương pháp luận phát triển phần mềm truyền thống và tuyến tính. Trong đó, quá trình phát triển được chia thành các giai đoạn nối tiếp nhau một cách nghiêm ngặt. Giai đoạn sau chỉ được bắt đầu khi giai đoạn trước đó đã hoàn thành hoàn toàn và được phê duyệt. 

## **Các giai đoạn cốt lõi của Waterfall:** 

1. **Phân tích yêu cầu (Requirements):** Thu thập toàn bộ yêu cầu từ khách hàng một cách chi tiết nhất ngay từ đầu dự án. Tài liệu hóa thành SRS (Software Requirement Specification). 

2. **Thiết kế hệ thống (Design):** Thiết kế kiến trúc phần mềm, cấu trúc dữ liệu, sơ đồ hệ thống và giao diện người dùng (UI/UX). 

3. **Triển khai/Lập trình (Implementation):** Đội ngũ lập trình viên tiến hành viết mã nguồn dựa trên tài liệu thiết kế đã duyệt. 

4. **Kiểm thử (Verification/Testing):** Tích hợp các module phần mềm và tiến hành kiểm thử toàn diện (hệ thống, hiệu năng, bảo mật) để tìm và sửa lỗi. 

5. **Triển khai & Bảo trì (Deployment & Maintenance):** Bàn giao sản phẩm cho khách hàng cài đặt, vận hành và tiến hành các hoạt động cập nhật, sửa lỗi phát sinh trong thực tế. 

## **Ưu điểm và Nhược điểm:** 

- **Ưu điểm:** Quản lý dễ dàng nhờ cấu trúc rõ ràng; mốc thời gian và sản phẩm đầu ra của từng giai đoạn được định nghĩa tường minh; phù hợp với các dự án có yêu cầu cố định, ít thay đổi và công nghệ quen thuộc. 

- **Nhược điểm:** Thiếu linh hoạt, cực kỳ khó và tốn kém khi muốn thay đổi yêu cầu ở giai đoạn muộn; khách hàng chỉ thấy được sản phẩm hoàn chỉnh ở cuối chu kỳ; rủi ro cao nếu việc phân tích yêu cầu ban đầu có sai sót. 

Tài liệu Quy trình Phần mềm & DevOps 

Trang 1 

## **2. Triết lý Agile (Agile Methodology)** 

Nhằm khắc phục những hạn chế của Waterfall trong kỷ nguyên công nghệ thay đổi nhanh chóng, triết lý Agile (Phát triển phần mềm linh hoạt) đã ra đời vào năm 2001 thông qua "Tuyên ngôn Agile" (Agile Manifesto). Agile không phải là một quy trình cụ thể, mà là một tập hợp các nguyên lý định hướng cho việc phát triển phần mềm dựa trên tính lặp (iterative) và tăng trưởng (incremental). 

## **4 Tôn chỉ cốt lõi của Agile:** 

- Cá nhân và sự tương tác quan trọng hơn quy trình và công cụ. 

- Phần mềm chạy tốt quan trọng hơn tài liệu đầy đủ. 

- Cộng tác với khách hàng quan trọng hơn đàm phán hợp đồng. 

- Phản hồi với thay đổi quan trọng hơn việcđi theo một kế hoạch vạch sẵn. 

Trong Agile, dự án được chia nhỏ thành các chu kỳ ngắn (vài tuần). Mỗi chu kỳ đềuđi qua các bước lập kế hoạch, phân tích, thiết kế, lập trình và kiểm thử để cho ra một phần tăng trưởng sản phẩm (increment) có thể chạy được và bàn giao ngay cho khách hàng lấy phản hồi. 

## **3. Khung Làm Việc Scrum (Scrum Framework)** 

Scrum là khung làm việc (framework) phổ biến nhất được xây dựng dựa trên triết lý Agile. Scrum định nghĩa rõ ràng các vai trò, sự kiện và tạo tác nhằm giúp đội ngũ phát triển cộng tác hiệu quả tốiđa. 

## **Các Vai trò trong Scrum (Scrum Roles):** 

- **Product Owner (PO):** Người sở hữu sản phẩm, chịu trách nhiệm tối ưu hóa giá trị của sản phẩm, quản lý và ưu tiên các hạng mục trong Product Backlog dựa trên nhu cầu của khách hàng và doanh nghiệp. 

- **Scrum Master (SM):** Người đảm bảo đội ngũ hiểu và tuân thủ các nguyên lý của Scrum; đóng vai trò là người hỗ trợ (facilitator), loại bỏ các rào cản cản trở tiến độ của đội sản xuất. 

- **Đội ngũ Phát triển (Developers):** Nhóm liên chức năng (cross-functional) bao gồm lập trình viên, tester, designer, có toàn quyền tự quản lý để chuyển đổi các hạng mục backlog thành sản phẩm chạy được sau mỗi chu kỳ. 

## **Các Sự kiện trong Scrum (Scrum Events):** 

- **Sprint (Chu kỳ lặp):** Khung thời gian cố định (thường từ 1 đến 4 tuần) nơi một phần tăng trưởng sản phẩm có khả năng phát hành được tạo ra. 

- **Sprint Planning (Lập kế hoạch Sprint):** Cuộc họp đầu Sprint để xác định mục tiêu Sprint và lựa chọn các công việc từ Product Backlog đưa vào Sprint Backlog. 

- **Daily Scrum (Họp hằng ngày):** Cuộc họp ngắn tốiđa 15 phút mỗi ngày để đội ngũ cập nhật tiến độ, kế hoạch trong ngày và nêu ra các khó khăn (nếu có). 

Tài liệu Quy trình Phần mềm & DevOps 

Trang 2 

- **Sprint Review (Sơ kết Sprint):** Diễn raởcuối Sprint để trình diễn phần sản phẩm đã hoàn thành cho các bên liên quan và thu thập phản hồi. 

- **Sprint Retrospective (Cải tiến Sprint):** Cuộc họp nội bộ đội ngũ sau Sprint Review để đánh giá lại quy trình làm việc, tìm ra cácđiểm tốt và điểm cần cải tiến cho Sprint tiếp theo. 

## **Các Tạo tác trong Scrum (Scrum Artifacts):** 

- **Product Backlog:** Danh sách tập trung, liên tục cập nhật chứa tất cả các tính năng, yêu cầu, cải tiến của sản phẩm. 

- **Sprint Backlog:** Tập hợp các hạng mục được chọn từ Product Backlog để phát triển trong Sprint hiện tại kèm theo kế hoạch triển khai. 

- **Increment (Phần tăng trưởng):** Tổng hợp tất cả các hạng mục backlog đã hoàn thành (đạt định nghĩa hoàn thành - Definition of Done) trong Sprint hiện tại và các Sprint trước đó. 

## **Bảng so sánh tóm tắt Waterfall và Agile/Scrum:** 

|**TIÊU CHÍ**|**MÔ HÌNH WATERFALL**|**KHUNG LÀM VIỆC SCRUM**<br>**(AGILE)**|
|---|---|---|
|**Cách tiếp cận**|Tuyến tính, tuần tự từng giai đoạn cố định.|Lặp đi lặp lại, tăng trưởng theo từng chu<br>kỳ ngắn.|
|**Quản lý yêu cầu**|Định nghĩa nghiêm ngặt ngay từđầu dự<br>án.|Linh hoạt, có thểthay đổi và cập nhật liên<br>tục.|
|**Sự tham gia của**<br>**khách hàng**|chủ yếu ở giai đoạn đầu (yêu cầu) và cuối<br>(bàn giao).|Liên tục tham gia qua các buổi Review<br>cuối mỗi Sprint.|
|**Kiểm thử (Testing)**|Thực hiện tập trung sau khi hoàn thành<br>code.|Thực hiện liên tục và đồng thời trong từng<br>Sprint.|



## **4. Quy Trình Thiết Lập Workflow Do DevOps Đảm Nhiệm** 

Nếu Agile/Scrum tối ưu hóa quy trình quản lý và giao tiếp giữa con người, thì DevOps (Development & Operations) là cầu nối kỹ thuật giúp tự động hóa và tối ưu hóa việc chuyển giao phần mềm từ môi trường phát triển (Dev) sang môi trường vận hànhổn định (Ops). Một kỹ sư DevOps chịu trách nhiệm thiết lập toàn bộ hạ tầng kỹ thuật, tự động hóa chuỗi cungứng phần mềm thông qua pipeline CI/CD (Continuous Integration / Continuous Deployment). 

Dưới đây là quy trình chi tiết các bước thiết lập một Workflow DevOps tiêu chuẩn do kỹ sư DevOps thực hiện: 

Tài liệu Quy trình Phần mềm & DevOps 

Trang 3 

## **Bước 1: Thiết lập Hệ thống Quản lý Phiên bản & Chiến lược Phân nhánh (Git branching strategy)** 

DevOps phối hợp với Tech Lead để thiết lập hệ thống lưu trữ mã nguồn (GitHub, GitLab, Bitbucket) và chuẩn hóa cách thức quản lý code của lập trình viên: 

- Áp dụng các mô hình phân nhánh như **Gitflow** hoặc **Trunk-based Development** . 

- Cấu hình các luật bảo vệ nhánh chính ( **Branch Protection Rules** ), yêu cầu tối thiểu số lượng người phê duyệt (Pull Request Approval) và bắt buộc phải pass qua kiểm tra tự động trước khi merge code vào nhánh `main` hoặc `develop` . 

## **Bước 2: Xây dựng Pipeline Tích hợp Tự động (Continuous Integration - CI)** 

Khi lập trình viên đẩy code lên Git, hệ thống CI (Jenkins, GitLab CI, GitHub Actions) sẽ tự động kích hoạt một chuỗi các tác vụ (Job) nhằm đảm bảo code mới không làm hỏng hệ thống hiện tại: 

- **Linting & Code Style Check:** Kiểm tra cú pháp, định dạng code theo chuẩn chung của dự án. 

- **Static Application Security Testing (SAST):** Sử dụng công cụ (như SonarQube) để quét mã nguồn, phát hiện các lỗ hổng bảo mật hoặc code smell sớm. 

- **Automated Unit Testing:** Tự động chạy toàn bộ các bài kiểm thử đơn vị (Unit Test). Nếu tỷ lệ bao phủ (Coverage) không đạt ngưỡng quy định, pipeline sẽ bị hủy. 

- **Build Artifact/Docker Image:** Đóng góiứng dụng thành các file thực thi hoặc Docker Image. Tiến hành quét bảo mật lỗ hổng của các thư viện phụ thuộc (Dependency scanning - Trivy, Snyk). 

- **Push Image to Registry:** Đẩy Docker Image đã kiểm tra an toàn lên các kho lưu trữ tập trung (Docker Hub, AWS ECR, GitLab Container Registry) và gắn tag phiên bản rõ ràng. 

## **Bước 3: Tự động hóa Hạ tầng dưới dạng Mã nguồn (Infrastructure as Code - IaC)** 

Kỹ sư DevOps không khởi tạo máy chủ hay dịch vụ cloud theo cách thủ công. Thay vào đó, toàn bộ hạ tầng (mạng VPC, cụm Kubernetes, cơ sở dữ liệu) được định nghĩa hoàn toàn bằng mã nguồn: 

- Sử dụng các công cụ như **Terraform** , OpenTofu hoặc Ansible để viết script định nghĩa hạ tầng. 

- Mã nguồn hạ tầng cũng được quản lý trên Git, cho phép tracking, rollback và tái cấu trúc môi trường (Staging, UAT, Production) một cách đồng nhất, nhanh chóng và tránh sai sót do con người. 

## **Bước 4: Quản lý Cấu hình và Bảo mật Bí mật (Configuration & Secrets Management)** 

Tách biệt hoàn toàn mã nguồn củaứng dụng khỏi các thông tin cấu hình và thông tin nhạy cảm: 

- Thiết lập hệ thống quản lý biến môi trường theo từng môi trường đích riêng biệt. 

- Sử dụng các công cụ quản lý bảo mật tập trung như **HashiCorp Vault** , AWS Secrets Manager hoặc Kubernetes Secrets để mã hóa và cấp phát động các thông tin nhạy cảm (Database Credentials, API Keys, Private Keys) choứng dụng khi chạy, tuyệt đối không hardcode vào mã nguồn. 

Tài liệu Quy trình Phần mềm & DevOps 

Trang 4 

## **Bước 5: Xây dựng Pipeline Triển khai Tự động (Continuous Delivery / Deployment - CD)** 

Sau khi CI hoàn thành và cho ra Artifact an toàn, DevOps thiết lập quy trình đưaứng dụng lên các môi trường máy chủ: 

- **Áp dụng mô hình GitOps (Xu hướng hiện đại):** Sử dụng các công cụ như **ArgoCD** hoặc FluxCD. Mọi trạng thái mong muốn của hệ thống trên Kubernetes được định nghĩa trong Git. Công cụ GitOps sẽ tự động đồng bộ (sync) trạng thái của cụm máy chủ thực tế trùng khớp với Git. 

- **Chiến lược triển khai giảm thiểu Downtime:** Cấu hình các cơ chế như **Rolling Update** , **Blue-Green Deployment** (chạy song song hai môi trường cũ-mới để chuyển đổi traffic), hoặc **Canary Deployment** (triển khai cho một nhóm nhỏ người dùng trước để kiểm tra lỗi) nhằm đảm bảo dịch vụ không bị giánđoạn khi cập nhật phiên bản mới. 

## **Bước 6: Thiết lập Hệ thống Giám sát & Cảnh báo (Monitoring, Logging & Alerting)** 

Công việc của DevOps không dừng lại sau khi deploy thành công. Để đảm bảo hệ thống vận hànhổn định bền vững, DevOps thiết lập hệ thống quan sát toàn diện: 

- **Centralized Logging:** Thu thập toàn bộ log của cácứng dụng và máy chủ về một nơi tập trung bằng các bộ công cụ như ELK Stack (Elasticsearch, Logstash, Kibana) hoặc LGTM Stack (Loki, Grafana, Promtail). 

- **Metrics Monitoring:** Thu thập các chỉ số phần cứng (CPU, RAM, Disk) và chỉ số ứng dụng (Request/second, Error rate, Latency) thông qua **Prometheus** và hiển thị trực quan lên dashboard của **Grafana** . 

- **Alerting:** Cấu hình các ngưỡng cảnh báo (ví dụ: CPU > 85% kéo dài 5 phút, hoặc tỷ lệ lỗi 5xx tăng đột biến). Hệ thống tự động gửi tin nhắn cảnh báo qua Slack, Telegram, Discord hoặc cuộc gọi khẩn cấp (PagerDuty) đến đội ngũ kỹ sư trực ca (SRE/DevOps) để xử lý kịp thời. 

## **Kết luận** 

Sự kết hợp nhuần nhuyễn giữa phương pháp quản lý linh hoạt của **Agile/Scrum** và năng lực tự động hóa hạ tầng, kiểm thử, triển khai của **DevOps** tạo nên một bộ máy sản xuất phần mềm hiện đại và tối ưu. Quy trình này giúp doanh nghiệp rút ngắn tốiđa thời gian đưa sản phẩm ra thị trường (Time-to-Market), giảm thiểu lỗi con người, tăng cường tính bảo mật và đảm bảo hệ thống luôn vận hànhởtrạng tháiổn định và sẵn sàng cao nhất. 

Tài liệu Quy trình Phần mềm & DevOps 

Trang 5 
