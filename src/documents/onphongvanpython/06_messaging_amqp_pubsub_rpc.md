# 6. Messaging: Queue, Pub/Sub, AMQP và RPC

Messaging giúp các thành phần giao tiếp bất đồng bộ, giảm coupling và xử lý tác vụ nền. Tuy nhiên nó làm hệ thống phức tạp hơn vì phải xử lý retry, duplicate message, ordering và quan sát luồng xử lý.

## 1. Giao tiếp đồng bộ và bất đồng bộ

### Đồng bộ

```text
Client → API A → API B → API C → Response
```

API A phải chờ B và C. Ưu điểm là luồng đơn giản, caller nhận kết quả ngay. Nhược điểm là latency cộng dồn và failure lan truyền.

Ứng dụng: kiểm tra số dư trước khi trả kết quả chuyển tiền, lấy thông tin cần thiết để hoàn thành request hiện tại.

### Bất đồng bộ

```text
API → Broker → Worker
```

API publish message rồi có thể trả response. Worker xử lý sau.

Ứng dụng:

- Gửi email.
- Resize ảnh.
- Đồng bộ CRM.
- Sinh báo cáo.
- Gửi webhook.
- Xử lý sự kiện sau thanh toán.

Không nên bất đồng bộ hóa tác vụ mà user bắt buộc phải biết kết quả ngay.

---

## 2. Message broker là gì?

Broker nhận, lưu tạm và phân phối message giữa producer và consumer.

```text
Producer → Exchange/Topic/Queue → Consumer
```

Broker giúp producer không cần biết consumer đang ở đâu. Consumer có thể tạm dừng rồi xử lý message sau nếu broker còn lưu message.

Các sản phẩm phổ biến:

- RabbitMQ: queue và routing linh hoạt, hỗ trợ AMQP.
- Kafka: distributed log, throughput cao, lưu event lâu dài.
- Redis Streams: stream messaging gọn nhẹ trong hệ sinh thái Redis.
- Cloud queue: SQS, Pub/Sub, Service Bus.

---

## 3. Queue và competing consumers

Với queue, một message thường được một consumer trong nhóm xử lý.

```text
             ┌→ Worker 1
Producer → Queue → Worker 2
             └→ Worker 3
```

Khi tăng worker, throughput có thể tăng nếu broker, database và dependency chịu được tải.

Ví dụ message:

```json
{
  "event_id": "evt_1001",
  "type": "send_welcome_email",
  "user_id": 10,
  "email": "hau@example.com"
}
```

Consumer chỉ ack sau khi xử lý thành công.

---

## 4. Pub/Sub

Trong Pub/Sub, một event được nhiều subscriber độc lập nhận.

```text
Order Paid
   ├→ Email subscriber
   ├→ Loyalty subscriber
   ├→ Analytics subscriber
   └→ CRM subscriber
```

Mỗi subscriber có mục đích riêng. Nếu analytics lỗi, email vẫn có thể xử lý bình thường.

Queue tập trung phân phối công việc giữa các worker cùng chức năng. Pub/Sub phát cùng một sự kiện đến nhiều nhóm xử lý khác nhau.

---

## 5. RabbitMQ và AMQP

RabbitMQ thường dùng các khái niệm:

- Producer.
- Exchange.
- Binding.
- Queue.
- Consumer.
- Routing key.

Producer thường publish vào exchange, không publish trực tiếp vào queue.

### Direct exchange

Routing key phải khớp chính xác.

```text
routing key: email.send → email_queue
```

### Topic exchange

Hỗ trợ pattern:

```text
order.*
order.paid
order.#
```

### Fanout exchange

Phát message đến tất cả queue được bind, bỏ qua routing key.

### Headers exchange

Routing dựa trên header. Ít phổ biến hơn direct và topic.

---

## 6. Acknowledgement và durability

Consumer ack để broker biết message đã xử lý xong.

- Auto ack: broker coi message thành công ngay khi gửi; có nguy cơ mất message khi worker crash.
- Manual ack: consumer ack sau khi xử lý thành công.
- Nack/reject: báo xử lý thất bại; có thể requeue hoặc chuyển DLQ.

Để tăng độ bền cần xem đồng thời:

- Queue durable.
- Message persistent.
- Publisher confirm.
- Broker replication phù hợp.

Không có một flag duy nhất biến hệ thống thành “không bao giờ mất message”.

---

## 7. Retry và Dead Letter Queue

Không retry vô hạn ngay lập tức vì có thể tạo retry storm.

Chiến lược:

```text
Main Queue
   ↓ lỗi
Retry 10s
   ↓ lỗi
Retry 1m
   ↓ lỗi
Retry 10m
   ↓ lỗi
Dead Letter Queue
```

