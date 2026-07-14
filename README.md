# CodeLife Study Backend

Backend Go theo Clean Architecture, tối ưu cho VPS 2 GB RAM. Markdown hiện được đặt trong `src/documents` và nhúng vào binary khi build.

## Kiến Trúc

```text
cmd/api                    Composition root, ghép dependency
internal/domain            Entity và interface, không phụ thuộc framework
internal/usecase           Luồng nghiệp vụ
internal/adapter           PostgreSQL, Redis, embedded document repository
internal/delivery/httpapi  HTTP handler và middleware
src/documents              Tài liệu Markdown tạm thời
```

Chiều phụ thuộc luôn hướng vào trong:

```text
HTTP -> Use case -> Domain <- Repository/Cache adapters
```

PostgreSQL được giữ pool tối đa 5 connection. Redis dùng cache-aside với TTL và giới hạn 64 MB trong Compose. API vẫn chạy được khi không cấu hình PostgreSQL/Redis; khi có cấu hình, `/readyz` kiểm tra cả hai dịch vụ.

## API

- `GET /healthz`
- `GET /readyz`
- `GET /api/v1/documents`
- `GET /api/v1/documents/{slug}`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/verify-email`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`
- `GET /api/v1/progress` (Bearer token required)
- `PUT /api/v1/progress/{documentSlug}` (Bearer token required)

Learning progress is stored per authenticated user in PostgreSQL. Supported fields are `status`, `scrollPosition`, `note`, and `checkedFlashcards`. The API derives the user ID from the signed token and never accepts a user ID from the request body.

## Database migrations

Migrations are embedded in the Go binary and run automatically at startup. To change the database schema, add exactly one new file under `migrations` using this format:

```text
YYYYMMDDHHMMSS_short_description.sql
```

Example: `20260714153000_add_user_preferences.sql`. Do not copy SQL into `postgres.go`; the runner applies pending files in timestamp order inside transactions and records their names and checksums in `schema_migrations`. Never edit a migration after it has been deployed—add a new migration instead.

Response tài liệu có dạng:

```json
{
  "data": {
    "slug": "cam_nang_redis",
    "title": "Cẩm Nang Toàn Diện Về Redis",
    "category": "database",
    "wordCount": 1200,
    "readingTime": 6,
    "content": "# Cẩm Nang Toàn Diện Về Redis\n..."
  }
}
```

Endpoint danh sách trả cùng metadata nhưng không trả `content`, để frontend tải nhẹ hơn. Endpoint chi tiết mới trả nội dung Markdown đầy đủ.

Auth flow đơn giản:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

curl -X POST http://localhost:8080/api/v1/auth/verify-email \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","otp":"123456"}'

curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

Nếu chưa cấu hình SMTP, backend sẽ log OTP ra stdout để test local. Khi deploy thật, cấu hình `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM` và đổi `AUTH_TOKEN_SECRET` thành chuỗi random dài.

## Chạy Local

```bash
cp .env.example .env
docker compose up --build
```

Không dùng Docker:

```bash
go run ./cmd/api
```

## Kiểm Thử

```bash
go test ./...
```

## Ghi Chú VPS 2 GB

Tổng giới hạn bộ nhớ của API + PostgreSQL + Redis trong Compose khoảng 576 MB. Phần còn lại dành cho hệ điều hành, reverse proxy và frontend. Không dùng Kubernetes hoặc microservice ở quy mô này vì chi phí vận hành và RAM không cần thiết.
