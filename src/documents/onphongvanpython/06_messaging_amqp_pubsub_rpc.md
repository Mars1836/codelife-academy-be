# 6. Messaging: RabbitMQ, Kafka, Pub/Sub, AMQP và RPC

Messaging giúp các thành phần giao tiếp bất đồng bộ, giảm coupling và xử lý tác vụ nền. Đổi lại, hệ thống phải xử lý retry, duplicate message, ordering, backpressure, schema evolution và observability.

## Mục lục

1. Giao tiếp đồng bộ và bất đồng bộ
2. Queue và Pub/Sub
3. RabbitMQ và AMQP
4. Kafka và distributed log
5. So sánh Kafka với RabbitMQ
6. Delivery guarantee và idempotency
7. Retry, DLQ và poison message
8. Ordering, partition và consumer group
9. Outbox pattern
10. Schema evolution
11. Kết nối RabbitMQ từ Python
12. Kết nối Kafka từ Python
13. Docker Compose cho RabbitMQ và Kafka
14. Monitoring và vận hành
15. RPC và REST
16. Câu hỏi phỏng vấn
17. Bài tập thực hành

---

## 1. Giao tiếp đồng bộ và bất đồng bộ

### Đồng bộ

```text
Client → API A → API B → Response
```

API A phải chờ API B trả kết quả.

Ưu điểm:

- Dễ hiểu.
- Caller nhận kết quả ngay.
- Phù hợp tác vụ cần kết quả trong request hiện tại.

Nhược điểm:

- Latency cộng dồn.
- Dependency lỗi có thể gây lỗi dây chuyền.
- Coupling cao hơn.

### Bất đồng bộ

```text
API → Broker → Consumer/Worker
```

API publish message rồi tiếp tục xử lý hoặc trả response.

Ứng dụng thực tế:

- Gửi email.
- Đồng bộ CRM.
- Resize ảnh.
- Sinh báo cáo.
- Gửi webhook.
- Ghi analytics event.
- Xử lý sự kiện sau thanh toán.

Không nên chuyển sang async nếu người dùng bắt buộc phải nhận kết quả ngay, ví dụ xác nhận giao dịch có được chấp nhận hay không.

---

## 2. Queue và Pub/Sub

### Queue

Một message thường được một consumer trong nhóm xử lý.

```text
             ┌→ Worker 1
Producer → Queue → Worker 2
             └→ Worker 3
```

Phù hợp:

- Gửi email.
- Xử lý ảnh.
- Chạy background job.
- Phân phối task giữa nhiều worker.

### Pub/Sub

Một event được nhiều subscriber độc lập nhận.

```text
payment.completed
   ├→ Email consumer
   ├→ CRM consumer
   ├→ Analytics consumer
   └→ Loyalty consumer
```

Queue tập trung vào **phân phối công việc**. Pub/Sub tập trung vào **phát sự kiện cho nhiều hệ thống quan tâm**.

RabbitMQ và Kafka đều có thể triển khai các mô hình này, nhưng cách hoạt động và điểm mạnh khác nhau.

---

## 3. RabbitMQ và AMQP

RabbitMQ là message broker thường dùng cho task queue và routing nghiệp vụ.

Luồng cơ bản:

```text
Producer → Exchange → Queue → Consumer
```

Các thành phần:

- Producer: gửi message.
- Exchange: định tuyến message.
- Binding: quy tắc nối exchange với queue.
- Queue: lưu message chờ xử lý.
- Consumer: đọc và xử lý message.
- Routing key: khóa dùng để định tuyến.

### Các loại exchange

#### Direct exchange

Routing key khớp chính xác.

```text
email.send → email_queue
```

#### Topic exchange

Routing theo pattern.

```text
order.created
order.paid
order.*
order.#
```

#### Fanout exchange

Gửi message tới mọi queue được bind.

#### Headers exchange

Định tuyến theo header.

### Acknowledgement

Consumer chỉ ACK sau khi xử lý thành công.

```text
Nhận message
→ xử lý
→ commit database
→ ACK
```

Nếu worker crash trước ACK, broker có thể giao lại message.

### Durability

Để giảm nguy cơ mất message cần kết hợp:

- Durable queue.
- Persistent message.
- Publisher confirm.
- Broker replication phù hợp.
- Backup và monitoring.

