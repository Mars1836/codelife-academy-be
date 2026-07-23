# 12. Kế hoạch ôn cấp tốc Python Backend

Kế hoạch này dành cho trường hợp chỉ còn ít ngày trước phỏng vấn. Mục tiêu không phải học hết mọi công nghệ, mà là nắm chắc phần cốt lõi, trả lời có cấu trúc và chứng minh được tư duy xử lý vấn đề.

## 1. Nguyên tắc ưu tiên

Ưu tiên theo thứ tự:

1. Nội dung xuất hiện trực tiếp trong JD.
2. Kiến thức nền Backend bắt buộc.
3. Phần liên quan domain công ty, đặc biệt transaction và consistency nếu là fintech.
4. Công nghệ bạn đã ghi trong CV.
5. System Design và behavioral.
6. Công nghệ phụ chỉ cần hiểu khái niệm và trade-off.

Không dành phần lớn thời gian học một công nghệ mới chỉ vì tên công nghệ đó “hot”.

---

## 2. Phương pháp học mỗi chủ đề

Với mỗi chủ đề, trả lời đủ sáu câu:

```text
1. Nó là gì?
2. Nó giải quyết vấn đề gì?
3. Nó hoạt động thế nào?
4. Khi nào nên dùng?
5. Khi nào không nên dùng?
6. Ví dụ thực tế và lỗi thường gặp?
```

Sau đó luyện một câu trả lời 60–90 giây.

Ví dụ với index:

> Index là cấu trúc dữ liệu phụ giúp database tìm dòng nhanh hơn thay vì quét toàn bảng. B-tree index phù hợp equality, range và ordering theo cột index. Đổi lại index tốn dung lượng và làm thao tác ghi chậm hơn. Với composite index `(user_id, created_at)`, query theo `user_id` hoặc cả hai cột thường dùng tốt, còn chỉ theo `created_at` thường không tối ưu vì nguyên tắc leftmost prefix. Em sẽ xác nhận bằng `EXPLAIN ANALYZE` thay vì chỉ đoán.

---

## 3. Kế hoạch 3 ngày

### Ngày 1: Python và Database

#### Buổi sáng — Python core

Ôn:

- list, tuple, set, dict.
- Mutable và immutable.
- `is` và `==`.
- Shallow copy và deep copy.
- Default mutable argument.
- `*args`, `**kwargs`.
- Decorator.
- Context manager.
- Exception handling.
- OOP, property, classmethod, staticmethod.

Thực hành:

- Viết decorator đo thời gian.
- Viết context manager transaction giả lập.
- Giải thích kết quả các đoạn code mutable/default argument.

#### Buổi chiều — Concurrency và framework

Ôn:

- Thread, process, async.
- GIL.
- I/O-bound và CPU-bound.
- Event loop.
- Asyncio so với Node.js.
- FastAPI/Django request lifecycle.
- Dependency injection.
- Middleware.
- Validation.

Câu trả lời cần chắc:

- Khi nào dùng async?
- Vì sao gọi blocking function trong async làm nghẽn event loop?
- Thread lock hoạt động thế nào?
- GIL có nghĩa Python không chạy đồng thời hay không?

#### Buổi tối — SQL và transaction

Ôn:

- JOIN.
- GROUP BY, WHERE, HAVING.
- CTE và subquery.
- Index và composite index.
- Constraint.
- ACID.
- Isolation level.
- Lock, deadlock, lost update.
- N+1 query.
- `EXPLAIN ANALYZE`.

Thực hành bắt buộc:

```sql
UPDATE accounts
SET balance = balance - :amount
WHERE id = :id
  AND balance >= :amount
RETURNING balance;
```

Giải thích vì sao thao tác atomic này tránh kiểu đọc rồi ghi không an toàn.

---

### Ngày 2: Docker, Linux, Git và CI/CD

#### Buổi sáng — Docker

Ôn:

- Image và container.
- Dockerfile.
- Layer và cache.
- CMD và ENTRYPOINT.
- Port mapping.
- Volume.
- Network và `localhost` trong container.
- Docker Compose.
- Healthcheck.
- Multi-stage build.
- Non-root.

Thực hành:

- Viết Dockerfile FastAPI.
- Viết Compose gồm API và PostgreSQL.
- Debug API không kết nối được database vì dùng `localhost`.

#### Buổi chiều — Linux và Git

Linux:

- Process, signal.
- `systemctl`, `journalctl`.
- `ss`, `curl`, `nc`.
- CPU, RAM, disk.
- Permission.
- Nginx.

Git:

- Working tree, staging, commit.
- Branch và HEAD.
- Merge và rebase.
- Reset và revert.
- Fetch và pull.
- Conflict.
- Force-with-lease.

Thực hành:

- Tạo conflict và giải quyết.
- Dùng reflog cứu commit.
- Điều tra một service không truy cập được theo checklist.

#### Buổi tối — CI/CD

Ôn:

- CI, Continuous Delivery, Continuous Deployment.
- Artifact.
- Build once, promote.
- Secret management.
- Migration trong pipeline.
- Rolling, blue-green và canary.
- Health check và rollback.
- Logs, metrics, traces.

