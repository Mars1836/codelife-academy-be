# 7. Linux cho Backend Developer

Linux là môi trường phổ biến để chạy backend, database, reverse proxy và container. Backend developer không nhất thiết phải là system administrator, nhưng cần đọc log, kiểm tra process, port, tài nguyên, quyền file và network để tự điều tra sự cố.

## 1. Filesystem và đường dẫn

Các thư mục thường gặp:

- `/etc`: cấu hình hệ thống và service.
- `/var/log`: log.
- `/var/lib`: dữ liệu trạng thái của service.
- `/opt`: phần mềm cài thêm.
- `/home`: thư mục người dùng.
- `/tmp`: file tạm.
- `/proc`: thông tin process và kernel dạng virtual filesystem.

Lệnh cơ bản:

```bash
pwd
ls -la
cd /var/log
find /var/log -type f -name "*.log"
du -sh /var/*
df -h
```

`df` kiểm tra dung lượng filesystem; `du` kiểm tra dữ liệu đang chiếm bao nhiêu trong một thư mục.

---

## 2. Đọc và tìm kiếm file

```bash
cat app.log
less app.log
head -n 100 app.log
tail -n 100 app.log
tail -f app.log
grep -n "ERROR" app.log
grep -R "DATABASE_URL" /etc/myapp
```

Với log lớn, ưu tiên `less`, `tail`, `grep`, không mở bằng editor nặng.

Tìm nhiều pattern:

```bash
grep -E "ERROR|CRITICAL|Traceback" app.log
```

Xem context:

```bash
grep -n -B 5 -A 10 "Traceback" app.log
```

---

## 3. Quyền file

Linux có ba nhóm quyền:

```text
owner | group | others
```

Mỗi nhóm có:

- `r`: read.
- `w`: write.
- `x`: execute hoặc traverse directory.

```bash
ls -l
chmod 640 .env
chmod +x deploy.sh
chown appuser:appgroup /opt/myapp
```

Không dùng `chmod 777` như cách sửa lỗi mặc định. Nó cấp quyền quá rộng và che mất nguyên nhân thật.

Với directory, quyền `x` cho phép đi xuyên vào directory. Có `r` nhưng thiếu `x` vẫn không thể truy cập file bên trong bình thường.

---

## 4. Process

```bash
ps aux
ps aux | grep gunicorn
pgrep -af gunicorn
top
htop
```

Các trạng thái quan trọng:

- Running.
- Sleeping.
- Zombie.
- Uninterruptible sleep, thường liên quan I/O.

Gửi signal:

```bash
kill -TERM <pid>
kill -KILL <pid>
```

`SIGTERM` cho process cơ hội shutdown an toàn. `SIGKILL` dừng ngay, không cleanup. Chỉ dùng `-9` khi process không phản hồi.

Backend nên xử lý SIGTERM để:

- Dừng nhận request mới.
- Hoàn thành request đang chạy.
- Đóng database pool.
- Flush log nếu cần.

---

## 5. systemd và journalctl

```bash
systemctl status nginx
systemctl restart nginx
systemctl enable nginx
journalctl -u nginx
journalctl -u nginx -f
journalctl -u myapp --since "30 minutes ago"
```

Ví dụ unit:

