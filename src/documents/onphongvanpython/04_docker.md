# 4. Docker cho Python Backend

Tài liệu này không chỉ liệt kê câu lệnh Docker mà tập trung giải thích Docker hoạt động như thế nào, vì sao Backend Developer cần dùng Docker, cách đóng gói một Python API, cách kết nối PostgreSQL và cách xử lý các lỗi thường gặp khi triển khai.

---

## Mục lục

1. Docker giải quyết vấn đề gì?
2. Image và container
3. Container khác máy ảo như thế nào?
4. Dockerfile
5. Layer và build cache
6. CMD và ENTRYPOINT
7. Port mapping và bind address
8. Volume và bind mount
9. Environment variable và secret
10. Docker network
11. Vì sao không dùng localhost giữa các container?
12. Healthcheck và service readiness
13. Docker Compose
14. Ví dụ FastAPI + PostgreSQL + Redis
15. Multi-stage build
16. Chạy container bằng non-root user
17. Resource limit
18. Logging
19. Quy trình build và deploy thực tế
20. Debug container theo từng bước
21. Các lỗi thường gặp
22. Câu hỏi phỏng vấn và câu trả lời mẫu
23. Bài tập thực hành
24. Checklist ôn tập

---

# 1. Docker giải quyết vấn đề gì?

Một ứng dụng Python thường phụ thuộc vào nhiều yếu tố:

- Phiên bản Python.
- Các package trong `requirements.txt` hoặc `pyproject.toml`.
- Thư viện hệ điều hành.
- Biến môi trường.
- Cấu trúc thư mục.
- Lệnh khởi động.

Nếu chỉ gửi source code cho người khác, ứng dụng có thể chạy trên máy của lập trình viên nhưng không chạy trên server vì môi trường khác nhau.

Ví dụ:

- Máy developer dùng Python 3.12 nhưng server dùng Python 3.10.
- Máy developer đã cài `libpq-dev` nhưng server chưa cài.
- Máy developer có biến `DATABASE_URL`, server lại thiếu.
- Developer chạy ứng dụng ở thư mục project, server chạy sai working directory.

Docker giúp đóng gói ứng dụng và các dependency thành một image có thể chạy nhất quán ở nhiều môi trường.

```text
Source code
+ Python runtime
+ Python packages
+ OS libraries
+ Start command
        ↓
    Docker image
        ↓
Dev / Staging / Production
```

Điều Docker mang lại không phải là “mọi môi trường hoàn toàn giống nhau”, vì kernel và hạ tầng vẫn có thể khác. Tuy nhiên, phần runtime của ứng dụng được kiểm soát tốt hơn rất nhiều.

## Ứng dụng thực tế

Trong CI/CD, ta thường chỉ build image một lần:

```text
Git push
   ↓
CI chạy test
   ↓
Build Docker image
   ↓
Push image lên registry
   ↓
Staging pull đúng image đó
   ↓
Production pull đúng image đó
```

Cách này tốt hơn việc SSH vào từng server rồi chạy `git pull` và cài dependency thủ công.

---

# 2. Image và container

## 2.1. Docker image là gì?

Docker image là một template bất biến dùng để tạo container.

Một image Python Backend có thể chứa:

- Base filesystem từ Debian hoặc Alpine.
- Python runtime.
- Package Python.
- Source code.
- User chạy ứng dụng.
- Lệnh mặc định.

Ví dụ luồng tạo image:

```text
Dockerfile → docker build → Image
```

Lệnh build:

```bash
docker build -t my-fastapi-app:1.0.0 .
```

Trong đó:

- `-t`: đặt tên và tag cho image.
- `my-fastapi-app`: tên image.
- `1.0.0`: tag.
- `.`: build context hiện tại.

## 2.2. Container là gì?

Container là một instance đang chạy được tạo từ image.

```text
Image → docker run → Container
```

Ví dụ:

```bash
docker run --name api-1 my-fastapi-app:1.0.0
```

Có thể tạo nhiều container từ cùng một image:

```bash
docker run --name api-1 my-fastapi-app:1.0.0
docker run --name api-2 my-fastapi-app:1.0.0
```