Chuẩn bị nói được luồng:

```text
PR
→ test/lint
→ build image theo SHA
→ push registry
→ deploy staging
→ smoke test
→ approval
→ migration tương thích ngược
→ deploy production
→ monitor/rollback
```

---

### Ngày 3: Kubernetes, Messaging, System Design và Behavioral

#### Buổi sáng — Kubernetes

Ôn:

- Pod.
- Deployment và ReplicaSet.
- Service.
- Ingress.
- ConfigMap và Secret.
- Readiness, liveness, startup probe.
- Request và limit.
- Rolling update.
- Job và StatefulSet.
- Các lỗi Pending, CrashLoopBackOff, ImagePullBackOff.

Không cần giả vờ có kinh nghiệm vận hành production nếu chưa có. Cần đọc manifest và giải thích được request đi qua hệ thống.

#### Buổi trưa — Messaging

Ôn:

- Sync và async.
- Queue và Pub/Sub.
- RabbitMQ exchange/routing/ack.
- Retry và DLQ.
- At-least-once.
- Idempotent consumer.
- Outbox pattern.
- RabbitMQ và Kafka.
- RPC và REST.

Câu quan trọng:

> Consumer xử lý thành công nhưng crash trước ack thì message có thể được giao lại. Vì vậy consumer phải idempotent.

#### Buổi chiều — System Design

Luyện thiết kế CRM hoặc payment system theo thứ tự:

1. Làm rõ requirement.
2. Nêu quy mô giả định.
3. Vẽ kiến trúc tổng quát.
4. Thiết kế API.
5. Thiết kế schema và index.
6. Transaction và consistency.
7. Cache, queue và search.
8. Scaling.
9. Failure handling.
10. Security và observability.
11. Trade-off.

Không bắt đầu ngay bằng microservice.

#### Buổi tối — Behavioral và mock interview

Chuẩn bị:

- Giới thiệu bản thân.
- Dự án tự hào nhất.
- Một lỗi đã mắc.
- Một production incident.
- Một bất đồng kỹ thuật.
- Một deadline gấp.
- Một feedback đã nhận.

Mock interview 45–60 phút, ghi âm và sửa các câu trả lời dài hoặc thiếu ví dụ.

---

## 4. Kế hoạch 7 ngày

### Ngày 1 — Python Core

- Data structure.
- Mutable/immutable.
- Function, decorator, context manager.
- Exception.
- OOP.
- 15 câu code output.

### Ngày 2 — Concurrency và Framework

- Thread/process/async.
- GIL.
- FastAPI/Django lifecycle.
- Auth, validation, middleware.
- Viết API nhỏ có async database.

### Ngày 3 — Database

- SQL.
- Index.
- Transaction.
- Isolation.
- Lock/deadlock.
- Migration.
- `EXPLAIN ANALYZE` trên dữ liệu test.

### Ngày 4 — Docker, Linux và Git

- Dockerfile/Compose.
- Debug network/volume.
- Linux troubleshooting.
- Git merge/rebase/reset/revert.

### Ngày 5 — CI/CD, Cloud và Kubernetes

- Pipeline.
- Artifact và secret.
- Deployment strategy.
- Kubernetes manifest và debug.

### Ngày 6 — Messaging và System Design

- Queue/PubSub.
- Retry, idempotency, outbox.
- Thiết kế CRM/payment system.
- Trình bày 5 phút.

### Ngày 7 — Tổng ôn

- Mock interview kỹ thuật.
- Mock behavioral.
- Xem lại câu trả lời yếu.
- Chuẩn bị câu hỏi cho nhà tuyển dụng.
- Nghỉ đủ và kiểm tra lịch/phương tiện.

---

## 5. Kế hoạch một ngày trước phỏng vấn

Không học thêm chủ đề lớn.

Ôn nhanh:

- 10 câu Python.
- 10 câu Database.
- Luồng CI/CD.
- Kiến trúc System Design 5 phút.
- 6 câu chuyện STAR.
- Thông tin công ty và JD.

Chuẩn bị:

- CV.
- Link GitHub/dự án.
- Máy tính và mạng nếu phỏng vấn online.
- Địa chỉ và thời gian nếu onsite.
- 3–5 câu hỏi cho interviewer.

Ngủ đủ quan trọng hơn cố học thêm đến sáng.

---

## 6. Kế hoạch 60 phút ngay trước phỏng vấn

### 0–15 phút

Đọc lại JD và mapping:

```text
Yêu cầu → kinh nghiệm/ví dụ của tôi
```

### 15–30 phút

Ôn câu trả lời trọng tâm:

- Giới thiệu bản thân.
- Index.
- Transaction và lock.
- Async.
- Docker network.
- CI/CD.

### 30–45 phút

Ôn hai câu chuyện STAR mạnh nhất.

### 45–60 phút

Dừng học, kiểm tra thiết bị, uống nước và giữ đầu óc ổn định.

---

## 7. Bộ câu hỏi tự kiểm tra Python

