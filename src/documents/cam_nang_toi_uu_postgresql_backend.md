## **CẨM NANG TOÀN DIỆN VỀ tối ưu HÓA POSTGRESQL VÀ KIẾN TRÚC HỆ THỐNG BACKEND** 

_Tài liệu phân tích chuyên sâu các giải pháp tối ưu hóa từ mức truy vấn đến kiến trúc hệ thống lớn_ 

## **LỜI NÓI ĐẦU** 

Trong các hệ thống phần mềm hiện đại, cơ sở dữ liệu thường là điểm nghẽn (bottleneck) lớn nhấtảnh hưởng đến hiệu năng tổng thể. Việc tối ưu hóa một hệ thống sử dụng PostgreSQL đòi hỏi kỹ sư phải có cái nhìn toàn diện, từ cách viết từng câu lệnh SQL, thiết kế cấu trúc bảng (Schema), quản lý chỉ mục (Index), điều phối kết nối (Connection Pooling), cho đến việcáp dụng các mô hình kiến trúc phân tán như Read Replica, Caching, Partitioning, Sharding và CQRS. Tài liệu này mở rộng chi tiết dựa trên các nguyên lý cốt lõi, cung cấp giải pháp thực tiễn giúp hệ thống vận hànhổn định dưới tải lượng lớn. 

## **1. tối ưu HÓA TRUY VẤN (QUERY OPTIMIZATION)** 

_Bản chất của tối ưu hóa truy vấn là giảm thiểu số lượng dữ liệu phải đọc từ ổ đĩa (I/O), tiết kiệm dung lượng RAM xử lý và tối giản băng thông mạng (Network)._ 

## **1.1 Cơ chế quét dữ liệu: Sequential Scan vs Index Scan** 

Khi thực hiện một câu lệnh tìm kiếm thông thường: 

```
-- Trước khi tối ưu: Quét tuần tự toàn bộ bảng (Sequential Scan)
SELECT * FROM orders WHERE customer_email = 'abc@gmail.com';
```

Nếu bảng `orders` có 10 triệu bản ghi và không được cấu hình chỉ mục (index), PostgreSQL buộc phải duyệt qua từng block dữ liệu trên đĩa từ đầu đến cuối bảng. Đây là tác vụ cực kỳ nặng về I/O. Khi thêm chỉ mục: 

```
CREATE INDEX idx_orders_customer_email ON orders(customer_email);
```

PostgreSQL sẽ chuyển sang cơ chế **Index Scan** hoặc **Bitmap Index Scan** . Hệ thống chỉ cần tìm kiếm trên cây chỉ mục B-tree với độ phức tạp O(log N), sau đó truy xuất trực tiếp đến con trỏ vật lý của bản ghi, giúp giảm thời gian phản hồi từ hàng giây xuống hàng miligiây. 

Tài liệu kỹ thuật - tối ưu hóa PostgreSQL & Backend System 

Trang 1 / 12 

## **1.2 Loại bỏ thói quen sử dụng SELECT *** 

Việc sử dụng `SELECT *` gây lãng phí tài nguyên nghiêm trọng trên ba phương diện: 

- **I/O đĩa:** Buộc PostgreSQL phải đọc toàn bộ các cột dữ liệu từ bộ nhớ đệm hoặc đĩa vào bộ nhớ làm việc. Đặc biệt nguy hiểm nếu bảng chứa các trường dữ liệu lớn như `TEXT` , `BLOB` hoặc `JSONB` (các trường này thường được lưuởvùng lưu trữ ngoài TOAST). 

- **RAM hệ thống:** Bộ nhớ đệm của Backend và Database bị chiếm dụng bởi các dữ liệu thừa không bao giờ dùng tới. 

• **Băng thông mạng (Network):** Tăng kích thước gói tin truyền tải giữa Database Server và Backend Server, gây nghẽn mạch khi số lượng request đồng thời tăng cao. 

> `--` x **`Không khuyến khích:`** `Lấy toàn bộtr ườ ng dữ liệu thừa SELECT * FROM users;` 

`--` **`Khuyến khích:`** `Chỉlấy chính xác những tr ườ ng cần hiển thịSELECT id, username FROM users;` 

## **1.3 Khắc phục triệt để lỗi N+1 Query** 