Không có một cấu hình đơn lẻ bảo đảm “không bao giờ mất message”.

### Prefetch

Prefetch giới hạn số message chưa ACK mà consumer được giữ.

Prefetch quá lớn có thể khiến một worker giữ quá nhiều message, phân phối không đều và tăng thời gian chờ của message khác.

---

## 4. Kafka và distributed log

Kafka là nền tảng event streaming dựa trên distributed append-only log.

Luồng cơ bản:

```text
Producer → Topic → Partition → Consumer Group
```

### Topic

Topic là luồng sự kiện có tên, ví dụ:

```text
payment-events
order-events
user-events
```

### Partition

Mỗi topic được chia thành một hoặc nhiều partition.

```text
order-events
   ├── partition 0
   ├── partition 1
   └── partition 2
```

Mỗi partition là một log có thứ tự.

Kafka chỉ bảo đảm ordering **bên trong một partition**, không bảo đảm thứ tự toàn topic.

### Offset

Mỗi record trong partition có offset.

```text
partition 0:
0 → 1 → 2 → 3 → 4
```

Consumer lưu offset để biết đã đọc đến đâu.

Khác RabbitMQ, Kafka thường không xóa record ngay sau khi consumer đọc. Record được giữ theo retention policy.

### Consumer group

Các consumer cùng `group.id` chia nhau các partition.

```text
Topic có 3 partition
Consumer group có 3 consumer

Consumer 1 → partition 0
Consumer 2 → partition 1
Consumer 3 → partition 2
```

Trong cùng consumer group, một partition tại một thời điểm chỉ được một consumer xử lý.

Nếu có 5 consumer nhưng topic chỉ có 3 partition, 2 consumer sẽ không nhận partition.

### Nhiều consumer group

```text
order-events
   ├→ group email-service
   ├→ group analytics-service
   └→ group crm-service
```

Mỗi group đọc độc lập và có offset riêng.

Đây là điểm mạnh quan trọng của Kafka cho event-driven architecture.

### Retention

Kafka giữ record theo:

- Thời gian, ví dụ 7 ngày.
- Dung lượng.
- Log compaction theo key.

Nhờ retention, consumer có thể replay dữ liệu cũ.

### Key và ordering

Producer thường gửi key:

```text
key = order_id
```

Các event có cùng key được đưa vào cùng partition, nhờ đó giữ thứ tự theo order.

```text
order.created
→ order.paid
→ order.shipped
```

### Replication

Partition có leader và replica trên các broker khác nhau.

Producer và consumer làm việc với leader. Replica giúp tăng khả năng chịu lỗi.

Các khái niệm cần biết:

- Replication factor.
- Leader.
- Follower replica.
- In-sync replica.
- `acks` của producer.

### Producer `acks`

- `acks=0`: không chờ broker xác nhận, nhanh nhưng dễ mất message hơn.
- `acks=1`: leader xác nhận.
- `acks=all`: chờ toàn bộ in-sync replica cần thiết xác nhận.

Với event quan trọng thường dùng `acks=all`, idempotent producer và cấu hình replication phù hợp.

---

## 5. So sánh Kafka với RabbitMQ

| Tiêu chí | RabbitMQ | Kafka |
|---|---|---|
| Mô hình chính | Message broker và task queue | Distributed event log và streaming platform |
| Đơn vị tổ chức | Exchange, queue, binding | Topic, partition, offset |
| Sau khi consumer xử lý | Message thường bị xóa sau ACK | Record vẫn được giữ theo retention |
| Replay | Không phải điểm mạnh mặc định | Hỗ trợ tự nhiên bằng cách reset offset |
| Routing | Rất linh hoạt qua exchange và routing key | Chủ yếu theo topic và partition key |
| Consumer | Broker push message đến consumer | Consumer chủ động poll record |
| Scale consumer | Thêm worker theo queue | Bị giới hạn bởi số partition trong một group |
| Ordering | Có thể giữ thứ tự trong queue nhưng nhiều consumer làm completion lệch thứ tự | Bảo đảm thứ tự trong từng partition |
| Throughput | Tốt cho workflow và task queue | Rất tốt cho event stream throughput lớn |
| Retention dài | Không phải use case chính | Là tính năng cốt lõi |
| Nhiều nhóm đọc độc lập | Cần queue riêng cho từng subscriber | Consumer group hỗ trợ tự nhiên |
| Retry/DLQ | Hỗ trợ rất tự nhiên bằng queue/exchange | Thường tự thiết kế retry topic và DLQ topic |
| Request/reply | Hỗ trợ thuận tiện | Không phải use case chính |
| Độ phức tạp vận hành | Thường đơn giản hơn ở quy mô nhỏ | Cao hơn vì partition, replication, rebalance, lag |
| Use case điển hình | Email, task nền, workflow nghiệp vụ | Event streaming, audit stream, CDC, analytics, replay |

