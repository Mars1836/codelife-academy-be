## **NGUYÊN LÝ SOLID & DESIGN PATTERNS** 

## **Hướng Dẫn Ứng Dụng Trong Ngôn Ngữ Go (Golang)** 

**Giới thiệu chung:** Trong phát triển phần mềm, việc viết mã nguồn chạy được mới chỉ là bước đầu tiên. Mã nguồn tốt đòi hỏi phải dễ đọc, dễ bảo trì, dễ kiểm thử và linh hoạt trước các thay đổi. Tài liệu này cung cấp cái nhìn chi tiết về 5 nguyên lý thiết kế hệ thống **SOLID** và các **Design Patterns** (Mẫu thiết kế) phổ biến nhất, đi kèm với cách triển khai đặc trưng bằng ngôn ngữ **Go** tận dụng cơ chế cấu trúc (struct) và giao diện (interface) tường minh. 

## **Phần 1: 5 Nguyên Lý Thiết Kế SOLID** 

SOLID là tập hợp 5 nguyên lý thiết kế hướng đối tượng được định hình bởi Robert C. Martin (Uncle Bob). Mặc dù Go không phải là ngôn ngữ hướng đối tượng truyền thống (không có tính kế thừa lớp), các nguyên lý này vẫn áp dụng hoàn hảo thông qua cơ chế `composition` (cấu thành) và `interfaces`. 

## **1. Single Responsibility Principle (SRP) - Nguyên lý đơn trách nhiệm** 

_Mỗi struct/module chỉ nên đảm nhiệm một trách nhiệm duy nhất và chỉ có một lý do duy nhất để thay đổi._ 

Việc gom nhiều logic khác nhau (như vừa xử lý dữ liệu vừa ghi log hoặc gửi email) vào một cấu trúc dữ liệu sẽ tạo ra sự ràng buộc chặt chẽ, gây khó khăn cho quá trình bảo trì và mở rộng. 

H ướ ng dẫn Nguyên lý SOLID & Design Patterns trong Go 

1 

```
package main
import "fmt"
// SAI: Struct đảm nhiệm cả lưu thông tin và inấn định dạng
type BadUser struct {
    Name  string
    Email string
}
func (u *BadUser) FormatJSON() string {
    return fmt.Sprintf("{\"name\": \"%s\", \"email\": \"%s\"}", u.Name, u.Email)
}
// ĐÚNG: Tách biệt trách nhiệm cấu trúc dữ liệu và xử lý định dạng
type User struct {
    Name  string
    Email string
}
type UserOutputFormatter struct{}
func (f *UserOutputFormatter) ToJSON(u *User) string {
    return fmt.Sprintf("{\"name\": \"%s\", \"email\": \"%s\"}", u.Name, u.Email)
}
```

## **2. Open/Closed Principle (OCP) - Nguyên lý M ở /Đóng** 

_Mã nguồn nên mở rộng cho việc phát triển thêm tính năng mới, nhưng đóng trước việc sửa đổi mã nguồn hiện tại._ 

Trong Go, ta thực hiện nguyên lý này thông qua `interface` . Khi muốn thêm hành vi, ta tạo struct mới hiện thực hóa interface thay vì sửa đổi hàm hoặc struct gốc bằng các câu lệnh rẽ nhánh `if-else` hoặc `switch-case` dài dòng. 

H ướ ng dẫn Nguyên lý SOLID & Design Patterns trong Go 

2 

```
package main
// Định nghĩa một interface chung cho việc tính toán hình học
type Shape interface {
    Area() float64
}
// Hình chữ nhật kế thừa hành vi của Shape
type Rectangle struct {
    Width, Height float64
}
func (r Rectangle) Area() float64 { return r.Width * r.Height }
// Khi muốn thêm Hình tròn, ta CHỈ CẦN thêm struct mới mà không cần sửa mã nguồn
cũ
type Circle struct {
    Radius float64
}
func (c Circle) Area() float64 { return 3.1415926535 * c.Radius * c.Radius }
// Hàm tính tổng diện tích đóng trước sự thay đổi khi thêm hình mới
func AreaCalculator(shapes []Shape) float64 {
    var totalArea float64
    for _, shape := range shapes {
        totalArea += shape.Area()
    }
    return totalArea
}
```

## **3. Liskov Substitution Principle (LSP) - Nguyên lý thay thếLiskov** 

_Các đối tượng thuộc lớp con phải có khả năng thay thế hoàn toàn cho đối tượng thuộc lớp cha mà không làm thay đổi tính đúng đắn của chương trình._ 

Vì Go sử dụng hệ thống kiểu dữ liệu ngầm định (duck typing) thông qua interface, LSP có nghĩa là mọi struct triển khai một interface phải tuân thủ nghiêm ngặt hợp đồng (contract) và hành vi kỳ vọng mà interface đó quy định, không được gây lỗi hoặc ném ra biệt lệ bất thường. 

H ướ ng dẫn Nguyên lý SOLID & Design Patterns trong Go 

3 