Lỗi N+1 xảy ra khi sử dụng các công cụ ORM (như TypeORM, Prisma, Hibernate) ở chế độ Lazy Loading. Hệ thống thực hiện 1 câu lệnh để lấy danh sách cha, sau đó lặp qua từng phần tử và chạy thêm N câu lệnh phụ để lấy dữ liệu con liên quan. 

> `--` x **`Mã nguồn dính lỗi N+1 (1000 ng ườ i dùng sẽtạo ra 1001 câu lệnh SQL):`** `const users = await userRepo.find(); for (const user of users) { await postRepo.find({ where: { userId: user.id } }); }` 

**Giải pháp tối ưu:** Tận dụng sức mạnh của các phép liên kết dữ liệu (JOIN) ở mức cơ sở dữ liệu để gộp thành một truy vấn duy nhất, hoặc sử dụng kĩ thuật Batching (như DataLoader trong GraphQL). 

`--` **`Câu lệnh SQL tối ưu sử dụng JOIN:`** `SELECT u.id, u.username, p.id as post_id, p.title FROM users u LEFT JOIN posts p ON p.user_id = u.id;` 

## **1.4 Bổ sung kĩ thuật: Phân trang hiệu năng cao (Pagination Optimization)** 

Thông thường, các nhà phát triển hayáp dụng cấu trúc phân trang dựa trên `OFFSET` : 

Tài liệu kỹ thuật - tối ưu hóa PostgreSQL & Backend System 

Trang 2 / 12 

> `--` x **`Kém hiệu năng khiứng dụng nhảy đ ế n các trang xa:`** `SELECT id, title FROM articles ORDER BY created_at DESC LIMIT 10 OFFSET 500000;` 

Với câu lệnh trên, PostgreSQL vẫn phải quét và sắp xếp qua 500.000 dòng đầu tiên trước khi bỏ qua chúng để lấy 10 dòng tiếp theo. Biện pháp thay thế là sử dụng **Phân trang theo con trỏ (Cursor-based / Keyset Pagination)** bằng cách dựa vào giá trị của bản ghi cuối cùng của trang trước: 

> `--` Vv **`tối ưu hoàn hảo nhờtận dụng Index trên cột created_at:`** `SELECT id, title FROM articles WHERE created_at < '2026-06-17 10:00:00' ORDER BY created_at DESC LIMIT 10;` 

## **2. CHIẾN LƯỢC QUẢN LÝ VÀ tối ưu CHỈ MỤC (INDEX OPTIMIZATION)** 

_Index là con dao hai lưỡi. Cấu hình đúng sẽ tăng tốc độ đọc dữ liệu vượt trội, nhưng lạm dụng hoặc thiết kế sai sẽ làm tê liệt hiệu năng ghi._ 

## **2.1 B-tree Index (Chỉ mục cây cân bằng)** 

Đây là kiểu chỉ mục mặc định và phổ biến nhất trong PostgreSQL. Nó duy trì một cấu trúc cây tự cân bằng giúp các thao tác tìm kiếm, thêm, xóa phần tử đạt độ phức tạp thời gian cực kỳ ổn định. 

**Phạm viáp dụng tối ưu:** Các toán tử so sánh chính xác ( `=` ), so sánh khoảng ( `<` , `>` , `<=` , `>=` ), mệnh đề `BETWEEN` , toán tử liệt kê `IN` , và đặc biệt là tối ưu hóa sắp xếp cho mệnh đề `ORDER BY` . 

`CREATE INDEX idx_user_email ON users(email);` 

## **2.2 Composite Index (Chỉ mục tổ hợp nhiều cột)** 

Khi một câu lệnh truy vấn thường xuyên lọc dữ liệu theo nhiềuđiều kiện đồng thời, việc tạo một chỉ mục tổ hợp là giải pháp bắt buộc. 

`-- Câu truy vấn mục tiêu: SELECT * FROM orders WHERE status = 'SUCCESS' AND created_at > NOW() - INTERVAL '7 days'; -- Tạo Composite Index tối ưu: CREATE INDEX idx_order_status_created ON orders(status, created_at);` 

Tài liệu kỹ thuật - tối ưu hóa PostgreSQL & Backend System 

Trang 3 / 12 

