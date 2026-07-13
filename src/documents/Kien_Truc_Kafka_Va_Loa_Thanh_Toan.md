## **HỆ THỐNG KAFKA & XỬ LÝ IDEMPOTENCY CHO E-COMMERCE & LOA THANH TOÁN** 

## **Phần 1: Các Khái Niệm Cơ Bản Về Apache Kafka** 

Kafka là một nền tảng Event Streaming phân tán, được thiết kế để xử lý lượng lớn dữ liệu theo thời gian thực. Để xây dựng hệ thống bền bỉ, chúng ta cần nắm vững các khái niệm cốt lõi: 

- **Broker:** Là một server Kafka (node) lưu trữ dữ liệu. Một cụm Kafka (Cluster) bao gồm nhiều Brokers. 

- **Topic:** Là một kênh (channel) hoặc danh mục để phân loại message. Trong hệ thống thanh toán, ta có thể có topic `payment.completed.events` . 

- **Partition:** Một Topic được chia thành nhiều Partitions. Partition cho phép Topic scale ngang. Dữ liệu trong partition được sắp xếp theo thứ tự và có một ID tuần tự gọi là _offset_ . 

- **Producer:** Ứng dụng đẩy (publish) message vào Topic (ví dụ: NestJS PayboxService sau khi nhận webhook từ Paybox). 

- **Consumer:** Ứng dụng kéo (subscribe) dữ liệu từ Topic về để xử lý. 

- **Consumer Group:** Một nhóm các consumers cùng đọc chung một Topic. Mỗi message trong một partition chỉ được xử lý bởi _duy nhất một consumer_ trong một consumer group. Điều này đảm bảo tính toán song song mà không bị trùng lặp message. 

## **1.1. Khi Nào Nên Scale Kafka?** 

Việc scale hệ thống cần được thực hiện khi gặp các tình trạng sau: 

- **Consumer Lag tăng cao:** Tốc độ Producer đẩy message vào nhanh hơn tốc độ Consumer xử lý. Dẫn đến độ trễ hệ thống (ví dụ: người dùng thanh toán xong nhưng 5 phút sau loa mới báo hoặc đơn hàng mới cập nhật). 

- **Tăng lượng giao dịch đột biến (Spike Traffic):** Lượng throughput giao dịch thanh toán tạo ra quá lớn, một broker hoặc một consumer không thể gánh nổi CPU/Memory. 

## **1.2. Cách Scale Kafka Hợp Lý** 

Để scale hệ thống Kafka, quy tắc vàng là: **Tăng số lượng Partition của Topic, sau đó tăng số lượng Consumer trong Consumer Group.** 

## **Quy tắc Ràng buộc giữa Partition và Consumer:** 

Số lượng Consumer đang chạy trong một Group **không bao giờ được lớn hơn** số lượng Partition của Topic đó. 

- Topic có 3 Partitions + 3 Consumers → Lý tưởng (Mỗi consumer đọc 1 partition). 
- Topic có 3 Partitions + 4 Consumers → 1 Consumer sẽ bị rảnh rỗi (Idle), gây lãng phí resource. 
- Topic có 3 Partitions + 2 Consumers → 1 Consumer sẽ phải đọc 2 partitions. 

## **Lệnh mở rộng partition (Command Line):** 

```bash
kafka-topics.sh --alter --bootstrap-server localhost:9092 \
  --topic payment.completed.events --partitions 6
```

Sau khi chạy lệnh trên, có thể scale ứng dụng Consumer lên tối đa 6 instances (ví dụ dùng Docker / K8s) để tăng gấp đôi tốc độ xử lý. 

## **Phần 2: Áp Dụng Cho Kịch Bản E-commerce & Loa Thanh Toán** 

Dựa trên thiết kế kiến trúc Paybox và quy trình ở hình ảnh, hệ thống đang phải xử lý giao dịch và thông báo tới loa (terminal). Khi áp dụng Event-Driven Architecture với Kafka, quy trình sẽ được tách bạch (decouple) như sau: 

## **2.1. Quy trình thực tế (Tích hợp Kafka)** 

1. **Tạo đơn và QR:** Khách chọn thanh toán, Client gọi API tạo Order, hệ thống gọi Paybox ( `createCustomQr` ) lấy `billNumber` và mã QR trả về Client. 
2. **Khách hàng quét mã:** Thanh toán chuyển khoản thành công, Paybox đẩy callback về hệ thống qua endpoint `POST /paybox/payment/callback/tomotek` . 
3. **Backend xử lý (Producer):** 
   - Tại `handleTransactionCallback` , hệ thống validate dữ liệu, dùng Pessimistic Lock tìm entity qua `billNumber` và update trạng thái `PaymentStatus.PAID` . 
   - Sau khi lưu DB thành công, thay vì gọi API loa trực tiếp, hệ thống publish một event (ví dụ `OrderPaidEvent` ) vào Kafka topic `payment-success-events` . 
4. **Các Consumer Groups cùng lắng nghe sự kiện:** 
   - **Group A (ecommerce-order-processor):** Lắng nghe event để cập nhật logic ecommerce (cập nhật tồn kho, gửi email/SMS cho khách, invalidate Redis cache). 
   - **Group B (iot-speaker-notifier):** Lắng nghe event và gửi tín hiệu (qua MQTT/WebSocket) xuống đúng terminalID (Loa) tại shop. Loa nhận lệnh và phát âm thanh _"Nhận được {amount} nghìn đồng"_ . 

