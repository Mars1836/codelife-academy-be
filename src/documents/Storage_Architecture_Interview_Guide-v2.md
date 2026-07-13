## **KIẾN TRÚC STORAGE & OBJECT STORE** 

_Phân tích chuyên sâu về MinIO, SeaweedFS, RustFS và Câu hỏi Phỏng vấn (Cập nhật Pre-signed URL)_ 

## **1. Tổ ề ng quan v Distributed Object Storage** 

Trong các hệ thống phân tán hiện đại, Object Storage (như Amazon S3, MinIO) được sử dụng rộng rãi để lưu trữ dữ liệu phi cấu trúc vì khả năng mở rộng vô hạn (horizontal scaling), quản lý metadata linh hoạt, và giao tiếp qua HTTP/REST API thay vì các giao thức file system truyền thống như NFS hay SMB. 

## **2. Cách hoạt động của các Hệ thống Storage** 

## **2.1. MinIO** 

MinIO là một Object Storage server hiệu năng cao, mã nguồn mở, tương thích hoàn toàn với API của Amazon S3. 

- **Kiến trúc Single-Binary:** MinIO không có kiến trúc Master-Slave. Mọi node trong cụm đều bình đẳng. 

- **Erasure Coding:** Thay vì nhân bản dữ liệu, MinIO chia một object thành các khối dữ liệu và khối chẵn lẻ (parity blocks), giúp khôi phục dữ liệu ngay cả khi mất tới phân nửa số đĩa. 

- **Bitrot Protection:** Phát hiện và tự động sửa các lỗi suy thoái dữ liệuâm thầmởcấp độ phần cứng bằng Hashing. 

## **2.2. SeaweedFS** 

SeaweedFS là hệ thống được thiết kế đặc biệt để giải quyết **vấn đề file nhỏ (Small File Problem)** . 

- **Tách biệt Metadata và Dữ liệu:** Gồm Master Server, Volume Server và Filer. 

- **Xử lý File nhỏ (O(1) Disk Read):** Gom hàng triệu file nhỏ vào một file khổng lồ gọi là Volume. Metadata (offset, size) được lưu trên RAM, cho phép đọc đĩa với đúng 1 thao tác I/O. 

## **2.3. RustFS** 

RustFS là Object Storage thế hệ mới viết bằng Rust, tối ưu hiệu năng và tài nguyên. 

- **Không Garbage Collection:** Nhờ Rust, hệ thống không bị độ trễ do dọn rác như các hệ thống viết bằng Go (MinIO), giúp p99 latency cực kỳ thấp và ổn định. 

- **tối ưu I/O:** Xử lý payload nhỏ rất tốt. 

## **3. Câu Hỏi Phỏng Vấn (Phần 1: Kiến trúc cốt lõi)** 

**Câu 1: Tại sao chúng ta dùng Object Storage thay vì Block Storage hay File Storage để lưuảnh user?** 

File Storage gặp giới hạn về số lượng node/thư mục, càng nhiều file thì duyệt metadata càng chậm. Block Storage đắt đỏ và thiếu linh hoạt để chia sẻ qua mạng. Object Storage là không gian phẳng, scale vô hạn, truy cập qua HTTP/REST API, chi phí cực rẻ và dễ dàng tích hợp trực tiếp với CDN. 

## **Câu 2: Cơ chế Erasure Coding trong MinIO là gì? Tại sao không dùng Replication?** 

Replication (x3) tốn 300% dung lượng. Erasure Coding (EC) cắt file thành N khối dữ liệu và M khối mã hóa. EC chỉ tốn khoảng 1.5x dung lượng đĩa để đạt được mức chịu lỗi tương đương. EC đánh đổi CPU để tiết kiệm dung lượng lưu trữ khổng lồ. 

**Câu 3: Hệ thống có 1 tỷ ảnh avatar rất nhỏ. MinIO truyền thống gặp vấn đề gì và SeaweedFS giải quyết ra sao?** 

1 tỷ file nhỏ sẽ làm cạn kiệt bộ nhớ Inode của Linux và làm chậm đĩa cơ học do random seek. SeaweedFS giải quyết bằng cách gom chúng vào các Volume lớn, đưa cấu trúc metadata của file con (Offset, Size) lên RAM. Khi cần đọc, đĩa cứng chỉ cần di chuyển kim từ từ đến đúng 1 offset duy nhất (O(1) IO). 

## **4. Câu Hỏi Phỏng Vấn (Phần 2: Pre-signed URL & Bảo mật)** 

**Câu 4: Pre-signed URL trong Object Storage là gì? Nêu các kịch bản bắt buộc phải sử dụng nó trong System Design?** 