**Nguyên tắc tiền tố bên trái (Left-most Prefix Rule):** Thứ tự khai báo cột trong Composite Index quyết định khả năng tái sử dụng của nó. Chỉ mục `(status, created_at)` sẽ hỗ trợ tốt cho truy vấn lọc theo `(status, created_at)` hoặc chỉ lọc theo `(status)` . Tuy nhiên, nếu câu truy vấn chỉ lọc duy nhất theo cột `(created_at)` , PostgreSQL sẽ không thể tận dụng chỉ mục này một cách hiệu quả. 

## **2.3 Bổ sung kĩ thuật nâng cao: Partial Index & Covering Index** 

• **Partial Index (Chỉ mục một phần):** Chỉ đánh chỉ mục cho các bản ghi thỏa mãn mộtđiều kiện nhất định. Giúp tiết kiệm dung lượng đĩa và RAM rất lớn. 

```
-- Chỉ tạo index cho người dùngđang hoạt động, bỏ qua các tài khoản đã bị khóa
CREATE INDEX idx_active_users_email ON users(email) WHERE is_active = true;
```

• **Covering Index (Chỉ mục bao phủ với mệnh đề INCLUDE):** Cho phép đính kèm dữ liệu của các cột phụ vào tầng lá của index mà không biến chúng thành một phần của khóa tìm kiếm. Từ đó giúp truy vấn thực hiện _Index Only Scan_ (không cần quay lại Table Heap để đọc dữ liệu). 

```
CREATE INDEX idx_users_username_incl_email ON users(username) INCLUDE(email);
```

## **2.4 Tác hại của việc thừa chỉ mục và cách dọn dẹp** 

Mỗi khi một hành động `INSERT` , `UPDATE` , hoặc `DELETE` diễn ra, PostgreSQL không chỉ ghi dữ liệu vào bảng gốc mà còn phải cập nhật lại toàn bộ các cây chỉ mục liên quan. Nếu một bảng có 10 chỉ mục, một thao tác ghi dữ liệu sẽ nhân thêm 10 lần công việc xử lý cho hệ thống, gây sụt giảm nghiêm trọng throughput ghi. 

Nhà phát triển cần định kỳ kiểm tra các chỉ mục không bao giờ được sử dụng bằng cách truy vấn view hệ thống để tiến hành xóa bỏ (Drop Index): 

```
SELECT indexrelname, idx_scan FROM pg_stat_user_indexes WHERE idx_scan = 0;
```

## **3. PHÂN TÍCH KẾ HOẠCH THỰC THI (EXECUTION PLAN ANALYSIS)** 

_Trước khi tiến hành bất kỳ chỉnh sửa nào liên quan đến SQL hoặc Index, kỹ sư cần phải hiểu rõ cách thức mà Database Planner vận hành câu lệnh đó thông qua công cụ giải trình._ 

Tài liệu kỹ thuật - tối ưu hóa PostgreSQL & Backend System 

Trang 4 / 12 

## **3.1 Sự khác biệt giữa EXPLAIN và EXPLAIN ANALYZE** 

|**Công cụ**|**Cơ chế hoạt động**|**Trường hợp sử dụng**|
|---|---|---|
|**EXPLAIN**|Dựa trên sốliệu thống kê nội bộ(Statistics)<br>đểđưa ra kếhoạch thực thi dựkiến và chi<br>phí (Cost)ước tính. Không chạy câu lệnh<br>thực tế.|Dùng đểkiểm tra nhanh cấu trúc truy vấn<br>hoặc khi câu lệnh quá nặng, có rủi ro làm<br>sập hệ thống nếu thực thi thật.|
|**EXPLAIN**<br>**ANALYZE**|Bắt buộc PostgreSQL phải chạy câu lệnh<br>SQL trên dữ liệu thật, đo đạc chính xác thời<br>gian xử lý (Actual Time) và sốlượng dòng<br>trảvề ởtừng node tác vụ.|Dùng trong môi trường Staging/<br>Development đểđo lường chính xác hiệu<br>năng thực tếcủa câu lệnh cần tối ưu.|



## **3.2 Các thuật ngữ trọng tâm cần lưuýkhi đọc kế hoạch** 

- **Seq Scan (Sequential Scan):** Dấu hiệu cảnh báo hệ thốngđang phải quét toàn bộ bảng dữ liệu. Cần xem xét bổ sung Index phù hợp. 

