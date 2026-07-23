# 10. System Design cho hệ thống CRM Backend

Tài liệu này trình bày cách phân tích và thiết kế một hệ thống CRM theo cách có thể dùng trong phỏng vấn Backend/System Design. Mục tiêu không phải tạo “kiến trúc hoàn hảo”, mà là đưa ra thiết kế hợp lý, nêu rõ giả định, trade-off, rủi ro và hướng mở rộng.

## 1. Bài toán

CRM quản lý:

- Khách hàng và liên hệ.
- Công ty/tổ chức.
- Lead và cơ hội bán hàng.
- Pipeline và stage.
- Hoạt động: cuộc gọi, email, cuộc hẹn, ghi chú.
- Task và nhắc việc.
- Phân quyền theo đội nhóm.
- Import/export dữ liệu.
- Báo cáo và dashboard.
- Đồng bộ hệ thống ngoài.

Không nên bắt đầu bằng việc vẽ microservice. Trước tiên cần làm rõ yêu cầu và quy mô.

---

## 2. Functional requirements

Yêu cầu cốt lõi:

1. Tạo, sửa, tìm kiếm customer/contact.
2. Tạo lead và chuyển lead thành opportunity.
3. Di chuyển opportunity qua các stage.
4. Ghi nhận activity và task.
5. Phân quyền theo organization, team và owner.
6. Lọc, sắp xếp, phân trang dữ liệu.
7. Import CSV số lượng lớn.
8. Gửi notification và email.
9. Xem báo cáo doanh số.
10. Audit lịch sử thay đổi.

Có thể chia MVP và phần mở rộng. Trong phỏng vấn nên ưu tiên luồng chính trước.

---

## 3. Non-functional requirements

- Availability cao cho CRUD và tìm kiếm.
- Dữ liệu tenant không bị lẫn.
- Transaction đúng cho thay đổi nghiệp vụ quan trọng.
- Audit được ai sửa gì, lúc nào.
- API p95 dưới mục tiêu, ví dụ 300 ms cho CRUD thông thường.
- Import lớn không làm nghẽn API.
- Có thể scale ngang API và worker.
- Backup, restore và disaster recovery.
- Bảo vệ dữ liệu cá nhân.

Trade-off: báo cáo có thể chậm vài phút so với dữ liệu giao dịch nếu dùng pipeline analytics riêng.

---

## 4. Ước lượng quy mô

Giả định:

- 10.000 organization.
- Trung bình 50 user/organization.
- 100 triệu contact.
- 20 triệu opportunity.
- 1 tỷ activity sau vài năm.
- Peak 5.000 request/giây.
- Import tối đa 1 triệu dòng/file.

Ước lượng giúp quyết định index, partition, cache, queue và search engine. Nếu quy mô chỉ vài nghìn record, kiến trúc đơn giản hơn sẽ tốt hơn.

---

## 5. Kiến trúc tổng quát

```text
Web / Mobile
    ↓
CDN / WAF / Load Balancer
    ↓
API Gateway hoặc Nginx
    ↓
Python API stateless
    ├── PostgreSQL
    ├── Redis
    ├── Object Storage
    ├── Search Engine
    └── Message Broker
             ↓
          Workers
             ├── Email
             ├── Import
             ├── CRM Sync
             └── Report jobs

Logs / Metrics / Traces → Observability platform
Analytics events → Warehouse / Superset
```

Có thể bắt đầu bằng modular monolith thay vì microservice:

```text
CRM Application
  ├── identity
  ├── contacts
  ├── opportunities
  ├── activities
  ├── notifications
  └── reporting
```

Modular monolith giảm network call, transaction phân tán và chi phí vận hành. Tách service khi có lý do rõ: scale độc lập, ownership khác, boundary ổn định hoặc yêu cầu bảo mật riêng.

---

## 6. API design

Ví dụ REST:

```text
POST   /v1/contacts
GET    /v1/contacts/{id}
GET    /v1/contacts?owner_id=10&status=active&cursor=...
PATCH  /v1/contacts/{id}
POST   /v1/opportunities
POST   /v1/opportunities/{id}/move-stage
POST   /v1/imports
GET    /v1/imports/{id}
```

Tạo contact:

```json
{
  "first_name": "Công Hậu",
  "last_name": "Vũ",
  "email": "hau@example.com",
  "phone": "+84901234567",
  "owner_id": 100
}
```

Response không nên lộ internal field hoặc dữ liệu tenant khác.

Dùng versioning khi public API cần ổn định. Internal API vẫn cần contract và compatibility.

---

## 7. Pagination

Offset pagination:

```sql
SELECT *
FROM contacts
WHERE organization_id = $1
ORDER BY created_at DESC
LIMIT 50 OFFSET 10000;
```

Đơn giản nhưng offset lớn chậm và dữ liệu có thể nhảy khi có insert mới.

