## **CẨM NANG TOÀN DIỆN VỀ LOCK TRONG POSTGRESQL** 

_Phân loại khóa, cơ chế xung đột, phân tích sự cố Deadlock, Race Condition và chiến lược tối ưu hóa hệ thống_ 

## **1. TỔNG QUAN VỀ CƠ CHẾ KHÓA (LOCKING) TRONG POSTGRESQL** 

Trong các hệ quản trị cơ sở dữ liệu quan hệ (RDBMS) nói chung và PostgreSQL nói riêng, cơ chế khóa (locking) đóng vai trò quyết định trong việc bảo đảm tính toàn vẹn dữ liệu và thực thi các thuộc tính **ACID** (Atomicity, Consistency, Isolation, Durability), đặc biệt là tính Cô lập (Isolation) khi có nhiều giao dịch (transactions) đồng thời cùng thao tác trên một vùng dữ liệu. 

PostgreSQL sử dụng kiến trúcđiều khiển đồng thờiđa phiên bản ( **MVCC - Multi-Version Concurrency Control** ). Nhờ MVCC, các thao tác đọc dữ liệu (đọc phiên bản cũ của dữ liệu) nói chung không bị chặn bởi các thao tác ghi dữ liệu (tạo phiên bản mới), và ngược lại ("readers don't block writers, and writers don't block readers"). Tuy nhiên, đối với một số thao tác thay đổi cấu trúc dữ liệu, cập nhật trùng hàng, hoặc khiứng dụng yêu cầu tính nhất quán tuyệt đối, PostgreSQL vẫn bắt buộc phải sử dụng các loại Lock để đồng bộ hóa. 

PostgreSQL chia Lock thành 4 tầng quản lý chính: 

- **Table-level locks (Khóa cấp bảng):** Áp dụng trên toàn bộ cấu trúc của một bảng dữ liệu. 

- **Row-level locks (Khóa cấp dòng):** Áp dụng trên từng bản ghi (hàng) cụ thể bên trong bảng. 

- 

- **Page-level locks (Khóa cấp trang):** Áp dụngởmức độ vật lý trên các data page trong ổ đĩa (thường mang tính nội bộ của hệ thống Engine). 

- **Advisory locks (Khóa khuyến nghị):** Do lập trình viên tự định nghĩa và quản lý thông qua logicứng dụng. 

- 

## **2. CHI TIẾT CÁC LOẠI TABLE-LEVEL LOCKS (KHÓA CẤP BẢNG)** 

Mặc dù tên gọi là khóa cấp bảng, các khóa này không nhất thiết phải khóa toàn bộ bảng khiến các giao dịch khác không thể truy cập. Ý nghĩa thực sự của chúng là khai báo mức độ tương thích của giao dịch hiện tại đối với bảng đó. PostgreSQL hỗ trợ 8 chế độ khóa cấp bảng với mức độ nghiêm ngặt tăng dần: 

1. **ACCESS SHARE:** Thường được kích thọc tự động bởi câu lệnh `SELECT` . Chế độ này chỉ xung đột với chế độ khóa nghiêm ngặt nhất là `ACCESS EXCLUSIVE` . 

2. **ROW SHARE:** Kích hoạt bởi câu lệnh `SELECT FOR SHARE` hoặc `SELECT FOR UPDATE` . Chế độ này cho thấy giao dịch muốn thao tác trên một số dòng cụ thể của bảng. 

3. **ROW EXCLUSIVE:** Kích hoạt bởi các lệnh thay đổi dữ liệu như `UPDATE` , `DELETE` , và `INSERT` . Nó cho phép nhiều giao dịch cùng sửa đổi các dòng khác nhau trên cùng một bảng một cách đồng thời. 

Cẩm nang PostgreSQL Lock & Xử lý Sự cố 

1 

4. **SHARE UPDATE EXCLUSIVE:** Áp dụng cho các lệnh bảo trì hệ thống như `VACUUM` (thông thường), `ANALYZE` , `CREATE INDEX CONCURRENTLY` , và `ALTER TABLE VALIDATE CONSTRAINT` . Khóa này tự xung đột với chính nó nhằm ngăn chặn hai tiến trình bảo trì chạy đồng thời trên một bảng. 