### Khi nên dùng RabbitMQ?

- Gửi email hoặc notification.
- Background task.
- Workflow cần routing linh hoạt.
- Retry theo thời gian và DLQ rõ ràng.
- Request/reply.
- Message thường chỉ cần xử lý rồi bỏ.
- Hệ thống nhỏ hoặc trung bình cần vận hành đơn giản.

### Khi nên dùng Kafka?

- Event cần lưu nhiều ngày hoặc nhiều tháng.
- Cần replay để rebuild search index hoặc analytics.
- Nhiều service độc lập cần đọc cùng event.
- Throughput lớn.
- CDC và data pipeline.
- Audit/event log.
- Stream processing.
- Cần ordering theo entity thông qua partition key.

### Ví dụ chọn RabbitMQ

```text
API tạo tài khoản
→ queue send_welcome_email
→ một worker gửi email
→ thành công thì ACK
```

Message không cần lưu lâu sau khi gửi email thành công.

### Ví dụ chọn Kafka

```text
payment.completed
   ├→ fraud detection
   ├→ ledger projection
   ├→ notification
   ├→ analytics
   └→ data warehouse
```

Nhiều hệ thống đọc độc lập và có thể cần replay.

### Có thể dùng cả hai không?

Có.

Ví dụ:

```text
Kafka: lưu business event lâu dài
RabbitMQ: thực thi task cụ thể
```

Luồng:

```text
payment.completed trên Kafka
→ notification service đọc event
→ tạo task gửi email trong RabbitMQ
→ email worker xử lý
```

Không nên dùng cả hai nếu hệ thống chưa có nhu cầu rõ, vì tăng chi phí vận hành và debug.

### Sai lầm phổ biến

- Chọn Kafka chỉ vì “scale lớn”.
- Dùng Kafka như task queue đơn giản nhưng không cần retention/replay.
- Dùng RabbitMQ để lưu event lịch sử rất lâu.
- Tạo quá ít partition rồi không scale consumer được.
- Tạo quá nhiều partition làm tăng metadata, rebalance và chi phí vận hành.
- Nghĩ Kafka tự động cung cấp exactly-once cho mọi side effect.

---

## 6. Delivery guarantee và idempotency

### At-most-once

Message có thể mất nhưng không xử lý lặp.

### At-least-once

Message có thể được giao nhiều lần.

Đây là mô hình phổ biến của RabbitMQ và Kafka consumer khi commit/ACK sau xử lý.

### Exactly-once

Exactly-once end-to-end rất khó, đặc biệt khi có side effect ngoài broker như:

- Ghi PostgreSQL.
- Gọi payment gateway.
- Gửi email.

Kafka có transaction và exactly-once semantics trong phạm vi Kafka hoặc một số pipeline hỗ trợ, nhưng application vẫn cần idempotency cho side effect bên ngoài.

### Idempotent consumer

```sql
CREATE TABLE processed_events (
    consumer_name VARCHAR(100) NOT NULL,
    event_id VARCHAR(100) NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (consumer_name, event_id)
);
```

Trong cùng transaction:

1. Insert event ID.
2. Nếu đã tồn tại, bỏ qua.
3. Thực hiện business update.
4. Commit.
5. ACK RabbitMQ hoặc commit Kafka offset.

---

## 7. Retry, DLQ và poison message

### RabbitMQ

```text
Main Queue
→ Retry Queue 10s
→ Retry Queue 1m
→ Retry Queue 10m
→ DLQ
```

RabbitMQ hỗ trợ retry queue bằng TTL và dead-letter exchange.

### Kafka

Kafka không có delayed queue giống RabbitMQ theo cách mặc định. Một chiến lược phổ biến:

```text
main-topic
→ retry-10s-topic
→ retry-1m-topic
→ retry-10m-topic
→ dead-letter-topic
```