Cursor pagination:

```sql
SELECT *
FROM contacts
WHERE organization_id = $1
  AND (created_at, id) < ($2, $3)
ORDER BY created_at DESC, id DESC
LIMIT 50;
```

Phù hợp danh sách lớn và feed thay đổi liên tục. Cursor cần dựa trên sort key ổn định, thường thêm `id` để xử lý timestamp trùng.

---

## 8. Database schema

```sql
CREATE TABLE organizations (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    email TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, email)
);

CREATE TABLE contacts (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    owner_id BIGINT REFERENCES users(id),
    first_name TEXT,
    last_name TEXT,
    email TEXT,
    phone TEXT,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_contacts_org_created
ON contacts(organization_id, created_at DESC, id DESC)
WHERE deleted_at IS NULL;
```

Mọi bảng tenant-scoped nên có `organization_id` và query luôn giới hạn tenant.

---

## 9. Multi-tenancy

Ba mô hình phổ biến:

### Shared database, shared schema

Mọi tenant dùng chung bảng và có `organization_id`.

Ưu điểm:

- Chi phí thấp.
- Dễ vận hành.
- Dễ query cross-tenant cho admin nội bộ.

Nhược điểm:

- Rủi ro lộ dữ liệu nếu quên tenant filter.
- Tenant lớn có thể ảnh hưởng tenant nhỏ.

### Shared database, schema riêng

Mỗi tenant có schema riêng. Isolation tốt hơn nhưng migration nhiều schema phức tạp.

### Database riêng mỗi tenant

Isolation và custom cao, nhưng tốn chi phí và khó vận hành khi tenant nhiều.

MVP thường dùng shared schema, kết hợp:

- Tenant context bắt buộc.
- Repository tự thêm organization filter.
- Row-Level Security nếu phù hợp.
- Test chống cross-tenant access.

Ví dụ PostgreSQL RLS:

```sql
ALTER TABLE contacts ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON contacts
USING (organization_id = current_setting('app.organization_id')::BIGINT);
```

RLS là lớp phòng vệ bổ sung, không thay thế thiết kế authentication và connection handling đúng.

---

## 10. Authorization

### RBAC

Role như admin, manager, sales, viewer.

### Resource ownership

Sales chỉ sửa contact/opportunity do mình sở hữu.

### ABAC

Quyết định dựa trên thuộc tính user, resource và context:

```text
user.organization_id == contact.organization_id
AND (
  user.role == admin
  OR contact.owner_id == user.id
  OR user.team_id == contact.team_id
)
```

Authorization phải kiểm tra ở backend. Ẩn button trên frontend không phải security.

---

## 11. Transaction và optimistic locking

Hai user có thể cùng sửa contact. Nếu update cuối cùng ghi đè update trước, dữ liệu bị lost update.

```sql
UPDATE contacts
SET phone = $1,
    version = version + 1,
    updated_at = NOW()
WHERE id = $2
  AND organization_id = $3
  AND version = $4;
```

Nếu affected rows bằng 0, trả conflict để client tải bản mới.

Di chuyển opportunity stage và ghi audit nên nằm trong cùng transaction:

```text
BEGIN
UPDATE opportunity
INSERT stage_history
INSERT audit_log
INSERT outbox_event
COMMIT
```

---

## 12. Audit log

Audit log cần trả lời:

- Ai thay đổi?
- Thay đổi resource nào?
- Trước và sau là gì?
- Khi nào?
- Từ request/IP nào?

Schema:

```sql
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    actor_id BIGINT,
    action TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT NOT NULL,
    before_data JSONB,
    after_data JSONB,
    request_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Không cho user bình thường sửa audit log. Dữ liệu nhạy cảm cần redact trước khi lưu.

---

## 13. Search

PostgreSQL có thể đáp ứng giai đoạn đầu bằng:

- B-tree cho exact/filter.
- Trigram index cho fuzzy text.
- Full-text search.

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_contacts_name_trgm
ON contacts USING gin ((first_name || ' ' || last_name) gin_trgm_ops);
```

Khi cần search phức tạp, ranking, typo tolerance và aggregation lớn, có thể dùng OpenSearch/Elasticsearch.

Luồng đồng bộ:

```text
PostgreSQL commit
→ Outbox
→ Search indexer
→ Search Engine
```

PostgreSQL là source of truth; search index có thể eventually consistent và rebuild được.

---

## 14. Cache

Redis dùng cho:

- Session/token metadata nếu kiến trúc cần.
- Rate limiting.
- Cache permission hoặc reference data.
- Distributed lock trong use case phù hợp.
- Job metadata.

Cache-aside:

```text
Read cache
→ miss
Read database
→ set cache TTL
→ return
```

Rủi ro:

- Stale data.
- Cache stampede.
- Invalidation phức tạp.