5. **SHARE:** Kích hoạt bởi lệnh `CREATE INDEX` (thông thường). Cho phép nhiều tiến trình cùng đọc dữ liệu nhưng ngăn chặn hoàn toàn việc ghi/chỉnh sửa dữ liệu (xung đột với `ROW EXCLUSIVE` ). 

6. **SHARE ROW EXCLUSIVE:** Thường được sử dụng bởi một số lệnh `ALTER TABLE` hoặc các thao tác nội bộ liên quan tới ràng buộc toàn vẹn. Xung đột với tất cả trừ `ACCESS SHARE` và `ROW SHARE` . 

7. **EXCLUSIVE:** Kích hoạt bởi lệnh `REFRESH MATERIALIZED VIEW CONCURRENTLY` . Cho phép đọc dữ liệu đồng thời ( `ACCESS SHARE` ) nhưng chặn đứng tất cả các tiến trình ghi ghi dữ liệu ( `ROW EXCLUSIVE` ). 

8. **ACCESS EXCLUSIVE:** Mức khóa tối cao, chặn hoàn toàn mọi truy cập khác (kể cả lệnh đọc `SELECT` đơn giản). Kích hoạt bởi các lệnh làm thay đổi cấu trúc vật lý bảng như: `DROP TABLE` , `TRUNCATE` , `VACUUM FULL REINDEX` , `CLUSTER` , và hầu hết các biến thể của `ALTER TABLE` . 

## **Ma trận tương thích giữa các chế độ khóa cấp bảng (Lock Compatibility Matrix)** 

Dưới đây là bảng tra cứu tính xung đột giữa các loại khóa cấp bảng. Dấu **[X]** biểu thị hai chế độ khóa **xung đột/ chặn lẫn nhau** (nếu tiến trình A giữ khóa này thì tiến trình B phải chờ và ngược lại). Ô trống biểu thị sự tương thích song song. 

|**Chế độ khóa**<br>**ứng dụng yêu**<br>**cầu**|**ACCESS**<br>**SHARE**|**ROW**<br>**SHARE**|**ROW**<br>**EXCLUSIVE**|**SHARE**<br>**UPDATE**<br>**EXCL**|**SHARE**|**SHARE**<br>**ROW**<br>**EXCL**|**EXCLUSIVE**|**ACCESS**<br>**EXCLUSIVE**|
|---|---|---|---|---|---|---|---|---|
|**ACCESS**<br>**SHARE**||||||||**[X]**|
|**ROW SHARE**|||||||**[X]**|**[X]**|
|**ROW**<br>**EXCLUSIVE**|||||**[X]**|**[X]**|**[X]**|**[X]**|
|**SHARE**<br>**UPDATE**<br>**EXCL**||||**[X]**|**[X]**|**[X]**|**[X]**|**[X]**|
|**SHARE**|||**[X]**|**[X]**||**[X]**|**[X]**|**[X]**|
|**SHARE ROW**<br>**EXCL**|||**[X]**|**[X]**|**[X]**|**[X]**|**[X]**|**[X]**|
|**EXCLUSIVE**||**[X]**|**[X]**|**[X]**|**[X]**|**[X]**|**[X]**|**[X]**|
|**ACCESS**<br>**EXCLUSIVE**|**[X]**|**[X]**|**[X]**|**[X]**|**[X]**|**[X]**|**[X]**|**[X]**|



Cẩm nang PostgreSQL Lock & Xử lý Sự cố 

2 

## **3. CHI TIẾT CÁC LOẠI ROW-LEVEL LOCKS (KHÓA CẤP DÒNG)** 

Khóa cấp dòng tự động được kích hoạt khi có sự thay đổi dữ liệu hoặc được gọi một cách tường minh qua các cú pháp khóa. PostgreSQL không lưu trữ thông tin Row lock trong bộ nhớ RAM (như các hệ quản trị dữ liệu khác để tránh tràn bộ nhớ), mà ghi trực tiếp thông tin khóa vào tiêu đề (header) của tuple dữ liệu trên đĩa. PostgreSQL hỗ trợ 4 chế độ Row lock: 

