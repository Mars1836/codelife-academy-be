# 9. DevOps, CI/CD và Cloud

DevOps là văn hóa và tập hợp thực hành giúp Development và Operations phối hợp để đưa phần mềm ra production nhanh, ổn định, có thể quan sát và có khả năng khôi phục khi xảy ra sự cố.

## 1. DevOps giải quyết vấn đề gì?

Nếu development chỉ quan tâm “code chạy trên máy em”, còn operations chỉ nhận artifact cuối cùng để triển khai thủ công, hệ thống dễ gặp:

- Môi trường không đồng nhất.
- Deploy phụ thuộc cá nhân.
- Khó rollback.
- Không biết thay đổi nào gây lỗi.
- Thiếu monitoring và ownership.

DevOps hướng tới:

```text
Plan → Code → Build → Test → Release → Deploy → Operate → Observe → Feedback
```

DevOps không chỉ là Docker, Kubernetes hay CI/CD. Công cụ chỉ hỗ trợ quy trình và văn hóa chịu trách nhiệm xuyên suốt vòng đời phần mềm.

---

## 2. CI là gì?

Continuous Integration là việc merge thay đổi thường xuyên vào nhánh chung và tự động kiểm tra chất lượng.

Pipeline CI cho Python Backend:

```text
Checkout
→ Install dependency
→ Lint
→ Type check
→ Unit test
→ Integration test
→ Security/dependency scan
→ Build Docker image
```

Ví dụ GitHub Actions:

```yaml
name: CI

on:
  pull_request:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: test_db
        ports:
          - 5432:5432
        options: >-
          --health-cmd "pg_isready -U test"
          --health-interval 5s
          --health-timeout 5s
          --health-retries 10

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with:
          python-version: "3.12"
      - run: pip install -r requirements.txt
      - run: ruff check .
      - run: pytest
```

CI tốt phải nhanh đủ để developer sử dụng thường xuyên, đồng thời đáng tin cậy. Test flaky làm team mất niềm tin vào pipeline.

---

## 3. Continuous Delivery và Continuous Deployment

### Continuous Delivery

Mọi thay đổi đạt chuẩn đều sẵn sàng deploy, nhưng production có thể cần approval thủ công.

### Continuous Deployment

Mọi thay đổi vượt qua pipeline được tự động deploy production.

Không phải hệ thống nào cũng nên tự động deploy production ngay. Fintech hoặc hệ thống có rủi ro cao có thể cần approval, change window hoặc kiểm tra bổ sung.

---

## 4. Artifact là gì?

Artifact là đầu ra bất biến của quá trình build, ví dụ:

- Docker image.
- Python wheel.
- Binary.
- Frontend bundle.
- Migration package.

Nguyên tắc quan trọng:

```text
Build once, promote the same artifact
```

Không build lại source riêng cho staging và production vì có thể tạo artifact khác nhau dù cùng commit.

Ví dụ:

```text
ghcr.io/company/api:sha-abc123
```

Staging và production dùng cùng image digest, chỉ khác configuration.

---

## 5. Pipeline triển khai hoàn chỉnh

```text
Pull Request
→ lint + test + scan
→ merge main
→ build image
→ push registry
→ deploy staging
→ smoke test
→ approval
→ backup/check migration
→ deploy production
→ health check
→ monitor
→ rollback nếu cần
```

Mỗi bước cần có input, output và điều kiện thất bại rõ ràng.

---

## 6. Quản lý environment

Các môi trường thường gặp:

- Local development.
- Test/CI.
- Staging.
- Production.

Không nên tạo code branch riêng cho từng môi trường. Source code giống nhau, configuration khác nhau.

Ví dụ environment variable:

```text
APP_ENV=production
DATABASE_URL=...
REDIS_URL=...
LOG_LEVEL=info
```

Configuration không nhạy cảm có thể lưu trong repo. Secret không nên commit vào Git.

---

## 7. Quản lý secret

Các lựa chọn:

- GitHub Actions Secrets.
- Environment secrets có approval.
- HashiCorp Vault.
- AWS Secrets Manager.
- GCP Secret Manager.
- Azure Key Vault.
- Kubernetes External Secrets.

Nguyên tắc:

- Least privilege.
- Rotate định kỳ hoặc khi nghi ngờ lộ.
- Không in secret vào log.
- Không truyền secret qua command line nếu dễ xuất hiện trong process list.
- Tách secret theo môi trường.
- Audit ai đã truy cập hoặc thay đổi.

Với hệ thống nhỏ, GitHub Secrets kết hợp secret file trên VPS có permission chặt có thể đủ. Vault hữu ích khi số lượng secret, service và yêu cầu dynamic credential tăng, nhưng nó cũng tạo thêm chi phí vận hành.

---