- **Index Scan / Index Only Scan:** Câu truy vấn đã tận dụng tốt chỉ mục, định vị bản ghi đích cực nhanh. 

- **Bitmap Index Scan + Bitmap Heap Scan:** Cơ chế trung gian, tối ưu khi cần lấy một lượng dữ liệu vừa phải nằm rải rác trong nhiều trang (pages). 

## **4. QUẢN LÝ KẾT NỐI QUA CONNECTION POOL (CONNECTION POOLING)** 

_Trong kiến trúc PostgreSQL, mỗi một kết nối (connection) từ ứng dụng backend sẽ thiết lập một tiến trình riêng biệt (Process-per-connection) ở phía Server._ 

Một tiến trình kết nối tiêu tốn khoảng 10MB RAM cố định cùng với chi phí quản lý context switch của CPU. Nếu cấu hình tham số `max_connections = 1000` trực tiếp trên PostgreSQL và có 1000 request đồng thời đổ vào, hệ thống Server sẽ ngay lập tức rơi vào trạng thái cạn kiệt tài nguyên RAM và treo hoàn toàn do CPU bận rộn tranh chấp tài nguyên. 

## **4.1 Giải pháp tối ưu: Sử dụng Connection Pooler** 

Thay vì kết nối trực tiếp, hệ thống sử dụng một lớp trung gian để quản lý và tái sử dụng các kết nối sẵn có, giới hạn số lượng kết nối thực tế tới databaseởmột ngưỡng tối ưu (thường từ 20-50 kết nối cho mỗi instance backend). 

- **PgBouncer:** Công cụ pooling chuyên dụng gọn nhẹ mức hệ thống dành cho PostgreSQL. Khuyến khích sử dụng chế độ **Transaction Pooling** (PgBouncer sẽ mượn kết nối khi một transaction bắt đầu và trả lại ngay khi transaction kết thúc, giúp tối ưu hóa số lượng client vượt trội). 

Tài liệu kỹ thuật - tối ưu hóa PostgreSQL & Backend System 

Trang 5 / 12 

- **ORM Pool (TypeORM / Prisma Pool):** Cấu hình trực tiếp trong mã nguồnứng dụng để kiểm soát số lượng kết nối tốiđa mà một instance backend được phép khởi tạo. 

```
// Cấu hình Connection Pool mẫu trongứng dụng Node.js/TypeORM
{
    type: "postgres",
    host: "localhost",
    username: "postgres",
    database: "production_db",
    extra: {
```

```
        max: 30,              // Số lượng connection tốiđa trong pool của instance
này
        min: 5,               // Số lượng connection tối thiểu duy trì
        idleTimeoutMillis: 30000, // Đóng kết nối nếu không sử dụng sau 30 giây
        connectionTimeoutMillis: 2000 // Giới hạn thời gian chờ kết nối tốiđa 2
giây
    }
}
```

## **5. CHIẾN LƯỢC CACHING NÂNG CAO (CACHING STRATEGIES)** 

_Nguyên lý cốt lõi: Truy vấn nhanh nhất là truy vấn không phải chạy vào Cơ sở dữ liệu gốc._ 

## **5.1 Mô hình Cache-Aside (Lazy Loading) với Redis** 

Áp dụng cho các vùng dữ liệu có tần suất đọc cực cao nhưngít khi thay đổi biến động (như thông tin cấu hình hệ thống, thông tin danh mục sản phẩm, hoặc hồ sơ cá nhân người dùng). 

```
 Luồng xử lý:
 Request ──> [ Kiểm tra Redis ] ──(Có dữ liệu - Cache Hit)──> Trả về Client
                    │
            (Không có - Cache Miss)
                    ▼
          [ Truy vấn PostgreSQL ] ──> Cập nhật vào Redis (kèm TTL) ──> Trả về
Client
```

Mô hình này giúp giảm tải từ 80% đến 90% áp lực truy vấn trực tiếp xuống PostgreSQL, bảo vệ DB khỏi tình trạng quá tải. 

## **5.2 Phòng chống các lỗi nghiêm trọng về Cache** 

