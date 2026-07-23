# 4. Docker

## 1. Image và container

Image là template bất biến dùng để tạo container.

Container là instance đang chạy của image.

```text
Dockerfile → docker build → Image → docker run → Container
```

## 2. Dockerfile

```dockerfile
FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .

RUN python -m pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["python", "main.py"]
```

Các lệnh chính:

- `FROM`: image nền.
- `WORKDIR`: thư mục làm việc.
- `COPY`: sao chép file.
- `RUN`: chạy lúc build.
- `CMD`: lệnh mặc định khi container chạy.
- `ENTRYPOINT`: executable chính.
- `ENV`: biến môi trường.
- `EXPOSE`: mô tả port ứng dụng dùng.

## 3. Layer và cache

Mỗi lệnh trong Dockerfile tạo một layer.

Nên copy dependency trước source:

```dockerfile
COPY requirements.txt .
RUN python -m pip install -r requirements.txt
COPY . .
```

Nếu chỉ source đổi, Docker có thể dùng cache của layer cài dependency.

## 4. CMD và ENTRYPOINT

```dockerfile
ENTRYPOINT ["python"]
CMD ["main.py"]
```

Khi chạy:

```bash
docker run myapp worker.py
```

lệnh thực tế là:

```text
python worker.py
```

`ENTRYPOINT` thường là executable cố định; `CMD` là tham số mặc định.

## 5. Port mapping

```bash
docker run -p 8000:8000 myapp
```

Cú pháp:

```text
host_port:container_port
```

Ứng dụng trong container phải bind `0.0.0.0`, không chỉ `127.0.0.1`.

## 6. Volume và bind mount

### Named volume

```bash
docker volume create pgdata
```

```yaml
volumes:
  - pgdata:/var/lib/postgresql/data
```

### Bind mount

```yaml
volumes:
  - ./src:/app/src
```

Named volume do Docker quản lý, phù hợp dữ liệu bền vững. Bind mount ánh xạ trực tiếp thư mục host, phù hợp development.

## 7. Environment variable

```yaml
environment:
  DB_HOST: postgres
  DB_PORT: 5432
```

Không hard-code secret vào image hoặc commit lên Git.

## 8. Docker network

Các container trong cùng network có thể gọi nhau bằng service name.

```yaml
services:
  api:
    depends_on:
      - postgres

  postgres:
    image: postgres:16
```

API dùng:

```text
DB_HOST=postgres
```

không dùng `localhost`.

## 9. localhost trong container

`localhost` bên trong container trỏ chính container đó.

Muốn gọi container PostgreSQL khác, dùng tên service:

```text
postgres:5432
```

Muốn gọi host trên Docker Desktop:

```text
host.docker.internal
```

## 10. Healthcheck

```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U postgres"]
  interval: 10s
  timeout: 5s
  retries: 5
```

Healthcheck giúp xác định service đã thực sự sẵn sàng, không chỉ process đang tồn tại.

## 11. Docker Compose

```yaml
services:
  api:
    build: .
    ports:
      - "8000:8000"
    environment:
      DB_HOST: postgres
    depends_on:
      postgres:
        condition: service_healthy

  postgres:
    image: postgres:16
    environment:
      POSTGRES_PASSWORD: secret
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]

volumes:
  pgdata:
```

## 12. Multi-stage build

```dockerfile
FROM python:3.12-slim AS builder

WORKDIR /app
COPY requirements.txt .
RUN python -m pip wheel --wheel-dir /wheels -r requirements.txt

FROM python:3.12-slim

WORKDIR /app
COPY --from=builder /wheels /wheels
COPY requirements.txt .
RUN python -m pip install --no-cache-dir --no-index --find-links=/wheels -r requirements.txt

COPY . .

CMD ["python", "main.py"]
```

Giúp image cuối không chứa tool build không cần thiết.

## 13. Chạy non-root

```dockerfile
RUN useradd --create-home appuser
USER appuser
```

Giảm rủi ro khi container bị khai thác.

## 14. Debug container

```bash
docker ps
```

```bash
docker logs -f api
```

```bash
docker exec -it api sh
```

```bash
docker inspect api
```

```bash
docker stats
```

```bash
docker compose ps
```

```bash
docker compose logs -f api
```

## Câu hỏi phỏng vấn

### Ứng dụng chạy trên máy nhưng không chạy trong container?

Kiểm tra:

- Bind đúng `0.0.0.0`.
- Port mapping.
- Dependency đã cài.
- File đã copy.
- Environment variable.
- Permission.
- Database hostname.
- Working directory.
- Log container.

### Vì sao container không kết nối PostgreSQL?

Các nguyên nhân:

- Dùng `localhost`.
- Không cùng network.
- Sai service name.
- PostgreSQL chưa ready.
- Sai user/password/database.
- Firewall hoặc port.
- PostgreSQL chỉ listen local.

### Làm sao giảm kích thước Python image?

- Dùng `python:slim`.
- Multi-stage build.
- `--no-cache-dir`.
- `.dockerignore`.
- Không copy test, `.git`, cache.
- Gộp và xóa package build khi phù hợp.