• **FOR UPDATE:** Chế độ khóa dòng nghiêm ngặt nhất. Nó ngăn chặn các giao dịch khác khóa dòng này bằng bất kỳ cơ chế nào (chặn hoàn toàn `FOR UPDATE` , `FOR NO KEY UPDATE` , `FOR SHARE` , `FOR KEY SHARE` ) và chặn mọi lệnh sửa đổi dòng ( `UPDATE/DELETE` ). 

- **FOR NO KEY UPDATE:** Tương tự như `FOR UPDATE` , nhưng nó có độ ưu tiên thấp hơn, cho phép các tiến trình yêu cầu `FOR KEY SHARE` đi qua. Lệnh `UPDATE` thông thường trong PostgreSQL không thay đổi các cột thuộc khóa chính/khóa duy nhất (Foreign Key/Primary Key) sẽ tự động kích hoạt loại khóa này. 

- • **FOR SHARE:** Tương tự như khóa đọc chung (Shared Lock). Nhiều giao dịch cùng lúc có thể giữ khóa `FOR SHARE` trên một hàng. Nó ngăn chặn các tiến trình muốn thay đổi hàng (chặn lệnh `UPDATE/DELETE` và khóa `FOR UPDATE` ). 

• **FOR KEY SHARE:** Mức khóa dòng nhẹ nhất. Nó được hệ thống tự động kích hoạt khi kiểm tra tính toàn vẹn của khóa ngoại (Foreign Key constraint). Cho phép các giao dịch khác chỉnh sửa các trường dữ liệu bình thường của hàng đó miễn là không chạm vào các trường Khóa (Key). 

## **4. SỰ CỐ DEADLOCK (KHÓA CHẾT)** 

## **4.1. Bản chất và Cơ chế xảy ra** 

**Deadlock** xảy ra khi hai hoặc nhiều giao dịch đồng thời đang nắm giữ các khóa trên các tài nguyên khác nhau, và mỗi giao dịch lại yêu cầu một khóa bổ sung trên tài nguyên mà giao dịch kia đang nắm giữ. Điều này tạo ra một vòng lặp phụ thuộc khép kín (circular dependency), khiến tất cả các giao dịch liên quan rơi vào trạng thái chờ đợi vô thời hạn nếu hệ thống không can thiệp. 

_Kịch bản kinh điển dẫn đến Deadlock giữa 2 giao dịch:_ 

- **Thờiđiểm T1:** Giao dịch A thực hiện cập nhật Hàng số 1 trong bảng dữ liệu (Giữ khóa độc quyền trên Hàng 1). 

- • **Thờiđiểm T2:** Giao dịch B thực hiện cập nhật Hàng số 2 trong bảng dữ liệu (Giữ khóa độc quyền trên Hàng 2). • **Thờiđiểm T3:** Giao dịch A cố gắng cập nhật Hàng số 2 (Giao dịch A bị block và phải chờ Giao dịch B nhả khóa). 

- **Thờiđiểm T4:** Giao dịch B cố gắng cập nhật Hàng số 1 (Giao dịch B bị block và phải chờ Giao dịch A nhả khóa). 

Tại thờiđiểm T4, đồ thị phụ thuộc của hệ thống xuất hiện một vòng lặp vô hạn: `Giao dịch A → Đ ợ i B → Đ ợ i A` . 

Cẩm nang PostgreSQL Lock & Xử lý Sự cố 

3 

## **4.2. Cách nhận biết (Detection)** 

Khác với các sự cố hiệu năng thông thường, PostgreSQL sở hữu một cơ chế phát hiện Deadlock nội bộ rất mạnh mẽ. Khi một câu lệnh bị block, PostgreSQL sẽ kích hoạt bộ đếm thời gian dựa trên tham số `deadlock_timeout` (mặc định là 1 giây). 

