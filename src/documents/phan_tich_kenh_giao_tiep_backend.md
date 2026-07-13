## **PHÂN TÍCH SÂU CÁC KÊNH GIAO TIẾP TRONG KIẾN TRÚC HỆ THỐNG** 

_Tài liệu phỏng vấn cấp độ Senior Backend Engineer | Tiêu chuẩn đánh giá kiến trúc hạ tầng thiết kế hệ thống phân tán_ 

Trong thiết kế hệ thống phân tán và microservices, việc lựa chọn kênh giao tiếp (Communication Channel/ Protocol) phù hợp không đơn thuần là việc chọn một công nghệ, mà là một quyết định kiến trúc ảnh hưởng trực tiếp đến hiệu năng (throughput, latency), khả năng mở rộng (scalability), độ tin cậy (reliability), và tài nguyên phần cứng (CPU, Memory, Network Bandwidth). Dưới đây là phân tích chuyên sâu mang tính thực chiến về ưu, nhược điểm và cơ chế hoạt động của 7 giải pháp giao tiếp cốt lõi. 

## **1. HTTP (Hypertext Transfer Protocol - Trọng tâm HTTP/1.1 và HTTP/2)** 

Giao thức lớpứng dụng (Layer 7) hoạt động theo mô hình Request-Response truyền thống trên nền tảng TCP. 

## **Cơ chế cốt lõi dưới góc nhìn Senior:** 

- **HTTP/1.1:** Sử dụng cơ chế _Keep-Alive_ để tái sử dụng kết nối TCP, nhưng bị giới hạn bởi hiện tượng **Head-of-Line (HoL) Blocking** ở tầng ứng dụng (một request bị nghẽn sẽ chặn toàn bộ các request phía sau trên cùng một kết nối). Trình duyệt buộc phải mở tối đa 6 kết nối TCP song song để giảm thiểu điều này. 

- **HTTP/2:** Giới thiệu **Binary Framing Layer** và cơ chế **Multiplexing** (đa luồng hóa), cho phép gửi đồng thời nhiều request/response trên duy nhất một kết nối TCP, giải quyết triệt để vấn đề HoL ở tầng ứng dụng. Nó cũng tích hợp nén header HPACK và Server Push. 

## **Ưu điểm:** 

- **Tính phổ quát và chuẩn hóa cao:** Được hỗ trợ bởi tất cả các ngôn ngữ, framework. Dễ dàng debug thông qua công cụ cURL, Postman. 

- **tối ưu hóa bộ nhớ đệm (Caching):** Cơ chế kiểm soát bộ nhớ đệm cực kỳ mạnh mẽ thông qua các HTTP Headers chuẩn hóa tại tầng CDN, Reverse Proxy hoặc Client. 

- **Thân thiện với hạ tầng mạng:** Chạy mặc định trên các cổng tiêu chuẩn (80/443), dễ dàngđi qua các hệ thống Firewall và NAT. 

- **Stateless:** Giúp việc scale ngang (Horizontal Scaling) cácứng dụng Backend trở nên cực kỳ đơn giản vì không cần đồng bộ trạng thái. 

Tài liệu Phân Tích Chuyên Sâu Kênh Giao Tiếp 

1 

## **Nhược điểm:** 

- **Độ trễ cao cho các tác vụ thời gian thực (Real-time):** Bản chất là Unidirectional (đơn hướng). Để cập nhật dữ liệu liên tục, client phải dùng Polling hoặc Long-Polling, gây lãng phí tài nguyên. 

- **Overhead của dữ liệu truyền tải:** HTTP/1.1 Headers dạng text chiếm dung lượng lớn, tỷ lệ dữ liệu hữuích trên băng thông có thể rất thấp. 

## **2. HTTPS (Hypertext Transfer Protocol Secure)** 

Bản chất là HTTP chạy trên lớp mã hóa dữ liệu **TLS (Transport Layer Security)** , nằm giữa tầng TCP và HTTP. 

## **Cơ chế cốt lõi dưới góc nhìn Senior:** 

Bảo vệ dữ liệu qua hai giaiđoạn: **Asymmetric Cryptography** (mã hóa bất đối xứng) trong quá trình TLS Handshake để xác thực và trao đổi Session Key; sau đó dùng **Symmetric Cryptography** (mã hóa đối xứng) để mã hóa toàn bộ dữ liệu luân chuyển. 

## **Ưu điểm:** 

- **Bảo mật toàn vẹn dữ liệu:** Đảm bảo 3 yếu tố cốt lõi: Bảo mật (Encryption), Toàn vẹn (Integrity), và Xác thực (Authentication). 

- **Tiêu  chuẩn  bắt  buộc:** Cần  thiết  để sử dụng  các  API  hiện  đại  của  trình  duyệt  (Service  Workers, Geolocation) và hỗ trợ HTTP/2. 

## **Nhược điểm:** 

- **Độ trễ tăng do TLS Handshake:** Việc thiết lập kết nối đòi hỏi thêm các lượt bắt tay (RTTs), làm tăng connection latency ban đầu. 