DLQ lưu message không xử lý được để điều tra hoặc replay.

Message nên có metadata:

```json
{
  "event_id": "evt_1001",
  "attempt": 3,
  "occurred_at": "2026-07-23T08:00:00Z",
  "correlation_id": "req_abc"
}
```

Phân biệt lỗi:

- Transient: timeout, service tạm lỗi → retry.
- Permanent: email sai format, resource không tồn tại → không retry vô hạn.
- Bug: payload không tương thích → DLQ và cảnh báo.

---

## 8. At-most-once, at-least-once và exactly-once

### At-most-once

Message có thể mất nhưng không xử lý lặp. Phù hợp telemetry không quan trọng.

### At-least-once

Message không dễ mất nhưng có thể xử lý nhiều lần. Đây là mô hình phổ biến, yêu cầu consumer idempotent.

### Exactly-once

“Đúng một lần” end-to-end rất khó. Một số nền tảng cung cấp guarantee trong phạm vi cụ thể, nhưng side effect ngoài hệ thống vẫn cần idempotency.

Ví dụ gửi request thanh toán thành công nhưng worker crash trước ack. Broker gửi lại message; consumer phải nhận ra giao dịch đã xử lý.

---

## 9. Idempotent consumer

Tạo bảng lưu event đã xử lý:

```sql
CREATE TABLE processed_events (
    consumer_name VARCHAR(100) NOT NULL,
    event_id VARCHAR(100) NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (consumer_name, event_id)
);
```

Trong cùng transaction:

1. Insert `(consumer_name, event_id)`.
2. Nếu conflict, bỏ qua vì đã xử lý.
3. Cập nhật nghiệp vụ.
4. Commit.
5. Ack message.

Nếu side effect là gọi API bên ngoài, cần idempotency key ở API đích hoặc cơ chế reconcile.

---

## 10. Outbox pattern

Vấn đề dual write:

```text
1. Commit order vào PostgreSQL thành công
2. Publish order.paid thất bại
```

Dữ liệu đã thay đổi nhưng event bị mất.

Outbox pattern ghi dữ liệu nghiệp vụ và event vào cùng transaction:

```sql
BEGIN;

UPDATE orders SET status = 'paid' WHERE id = 1001;

INSERT INTO outbox_events(id, event_type, payload)
VALUES ('evt_1001', 'order.paid', '{"order_id":1001}');

COMMIT;
```

Một publisher riêng đọc outbox và publish broker. Sau khi publish thành công, đánh dấu event đã gửi.

Outbox không tự loại bỏ duplicate; consumer vẫn cần idempotency.

---

## 11. Ordering

Ordering toàn cục làm giảm khả năng scale. Thường chỉ cần ordering theo entity, ví dụ theo `account_id` hoặc `order_id`.

Kafka partition theo key:

```text
key = order_id
```

Các event cùng order vào cùng partition và giữ thứ tự trong partition.

RabbitMQ với nhiều consumer có thể hoàn thành message không đúng thứ tự dù broker giao lần lượt. Nếu nghiệp vụ phụ thuộc version, consumer nên kiểm tra sequence/version.

---

## 12. RabbitMQ và Kafka

### RabbitMQ

Phù hợp:

- Task queue.
- Routing linh hoạt.
- Request/reply.
- Message thường được xóa sau ack.
- Latency thấp cho workflow nghiệp vụ.

### Kafka

Phù hợp:

- Event streaming.
- Throughput lớn.
- Retention dài.
- Replay event.
- Nhiều consumer group đọc độc lập.

Không nên thay RabbitMQ bằng Kafka chỉ vì Kafka “scale lớn hơn”. Cần xem use case là task distribution hay event log.

---

## 13. Backpressure và prefetch

Nếu producer nhanh hơn consumer, queue tăng liên tục, latency và storage tăng.

Cách xử lý:

- Scale consumer có giới hạn.
- Rate limit producer.
- Batch processing.
- Prefetch hợp lý.
- Giới hạn queue và alert queue depth.
- Tối ưu dependency chậm.

RabbitMQ prefetch giới hạn số message chưa ack trên consumer. Prefetch quá lớn làm một worker giữ quá nhiều message và phân phối không đều.

---

## 14. Poison message

Poison message luôn làm consumer lỗi, ví dụ schema sai hoặc dữ liệu không hợp lệ.

Nếu requeue vô hạn, nó chiếm tài nguyên và làm log nhiễu. Cần giới hạn retry rồi đưa DLQ kèm lý do lỗi.

Consumer cần validate schema trước xử lý.

---

## 15. Schema evolution

Message là contract giữa producer và consumer.

Nguyên tắc:

- Thêm field optional thay vì đổi nghĩa field cũ.
- Không xóa field ngay khi consumer cũ còn chạy.
- Có `event_version`.
- Consumer bỏ qua field không biết.
- Dùng schema registry khi hệ thống lớn.

