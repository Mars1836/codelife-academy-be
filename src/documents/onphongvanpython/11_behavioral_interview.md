# 11. Behavioral Interview cho Backend Developer

Behavioral interview đánh giá cách ứng viên xử lý tình huống thực tế, giao tiếp, chịu trách nhiệm, học hỏi và phối hợp với người khác. Nhà tuyển dụng không chỉ nghe kết quả, họ còn quan tâm cách bạn suy nghĩ và hành động.

## 1. STAR là gì?

STAR là cấu trúc trả lời gồm:

- Situation: bối cảnh cụ thể.
- Task: trách nhiệm hoặc mục tiêu của bạn.
- Action: chính bạn đã làm gì.
- Result: kết quả và bài học.

Một câu trả lời tốt thường dài 1–2 phút, đủ cụ thể nhưng không kể lan man.

---

## 2. Ví dụ STAR hoàn chỉnh

### Câu hỏi

“Hãy kể về một lần hệ thống production gặp sự cố.”

### Trả lời

**Situation:** Trong một lần deploy API mới, tỷ lệ lỗi 502 tăng mạnh ngay sau release. Hệ thống dùng Nginx reverse proxy đến container backend.

**Task:** Tôi chịu trách nhiệm xác định nguyên nhân và khôi phục dịch vụ nhanh nhất có thể.

**Action:** Tôi kiểm tra dashboard và xác nhận lỗi bắt đầu đúng thời điểm deploy. Sau đó tôi kiểm tra Nginx error log, container log và port đang listen. Tôi phát hiện phiên bản mới chỉ bind vào `127.0.0.1` bên trong container thay vì `0.0.0.0`, nên Nginx ở container khác không kết nối được. Tôi rollback image trước để khôi phục dịch vụ, sửa cấu hình bind, bổ sung smoke test kiểm tra endpoint qua Docker network rồi deploy lại.

**Result:** Dịch vụ hoạt động lại sau khoảng 10 phút. Sau sự cố, pipeline có thêm smoke test nên lỗi tương tự được phát hiện ở staging trước production.

Điểm tốt của câu trả lời:

- Có bối cảnh rõ.
- Nêu vai trò cá nhân.
- Có quy trình debug.
- Có hành động giảm thiểu ngay và cải tiến dài hạn.
- Không đổ lỗi.

---

## 3. Chuẩn bị ngân hàng câu chuyện

Nên chuẩn bị ít nhất 6 câu chuyện có thể tái sử dụng:

1. Một dự án thành công.
2. Một sự cố production.
3. Một lần mắc lỗi.
4. Một bất đồng kỹ thuật.
5. Một lần deadline gấp.
6. Một lần học công nghệ mới.
7. Một lần cải thiện hiệu năng hoặc quy trình.
8. Một lần hỗ trợ đồng đội.

Mỗi câu chuyện nên ghi:

```text
Bối cảnh:
Mục tiêu:
Vai trò của tôi:
Hành động cụ thể:
Khó khăn:
Kết quả đo được:
Bài học:
```

---

## 4. Câu hỏi “Giới thiệu bản thân”

Cấu trúc:

```text
Hiện tại → kinh nghiệm liên quan → điểm mạnh → lý do quan tâm vị trí
```

Ví dụ:

> Em là Vũ Công Hậu, định hướng Backend Developer. Em tập trung vào Python, PostgreSQL, Docker và triển khai ứng dụng trên Linux. Trong các dự án gần đây, em quan tâm nhiều đến thiết kế API, transaction database, CI/CD và cách vận hành hệ thống ổn định. Em thích tìm hiểu bản chất vấn đề thay vì chỉ làm cho code chạy, đặc biệt ở các phần concurrency, database và deployment. Em quan tâm vị trí này vì công việc có nhiều bài toán backend thực tế và cơ hội phát triển trong hệ thống fintech.

Không kể toàn bộ CV theo thứ tự thời gian. Tập trung thông tin liên quan vị trí.

---

## 5. Điểm mạnh

Điểm mạnh cần có bằng chứng.

Cách trả lời:

```text
Điểm mạnh → ví dụ → tác động
```