_Ưu điểm:_ Không để việc gọi Loa (vốn dễ lỗi kết nối mạng) làm chậm quá trình callback từ ngân hàng. Giao dịch vẫn update tức thời vào DB, và các Consumer tự xử lý tác vụ của mình theo tốc độ riêng (asynchronous). 

## **Phần 3: Xử lý Idempotency (Từ Frontend → NestJS → Kafka)** 

Idempotency (Tính lũy đẳng) đảm bảo rằng một thao tác có thể được gọi nhiều lần nhưng hệ thống chỉ xử lý (như trừ tiền, thay đổi trạng thái) ở lần đầu tiên. Điều này cực kỳ quan trọng trong thanh toán E-commerce. 

## **3.1. Tại Frontend** 

Khi User submit form "Tạo đơn hàng / Thanh toán", Frontend sinh ra một `Idempotency-Key` (thường là UUID v4) gắn vào Header của Request. Đồng thời disable nút bấm sau khi click để giảm tải phía UI. 

```javascript
// Gửi request kèm Idempotency-Key
const idempotencyKey = crypto.randomUUID();
await axios.post('/api/orders/pay', data, {
    headers: { 'X-Idempotency-Key': idempotencyKey }
});
``` 

## **3.2. Tại Backend NestJS (Interceptor / Middleware)** 

Sử dụng Redis để khóa (lock) theo `Idempotency-Key` . Nếu key đã tồn tại, ta trả về kết quả đã xử lý trước đó thay vì chạy lại logic tạo QR hay tạo Order. 

```typescript
import { Injectable, NestInterceptor, ExecutionContext, CallHandler, HttpException, HttpStatus } from '@nestjs/common';
import { RedisService } from '@/services/redis/redis.service';
import { of } from 'rxjs';
import { tap } from 'rxjs/operators';

@Injectable()
export class IdempotencyInterceptor implements NestInterceptor {
  constructor(private readonly redisService: RedisService) {}
  
  async intercept(context: ExecutionContext, next: CallHandler) {
    const request = context.switchToHttp().getRequest();
    const key = request.headers['x-idempotency-key'];
    if (!key) return next.handle();

    const cacheKey = `idempotency:${key}`;
    const cachedResponse = await this.redisService.get(cacheKey);
    if (cachedResponse) {
      return of(JSON.parse(cachedResponse)); // Trả kết quả cũ nếu duplicate
    }
    
    // Lock key bằng SETNX (tránh race condition)
    const isSet = await this.redisService.setNX(cacheKey, 'PROCESSING', 60);
    if (!isSet) throw new HttpException('Request is processing', HttpStatus.CONFLICT);
    
    return next.handle().pipe(
      tap(async (response) => {
        // Lưu response thành công
        await this.redisService.set(cacheKey, JSON.stringify(response), 86400);
      }),
    );
  }
}
```

## **3.3. Tại Webhook Callback (Database Idempotency)** 

Phía đối tác (Paybox/Ngân hàng) có thể gọi callback 2, 3 lần nếu có lỗi network. Code hiện tại của bạn đã xử lý phần này rất tốt bằng cách dùng `pessimistic_write` và kiểm tra trạng thái cũ: 

```typescript
// 1. Pessimistic Lock Record trên PostgreSQL
const order = await manager.findOne(Order, {
  where: { billNumber: payload.billNumber },
  lock: { mode: 'pessimistic_write' },
});

// 2. Idempotency Check: Đã PAID rồi thì return success ngay
if (order.paymentStatus === PaymentStatus.PAID) {
  this.logger.log(`Order ${order.id} already processed — idempotent success`);
  return { resCode: '00', resDesc: 'Success' };
}
```

## **3.4. Tại Kafka Producer & Consumer** 

**Producer (NestJS):** Cần bật tính năng `enable.idempotence = true` . Tránh trường hợp Producer gửi message thành công lên Broker, nhưng Broker phản hồi ACK chậm do mạng, khiến Producer retry gửi message thứ 2. 

```typescript
// Cấu hình Producer KafkaJS trong NestJS
{
  clientId: 'tomotek-ecommerce',
  brokers: ['localhost:9092'],
  producer: {
    idempotent: true, // Chống duplicate từ phía Producer
    maxInFlightRequests: 5,
    retries: 5
  }
}
```

**Consumer (xử lý Loa):** Kafka đảm bảo at-least-once delivery, nên consumer báo Loa vẫn có tỷ lệ rất nhỏ nhận message 2 lần. Cần thêm bảng DB lưu `eventId` (ví dụ: `processed_events`) hoặc dùng Redis kiểm tra. 

```typescript
async handlePaymentSuccess(payload: PaymentEventPayload) {
  const isProcessed = await this.redisService.setNX(`event:${payload.billNumber}`, 'PROCESSED', 86400);
  if (!isProcessed) {
    this.logger.warn(`Event ${payload.billNumber} đã được thông báo ra loa, bỏ qua.`);
    return; // Đảm bảo Idempotent, Loa không đọc 2 lần
  }

  // Đẩy lệnh xuống loa qua websocket
  await this.speakerService.notifyTerminal(payload.terminalID, payload.transactionAmount);
}
```