```ini
[Unit]
Description=Python Backend
After=network.target

[Service]
User=appuser
WorkingDirectory=/opt/myapp
EnvironmentFile=/etc/myapp/app.env
ExecStart=/opt/myapp/.venv/bin/gunicorn app.main:app -k uvicorn.workers.UvicornWorker -b 0.0.0.0:8000
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Sau khi sửa unit:

```bash
sudo systemctl daemon-reload
sudo systemctl restart myapp
```

---

## 6. Port và socket

```bash
ss -lntp
ss -lunp
lsof -i :8000
```

- `LISTEN`: process đang chờ kết nối.
- `127.0.0.1:8000`: chỉ truy cập từ máy local.
- `0.0.0.0:8000`: lắng nghe mọi IPv4 interface.

Ứng dụng chạy nhưng không truy cập từ ngoài thường do bind `127.0.0.1`, firewall, Docker port mapping hoặc reverse proxy.

Kiểm tra local:

```bash
curl -v http://127.0.0.1:8000/health
```

---

## 7. DNS và network

```bash
ip addr
ip route
ping 10.8.0.9
curl -v https://api.example.com
nslookup api.example.com
dig api.example.com
traceroute api.example.com
```

Kiểm tra port:

```bash
nc -vz 10.8.0.9 5432
```

`ping` thành công không có nghĩa port ứng dụng mở. ICMP và TCP là hai thứ khác nhau.

Luồng debug:

```text
DNS đúng?
→ Route đến IP được?
→ Port có mở?
→ Process có listen?
→ Firewall có chặn?
→ Reverse proxy có route đúng?
→ Application có trả lỗi?
```

---

## 8. Firewall

Ubuntu thường dùng UFW:

```bash
sudo ufw status
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
```

Không nên mở PostgreSQL `5432` ra toàn Internet. Hạn chế theo private network, VPN hoặc source IP cụ thể.

Ngoài firewall trong VPS còn có thể có cloud security group hoặc firewall của nhà cung cấp.

---

## 9. CPU, RAM và load average

```bash
top
free -h
uptime
vmstat 1
```

Load average không chỉ là CPU utilization; nó thể hiện số task đang chạy hoặc chờ tài nguyên, tùy hệ thống. Load cao nhưng CPU thấp có thể liên quan I/O wait.

Kiểm tra process dùng RAM:

```bash
ps aux --sort=-%mem | head
```

Kiểm tra CPU:

```bash
ps aux --sort=-%cpu | head
```

Nếu process bị kill không rõ lý do, kiểm tra OOM:

```bash
dmesg | grep -i "out of memory\|killed process"
journalctl -k | grep -i oom
```

---

## 10. Disk và inode

```bash
df -h
df -i
du -sh /var/log/* | sort -h
```

Disk còn dung lượng nhưng không tạo được file có thể do hết inode vì quá nhiều file nhỏ.

Log không rotate có thể làm đầy disk, khiến database hoặc Docker dừng ghi dữ liệu.

Dùng logrotate hoặc logging driver phù hợp.

---

## 11. Environment variable

```bash
printenv
env | sort
echo "$DATABASE_URL"
```

Environment của shell hiện tại không nhất thiết giống environment của systemd, Docker hoặc CI runner.

Kiểm tra systemd:

```bash
systemctl show myapp --property=Environment
```

Không in secret vào log hoặc chia sẻ output chứa password/token.

---

## 12. SSH

```bash
ssh user@server
ssh -i ~/.ssh/id_ed25519 user@server
scp file.txt user@server:/tmp/
rsync -avz ./dist/ user@server:/opt/myapp/
```

Quyền private key nên chặt:

```bash
chmod 600 ~/.ssh/id_ed25519
```

Cấu hình `~/.ssh/config`:

```text
Host staging
  HostName 10.8.0.20
  User xadmin
  IdentityFile ~/.ssh/id_ed25519_company
```

Sau đó:

```bash
ssh staging
```

Production nên ưu tiên SSH key, tắt root login và giới hạn nguồn truy cập.

---

## 13. Nginx reverse proxy

```nginx
server {
    listen 80;
    server_name api.example.com;

    location / {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Kiểm tra cấu hình:

```bash
sudo nginx -t
sudo systemctl reload nginx
```

`reload` ưu tiên hơn restart khi cấu hình hợp lệ vì ít gián đoạn hơn.

Log:

```bash
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

Lỗi `502 Bad Gateway` thường nghĩa Nginx không kết nối được upstream hoặc upstream đóng kết nối bất thường.

---

## 14. Docker trên Linux

```bash
docker ps
docker logs -f api
docker inspect api
docker stats
docker exec -it api sh
docker compose ps
docker compose logs -f api
```

Kiểm tra port mapping:

```bash
docker port api
```

Kiểm tra network:

```bash
docker network ls
docker network inspect <network>
```

`localhost` bên trong container là container đó, không phải host hoặc container database.

---

## 15. Quy trình debug API không truy cập được

### Bước 1: Xác định phạm vi

- Chỉ một user hay tất cả user?
- Local trên server có truy cập được không?
- Lỗi DNS, timeout, connection refused, 502 hay 500?

### Bước 2: Kiểm tra process/container

```bash
systemctl status myapp
docker ps
```

### Bước 3: Kiểm tra log

```bash
journalctl -u myapp -n 200
docker logs --tail 200 api
```

### Bước 4: Kiểm tra port

```bash
ss -lntp | grep 8000
curl -v http://127.0.0.1:8000/health
```

### Bước 5: Kiểm tra reverse proxy

```bash
nginx -t
curl -v -H "Host: api.example.com" http://127.0.0.1
```

### Bước 6: Kiểm tra DNS và firewall

```bash
dig api.example.com
sudo ufw status
```

### Bước 7: Kiểm tra dependency

```bash
nc -vz postgres.internal 5432
nc -vz redis.internal 6379
```

Đi từ ngoài vào trong hoặc từ dưới lên trên một cách có hệ thống, không restart ngẫu nhiên mọi service.

---

## 16. Quy trình debug server chậm

1. Kiểm tra CPU, RAM, load và disk.
2. Xác định process tiêu thụ tài nguyên.
3. Kiểm tra application latency và error rate.
4. Kiểm tra database slow query và connection pool.
5. Kiểm tra queue backlog.
6. Kiểm tra network và dependency.
7. So sánh với thời điểm triển khai gần nhất.

```bash
uptime
free -h
df -h
top
ss -s
```

Không kết luận “thiếu RAM” chỉ từ một chỉ số. Linux dùng RAM làm page cache; cần xem available memory, swap và OOM event.

---

## 17. Câu hỏi phỏng vấn

### Process khác thread?

Process có không gian địa chỉ riêng. Thread trong cùng process chia sẻ memory và tài nguyên process. Thread nhẹ hơn nhưng cần đồng bộ khi truy cập dữ liệu chung.

### SIGTERM khác SIGKILL?

SIGTERM yêu cầu process dừng và có thể cleanup. SIGKILL bị kernel dừng ngay, không thể bắt hoặc xử lý.

### Ứng dụng listen `127.0.0.1` và `0.0.0.0` khác nhau thế nào?

`127.0.0.1` chỉ nhận kết nối loopback trong máy/network namespace. `0.0.0.0` listen trên tất cả IPv4 interface.

### Vì sao `ping` được nhưng không kết nối port?

Ping dùng ICMP, còn ứng dụng dùng TCP/UDP. Port có thể không listen hoặc bị firewall chặn dù host vẫn trả lời ping.

### `df` và `du` khác nhau thế nào?

`df` hiển thị mức dùng của filesystem. `du` cộng kích thước file nhìn thấy trong directory. Chênh lệch có thể do file đã xóa nhưng process vẫn đang giữ file descriptor.

### 502 và 504 khác nhau thế nào?

502 thường là reverse proxy nhận phản hồi không hợp lệ hoặc không kết nối được upstream. 504 thường là upstream không trả lời trước timeout.

---

## 18. Bài tập thực hành

1. Chạy FastAPI bằng systemd và đọc log bằng journalctl.
2. Cấu hình Nginx reverse proxy đến port 8000.
3. Cố tình bind API vào `127.0.0.1`, sau đó giải thích phạm vi truy cập.
4. Tạo file chỉ owner được đọc bằng `chmod 600`.
5. Mô phỏng disk đầy bằng thư mục test và tìm file lớn.
6. Điều tra container không kết nối được PostgreSQL.
7. Viết checklist xử lý lỗi 502.

## Checklist

- [ ] Dùng được `less`, `tail`, `grep`, `find`, `df`, `du`.
- [ ] Hiểu permission và không lạm dụng `777`.
- [ ] Kiểm tra process, signal và systemd.
- [ ] Kiểm tra port bằng `ss`, `lsof`, `curl`, `nc`.
- [ ] Phân biệt DNS, route, firewall và application error.
- [ ] Kiểm tra CPU, RAM, disk, inode và OOM.
- [ ] Debug được Nginx và Docker ở mức cơ bản.