Pre-signed URL là một đường link HTTP(s) chứa sẵn các tham số xác thực được ký bằng thuật toán mã hóa. Nó cho phép một người dùng bất kỳ thực hiện một hành động cụ thể (như tải lên hoặc tải xuống một file nhất định) trong một thời gian ngắn giới hạn mà **không cần phải biết Access Key / Secret Key** của hệ thống. 

**Kịch bản sử dụng cốt lõi (Direct Upload Pattern):** Khi user muốn upload video nặng 500MB, nếu upload qua Backend Server thì Server sẽ bị nghẽn băng thông và RAM. Thay vào đó, Backend tạo một Pre-signed URL cấp quyền PUT và trả về cho Client. Client dùng URL đó đẩy video **trực tiếp** lên MinIO/S3. Điều này giúp hệ thống Backend không bị quá tải. 

**Câu 5: Thuật toán và cơ chế đằng sau Pre-signed URL (AWS Signature v4) hoạt động như thế nào? Chữ ký được tạo ra sao?** 

Bản chất của thuật toán là sử dụng **HMAC-SHA256** (Hash-based Message Authentication Code). Quá trình diễn ra hoàn toànởBackend (nơi giữ Secret Key): 

- **Bước 1 - Tạo StringToSign:** Gom tất cả các thông tin của request lại thành một chuỗi duy nhất: HTTP Method (GET/PUT), đường dẫn file (Bucket + Object Key), Query Parameters (đặc biệt là `X-Amz-Expires` - thời gian sống của link) và Timestamp hiện tại. 

- **Bước 2 - Tính toán Signing Key:** Không dùng trực tiếp Secret Key, hệ thống dùng Secret Key để HMAC tuần tự với các chuỗi: Ngày tháng (Date) -> Region -> Service (s3). Kết quả tạo ra một Signing Key dùng cho ngày hôm đó. 

- **Bước 3 - Tạo Signature:** Dùng Signing Key vừa tạo kết hợp thuật toán HMAC-SHA256 để băm chuỗi `StringToSign` ở Bước 1. Output sẽ là một chuỗi hex gọi là Signature. 

Chuỗi Signature này sẽ được đính kèm vào URL. Khi Storage Server nhận được request, nó tự lấy Secret Key nội bộ làm lại quy trình trên. Nếu 2 Signature khớp nhau -> Hợp lệ. 

**Câu 6: Nếu Hacker bắt được gói tin chứa Pre-signed URL và cố tình chỉnh sửa tham số** **`X-Amz-Expires` để link sống lâu hơn, hoặc đổi tên file tải lên, chuyện gì sẽ xảy ra?** 

Hệ thống Storage sẽ từ chối request và trả về lỗi **SignatureDoesNotMatch (HTTP 403)** . 

Lý do: Chữ ký (Signature) được sinh ra dựa trên thuật toán Hash một chiều có tính bao hàm (hash của HTTP Method, URI, Tham số). Bất kỳ sự thay đổi nhỏ nào (đổi tên file từ `avatar.jpg` sang `shell.php` , hay đổi thời gian sống từ 300s lên 3000s) đều khiến chuỗi StringToSign bị thay đổi. Khi Storage Server nhận URL bị sửa, nó tính toán lại Signature và thấy không khớp với Signature do Backend bạn ký ban đầu. Đây là cốt lõi của tính **Toàn vẹn (Integrity)** trong quá trình tạo chữ ký. 

**Câu 7: Khi dùng Pre-signed URL để cấp quyền cho User upload file, làm sao để chặn User đẩy file quá lớn (ví dụ chống đẩy file 50GB làm đầy đĩa) hoặc sai định dạng?** 

Với một link PUT URL đơn giản, ta **không thể** giới hạn kích thước hoặc định dạng trực tiếp trên URL. Đây là một lỗ hổng rất hay bị hỏi trong phỏng vấn. 

**Cách giải quyết:** Phải sử dụng cơ chế **Pre-signed POST Policy** thay vì PUT thông thường. 

POST Policy là một tài liệu JSON được Backend mã hóa Base64 và ký xác nhận. Trong Policy JSON này, ta có thể khai báo các quy tắc (Conditions) nghiêm ngặt mà Client phải tuân theo khi upload: 

- `["content-length-range", 0, 5242880]` : Giới hạn file từ 0 đến 5MB. 

- `["starts-with", "$key", "user-uploads/123/"]` : Bắt buộc file phải được upload vào đúng 

- thư mục của user có ID 123. 

- `["eq", "$Content-Type", "image/png"]` : Bắt buộc định dạng là PNG. 

Nếu Client không tuân thủ bất kỳ rules nào trong Policy đã ký, Storage Server sẽ ngay lập tức hủy kết nối. 