Hai container này dùng chung image nền nhưng có:

- Process riêng.
- Network namespace riêng.
- Writable layer riêng.
- Environment variable riêng.
- Lifecycle riêng.

Có thể hình dung gần đúng:

```text
Image giống class.
Container giống object được tạo từ class.
```

Đây chỉ là cách so sánh để dễ hiểu, không phải định nghĩa kỹ thuật chính xác tuyệt đối.

## 2.3. Image có chạy được không?

Image không phải process đang chạy. Muốn ứng dụng hoạt động phải tạo container từ image.

## 2.4. Xóa container có xóa image không?

Không.

```bash
docker rm api-1
```

Lệnh trên chỉ xóa container. Image vẫn tồn tại và có thể dùng để tạo container khác.

---

# 3. Container khác máy ảo như thế nào?

Máy ảo thường có hệ điều hành và kernel riêng. Container dùng chung kernel của host nhưng cô lập process, network, filesystem và resource bằng các cơ chế của Linux như namespace và cgroup.

```text
Máy ảo:
Hardware
  ↓
Host OS
  ↓
Hypervisor
  ↓
Guest OS + App

Container:
Hardware
  ↓
Host OS kernel
  ↓
Docker runtime
  ↓
Containerized process
```

## So sánh nhanh

| Tiêu chí | Container | Máy ảo |
|---|---|---|
| Kernel | Dùng chung kernel host | Có kernel riêng |
| Khởi động | Nhanh | Chậm hơn |
| Dung lượng | Thường nhỏ hơn | Thường lớn hơn |
| Mức cô lập | Process-level | Mạnh hơn ở mức máy |
| Use case | Deploy application | Chạy nhiều hệ điều hành, cô lập mạnh |

Container không phải giải pháp bảo mật tuyệt đối. Nếu container chạy quyền root, mount socket Docker hoặc được cấp capability nguy hiểm thì rủi ro vẫn rất lớn.

---

# 4. Dockerfile

Dockerfile mô tả cách Docker build image.

Ví dụ cơ bản cho FastAPI:

```dockerfile
FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .

RUN python -m pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

## 4.1. FROM

```dockerfile
FROM python:3.12-slim
```

`FROM` chọn base image.

`python:3.12-slim` thường phù hợp hơn `python:3.12` vì nhỏ hơn nhưng vẫn sử dụng Debian, dễ cài thư viện hệ thống hơn Alpine trong nhiều dự án Python.

Không nên dùng tag quá chung như:

```dockerfile
FROM python:latest
```

Vì lần build sau có thể lấy phiên bản Python khác và gây lỗi không dự đoán được.

## 4.2. WORKDIR

```dockerfile
WORKDIR /app
```

Thiết lập thư mục làm việc cho các lệnh sau.

Sau lệnh này:

```dockerfile
COPY requirements.txt .
```

sẽ copy file vào `/app/requirements.txt`.

## 4.3. COPY

```dockerfile
COPY requirements.txt .
COPY . .
```

`COPY` đưa file từ build context vào image.

Không nên copy toàn bộ source trước khi cài dependency nếu muốn tận dụng cache tốt.

## 4.4. RUN

```dockerfile
RUN python -m pip install --no-cache-dir -r requirements.txt
```

`RUN` chạy trong quá trình build image.

Kết quả của lệnh được lưu vào image layer.

## 4.5. CMD

```dockerfile
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

`CMD` là lệnh mặc định khi container khởi động.

Nên dùng exec form dạng JSON thay vì shell form:

```dockerfile
CMD uvicorn app.main:app --host 0.0.0.0 --port 8000
```

Exec form giúp process nhận signal trực tiếp tốt hơn, quan trọng khi container cần shutdown graceful.

## 4.6. EXPOSE

```dockerfile
EXPOSE 8000
```

`EXPOSE` chỉ mô tả rằng ứng dụng dự kiến dùng port 8000. Nó không tự publish port ra host.

Muốn truy cập từ host vẫn cần:

```bash
docker run -p 8000:8000 my-fastapi-app
```