1. **Cảnh báo trong Log file:** Nếu hết thời gian cấu hình mà khóa chưa được giải phóng, PostgreSQL sẽ tính toán đồ thị phụ thuộc khóa (Lock Dependency Graph). Nếu phát hiện vòng lặp, hệ thống sẽ chủ động hy sinh một trong  các  giao  dịch  bằng  cách  hủy  bỏ (abort)  nó  và  bắn  lỗi  với  mã  trạng  thái  SQLSTATE: `40P01 (deadlock_detected)` . 

2. **Dấu hiệu trongỨng dụng:** Tầngứng dụng (Backend code) nhận lại một Exception/Error với nội dung tương tự: `ERROR: deadlock detected. DETAIL: Process 12345 waits for ShareLock on transaction 98765; blocked by process 54321...` . 

## **4.3. Cách phòng tránh và khắc phục** 

Để triệt tiêu khả năng xảy ra Deadlock, thiết kế hệ thống phải tuân thủ nghiêm ngặt các nguyên tắc sau: 

- **Nhất quán về thứ tự thao tác:** Đây là giải pháp cốt lõi nhất. Đảm bảo tất cả các hàm, API và tiến trình trongứng dụng luôn thực hiện cập nhật dữ liệu trên các bảng/các hàng theo một thứ tự cố định (Ví dụ: Luôn cập nhật Bảng Tài khoản trước, Bảng Lịch sử giao dịch sau; hoặc sắp xếp ID của mảng các dòng cần cập nhật theo thứ tự tăng dần trước khi thực hiện vòng lặp cập nhật trong SQL). 

- **Thu hẹp phạm vi Giao dịch:** Giữ cho các giao dịch luôn ngắn nhất có thể. Không đặt các tác vụ xử lý tính toán nặng, tương tác I/O mạng (gọi API bên thứ ba) ở bên trong khối dữ liệu `BEGIN...COMMIT` nhằm giải phóng khóa nhanh nhất có thể. 

• **Sử dụng cơ chế khóa không chờ (Non-blocking Lock):** Khiáp dụng khóa tường minh thông qua câu lệnh SQL, nênđi kèm từ khóa `NOWAIT` hoặc `SKIP LOCKED` . Ví dụ: `SELECT * FROM invoice WHERE status = 'PENDING' FOR UPDATE NOWAIT;` . Nếu hàng đã bị khóa, câu lệnh sẽ lập tức ném ra lỗi thay vì treo máy chờ đợi, giúpứng dụng có thể bắt được ngoại lệ và xử lý logic thử lại (retry logic) một cách chủ động. 

- **Cấu hình thời gian chờ tốiđa (Lock Timeout):** Đặt giá trị cho tham số `lock_timeout` . Tham số này quy định thời gian tốiđa một câu lệnh SQL được phép đứng chờ khóa. Nếu quá thời gian, lệnh sẽ bị hủy mà không cần đợi đến khi thuật toán quét deadlock chạy. 

**Chiến lược xử lý ở tầngỨng dụng:** Vì bản chất phân tán và đồng thời, việc phòng tránh tuyệt đối Deadlock trong một hệ thống tải cao cực kỳ khó khăn. Ứng dụng bắt buộc phải thiết lập **Cơ chế Thử lại (Retry Mechanism)** . Khi nhận được mã lỗi hệ thống là `40P01` , ứng dụng nên rollback giao dịch hiện tại hoàn toàn, đợi một khoảng thời gian ngắn (ngẫu nhiên - exponential backoff) rồi thực hiện lại toàn bộ giao dịch đó. 

Cẩm nang PostgreSQL Lock & Xử lý Sự cố 

4 

## **5. SỰ CỐ RACE CONDITION (TÌNH TRẠNGĐUA TRANH DỮ LIỆU)** 

## **5.1. Bản chất và Cơ chế xảy ra** 

**Race Condition** trong cơ sở dữ liệu xảy ra khi kết quả cuối cùng của một chuỗi các thao tác dữ liệu phụ thuộc hoàn toàn vào tiến trình thực thi, tốc độ, hoặc thứ tự đan xen (interleaving) ngẫu nhiên của các luồng (threads/processes) đồng thời. Nó thường xuất hiện dưới mô hình phản mẫu mang tên **Read-Modify-Write** (Đọc dữ liệu ra → Tính toán ở tầng ứng dụng → Ghi kết quả ngược lại DB). 

