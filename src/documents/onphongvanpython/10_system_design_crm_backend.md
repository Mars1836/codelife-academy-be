# 10. System Design cho hệ thống CRM Backend

## 1. Kiến trúc tổng quát

```text
Client
  ↓
Nginx / Load Balancer
  ↓
Python API
  ├── PostgreSQL
  ├── Redis
  ├── RabbitMQ
  └── External Services
          ↓
        Worker
```

## 2. Vai trò từng thành phần

### Nginx / Load Balancer

- TLS termination.
- Reverse proxy.
- Load balancing.
- Rate limiting cơ bản.
- Routing.

### Python API

- Routing.
- Authentication.
- Authorization.
- Validation.
- Business logic.
- Transaction.
- Gọi external service.
- Publish event.

### PostgreSQL

Lưu dữ liệu chính cần transaction và consistency:

- User.
- Customer.
- Order.
- Payment.
- Permission.
- Audit log.

### Redis

Có thể dùng cho:

- Cache.
- Session.
- Rate limit.
- Distributed lock.
- Temporary data.
- Celery broker trong một số hệ thống.

Redis không nên mặc định là source of truth cho dữ liệu tài chính.

### RabbitMQ

- Background job.
- Event distribution.
- Retry.
- Buffer tải.
- Tách service.

### Worker

- Gửi email.
- Đồng bộ CRM.
- Export report.
- Gọi API chậm.
- Xử lý message.

## 3. Luồng tạo order

```text
Client gửi POST /orders
→ API authenticate
→ validate request
→ kiểm tra idempotency key
→ BEGIN transaction
→ tạo order
→ tạo outbox event
→ COMMIT
→ trả response
→ outbox worker publish order_created
→ email worker gửi mail
→ CRM worker đồng bộ
→ analytics worker cập nhật dữ liệu
```

## 4. Vì sao dùng outbox?

Đảm bảo không xảy ra trạng thái:

```text
Order đã lưu nhưng event không được gửi
```

Order và outbox event được commit trong cùng transaction.

## 5. Cache

Cache-aside:

```text
API đọc Redis
→ nếu có: trả cache
→ nếu không: đọc DB
→ ghi Redis
→ trả dữ liệu
```

Vấn đề cần xử lý:

- TTL.
- Cache invalidation.
- Cache stampede.
- Dữ liệu cũ.
- Không cache dữ liệu nhạy cảm không cần thiết.

## 6. Scale

### API

Stateless để scale ngang.

### Database

- Index.
- Connection pool.
- Read replica.
- Partitioning nếu cần.
- Query optimization.

### Worker

Scale theo queue depth.

### Redis và RabbitMQ

Triển khai HA tùy yêu cầu.

## 7. Security

- HTTPS.
- Password hash.
- JWT/session an toàn.
- RBAC/ABAC.
- Audit log.
- Encrypt secret.
- Input validation.
- SQL parameterization.
- Rate limiting.
- Principle of least privilege.

## 8. Reliability

- Timeout.
- Retry có backoff.
- Circuit breaker.
- DLQ.
- Idempotency.
- Healthcheck.
- Graceful shutdown.
- Backup.
- Disaster recovery.

## 9. Observability

Mỗi request có request ID.

Theo dõi:

- API latency.
- Error rate.
- Database slow query.
- Queue depth.
- Consumer failure.
- External service latency.
- Pod restart.
- Business metric.

## Bài nói 5 phút

> Em sẽ thiết kế hệ thống theo hướng client đi qua load balancer hoặc Nginx đến các Python API stateless. PostgreSQL là nguồn dữ liệu chính vì CRM và fintech cần transaction và consistency. Redis dùng cho cache, session hoặc rate limiting, không dùng làm nguồn dữ liệu tài chính chính. Những tác vụ chậm như gửi email, đồng bộ CRM và báo cáo được đưa qua RabbitMQ để worker xử lý bất đồng bộ. Với sự kiện quan trọng, em dùng outbox pattern để tránh trường hợp database commit nhưng publish message thất bại. Toàn bộ ứng dụng được đóng gói bằng Docker và triển khai bằng Kubernetes với Deployment, Service, Ingress, readiness probe và rolling update. CI/CD chạy lint, test, build image, deploy staging rồi production. Hệ thống cần timeout, retry có giới hạn, idempotency, logging, metrics và distributed tracing để dễ vận hành.