## 4.7. ENV

```dockerfile
ENV PYTHONUNBUFFERED=1
```

Biến này giúp log Python được ghi ngay ra stdout thay vì bị buffer lâu.

Không nên lưu password trực tiếp trong Dockerfile:

```dockerfile
ENV DB_PASSWORD=secret
```

Vì secret có thể xuất hiện trong image metadata hoặc lịch sử build.

---

# 5. Layer và build cache

Docker image được tạo từ nhiều layer.

Ví dụ:

```dockerfile
FROM python:3.12-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY . .
```

Khi build lại, Docker có thể tái sử dụng layer cũ nếu instruction và dữ liệu đầu vào không đổi.

## Cách viết chưa tối ưu

```dockerfile
COPY . .
RUN pip install -r requirements.txt
```

Chỉ cần sửa một file source code, layer `COPY . .` thay đổi và bước cài dependency phải chạy lại.

## Cách tốt hơn

```dockerfile
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY . .
```

Nếu source code thay đổi nhưng `requirements.txt` không đổi, Docker có thể dùng cache của bước cài dependency.

## .dockerignore

File `.dockerignore` giúp loại bỏ file không cần thiết khỏi build context.

```text
.git
.venv
__pycache__
.pytest_cache
.env
*.pyc
tests
README.md
```

Lợi ích:

- Build context nhỏ hơn.
- Build nhanh hơn.
- Tránh vô tình đưa `.env`, Git history hoặc file nhạy cảm vào image.

---

# 6. CMD và ENTRYPOINT

## CMD

`CMD` cung cấp lệnh hoặc tham số mặc định.

```dockerfile
CMD ["main.py"]
```

## ENTRYPOINT

`ENTRYPOINT` định nghĩa executable chính.

```dockerfile
ENTRYPOINT ["python"]
CMD ["main.py"]
```

Khi chạy:

```bash
docker run myapp worker.py
```

Lệnh thực tế trở thành:

```text
python worker.py
```

## Khi nào dùng?

- Dùng `ENTRYPOINT` khi container luôn phải chạy một executable cố định.
- Dùng `CMD` cho lệnh mặc định dễ override.
- Với API Python đơn giản, chỉ dùng `CMD` thường đã đủ.

## Lỗi thường gặp với entrypoint script

Nếu dùng script:

```dockerfile
ENTRYPOINT ["/app/entrypoint.sh"]
```

thì script nên kết thúc bằng `exec`:

```bash
#!/bin/sh
set -e
exec "$@"
```

`exec` thay shell bằng process ứng dụng, giúp ứng dụng nhận SIGTERM trực tiếp khi Docker dừng container.

---

# 7. Port mapping và bind address

Lệnh:

```bash
docker run -p 8000:8000 my-fastapi-app
```

Có nghĩa:

```text
host_port:container_port
```

Request đi theo luồng:

```text
Client → Host port 8000 → Container port 8000 → FastAPI
```

## Vì sao ứng dụng phải bind 0.0.0.0?

Sai:

```bash
uvicorn app.main:app --host 127.0.0.1 --port 8000
```

`127.0.0.1` bên trong container chỉ nhận kết nối từ chính container.

Đúng:

```bash
uvicorn app.main:app --host 0.0.0.0 --port 8000
```

`0.0.0.0` cho phép process lắng nghe trên các interface mạng của container.

## Chỉ cho phép truy cập từ host

```bash
docker run -p 127.0.0.1:8000:8000 my-fastapi-app
```

Khi đó port chỉ bind vào loopback của host, phù hợp nếu Nginx trên cùng máy sẽ proxy vào ứng dụng.

---

# 8. Volume và bind mount

Writable layer của container không nên được dùng cho dữ liệu cần lưu lâu dài. Khi container bị xóa, dữ liệu trong writable layer có thể mất.

## 8.1. Named volume

```bash
docker volume create pgdata
```

Dùng trong Compose:

```yaml
services:
  postgres:
    image: postgres:16
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

Named volume do Docker quản lý, phù hợp với dữ liệu database trên môi trường đơn giản.

## 8.2. Bind mount

```yaml
services:
  api:
    volumes:
      - ./src:/app/src
```

Bind mount ánh xạ trực tiếp thư mục host vào container.

Phù hợp cho development vì sửa source trên host sẽ phản ánh ngay trong container.

## So sánh

| Tiêu chí | Named volume | Bind mount |
|---|---|---|
| Docker quản lý | Có | Không |
| Phụ thuộc đường dẫn host | Ít | Có |
| Phù hợp database | Có | Có thể, nhưng cần cẩn thận permission |
| Phù hợp hot reload | Không tối ưu | Rất phù hợp |

## Lưu ý database

Không mount một thư mục PostgreSQL đang được sử dụng đồng thời bởi nhiều container. PostgreSQL không được thiết kế để nhiều instance cùng ghi trực tiếp vào một data directory.

---

# 9. Environment variable và secret

Ví dụ:

```yaml
services:
  api:
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_NAME: app
```

Có thể dùng file env:

```yaml
env_file:
  - .env
```

`.env`:

```text
DB_USER=app_user
DB_PASSWORD=change-me
```

Không nên commit `.env` thật lên Git.

## Production secret

Docker Compose environment variable vẫn có thể bị xem bởi người có quyền Docker:

```bash
docker inspect api
```

Vì vậy, với production cần:

- Giới hạn quyền truy cập Docker daemon.
- Dùng secret manager khi hệ thống lớn hơn.
- Rotate secret định kỳ.
- Không ghi secret vào log.
- Không bake secret vào image.

---

# 10. Docker network

Docker network cho phép các container giao tiếp với nhau.

Trong Docker Compose, các service mặc định được đưa vào cùng một network và có thể gọi nhau bằng service name.

```yaml
services:
  api:
    build: .

  postgres:
    image: postgres:16
```

API sử dụng:

```text
DB_HOST=postgres
DB_PORT=5432
```

Docker cung cấp DNS nội bộ để phân giải tên `postgres` thành IP của container tương ứng.

## Không nên dùng IP container cố định

IP container có thể thay đổi sau khi recreate. Service name ổn định hơn và được Docker DNS quản lý.

## Tạo network riêng

```yaml
services:
  api:
    networks:
      - backend

  postgres:
    networks:
      - backend

networks:
  backend:
```

Có thể tách network public và private:

```text
Internet
   ↓
Nginx network
   ↓
API
   ↓
Backend private network
   ↓
PostgreSQL
```

PostgreSQL thường không cần publish port ra Internet.

---

# 11. Vì sao không dùng localhost giữa các container?

Mỗi container có network namespace riêng.

Bên trong container API:

```text
localhost = container API hiện tại
```

Nó không trỏ đến container PostgreSQL.

Sai:

```text
DB_HOST=localhost
```

Đúng:

```text
DB_HOST=postgres
```

Muốn gọi service chạy trên host trong Docker Desktop có thể dùng:

```text
host.docker.internal
```

Trên Linux, cách này phụ thuộc cấu hình và có thể cần thêm host mapping.

---

# 12. Healthcheck và service readiness

Container đang ở trạng thái `running` không có nghĩa ứng dụng đã sẵn sàng nhận request.

Ví dụ PostgreSQL process đã khởi động nhưng vẫn đang recovery hoặc chưa nhận connection.

Healthcheck PostgreSQL:

```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U app_user -d app"]
  interval: 10s
  timeout: 5s
  retries: 5
  start_period: 10s
```

Healthcheck API:

```yaml
healthcheck:
  test: ["CMD", "python", "-c", "import urllib.request; urllib.request.urlopen('http://localhost:8000/health')"]
  interval: 10s
  timeout: 3s
  retries: 3
```

## depends_on có đủ không?

```yaml
depends_on:
  - postgres
```

Cấu hình này chỉ đảm bảo container PostgreSQL được start trước, không đảm bảo database đã sẵn sàng.

Tốt hơn:

```yaml
depends_on:
  postgres:
    condition: service_healthy
