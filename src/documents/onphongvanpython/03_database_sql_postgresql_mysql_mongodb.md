# 3. PostgreSQL, MySQL và MongoDB

Tài liệu này tập trung vào các phần thường được hỏi khi phỏng vấn Backend Python, đặc biệt trong hệ thống fintech.

## 1. JOIN

### INNER JOIN

Chỉ lấy các bản ghi khớp ở cả hai bảng.

```sql
SELECT u.id, u.name, o.id AS order_id
FROM users u
INNER JOIN orders o ON o.user_id = u.id;
```

### LEFT JOIN

Lấy toàn bộ bản ghi bên trái, kể cả khi không có dữ liệu khớp bên phải.

```sql
SELECT u.id, u.name, o.id AS order_id
FROM users u
LEFT JOIN orders o ON o.user_id = u.id;
```

### RIGHT JOIN

Lấy toàn bộ bản ghi bên phải. Trong thực tế thường có thể viết lại bằng `LEFT JOIN` để dễ đọc hơn.

## 2. GROUP BY và HAVING

```sql
SELECT user_id, COUNT(*) AS total_orders
FROM orders
GROUP BY user_id;
```

`WHERE` lọc trước khi nhóm, `HAVING` lọc sau khi nhóm.

```sql
SELECT user_id, COUNT(*) AS total_orders
FROM orders
WHERE status = 'completed'
GROUP BY user_id
HAVING COUNT(*) >= 5;
```

## 3. Subquery và CTE

### Subquery

```sql
SELECT *
FROM users
WHERE id IN (
    SELECT user_id
    FROM orders
    WHERE total_amount > 1000000
);
```

### CTE

```sql
WITH high_value_users AS (
    SELECT user_id, SUM(total_amount) AS total_spent
    FROM orders
    GROUP BY user_id
    HAVING SUM(total_amount) > 10000000
)
SELECT u.id, u.name, h.total_spent
FROM users u
JOIN high_value_users h ON h.user_id = u.id;
```

CTE thường dễ đọc hơn khi query có nhiều bước.

## 4. Index hoạt động thế nào?

Index giống mục lục của cuốn sách. Database có thể tìm nhanh vị trí bản ghi thay vì quét toàn bộ bảng.

```sql
CREATE INDEX idx_orders_user_id ON orders(user_id);
```

Index giúp tăng tốc đọc nhưng có nhược điểm:

- Tốn dung lượng lưu trữ.
- Làm `INSERT`, `UPDATE`, `DELETE` chậm hơn.
- Cần bảo trì.
- Quá nhiều index có thể làm optimizer chọn kế hoạch chưa tối ưu.

## 5. Composite index

```sql
CREATE INDEX idx_orders_user_created
ON orders(user_id, created_at);
```

Thường hiệu quả cho:

```sql
WHERE user_id = 10;
```

và:

```sql
WHERE user_id = 10 AND created_at >= '2026-01-01';
```

Thường không tối ưu cho:

```sql
WHERE created_at >= '2026-01-01';
```

Lý do là quy tắc leftmost prefix: index bắt đầu bằng `user_id`.

## 6. Khi nào index không được sử dụng?

Các trường hợp phổ biến:

```sql
WHERE LOWER(email) = 'a@example.com';
```

Nếu chỉ có index trên `email`, việc bọc cột trong hàm có thể làm index không dùng được. Có thể tạo functional index trong PostgreSQL.

```sql
CREATE INDEX idx_users_lower_email ON users(LOWER(email));
```

Một số trường hợp khác:

- Dùng `%keyword` ở đầu với `LIKE`.
- Ép kiểu không phù hợp.
- Bảng quá nhỏ.
- Query trả về phần lớn bảng.
- Kiểu dữ liệu hai phía không khớp.
- Statistics cũ.
- Thứ tự cột composite index không phù hợp.

## 7. Constraint

### Primary key

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL
);
```

### Foreign key

```sql
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id)
);
```

### Unique constraint

```sql
ALTER TABLE users
ADD CONSTRAINT uq_users_email UNIQUE(email);
```

Constraint bảo vệ tính đúng đắn dữ liệu tại database, không chỉ ở application.

## 8. Chuẩn hóa dữ liệu

### 1NF

Mỗi ô chứa một giá trị nguyên tử.

Không nên:

```text
skills = "Python,Docker,PostgreSQL"
```

### 2NF

Mọi cột không khóa phải phụ thuộc đầy đủ vào khóa chính.

### 3NF

Không để cột không khóa phụ thuộc gián tiếp vào khóa chính.

Chuẩn hóa giúp:

- Giảm trùng lặp.
- Tránh update anomaly.
- Giữ dữ liệu nhất quán.

Nhưng hệ thống báo cáo đôi khi chủ động denormalize để tăng tốc đọc.

## 9. Transaction

Transaction nhóm nhiều thao tác thành một đơn vị logic.

```sql
BEGIN;

UPDATE accounts
SET balance = balance - 100000
WHERE id = 1;

UPDATE accounts
SET balance = balance + 100000
WHERE id = 2;