- **Cache Penetration (Xuyên thủng cache):** Xảy ra khi tin tặc liên tục gửi yêu cầu với các mã ID không hề tồn tại trong hệ thống. Hệ thống kiểm tra cache không thấy sẽ liên tục chuyển hướng truy vấn xuống 

Tài liệu kỹ thuật - tối ưu hóa PostgreSQL & Backend System 

Trang 6 / 12 

DB. _Giải pháp:_ Sử dụng Bloom Filterởtầng trước cache hoặc lưu trữ cả các key trống với thời gian hết hạn (TTL) rất ngắn (ví dụ 1-2 phút). 

- **Cache Stampede / Thundering Herd (Bão sập cache):** Xảy ra khi một key cực kỳ hot (ví dụ trang chủ sự kiện lớn) hết hạn TTL. Hàng vạn request cùng lúc nhận thấy cache trống và đồng loạtép vào Database. _Giải pháp:_ Áp dụng cơ chế khóa phân tán (Mutex Lock) để chỉ cho phép 1 request duy nhất xuống DB cập nhật cache, các request khác chờ đợi; hoặc chạy worker cập nhật cache nền trước khi TTL chính thức hết hạn. 

## **6. PHÂN VÙNG DỮ LIỆU THEO CHIỀU NGANG (TABLE PARTITIONING)** 

_Khi một bảng dữ liệu tích lũy đến hàng trăm triệu hoặc hàng tỷ bản ghi (ví dụ bảng chứa logs hệ thống, lịch sử thanh toán), hiệu năng của các cây index cũng bắt đầu suy giảm do kích thước quá lớn vượt ngoài dung lượng RAM._ 

## **6.1 Declarative Partitioning trong PostgreSQL** 

Giải pháp là chia nhỏ một bảng logic khổng lồ thành các bảng vật lý nhỏ hơn dựa trên một tiêu chí phân vùng cố định. PostgreSQL hỗ trợ 3 phương thức phân vùng chính từ phiên bản 10 trở đi: 

- **PARTITION BY RANGE:** Chia theo phạm vi, cực kỳ thích hợp cho dữ liệu chuỗi thời gian (Timeseries). Ví dụ phân vùng dữ liệu logs theo từng tháng: `logs_2026_01` , `logs_2026_02` ... 

- **PARTITION BY LIST:** Chia theo một danh mục giá trị cụ thể, ví dụ phân vùng hóa đơn theo khu vực địa lý hoặc trạng thái quốc gia. 

- **PARTITION BY HASH:** Sử dụng hàm băm để phân bổ đều dữ liệu vào một số lượng bảng con định sẵn nhằm mục đích chia đều tải trọng ghi. 

## **6.2 Cơ chế loại trừ phân vùng (Partition Pruning)** 

Khi thực hiện câu lệnh truy vấn có kèm theođiều kiện lọc của khóa phân vùng: 

```
SELECT * FROM system_logs WHERE created_at >= '2026-06-01' AND created_at <=
'2026-06-15';
```

Bộ tối ưu hóa của PostgreSQL sẽ kích hoạt tính năng **Partition Pruning** . Nó ngay lập tức xác định bản ghi chỉ nằm trong bảng vật lý `logs_2026_06` và hoàn toàn bỏ qua việc quét tất cả các bảng phân vùng của các tháng khác, giúp tốc độ truy vấn nhanh tương đương với việc truy vấn trên một bảng nhỏ độc lập. 

Tài liệu kỹ thuật - tối ưu hóa PostgreSQL & Backend System 

Trang 7 / 12 

## **7. KIẾN TRÚC BẢN SAO CHỈ ĐỌC (READ REPLICA ARCHITECTURE)** 

_Hầu hết các hệ thống trực tuyến đều có đặc thù tải lượng thiên về Đọc (Read chiếm 80-90%, Write chỉ chiếm 10-20%). Kiến trúc Read Replica ra đời để giải quyết bài toán mất cân bằng này._ 

```
 Sơ đồ luồng:
                           [ Ứng dụng Backend ]
                            /                             (Thao tác
Ghi:                   (Thao tác Đọc:
          INSERT, UPDATE, DELETE)               SELECT)
                          /
▼                      ▼
                [ Primary Node ] ───(WAL)───> [ Replica Node ]
```

## **7.1 Cơ chế hoạt động và cấu hìnhđiều hướng tầng Backend** 