```

Tuy nhiên ứng dụng vẫn nên có retry kết nối database vì service có thể restart sau khi API đã chạy.

---

# 13. Docker Compose

Docker Compose mô tả và chạy nhiều service liên quan trong cùng một project.

Ví dụ:

```yaml
services:
  api:
    build:
      context: .
    ports:
      - "8000:8000"
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_NAME: app
      DB_USER: app_user
      DB_PASSWORD: app_password
      REDIS_URL: redis://redis:6379/0
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    restart: unless-stopped

  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: app
      POSTGRES_USER: app_user
      POSTGRES_PASSWORD: app_password
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app_user -d app"]
      interval: 5s
      timeout: 3s
      retries: 10
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    restart: unless-stopped

volumes:
  pgdata:
```

Chạy:

```bash
docker compose up -d --build
```

Xem trạng thái:

```bash
docker compose ps
```

Xem log:

```bash
docker compose logs -f api
```

Dừng:

```bash
docker compose down
```

Dừng và xóa volume:

```bash
docker compose down -v
```

Cần cẩn thận với `-v` vì có thể xóa dữ liệu database.

---

# 14. Ví dụ FastAPI + PostgreSQL + Redis

## Luồng request

```text
Client
  ↓
Host port 8000
  ↓
FastAPI container
  ├── PostgreSQL qua postgres:5432
  └── Redis qua redis:6379
```

## Cấu hình Python

```python
import os

DATABASE_URL = (
    f"postgresql://{os.environ['DB_USER']}:"
    f"{os.environ['DB_PASSWORD']}@"
    f"{os.environ['DB_HOST']}:"
    f"{os.environ['DB_PORT']}/"
    f"{os.environ['DB_NAME']}"
)

REDIS_URL = os.environ["REDIS_URL"]
```

## Health endpoint

```python
from fastapi import FastAPI

app = FastAPI()


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}
```

Trong production, health endpoint có thể tách:

- Liveness: process còn sống.
- Readiness: ứng dụng có sẵn sàng nhận traffic hay không.

Không nên để healthcheck thực hiện query quá nặng hoặc gọi quá nhiều dependency ngoài.

---

# 15. Multi-stage build

Multi-stage build dùng nhiều stage nhưng image cuối chỉ giữ những thành phần cần chạy.

```dockerfile
FROM python:3.12-slim AS builder

WORKDIR /build

COPY requirements.txt .

RUN python -m pip wheel --wheel-dir /wheels -r requirements.txt

FROM python:3.12-slim AS runtime

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

WORKDIR /app

COPY --from=builder /wheels /wheels
COPY requirements.txt .

RUN python -m pip install --no-cache-dir --no-index --find-links=/wheels -r requirements.txt \
    && rm -rf /wheels

COPY app ./app

CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

Lợi ích:

- Image runtime không cần giữ compiler và file build tạm.
- Giảm kích thước image.
- Giảm bề mặt tấn công.

Multi-stage không phải lúc nào cũng làm image nhỏ hơn đáng kể. Cần đo thực tế bằng:

```bash
docker images
```

---

# 16. Chạy container bằng non-root user

Mặc định nhiều image chạy process bằng root trong container.

Tạo user riêng:

```dockerfile
RUN useradd --create-home --uid 10001 appuser

WORKDIR /app

COPY --chown=appuser:appuser . .

USER appuser
```

Lợi ích:

- Giảm tác động nếu ứng dụng bị khai thác.
- Hạn chế ghi vào các thư mục hệ thống.

Lưu ý:

- Thư mục ứng dụng phải có permission phù hợp.
- Port nhỏ hơn 1024 thường cần quyền đặc biệt.
- Volume mount từ host có thể gây lỗi owner/permission.

---

# 17. Resource limit

Container không tự động được bảo vệ khỏi việc sử dụng quá nhiều CPU hoặc RAM.

Docker run:

```bash
docker run --memory=512m --cpus=1.0 my-fastapi-app
```

Nếu ứng dụng vượt memory limit, process có thể bị OOM kill.

Kiểm tra:

```bash
docker inspect api
```

```bash
docker stats
```