```
package main
type Document interface {
    Read() string
}
type EditableDocument interface {
    Document
    Write(content string)
}
// Định nghĩa một tài liệu chỉ đọc
type ReadOnlyFile struct {
    Content string
}
func (ro ReadOnlyFile) Read() string { return ro.Content }
// Nếu ta cố tìnhép ReadOnlyFile thực hiện một phương thức Write rỗng hoặc lỗi,
ta vi phạm LSP.
```

```
// Giải pháp: Tách nhỏ interface như trên để đảm bảo tính đúng đắn.
```

## **4. Interface Segregation Principle (ISP) - Nguyên lý phân tách Interface** 

_Thà thiết kế nhiều interface nhỏ, tập trung vào mục đích cụ thể còn hơn thiết kế một interface lớn, đa năng nhưngép client phụ thuộc vào các phương thức họ không sử dụng._ 

Đây là triết lý cốt lõi của Golang. Các interface tiêu chuẩn trong Go thường rất nhỏ (ví dụ: `io.Reader` chỉ chứa một phương thức duy nhất là `Read` ). 

H ướ ng dẫn Nguyên lý SOLID & Design Patterns trong Go 

4 

```
package main
// SAI: Interface quá lớn, ép buộc các đối tượng triển khai cả hai hành vi
type Worker interface {
    Work()
    Eat()
}
// ĐÚNG: Phân tách thành hai interface nhỏ gọn chuyên biệt
type SimpleWorker interface {
    Work()
}
type Eater interface {
    Eat()
}
// Struct Robot chỉ cần thực hiện SimpleWorker
type Robot struct{}
func (r Robot) Work() {} // Robot hoạt động tốt, không cần phải hiện thực hàm
Eat() vô nghĩa
```

## **5. Dependency Inversion Principle (DIP) - Nguyên lý đ ả o ng ượ c phụthuộc** 

_Các module cấp cao không nên phụ thuộc trực tiếp vào các module cấp thấp. Cả hai nên phụ thuộc vào sự trừu tượng (Abstraction)._ 

Nói cách khác, các thành phần xử lý logic nghiệp vụ chính (Business Logic) nên giao tiếp với tầng cơ sở dữ liệu hoặc dịch vụ bên ngoài thông qua interface, chứ không phụ thuộc vào một thư viện hoặc driver cụ thể. 

H ướ ng dẫn Nguyên lý SOLID & Design Patterns trong Go 

5 

```
package main
// Trừu tượng hóa tầng lưu trữ dữ liệu
type DBConnection interface {
    Query(sql string) string
}
// Module cấp thấp triển khai trừu tượng
type MySQL struct{}
func (db MySQL) Query(sql string) string { return "MySQL Data" }
// Module cấp cao phụ thuộc hoàn toàn vào trừu tượng
type UserServices struct {
    db DBConnection
}
func NewUserService(db DBConnection) *UserServices {
    return &UserServices{db: db} // Đóng gói thông qua Dependency Injection
}
```

## **Phần 2: Các Mẫu Thiết K ế (Design Patterns) Cần Thiết & Th ườ ng Dùng** 

Design Patterns là các giải pháp đã được chuẩn hóa để giải quyết các vấn đề phổ biến trong kiến trúc phần mềm. Chúng được chia làm 3 nhóm chính: Khởi tạo (Creational), Cấu trúc (Structural), và Hành vi (Behavioral). 

## **1. Singleton Pattern (Nhóm Khởi tạo)** 

**Mục đích:** Đảm bảo một cấu trúc dữ liệu (struct) chỉ có duy nhất một thực thể (instance) trong suốt vòng đời củaứng dụng và cung cấp mộtđiểm truy cập toàn cục tới nó. Trong Go, Singleton thường được triển khai an toàn trong môi trườngđa luồng (concurrency) bằng cách sử dụng gói `sync.Once` . 

H ướ ng dẫn Nguyên lý SOLID & Design Patterns trong Go 

6 

```
package main
import (
    "fmt"
    "sync"
)
type database struct {
    connectionString string
}
var (
    instance *database
    once     sync.Once
)
// GetDatabaseInstance đảm bảo việc khởi tạo thread-safe độc bản
func GetDatabaseInstance() *database {
    once.Do(func() {
        instance = &database{
            connectionString: "db://root:secret@localhost:5432/main",
        }
    })
    return instance
}
```

## **2. Factory Method Pattern (Nhóm Khởi tạo)** 

**Mục đích:** Cung cấp một giao diện chung để tạo các đối tượng mà không cần chỉ định chính xác lớp/struct cụ thể nào sẽ được khởi tạo. Việc quyết định khởi tạo kiểu đối tượng nào sẽ phụ thuộc vào tham số truyền vào đầu hàm. 

H ướ ng dẫn Nguyên lý SOLID & Design Patterns trong Go 

7 

```
package main
import "errors"
type PaymentMethod interface {
    Pay(amount float64) string
}
type CreditCard struct{}
func (cc *CreditCard) Pay(amount float64) string { return "Thanh toán bằng Thẻ
Tín Dụng" }
```