COMMIT;
```

Nếu lỗi:

```sql
ROLLBACK;
```

## 10. ACID

- Atomicity: hoặc thành công toàn bộ, hoặc rollback toàn bộ.
- Consistency: dữ liệu luôn thỏa constraint và business rule.
- Isolation: transaction đồng thời không gây sai lệch ngoài mức cho phép.
- Durability: commit rồi thì dữ liệu phải được lưu bền vững.

## 11. Isolation level

### Read Uncommitted

Có thể đọc dữ liệu chưa commit.

### Read Committed

Chỉ đọc dữ liệu đã commit. PostgreSQL mặc định dùng mức này.

### Repeatable Read

Trong cùng transaction, đọc lại cùng dữ liệu cho kết quả ổn định hơn.

### Serializable

Mức cô lập cao nhất, hành vi gần như các transaction chạy tuần tự.

Các hiện tượng:

- Dirty read.
- Non-repeatable read.
- Phantom read.
- Lost update.

## 12. Lock

### Pessimistic lock

```sql
BEGIN;

SELECT balance
FROM accounts
WHERE id = 1
FOR UPDATE;
```

Dòng bị khóa đến khi transaction kết thúc.

### Optimistic lock

Thêm cột version:

```sql
UPDATE accounts
SET balance = 900000, version = version + 1
WHERE id = 1 AND version = 5;
```

Nếu số dòng cập nhật bằng 0, dữ liệu đã bị thay đổi bởi request khác.

## 13. Deadlock

Ví dụ:

- Transaction A khóa account 1 rồi chờ account 2.
- Transaction B khóa account 2 rồi chờ account 1.

Cách hạn chế:

- Luôn khóa tài nguyên theo cùng thứ tự.
- Giữ transaction ngắn.
- Không gọi API ngoài trong transaction.
- Có retry khi database phát hiện deadlock.
- Tạo index để giảm phạm vi lock.

## 14. Tránh hai người cùng cập nhật số dư

Cách an toàn:

```sql
UPDATE accounts
SET balance = balance - 100000
WHERE id = 1 AND balance >= 100000;
```

Sau đó kiểm tra số dòng bị ảnh hưởng.

Hoặc:

```sql
SELECT balance
FROM accounts
WHERE id = 1
FOR UPDATE;
```

Không nên:

```python
balance = get_balance()
balance -= amount
save(balance)
```

nếu không có transaction hoặc lock, vì dễ lost update.

## 15. N+1 query

Ví dụ:

```python
users = session.query(User).all()

for user in users:
    print(user.orders)
```

Nếu mỗi user tạo thêm một query, 100 user sẽ sinh 101 query.

Giải pháp:

- Eager loading.
- Join fetch.
- Batch loading.
- Query trực tiếp dữ liệu cần dùng.
- Theo dõi SQL log.

## 16. EXPLAIN

PostgreSQL:

```sql
EXPLAIN ANALYZE
SELECT *
FROM orders
WHERE user_id = 10;
```

Quan sát:

- Seq Scan.
- Index Scan.
- Rows estimated và actual rows.
- Cost.
- Execution time.
- Sort.
- Nested Loop, Hash Join, Merge Join.

`EXPLAIN ANALYZE` thực sự chạy query, nên cần cẩn thận với `UPDATE` hoặc `DELETE`.

## 17. Migration database

Migration là các thay đổi schema có phiên bản.

Ví dụ:

```text
001_create_users.sql
002_add_users_email.sql
003_create_orders.sql
```

Nguyên tắc:

- Migration phải được lưu trong Git.
- Không sửa migration đã chạy production.
- Tạo migration mới để thay đổi tiếp.
- Backup trước thay đổi rủi ro.
- Hạn chế migration khóa bảng lâu.
- Tách thay đổi breaking thành nhiều bước.

## 18. PostgreSQL, MySQL hay MongoDB?

### PostgreSQL

Phù hợp khi:

- Transaction phức tạp.
- Tính nhất quán cao.
- Quan hệ dữ liệu rõ.
- Query phân tích tốt.
- Cần JSON nhưng vẫn muốn hệ SQL mạnh.

### MySQL

Phù hợp khi:

- Web application phổ biến.
- Hệ sinh thái và vận hành quen thuộc.
- Workload CRUD tiêu chuẩn.
- Đội ngũ đã có kinh nghiệm MySQL.

### MongoDB

Phù hợp khi:

- Dữ liệu document thay đổi linh hoạt.
- Dữ liệu lồng nhau tự nhiên.
- Cần scale ngang theo document workload.
- Truy cập chủ yếu theo aggregate/document.

Không nên chọn MongoDB chỉ vì “không cần schema”. MongoDB vẫn cần thiết kế schema ở cấp ứng dụng.

## 19. Idempotency trong fintech

Một request chuyển tiền bị gửi lại không được trừ tiền hai lần.

Ví dụ client gửi:

```text
Idempotency-Key: transfer-20260724-001
```

Server lưu key cùng kết quả. Nếu nhận lại cùng key, trả lại kết quả cũ thay vì thực hiện giao dịch lần nữa.

## Câu hỏi phỏng vấn

### Index giúp truy vấn nhanh hơn nhưng có nhược điểm gì?

Tốn dung lượng, làm chậm ghi, cần bảo trì và có thể không được optimizer dùng nếu thiết kế sai.

### Index `(user_id, created_at)` có dùng khi chỉ lọc `created_at` không?

Thông thường không hiệu quả vì `user_id` là cột đầu tiên của composite index.

### WHERE và HAVING khác nhau thế nào?

`WHERE` lọc dòng trước `GROUP BY`; `HAVING` lọc nhóm sau khi aggregate.

### Khi nào dùng PostgreSQL, khi nào dùng MongoDB?

PostgreSQL phù hợp transaction, consistency và quan hệ phức tạp. MongoDB phù hợp document linh hoạt và pattern truy cập theo document. Quyết định phải dựa trên dữ liệu và truy vấn, không dựa vào xu hướng.
