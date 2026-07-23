# 7. Linux cho Backend Developer

## 1. Điều hướng và file

```bash
pwd
```

```bash
ls -lah
```

```bash
cd /var/log
```

```bash
cp source.txt target.txt
```

```bash
mv old.txt new.txt
```

```bash
rm file.txt
```

```bash
rm -rf directory
```

Cẩn thận với `rm -rf`.

## 2. Đọc file và log

```bash
cat app.log
```

```bash
less app.log
```

```bash
head -n 20 app.log
```

```bash
tail -n 100 app.log
```

```bash
tail -f app.log
```

## 3. Tìm kiếm

```bash
grep "ERROR" app.log
```

```bash
grep -R "DATABASE_URL" /opt/app
```

```bash
find /var/log -name "*.log"
```

```bash
find /opt/app -type f -size +100M
```

## 4. Process

```bash
ps aux
```

```bash
ps aux | grep python
```

```bash
top
```

```bash
kill 1234
```

```bash
kill -15 1234
```

`SIGTERM` cho process cơ hội shutdown sạch.

```bash
kill -9 1234
```

`SIGKILL` dừng ngay, chỉ dùng khi cần.

## 5. Permission

```bash
chmod 644 file.txt
```

```bash
chmod 755 script.sh
```

```bash
chown appuser:appgroup file.txt
```

Ý nghĩa:

- Read = 4.
- Write = 2.
- Execute = 1.

`755`:

- Owner: rwx.
- Group: r-x.
- Other: r-x.

## 6. Network

```bash
curl http://localhost:8000/health
```

```bash
curl -I https://example.com
```

```bash
wget https://example.com/file.zip
```

```bash
ss -lntp
```

```bash
sudo ss -lntp | grep 8000
```

```bash
netstat -lntp
```

`ss` thường được ưu tiên hơn `netstat` trên hệ Linux hiện đại.

## 7. Disk và memory

```bash
df -h
```

Xem dung lượng filesystem.

```bash
du -sh /var/log/*
```

Xem thư mục nào chiếm dung lượng.

```bash
free -h
```

Xem RAM và swap.

## 8. systemd

```bash
sudo systemctl status myapp
```

```bash
sudo systemctl restart myapp
```

```bash
sudo systemctl enable myapp
```

```bash
sudo journalctl -u myapp -n 100 --no-pager
```

```bash
sudo journalctl -u myapp -f
```

## 9. Tình huống phỏng vấn

### Tìm process chiếm port

```bash
sudo ss -lntp | grep 8000
```

### Xem log realtime

```bash
tail -f /var/log/nginx/error.log
```

### Kiểm tra ổ đĩa

```bash
df -h
```

Sau đó:

```bash
du -sh /var/log/* | sort -h
```

### Kiểm tra service

```bash
sudo systemctl status myapp
```

### Kiểm tra API từ server

```bash
curl -v http://127.0.0.1:8000/health
```

### Tìm lỗi trong log

```bash
grep -n "ERROR" app.log
```

## 10. Quy trình debug server

1. Kiểm tra process.
2. Kiểm tra port.
3. Kiểm tra log.
4. Kiểm tra CPU/RAM/disk.
5. Kiểm tra network.
6. Kiểm tra environment.
7. Kiểm tra database.
8. Kiểm tra reverse proxy.
