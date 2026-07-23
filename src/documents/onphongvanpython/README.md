# Bộ tài liệu ôn phỏng vấn Python Backend

Bộ tài liệu được viết theo hướng **hiểu bản chất, trả lời phỏng vấn và áp dụng thực tế**, không chỉ liệt kê định nghĩa.

## Cách sử dụng

Với mỗi chủ đề:

1. Đọc định nghĩa và vấn đề công nghệ giải quyết.
2. Hiểu cơ chế hoạt động.
3. Chạy lại ví dụ.
4. Liên hệ tình huống Backend thực tế.
5. Tự trả lời câu hỏi phỏng vấn.
6. Làm bài tập và kiểm tra checklist.

## Danh sách tài liệu

1. [`python_backend_interview_guide.md`](./python_backend_interview_guide.md) — Python core, OOP, concurrency, async và framework.
2. [`03_database_sql_postgresql_mysql_mongodb.md`](./03_database_sql_postgresql_mysql_mongodb.md) — SQL, index, transaction, locking, migration và idempotency.
3. [`04_docker.md`](./04_docker.md) — Dockerfile, network, volume, Compose, bảo mật và debug.
4. [`05_kubernetes.md`](./05_kubernetes.md) — Pod, Deployment, Service, Ingress, probe, resource và troubleshooting.
5. [`06_messaging_amqp_pubsub_rpc.md`](./06_messaging_amqp_pubsub_rpc.md) — Queue, Pub/Sub, RabbitMQ, Kafka, retry, DLQ, outbox và RPC.
6. [`07_linux.md`](./07_linux.md) — Process, systemd, port, DNS, firewall, Nginx và điều tra sự cố.
7. [`08_git.md`](./08_git.md) — Commit, branch, merge, rebase, reset, revert, conflict và reflog.
8. [`09_devops_cicd_cloud.md`](./09_devops_cicd_cloud.md) — CI/CD, artifact, secret, migration, deployment, cloud và observability.
9. [`10_system_design_crm_backend.md`](./10_system_design_crm_backend.md) — CRM, multi-tenancy, authorization, search, cache, import và scaling.
10. [`11_behavioral_interview.md`](./11_behavioral_interview.md) — STAR, incident, conflict, ownership và câu hỏi behavioral.
11. [`12_ke_hoach_on_cap_toc.md`](./12_ke_hoach_on_cap_toc.md) — Kế hoạch ôn 3 ngày, 7 ngày và bộ câu hỏi tự kiểm tra.

## Thứ tự học đề xuất

```text
Python → Database → Docker/Linux → Git/CI-CD
→ Kubernetes/Messaging → System Design → Behavioral
```

Nếu phỏng vấn fintech, ưu tiên transaction, isolation, locking, deadlock, idempotency, audit, outbox, migration an toàn, security và consistency.

Mỗi file gồm định nghĩa, cơ chế, ví dụ, ứng dụng thực tế, lỗi thường gặp, best practice, trade-off, câu hỏi phỏng vấn, bài tập và checklist.