- **Tiêu tốn tài nguyên CPU xử lý mã hóa:** Ở quy mô lớn, backend bắt buộc phải sử dụng giải pháp **TLS Termination** tại API Gateway/Load Balancer để giảm tải cho app. 

## **3. SSE (Server-Sent Events)** 

Cho phép Server chủ động đẩy dữ liệu (Push Notification) về Client theo thời gian thực dựa trên một kết nối HTTP duy nhất kéo dài bền vững. 

## **Cơ chế cốt lõi dưới góc nhìn Senior:** 

Client gửi HTTP Request có header `Accept: text/event-stream` . Server giữ kết nối mở và gửi dữ liệu dưới dạng luồng text phân tách bằng dòng trống. 

Tài liệu Phân Tích Chuyên Sâu Kênh Giao Tiếp 

2 

## **Ưu điểm:** 

- **Hoạt động mượt mà trên hạ tầng HTTP:** Không cần nâng cấp giao thức phức tạp, xuyên tường lửa và proxy dễ dàng. 

- **Tự động kết nối lại (Built-in Auto-reconnect):** Tích hợp sẵn cơ chế reconnect và Last-Event-ID để không mất luồng tin nhắn khi mạng chập chờn. 

- **Tiết kiệm tài nguyên Server:** Rất tối ưu khiứng dụng chỉ có nhu cầu đẩy dữ liệu một chiều (Server -> Client). 

## **Nhược điểm:** 

- **Giao tiếp một chiều (Unidirectional):** Client không thể gửi dữ liệu lên qua kênh SSE này. 

- **Định dạng giới hạn:** Chỉ hỗ trợ Text/UTF-8. Kém hiệu quả với dữ liệu nhị phân nguyên bản. 

- **Giới hạn kết nối (HTTP/1.1):** Trình duyệt giới hạn tốiđa 6 kết nối song song trên mỗi domain. 

## **4. WebSocket** 

Kênh giao tiếp song công toàn phần ( **Full-Duplex** ) trên duy nhất một kết nối TCP duy trì lâu dàiởtầng 7. 

## **Cơ chế cốt lõi dưới góc nhìn Senior:** 

Bắt đầu bằng HTTP Upgrade Request ( `Upgrade: websocket` ). Sau khi bắt tay, kết nối chuyển sang kênh nhị phân trên TCP độc lập. 

## **Ưu điểm:** 

- **Giao tiếp hai chiều thời gian thực:** Cả Client và Server đều có thể chủ động Push dữ liệu không độ trễ. 

- **Overhead Frame cực thấp:** Mỗi Frame chỉ tốn 2 đến 10 bytes cho Header, siêu tiết kiệm so với HTTP Headers. 

- **Hỗ trợ Native Binary Payload:** Truyền tải thẳng luồng byte nhị phân thô, âm thanh, hìnhảnh. 

## **Nhược điểm:** 

- **Scale Ngang phức tạp:** Vì là kết nối Stateful, khi hệ thống phân tán cần dùng thêm Pub/Sub (Redis/ Kafka) để định tuyến tin nhắn giữa các node server. 

- **Không thể Caching:** Kết nối WebSocket không được đệm qua CDN hay Proxy. 

- **Proxy/Firewall Timeout:** Nếu không duy trì Ping-Pong liên tục, các thiết bị mạng trung gian thường tự động cắt đứt kết nối TCP. 

Tài liệu Phân Tích Chuyên Sâu Kênh Giao Tiếp 

3 

## **5. RPC (Remote Procedure Call - Đại diện tiêu biểu: gRPC)** 

Cho phép một chương trình gọi một hàmởmột máy tính khác (mạng nội bộ/phân tán) như thể gọi hàm cục bộ. 

## **Cơ chế cốt lõi dưới góc nhìn Senior:** 

gRPC sử dụng **HTTP/2** làm giao thức truyền tải kết hợp với mã hóa nhị phân bằng **Protocol Buffers (ProtoBuf)** , tự động sinh code (Stub) từ schema định sẵn (IDL). 

## **Ưu điểm:** 

- **Hiệu năng tốiđa:** Dữ liệu nhị phân siêu đặc, tốc độ Serialization/Deserialization nhanh gấp nhiều lần JSON. 

- **Hợp đồng chặt chẽ (Strict Contract Typing):** Ràng buộc cực mạnh về schema thông qua file `.proto` , hạn chế lỗi tích hợp giữa các team. 

- **Hỗ trợ Streaming 4 chiều:** Unary, Server Streaming, Client Streaming, Bidirectional Streaming. 

## **Nhược điểm:** 

- **Khó debug bằng mắt thường:** Dữ liệu nhị phân cần các công cụ chuyên biệt (gRPCurl) để kiểm thử, không thân thiện với cURL thông thường. 

- **Khả năng tương thích Web Browser kém:** Trình duyệt không can thiệp sâu được vào frame HTTP/2, cần gRPC-Web proxy để hoạt động. 

## **6. TCP (Transmission Control Protocol)** 