Ứng dụng Python cần cấu hình số worker phù hợp với CPU và memory. Không nên tăng Gunicorn worker tùy ý vì mỗi worker là một process riêng và sử dụng thêm RAM.

---

# 18. Logging

Container nên ghi log ứng dụng ra stdout và stderr.

Python:

```python
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

logger.info("Application started")
```

Xem log:

```bash
docker logs -f api
```

Không nên chỉ ghi log vào file bên trong container vì:

- File có thể mất khi container bị xóa.
- Khó thu thập tập trung.
- Có thể làm đầy filesystem container.

Trong production, log thường được chuyển đến Loki, Elasticsearch, CloudWatch hoặc hệ thống log tập trung khác.

Không ghi password, token, số thẻ hoặc dữ liệu nhạy cảm vào log.

---

# 19. Quy trình build và deploy thực tế

Một quy trình đơn giản:

```text
Developer push code
       ↓
CI chạy lint và test
       ↓
Build image với tag commit SHA
       ↓
Scan vulnerability
       ↓
Push image lên registry
       ↓
Server pull image mới
       ↓
Chạy migration
       ↓
Recreate container
       ↓
Healthcheck
       ↓
Rollback nếu lỗi
```

## Vì sao nên dùng tag commit SHA?

Không nên chỉ dùng:

```text
myapp:latest
```

Nên có tag bất biến:

```text
myapp:9f21c0a
```

Khi cần rollback, server có thể chạy lại chính xác image cũ.

Có thể gắn đồng thời nhiều tag:

```text
myapp:9f21c0a
myapp:staging
```

Trong đó commit SHA dùng để truy vết, còn tag môi trường là alias tiện dụng.

---

# 20. Debug container theo từng bước

Khi API không truy cập được, không nên thử ngẫu nhiên. Hãy kiểm tra theo luồng.

## Bước 1: Container có chạy không?

```bash
docker ps -a
```

Nếu container `Exited`, xem log:

```bash
docker logs api
```

## Bước 2: Process đang listen port nào?

```bash
docker exec -it api sh
```

Trong container:

```bash
ss -lntp
```

Nếu image không có `ss`, có thể kiểm tra endpoint từ bên trong container.

## Bước 3: Ứng dụng bind 0.0.0.0 chưa?

Nếu chỉ bind `127.0.0.1`, port mapping sẽ không hoạt động như mong đợi.

## Bước 4: Port mapping đúng chưa?

```bash
docker port api
```

```bash
docker inspect api
```

## Bước 5: Test từ host

```bash
curl http://127.0.0.1:8000/health
```

## Bước 6: Kiểm tra network và DNS

```bash
docker network inspect project_default
```

Từ container API:

```bash
getent hosts postgres
```

## Bước 7: Kiểm tra kết nối database

Xác nhận:

- Hostname.
- Port.
- Database name.
- Username/password.
- PostgreSQL đã healthy.
- Cùng network.

## Bước 8: Kiểm tra environment variable

```bash
docker exec api env
```

Không đưa output chứa secret lên ticket công khai.

## Bước 9: Kiểm tra resource

```bash
docker stats
```

Kiểm tra container có bị OOMKilled:

```bash
docker inspect api
```

## Bước 10: Kiểm tra reverse proxy và firewall

Nếu curl localhost thành công nhưng bên ngoài không truy cập được, kiểm tra:

- Nginx.
- Host firewall.
- Security group.
- DNS.
- TLS certificate.
- Port publish chỉ bind `127.0.0.1` hay `0.0.0.0`.

---

# 21. Các lỗi thường gặp

## 21.1. Ứng dụng chạy trên máy nhưng không chạy trong container

Nguyên nhân thường gặp:

- Dependency chưa được khai báo.
- Sai working directory.
- Thiếu file do `.dockerignore`.
- Sai path import.
- Biến môi trường thiếu.
- Permission.
- Package cần thư viện hệ thống.
- Lệnh khởi động sai.

## 21.2. Không truy cập được API dù đã map port

Kiểm tra:

- Ứng dụng bind `0.0.0.0`.
- Container còn chạy.
- Container listen đúng port.
- Mapping đúng `host:container`.
- Firewall.

## 21.3. API không kết nối được PostgreSQL

Nguyên nhân:

- Dùng `localhost` thay vì service name.
- Không cùng network.
- PostgreSQL chưa ready.
- Sai credential.
- Sai database name.
- PostgreSQL container restart liên tục.
- Volume permission lỗi.

## 21.4. Container restart liên tục

Kiểm tra:

```bash
docker ps -a
```

```bash
docker logs --tail 200 api
```

Nguyên nhân thường gặp:

- Process chính thoát.
- Exception lúc startup.
- Migration lỗi.
- Không kết nối được dependency.
- Healthcheck sai kết hợp restart policy.
- OOM.

## 21.5. Image quá lớn

Cách cải thiện:

- Dùng base image phù hợp như `python:3.12-slim`.
- Dùng multi-stage build.
- Dùng `.dockerignore`.
- `pip install --no-cache-dir`.
- Không copy `.git`, test artifact, virtualenv.
- Xóa package build không cần ở runtime.
- Kiểm tra layer lớn thay vì chỉ đoán.

## 21.6. Thay source nhưng container vẫn chạy code cũ

Có thể do:

- Chưa rebuild image.
- Compose đang dùng image cũ.
- Build cache giữ layer do copy sai.
- Container chưa được recreate.
- Bind mount ghi đè source trong image.

Thử kiểm tra:

```bash
docker compose up -d --build --force-recreate
```

Không nên dùng `--no-cache` mặc định cho mọi build vì sẽ làm build chậm và che giấu vấn đề tổ chức layer.

---

# 22. Câu hỏi phỏng vấn và câu trả lời mẫu

## Docker image khác container như thế nào?

Image là template bất biến chứa filesystem và cấu hình dùng để tạo container. Container là instance đang chạy của image, có process, network namespace và writable layer riêng. Một image có thể tạo nhiều container.

## Container khác máy ảo như thế nào?

Máy ảo có guest OS và kernel riêng, còn container dùng chung kernel host và cô lập process bằng namespace, resource bằng cgroup. Vì vậy container thường nhẹ và khởi động nhanh hơn, nhưng mức cô lập không giống máy ảo.

## Vì sao ứng dụng trong container phải bind 0.0.0.0?

Vì `127.0.0.1` chỉ nhận kết nối từ loopback bên trong container. Port mapping chuyển traffic vào interface của container, nên ứng dụng phải listen trên `0.0.0.0` hoặc interface phù hợp.

## EXPOSE có publish port không?

Không. `EXPOSE` chỉ là metadata mô tả port dự kiến. Muốn publish port phải dùng `docker run -p` hoặc cấu hình `ports` trong Compose.

## Vì sao container API không dùng localhost để gọi PostgreSQL?

Mỗi container có network namespace riêng. `localhost` trong container API chỉ trỏ về container API đó. Phải dùng service name như `postgres`, được Docker DNS phân giải đến container database.

## Volume khác bind mount thế nào?

Named volume do Docker quản lý, ít phụ thuộc đường dẫn host và phù hợp lưu dữ liệu bền vững. Bind mount ánh xạ trực tiếp file hoặc thư mục host, rất tiện cho development nhưng phụ thuộc filesystem và permission của host.

## CMD khác ENTRYPOINT thế nào?

`ENTRYPOINT` xác định executable chính, còn `CMD` cung cấp lệnh hoặc tham số mặc định có thể override dễ hơn. Trong nhiều API container đơn giản, chỉ `CMD` là đủ.

## Làm sao giảm kích thước image Python?

Tôi sẽ dùng base image phù hợp, `.dockerignore`, copy dependency trước source, `pip --no-cache-dir`, multi-stage build khi có package cần compile và chỉ copy artifact runtime cần thiết. Tôi cũng đo từng layer thay vì tối ưu theo cảm tính.

## Docker Compose có phù hợp production không?