Ví dụ:

> Điểm mạnh của em là khả năng tự điều tra vấn đề theo từng lớp. Khi API không truy cập được, em không restart ngẫu nhiên mà kiểm tra từ DNS, port, process, reverse proxy, container log đến dependency. Cách này giúp em tìm nguyên nhân nhanh hơn và ghi lại được runbook cho lần sau.

Tránh trả lời chung chung như “em chăm chỉ” mà không có ví dụ.

---

## 6. Điểm yếu

Một điểm yếu tốt phải:

- Có thật.
- Không phá hủy yêu cầu cốt lõi của vị trí.
- Có hành động cải thiện.
- Cho thấy tiến bộ.

Ví dụ:

> Trước đây em có xu hướng đi quá sâu vào chi tiết kỹ thuật trước khi xác nhận mức độ ưu tiên. Em nhận ra điều này có thể làm chậm tiến độ, nên hiện tại em thường làm rõ mục tiêu, deadline và mức độ rủi ro trước, sau đó mới quyết định phần nào cần tối ưu sâu. Cách này giúp em cân bằng tốt hơn giữa chất lượng và thời gian.

Không dùng câu giả yếu như “em quá cầu toàn” nếu không giải thích ảnh hưởng và cách cải thiện.

---

## 7. Kể về một lần mắc lỗi

Nhà tuyển dụng muốn xem:

- Bạn có nhận trách nhiệm không?
- Bạn phát hiện và sửa thế nào?
- Bạn ngăn lỗi lặp lại ra sao?

Ví dụ khung trả lời:

> Em từng viết migration thêm constraint trực tiếp trên bảng lớn mà chưa kiểm tra thời gian lock ở môi trường gần production. Khi chạy staging, migration kéo dài và ảnh hưởng request. Em chủ động dừng rollout, kiểm tra execution plan và chuyển sang quy trình backfill theo batch rồi tạo index phù hợp trước khi thêm constraint. Sau đó em bổ sung checklist migration gồm kích thước bảng, lock, rollback và thời gian dự kiến. Bài học của em là thay đổi schema phải được đánh giá như thay đổi production, không chỉ như một câu SQL.

Không chọn lỗi quá nhỏ như typo nếu câu hỏi muốn đánh giá ownership.

---

## 8. Bất đồng với đồng đội

Tập trung vào vấn đề, không công kích con người.

Cấu trúc:

1. Nêu điểm hai bên khác nhau.
2. Xác định tiêu chí quyết định.
3. Thu thập dữ liệu hoặc làm thử nghiệm nhỏ.
4. Thống nhất và cam kết với quyết định.

Ví dụ:

> Trong một dự án, em đề xuất PostgreSQL còn đồng đội muốn dùng MongoDB vì schema linh hoạt. Thay vì tranh luận theo sở thích, chúng em liệt kê yêu cầu về transaction, quan hệ, query báo cáo và tốc độ thay đổi schema. Sau một prototype nhỏ, nhóm chọn PostgreSQL với JSONB cho custom field. Dù giải pháp không hoàn toàn giống đề xuất ban đầu của từng người, nó đáp ứng tốt nhất tiêu chí chung.

---

## 9. Deadline gấp

Câu trả lời nên thể hiện khả năng ưu tiên và giao tiếp rủi ro.

Ví dụ:

> Khi deadline ngắn, em chia yêu cầu thành phần bắt buộc và phần có thể hoãn. Em xác nhận với product về acceptance criteria, ưu tiên luồng chính, thêm test cho rủi ro lớn và tránh refactor không cần thiết. Em cập nhật sớm nếu có blocker thay vì chờ sát deadline. Sau khi release ổn định, em tạo task xử lý technical debt còn lại.

Không nói “em làm xuyên đêm” như giải pháp mặc định. Làm thêm giờ đôi khi cần thiết nhưng không thay thế planning.

---

## 10. Khi không biết câu trả lời kỹ thuật

Cách phản hồi tốt:

> Phần này em chưa trực tiếp triển khai production nên em không muốn khẳng định quá mức. Theo hiểu biết hiện tại của em, cơ chế là... Nếu xử lý thực tế, em sẽ kiểm tra tài liệu chính thức, dựng một thử nghiệm nhỏ và xác nhận bằng metric/log trước khi áp dụng.

Điều này tốt hơn đoán chắc chắn.

Có thể suy luận thành tiếng ở mức tóm tắt:

- Làm rõ yêu cầu.
- Nêu giả định.
- Đưa ra cách tiếp cận.
- Nêu trade-off.
- Xác định cách kiểm chứng.

---

## 11. Khi chưa có kinh nghiệm đúng công nghệ

Kết nối với kiến thức tương đương:

> Em chưa vận hành Kafka production, nhưng em đã làm với RabbitMQ và hiểu các vấn đề chung như delivery guarantee, retry, duplicate message, idempotency và DLQ. Em biết Kafka khác ở mô hình distributed log, partition và consumer group. Nếu tham gia dự án, em có thể bắt đầu từ use case cụ thể và học phần vận hành còn thiếu.

Không nói “em biết” nếu chỉ đọc qua. Phân biệt:

- Đã dùng production.
- Đã dùng trong dự án cá nhân/lab.
- Hiểu khái niệm.
- Chưa biết.

---

## 12. Dự án tự hào nhất

Cần giải thích:

- Vấn đề thực tế.
- Quy mô hoặc constraint.
- Quyết định kỹ thuật của bạn.
- Trade-off.
- Kết quả.

Ví dụ:

> Dự án em tự hào nhất là xây dựng quy trình triển khai ứng dụng backend bằng Docker Compose và GitHub Actions. Trước đó deploy phụ thuộc thao tác thủ công nên khó lặp lại. Em chuẩn hóa image, tách environment, thêm healthcheck và pipeline test/build/deploy. Sau đó việc release rõ version hơn, rollback nhanh hơn và giảm lỗi do khác biệt môi trường. Điều em học được là CI/CD không chỉ là script deploy mà còn cần artifact bất biến, secret management và kiểm tra sau release.

---

## 13. Câu hỏi về ownership

Ownership không có nghĩa tự làm tất cả. Nó có nghĩa:

- Nhận trách nhiệm về kết quả.
- Chủ động làm rõ vấn đề.
- Kéo đúng người vào khi cần.
- Theo dõi đến khi hoàn thành.
- Cải thiện hệ thống sau sự cố.

Ví dụ:

> Khi phát hiện lỗi nằm ở dependency của team khác, em vẫn thu thập log, request ID và bước tái hiện trước khi liên hệ. Em phối hợp theo dõi đến khi fix được deploy và xác nhận metric phục hồi, thay vì chỉ chuyển ticket rồi coi như xong.

---

## 14. Câu hỏi về ưu tiên công việc

Khung quyết định:

- Mức ảnh hưởng user/business.
- Độ khẩn cấp.
- Rủi ro dữ liệu hoặc bảo mật.
- Dependency/blocker.
- Chi phí và thời gian xử lý.

Ví dụ thứ tự:

```text
Data loss/security incident
→ production outage
→ blocker của nhiều người
→ deadline cam kết
→ cải tiến và technical debt
```

Cần trao đổi với manager/product nếu ưu tiên xung đột, không tự đoán âm thầm.

---

## 15. Câu hỏi về feedback

Ví dụ:

> Em từng nhận feedback rằng pull request của em quá lớn nên khó review. Em xem lại và nhận thấy mình thường gom refactor với feature. Sau đó em tách thay đổi thành PR nhỏ, viết rõ cách test và rủi ro. Thời gian review giảm và conflict cũng ít hơn. Em học được rằng code tốt nhưng khó review vẫn làm chậm cả team.

Thể hiện rằng bạn nghe, kiểm chứng và thay đổi hành vi.

---

## 16. Câu hỏi “Vì sao muốn nghỉ công việc cũ?”

Giữ thái độ chuyên nghiệp, tập trung vào hướng đi tương lai.

Ví dụ:

> Em muốn tìm môi trường có nhiều bài toán backend và cơ hội phát triển sâu hơn về database, distributed system và vận hành production. Em trân trọng kinh nghiệm ở công việc hiện tại, nhưng vị trí này phù hợp hơn với hướng phát triển dài hạn của em.

Không nói xấu công ty, quản lý hoặc đồng nghiệp cũ.

---

## 17. Câu hỏi “Vì sao chúng tôi nên tuyển bạn?”

Kết nối yêu cầu vị trí với bằng chứng:

> Em phù hợp vì có nền tảng Python Backend, SQL, Docker và Linux, đồng thời em đang tập trung sâu vào transaction, concurrency và CI/CD — những phần quan trọng với hệ thống fintech. Điểm em có thể đóng góp là khả năng học nhanh, điều tra vấn đề có hệ thống và chủ động biến bài học thành tài liệu hoặc quy trình để lỗi không lặp lại. Những phần em chưa có nhiều kinh nghiệm production, em sẽ nói rõ và có kế hoạch học bằng thử nghiệm thực tế.

---

## 18. Câu hỏi dành cho nhà tuyển dụng

Nên hỏi câu giúp hiểu công việc:

- Backend hiện là monolith hay microservice?
- Những vấn đề kỹ thuật lớn nhất team đang giải quyết là gì?
- Quy trình review, test và deploy hiện tại ra sao?
- Team theo dõi chất lượng production bằng metric nào?
- Vai trò này được kỳ vọng đạt kết quả gì trong ba tháng đầu?
- Cách team xử lý incident và chia sẻ bài học?
- Cơ hội mentoring và ownership dự án thế nào?

Tránh chỉ hỏi quyền lợi khi interviewer đang đánh giá chuyên môn; có thể hỏi HR ở vòng phù hợp.

---

## 19. Những lỗi thường gặp

- Trả lời quá dài, không có điểm chính.
- Nói “chúng em” suốt nhưng không rõ cá nhân làm gì.
- Đổ lỗi đồng đội.
- Bịa số liệu hoặc kinh nghiệm.
- Chỉ kể hành động mà không có kết quả.
- Không nêu bài học.
- Dùng cùng một câu chuyện cho mọi câu hỏi dù không phù hợp.
- Học thuộc từng chữ khiến câu trả lời thiếu tự nhiên.

---

## 20. Cách luyện tập

1. Viết bullet cho mỗi câu chuyện STAR.
2. Ghi âm câu trả lời 2 phút.
3. Nghe lại và bỏ chi tiết thừa.
4. Bổ sung số liệu có thật.
5. Luyện follow-up: “Vì sao?”, “Trade-off?”, “Bạn sẽ làm khác gì?”.
6. Không học thuộc từng câu; nhớ cấu trúc và dữ kiện.

Một câu chuyện tốt có thể trả lời nhiều dạng câu hỏi, nhưng cần điều chỉnh trọng tâm.

---

## 21. Bộ câu hỏi luyện tập

1. Giới thiệu bản thân.
2. Dự án bạn tự hào nhất.
3. Một lỗi bạn đã gây ra.
4. Một incident production.
5. Một lần bất đồng kỹ thuật.
6. Một deadline khó.
7. Một lần yêu cầu không rõ ràng.
8. Một lần cải thiện hiệu năng.
9. Một lần nhận feedback khó nghe.
10. Một lần hỗ trợ đồng đội.
11. Một công nghệ bạn phải học nhanh.
12. Một quyết định có trade-off.
13. Khi nào bạn chủ động nói “không”.
14. Bạn ưu tiên bug và feature thế nào?
15. Bạn làm gì khi không biết câu trả lời?

## Checklist trước phỏng vấn

- [ ] Có phần giới thiệu 60–90 giây.
- [ ] Có ít nhất 6 câu chuyện STAR.
- [ ] Mỗi câu chuyện nêu rõ vai trò cá nhân.
- [ ] Có kết quả hoặc tác động cụ thể.
- [ ] Không đổ lỗi.
- [ ] Có bài học và hành động phòng ngừa.
- [ ] Trung thực về mức kinh nghiệm.
- [ ] Có 3–5 câu hỏi dành cho nhà tuyển dụng.