Consumer retry topic cần kiểm soát thời gian, số lần thử và ordering.

### Poison message

Poison message luôn làm consumer lỗi, ví dụ schema sai hoặc dữ liệu không hợp lệ.

Không retry vô hạn. Cần đưa sang DLQ/DLT kèm:

- Event ID.
- Lý do lỗi.
- Stack trace rút gọn.
- Số lần retry.
- Timestamp.
- Consumer name.

---

## 8. Ordering, partition và consumer group

### RabbitMQ

Một queue có thể giao message theo thứ tự, nhưng khi có nhiều consumer:

```text
Message 1 → Worker A xử lý 5 giây
Message 2 → Worker B xử lý 1 giây
```

Message 2 có thể hoàn thành trước message 1.

### Kafka

Ordering được bảo đảm trong partition.

Muốn giữ thứ tự theo account:

```text
key = account_id
```

Mọi event cùng account vào cùng partition.

Trade-off:

- Ordering mạnh hơn làm giảm parallelism cho cùng key.
- Một hot key có thể làm một partition quá tải.

### Rebalance

Khi consumer join hoặc leave group, Kafka phân phối lại partition.

Rebalance có thể gây pause tạm thời. Consumer cần:

- Xử lý message trong thời gian hợp lý.
- Cấu hình poll interval phù hợp.
- Graceful shutdown.
- Commit offset đúng thời điểm.

---

## 9. Outbox pattern

Vấn đề dual write:

```text
1. Commit database thành công
2. Publish message thất bại
```

Outbox ghi business data và event trong cùng transaction:

```sql
BEGIN;

UPDATE payments
SET status = 'completed'
WHERE id = 1001;

INSERT INTO outbox_events(id, event_type, payload)
VALUES (
    'evt_1001',
    'payment.completed',
    '{"payment_id":1001}'
);

COMMIT;
```

Publisher đọc outbox rồi gửi sang RabbitMQ hoặc Kafka.

Các cách publish outbox:

- Polling publisher.
- CDC bằng Debezium rồi đẩy vào Kafka.

Outbox không loại bỏ duplicate. Consumer vẫn phải idempotent.

---

## 10. Schema evolution

Event là contract giữa producer và consumer.

Nguyên tắc:

- Thêm field optional.
- Không đổi nghĩa field cũ.
- Không xóa field khi consumer cũ còn chạy.
- Có `event_version`.
- Consumer bỏ qua field không biết.
- Dùng Avro, Protobuf hoặc JSON Schema khi hệ thống lớn.
- Kafka thường kết hợp Schema Registry.

Ví dụ:

```json
{
  "event_id": "evt_1001",
  "event_type": "payment.completed",
  "event_version": 2,
  "occurred_at": "2026-07-23T10:00:00Z",
  "data": {
    "payment_id": 1001,
    "amount": 500000,
    "currency": "VND"
  }
}
```

---

## 11. Kết nối RabbitMQ từ Python

Cài thư viện:

```bash
pip install aio-pika
```

Biến môi trường:

```env
RABBITMQ_URL=amqp://app:password@rabbitmq:5672/app
```

Kết nối:

```python
import os
import aio_pika


async def connect_rabbitmq() -> aio_pika.RobustConnection:
    return await aio_pika.connect_robust(
        os.environ["RABBITMQ_URL"],
        timeout=5,
    )
```

Publish:

```python
import json


async def publish_event(channel, event: dict) -> None:
    exchange = await channel.declare_exchange(
        "business-events",
        aio_pika.ExchangeType.TOPIC,
        durable=True,
    )

    message = aio_pika.Message(
        body=json.dumps(event).encode(),
        delivery_mode=aio_pika.DeliveryMode.PERSISTENT,
        message_id=event["event_id"],
    )

    await exchange.publish(
        message,
        routing_key=event["event_type"],
    )
```

Consumer nên ACK sau khi transaction thành công.

---

## 12. Kết nối Kafka từ Python

Có thể dùng `aiokafka` cho asyncio.

Cài thư viện:

```bash
pip install aiokafka
```

Biến môi trường:

```env
KAFKA_BOOTSTRAP_SERVERS=kafka:9092
KAFKA_TOPIC=payment-events
KAFKA_CONSUMER_GROUP=notification-service
```

### Producer