Compose có thể phù hợp với hệ thống nhỏ chạy trên một host nếu có backup, monitoring, restart policy và quy trình deploy rõ ràng. Khi cần multi-host orchestration, autoscaling, scheduling và self-healing phức tạp hơn thì Kubernetes hoặc nền tảng managed thường phù hợp hơn.

## depends_on có đảm bảo database đã sẵn sàng không?

Không nếu chỉ khai báo thứ tự start. Có thể kết hợp healthcheck và `condition: service_healthy`, nhưng ứng dụng vẫn cần retry vì database có thể restart trong lúc hệ thống đang chạy.

## Tại sao không nên chạy root trong container?

Nếu ứng dụng bị khai thác, process root có quyền lớn hơn trong container và có thể làm tăng tác động, đặc biệt khi có mount hoặc capability nguy hiểm. Chạy non-root là một lớp giảm thiểu rủi ro, không thay thế các biện pháp bảo mật khác.

## Cách trả lời tình huống API không chạy sau khi deploy

> Em kiểm tra theo luồng từ container đến network. Đầu tiên xem container có running hay restart không, sau đó đọc log và kiểm tra process có listen đúng port, đúng `0.0.0.0` hay không. Tiếp theo em kiểm tra port mapping, environment variable, DNS nội bộ và kết nối database. Nếu localhost trên server hoạt động nhưng bên ngoài không vào được, em kiểm tra Nginx, firewall, DNS và TLS. Em tránh sửa ngẫu nhiên mà xác định chính xác request đang hỏng ở lớp nào.

---

# 23. Bài tập thực hành

## Bài 1: Đóng gói FastAPI

Yêu cầu:

- Tạo endpoint `/health`.
- Viết Dockerfile.
- Chạy non-root.
- Publish port 8000.
- Kiểm tra bằng curl.

## Bài 2: FastAPI và PostgreSQL

Yêu cầu:

- Viết `compose.yaml`.
- API gọi database bằng hostname `postgres`.
- PostgreSQL dùng named volume.
- Có healthcheck.
- API chỉ start sau khi database healthy.
- API vẫn có retry kết nối.

## Bài 3: Debug lỗi localhost

Cố tình cấu hình:

```text
DB_HOST=localhost
```

Quan sát lỗi rồi sửa thành:

```text
DB_HOST=postgres
```

Giải thích vì sao.

## Bài 4: Tối ưu image

- Build image ban đầu.
- Ghi lại kích thước.
- Thêm `.dockerignore`.
- Dùng slim image.
- Dùng multi-stage build.
- So sánh kích thước và thời gian build.

## Bài 5: Mô phỏng deploy

- Tag image bằng commit SHA giả lập.
- Chạy version 1.
- Deploy version 2.
- Healthcheck thất bại.
- Rollback version 1.

---

# 24. Checklist ôn tập

- [ ] Giải thích được Docker giải quyết vấn đề gì.
- [ ] Phân biệt image và container.
- [ ] Phân biệt container và máy ảo.
- [ ] Hiểu `FROM`, `WORKDIR`, `COPY`, `RUN`, `CMD`, `ENTRYPOINT`.
- [ ] Hiểu Docker layer và build cache.
- [ ] Biết dùng `.dockerignore`.
- [ ] Hiểu `EXPOSE` không tự publish port.
- [ ] Giải thích được `0.0.0.0` và `127.0.0.1`.
- [ ] Phân biệt named volume và bind mount.
- [ ] Biết vì sao container không gọi nhau bằng localhost.
- [ ] Hiểu Docker DNS và service name.
- [ ] Biết viết Docker Compose cho API, PostgreSQL và Redis.
- [ ] Hiểu healthcheck và giới hạn của `depends_on`.
- [ ] Biết dùng multi-stage build.
- [ ] Biết chạy process bằng non-root user.
- [ ] Biết xem log, inspect, stats và network.
- [ ] Biết debug lỗi API không truy cập được.
- [ ] Biết debug lỗi kết nối PostgreSQL.
- [ ] Hiểu cách tag image và rollback.
- [ ] Có thể trả lời các câu hỏi phỏng vấn bằng ví dụ thực tế.