Hệ thống bao gồm một node chính (Primary) nhận toàn bộ các tác vụ thay đổi dữ liệu, sau đó bất đồng bộ hoặc đồng bộ các bản ghi nhật ký WAL (Write-Ahead Logging) sang một hoặc nhiều máy chủ sao lưu (Replica). Các máy chủ Replica này cấu hìnhởchế độ Read-Only. 

Tại tầng mã nguồnứng dụng (ví dụ trong NestJS với TypeORM), kỹ sư cần cấu hình driver để tự động phân tách luồng truy vấn: 

Tài liệu kỹ thuật - tối ưu hóa PostgreSQL & Backend System 

Trang 8 / 12 

```
// Cấu hình phân tách Read/Write Replication trong TypeORM
{
    type: "postgres",
    replication: {
        master: {
            host: "primary-db.domain.com",
            port: 5432,
            username: "repl_master",
            password: "password"
        },
        slaves: [
            {
                host: "replica-db-1.domain.com",
                port: 5432,
                username: "repl_slave",
                password: "password"
            },
            {
                host: "replica-db-2.domain.com",
                port: 5432,
                username: "repl_slave",
                password: "password"
            }
        ]
    }
}
```

## **8. PHÂN MẢNH CƠ SỞ DỮ LIỆU HỆ THỐNG LỚN (DATABASE SHARDING)** 

_Khi một máy chủ vật lý đơn lẻ đạt tới giới hạn phần cứng (kể cả khi đã tối ưu Read Replica và Partitioning), Sharding là giải pháp tối thượngởmức độ mở rộng theo chiều ngang (Horizontal Scaling)._ 

Sharding tuân thủ kiến trúc "Không chia sẻ tài nguyên" (Share-nothing architecture). Dữ liệu của một bảng được phân tán ra nhiều máy chủ Database độc lập hoàn toàn về mặt vật lý (mỗi máy gọi là một Shard). 

## **8.1 Tiêu chí lựa chọn Sharding Key và Thách thức lớn** 

- **Lựa chọn Sharding Key:** Cần chọn một trường dữ liệu đóng vai trò phân loại gốc (ví dụ `tenant_id` trong hệ thống SaaS hoặc `user_id` trong mạng xã hội) để thuật toán băm phân phối dữ liệu đồng đều, ngăn ngừa hiện tượng "Hot Shard" (một shard chịu tải quá lớn trong khi các shard khác nhàn rỗi). 

- **Thách thức kỹ thuật:** Hệ thống mấtđi khả năng hỗ trợ các transaction ACID trên toàn cục (Cross-shard Transactions). Việc thực hiện các phép kết hợp `JOIN` dữ liệu giữa các máy chủ khác nhau là cực kỳ tốn 

Tài liệu kỹ thuật - tối ưu hóa PostgreSQL & Backend System 

Trang 9 / 12 

kém và phức tạp. Thông thường, các logic này phải được xử lý một cách thủ côngởtầngứng dụng hoặc thông qua các giải pháp middleware trung gian chuyên dụng như Citus Data cho PostgreSQL. 

## **9. tối ưu HÓA THIẾT KẾ SCHEMA DỮ LIỆU (SCHEMA DESIGN OPTIMIZATION)** 

_Một thiết kế Schema tồi là nguồn gốc rễ của mọi vấn đề suy giảm hiệu năng không thể cứu vãn bằng Index._ 

## **9.1 Lựa chọn kiểu dữ liệu (Datatypes) chính xác** 

- **Dữ liệu tiền tệ, số thập phân chính xác:** Tuyệt đối không dùng các kiểu số thực dấu phẩy động như `FLOAT` hoặc `REAL` do rủi ro sai số làm tròn trong tính toán kế toán. Phải sử dụng kiểu dữ liệu `NUMERIC(precision, scale)` hoặc `DECIMAL` . 

- **Khóa chính (Primary Key):** Tránh việc sử dụng chuỗi mã hóa ngẫu nhiên hoàn toàn như `UUIDv4` làm khóa chính. Do tính ngẫu nhiên, khi chèn bản ghi mới, cây index B-tree sẽ bị phân mảnh liên tục, ép hệ thống phải dịch chuyển các page dữ liệu trên đĩa liên tục. Thay vào đó, hãy sử dụng kiểu tăng tuần tự `BIGSERIAL` (8-byte) hoặc cấu trúc mã định danh thế hệ mới sắp xếp theo thời gian như `UUIDv7` . 