```
type PayPal struct{}
func (p *PayPal) Pay(amount float64) string { return "Thanh toán bằng PayPal" }
// Factory Function
func GetPaymentMethod(methodType string) (PaymentMethod, error) {
    switch methodType {
    case "creditcard":
        return &CreditCard{}, nil
    case "paypal":
        return &PayPal{}, nil
    default:
        return nil, errors.New("phương thức không hợp lệ")
    }
}
```

## **3. Adapter Pattern (Nhóm Cấu trúc)** 

**Mục đích:** Đóng vai trò là cầu nối trung gian, cho phép hai interface không tương thích có thể làm việc chung với nhau mà không cần sửa đổi mã nguồn gốc của cả hai bên. 

H ướ ng dẫn Nguyên lý SOLID & Design Patterns trong Go 

8 

```
package main
```

```
// Giao diện hiện tại hệ thốngđang sử dụng
type ModernPrinter interface {
    PrintModern() string
}
// Hệ thống cũ (Legacy) có giao diện khác biệt hoàn toàn
type LegacyPrinter struct{}
func (lp *LegacyPrinter) PrintOldWay() string { return "Dữ liệu từ máy in cũ" }
// PrinterAdapter bao bọc đối tượng cũ để tương thích với giao diện mới
type PrinterAdapter struct {
    legacyPrinter *LegacyPrinter
}
func (adapter *PrinterAdapter) PrintModern() string {
    return adapter.legacyPrinter.PrintOldWay()
}
```

## **4. Decorator Pattern (Nhóm Cấu trúc)** 

**Mục đích:** Cho phép bổ sung thêm các hành vi hoặc trách nhiệm mới vào một đối tượng một cách động mà không làmảnh hưởng đến cấu trúc bên trong đối tượng gốc hay phá vỡ các lớp khác. 

```
package main
type Coffee interface {
    GetCost() int
}
type SimpleCoffee struct{}
func (c *SimpleCoffee) GetCost() int { return 20000 }
```

```
// MilkDecorator thêm tính năng và giá tiền mà không đổi cấu trúc SimpleCoffee
type MilkDecorator struct {
    coffee Coffee
}
func (m *MilkDecorator) GetCost() int {
    return m.coffee.GetCost() + 5000 // Thêm giá trị của sữa vào ly cà phê
}
```

H ướ ng dẫn Nguyên lý SOLID & Design Patterns trong Go 

9 

## **5. Strategy Pattern (Nhóm Hành vi)** 

**Mục đích:** Định nghĩa một tập hợp các thuật toán tương tự nhau, bao gói từng thuật toán lại và giúp chúng có thể hoán đổi linh hoạt cho nhau ngay trong quá trình thực thi chương trình (runtime). 

```
package main
// Chiến lược định giá giảm giá tổng thể
type DiscountStrategy interface {
    ApplyDiscount(price float64) float64
}
type NormalCustomer struct{}
func (n *NormalCustomer) ApplyDiscount(price float64) float64 { return price }
type VIPCustomer struct{}
func (v *VIPCustomer) ApplyDiscount(price float64) float64 { return price *
0.8 } // Giảm giá 20%
// Bộ thanh toán sử dụng chiến lược linh động
type Checkout struct {
    strategy DiscountStrategy
}
func (c *Checkout) SetStrategy(s DiscountStrategy) { c.strategy = s }
func (c *Checkout) Calculate(amount float64) float64 { return
c.strategy.ApplyDiscount(amount) }
```

## **6. Observer Pattern (Nhóm Hành vi)** 

**Mục đích:** Xây dựng mối quan hệ phụ thuộc "một-nhiều". Khi trạng thái của một đối tượng (Subject) thay đổi, tất cả các thành phần đăng ký theo dõi nó (Observers) sẽ nhận được thông báo tự động. 

H ướ ng dẫn Nguyên lý SOLID & Design Patterns trong Go 

10 

```
package main
import "fmt"
type Observer interface {
    Update(message string)
}
type CustomerObserver struct{ Name string }
func (co *CustomerObserver) Update(msg string) {
    fmt.Printf("Khách hàng %s nhận thông báo: %s\n", co.Name, msg)
}
type ItemSubject struct {
    observers []Observer
    name      string
}
func (item *ItemSubject) Register(o Observer) { item.observers =
append(item.observers, o) }
func (item *ItemSubject) NotifyAll(msg string) {
    for _, observer := range item.observers {
        observer.Update(msg)
    }
}
```

## **Kết luận** 

Việcáp dụng nhuần nhuyễn **SOLID** và các **Design Patterns** trong ngôn ngữ Go yêu cầu tư duy lập trình tập trung vào Interfaces và Composition thay vì Kế thừa. Hãy bắt đầu từ những cấu trúc đơn giản nhất, tránh việc lạm dụng quá sớm (Over-engineering) để mã nguồn của bạn vừa sạch, vừa tối ưu, vừa đúng bản chất tối giản đặc trưng của Golang. 

H ướ ng dẫn Nguyên lý SOLID & Design Patterns trong Go 

11 
