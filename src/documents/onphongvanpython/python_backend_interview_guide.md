# Python Backend Interview Guide

Tài liệu này giúp bạn ôn tập các kiến thức Python Backend thường gặp khi phỏng vấn. Nội dung được viết theo hướng dễ hiểu, có ví dụ minh họa, câu hỏi phỏng vấn và các lưu ý thực tế.

---

# Mục lục

1. [Kiểu dữ liệu: list, tuple, set, dict](#1-kiểu-dữ-liệu-list-tuple-set-dict)
2. [Mutable và immutable](#2-mutable-và-immutable)
3. [`*args` và `**kwargs`](#3-args-và-kwargs)
4. [List comprehension và generator](#4-list-comprehension-và-generator)
5. [Decorator](#5-decorator)
6. [Context manager và `with`](#6-context-manager-và-with)
7. [Exception handling](#7-exception-handling)
8. [OOP trong Python](#8-oop-trong-python)
9. [`classmethod`, `staticmethod`, `property`](#9-classmethod-staticmethod-property)
10. [Thread, process và async](#10-thread-process-và-async)
11. [GIL là gì?](#11-gil-là-gì)
12. [`async/await` phù hợp trường hợp nào?](#12-asyncawait-phù-hợp-trường-hợp-nào)
13. [Virtual environment và dependency](#13-virtual-environment-và-quản-lý-dependency)
14. [Câu hỏi phỏng vấn thường gặp](#14-câu-hỏi-phỏng-vấn-thường-gặp)
15. [Kiến thức chung về Python framework](#15-kiến-thức-chung-về-python-framework)
16. [Bài tập thực hành](#16-bài-tập-thực-hành)
17. [Checklist trước phỏng vấn](#17-checklist-trước-phỏng-vấn)

---

# 1. Kiểu dữ liệu: list, tuple, set, dict

## 1.1. List

`list` là kiểu dữ liệu dùng để lưu một danh sách phần tử có thứ tự.

```python
users = ["Hậu", "Nam", "Linh"]

print(users[0])
# Hậu
```

Đặc điểm:

- Có thứ tự.
- Cho phép phần tử trùng nhau.
- Có thể thay đổi sau khi tạo.
- Có thể chứa nhiều kiểu dữ liệu khác nhau.

```python
data = ["Hậu", 25, True, {"role": "backend"}]
```

Thêm phần tử:

```python
users = ["Hậu", "Nam"]

users.append("Linh")

print(users)
# ['Hậu', 'Nam', 'Linh']
```

Chèn phần tử vào vị trí cụ thể:

```python
users.insert(1, "Minh")

print(users)
```

Xóa phần tử:

```python
users.remove("Nam")
```

Xóa theo vị trí:

```python
removed_user = users.pop(0)
```

Duyệt list:

```python
for user in users:
    print(user)
```

Cắt list bằng slicing:

```python
numbers = [1, 2, 3, 4, 5]

print(numbers[1:4])
# [2, 3, 4]
```

Độ phức tạp thường gặp:

| Thao tác | Độ phức tạp trung bình |
|---|---:|
| Truy cập theo index | O(1) |
| Thêm cuối list | O(1) |
| Tìm kiếm phần tử | O(n) |
| Chèn đầu list | O(n) |
| Xóa phần tử giữa list | O(n) |

Không nên dùng `list` nếu cần liên tục kiểm tra một giá trị có tồn tại hay không trong một tập dữ liệu rất lớn. Trường hợp đó, `set` thường phù hợp hơn.

---

## 1.2. Tuple

`tuple` gần giống `list`, nhưng không thể thay đổi sau khi tạo.

```python
point = (10, 20)

print(point[0])
# 10
```

Không thể làm như sau:

```python
point[0] = 99
# TypeError
```

Tuple phù hợp khi:

- Dữ liệu không nên bị thay đổi.
- Cần dùng làm key của dictionary.
- Muốn biểu diễn một nhóm giá trị cố định.

Ví dụ:

```python
DATABASE_CONFIG = ("localhost", 5432, "nocobase")
```

Tuple một phần tử phải có dấu phẩy:

```python
value = (10,)

print(type(value))
# <class 'tuple'>
```

Nếu viết:

```python
value = (10)

print(type(value))
# <class 'int'>
```

---

## 1.3. Set

`set` là tập hợp không có thứ tự và không chứa phần tử trùng nhau.

```python
roles = {"admin", "user", "user"}

print(roles)
# {'admin', 'user'}
```

Set phù hợp khi:

- Loại bỏ phần tử trùng.
- Kiểm tra phần tử tồn tại nhanh.
- Thực hiện phép hợp, giao, hiệu.

Ví dụ loại bỏ phần tử trùng:

```python
emails = [
    "a@example.com",
    "b@example.com",
    "a@example.com",
]

unique_emails = set(emails)

print(unique_emails)
```

Kiểm tra tồn tại:

```python
allowed_roles = {"admin", "manager"}

if "admin" in allowed_roles:
    print("Có quyền truy cập")
```

Các phép toán tập hợp:

```python
backend_skills = {"Python", "PostgreSQL", "Docker"}
devops_skills = {"Docker", "Kubernetes", "Linux"}

print(backend_skills | devops_skills)
# Hợp

print(backend_skills & devops_skills)
# Giao

print(backend_skills - devops_skills)
# Hiệu
```

Lưu ý: phần tử trong `set` phải là kiểu hashable, thường là kiểu immutable.

Không thể thêm `list` vào `set`:

```python
my_set = {[1, 2, 3]}
# TypeError: unhashable type: 'list'
```

Có thể thêm `tuple` nếu mọi phần tử bên trong đều hashable:

```python
my_set = {(1, 2), (3, 4)}
```

---

## 1.4. Dict

`dict` lưu dữ liệu dạng key-value.

```python
user = {
    "id": 1,
    "name": "Hậu",
    "role": "backend",
}

print(user["name"])
# Hậu
```

Truy cập an toàn bằng `get`:

```python
print(user.get("email"))
# None
```

Có thể truyền giá trị mặc định:

```python
print(user.get("email", "Chưa có email"))
```

Cập nhật dữ liệu:

```python
user["name"] = "Vũ Công Hậu"
user["email"] = "hau@example.com"
```

Duyệt dictionary:

```python
for key, value in user.items():
    print(key, value)
```

Các phương thức thường dùng:

```python
print(user.keys())
print(user.values())
print(user.items())
```

Xóa phần tử:

```python
email = user.pop("email", None)
```

Key của dictionary phải là hashable.

Hợp lệ:

```python
data = {
    "name": "Hậu",
    1: "one",
    (10, 20): "position",
}
```

Không hợp lệ:

```python
data = {
    [1, 2]: "value"
}
# TypeError
```

---

## So sánh nhanh

| Kiểu | Có thứ tự | Cho phép trùng | Có thể thay đổi | Truy cập |
|---|---|---|---|---|
| `list` | Có | Có | Có | Index |
| `tuple` | Có | Có | Không | Index |
| `set` | Không đảm bảo | Không | Có | Không dùng index |
| `dict` | Có thứ tự chèn | Key không trùng | Có | Key |

---

# 2. Mutable và immutable

## 2.1. Mutable là gì?

Mutable nghĩa là đối tượng có thể bị thay đổi sau khi được tạo.

Một số kiểu mutable:

- `list`
- `dict`
- `set`
- Phần lớn object do người dùng tự tạo

Ví dụ:

```python
numbers = [1, 2, 3]

numbers.append(4)

print(numbers)
# [1, 2, 3, 4]
```

Đối tượng ban đầu đã bị thay đổi.

---

## 2.2. Immutable là gì?

Immutable nghĩa là không thể thay đổi đối tượng sau khi tạo.

Một số kiểu immutable:

- `int`
- `float`
- `bool`
- `str`
- `tuple`
- `frozenset`

Ví dụ:

```python
name = "Hậu"

name = name + " Backend"
```

Chuỗi ban đầu không bị chỉnh sửa. Python tạo một chuỗi mới và gán lại cho biến `name`.

---

## 2.3. Hiểu về tham chiếu

```python
a = [1, 2, 3]
b = a

b.append(4)

print(a)
# [1, 2, 3, 4]
```

`a` và `b` cùng trỏ đến một object.

Kiểm tra:

```python
print(a is b)
# True
```

Sao chép nông:

```python
a = [1, 2, 3]
b = a.copy()

b.append(4)

print(a)
# [1, 2, 3]

print(b)
# [1, 2, 3, 4]
```

---

## 2.4. Shallow copy và deep copy

Shallow copy chỉ sao chép object ngoài cùng.

```python
import copy

original = {
    "user": {
        "name": "Hậu"
    }
}

shallow = copy.copy(original)

shallow["user"]["name"] = "Nam"

print(original["user"]["name"])
# Nam
```

Object lồng nhau vẫn dùng chung tham chiếu.

Deep copy sao chép toàn bộ:

```python
import copy

original = {
    "user": {
        "name": "Hậu"
    }
}

deep = copy.deepcopy(original)

deep["user"]["name"] = "Nam"

print(original["user"]["name"])
# Hậu
```

---

## 2.5. Lỗi phổ biến với default argument

Không nên viết:

```python
def add_item(item, items=[]):
    items.append(item)
    return items
```

Gọi hàm:

```python
print(add_item(1))
# [1]

print(add_item(2))
# [1, 2]
```

List mặc định chỉ được tạo một lần khi Python định nghĩa hàm.

Cách đúng:

```python
def add_item(item, items=None):
    if items is None:
        items = []

    items.append(item)
    return items
```

Đây là câu hỏi phỏng vấn Python rất phổ biến.

---

# 3. `*args` và `**kwargs`

## 3.1. `*args`

`*args` cho phép hàm nhận nhiều positional arguments.

```python
def calculate_sum(*args):
    return sum(args)

print(calculate_sum(1, 2, 3, 4))
# 10
```

Bên trong hàm, `args` là một tuple.

```python
def show_args(*args):
    print(type(args))
    print(args)

show_args("Python", "Docker", "PostgreSQL")
```

Tên `args` chỉ là quy ước. Điều quan trọng là dấu `*`.

```python
def show_values(*values):
    print(values)
```

---

## 3.2. `**kwargs`

`**kwargs` cho phép hàm nhận nhiều keyword arguments.

```python
def create_user(**kwargs):
    print(kwargs)

create_user(
    name="Hậu",
    role="backend",
    age=25,
)
```

Bên trong hàm, `kwargs` là một dictionary.

---

## 3.3. Kết hợp tham số

Thứ tự thông thường:

```python
def example(normal_arg, *args, default_arg=True, **kwargs):
    pass
```

Ví dụ:

```python
def log_event(event_name, *tags, level="INFO", **metadata):
    print("Event:", event_name)
    print("Tags:", tags)
    print("Level:", level)
    print("Metadata:", metadata)


log_event(
    "user_created",
    "user",
    "crm",
    level="SUCCESS",
    user_id=100,
    created_by="admin",
)
```

---

## 3.4. Unpacking

Unpack list hoặc tuple:

```python
numbers = [1, 2, 3]

print(*numbers)
# 1 2 3
```

Truyền list vào hàm:

```python
def add(a, b, c):
    return a + b + c

numbers = [1, 2, 3]

print(add(*numbers))
```

Unpack dictionary:

```python
def create_user(name, role):
    return {
        "name": name,
        "role": role,
    }

data = {
    "name": "Hậu",
    "role": "backend",
}

user = create_user(**data)
```

Ứng dụng thực tế:

- Wrapper function.
- Decorator.
- Framework gọi handler.
- Truyền cấu hình động.
- Viết hàm linh hoạt.

---

# 4. List comprehension và generator

## 4.1. List comprehension

Cách thông thường:

```python
squares = []

for number in range(1, 6):
    squares.append(number ** 2)

print(squares)
```

Dùng list comprehension:

```python
squares = [number ** 2 for number in range(1, 6)]

print(squares)
```

Thêm điều kiện:

```python
even_squares = [
    number ** 2
    for number in range(1, 11)
    if number % 2 == 0
]
```

Có `if-else`:

```python
labels = [
    "even" if number % 2 == 0 else "odd"
    for number in range(1, 6)
]
```

Dictionary comprehension:

```python
users = ["Hậu", "Nam", "Linh"]

user_map = {
    index: name
    for index, name in enumerate(users, start=1)
}

print(user_map)
```

Set comprehension:

```python
unique_lengths = {
    len(name)
    for name in ["Hậu", "Nam", "Linh", "Minh"]
}
```

Không nên dùng comprehension quá phức tạp vì làm code khó đọc.

---

## 4.2. Generator expression

List comprehension tạo toàn bộ dữ liệu trong bộ nhớ:

```python
numbers = [number ** 2 for number in range(1_000_000)]
```

Generator expression tạo giá trị từng phần khi cần:

```python
numbers = (number ** 2 for number in range(1_000_000))
```

Khác biệt chính:

```python
list_data = [x * 2 for x in range(5)]
generator_data = (x * 2 for x in range(5))

print(list_data)
# [0, 2, 4, 6, 8]

print(generator_data)
# <generator object ...>
```

Lấy dữ liệu generator:

```python
for value in generator_data:
    print(value)
```

Hoặc:

```python
generator_data = (x * 2 for x in range(5))

print(next(generator_data))
print(next(generator_data))
```

Generator sẽ bị tiêu thụ sau khi duyệt.

```python
generator_data = (x for x in range(3))

print(list(generator_data))
# [0, 1, 2]

print(list(generator_data))
# []
```

---

## 4.3. Hàm generator với `yield`

```python
def count_up_to(limit):
    current = 1

    while current <= limit:
        yield current
        current += 1


for number in count_up_to(5):
    print(number)
```

Khi gặp `yield`, hàm tạm dừng và lưu trạng thái. Lần gọi tiếp theo, hàm tiếp tục chạy từ vị trí cũ.

Ứng dụng:

- Đọc file lớn theo từng dòng.
- Stream dữ liệu.
- Xử lý query nhiều record.
- Sinh dữ liệu vô hạn.
- Tránh đưa toàn bộ dữ liệu vào RAM.

Ví dụ đọc file lớn:

```python
def read_large_file(file_path):
    with open(file_path, "r", encoding="utf-8") as file:
        for line in file:
            yield line.strip()
```

---

## 4.4. So sánh bộ nhớ

```python
import sys

list_data = [x for x in range(1_000_000)]
generator_data = (x for x in range(1_000_000))

print(sys.getsizeof(list_data))
print(sys.getsizeof(generator_data))
```

Generator thường dùng ít bộ nhớ hơn nhiều vì không giữ toàn bộ phần tử cùng lúc.

Lưu ý: generator không phải lúc nào cũng nhanh hơn. Nó chủ yếu giúp giảm bộ nhớ và hỗ trợ lazy evaluation.

---

# 5. Decorator

Decorator là một hàm nhận vào một hàm khác và trả về một hàm mới.

## 5.1. Ví dụ cơ bản

```python
def my_decorator(func):
    def wrapper():
        print("Trước khi chạy hàm")
        func()
        print("Sau khi chạy hàm")

    return wrapper


@my_decorator
def say_hello():
    print("Xin chào")


say_hello()
```

Dòng:

```python
@my_decorator
```

tương đương:

```python
say_hello = my_decorator(say_hello)
```

---

## 5.2. Decorator hỗ trợ tham số bất kỳ

```python
from functools import wraps


def log_function(func):
    @wraps(func)
    def wrapper(*args, **kwargs):
        print(f"Đang gọi hàm: {func.__name__}")
        result = func(*args, **kwargs)
        print(f"Kết thúc hàm: {func.__name__}")
        return result

    return wrapper


@log_function
def add(a, b):
    return a + b


print(add(2, 3))
```

`@wraps(func)` giúp giữ lại:

- Tên hàm.
- Docstring.
- Metadata của hàm gốc.

---

## 5.3. Decorator có tham số

```python
from functools import wraps


def require_role(required_role):
    def decorator(func):
        @wraps(func)
        def wrapper(user, *args, **kwargs):
            if user.get("role") != required_role:
                raise PermissionError("Không có quyền truy cập")

            return func(user, *args, **kwargs)

        return wrapper

    return decorator


@require_role("admin")
def delete_user(user, user_id):
    return f"Đã xóa user {user_id}"


admin = {
    "id": 1,
    "role": "admin",
}

print(delete_user(admin, 100))
```

Ứng dụng trong backend:

- Authentication.
- Authorization.
- Logging.
- Đo thời gian chạy.
- Retry.
- Caching.
- Transaction.
- Rate limiting.

---

## 5.4. Decorator đo thời gian

```python
import time
from functools import wraps


def measure_time(func):
    @wraps(func)
    def wrapper(*args, **kwargs):
        start_time = time.perf_counter()

        try:
            return func(*args, **kwargs)
        finally:
            elapsed = time.perf_counter() - start_time
            print(f"{func.__name__} chạy trong {elapsed:.4f} giây")

    return wrapper


@measure_time
def slow_function():
    time.sleep(1)


slow_function()
```

---

# 6. Context manager và `with`

Context manager giúp quản lý tài nguyên theo một vòng đời rõ ràng:

1. Mở hoặc khởi tạo tài nguyên.
2. Sử dụng tài nguyên.
3. Giải phóng tài nguyên dù thành công hay có lỗi.

## 6.1. Ví dụ đọc file

Không nên:

```python
file = open("data.txt", "r", encoding="utf-8")
content = file.read()
file.close()
```

Nếu lỗi xảy ra trước `file.close()`, file có thể không được đóng.

Nên dùng:

```python
with open("data.txt", "r", encoding="utf-8") as file:
    content = file.read()
```

Sau khi thoát khỏi block `with`, file được đóng tự động.

---

## 6.2. Tạo context manager bằng class

```python
class DatabaseConnection:
    def __enter__(self):
        print("Mở kết nối database")
        self.connection = "database connection"
        return self.connection

    def __exit__(self, exc_type, exc_value, traceback):
        print("Đóng kết nối database")

        if exc_type is not None:
            print("Có lỗi:", exc_value)

        return False


with DatabaseConnection() as connection:
    print("Đang sử dụng:", connection)
```

`__enter__` chạy khi vào block `with`.

`__exit__` chạy khi thoát block, kể cả khi có exception.

Nếu `__exit__` trả về `True`, exception được xem là đã xử lý. Thông thường nên trả về `False` để exception tiếp tục được ném ra.

---

## 6.3. Tạo context manager bằng `contextlib`

```python
from contextlib import contextmanager


@contextmanager
def database_transaction():
    print("BEGIN")

    try:
        yield
        print("COMMIT")
    except Exception:
        print("ROLLBACK")
        raise
    finally:
        print("CLOSE CONNECTION")


with database_transaction():
    print("Cập nhật dữ liệu")
```

Ứng dụng:

- File.
- Database connection.
- Database transaction.
- Lock.
- Network connection.
- Temporary directory.
- Đo thời gian.
- Trạng thái tài nguyên.

---

# 7. Exception handling

## 7.1. Cấu trúc cơ bản

```python
try:
    value = int("abc")
except ValueError:
    print("Giá trị không hợp lệ")
```

Không nên bắt exception quá rộng:

```python
try:
    value = int("abc")
except:
    print("Có lỗi")
```

Cách này có thể che giấu lỗi không mong muốn.

Nên bắt lỗi cụ thể:

```python
try:
    value = int("abc")
except ValueError as error:
    print("Lỗi chuyển kiểu:", error)
```

---

## 7.2. `else` và `finally`

```python
try:
    result = 10 / 2
except ZeroDivisionError:
    print("Không thể chia cho 0")
else:
    print("Kết quả:", result)
finally:
    print("Luôn chạy")
```

- `try`: đoạn code có thể lỗi.
- `except`: xử lý lỗi.
- `else`: chạy nếu không có lỗi.
- `finally`: luôn chạy.

---

## 7.3. Ném exception bằng `raise`

```python
def withdraw(balance, amount):
    if amount <= 0:
        raise ValueError("Số tiền phải lớn hơn 0")

    if amount > balance:
        raise ValueError("Số dư không đủ")

    return balance - amount
```

---

## 7.4. Custom exception

```python
class InsufficientBalanceError(Exception):
    pass


def withdraw(balance, amount):
    if amount > balance:
        raise InsufficientBalanceError("Số dư không đủ")

    return balance - amount
```

Custom exception giúp domain logic rõ ràng hơn.

---

## 7.5. Exception trong API

Ví dụ khái niệm:

```python
class UserNotFoundError(Exception):
    pass


def get_user(user_id):
    user = None

    if user is None:
        raise UserNotFoundError(f"Không tìm thấy user {user_id}")

    return user
```

Ở tầng API, exception có thể được chuyển thành HTTP response:

```json
{
  "error": "USER_NOT_FOUND",
  "message": "Không tìm thấy user 10"
}
```

Status code:

```text
404 Not Found
```

Không nên trả toàn bộ stack trace cho client vì có thể lộ thông tin hệ thống.

---

## 7.6. Logging exception

```python
import logging

logger = logging.getLogger(__name__)


def process_order(order_id):
    try:
        raise RuntimeError("Payment service unavailable")
    except RuntimeError:
        logger.exception(
            "Không thể xử lý đơn hàng order_id=%s",
            order_id,
        )
        raise
```

`logger.exception()` tự ghi cả stack trace khi gọi trong block `except`.

---

# 8. OOP trong Python

Bốn tính chất chính:

1. Encapsulation.
2. Inheritance.
3. Polymorphism.
4. Abstraction.

---

## 8.1. Encapsulation

Encapsulation là đóng gói dữ liệu và hành vi trong một object.

```python
class BankAccount:
    def __init__(self, owner, balance=0):
        self.owner = owner
        self._balance = balance

    def deposit(self, amount):
        if amount <= 0:
            raise ValueError("Số tiền phải lớn hơn 0")

        self._balance += amount

    def get_balance(self):
        return self._balance
```

Trong Python:

- `name`: public.
- `_name`: quy ước internal/protected.
- `__name`: name mangling.

Ví dụ:

```python
class User:
    def __init__(self):
        self.__password = "secret"
```

Python đổi tên nội bộ thành dạng gần giống:

```text
_User__password
```

Đây không phải cơ chế bảo mật tuyệt đối, chủ yếu để hạn chế truy cập nhầm.

---

## 8.2. Inheritance

Inheritance là kế thừa thuộc tính và hành vi từ class cha.

```python
class Employee:
    def __init__(self, name):
        self.name = name

    def work(self):
        return f"{self.name} đang làm việc"


class BackendDeveloper(Employee):
    def code_api(self):
        return f"{self.name} đang viết API"


developer = BackendDeveloper("Hậu")

print(developer.work())
print(developer.code_api())
```

Gọi constructor class cha:

```python
class BackendDeveloper(Employee):
    def __init__(self, name, language):
        super().__init__(name)
        self.language = language
```

Không nên lạm dụng inheritance quá sâu. Trong nhiều trường hợp, composition dễ bảo trì hơn.

---

## 8.3. Polymorphism

Polymorphism cho phép nhiều object có cùng interface nhưng hành vi khác nhau.

```python
class EmailNotification:
    def send(self, message):
        print("Gửi email:", message)


class SMSNotification:
    def send(self, message):
        print("Gửi SMS:", message)


def notify(notification_service, message):
    notification_service.send(message)


notify(EmailNotification(), "Đơn hàng thành công")
notify(SMSNotification(), "Đơn hàng thành công")
```

Python thường sử dụng duck typing:

> Nếu object có hành vi cần thiết, không nhất thiết phải thuộc cùng một class cha.

---

## 8.4. Abstraction

Abstraction ẩn chi tiết cài đặt và chỉ cung cấp interface cần thiết.

```python
from abc import ABC, abstractmethod


class PaymentGateway(ABC):
    @abstractmethod
    def charge(self, amount):
        pass


class VNPayGateway(PaymentGateway):
    def charge(self, amount):
        return f"Thanh toán {amount} qua VNPay"


class StripeGateway(PaymentGateway):
    def charge(self, amount):
        return f"Thanh toán {amount} qua Stripe"
```

Không thể khởi tạo abstract class nếu chưa cài đặt abstract method:

```python
gateway = PaymentGateway()
# TypeError
```

---

## 8.5. Composition

Composition là một object chứa và sử dụng object khác.

```python
class EmailService:
    def send(self, message):
        print("Email:", message)


class UserService:
    def __init__(self, email_service):
        self.email_service = email_service

    def create_user(self, name):
        print(f"Tạo user {name}")
        self.email_service.send("Chào mừng người dùng mới")


service = UserService(EmailService())
service.create_user("Hậu")
```

Composition thường linh hoạt hơn inheritance vì dependency có thể thay thế dễ dàng.

---

# 9. `classmethod`, `staticmethod`, `property`

## 9.1. Instance method

```python
class User:
    def __init__(self, name):
        self.name = name

    def introduce(self):
        return f"Tôi là {self.name}"
```

Instance method nhận `self` và thao tác với instance.

---

## 9.2. `classmethod`

`classmethod` nhận `cls`, đại diện cho class.

```python
class User:
    default_role = "user"

    def __init__(self, name, role):
        self.name = name
        self.role = role

    @classmethod
    def create_default(cls, name):
        return cls(name=name, role=cls.default_role)


user = User.create_default("Hậu")

print(user.role)
```

Ứng dụng phổ biến: alternative constructor.

```python
class User:
    def __init__(self, name, age):
        self.name = name
        self.age = age

    @classmethod
    def from_dict(cls, data):
        return cls(
            name=data["name"],
            age=data["age"],
        )


user = User.from_dict({
    "name": "Hậu",
    "age": 25,
})
```

---

## 9.3. `staticmethod`

`staticmethod` không nhận `self` hoặc `cls`.

```python
class PasswordValidator:
    @staticmethod
    def is_valid(password):
        return len(password) >= 8


print(PasswordValidator.is_valid("12345678"))
```

Dùng khi hàm có liên quan logic tới class nhưng không cần truy cập instance hoặc class state.

Tuy nhiên, nếu logic hoàn toàn độc lập, một module-level function cũng có thể rõ ràng hơn.

---

## 9.4. `property`

`property` cho phép sử dụng method giống như attribute.

```python
class User:
    def __init__(self, first_name, last_name):
        self.first_name = first_name
        self.last_name = last_name

    @property
    def full_name(self):
        return f"{self.first_name} {self.last_name}"


user = User("Vũ Công", "Hậu")

print(user.full_name)
```

Không cần gọi:

```python
user.full_name()
```

---

## Setter với `property`

```python
class Product:
    def __init__(self, price):
        self.price = price

    @property
    def price(self):
        return self._price

    @price.setter
    def price(self, value):
        if value < 0:
            raise ValueError("Giá không thể âm")

        self._price = value
```

Sử dụng:

```python
product = Product(100)

product.price = 200

print(product.price)
```

---

# 10. Thread, process và async

Đây là phần rất quan trọng trong phỏng vấn backend.

Trước tiên cần phân biệt:

- CPU-bound: tác vụ tiêu tốn CPU.
- I/O-bound: tác vụ chờ mạng, file, database hoặc service khác.

Ví dụ CPU-bound:

- Mã hóa video.
- Xử lý ảnh.
- Tính toán số học lớn.
- Machine learning inference nặng.

Ví dụ I/O-bound:

- Gọi API.
- Query database.
- Đọc file.
- Gửi email.
- Chờ message queue.

---

## 10.1. Threading

Thread là các luồng chạy trong cùng một process và chia sẻ bộ nhớ.

Ví dụ:

```python
import threading
import time


def download_file(file_name):
    print(f"Bắt đầu tải {file_name}")
    time.sleep(2)
    print(f"Hoàn thành {file_name}")


threads = []

for file_name in ["a.zip", "b.zip", "c.zip"]:
    thread = threading.Thread(
        target=download_file,
        args=(file_name,),
    )
    thread.start()
    threads.append(thread)

for thread in threads:
    thread.join()
```

Threading phù hợp với tác vụ I/O-bound.

Ưu điểm:

- Chia sẻ bộ nhớ dễ.
- Nhẹ hơn process.
- Phù hợp với code đồng bộ cần chạy song song tác vụ I/O.

Nhược điểm:

- Dễ xảy ra race condition.
- Cần lock khi cập nhật dữ liệu dùng chung.
- Không tăng tốc tốt cho CPU-bound trong CPython do GIL.

---

## 10.2. Race condition

```python
import threading

counter = 0


def increment():
    global counter

    for _ in range(100_000):
        counter += 1
```

Nhiều thread cùng cập nhật `counter` có thể làm kết quả không chính xác.

Dùng lock:

```python
import threading

counter = 0
lock = threading.Lock()


def increment():
    global counter

    for _ in range(100_000):
        with lock:
            counter += 1
```

Lock bảo đảm tại một thời điểm chỉ một thread chạy đoạn critical section.

---

## 10.3. Multiprocessing

Mỗi process có bộ nhớ và Python interpreter riêng.

```python
from multiprocessing import Process


def calculate():
    result = sum(number * number for number in range(10_000_000))
    print(result)


processes = [
    Process(target=calculate)
    for _ in range(2)
]

for process in processes:
    process.start()

for process in processes:
    process.join()
```

Multiprocessing phù hợp với CPU-bound.

Ưu điểm:

- Tận dụng nhiều CPU core.
- Không bị giới hạn bởi GIL giống threading cho Python bytecode.

Nhược điểm:

- Tốn bộ nhớ hơn.
- Chi phí tạo process lớn hơn.
- Chia sẻ dữ liệu phức tạp hơn.
- Cần IPC, queue hoặc shared memory.

Có thể dùng `ProcessPoolExecutor`:

```python
from concurrent.futures import ProcessPoolExecutor


def square(number):
    return number ** 2


with ProcessPoolExecutor() as executor:
    results = list(executor.map(square, range(10)))

print(results)
```

---

## 10.4. Asyncio

Asyncio sử dụng event loop để xử lý nhiều tác vụ I/O mà không cần tạo một thread cho mỗi tác vụ.

```python
import asyncio


async def fetch_data(name, delay):
    print(f"Bắt đầu {name}")
    await asyncio.sleep(delay)
    print(f"Hoàn thành {name}")
    return name


async def main():
    results = await asyncio.gather(
        fetch_data("service-a", 2),
        fetch_data("service-b", 1),
        fetch_data("service-c", 3),
    )

    print(results)


asyncio.run(main())
```

Tổng thời gian gần bằng tác vụ lâu nhất, thay vì tổng thời gian của tất cả tác vụ.

Async phù hợp khi:

- Có nhiều request I/O đồng thời.
- Gọi nhiều API ngoài.
- Query database bằng async driver.
- WebSocket.
- Streaming.
- High-concurrency backend.

Async không tự làm CPU-bound nhanh hơn.

Ví dụ không tốt:

```python
async def calculate_heavy():
    return sum(i * i for i in range(100_000_000))
```

Hàm này vẫn chặn event loop vì không có điểm `await` thực sự nhường quyền.

---

## 10.5. So sánh

| Mô hình | Phù hợp | Bộ nhớ | Chia sẻ dữ liệu | GIL |
|---|---|---:|---|---|
| Thread | I/O-bound | Trung bình | Dễ | Bị ảnh hưởng |
| Process | CPU-bound | Cao | Khó hơn | Mỗi process có GIL riêng |
| Async | Nhiều I/O đồng thời | Thấp | Cùng process | Không giải quyết CPU-bound |

---

# 11. GIL là gì?

GIL là viết tắt của Global Interpreter Lock.

Trong CPython, GIL bảo đảm tại một thời điểm chỉ một thread thực thi Python bytecode trong một process.

Ví dụ có 4 thread tính toán CPU-bound không có nghĩa là cả 4 thread thực sự chạy Python bytecode song song trên 4 CPU core.

## Vì sao Python có GIL?

Một lý do lịch sử quan trọng là giúp quản lý bộ nhớ và reference counting đơn giản, an toàn hơn.

Python thường quản lý object bằng reference count:

```python
import sys

data = []

print(sys.getrefcount(data))
```

Nếu nhiều thread cùng thay đổi reference count mà không được bảo vệ, có thể phát sinh lỗi phức tạp.

GIL giúp giảm độ phức tạp của việc đồng bộ nội bộ interpreter, nhưng đánh đổi bằng giới hạn parallelism cho CPU-bound Python code.

---

## GIL có làm threading vô dụng không?

Không.

Với I/O-bound, thread thường nhường GIL khi chờ:

- Network.
- Disk.
- Database.
- Một số system call.

Do đó threading vẫn hữu ích khi tác vụ chủ yếu là chờ I/O.

---

## Cách xử lý CPU-bound

- Dùng multiprocessing.
- Dùng native extension giải phóng GIL.
- Dùng NumPy hoặc thư viện tính toán tối ưu.
- Đưa tác vụ nặng sang worker riêng.
- Dùng ngôn ngữ hoặc service khác nếu cần.

---

# 12. `async/await` phù hợp trường hợp nào?

`async/await` phù hợp nhất với I/O-bound và concurrency cao.

## 12.1. Ví dụ API gọi nhiều dịch vụ ngoài

Giả sử một endpoint cần gọi:

- User service: 1 giây.
- Order service: 2 giây.
- Payment service: 1.5 giây.

Chạy tuần tự:

```python
async def get_dashboard():
    user = await get_user()
    orders = await get_orders()
    payments = await get_payments()

    return {
        "user": user,
        "orders": orders,
        "payments": payments,
    }
```

Nếu mỗi hàm chờ xong mới chạy hàm tiếp theo, tổng thời gian khoảng 4.5 giây.

Chạy đồng thời:

```python
import asyncio


async def get_dashboard():
    user, orders, payments = await asyncio.gather(
        get_user(),
        get_orders(),
        get_payments(),
    )

    return {
        "user": user,
        "orders": orders,
        "payments": payments,
    }
```

Tổng thời gian gần với request lâu nhất, khoảng 2 giây.

---

## 12.2. Timeout

Không nên chờ vô hạn một dịch vụ ngoài.

```python
import asyncio


async def call_external_service():
    await asyncio.sleep(10)
    return {"status": "ok"}


async def main():
    try:
        result = await asyncio.wait_for(
            call_external_service(),
            timeout=2,
        )
        print(result)
    except asyncio.TimeoutError:
        print("Dịch vụ phản hồi quá lâu")


asyncio.run(main())
```

Trong Python hiện đại có thể dùng:

```python
async with asyncio.timeout(2):
    result = await call_external_service()
```

---

## 12.3. Xử lý lỗi riêng từng service

`asyncio.gather` mặc định sẽ ném exception nếu một task lỗi.

```python
results = await asyncio.gather(
    call_service_a(),
    call_service_b(),
    call_service_c(),
    return_exceptions=True,
)
```

Sau đó cần kiểm tra từng kết quả:

```python
for result in results:
    if isinstance(result, Exception):
        print("Có service bị lỗi:", result)
```

Trong hệ thống thực tế, nên xác định:

- Service nào bắt buộc.
- Service nào có thể trả fallback.
- Có retry hay không.
- Timeout bao lâu.
- Có circuit breaker hay không.
- Có cache dữ liệu cũ hay không.

---

## 12.4. Giới hạn concurrency

Không nên gọi hàng nghìn request cùng lúc không giới hạn.

```python
import asyncio

semaphore = asyncio.Semaphore(10)


async def limited_call(item):
    async with semaphore:
        return await call_external_service(item)
```

Semaphore giới hạn tối đa 10 tác vụ chạy đồng thời.

---

## 12.5. Blocking code trong async

Không nên:

```python
import time


async def handler():
    time.sleep(5)
```

`time.sleep()` chặn toàn bộ event loop.

Nên dùng:

```python
import asyncio


async def handler():
    await asyncio.sleep(5)
```

Nếu buộc phải gọi blocking function:

```python
result = await asyncio.to_thread(blocking_function, argument)
```

---

# 13. Virtual environment và quản lý dependency

## 13.1. Tại sao cần virtual environment?

Mỗi project có thể cần phiên bản thư viện khác nhau.

Ví dụ:

```text
Project A cần Django 4
Project B cần Django 5
```

Nếu cài tất cả thư viện vào môi trường global, dễ xảy ra xung đột.

Virtual environment tạo môi trường Python riêng cho từng project.

---

## 13.2. Tạo virtual environment

Windows:

```powershell
python -m venv .venv
```

Kích hoạt trong PowerShell:

```powershell
.\.venv\Scripts\Activate.ps1
```

Kích hoạt trong Command Prompt:

```cmd
.venv\Scripts\activate.bat
```

Linux/macOS:

```bash
python3 -m venv .venv && source .venv/bin/activate
```

Thoát virtual environment:

```bash
deactivate
```

Kiểm tra Python đang dùng:

```bash
python -c "import sys; print(sys.executable)"
```

---

## 13.3. Cài dependency

```bash
python -m pip install fastapi uvicorn
```

Nên dùng:

```bash
python -m pip
```

thay vì chỉ `pip`, vì bảo đảm pip thuộc đúng Python interpreter.

---

## 13.4. `requirements.txt`

Xuất dependency:

```bash
python -m pip freeze > requirements.txt
```

Cài lại:

```bash
python -m pip install -r requirements.txt
```

Ví dụ:

```text
fastapi==0.115.0
uvicorn==0.30.6
sqlalchemy==2.0.35
```

Nhược điểm của `pip freeze` là ghi cả dependency trực tiếp và dependency gián tiếp, đôi khi khó quản lý.

---

## 13.5. Dependency trực tiếp và dependency gián tiếp

Ví dụ bạn cài:

```text
fastapi
```

FastAPI lại phụ thuộc vào:

```text
starlette
pydantic
```

- FastAPI là dependency trực tiếp.
- Starlette và Pydantic là dependency gián tiếp.

---

## 13.6. Version constraint

```text
fastapi==0.115.0
```

Chính xác một phiên bản.

```text
fastapi>=0.115.0
```

Cho phép phiên bản mới hơn.

```text
fastapi>=0.115.0,<1.0.0
```

Cho phép nâng cấp trong phạm vi kiểm soát.

Pin version quá lỏng có thể khiến build hôm nay và build tháng sau khác nhau.

---

## 13.7. Các công cụ phổ biến

### pip-tools

Có thể quản lý file input và lock dependency:

```text
requirements.in
requirements.txt
```

### Poetry

Quản lý:

- Dependency.
- Virtual environment.
- Package.
- Lock file.

Các file thường gặp:

```text
pyproject.toml
poetry.lock
```

### uv

Công cụ quản lý Python project và dependency có tốc độ nhanh.

Các file thường gặp:

```text
pyproject.toml
uv.lock
```

Khi phỏng vấn, không nhất thiết phải dùng tất cả. Quan trọng là hiểu:

- Môi trường cô lập.
- Lock phiên bản.
- Reproducible build.
- Dependency trực tiếp và gián tiếp.
- Tránh cài package trực tiếp lên production bằng thao tác thủ công.

---

## 13.8. Dependency trong Docker

Ví dụ Dockerfile:

```dockerfile
FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .

RUN python -m pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["python", "main.py"]
```

Nên copy file dependency trước source code để tận dụng Docker layer cache.

Nếu source code thay đổi nhưng `requirements.txt` không đổi, Docker không cần cài lại toàn bộ package.

---

# 14. Câu hỏi phỏng vấn thường gặp

## 14.1. List và tuple khác nhau thế nào?

Câu trả lời gợi ý:

> List và tuple đều là sequence có thứ tự, cho phép phần tử trùng nhau và truy cập bằng index. Điểm khác biệt chính là list mutable còn tuple immutable. List phù hợp với dữ liệu cần thêm, xóa hoặc cập nhật. Tuple phù hợp với dữ liệu cố định, có thể dùng làm dictionary key nếu toàn bộ phần tử bên trong hashable. Tuple cũng thể hiện rõ ý nghĩa rằng dữ liệu không nên bị thay đổi.

Ví dụ:

```python
users = ["Hậu", "Nam"]
users.append("Linh")

point = (10, 20)
```

---

## 14.2. Generator giúp tiết kiệm bộ nhớ ra sao?

Câu trả lời gợi ý:

> List tạo và lưu toàn bộ phần tử trong bộ nhớ ngay lập tức. Generator sử dụng lazy evaluation, chỉ tạo ra từng giá trị khi được yêu cầu. Vì vậy generator phù hợp khi xử lý file lớn, query nhiều record hoặc stream dữ liệu. Đổi lại, generator thường chỉ duyệt được một lần và không hỗ trợ truy cập ngẫu nhiên bằng index như list.

Ví dụ:

```python
users = [load_user(user_id) for user_id in range(1_000_000)]
```

Có thể tốn nhiều RAM.

```python
users = (load_user(user_id) for user_id in range(1_000_000))
```

Generator tạo từng user khi cần.

---

## 14.3. Decorator hoạt động thế nào?

Câu trả lời gợi ý:

> Decorator là một callable nhận vào một hàm hoặc class và trả về một callable mới. Nó thường dùng để bổ sung hành vi trước hoặc sau hàm gốc mà không sửa trực tiếp code của hàm đó. Cú pháp `@decorator` là syntactic sugar cho việc gán lại `function = decorator(function)`. Trong backend, decorator thường dùng cho authentication, authorization, logging, transaction, retry và caching.

---

## 14.4. Vì sao Python có GIL?

Câu trả lời gợi ý:

> Trong CPython, GIL bảo đảm một thời điểm chỉ một thread thực thi Python bytecode. Nó giúp đơn giản hóa việc quản lý bộ nhớ, đặc biệt là reference counting, và giảm độ phức tạp đồng bộ nội bộ interpreter. Nhược điểm là thread không tăng tốc tốt cho CPU-bound Python code. Tuy nhiên threading vẫn có ích cho I/O-bound vì thread thường nhường GIL khi chờ network hoặc disk. Với CPU-bound, có thể dùng multiprocessing hoặc native library có khả năng giải phóng GIL.

---

## 14.5. Khi nào dùng threading, multiprocessing và asyncio?

Câu trả lời gợi ý:

> Threading phù hợp với I/O-bound khi đang dùng thư viện đồng bộ hoặc cần chia sẻ bộ nhớ trong cùng process. Multiprocessing phù hợp với CPU-bound vì mỗi process có Python interpreter và GIL riêng, cho phép tận dụng nhiều CPU core. Asyncio phù hợp với lượng lớn tác vụ I/O đồng thời như gọi nhiều API, WebSocket hoặc database async, với chi phí thấp hơn việc tạo nhiều thread. Asyncio không làm CPU-bound nhanh hơn và blocking code có thể chặn toàn bộ event loop.

Ví dụ:

```text
Gọi 100 API ngoài       -> asyncio
Đọc nhiều file          -> threading hoặc asyncio tùy thư viện
Resize 1.000 ảnh        -> multiprocessing
Tính toán số học lớn    -> multiprocessing
WebSocket server        -> asyncio
```

---

## 14.6. Làm thế nào xử lý một API gọi nhiều dịch vụ bên ngoài?

Một câu trả lời tốt nên đề cập:

1. Gọi song song các service độc lập.
2. Đặt timeout cho từng service.
3. Retry có giới hạn cho lỗi tạm thời.
4. Không retry mù quáng với mọi lỗi.
5. Dùng exponential backoff và jitter.
6. Xác định service bắt buộc và service tùy chọn.
7. Có fallback hoặc cache nếu phù hợp.
8. Dùng circuit breaker.
9. Logging và distributed tracing.
10. Request ID hoặc correlation ID.
11. Giới hạn concurrency.
12. Không giữ database transaction trong lúc chờ API ngoài.
13. Bảo đảm idempotency nếu request có thể bị gửi lại.

Ví dụ:

```python
import asyncio


async def get_customer_dashboard(customer_id):
    user_task = get_user_service(customer_id)
    order_task = get_order_service(customer_id)
    score_task = get_credit_score_service(customer_id)

    user, orders, score = await asyncio.gather(
        user_task,
        order_task,
        score_task,
    )

    return {
        "user": user,
        "orders": orders,
        "credit_score": score,
    }
```

Phiên bản có timeout:

```python
import asyncio


async def get_customer_dashboard(customer_id):
    async with asyncio.timeout(3):
        user, orders, score = await asyncio.gather(
            get_user_service(customer_id),
            get_order_service(customer_id),
            get_credit_score_service(customer_id),
        )

    return {
        "user": user,
        "orders": orders,
        "credit_score": score,
    }
```

Không nên mở database transaction rồi chờ nhiều external service:

```python
async with database.transaction():
    await update_order()
    await call_slow_external_service()
```

Transaction kéo dài có thể giữ lock lâu, giảm khả năng phục vụ và tăng nguy cơ deadlock.

---

# 15. Kiến thức chung về Python framework

JD có thể không yêu cầu cụ thể Django, Flask hay FastAPI, nhưng bạn cần hiểu các khái niệm chung.

---

## 15.1. Routing

Routing ánh xạ HTTP request vào handler.

Ví dụ khái niệm:

```text
GET /users/10
```

được ánh xạ tới:

```python
def get_user(user_id):
    pass
```

Ví dụ FastAPI:

```python
from fastapi import FastAPI

app = FastAPI()


@app.get("/users/{user_id}")
def get_user(user_id: int):
    return {
        "id": user_id,
        "name": "Hậu",
    }
```

Cần hiểu:

- Path parameter.
- Query parameter.
- Request body.
- HTTP method.
- Status code.

---

## 15.2. Middleware

Middleware là lớp xử lý nằm giữa request và handler.

Luồng:

```text
Client
  ↓
Middleware 1
  ↓
Middleware 2
  ↓
Route handler
  ↓
Middleware 2
  ↓
Middleware 1
  ↓
Response
```

Ứng dụng:

- Logging.
- Request ID.
- Authentication.
- CORS.
- Metrics.
- Timing.
- Error handling.

Ví dụ:

```python
import time
from fastapi import FastAPI, Request

app = FastAPI()


@app.middleware("http")
async def add_process_time(request: Request, call_next):
    start = time.perf_counter()

    response = await call_next(request)

    elapsed = time.perf_counter() - start
    response.headers["X-Process-Time"] = str(elapsed)

    return response
```

---

## 15.3. Validation

Validation kiểm tra dữ liệu đầu vào.

Ví dụ:

```python
from pydantic import BaseModel, Field


class CreateUserRequest(BaseModel):
    name: str = Field(min_length=2, max_length=100)
    age: int = Field(ge=18, le=100)
    email: str
```

Không nên tin dữ liệu client gửi lên.

Cần validate:

- Kiểu dữ liệu.
- Độ dài.
- Khoảng giá trị.
- Format.
- Business rule.
- Quyền truy cập.

Validation schema không thay thế hoàn toàn business validation.

Ví dụ:

```text
age >= 18
```

là validation đơn giản.

```text
Người dùng chỉ được tạo tối đa 3 tài khoản phụ
```

là business rule, thường cần kiểm tra database.

---

## 15.4. ORM

ORM ánh xạ object trong code với bảng database.

Ví dụ khái niệm:

```python
user = User(
    name="Hậu",
    email="hau@example.com",
)

session.add(user)
session.commit()
```

Thay vì tự viết:

```sql
INSERT INTO users (name, email)
VALUES ('Hậu', 'hau@example.com');
```

Ưu điểm:

- Code dễ đọc.
- Tái sử dụng model.
- Hỗ trợ relation.
- Migration.
- Parameter binding giúp giảm SQL injection nếu dùng đúng cách.

Nhược điểm:

- Có thể sinh query không tối ưu.
- Dễ gặp N+1 query.
- Vẫn cần biết SQL.
- Query phức tạp đôi khi dùng raw SQL rõ ràng hơn.

Ví dụ N+1:

```python
users = get_all_users()

for user in users:
    print(user.orders)
```

Nếu mỗi lần truy cập `user.orders` lại chạy một query, 100 user có thể tạo 101 query.

Cần dùng eager loading hoặc query phù hợp.

---

## 15.5. Authentication

Authentication trả lời câu hỏi:

> Bạn là ai?

Ví dụ:

- Username/password.
- Session cookie.
- JWT.
- OAuth 2.0.
- API key.

Luồng JWT đơn giản:

```text
Client gửi username/password
  ↓
Server xác thực
  ↓
Server cấp access token
  ↓
Client gửi Bearer token
  ↓
Server kiểm tra chữ ký và claim
```

Lưu ý:

- Không lưu password dạng plain text.
- Dùng password hashing như Argon2 hoặc bcrypt.
- Access token nên có thời gian sống ngắn.
- Refresh token cần quản lý cẩn thận.
- JWT không tự động hỗ trợ logout hoàn hảo.
- Không đưa thông tin nhạy cảm vào payload JWT.

---

## 15.6. Authorization

Authorization trả lời câu hỏi:

> Bạn được phép làm gì?

Ví dụ:

```text
Admin được xóa user.
Manager được xem báo cáo.
Nhân viên chỉ được xem dữ liệu của mình.
```

Các mô hình:

- RBAC: Role-Based Access Control.
- ABAC: Attribute-Based Access Control.
- Permission-based.
- Resource ownership.

Không chỉ kiểm tra user đã đăng nhập. Phải kiểm tra user có quyền trên resource cụ thể hay không.

Ví dụ lỗi bảo mật:

```text
GET /orders/100
```

Nếu user A đổi thành:

```text
GET /orders/101
```

và xem được đơn hàng user B, đây là lỗi authorization.

---

## 15.7. Dependency injection

Dependency injection là truyền dependency từ bên ngoài thay vì tự tạo cứng bên trong class hoặc function.

Không tốt:

```python
class UserService:
    def __init__(self):
        self.repository = PostgresUserRepository()
```

`UserService` bị gắn chặt với PostgreSQL repository.

Tốt hơn:

```python
class UserService:
    def __init__(self, repository):
        self.repository = repository
```

Sử dụng:

```python
repository = PostgresUserRepository()
service = UserService(repository)
```

Khi test:

```python
fake_repository = FakeUserRepository()
service = UserService(fake_repository)
```

Lợi ích:

- Dễ test.
- Giảm coupling.
- Dễ thay implementation.
- Quản lý lifecycle dependency.
- Code rõ trách nhiệm hơn.

Ví dụ FastAPI:

```python
from fastapi import Depends, FastAPI

app = FastAPI()


def get_database():
    database = Database()

    try:
        yield database
    finally:
        database.close()


@app.get("/users/{user_id}")
def get_user(
    user_id: int,
    database=Depends(get_database),
):
    return database.get_user(user_id)
```

---

## 15.8. Cấu trúc project backend

Một cấu trúc đơn giản:

```text
app/
├── main.py
├── api/
│   ├── routes/
│   │   └── users.py
│   └── dependencies.py
├── models/
│   └── user.py
├── schemas/
│   └── user.py
├── repositories/
│   └── user_repository.py
├── services/
│   └── user_service.py
├── core/
│   ├── config.py
│   ├── security.py
│   └── logging.py
└── database/
    ├── session.py
    └── migrations/
```

Luồng phổ biến:

```text
Route
  ↓
Service
  ↓
Repository
  ↓
Database
```

- Route xử lý HTTP.
- Service xử lý business logic.
- Repository xử lý truy cập dữ liệu.
- Model ánh xạ database.
- Schema định nghĩa request/response.

Không phải project nào cũng cần nhiều layer. Project nhỏ nên giữ cấu trúc vừa đủ, tránh over-engineering.

---

# 16. Bài tập thực hành

## Bài 1: Loại bỏ email trùng

Input:

```python
emails = [
    "a@example.com",
    "b@example.com",
    "a@example.com",
]
```

Yêu cầu:

- Loại bỏ email trùng.
- Giữ kết quả dễ đọc.
- Giải thích vì sao dùng `set`.

---

## Bài 2: Viết generator đọc file

Viết hàm:

```python
def read_lines(file_path):
    pass
```

Hàm phải:

- Đọc từng dòng.
- Không đưa toàn bộ file vào RAM.
- Loại bỏ newline.
- Tự đóng file.

---

## Bài 3: Decorator logging

Viết decorator:

```python
@log_execution
def create_user(name):
    return {
        "name": name,
    }
```

Decorator cần log:

- Tên hàm.
- Thời điểm bắt đầu.
- Thời gian chạy.
- Có lỗi hay không.

---

## Bài 4: Context manager transaction

Viết context manager in ra:

```text
BEGIN
COMMIT
```

Nếu lỗi:

```text
BEGIN
ROLLBACK
```

Cuối cùng luôn in:

```text
CLOSE
```

---

## Bài 5: OOP notification

Tạo interface:

```python
class NotificationService:
    def send(self, message):
        pass
```

Cài đặt:

- EmailNotification.
- SMSNotification.
- SlackNotification.

Sau đó truyền implementation vào `OrderService`.

---

## Bài 6: Async gọi nhiều service

Giả lập ba service bằng `asyncio.sleep`.

Yêu cầu:

- Chạy đồng thời.
- Timeout toàn bộ sau 3 giây.
- Một service lỗi không làm mất toàn bộ dữ liệu nếu service đó không bắt buộc.
- Trả kết quả fallback.

---

## Bài 7: Chọn concurrency model

Giải thích nên dùng gì cho từng tình huống:

1. Gọi 200 API.
2. Resize 10.000 ảnh.
3. Đọc 50 file log.
4. WebSocket server.
5. Export Excel rất lớn.
6. Gửi email nền.
7. Tính hash mật khẩu cho nhiều người dùng.

---

# 17. Checklist trước phỏng vấn

## Python core

- [ ] Phân biệt list, tuple, set, dict.
- [ ] Hiểu mutable và immutable.
- [ ] Hiểu reference, shallow copy và deep copy.
- [ ] Giải thích được lỗi default mutable argument.
- [ ] Sử dụng được `*args`, `**kwargs`.
- [ ] Viết được list comprehension.
- [ ] Viết được generator với `yield`.
- [ ] Viết được decorator.
- [ ] Hiểu context manager.
- [ ] Bắt exception cụ thể.
- [ ] Viết custom exception.

## OOP

- [ ] Giải thích inheritance.
- [ ] Giải thích encapsulation.
- [ ] Giải thích abstraction.
- [ ] Giải thích polymorphism.
- [ ] Biết composition.
- [ ] Phân biệt instance method, classmethod và staticmethod.
- [ ] Biết dùng property.

## Concurrency

- [ ] Phân biệt CPU-bound và I/O-bound.
- [ ] Biết khi nào dùng threading.
- [ ] Biết khi nào dùng multiprocessing.
- [ ] Biết khi nào dùng asyncio.
- [ ] Giải thích GIL.
- [ ] Biết race condition và lock.
- [ ] Không gọi blocking code trực tiếp trong event loop.
- [ ] Biết `asyncio.gather`.
- [ ] Biết timeout và giới hạn concurrency.

## Backend framework

- [ ] Hiểu routing.
- [ ] Hiểu middleware.
- [ ] Hiểu validation.
- [ ] Hiểu ORM và N+1.
- [ ] Phân biệt authentication và authorization.
- [ ] Hiểu dependency injection.
- [ ] Biết cách tổ chức route, service, repository.
- [ ] Biết cách xử lý exception thành HTTP response.
- [ ] Biết bảo vệ password và token.
- [ ] Biết timeout, retry, fallback khi gọi service ngoài.

## Dependency

- [ ] Biết tạo virtual environment.
- [ ] Biết kích hoạt trên Windows và Linux.
- [ ] Biết dùng `python -m pip`.
- [ ] Hiểu `requirements.txt`.
- [ ] Hiểu lock file.
- [ ] Hiểu dependency trực tiếp và gián tiếp.
- [ ] Biết cách tận dụng Docker layer cache.

---

# Kết luận

Khi phỏng vấn, không nên chỉ trả lời định nghĩa. Một câu trả lời tốt thường có cấu trúc:

1. Định nghĩa ngắn gọn.
2. Khi nào sử dụng.
3. Ví dụ thực tế.
4. Ưu điểm.
5. Nhược điểm hoặc trade-off.
6. Liên hệ với dự án đã làm.

Ví dụ khi được hỏi về asyncio:

> Asyncio là mô hình concurrency dựa trên event loop, phù hợp với I/O-bound và số lượng tác vụ đồng thời lớn. Em thường dùng khi một API cần gọi nhiều service ngoài hoặc xử lý WebSocket. Nó giúp giảm thời gian chờ bằng cách chạy các tác vụ I/O đồng thời, nhưng không phù hợp để tăng tốc CPU-bound. Khi dùng async, em chú ý timeout, retry, giới hạn concurrency và tránh gọi blocking function trực tiếp trong event loop.

Đây là cách trả lời vừa có kiến thức nền, vừa thể hiện tư duy triển khai thực tế.