- • **Dữ liệu thời gian:** Luônưu tiên dùng kiểu `TIMESTAMPTZ` thay vì `TIMESTAMP` để đảm bảo hệ thống lưu trữ đồng nhất thông tin múi giờ quốc tế UTC, tránh các lỗi logic khi hệ thống mở rộngđa quốc gia. 

## **9.2 Đánh đổi giữa Chuẩn hóa (Normalization) và Phi chuẩn hóa (Denormalization)** 

Trong thiết kế cơ sở dữ liệu truyền thống, việc tuân thủ các dạng chuẩn (1NF, 2NF, 3NF) là bắt buộc để triệt tiêu sự dư thừa dữ liệu. Tuy nhiên, trong các hệ thống tải cao, việc liên tục thực hiện các câu lệnh liên kết (JOIN) 5-6 bảng để hiển thị một thông tin cơ bản sẽ kéo sụt tốc độ phản hồi. 

**Giải pháp Denormalization:** Chủ động lưu trữ dư thừa thông tin có tính toán. Ví dụ, thay vì mỗi lần hiển thị bài viết phải chạy lệnh `COUNT(*) FROM comments WHERE post_id = x` , ta thêm trực tiếp một cột `comment_count` vào bảng `posts` . Kỹ sư chịu trách nhiệm đồng bộ trường số lượng này thông qua DB Trigger hoặc xử lý bất đồng bộ ở tầng Application Logic. 

## **10. tối ưu HÓAỞMỨC KIẾN TRÚC HỆ THỐNG TỔNG THỂ** 

_Khi hệ thống phát triển đến quy mô cực lớn, việcép PostgreSQL xử lý cả các tác vụ không thuộc thế mạnh của nó (như tìm kiếm toàn văn phức tạp) là một sai lầm về mặt kiến trúc._ 

## **10.1 Áp dụng mô hình CQRS kết hợp công nghệ CDC (Change Data Capture)** 

Hệ thống thực hiện tách biệt hoàn toàn cơ sở dữ liệu chuyên trách cho việc thay đổi dữ liệu (Ghi - Commands) và cơ sở dữ liệu tối ưu cho việc tìm kiếm, kết xuất dữ liệu (Đọc - Queries). 

Tài liệu kỹ thuật - tối ưu hóa PostgreSQL & Backend System 

Trang 10 / 12 

```
 Luồng kiến trúc tổng thể:
 [ Client Write ] ──> [ Backend ] ──> [ PostgreSQL (Primary) ]
                                               │
                                       (Ghi nhận nhật ký WAL)
                                               ▼
                                      [ Debezium / Kafka CDC ]
                                               │
                                       (Đẩy luồng sự kiện)
                                               ▼
 [ Client Read ]  <── [ Backend ] <── [ Elasticsearch Engine ]
```

Trong mô hình này, PostgreSQL đóng vai trò là "Nguồn chân lý" (Source of Truth) lưu trữ dữ liệu cốt lõi dưới dạng chuẩn hóa. Công cụ CDC (như Debezium) liên tục theo dõi file log WAL của PostgreSQL, bất đồng bộ đẩy các sự kiện thay đổi dữ liệu qua hệ thống hàng đợi tin nhắn Kafka để đồng bộ sang **Elasticsearch** hoặc **OpenSearch** . Toàn bộ các tác vụ tìm kiếm tự động, tìm kiếm mờ (Fuzzy Search), lọcđa tiêu chí (Faceted Search) sẽ được xử lý trên Elasticsearch, giải phóng 100% gánh nặng xử lý tìm kiếm cho PostgreSQL. 

## **11. PHẦN BỔ SUNG: CÁC YẾU TỐ tối ưu HÓA HỆ QUẢN TRỊ CỐT LÕI** 

_Đây là các khía cạnh kỹ thuật chuyên sâu về cấu hình và bảo trì hệ thống PostgreSQL được bổ sung nhằm đảm bảo cẩm nang đạt tính toàn diện cao nhất._ 

## **11.1 tối ưu hóa các tham số cấu hình PostgreSQL (PostgreSQL Tuning)** 

