# 11. Behavioral Interview

## 1. STAR là gì?

- Situation: bối cảnh.
- Task: trách nhiệm của bạn.
- Action: bạn đã làm gì.
- Result: kết quả đo được hoặc bài học.

Không kể lan man. Mỗi câu chuyện nên khoảng 1–2 phút.

## 2. Bug production

Khung trả lời:

```text
Situation:
Production xảy ra lỗi...

Task:
Tôi chịu trách nhiệm xác định nguyên nhân và khôi phục dịch vụ...

Action:
Tôi kiểm tra log, metric, thay đổi gần nhất...
Tôi rollback hoặc đưa hotfix...
Sau đó thêm test và cảnh báo...

Result:
Dịch vụ phục hồi sau...
Tỷ lệ lỗi giảm...
```

## 3. API hoặc database chậm

Nên đề cập:

- Xác định endpoint chậm.
- Dùng log/APM.
- Kiểm tra `EXPLAIN ANALYZE`.
- Phát hiện missing index hoặc N+1.
- Tối ưu query.
- Đo trước và sau.
- Thêm monitoring.

## 4. Bất đồng với đồng nghiệp

Không đổ lỗi.

Cấu trúc:

- Nêu khác biệt về giải pháp.
- Làm rõ tiêu chí quyết định.
- Dùng benchmark, tài liệu hoặc prototype.
- Đồng thuận theo dữ liệu.
- Tôn trọng quyết định cuối cùng.

## 5. Tự học công nghệ mới

Nêu:

- Vì sao cần học.
- Cách chia nhỏ.
- Tài liệu chính thức.
- Lab nhỏ.
- Áp dụng vào task.
- Chia sẻ lại với team.

## 6. Một nhiệm vụ làm chưa tốt

Chọn lỗi thật nhưng không quá nguy hiểm.

Nêu rõ:

- Bạn nhận trách nhiệm.
- Nguyên nhân.
- Cách sửa.
- Quy trình ngăn tái diễn.
- Bài học.

## 7. Xử lý áp lực

Câu trả lời tốt:

- Ưu tiên theo impact.
- Giao tiếp rõ với stakeholder.
- Chia task.
- Không thay đổi nhiều thứ cùng lúc.
- Ghi lại timeline.
- Sau sự cố có postmortem.

## 8. Giới thiệu bản thân

Khung:

```text
Tên và số năm kinh nghiệm
→ thế mạnh Backend Python
→ database và API
→ Docker/Linux/CI-CD
→ một dự án tiêu biểu
→ lý do phù hợp vị trí
```

Ví dụ:

> Em là Vũ Công Hậu, định hướng Backend Developer với Python. Em có kinh nghiệm làm API, làm việc với PostgreSQL, Docker, Linux và quy trình triển khai ứng dụng. Gần đây em tập trung nhiều vào hệ thống self-hosted, CI/CD, quản lý môi trường và xử lý dữ liệu giữa các service. Em quan tâm vị trí này vì công việc có cả backend, database, messaging và hạ tầng, phù hợp với hướng phát triển của em.

## 9. Lý do chuyển việc

Tập trung vào hướng phát triển:

> Em muốn tìm môi trường có hệ thống backend thực tế hơn, nhiều bài toán về database, tích hợp service và vận hành production để phát triển sâu hơn về Backend Engineering.

Không nói xấu công ty cũ.

## 10. Mức lương mong muốn

Cách trả lời:

> Dựa trên phạm vi công việc và năng lực hiện tại, em mong muốn mức khoảng X–Y. Tuy nhiên em vẫn sẵn sàng trao đổi dựa trên tổng package, trách nhiệm thực tế và cơ hội phát triển.

## Checklist

- [ ] Có 5 câu chuyện STAR.
- [ ] Có số liệu hoặc kết quả cụ thể.
- [ ] Không đổ lỗi.
- [ ] Mỗi câu trả lời dưới 2 phút.
- [ ] Giải thích được vai trò cá nhân.
- [ ] Có bài học sau mỗi tình huống.