_Ví dụ thực tế về Race Condition (Sự cố rút tiền):_ 

Giả sử Tài khoản Xđang có số dư là _**1.000.000đ**_ . Có hai yêu cầu rút tiền _**600.000đ**_ diễn ra cùng một mili-giây. 

- **Luồng 1 (Yêu cầu 1):** Thực hiện `SELECT balance FROM accounts WHERE id = X;` → Hệ thống trả về kết quả _**1.000.000đ**_ . 

- **Luồng 2 (Yêu cầu 2):** Đồng thời thực hiện câu lệnh đọc tương tự `SELECT` → Nhận lại kết quả _**1.000.000đ**_ vì Luồng 1 chưa cập nhật dữ liệu. 

- **Luồng 1:** Kiểm traởtầng Code thấy _**1.000.000đ ≥ 600.000đ**_ (Hợp lệ). Tiến hành trừ tiền trên Ram ra _**400.000đ**_ và gọi lệnh cập nhật: `UPDATE accounts SET balance = 400000 WHERE id = X;` . 

- **Luồng 2:** Kiểm traởtầng Code riêng cũng thấy _**1.000.000đ ≥ 600.000đ**_ (Hợp lệ). Tiến hành trừ tiền trên Ram ra _**400.000đ**_ và gọi lệnh cập nhật: `UPDATE accounts SET balance = 400000 WHERE id = X;` . 

**Hậu quả:** Khách hàng rút tổng cộng _**1.200.000đ**_ thành công, nhưng số dư tài khoản trong DB cuối cùng vẫn là _**400.000đ**_ thay vì phảiâm hoặc báo lỗi thiếu tiền. Đây gọi là hiện tượng **Lost Update (Mất mát dữ liệu cập nhật)** . 

## **5.2. Mối liên hệ với các mức độ Cô lập Giao dịch (Transaction Isolation Levels)** 

PostgreSQL cung cấp 3 mức độ cô lập cấu hình giúp kiểm soát Race Condition: 

1. **Read Committed (Mặc định):** Mỗi câu lệnh trong một giao dịch chỉ nhìn thấy dữ liệu đã được commit trước khi _câu lệnh đó_ bắt đầu. Mức này hoàn toàn không thể chống lại Race Condition theo mô hình Read-ModifyWrite. 

2. **Repeatable Read:** Giao dịch chỉ nhìn thấy dữ liệu đã được commit trước khi _toàn bộ giao dịch_ bắt đầu. Nếu Luồng 2 cố gắng cập nhật một hàng dữ liệu mà Luồng 1 đã thay đổi và commit trước đó, PostgreSQL sẽ chặn đứng Luồng 2 và ném ra lỗi: `ERROR: could not serialize access due to concurrent update (SQLSTATE: 40001)` . Ứng dụng phải bắt lỗi này và thực hiện retry. 

3. **Serializable:** Mức độ cô lập nghiêm ngặt nhất. Giả lập môi trường như thể tất cả các giao dịch được thực thi tuần tự từng cái một. Nó sử dụng cơ chế Khóa cô lập (SSI - Serializable Snapshot Isolation). Nếu phát hiện bất kỳ dấu hiệu xung đột phi tuần tự nào, hệ thống lập tức hủy giao dịch lỗi với mã lỗi `40001` . 

Cẩm nang PostgreSQL Lock & Xử lý Sự cố 

5 

## **5.3. Cách nhận biết, phòng tránh và khắc phục** 

Để nhận biết Race Condition, chúng ta không thể dựa vào log lỗi của DB (trừ khi dùng mức cô lập cao), mà phải dựa vào sự bất nhất về mặt logic kinh doanh dữ liệu (ví dụ: số dư âm, lượng hàng tồn kho bị sai lệchâm, thống kê không khớp). 

Có 3 phương pháp chính để phòng tránh và giải quyết triệt để Race Condition: 

## **Giải pháp 1: Thực hiện tính toán Nguyên tử (Atomic Operations) trực tiếp trong SQL** 