Các thông số cấu hình mặc định khi cài đặt PostgreSQL thường được thiết lậpởmức an toàn cho các máy chủ cấu  hình  thấp.  Để chạy  production  hiệu  năng  cao,  cầnđiều  chỉnh  lại  các  tham  số trong  file `postgresql.conf` dựa trên tài nguyên phần cứng vật lý: 

- `shared_buffers` : Lượng bộ nhớ RAM được PostgreSQL sử dụng làm bộ nhớ đệm để đọc và ghi dữ 

- liệu. Giá trị khuyến nghị đối với máy chủ chuyên dụng là khoảng **25% tổng dung lượng RAM** của hệ thống. 

- `work_mem` : Bộ nhớ cấp cho mỗi tác vụ sắp xếp nội bộ (như `ORDER BY` , `DISTINCT` ) và các phép toán 

- Join trước khi phải ghi đĩa tạm. Cần tính toán cẩn thận dựa trên số kết nối đồng thời nhằm tránh tình trạng cạn kiệt RAM đột ngột. 

- `maintenance_work_mem` : Bộ nhớ dành riêng cho các tác vụ quản trị, bảo trì hệ thống lớn như khởi tạo 

- Index ( `CREATE INDEX` ), dọn dẹp rác bảng ( `VACUUM` ). Nên đặt giá trị lớn hơn nhiều so với `work_mem` . 

- • `effective_cache_size` : Tuyên bố cho PostgreSQL biết dung lượng bộ nhớ khả dụng cho việc lưu cache file của Hệ điều hành. Thường cấu hìnhởmức **50% - 75% tổng RAM** nhằm giúp Optimizer có xu hướng chọn Index Scan nhiều hơn thay vì Seq Scan. 

Tài liệu kỹ thuật - tối ưu hóa PostgreSQL & Backend System 

Trang 11 / 12 

## **11.2 Cơ chế MVCC và Công tác bảo trì định kỳ tránh hiện tượng Bloat** 

PostgreSQL sử dụng cơ chế kiểm soátđa phiên bản (MVCC - Multi-Version Concurrency Control). Khi một dòng dữ liệu bị `DELETE` hoặc `UPDATE` , dòng cũ không bị xóa trực tiếp vật lý trên đĩa ngay lập tức mà chỉ được đánh dấu thành một "Dead Tuple" (Bản ghi chết) để đảm bảo các transaction khácđang chạy song song không bị ảnh hưởng dữ liệu. 

**Hiện tượng Phình bảng (Table Bloat):** Nếu không có biện pháp dọn dẹp, các Dead Tuple này sẽ tích lũy ngày qua ngày, chiếm dụng không gian đĩa vật lý bừa bãi và làm chậm các câu lệnh quét dữ liệu. Do đó, việc cấu hình tiến trình ngầm **Auto-Vacuum** hoạt động hiệu quả là bắt buộc để định kỳ thu hồi không gian và cập nhật lại thống kê phân phối dữ liệu ( `ANALYZE` ) giúp bộ lập lịch chọn đúng kế hoạch thực thi tối ưu nhất. 

## **11.3 Giám sát hiệu năng truy vấn qua pg_stat_statements** 

Để tối ưu hóa một hệ thốngđang vận hành, kỹ sư không thể đoán mò. Việc kích hoạt extension hệ thống `pg_stat_statements` là bắt buộc để thu thập số liệu thực tế: 

```
-- Kích hoạt trong postgresql.conf
shared_preload_libraries = 'pg_stat_statements'
-- Câu lệnh tìm ra top 5 câu truy vấn tiêu tốn nhiều thời gian xử lý nhất của hệ
thống:
SELECT query, calls, total_exec_time, mean_exec_time, rows
FROM pg_stat_statements
ORDER BY total_exec_time DESC
LIMIT 5;
```

Thông qua bảng số liệu thu được từ view này, kỹ sư sẽ định vị chính xác câu lệnh SQL nàođang chạy nhiều nhất hoặc tốn thời gian trung bình lâu nhất, từ đó tập trung nguồn lực tối ưu hóa mang lại hiệu quả cao nhất cho toàn bộ hệ thống backend. 

Tài liệu kỹ thuật - tối ưu hóa PostgreSQL & Backend System 

Trang 12 / 12 