Giao thức Giao vận (Layer 4), cung cấp luồng byte đáng tin cậy, hướng kết nối. 

## **Cơ chế cốt lõi dưới góc nhìn Senior:** 

Thực thi bảo đảm qua bắt tay 3 bước ( **3-way handshake** ), xác nhận gói tin (ACK), và các thuật toán kiểm soát dòng chảy (Sliding Window, Congestion Control). 

## **Ưu điểm:** 

- **Độ tin cậy tuyệt đối tại tầng mạng:** Cam kết dữ liệu đúng thứ tự, không mất, không lặp lại. 

- **Tự do thiết kế:** Lập trình viên có quyền thiết kế Protocol tùy biến trên Socket, cực kỳ tối ưu cho các database engine nội bộ. 

## **Nhược điểm:** 

- **Head-of-Line Blocking giao vận:** Chỉ một gói tin bị drop là toàn bộ luồng byte phía sau phải chờ. 

Tài liệu Phân Tích Chuyên Sâu Kênh Giao Tiếp 

4 

- **Độ phức tạp lập trình cực cao:** Tầngứng dụng phải tự định nghĩa ranh giới gói tin (Packet Framing) để tránh hiện tượng dính gói (Sticky Packets). 

## **7. Message Broker (Đại diện tiêu biểu: RabbitMQ, Apache Kafka)** 

Kiến trúc giao tiếp bất đồng bộ ( **Asynchronous Messaging Pattern** ) luân chuyển dữ liệu không kết nối trực tiếp giữa các service. 

## **Cơ chế cốt lõi dưới góc nhìn Senior:** 

- **RabbitMQ (AMQP):** Mô hình "Smart Broker, Dumb Consumer". Định tuyến logic phức tạp, xóa tin ngay khi có xác nhận (ACK). 

- **Kafka:** Mô hình "Dumb Broker, Smart Consumer". Lưu trữ tuần tự dạng Commit Log trên ổ đĩa. Dùng mạnh cho Event Streaming cường độ khổng lồ. 

## **Ưu điểm:** 

- **Loose Coupling tuyệt đối:** Producers và Consumers hoàn toàn không biết đến nhau, dễ dàng scale hệ thống độc lập. 

- **Hấp thụ tải (Traffic Shaving / Load Smoothing):** Bảo vệ hệ thống lõi không bị sập trước luồng traffic đột biến (Spike) bằng cách xếp hàng yêu cầu chờ xử lý. 

- **Khả năng chịu lỗi cao:** Consumer gặp sự cố, dữ liệu vẫn an toàn trên Broker. 

## **Nhược điểm:** 

- **Eventual Consistency:** Dữ liệu không được đồng bộ ngay lập tức, đòi hỏi thiết kế UX và quy trình nghiệp vụ phù hợp. 

- **Chi phí vận hành phức tạp:** Quản trị cụm cluster, xử lý Partition/Replication, và xử lý Dead Letter Queues yêu cầu kiến thức DevOps sâu. 

- **Rủi ro trùng lặp tin nhắn:** Thường phải thiết kế Idempotency kỹ lưỡngởphía Consumer để đề phòng Atleast-once delivery. 

Tài liệu Phân Tích Chuyên Sâu Kênh Giao Tiếp 

5 

## **Bảng Tổng Quan Đánh Giá Nhanh** 

|**Công nghệ**|**Tầng**<br>**OSI**|**Mô hình**|**Độ trễ**|**Use Case Lý Tưởng**|
|---|---|---|---|---|
|**HTTP/1.1**|Layer 7|Request-<br>Response|Trung bình|REST API truyền thống, Web tĩnh, tích hợp<br>Public.|
|**HTTP/2**|Layer 7|Multiplexing|Thấp|hệ thống API hiện đại cần tối ưu kết nối và<br>băng thông.|
|**HTTPS**|Layer 7|Bảo mật qua TLS|Tăng nhẹ|Chuẩn giao tiếp Internet bắt buộc mọi hệ<br>thống.|
|**SSE**|Layer 7|Server Push|Thấp|Dashboard sốliệu, theo dõi trạng thái đơn<br>hàng (Read-only real-time).|
|**WebSocket**|Layer 7|Full-Duplex|Cực thấp|Chat realtime, Game online, Sàn giao dịch tài<br>chính liên tục.|
|**gRPC**|Layer 7|RPC (Binary)|Cực thấp|Giao tiếp nội bộgiữa các microservices cần<br>hiệu suất xử lý cực đại.|
|**TCP**|Layer 4|Luồng Byte đáng<br>tin cậy|Thấp|Tạo protocol riêng biệt (vd: driver database,<br>cache engine).|
|**Message**<br>**Broker**|Ứng<br>dụng|Bất đồng bộ(Pub/<br>Sub)|Phụthuộc<br>hàng đợi|Tách rời logic, xử lý luồng dữ liệu khổng lồ, xử<br>lý email/thanh toán nền.|



Tài liệu Phân Tích Chuyên Sâu Kênh Giao Tiếp 

6 
