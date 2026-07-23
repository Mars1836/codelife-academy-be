# 9. DevOps, CI/CD và Cloud

## 1. DevOps là gì?

DevOps là văn hóa và tập hợp thực hành giúp Development và Operations phối hợp để đưa phần mềm ra production nhanh, ổn định và có thể quan sát.

Các yếu tố:

- Automation.
- CI/CD.
- Infrastructure as Code.
- Monitoring.
- Feedback nhanh.
- Ownership.
- Security tích hợp sớm.

## 2. CI

Continuous Integration:

```text
Push code
→ lint
→ unit test
→ security scan
→ build
```

Mục tiêu: phát hiện lỗi sớm và đảm bảo nhánh chính luôn có chất lượng.

## 3. CD

Continuous Delivery: luôn sẵn sàng deploy, production có thể cần approval.

Continuous Deployment: thay đổi đạt điều kiện được tự động deploy production.

## 4. Pipeline đề xuất

```text
Push code
→ lint
→ unit test
→ build Docker image
→ scan image
→ push registry
→ deploy staging
→ migration
→ integration test
→ approval
→ deploy production
→ healthcheck
→ rollback nếu lỗi
```

## 5. Quản lý environment

Tách:

```text
development
staging
production
```

Không commit secret vào Git.

Config nên truyền qua:

- Environment variable.
- Secret manager.
- ConfigMap/Secret.
- CI/CD secret store.

## 6. Chiến lược deploy

### Recreate

Dừng bản cũ rồi chạy bản mới. Có downtime.

### Rolling update

Thay dần instance.

### Blue-green

Có hai môi trường, chuyển traffic từ blue sang green.

### Canary

Đưa một phần nhỏ traffic sang bản mới rồi tăng dần.

## 7. Rollback

Rollback cần chuẩn bị trước:

- Image có version immutable.
- Manifest có version.
- Migration backward-compatible.
- Healthcheck rõ.
- Monitoring sau deploy.
- Có lệnh rollback đã kiểm thử.

## 8. VM

Virtual Machine là máy ảo có hệ điều hành riêng.

Ví dụ cloud:

- AWS EC2.
- Azure Virtual Machines.
- Google Compute Engine.

## 9. VPC và network

VPC là mạng riêng logic trên cloud.

Thường gồm:

- Public subnet.
- Private subnet.
- Route table.
- Internet gateway.
- NAT gateway.
- Security group/firewall.

Database thường đặt private subnet.

## 10. Security group và firewall

Quy tắc tối thiểu:

- Chỉ mở port cần thiết.
- PostgreSQL không mở công khai nếu không cần.
- Giới hạn source IP.
- Dùng principle of least privilege.

## 11. Load balancer

Phân phối traffic đến nhiều backend instance.

Vai trò:

- Healthcheck.
- TLS termination.
- High availability.
- Routing.
- Load distribution.

## 12. Object storage

Ví dụ:

- Amazon S3.
- Azure Blob Storage.
- Google Cloud Storage.

Phù hợp:

- File upload.
- Backup.
- Static asset.
- Log archive.

Không phải filesystem truyền thống.

## 13. Managed database

Ví dụ:

- Amazon RDS.
- Azure Database.
- Cloud SQL.

Cloud provider quản lý một phần:

- Backup.
- Patch.
- HA.
- Monitoring.

Nhưng application vẫn phải quản lý schema, query, index và connection.

## 14. Container registry

Lưu Docker image:

- Amazon ECR.
- Azure Container Registry.
- Google Artifact Registry.
- GitHub Container Registry.

Nên tag bằng version hoặc commit SHA, không chỉ dùng `latest`.

## 15. IAM

Identity and Access Management.

Nguyên tắc:

- Least privilege.
- Không dùng tài khoản admin cho ứng dụng.
- Ưu tiên role/service account.
- Rotate credential.
- Audit access.

## 16. Monitoring và logging

Cần theo dõi:

- Request rate.
- Error rate.
- Latency.
- CPU.
- Memory.
- Disk.
- Database connection.
- Queue depth.
- Pod restart.

Ba trụ cột observability:

- Logs.
- Metrics.
- Traces.

## 17. Secret management

Giải pháp:

- Cloud secret manager.
- HashiCorp Vault.
- Kubernetes Secret kết hợp encryption.
- CI/CD secret store.

Không:

- Commit `.env`.
- Ghi secret vào Docker image.
- In secret ra log.
- Dùng chung credential giữa môi trường.

## Câu hỏi phỏng vấn

### Pipeline tốt cần gì?

Fast feedback, test tự động, image immutable, quản lý secret, deploy staging, approval phù hợp, healthcheck, rollback và observability.

### Migration chạy ở đâu?

Có thể chạy bằng job riêng trong pipeline hoặc Kubernetes Job. Cần tránh nhiều replica cùng chạy migration.

### Làm sao giảm rủi ro production?

- Canary/blue-green.
- Feature flag.
- Backward-compatible migration.
- Monitoring.
- Automated rollback.
- Backup.
- Approval cho thay đổi rủi ro.