Không cache mọi thứ. Query PostgreSQL có index đúng có thể đủ nhanh và đơn giản hơn.

Không dùng Redis làm source of truth cho dữ liệu CRM quan trọng nếu không có thiết kế durability rõ ràng.

---

## 15. Import CSV lớn

Không xử lý file triệu dòng trong HTTP request.

Luồng:

```text
Client upload trực tiếp object storage bằng signed URL
→ API tạo import job
→ Queue
→ Worker đọc file theo stream/batch
→ Validate
→ Batch insert/upsert
→ Ghi lỗi theo dòng
→ Cập nhật progress
→ Notification hoàn thành
```

Import cần:

- Giới hạn kích thước.
- Virus/file type validation.
- Idempotency.
- Partial failure strategy.
- Error report.
- Rate limit theo tenant.
- Backpressure để không làm nghẽn database.

Không load toàn bộ CSV vào RAM.

---

## 16. Email và notification

API không chờ gửi email nếu email không quyết định kết quả request.

```text
Opportunity assigned
→ outbox event
→ broker
→ notification worker
→ email provider
```

Worker cần retry có backoff, idempotency và DLQ.

Notification trong app có thể lưu database:

```sql
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    type TEXT NOT NULL,
    payload JSONB NOT NULL,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 17. Reporting và Superset

Không nên để dashboard nặng query trực tiếp vào primary database nếu ảnh hưởng transactional workload.

Các mức:

1. Query primary cho báo cáo nhỏ.
2. Read replica cho truy vấn đọc.
3. Materialized view/pre-aggregation.
4. ETL/CDC sang data warehouse.
5. Superset query warehouse hoặc database reporting.

```text
CRM PostgreSQL
→ CDC/ETL
→ Warehouse
→ Superset
```

Superset là lớp BI/query/visualization, không thay thế backend nghiệp vụ và không nên dùng để cập nhật transaction CRM.

---

## 18. Messaging và outbox

Các event:

- `contact.created`.
- `opportunity.stage_changed`.
- `task.overdue`.
- `import.completed`.

Outbox bảo đảm business transaction và event được ghi cùng nhau. Consumer cần idempotent vì broker thường giao at-least-once.

Message nên chứa ID và dữ liệu cần thiết, nhưng tránh nhét toàn bộ record nhạy cảm nếu không cần.

---

## 19. File storage

Attachment lưu object storage, database chỉ lưu metadata:

```sql
CREATE TABLE attachments (
    id UUID PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    object_key TEXT NOT NULL,
    file_name TEXT NOT NULL,
    content_type TEXT,
    size_bytes BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Upload bằng signed URL giảm tải API. Download cần authorization trước khi cấp signed URL ngắn hạn.

---

## 20. Scale API

Python API nên stateless để scale ngang:

```text
Load Balancer
 ├→ API 1
 ├→ API 2
 └→ API 3
```

State dùng chung đặt ở PostgreSQL, Redis hoặc object storage.

Các giới hạn:

- Database connection pool: nhiều replica có thể tạo quá nhiều connection.
- Dependency downstream.
- Hot tenant.
- Query không có index.

Có thể dùng PgBouncer và đặt pool size theo capacity database.

---

## 21. Partitioning

Bảng activity/audit có thể tăng rất lớn. Partition theo thời gian hoặc organization tùy query pattern.

Ví dụ partition theo tháng:

```text
activities_2026_07
activities_2026_08
```

Partition không thay thế index. Chỉ dùng khi bảng đủ lớn và có lợi rõ ràng cho pruning, retention hoặc maintenance.

---

## 22. Availability và failure handling

### PostgreSQL lỗi

- API trả lỗi có kiểm soát.
- Không retry write mù quáng.
- Circuit breaker/connection timeout.
- Failover managed database nếu có.

### Redis lỗi

Nếu Redis chỉ là cache, hệ thống nên fallback database có giới hạn để tránh stampede.

### Broker lỗi

Outbox giữ event chưa publish. API nghiệp vụ vẫn có thể commit nếu thiết kế chấp nhận notification trễ.

### Search lỗi

Fallback tìm kiếm đơn giản từ PostgreSQL hoặc trả thông báo degraded mode.

---

## 23. Observability

Logs:

- request_id.
- organization_id.
- user_id.
- route.
- latency.
- error code.

Metrics:

- Request rate/error/latency.
- DB query latency.
- Pool saturation.
- Queue depth/message age.
- Import duration/failure.
- Search indexing lag.
- Business metric: lead conversion, opportunity value.

Trace theo request qua API, database, broker và worker.

Alert dựa trên triệu chứng user-facing và saturation, không chỉ CPU.

---

## 24. Security

- TLS mọi kết nối bên ngoài.
- Password hash bằng Argon2/bcrypt phù hợp.
- Access token ngắn hạn, refresh token có rotation nếu dùng.
- Rate limit login và API nhạy cảm.
- Input validation.
- Parameterized SQL.
- Tenant isolation test.
- Encryption at rest.
- Secret manager.
- Audit admin action.
- Data retention và quyền xóa dữ liệu.
- Redact PII trong log.

CSV export cần authorization và audit vì có thể chứa lượng lớn dữ liệu cá nhân.

---

## 25. Monolith hay microservice?

Bắt đầu modular monolith khi:

- Team nhỏ.
- Domain chưa ổn định.
- Cần transaction đơn giản.
- Scale chưa yêu cầu tách riêng.

Tách service khi:

- Boundary domain rõ.
- Team ownership độc lập.
- Workload cần scale khác nhau.
- Có yêu cầu isolation hoặc release độc lập.

Microservice tăng:

- Network failure.
- Distributed tracing.
- Eventual consistency.
- Deployment và schema contract.
- Chi phí vận hành.

Không dùng microservice để giải quyết codebase thiếu module và discipline.

---

## 26. Bài trình bày 5 phút

> Em bắt đầu bằng modular monolith Python API stateless để giảm độ phức tạp. PostgreSQL là source of truth vì CRM có nhiều quan hệ, transaction, phân quyền và audit. Mọi bảng nghiệp vụ đều có organization_id để cô lập tenant, kết hợp repository filter và có thể thêm Row-Level Security. Redis dùng cho cache và rate limit, không giữ dữ liệu CRM chính. Tác vụ chậm như import CSV, email và đồng bộ hệ thống ngoài được đưa qua broker cho worker xử lý. Em dùng outbox để tránh database commit nhưng publish event thất bại, và consumer phải idempotent. Search ban đầu dùng PostgreSQL, khi yêu cầu fuzzy search và scale tăng thì đồng bộ sang OpenSearch. Attachment lưu object storage. API scale ngang sau load balancer, đồng thời kiểm soát connection pool. Báo cáo nặng được chuyển qua read replica hoặc warehouse và Superset. Hệ thống cần audit log, tenant isolation, timeout, retry có giới hạn, metrics, logs, tracing, backup và restore test.

---

## 27. Câu hỏi phỏng vấn

### Vì sao chọn PostgreSQL?

CRM có dữ liệu quan hệ, constraint, transaction và query báo cáo. PostgreSQL phù hợp làm source of truth và vẫn hỗ trợ JSONB cho custom field có kiểm soát.

### Custom field thiết kế thế nào?

Có thể dùng bảng metadata + value tables, JSONB hoặc hybrid. JSONB nhanh để bắt đầu nhưng cần validation và index phù hợp; EAV linh hoạt nhưng query phức tạp. Quyết định dựa trên loại filter/report cần hỗ trợ.

### Làm sao tránh lộ dữ liệu giữa tenant?

Tenant ID bắt buộc trong auth context, mọi query filter theo tenant, unique/index gồm tenant ID, test cross-tenant, least privilege và có thể dùng PostgreSQL RLS làm lớp bảo vệ bổ sung.

### Import một triệu contact thế nào?

Upload object storage, tạo job, worker stream theo batch, validate, bulk insert/upsert, ghi error report, cập nhật progress và áp dụng backpressure/rate limit.

### Khi nào thêm Elasticsearch/OpenSearch?

Khi PostgreSQL search không còn đáp ứng yêu cầu typo tolerance, ranking, nhiều field và aggregation ở quy mô lớn. Search engine là read model eventually consistent, không phải source of truth.

### Làm sao xử lý hai người cùng sửa contact?

Optimistic locking bằng version và trả `409 Conflict` khi version cũ, hoặc pessimistic lock cho luồng thật sự cần serialize.

---

## 28. Bài tập thực hành

1. Thiết kế schema contact, company, opportunity, stage và activity.
2. Viết query cursor pagination theo `created_at, id`.
3. Thiết kế authorization cho admin, manager và sales owner.
4. Thiết kế luồng import 1 triệu contact.
5. Viết transaction đổi stage kèm audit và outbox.
6. Thiết kế dashboard mà không ảnh hưởng primary database.
7. Đưa ra kế hoạch tách notification worker khỏi modular monolith.
8. Liệt kê failure mode và cách degrade cho Redis, broker và search.

## Checklist

- [ ] Làm rõ functional và non-functional requirements.
- [ ] Có ước lượng quy mô và giả định.
- [ ] Thiết kế schema, index và pagination.
- [ ] Giải thích multi-tenancy và authorization.
- [ ] Hiểu transaction, optimistic locking, audit và outbox.
- [ ] Thiết kế import, notification, search và reporting.
- [ ] Nêu rõ cache, scaling, failure và observability.
- [ ] Biết trade-off monolith và microservice.