## 8. Docker image trong CI/CD

Dockerfile production nên:

- Dùng base image cụ thể.
- Multi-stage build khi phù hợp.
- Chạy non-root.
- Không chứa secret.
- Có `.dockerignore`.
- Pin dependency.
- Scan vulnerability.

Tag image:

```text
api:sha-abc123
api:v1.4.0
```

Không dựa duy nhất vào `latest`, vì khó biết chính xác phiên bản nào đang chạy và rollback.

---

## 9. Database migration trong pipeline

Migration là phần rủi ro cao vì database là shared state.

Chiến lược an toàn:

1. Backup hoặc xác nhận recovery plan.
2. Chạy migration backward-compatible.
3. Deploy code mới.
4. Backfill theo batch nếu cần.
5. Theo dõi lock, latency và error.
6. Xóa schema cũ ở release sau.

Không chạy migration phá tương thích trước khi toàn bộ instance code cũ dừng.

Ví dụ expand-contract:

```text
Release 1: thêm cột mới
Release 2: ghi cả cột cũ và mới
Release 3: đọc cột mới
Release 4: xóa cột cũ
```

Migration nên là Job riêng hoặc bước được kiểm soát, không để mọi Pod cùng chạy migration khi startup.

---

## 10. Deployment strategies

### Recreate

Dừng bản cũ rồi chạy bản mới. Đơn giản nhưng có downtime.

### Rolling update

Thay instance dần dần. Cần code và migration tương thích giữa version cũ và mới.

### Blue-green

```text
Blue = production hiện tại
Green = version mới
```

Deploy Green, test rồi chuyển traffic. Rollback nhanh bằng cách chuyển lại Blue, nhưng tốn tài nguyên gấp đôi trong thời gian chuyển đổi.

### Canary

Chuyển một phần nhỏ traffic sang version mới:

```text
95% version cũ
5% version mới
```

Theo dõi error rate, latency và business metric rồi tăng dần.

### Feature flag

Deploy code nhưng bật/tắt tính năng độc lập. Cần quản lý lifecycle; flag cũ không được dọn sẽ tạo technical debt.

---

## 11. Rollback

Rollback application thường là deploy lại image trước.

```text
api:sha-good
```

Rollback database khó hơn. Migration `DROP COLUMN` hoặc biến đổi dữ liệu có thể không khôi phục đơn giản.

Vì vậy chiến lược tốt thường là roll forward bằng migration sửa lỗi, kết hợp backup và backward compatibility.

Pipeline cần biết:

- Version đang chạy.
- Version trước.
- Image digest.
- Migration version.
- Người/commit kích hoạt deploy.

---

## 12. Health check và smoke test

Sau deploy không chỉ kiểm tra process còn chạy.

Smoke test có thể kiểm tra:

- `/health/live`.
- `/health/ready`.
- Login hoặc endpoint quan trọng.
- Kết nối database.
- Publish/consume một message test nếu phù hợp.

Không chạy test phá dữ liệu production. Dùng synthetic account hoặc endpoint kiểm tra riêng.

---

## 13. Observability

Ba trụ cột thường nhắc đến:

- Logs.
- Metrics.
- Traces.

### Logs

Log có cấu trúc:

```json
{
  "level": "error",
  "message": "payment failed",
  "request_id": "req_123",
  "user_id": 10,
  "error_code": "BANK_TIMEOUT"
}
```

Không log password, access token, số thẻ hoặc dữ liệu nhạy cảm không cần thiết.

### Metrics

Các metric backend:

- Request rate.
- Error rate.
- Latency p50/p95/p99.
- CPU/RAM.
- Database connection pool.
- Queue depth.
- Business metric như payment success rate.

### Tracing

Distributed trace theo request qua nhiều service, database và broker. Cần propagate trace ID/correlation ID.

---

## 14. SLI, SLO và SLA

- SLI: chỉ số đo, ví dụ tỷ lệ request thành công.
- SLO: mục tiêu nội bộ, ví dụ 99,9% request thành công trong 30 ngày.
- SLA: cam kết với khách hàng, có thể kèm bồi thường.

Error budget là phần lỗi được phép trong SLO. Khi tiêu thụ error budget quá nhanh, team nên ưu tiên độ ổn định hơn phát hành tính năng mới.

---

## 15. Infrastructure as Code

IaC quản lý hạ tầng bằng code, ví dụ Terraform.

```hcl
resource "aws_s3_bucket" "uploads" {
  bucket = "company-uploads-production"
}
```

Lợi ích:

- Review thay đổi hạ tầng qua PR.
- Tái tạo môi trường.
- Giảm thao tác thủ công.
- Có lịch sử và plan.

Cần bảo vệ state file vì nó có thể chứa thông tin nhạy cảm và là nguồn trạng thái quan trọng.