1. List khác tuple?
2. Set phù hợp khi nào?
3. Mutable default argument gây lỗi gì?
4. `is` khác `==`?
5. Shallow copy khác deep copy?
6. Decorator hoạt động thế nào?
7. Context manager gọi `__enter__` và `__exit__` khi nào?
8. Generator tiết kiệm memory thế nào?
9. GIL là gì?
10. Thread, process, async dùng khi nào?
11. Blocking call ảnh hưởng event loop thế nào?
12. `classmethod` khác `staticmethod`?
13. Exception nên catch ở đâu?
14. Dependency injection có lợi gì?
15. N+1 ở ORM là gì?

---

## 8. Bộ câu hỏi tự kiểm tra Database

1. INNER JOIN khác LEFT JOIN?
2. WHERE khác HAVING?
3. CTE khác subquery?
4. Index hoạt động thế nào?
5. Nhược điểm của index?
6. Composite index theo leftmost prefix là gì?
7. Khi nào optimizer không dùng index?
8. ACID là gì?
9. Isolation level khác nhau thế nào?
10. Lost update là gì?
11. `SELECT FOR UPDATE` làm gì?
12. Atomic update số dư viết thế nào?
13. Deadlock xảy ra ra sao?
14. Idempotency key dùng thế nào?
15. Migration zero-downtime cần gì?

---

## 9. Bộ câu hỏi tự kiểm tra DevOps

1. Image khác container?
2. `localhost` trong container là đâu?
3. Volume khác bind mount?
4. Healthcheck dùng để làm gì?
5. CI khác CD?
6. Artifact là gì?
7. Vì sao build once?
8. Secret nên lưu ở đâu?
9. Rolling khác blue-green/canary?
10. Rollback database khó ở đâu?
11. Readiness khác liveness?
12. Pod khác Deployment?
13. Service tìm Pod thế nào?
14. CrashLoopBackOff debug ra sao?
15. Logs, metrics và traces khác nhau thế nào?

---

## 10. Bộ câu hỏi tự kiểm tra System Design

1. Functional và non-functional requirement là gì?
2. Vì sao cần ước lượng scale?
3. Khi nào modular monolith tốt hơn microservice?
4. PostgreSQL là source of truth nghĩa là gì?
5. Redis dùng ở đâu?
6. Cache invalidation xử lý thế nào?
7. Outbox giải quyết gì?
8. Consumer idempotent thế nào?
9. Cursor pagination tốt hơn offset khi nào?
10. Multi-tenancy tránh lộ dữ liệu thế nào?
11. Import file lớn thiết kế ra sao?
12. Search engine đồng bộ thế nào?
13. Reporting có nên query primary DB?
14. Failure của Redis/broker/search xử lý thế nào?
15. Metric nào chứng minh hệ thống khỏe?

---

## 11. Cách trả lời câu kỹ thuật

Dùng cấu trúc:

```text
Định nghĩa
→ cơ chế
→ ví dụ thực tế
→ trade-off
→ cách kiểm chứng
```

Ví dụ async:

> Asyncio là mô hình concurrency dựa trên event loop, phù hợp I/O-bound với nhiều tác vụ chờ như gọi API hoặc database. Khi một coroutine `await`, event loop có thể chạy tác vụ khác thay vì block thread. Nó không tự tăng tốc CPU-bound; phần CPU nặng nên đưa sang process pool hoặc worker. Khi triển khai em chú ý timeout, giới hạn concurrency và tránh gọi thư viện blocking trực tiếp trong event loop.

---

## 12. Khi bị hỏi sâu hơn kiến thức hiện tại

Không đoán bừa. Trả lời:

```text
Phần em chắc chắn là...
Phần em chưa trực tiếp triển khai là...
Theo hiểu biết hiện tại, em sẽ tiếp cận bằng...
Em sẽ kiểm chứng bằng tài liệu chính thức/test/metric...
```

Sự trung thực cộng với tư duy kiểm chứng tốt hơn một câu trả lời tự tin nhưng sai.

---

## 13. Checklist cuối cùng

### Python

- [ ] Mutable/default argument.
- [ ] Decorator/context manager.
- [ ] OOP.
- [ ] Thread/process/async/GIL.
- [ ] Framework lifecycle.

### Database

- [ ] JOIN/GROUP BY/HAVING.
- [ ] Index/composite index.
- [ ] Transaction/ACID/isolation.
- [ ] Lock/deadlock/lost update.
- [ ] N+1/EXPLAIN/migration.
- [ ] Idempotency.

### Hạ tầng

- [ ] Dockerfile/Compose/network/volume.
- [ ] Linux process/port/log/resource.
- [ ] Git merge/rebase/reset/revert.
- [ ] CI/CD/artifact/secret/rollback.
- [ ] Kubernetes core và debug.

### Kiến trúc

- [ ] Queue/PubSub/retry/DLQ.
- [ ] Idempotent consumer/outbox.
- [ ] System Design 5 phút.
- [ ] Security/observability/failure mode.

### Behavioral

- [ ] Giới thiệu bản thân.
- [ ] Sáu câu chuyện STAR.
- [ ] Ba câu hỏi cho interviewer.
- [ ] Trung thực về phạm vi kinh nghiệm.
