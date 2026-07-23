# 3. PostgreSQL, MySQL và MongoDB

Tài liệu này giúp ôn phỏng vấn Backend Python theo hướng hiểu bản chất và áp dụng được vào hệ thống thực tế, đặc biệt là CRM, thương mại điện tử và fintech.

## Mục lục

1. Mô hình dữ liệu quan hệ và document
2. JOIN
3. GROUP BY, WHERE và HAVING
4. Subquery, CTE và window function
5. Index
6. Constraint và chuẩn hóa
7. Transaction và ACID
8. Isolation level và concurrency
9. Lock, deadlock và cập nhật số dư
10. N+1 query
11. EXPLAIN và tối ưu truy vấn
12. Migration database
13. PostgreSQL, MySQL và MongoDB
14. Idempotency trong fintech
15. Câu hỏi phỏng vấn
16. Bài tập thực hành

---

## 1. Mô hình dữ liệu quan hệ và document

### Cơ sở dữ liệu quan hệ là gì?

PostgreSQL và MySQL lưu dữ liệu theo bảng. Mỗi bảng có cột, kiểu dữ liệu, constraint và quan hệ với bảng khác.

Ví dụ:

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    status VARCHAR(30) NOT NULL,
    total_amount NUMERIC(18, 2) NOT NULL CHECK (total_amount >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Ưu điểm:

- Schema rõ ràng.
- Constraint bảo vệ tính đúng đắn của dữ liệu.
- Transaction mạnh.
- JOIN tốt.
- Phù hợp dữ liệu có nhiều quan hệ.

### MongoDB là gì?

MongoDB lưu dữ liệu dạng document BSON, gần giống JSON.

```json
{
  "_id": "order_1001",
  "customer": {
    "id": 10,
    "name": "Vũ Công Hậu"
  },
  "items": [
    {"sku": "BOOK-01", "quantity": 2, "price": 120000}
  ],
  "status": "paid"
}
```

MongoDB phù hợp khi dữ liệu thường được đọc và ghi theo một document hoàn chỉnh, schema thay đổi linh hoạt hoặc có cấu trúc lồng nhau tự nhiên.

Không nên chọn MongoDB chỉ vì “không cần thiết kế schema”. Dữ liệu production vẫn cần quy tắc, validation và chiến lược migration.

---

## 2. JOIN

Giả sử có hai bảng `users` và `orders`.

### INNER JOIN

Chỉ trả về các dòng khớp ở cả hai bảng.

```sql
SELECT u.id, u.full_name, o.id AS order_id, o.total_amount
FROM users AS u
INNER JOIN orders AS o ON o.user_id = u.id;
```

Ứng dụng: lấy danh sách khách hàng đã từng phát sinh đơn hàng.

### LEFT JOIN

Lấy toàn bộ dòng bên trái, kể cả khi bên phải không khớp.

```sql
SELECT u.id, u.full_name, COUNT(o.id) AS order_count
FROM users AS u
LEFT JOIN orders AS o ON o.user_id = u.id
GROUP BY u.id, u.full_name;
```

Ứng dụng: báo cáo toàn bộ khách hàng, bao gồm khách chưa có đơn.

### RIGHT JOIN

Lấy toàn bộ dòng bên phải. Trong thực tế thường đổi vị trí bảng và dùng `LEFT JOIN` để dễ đọc.

### Lỗi phổ biến khi JOIN

Nếu một user có nhiều order, JOIN tạo nhiều dòng cho user đó. Không nên thêm `DISTINCT` theo thói quen để che vấn đề; cần hiểu cardinality của quan hệ.

```sql
-- Sai mục tiêu nếu muốn mỗi user một dòng
SELECT u.*, o.*
FROM users u
JOIN orders o ON o.user_id = u.id;
```

Nếu muốn đơn gần nhất, có thể dùng window function hoặc `LATERAL` trong PostgreSQL.

```sql
SELECT u.id, u.full_name, latest.id AS latest_order_id
FROM users u
LEFT JOIN LATERAL (
    SELECT id
    FROM orders
    WHERE user_id = u.id
    ORDER BY created_at DESC
    LIMIT 1
) latest ON TRUE;
```

---

## 3. WHERE, GROUP BY và HAVING

Thứ tự logic đơn giản:

```text
FROM/JOIN → WHERE → GROUP BY → HAVING → SELECT → ORDER BY → LIMIT
```

`WHERE` lọc dòng trước khi nhóm. `HAVING` lọc kết quả sau khi aggregate.

```sql
SELECT user_id, COUNT(*) AS completed_orders
FROM orders
WHERE status = 'completed'
GROUP BY user_id
HAVING COUNT(*) >= 5;
```

Ứng dụng: tìm khách hàng có ít nhất 5 đơn hoàn thành.

Không nên dùng `HAVING` thay `WHERE` khi có thể lọc sớm, vì database phải xử lý nhiều dữ liệu hơn.

---

## 4. Subquery, CTE và window function

### Subquery

```sql
SELECT id, full_name
FROM users
WHERE id IN (
    SELECT user_id
    FROM orders
    WHERE total_amount >= 1000000
);
```

### CTE

CTE giúp chia truy vấn thành các bước có tên.

```sql
WITH spending AS (
    SELECT user_id, SUM(total_amount) AS total_spent
    FROM orders
    WHERE status = 'completed'
    GROUP BY user_id
)
SELECT u.id, u.full_name, s.total_spent
FROM users u
JOIN spending s ON s.user_id = u.id
WHERE s.total_spent >= 10000000;
```

### Window function

Window function tính toán trên tập dòng nhưng không gộp chúng thành một dòng như `GROUP BY`.

```sql
SELECT
    id,
    user_id,
    total_amount,
    ROW_NUMBER() OVER (
        PARTITION BY user_id
        ORDER BY created_at DESC
    ) AS order_rank
FROM orders;
```

Ứng dụng: lấy đơn gần nhất của mỗi khách hàng, xếp hạng doanh số, tính tổng lũy kế.

---

## 5. Index

### Index hoạt động thế nào?

Index là cấu trúc dữ liệu phụ giúp database tìm vị trí dòng nhanh hơn thay vì quét toàn bộ bảng. PostgreSQL thường dùng B-tree cho index mặc định.

```sql
CREATE INDEX idx_orders_user_id ON orders(user_id);
```

Query có thể chuyển từ `Seq Scan` sang `Index Scan`:

```sql
SELECT * FROM orders WHERE user_id = 100;
```

### Nhược điểm của index

- Tốn dung lượng.
- Làm INSERT, UPDATE, DELETE chậm hơn vì phải cập nhật index.
- Tăng chi phí bảo trì.
- Index thừa khiến optimizer có thêm lựa chọn không cần thiết.

### Composite index

```sql
CREATE INDEX idx_orders_user_created
ON orders(user_id, created_at DESC);
```

Phù hợp với:

```sql
SELECT *
FROM orders
WHERE user_id = 100
ORDER BY created_at DESC
LIMIT 20;
```

Theo nguyên tắc leftmost prefix, index trên `(user_id, created_at)` thường dùng tốt khi lọc `user_id`, hoặc `user_id` kết hợp `created_at`. Nó thường không tối ưu khi chỉ lọc `created_at`.

### Khi index có thể không được sử dụng?

- Bảng nhỏ, quét toàn bảng rẻ hơn.
- Điều kiện trả về phần lớn dữ liệu.
- Dùng hàm lên cột nhưng không có expression index.
- Ép kiểu làm sai kiểu cột.
- Dùng wildcard đầu chuỗi: `LIKE '%abc'`.
- Statistics cũ.
- Composite index không khớp cột đầu.

Ví dụ:

```sql
-- Index email thông thường có thể không được dùng
SELECT * FROM users WHERE LOWER(email) = 'a@example.com';

-- PostgreSQL expression index
CREATE INDEX idx_users_lower_email ON users(LOWER(email));
```

### Partial index

```sql
CREATE INDEX idx_orders_pending
ON orders(created_at)
WHERE status = 'pending';
```

Hữu ích khi thường xuyên truy vấn một tập con nhỏ.

---

## 6. Constraint và chuẩn hóa

### Constraint quan trọng

- `PRIMARY KEY`: định danh duy nhất và không null.
- `FOREIGN KEY`: bảo vệ quan hệ giữa bảng.
- `UNIQUE`: không cho phép trùng.
- `NOT NULL`: bắt buộc có giá trị.
- `CHECK`: kiểm tra điều kiện nghiệp vụ cơ bản.

```sql
CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id),
    balance NUMERIC(18, 2) NOT NULL DEFAULT 0 CHECK (balance >= 0)
);
```

Constraint là lớp bảo vệ cuối cùng. Validation trong Python không thay thế constraint vì dữ liệu có thể được ghi từ worker, script hoặc service khác.

### Chuẩn hóa

Chuẩn hóa giảm trùng lặp và anomaly khi cập nhật.

Không nên lưu tên khách hàng lặp lại trong mọi đơn nếu tên hiện tại phải đồng nhất. Tuy nhiên, invoice có thể cần snapshot tên và địa chỉ tại thời điểm mua. Đây là denormalization có chủ đích.

Thiết kế tốt dựa trên yêu cầu lịch sử và pattern truy cập, không áp dụng chuẩn hóa một cách máy móc.

---

## 7. Transaction và ACID

Transaction gom nhiều thao tác thành một đơn vị thành công hoặc thất bại.

```sql
BEGIN;

UPDATE accounts
SET balance = balance - 100000
WHERE id = 1;

UPDATE accounts
SET balance = balance + 100000
WHERE id = 2;

INSERT INTO transfers(from_account_id, to_account_id, amount)
VALUES (1, 2, 100000);

COMMIT;
```

Nếu có lỗi:

```sql
ROLLBACK;
```

### ACID

- Atomicity: tất cả hoặc không thao tác nào được áp dụng.
- Consistency: dữ liệu chuyển từ trạng thái hợp lệ này sang trạng thái hợp lệ khác.
- Isolation: transaction đồng thời không gây kết quả sai ngoài mức cho phép.
- Durability: đã commit thì dữ liệu phải tồn tại sau sự cố trong giới hạn bảo đảm của hệ thống.

Trong Python, transaction nên được quản lý rõ ràng:

```python
async with connection.transaction():
    await connection.execute(
        "UPDATE accounts SET balance = balance - $1 WHERE id = $2",
        amount,
        source_id,
    )
```

Không gọi API bên ngoài trong khi giữ transaction lâu, vì sẽ giữ connection và lock lâu.

---

## 8. Isolation level và concurrency

Các hiện tượng thường gặp:

- Dirty read: đọc dữ liệu chưa commit.
- Non-repeatable read: cùng một dòng đọc hai lần ra kết quả khác.
- Phantom read: cùng điều kiện query nhưng xuất hiện hoặc mất dòng.
- Lost update: hai transaction ghi đè kết quả của nhau.

Các isolation level phổ biến:

| Level | Ý nghĩa thực tế |
|---|---|
| Read Uncommitted | Cách ly rất thấp; PostgreSQL xử lý như Read Committed |
| Read Committed | Mỗi statement thấy dữ liệu đã commit trước statement đó |
| Repeatable Read | Transaction đọc một snapshot ổn định |
| Serializable | Kết quả tương đương chạy tuần tự, có thể phải retry |

Isolation càng cao thường giảm concurrency và tăng khả năng transaction bị abort hoặc chờ lock.

---

## 9. Lock, deadlock và cập nhật số dư

### Cách sai: đọc rồi ghi

```python
balance = await get_balance(account_id)
if balance >= amount:
    await update_balance(account_id, balance - amount)
```

Hai request có thể cùng đọc số dư 1.000.000 rồi cùng rút 800.000.

### Pessimistic locking

```sql
BEGIN;

SELECT balance
FROM accounts
WHERE id = 1
FOR UPDATE;

UPDATE accounts
SET balance = balance - 800000
WHERE id = 1;

COMMIT;
```

`FOR UPDATE` khóa dòng đến khi transaction kết thúc.

### Atomic update

Đây thường là cách ngắn gọn hơn:

```sql
UPDATE accounts
SET balance = balance - 800000
WHERE id = 1
  AND balance >= 800000
RETURNING balance;
```

Nếu không có dòng trả về, số dư không đủ hoặc tài khoản không tồn tại.

### Optimistic locking

```sql
UPDATE accounts
SET balance = 200000,
    version = version + 1
WHERE id = 1
  AND version = 5;
```

Nếu affected rows bằng 0, dữ liệu đã bị thay đổi và ứng dụng cần đọc lại hoặc retry.

### Deadlock

Deadlock xảy ra khi transaction A giữ lock 1 và chờ lock 2, còn transaction B giữ lock 2 và chờ lock 1.

Giảm deadlock bằng cách:

- Khóa tài nguyên theo cùng một thứ tự.
- Giữ transaction ngắn.
- Có index phù hợp để tránh khóa/quét quá nhiều dòng.
- Retry transaction khi database phát hiện deadlock.

Chuyển tiền nên khóa account theo ID tăng dần, bất kể chiều chuyển tiền.

---

## 10. N+1 query

N+1 xảy ra khi query một danh sách rồi query thêm một lần cho từng phần tử.

```python
users = await get_users()          # 1 query
for user in users:
    user.orders = await get_orders(user.id)  # N query
```

Cách xử lý:

- JOIN.
- Eager loading của ORM.
- Batch query với `WHERE user_id IN (...)`.
- DataLoader trong GraphQL.

Không phải mọi JOIN đều tốt. Nếu JOIN tạo lượng dữ liệu nhân lên rất lớn, batch query có thể dễ kiểm soát hơn.

---

## 11. EXPLAIN và tối ưu truy vấn

PostgreSQL:

```sql
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM orders
WHERE user_id = 100
ORDER BY created_at DESC
LIMIT 20;
```

Cần chú ý:

- `Seq Scan`, `Index Scan`, `Bitmap Heap Scan`.
- Estimated rows và actual rows.
- Execution time.
- Số loop.
- Sort có tràn ra disk không.
- Buffer hit và read.

`EXPLAIN ANALYZE` thực sự chạy query. Cẩn thận khi dùng với UPDATE hoặc DELETE; nên chạy trong transaction rồi rollback nếu cần.

Quy trình tối ưu:

1. Xác định query chậm bằng monitoring.
2. Lấy query và tham số đại diện.
3. Chạy `EXPLAIN ANALYZE` trên dữ liệu gần thực tế.
4. Kiểm tra index, cardinality và số dòng trung gian.
5. Sửa query hoặc index.
6. Đo lại thay vì đoán.

---

## 12. Migration database

Migration là thay đổi schema có version và có thể review.

Ví dụ thêm cột an toàn:

```sql
ALTER TABLE users ADD COLUMN phone VARCHAR(30);
```

Quy trình backward-compatible:

1. Thêm cột mới dạng nullable hoặc có default phù hợp.
2. Deploy code có thể đọc cả cột cũ và mới.
3. Backfill dữ liệu theo batch.
4. Chuyển toàn bộ traffic sang cột mới.
5. Thêm constraint nếu cần.
6. Xóa cột cũ ở release sau.

Không nên đổi tên hoặc xóa cột ngay trong cùng release nếu instance code cũ vẫn đang chạy.

Với bảng lớn, cần xem xét lock, thời gian rewrite bảng và online index creation.

PostgreSQL hỗ trợ:

```sql
CREATE INDEX CONCURRENTLY idx_orders_created_at
ON orders(created_at);
```

Lệnh này giảm blocking nhưng có quy tắc riêng và không chạy trong transaction block thông thường.

---

## 13. PostgreSQL, MySQL và MongoDB

### PostgreSQL

Phù hợp khi:

- Transaction và consistency quan trọng.
- Query phức tạp, CTE, window function.
- Cần kiểu dữ liệu phong phú như JSONB, array.
- Hệ thống CRM, ERP, fintech.

### MySQL

Phù hợp khi:

- Team đã có kinh nghiệm vận hành MySQL.
- Hệ sinh thái và hạ tầng đang chuẩn hóa theo MySQL.
- Workload CRUD quan hệ thông thường.

MySQL và PostgreSQL đều là lựa chọn tốt. Quyết định nên dựa trên yêu cầu, năng lực team và hệ sinh thái, không dựa vào tranh luận “database nào tốt nhất”.

### MongoDB

Phù hợp khi:

- Aggregate dữ liệu tự nhiên thành document.
- Cần schema linh hoạt có kiểm soát.
- Pattern đọc chủ yếu theo ID hoặc document.
- Dữ liệu lồng nhau được cập nhật cùng nhau.

Không phù hợp làm lựa chọn mặc định cho sổ cái tài chính có nhiều invariant và quan hệ phức tạp, trừ khi team hiểu rõ trade-off và thiết kế transaction phù hợp.

### Embed hay reference trong MongoDB?

Embed khi dữ liệu:

- Được đọc cùng nhau.
- Kích thước có giới hạn.
- Vòng đời giống nhau.

Reference khi dữ liệu:

- Được chia sẻ bởi nhiều document.
- Tăng trưởng không giới hạn.
- Thường xuyên cập nhật độc lập.

---

## 14. Idempotency trong fintech

Idempotency đảm bảo gửi lại cùng một request không tạo giao dịch trùng.

Client gửi header:

```text
Idempotency-Key: transfer-20260723-001
```

Database:

```sql
CREATE TABLE idempotency_keys (
    key VARCHAR(100) PRIMARY KEY,
    request_hash VARCHAR(64) NOT NULL,
    response_body JSONB,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Luồng xử lý:

1. Bắt đầu transaction.
2. Insert key duy nhất.
3. Nếu key đã tồn tại, trả kết quả cũ hoặc trạng thái đang xử lý.
4. Thực hiện nghiệp vụ.
5. Lưu response.
6. Commit.

Không chỉ cache idempotency key trong memory vì instance có thể restart hoặc request retry đi vào instance khác.

---

## 15. Câu hỏi phỏng vấn và cách trả lời

### Index giúp query nhanh hơn nhưng có nhược điểm gì?

Index giảm số dòng database phải quét, nhưng tốn dung lượng và làm chậm thao tác ghi vì mỗi thay đổi dữ liệu phải cập nhật index. Vì vậy chỉ tạo index theo query pattern đã xác định và cần đo bằng execution plan.

### Index `(user_id, created_at)` có dùng khi chỉ lọc `created_at` không?

Thông thường không tối ưu vì `user_id` là cột đầu. B-tree composite index hoạt động tốt theo leftmost prefix. Nếu query theo `created_at` thường xuyên, cần index riêng hoặc thiết kế index khác.

### WHERE và HAVING khác nhau thế nào?

`WHERE` lọc dòng trước `GROUP BY`; `HAVING` lọc nhóm sau aggregate. Nên dùng `WHERE` để giảm dữ liệu sớm khi điều kiện không phụ thuộc aggregate.

### Làm sao tránh hai người cùng cập nhật số dư?

Dùng transaction và thao tác atomic, ví dụ `UPDATE ... WHERE balance >= amount RETURNING ...`, hoặc khóa dòng bằng `SELECT ... FOR UPDATE`. Đồng thời cần idempotency để tránh retry tạo giao dịch trùng.

### Khi nào dùng PostgreSQL, khi nào dùng MongoDB?

PostgreSQL phù hợp dữ liệu quan hệ, transaction và consistency mạnh. MongoDB phù hợp aggregate dạng document, schema linh hoạt và pattern đọc theo document. Quyết định dựa trên invariant, quan hệ, query pattern, scale và kinh nghiệm vận hành.

### DELETE và soft delete khác nhau thế nào?

Hard delete xóa vật lý. Soft delete cập nhật `deleted_at`, hỗ trợ audit và khôi phục nhưng làm query và unique constraint phức tạp hơn. Không nên áp dụng soft delete cho mọi bảng theo mặc định.

---

## 16. Bài tập thực hành

1. Thiết kế schema cho user, account, transfer và ledger entry.
2. Viết query lấy 10 khách hàng có doanh số cao nhất trong 30 ngày.
3. Tạo composite index và giải thích query nào dùng được.
4. Mô phỏng hai transaction cùng rút tiền, sau đó sửa bằng atomic update.
5. Dùng `EXPLAIN ANALYZE` trước và sau khi thêm index.
6. Viết migration đổi cột `full_name` sang `display_name` theo hướng backward-compatible.
7. Thiết kế idempotency cho API tạo thanh toán.

## Checklist trước phỏng vấn

- [ ] Giải thích được JOIN và cardinality.
- [ ] Phân biệt WHERE và HAVING.
- [ ] Hiểu leftmost prefix của composite index.
- [ ] Giải thích được ACID và isolation level.
- [ ] Xử lý được lost update và deadlock.
- [ ] Nhận biết N+1 query.
- [ ] Đọc được execution plan cơ bản.
- [ ] Biết migration backward-compatible.
- [ ] So sánh PostgreSQL, MySQL và MongoDB theo use case.
- [ ] Trình bày được idempotency trong thanh toán.