---

## 16. Cloud cơ bản

### Compute

- VM: kiểm soát cao, tự vận hành OS.
- Managed container service: giảm quản lý cluster.
- Kubernetes: linh hoạt nhưng phức tạp.
- Serverless: trả theo invocation, scale tự động, có giới hạn runtime/cold start.

### Storage

- Block storage: disk cho VM/database.
- Object storage: file, backup, ảnh, artifact.
- File storage: filesystem chia sẻ.

### Managed database

Managed PostgreSQL giảm công sức backup, replication, patching nhưng không loại bỏ trách nhiệm thiết kế schema, query, capacity và recovery test.

---

## 17. Backup và disaster recovery

Backup chỉ có giá trị khi restore được.

Cần xác định:

- RPO: chấp nhận mất tối đa bao nhiêu dữ liệu.
- RTO: cần khôi phục trong bao lâu.

Ví dụ:

```text
RPO = 5 phút
RTO = 30 phút
```

Thực hành:

- Backup tự động.
- Lưu ở vị trí/tài khoản khác khi phù hợp.
- Mã hóa.
- Retention policy.
- Restore test định kỳ.
- Runbook rõ ràng.

---

## 18. Self-hosted runner

GitHub Actions self-hosted runner là agent chạy trên máy của bạn. GitHub gửi job đã được phân công đến runner, runner checkout code và thực thi command trên máy đó.

Rủi ro:

- Workflow có thể truy cập network và credential của runner.
- Workspace hoặc Docker cache có thể còn dữ liệu.
- Pull request không tin cậy có thể chạy code độc hại.

Best practices:

- Runner riêng theo môi trường.
- Không chạy PR từ fork không tin cậy trên runner production.
- Least privilege.
- Ephemeral runner khi có thể.
- Không dùng cùng máy production cho build không kiểm soát.

---

## 19. Ví dụ pipeline deploy VPS

```text
PR: lint + test
Merge main: build image + push registry
Deploy job:
  SSH vào VPS
  pull image theo SHA
  chạy migration có kiểm soát
  docker compose up -d
  health check
  rollback image nếu health check thất bại
```

Không copy `.env` từ repo. VPS giữ environment file hoặc lấy secret từ secret manager.

Lệnh deploy nên idempotent: chạy lại không làm hệ thống hỏng.

---

## 20. Câu hỏi phỏng vấn

### CI khác CD thế nào?

CI tập trung tích hợp thường xuyên và tự động kiểm tra code. Continuous Delivery đảm bảo artifact luôn sẵn sàng phát hành; Continuous Deployment tự động đưa mọi thay đổi đạt chuẩn lên production.

### Artifact là gì?

Artifact là đầu ra bất biến của build như Docker image hoặc wheel. Nên build một lần rồi promote cùng artifact qua các môi trường.

### Vì sao không dùng `latest`?

`latest` là tag có thể thay đổi, khó audit và rollback. Nên dùng tag theo commit SHA/version và tốt hơn nữa là image digest.

### Blue-green khác canary?

Blue-green chuẩn bị toàn bộ môi trường mới rồi chuyển traffic gần như một lần. Canary chuyển một phần traffic nhỏ trước và tăng dần dựa trên metric.

### Làm sao deploy migration không downtime?

Dùng expand-contract, migration backward-compatible, backfill theo batch và đảm bảo code cũ/mới cùng hoạt động trong thời gian rolling update.

### Khi nào cần Vault?

Khi hệ thống có nhiều service, nhiều môi trường, cần dynamic secret, rotation, audit và kiểm soát truy cập tập trung. Hệ thống nhỏ có thể bắt đầu với secret store đơn giản hơn để tránh tăng độ phức tạp không cần thiết.

---

## 21. Bài tập thực hành

1. Viết GitHub Actions chạy Ruff và pytest.
2. Build image có tag commit SHA và push registry.
3. Thiết kế pipeline staging → approval → production.
4. Mô phỏng health check thất bại và rollback.
5. Viết migration theo expand-contract.
6. Thiết kế secret management cho dev, staging và production.
7. Định nghĩa SLI/SLO cho API thanh toán.
8. Viết runbook khôi phục PostgreSQL từ backup.

## Checklist

- [ ] Phân biệt CI, Continuous Delivery và Continuous Deployment.
- [ ] Hiểu artifact bất biến và build once.
- [ ] Quản lý environment/secret an toàn.
- [ ] Hiểu rolling, blue-green, canary và feature flag.
- [ ] Biết migration backward-compatible và rollback.
- [ ] Hiểu logs, metrics, traces, SLI và SLO.
- [ ] Hiểu RPO, RTO và restore test.
- [ ] Nhận biết rủi ro của self-hosted runner.