Loại bỏ hoàn toàn bước trung gian đưa dữ liệu lên ứng dụng rồi tính toán. Tận dụng tối đa tính nguyên tử của bản thân câu lệnh SQL. 

```
-- Thay vì Read-Modify-Write, hãy gộp vào một câu lệnh duy nhất vớiđiều kiện chặt chẽ
UPDATE accounts
SET balance = balance - 600000
WHERE id = 'X' AND balance >= 600000;
```

Hệ thống sẽ thực hiện khóa dòng này, kiểm trađiều kiện `balance ≥ 600000` trực tiếp tại tầng đĩa cứng, đảm bảo an toàn tuyệt đối. 

## **Giải pháp 2: Cơ chế khóa Bi quan (Pessimistic Locking)** 

Chủ động chiếm quyền kiểm soát độc quyền dòng dữ liệu ngay từ bước đọc bằng cú pháp `FOR UPDATE` . Điều này buộc tất cả các luồng đồng thời khác khi gọi lệnh đọc dòng đó đều phải xếp hàng chờ đợi. 

## `BEGIN;` 

```
-- Khóa dòng này lại ngay lập tức, các tiến trình khác gọi đến dòng này phải xếp hàng
chờ
```

```
SELECT balance FROM accounts WHERE id = 'X' FOR UPDATE;
```

```
-- Thực hiện kiểm tra logic tại tầngỨng dụng một cách an toàn
```

```
-- Nếu hợp lệ, tiến hành cập nhật dữ liệu
UPDATE accounts SET balance = new_calculated_balance WHERE id = 'X';
COMMIT;
```

## **Giải pháp 3: Cơ chế khóa Lạc quan (Optimistic Locking)** 

Giải pháp này không dùng khóa của DB mà dùng một cột đánh dấu phiên bản (ví dụ: cột `version` kiểu số nguyên hoặc cột `last_updated` kiểu timestamp) để kiểm tra xem dữ liệu có bị thay đổi trong quá trìnhứng dụng xử lý hay không. 

Cẩm nang PostgreSQL Lock & Xử lý Sự cố 

6 

```
-- Bước 1: Đọc dữ liệu cùng số version hiện tại
SELECT balance, version FROM accounts WHERE id = 'X';
-- Giả sử lấy ra balance = 1.000.000, version = 5
-- Bước 2: Thực hiện tính toán ở code ứng dụng và cập nhật kèm điều kiện version
UPDATE accounts
SET balance = 400000, version = version + 1
WHERE id = 'X' AND version = 5;
```

Nếu Luồng khác đã cập nhật trước đó, số `version` trong DB đã tăng lên thành 6. Câu lệnh `UPDATE` của luồng hiện tại sẽ trả về số lượng hàng bị tác động là `0` . Ứng dụng sẽ nhận biết việc cập nhật thất bại, tiến hành rollback và thực hiện lại tiến trình đọc. 

## **6. CÁC CÂU LỆNH SQL HỮUÍCH ĐỂ GIÁM SÁT VÀ XỬ LÝ KHÓA** 

Trong quá trình vận hành hệ thống, khi xảy ra tình trạng nghẽn mạch (hệ thống chạy chậm, đứng kết nối), người quản trị cơ sở dữ liệu (DBA) hoặc Lập trình viên cần nhanh chóng tìm ra nguồn gốc gây nghẽn. 

## **6.1. Truy vấn tìm các Giao dịchđang bị Block (Bị treo) và Giao dịch gây Block** 

Đoạn script sau sử dụng các view hệ thống `pg_stat_activity` và `pg_locks` để chỉ ra chi tiết mã tiến trình (PID) nàođang chặn tiến trình nào, cùng nội dung câu lệnh SQL cụ thể: 

Cẩm nang PostgreSQL Lock & Xử lý Sự cố 

7 

## `SELECT` 

```
    blocked_locks.pid     AS blocked_pid,
    blocked_activity.usename  AS blocked_user,
    blocking_locks.pid    AS blocking_pid,
    blocking_activity.usename AS blocking_user,
    blocked_activity.query    AS blocked_statement,
    blocking_activity.query   AS current_statement_in_blocking_process
FROM  pg_catalog.pg_locks         blocked_locks
JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid =
blocked_locks.pid
JOIN pg_catalog.pg_locks         blocking_locks
    ON blocking_locks.locktype = blocked_locks.locktype
    AND blocking_locks.database IS NOT DISTINCT FROM blocked_locks.database
    AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
    AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
    AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
    AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
    AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
    AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
    AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
    AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
    AND blocking_locks.pid != blocked_locks.pid
```