```python
import json
import os

from aiokafka import AIOKafkaProducer


async def create_kafka_producer() -> AIOKafkaProducer:
    producer = AIOKafkaProducer(
        bootstrap_servers=os.environ["KAFKA_BOOTSTRAP_SERVERS"],
        acks="all",
        enable_idempotence=True,
        value_serializer=lambda value: json.dumps(value).encode("utf-8"),
        key_serializer=lambda value: value.encode("utf-8"),
    )
    await producer.start()
    return producer


async def publish_payment_event(producer, event: dict) -> None:
    await producer.send_and_wait(
        os.environ["KAFKA_TOPIC"],
        key=str(event["payment_id"]),
        value=event,
    )
```

Key bằng `payment_id` giúp event cùng payment vào cùng partition.

Producer phải được `stop()` khi application shutdown.

### Consumer

```python
import json
import os

from aiokafka import AIOKafkaConsumer


async def create_kafka_consumer() -> AIOKafkaConsumer:
    consumer = AIOKafkaConsumer(
        os.environ["KAFKA_TOPIC"],
        bootstrap_servers=os.environ["KAFKA_BOOTSTRAP_SERVERS"],
        group_id=os.environ["KAFKA_CONSUMER_GROUP"],
        enable_auto_commit=False,
        auto_offset_reset="earliest",
        value_deserializer=lambda value: json.loads(value.decode("utf-8")),
    )
    await consumer.start()
    return consumer
```

Xử lý:

```python
async def consume(consumer) -> None:
    try:
        async for record in consumer:
            event = record.value

            await process_idempotently(event)

            await consumer.commit()
    finally:
        await consumer.stop()
```

Không commit offset trước khi business transaction thành công, nếu không consumer có thể bỏ qua message khi worker crash.

### Kết nối Kafka có authentication

Production thường dùng TLS và SASL.

Ví dụ cấu hình khái niệm:

```python
consumer = AIOKafkaConsumer(
    "payment-events",
    bootstrap_servers="kafka.example.com:9093",
    security_protocol="SASL_SSL",
    sasl_mechanism="SCRAM-SHA-512",
    sasl_plain_username=os.environ["KAFKA_USERNAME"],
    sasl_plain_password=os.environ["KAFKA_PASSWORD"],
)
```

Không commit username/password vào Git.

---

## 13. Docker Compose cho RabbitMQ và Kafka

Ví dụ development:

```yaml
services:
  rabbitmq:
    image: rabbitmq:4-management
    environment:
      RABBITMQ_DEFAULT_USER: app
      RABBITMQ_DEFAULT_PASS: secret
      RABBITMQ_DEFAULT_VHOST: app
    ports:
      - "5672:5672"
      - "15672:15672"

  kafka:
    image: bitnami/kafka:latest
    environment:
      KAFKA_CFG_NODE_ID: 1
      KAFKA_CFG_PROCESS_ROLES: broker,controller
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: 1@kafka:9093
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
    ports:
      - "9092:9092"
```

Container Backend dùng:

```env
RABBITMQ_URL=amqp://app:secret@rabbitmq:5672/app
KAFKA_BOOTSTRAP_SERVERS=kafka:9092
```

Không dùng `localhost` từ container Backend để gọi RabbitMQ hoặc Kafka container khác.

Lưu ý: image và biến cấu hình Kafka có thể thay đổi theo phiên bản; production không nên dùng `latest` và cần cấu hình cluster nhiều broker phù hợp.

---

## 14. Monitoring và vận hành

### RabbitMQ metric

- Queue depth.
- Unacked messages.
- Publish rate.
- Deliver rate.
- Consumer count.
- Redelivery rate.
- DLQ size.
- Memory và disk alarm.

### Kafka metric

- Consumer lag.
- Records in/out.
- Under-replicated partitions.
- Offline partitions.
- Request latency.
- Rebalance rate.
- ISR shrink/expand.
- Disk usage.

Consumer lag là khoảng cách giữa latest offset và offset consumer đã xử lý.

Lag tăng liên tục có thể do:

- Consumer chậm.
- Partition không đủ.
- Database downstream chậm.
- Hot partition.
- Consumer crash/rebalance liên tục.

Log nên có:

- event_id.
- correlation_id.
- topic/queue.
- partition và offset với Kafka.
- routing key với RabbitMQ.
- consumer group.
- attempt.
- processing time.

