# 6. Messaging: AMQP, Pub/Sub và RPC

## 1. Synchronous và asynchronous

### Synchronous

Service A gọi Service B và chờ phản hồi.

```text
A → B → response → A
```

Ưu điểm: đơn giản, nhận kết quả ngay.

Nhược điểm: coupling cao, B chậm thì A chậm, dễ tạo lỗi dây chuyền.

### Asynchronous

Service A gửi message rồi tiếp tục xử lý.

```text
A → Broker → Consumer
```

Ưu điểm: giảm coupling, hấp thụ tải đột biến, retry linh hoạt.

Nhược điểm: eventual consistency, khó debug hơn, cần xử lý duplicate.

## 2. Message queue

Các thành phần:

- Producer gửi message.
- Broker lưu và định tuyến.
- Queue giữ message.
- Consumer xử lý message.

Ví dụ:

```text
API tạo đơn hàng
→ lưu order
→ publish order_created
→ worker gửi email
→ worker cập nhật CRM
→ worker cập nhật báo cáo
```

## 3. RabbitMQ

RabbitMQ là message broker thường dùng AMQP.

Luồng:

```text
Producer → Exchange → Queue → Consumer
```

## 4. Exchange

### Direct exchange

Routing key khớp chính xác.

### Topic exchange

Routing theo pattern:

```text
order.created
order.paid
order.cancelled
```

Pattern:

```text
order.*
```

### Fanout exchange

Gửi message đến mọi queue được bind.

### Headers exchange

Routing theo header.

## 5. Queue và routing key

Producer không nhất thiết gửi trực tiếp vào queue; thường publish vào exchange với routing key.

```python
channel.basic_publish(
    exchange="order-events",
    routing_key="order.created",
    body=message,
)
```

## 6. AMQP

AMQP là giao thức messaging quy định cách client và broker giao tiếp.

RabbitMQ hỗ trợ AMQP 0-9-1 phổ biến.

## 7. Acknowledgement

Consumer chỉ ACK sau khi xử lý thành công.

```text
Nhận message
→ xử lý
→ commit DB
→ ACK
```

Nếu consumer chết trước ACK, broker có thể giao lại message.

Điều này dẫn đến khả năng xử lý ít nhất một lần, vì vậy consumer phải idempotent.

## 8. Retry

Không retry vô hạn ngay lập tức.

Chiến lược:

- Giới hạn số lần.
- Exponential backoff.
- Jitter.
- Phân biệt lỗi tạm thời và lỗi vĩnh viễn.
- Sau ngưỡng retry, đưa vào DLQ.

## 9. Dead-letter queue

DLQ lưu message không xử lý được sau nhiều lần retry.

Dùng để:

- Điều tra lỗi.
- Sửa dữ liệu.
- Replay có kiểm soát.
- Tránh message độc chặn queue chính.

## 10. Message duplication

Duplicate có thể xảy ra khi:

1. Consumer cập nhật DB thành công.
2. Consumer chết trước khi ACK.
3. Broker giao lại message.

Giải pháp:

- Message ID duy nhất.
- Bảng processed_messages.
- Unique constraint.
- Idempotency key.
- Upsert.
- Business operation có điều kiện.

Ví dụ:

```sql
INSERT INTO processed_messages(message_id)
VALUES ('msg-001')
ON CONFLICT DO NOTHING;
```

Nếu không insert được vì đã tồn tại, bỏ qua xử lý.

## 11. Idempotent consumer

Một message xử lý nhiều lần vẫn cho kết quả giống một lần.

Ví dụ không idempotent:

```text
balance = balance + 100
```

Ví dụ tốt hơn:

```text
Ghi nhận payment_id duy nhất.
Nếu payment_id đã xử lý thì không cộng lại.
```

## 12. Eventual consistency

Khi tạo order:

- Order DB commit ngay.
- CRM có thể cập nhật sau vài giây.
- Báo cáo có thể cập nhật sau vài phút.

Hệ thống cuối cùng nhất quán, nhưng không đồng bộ tức thời.

## 13. Pub/Sub và queue

### Queue

Một message thường được một consumer trong consumer group xử lý.

### Pub/Sub

Một event có thể được nhiều subscriber độc lập nhận.

Ví dụ `order_created`:

- Email subscriber.
- CRM subscriber.
- Analytics subscriber.

## 14. RPC

RPC cho phép gọi procedure trên service khác như gọi hàm từ xa.

Ví dụ gRPC:

```text
GetUser(user_id) → UserResponse
```

RPC thường có contract rõ và hiệu năng tốt.

REST dùng HTTP resource-oriented:

```text
GET /users/10
```

## 15. Outbox pattern

Vấn đề:

1. Lưu order vào DB thành công.
2. Publish message thất bại.

Kết quả: order tồn tại nhưng không có event.

Outbox pattern:

- Trong cùng transaction, lưu order và lưu outbox event.
- Worker đọc outbox và publish.
- Đánh dấu đã publish.

```text
BEGIN
INSERT order
INSERT outbox_event
COMMIT
```

## Câu hỏi phỏng vấn

### Tại sao không gọi tất cả service trực tiếp?

Vì synchronous call tạo coupling và lỗi dây chuyền. Queue giúp tách rời, retry, hấp thụ tải và xử lý nền.

### Consumer xử lý xong nhưng chưa ACK rồi chết?

Message có thể được giao lại. Consumer cần idempotent.

### Làm sao tránh xử lý hai lần?

Message ID, unique constraint, processed message table, idempotency key và transaction phù hợp.

### Pub/Sub và queue khác nhau?

Queue thường phân phối công việc giữa consumer; Pub/Sub phát cùng event cho nhiều subscriber độc lập.

### RPC khác REST?

RPC tập trung hành động/procedure, contract mạnh và thường hiệu năng cao. REST tập trung resource và tận dụng semantics HTTP.