```
JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid =
blocking_locks.pid
```

```
WHERE NOT blocked_locks.granted;
```

## **6.2. Các câu lệnh can thiệp khẩn cấp để giải phóng Lock** 

Khi phát hiện một tiến trình (ví dụ có `PID = 54321` ) đang giữ khóa độc quyền quá lâu (do transaction bị treo, kết nối mạng từ client bị ngắt nhưng chưa giải phóng session), ta dùng các lệnh sau để can thiệp trực tiếp: 

- **Hủy bỏ câu lệnh hiện tại của Session (Giữ lại kết nối):** 

## • 

```
SELECT pg_cancel_backend(54321);
```

Lệnh này sẽ gửi một tín hiệu SIGINT đến tiến trình con, dừng câu lệnh SQLđang chạy của session đó nhưng không ngắt kết nối mạng của client. 

## **Ngắt hoàn toàn phiên làm việc (Giải phóng lập tức mọi khóa):** 

## • 

```
SELECT pg_terminate_backend(54321);
```

Lệnh này gửi tín hiệu SIGTERM, giết chết tiến trình con quản lý session này. Toàn bộ các khóađang nắm giữ bởi session này sẽ được giải phóng ngay lập tức. Đây là giải pháp triệt để nhất khi xử lý sự cố nghẽn khóa nghiêm trọng. 

Cẩm nang PostgreSQL Lock & Xử lý Sự cố 

8 

## **6.3. Cấu hình tham số hệ thống tối ưu để tự động phòng vệ** 

Để ngăn chặn việc hệ thống bị treo cứng diện rộng do một câu lệnh không tối ưu gây ra, nên cấu hình các tham số sau trong file `postgresql.conf` hoặc setởmức Session củaứng dụng: 

```
-- 1. Thoi gian toi da cho mot cau lenh dung cho Lock (Vi du: 5 giay)
-- Neu vuot qua 5s ma chua lay duoc lock, lenh se tu dong huy bo de bao ve he thong
SET lock_timeout = '5s';
-- 2. Thoi gian phat hien Deadlock (Vi du: 1 giay)
-- Quy dinh thoi gian PostgreSQL bat dau quet do thi lock de phat hien va xu ly
deadlock
SET deadlock_timeout = '1s';
```

```
-- 3. Thoi gian toi da de thuc thi mot cau lenh SQL bat ky (Vi du: 30 giay)
-- Ngan chan cac cau lenh SELECT hoac UPDATE chay qua lau làm can kiet tai nguyen
SET statement_timeout = '30s';
```

## **7. TỔNG KẾT QUY TRÌNH CHUẨNỨNG PHÓ SỰ CỐ VỀ LOCK** 

1. **Chủ động thiết lập cấu hình phòng vệ:** Cấu hình `lock_timeout` hợp lý trên môi trường Production để hệ thống tự động ngắt các truy vấn bị nghẽn, tránh hiệuứng Domino kéo sập toàn bộ cụm cơ sở dữ liệu. 

2. **Xây dựng cơ chế Logs và Alerting:** Theo dõi tần suất xuất hiện của mã lỗi `40P01` (Deadlock) và `40001` (Serialization Failure) trong log củaứng dụng để đưa ra các biện pháp tái cấu trúc code kịp thời. 

3. **Tối ưu hóa code ứng dụng:** Loại bỏ hoàn toàn tư duy thiết kế "Read-Modify-Write" không an toàn. Thay thế bằng các câu lệnh tính toán nguyên tử (Atomic updates) hoặc áp dụng cơ chế Khóa bi quan/Khóa lạc quan tùy thuộc vào mức độ xung đột dữ liệu của hệ thống. 

Cẩm nang PostgreSQL Lock & Xử lý Sự cố 

9 