---

## 15. RPC và REST

RPC cho phép gọi procedure trên service khác.

Ví dụ gRPC:

```proto
service AccountService {
  rpc GetBalance(GetBalanceRequest) returns (GetBalanceResponse);
}
```

RPC và REST thường là giao tiếp đồng bộ, khác với event messaging.

RPC phù hợp:

- Contract mạnh.
- Code generation.
- Internal service call.
- Streaming.

REST phù hợp:

- Resource-oriented API.
- Public API.
- Dễ debug bằng HTTP tooling.

Dù dùng RPC hay REST, cần timeout, retry có giới hạn, circuit breaker và idempotency khi phù hợp.

---

## 16. Câu hỏi phỏng vấn

### Kafka khác RabbitMQ như thế nào?

RabbitMQ là message broker mạnh về task queue, routing, ACK, retry và DLQ. Kafka là distributed log mạnh về event streaming, retention, replay, throughput lớn và nhiều consumer group đọc độc lập. RabbitMQ thường xóa message sau ACK, còn Kafka giữ record theo retention và consumer quản lý offset.

### Khi nào dùng RabbitMQ thay Kafka?

Khi cần background task, routing linh hoạt, retry/DLQ đơn giản, request/reply hoặc message chỉ cần xử lý một lần rồi bỏ.

### Khi nào dùng Kafka thay RabbitMQ?

Khi cần event lưu lâu, replay, nhiều hệ thống đọc độc lập, CDC, analytics hoặc throughput event rất lớn.

### Kafka có bảo đảm thứ tự không?

Có trong phạm vi một partition. Muốn giữ thứ tự theo entity phải dùng cùng partition key, ví dụ `account_id`.

### Vì sao số consumer trong một group không nên lớn hơn số partition?

Mỗi partition chỉ được một consumer trong group đọc tại một thời điểm. Consumer dư sẽ không có partition để xử lý.

### Offset là gì?

Offset là vị trí của record trong partition. Consumer commit offset để đánh dấu tiến độ đọc.

### Consumer xử lý xong nhưng crash trước khi ACK hoặc commit offset thì sao?

Message có thể được xử lý lại. Consumer cần idempotent.

### Kafka exactly-once có nghĩa không cần idempotency?

Không. Exactly-once của Kafka chỉ áp dụng trong phạm vi và cấu hình nhất định. Side effect ngoài Kafka như ghi database hoặc gọi API vẫn cần idempotency hoặc transaction phù hợp.

### Outbox pattern giải quyết gì?

Outbox giải quyết lỗi dual write giữa database và broker bằng cách lưu business change và event trong cùng transaction, rồi publish event sau.

### Có thể dùng Kafka làm queue không?

Có thể, nhưng cần cân nhắc. Nếu chỉ cần task queue đơn giản, delayed retry và DLQ, RabbitMQ thường dễ phù hợp hơn. Kafka đáng giá khi cần retention, replay, consumer group và streaming.

---

## 17. Bài tập thực hành

1. Tạo RabbitMQ queue gửi email có retry và DLQ.
2. Tạo Kafka topic `payment-events` có 3 partition.
3. Viết Python producer gửi key bằng `payment_id`.
4. Chạy hai consumer cùng group và quan sát partition assignment.
5. Chạy consumer khác group và xác nhận nó vẫn nhận toàn bộ event.
6. Reset offset để replay dữ liệu.
7. Mô phỏng consumer crash trước ACK/commit và xử lý duplicate.
8. Thiết kế outbox publish sang Kafka.
9. Thiết kế retry topic và dead-letter topic.
10. Chọn Kafka hay RabbitMQ cho email, analytics, audit log và CRM sync, kèm lý do.

## Checklist

- [ ] Phân biệt sync và async.
- [ ] Phân biệt queue và Pub/Sub.
- [ ] Hiểu RabbitMQ exchange, queue, routing key, ACK và prefetch.
- [ ] Hiểu Kafka topic, partition, offset và consumer group.
- [ ] Giải thích được retention và replay.
- [ ] Hiểu ordering theo partition key.
- [ ] So sánh Kafka với RabbitMQ theo use case.
- [ ] Hiểu at-least-once và idempotent consumer.
- [ ] Thiết kế retry, DLQ/DLT và outbox.
- [ ] Theo dõi queue depth và Kafka consumer lag.
