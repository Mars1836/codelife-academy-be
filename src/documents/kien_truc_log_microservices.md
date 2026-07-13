**TÀI LIỆU KỸ THUẬT | ARCHITECTURE COMPONENT** 

## **TÀI LIỆU: KIẾN TRÚC LƯU TRỮ LOG MICROSERVICES** 

(PRODUCTION-READY ARCHITECTURE - DEEP DIVE) 

## **1. TỔNG QUAN LUỒNG DỮ LIỆU (WORKFLOW)** 

**==> picture [511 x 29] intentionally omitted <==**

**----- Start of picture text -----**<br>
Ứng dụng (stdout) → K8s Node → Fluent Bit → Kafka → Logstash → OpenSearch → Dashboard/Alerting<br>**----- End of picture text -----**<br>


**==> picture [437 x 45] intentionally omitted <==**

**----- Start of picture text -----**<br>
App / K8s Fluent Bit Kafka Logstash OpenSearch<br>stdout JSON DaemonSet Buffer/Queue Transform Storage/ILM<br>**----- End of picture text -----**<br>


## **2. CHI TIẾT CẤU HÌNH VÀ TRÁCH NHIỆM TỪNG THÀNH PHẦN** 

## **2.1 Ứng dụng & K8s (Nguồn phát log)** 

- **Định dạng JSON:** App (NestJS, Python, Go...) chỉ ghi log ra `stdout/stderr` bằng format JSON kèm `trace_id` . 

- **Không dùng File Transport:** Khuyến nghị dùng Console transport (VD: Winston). K8s sẽ tự động gom stdout vào `/var/log/containers/*.log` và tự rotate. Nếu dùng File transport, log lưu trong Pod sẽ mất khi Pod chết/restart. 

- **Instance ID:** Trên K8s, sử dụng luôn **Pod name** làm Instance ID (K8s tự động quản lý và đảm bảo unique), không cần tự sinh UUID. 

## **2.2 Fluent Bit (Log Collector)** 

- **Triển khai:** Chạy dạng `DaemonSet` (1 Pod/Node). Tự động quét và đọc log từ `/var/log/containers/*.log` bằng plugin `tail` . 

- **Chống mất/trùng log:** Cấu hình lưu trữ file offset vào DB nội bộ (VD: `/var/log/flb_kube.db` ). Khi restart, Fluent Bit đọc DB này để chạy tiếp từ dòng log cuối. 

- **Enrich  Metadata:** Sử dụng `kubernetes  filter` để đính  kèm  metadata  như `pod_name` , `namespace` , `labels` vào cấu trúc log. Có thể kết hợp script Lua để mask dữ liệu nhạy cảm (password, token) trước khi xuất 

- ra. 

## **2.3 Kafka (Message Queue / Giảm sốc)** 

- **Cấu hình Cluster:** Tối thiểu 3 broker (Replication factor ≥ 3). Chịu tải write cực cao, ngăn chặn OpenSearch bị crash khi log sinh ra ồ ạt. 

Trang 1 / 3 

**TÀI LIỆU KỸ THUẬT | ARCHITECTURE COMPONENT** 

- **Phân vùng (Partitioning):** Partition theo `service_name` đảm bảo log cùng một service vào cùng partition (đảm bảo tính ordered). 

- **Topic Strategy:** Có thể tạo topic riêng rẽ ( `logs.payment` , `logs.order` ) để dễ quản lý retention, hoặc dùng chung ( `logs.production` ). 

## **2.4 Logstash / OpenSearch Ingestion (Data Processor)** 

- **Trách nhiệm:** Đóng vai trò Kafka Consumer. Parse lại chuẩn `@timestamp` , tách Index linh hoạt dựa trên service label (VD: `logs-%{service}-%{+YYYY.MM.dd}` ). 

- **Làm sạch:** Lọc bỏ các trường dư thừa (host, agent) trước khi lưu. 

## **2.5 OpenSearch (Storage, Indexing & Alerting)** 

- **Index Pattern:** Tách theo ngày và service (VD: `logs-payment-service-2026.06.17` ). 

- 

- **Index Lifecycle Management (ILM):** 

- 

   - _Hot (0-3 ngày):_ SSD, cho phép write/query nhanh. 

   - _Warm (3-15 ngày):_ HDD, read-only. 

   - _Cold (15-30 ngày):_ Chuyển snapshot xuống S3. 

   - 

   - _Delete (>30 ngày):_ Tự động xóa giải phóng dung lượng. 

- **Alerting:** Thiết lập OpenSearch Alerting quét log ERROR liên tục, trigger gửi thông báo qua Slack/Email. 

## **3. HỎI - ĐÁP TRỌNG TÂM (INTERVIEW Q&A)** 

## **Q: Làm sao để trace lỗi của một request khi nó gọi chéo qua nhiều service?** 

**A:** Dùng cơ chế Distributed Tracing. Inject một `trace_id` (thông qua OpenTelemetry/API Gateway) vào HTTP 

Header. Mọi service ghi log đều phải kèm mã này. Lên OpenSearch Dashboard filter theo `trace_id` sẽ ra toàn bộ luồng request. 

## **Q: Nếu Kafka bị down, hệ thống có mất log không?** 

**A:** Không. Kafka cluster thiết lập ≥ 3 broker, 1 broker down vẫn chạy bình thường. Nếu sập toàn bộ, Fluent Bit có cơ chế local buffer và retry policy, log được lưu tạmởK8s Node và sẽ được đẩy bù khi Kafka sống lại. 

## **Q: Nếu Pod Fluent Bit restart hoặc crash, có bị đọc trùng log không?** 

**A:** Không. Fluent Bit sử dụng file DB (VD: SQLite) để lưu lại vị trí offset của file logđang đọc. Khi khởi động lại, nó căn cứ vào file DB này để đọc tiếp tục tại dòng log bị ngắt quãng. 

## **Q: Log sinh ra ồ ạt quá tải thì OpenSearch có chịu nổi không?** 

**A:** Nhờ có Kafka đứng giữa làm "bộ đệm", OpenSearch sẽ không bị sốc tải. Ta có thể scale horizontal OpenSearch Cluster, tăng số lượng Logstash consumer threads và Kafka partitions để tăng tốc độ tiêu thụ log khi cần. 

Trang 2 / 3 

**TÀI LIỆU KỸ THUẬT | ARCHITECTURE COMPONENT** 

## **Q: Khi dùng Winston, tại sao KHÔNG dùng File transport ghi thẳng ra file vật lý?** 

**A:** Container là môi trường phù du (ephemeral). Nếu ghi file vật lý trong Pod, khi Pod chết log sẽ bay màu. Việc ghi ra stdout (Console transport) sẽ đẩy trách nhiệm quản lý file cho K8s. K8s gom log vào thư mục chuẩn của Node, đảm bảo tính bền vững cho Fluent Bit thu thập. 

## **Q: Xử lý dữ liệu nhạy cảm (Mật khẩu, Token) trong log như thế nào?** 

**A:** Tích hợp script Lua (Lua filter) trực tiếp trên Fluent Bit hoặc sử dụng filter của Logstash để mask/ẩn các trường nhạy cảm trước khi dữ liệu được đẩy vào Kafka và OpenSearch. 

Trang 3 / 3 

