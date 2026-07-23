# Bộ tài liệu ôn phỏng vấn Python Backend

Bộ tài liệu được viết theo hướng **hiểu bản chất, trả lời phỏng vấn và áp dụng thực tế**, không chỉ liệt kê định nghĩa.

## Cách sử dụng

Với mỗi chủ đề, nên học theo thứ tự:

1. Đọc định nghĩa và vấn đề công nghệ giải quyết.
2. Hiểu cơ chế hoạt động.
3. Chạy lại ví dụ.
4. Liên hệ với một tình huống Backend thực tế.
5. Tự trả lời phần câu hỏi phỏng vấn.
6. Làm ít nhất một bài tập thực hành.

## Danh sách tài liệu

1. [`python_backend_interview_guide.md`](./python_backend_interview_guide.md) — Python core, OOP, decorator, context manager, exception, concurrency, async và framework.
2. [`03_database_sql_postgresql_mysql_mongodb.md`](./03_database_sql_postgresql_mysql_mongodb.md) — SQL, index, transaction, isolation, locking, migration, PostgreSQL, MySQL, MongoDB và idempotency fintech.
3. [`04_docker.md`](./04_docker.md) — Image, container, Dockerfile, layer, network, volume, Compose, bảo mật, CI/CD và debug.
4. [`05_kubernetes.md`](./05_kubernetes.md) — Pod, Deployment, Service, Ingress, probe, resource, rollout, autoscaling và troubleshooting.
5. [`06_messaging_amqp_pubsub_rpc.md`](./06_messaging_amqp_pubsub_rpc.md) — Queue, Pub/Sub, RabbitMQ, Kafka, retry, DLQ, idempotent consumer, outbox và RPC.
6. [`07_linux.md`](./07_linux.md) — Process, systemd, permission, port, DNS, firewall, tài nguyên, Nginx, Docker và quy trình điều tra sự cố.
7. [`08_git.md`](./08_git.md) — Commit, branch, HEAD, merge, rebase, reset, revert, conflict, reflog và workflow Pull Request.
8. [`09_devops_cicd_cloud.md`](./09_devops_cicd_cloud.md) — CI/CD, artifact, secret, migration, deployment strategy, observability, cloud, backup và self-hosted runner.
9. [`10_system_design_crm_backend.md`](./10_system_design_crm_backend.md) — Thiết kế CRM, multi-tenancy, authorization, search, cache, import, reporting, scaling và failure handling.
10. [`11_behavioral_interview.md`](./11_behavioral_interview.md) — STAR, giới thiệu bản thân, điểm mạnh/yếu, incident, conflict, ownership và câu hỏi dành cho nhà tuyển dụng.
11. [`12_ke_hoach_on_cap_toc.md`](./12_ke_hoach_on_cap_toc.md) — Kế hoạch ôn 3 ngày, 7 ngày, một ngày trước phỏng vấn và bộ câu hỏi tự kiểm tra.

## Thứ tự học đề xuất

Nếu thời gian ngắn:

```text
Python → Database → Docker/Linux → Git/CI-CD
→ Kubernetes/Messaging → System Design → Behavioral
```

Nếu phỏng vấn công ty fintech, ưu tiên đặc biệt:

- Transaction và ACID.
- Isolation level.
- Lost update, locking và deadlock.
- Idempotency.
- Audit log.
- Outbox và duplicate message.
- Migration database an toàn.
- Security và consistency.

## Cấu trúc chung của tài liệu

Mỗi file cố gắng bao gồm:

- Định nghĩa.
- Vấn đề cần giải quyết.
- Cơ chế hoạt động.
- Ví dụ cơ bản.
- Ứng dụng Backend thực tế.
- Lỗi thường gặp.
- Best practices và trade-off.
- Câu hỏi phỏng vấn kèm định hướng trả lời.
- Bài tập và checklist.