```json
{
  "event_type": "order.paid",
  "event_version": 2,
  "data": {
    "order_id": 1001,
    "currency": "VND"
  }
}
```

---

## 16. RPC

RPC cho phép gọi procedure trên service khác như gọi hàm từ xa.

Ví dụ gRPC:

```proto
service AccountService {
  rpc GetBalance(GetBalanceRequest) returns (GetBalanceResponse);
}
```

Ưu điểm:

- Contract mạnh.
- Serialization hiệu quả.
- Code generation.
- Hỗ trợ streaming.

Nhược điểm:

- Coupling contract cao hơn event.
- Caller vẫn phụ thuộc availability của server.
- Debug thủ công khó hơn REST JSON.

RPC không phải messaging bất đồng bộ. RPC thường vẫn là request-response đồng bộ dù transport khác HTTP REST.

---

## 17. Timeout, retry và circuit breaker

Với RPC hoặc HTTP giữa service:

- Luôn có timeout.
- Retry chỉ với lỗi phù hợp.
- Dùng exponential backoff và jitter.
- Không retry thao tác không idempotent nếu không có idempotency key.
- Circuit breaker giúp ngừng gọi dependency đang lỗi liên tục.

Retry ở nhiều tầng có thể nhân số request. Ví dụ client retry 3 lần, API retry 3 lần và SDK retry 3 lần có thể tạo tới 27 lần gọi.

---

## 18. Quan sát hệ thống messaging

Metric quan trọng:

- Queue depth.
- Message age.
- Publish rate.
- Consume rate.
- Retry rate.
- DLQ size.
- Processing latency.
- Consumer error rate.

Log nên có:

- event_id.
- correlation_id.
- event_type.
- attempt.
- consumer_name.

Distributed tracing cần propagate trace context qua message header.

---

## 19. Ví dụ FastAPI publish task

Pseudo-code:

```python
@app.post("/users")
async def create_user(payload: CreateUserRequest):
    user = await user_service.create(payload)

    await broker.publish(
        routing_key="email.welcome",
        message={
            "event_id": str(uuid4()),
            "user_id": user.id,
            "email": user.email,
        },
    )

    return user
```

Trong production, nếu cần bảo đảm event không mất sau commit database, nên dùng outbox thay vì publish trực tiếp như ví dụ đơn giản trên.

---

## 20. Câu hỏi phỏng vấn

### Queue khác Pub/Sub thế nào?

Queue thường phân phối một công việc cho một consumer trong nhóm. Pub/Sub phát cùng event cho nhiều nhóm subscriber độc lập.

### Khi nào dùng RabbitMQ, khi nào dùng Kafka?

RabbitMQ phù hợp task queue và routing nghiệp vụ linh hoạt. Kafka phù hợp event stream, throughput lớn, retention và replay. Quyết định dựa trên semantics chứ không chỉ throughput.

### Consumer xử lý xong nhưng chưa ack rồi crash thì sao?

Broker có thể giao lại message. Vì vậy consumer cần idempotent để duplicate không tạo side effect lặp.

### Outbox giải quyết vấn đề gì?

Outbox giải quyết dual write giữa database và broker bằng cách lưu business change và event trong cùng transaction, sau đó publisher gửi event ra broker.

### DLQ dùng để làm gì?

DLQ chứa message đã vượt giới hạn retry hoặc không thể xử lý, giúp tránh vòng lặp vô hạn và hỗ trợ điều tra/replay.

### RPC khác REST?

RPC tập trung procedure và contract, thường dùng code generation. REST tập trung resource và semantics HTTP. Cả hai thường là giao tiếp đồng bộ.

---

## 21. Bài tập thực hành

1. Thiết kế queue gửi email với retry 10 giây, 1 phút và DLQ.
2. Viết idempotent consumer dùng bảng `processed_events`.
3. Thiết kế outbox cho sự kiện `payment.completed`.
4. So sánh RabbitMQ và Kafka cho hệ thống đồng bộ CRM.
5. Mô phỏng worker crash sau khi ghi database nhưng trước ack.
6. Thiết kế metric và alert cho queue backlog.
7. Version hóa schema event mà không phá consumer cũ.

## Checklist

- [ ] Phân biệt sync và async.
- [ ] Phân biệt queue và Pub/Sub.
- [ ] Hiểu exchange, routing key, ack và prefetch.
- [ ] Giải thích được at-least-once và idempotency.
- [ ] Hiểu retry, backoff và DLQ.
- [ ] Giải thích được outbox pattern.
- [ ] So sánh RabbitMQ và Kafka theo use case.
- [ ] Hiểu ordering, schema evolution và observability.
