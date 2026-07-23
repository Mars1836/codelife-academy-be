# 8. Git cho Backend Developer

Git là hệ thống quản lý phiên bản phân tán. Git không chỉ lưu “file mới nhất” mà lưu lịch sử snapshot, quan hệ giữa commit và các tham chiếu như branch, tag và HEAD.

## 1. Ba khu vực quan trọng

```text
Working tree → Staging area → Repository
```

- Working tree: file đang sửa.
- Staging area/index: nội dung sẽ vào commit tiếp theo.
- Repository: lịch sử commit trong `.git`.

```bash
git status
git add src/app.py
git commit -m "feat: add health endpoint"
```

`git add` không “đưa file lên GitHub”; nó đưa phiên bản hiện tại của file vào staging area.

---

## 2. Commit là gì?

Commit là object chứa:

- Tham chiếu đến tree đại diện snapshot.
- Parent commit.
- Author, committer, thời gian.
- Commit message.

Commit được nhận diện bằng hash. Nếu nội dung hoặc metadata thay đổi, hash thay đổi.

```bash
git show <commit>
git log --oneline --graph --decorate --all
```

Commit nên nhỏ, có mục đích rõ ràng và có thể review độc lập.

---

## 3. Branch và HEAD

Branch là con trỏ có thể di chuyển đến commit.

```text
main → C3
feature → C5
HEAD → feature
```

`HEAD` thường trỏ đến branch hiện tại.

```bash
git switch -c feature/payment
```

Khi commit mới, branch hiện tại di chuyển tới commit mới.

Detached HEAD xảy ra khi checkout trực tiếp một commit hoặc tag. Commit mới lúc này không thuộc branch nếu không tạo branch để giữ nó.

---

## 4. Remote, fetch và pull

```bash
git remote -v
git fetch origin
git pull origin main
git push origin feature/payment
```

`git fetch` tải commit và cập nhật remote-tracking branch như `origin/main`, nhưng không tự nhập thay đổi vào branch hiện tại.

`git pull` về cơ bản là fetch rồi merge hoặc rebase tùy cấu hình.

Nên fetch và kiểm tra thay đổi khi muốn kiểm soát rõ:

```bash
git fetch origin
git log --oneline HEAD..origin/main
git rebase origin/main
```

---

## 5. Merge

Merge kết hợp hai lịch sử.

```bash
git switch main
git merge feature/payment
```

Nếu lịch sử phân nhánh, Git có thể tạo merge commit.

Ưu điểm:

- Không viết lại lịch sử.
- Thể hiện rõ nhánh được hợp nhất.

Nhược điểm:

- Lịch sử có thể nhiều merge commit.

Fast-forward xảy ra khi branch đích chưa có commit mới và chỉ cần di chuyển con trỏ.

---

## 6. Rebase

Rebase phát lại commit của branch lên base mới.

```bash
git switch feature/payment
git fetch origin
git rebase origin/main
```

Lịch sử từ:

```text
A---B---C main
     \
      D---E feature
```

thành:

```text
A---B---C---D'---E' feature
```

Commit hash thay đổi vì parent thay đổi.

Không rebase lịch sử public mà người khác đang dựa vào nếu chưa thống nhất, vì sẽ buộc họ xử lý lịch sử bị viết lại.

---

## 7. Squash

Squash gộp nhiều commit thành ít commit hơn.

```bash
git rebase -i HEAD~4
```

Ứng dụng:

- Gộp commit “fix typo”, “fix again” trước khi merge.
- Tạo lịch sử dễ đọc.

Không nên squash mù quáng nếu các commit riêng biệt có giá trị cho review hoặc rollback.

---

## 8. Conflict

Conflict xảy ra khi Git không tự quyết định cách kết hợp thay đổi, thường do sửa cùng vùng file hoặc một bên xóa trong khi bên kia sửa.

```text
<<<<<<< HEAD
code của branch hiện tại
=======
code của branch còn lại
>>>>>>> feature
```

Quy trình:

```bash
git status
# sửa file
git add <file>
git commit
```

Khi rebase:

```bash
git add <file>
git rebase --continue
```

Hủy:

```bash
git merge --abort
git rebase --abort
```

Không chọn “ours” hoặc “theirs” mà chưa hiểu logic nghiệp vụ. Conflict được giải quyết đúng khi code cuối cùng giữ được ý định cần thiết của cả hai thay đổi.

---

## 9. Reset

### `--soft`

Di chuyển HEAD, giữ thay đổi trong staging.

```bash
git reset --soft HEAD~1
```

Dùng khi muốn viết lại commit gần nhất nhưng giữ toàn bộ nội dung đã stage.

### `--mixed`

Mặc định. Di chuyển HEAD, bỏ staging nhưng giữ file trong working tree.

```bash
git reset HEAD~1
```

### `--hard`

Di chuyển HEAD và xóa thay đổi ở staging/working tree.

```bash
git reset --hard HEAD~1
```

`--hard` có thể làm mất công việc chưa commit. Kiểm tra `git status` trước khi dùng.

---

## 10. Revert

```bash
git revert <commit>
```

Revert tạo commit mới đảo ngược thay đổi, không xóa lịch sử cũ. Đây là lựa chọn an toàn cho branch đã push và được nhiều người dùng.

So sánh:

- Reset: di chuyển lịch sử; phù hợp local/private history.
- Revert: thêm commit đảo ngược; phù hợp shared branch.

---

## 11. Restore và checkout

Khôi phục file chưa stage:

```bash
git restore src/app.py
```

Bỏ file khỏi staging:

```bash
git restore --staged src/app.py
```

`git checkout` trước đây làm nhiều nhiệm vụ. Git hiện đại tách thành `git switch` cho branch và `git restore` cho file để rõ nghĩa hơn.

---

## 12. Stash

```bash
git stash push -m "WIP payment"
git stash list
git stash pop
```

Stash hữu ích khi cần chuyển branch tạm thời nhưng chưa muốn commit.

Không nên dùng stash như kho lưu trữ dài hạn; dễ quên và khó review.

Bao gồm untracked file:

```bash
git stash -u
```

---

## 13. Cherry-pick

```bash
git cherry-pick <commit>
```

Cherry-pick áp dụng thay đổi của một commit lên branch hiện tại và tạo commit mới.

Ứng dụng:

- Đưa hotfix từ main sang release branch.
- Lấy một commit cụ thể mà không merge cả branch.

Lạm dụng cherry-pick có thể tạo các commit tương đương với hash khác nhau và làm lịch sử khó theo dõi.

---

## 14. Reflog

Reflog lưu lịch sử di chuyển của HEAD/ref trong local repository.

```bash
git reflog
```

Nếu reset nhầm:

```bash
git reflog
git reset --hard <old-head>
```

Reflog rất hữu ích để cứu commit local “bị mất”, nhưng không phải backup vĩnh viễn và không được push lên remote.

---

## 15. Tag

```bash
git tag -a v1.0.0 -m "release v1.0.0"
git push origin v1.0.0
```

Tag thường đánh dấu release. Annotated tag chứa metadata và message; lightweight tag chỉ là ref đơn giản.

Không nên di chuyển tag release đã công bố nếu artifact/deployment dựa vào tag đó.

---

## 16. `.gitignore` và `.gitattributes`

`.gitignore` bỏ qua file chưa được track:

```text
.env
.venv/
__pycache__/
*.pyc
```

Nếu file đã được track, thêm vào `.gitignore` không tự xóa khỏi repository:

```bash
git rm --cached .env
```

Nếu secret đã commit, cần rotate secret; xóa file ở commit mới không làm secret biến mất khỏi lịch sử.

`.gitattributes` điều khiển thuộc tính như line ending, diff và merge behavior:

```text
* text=auto
*.sh text eol=lf
*.bat text eol=crlf
```

---

## 17. Force push

Sau rebase, branch remote có lịch sử cũ nên có thể cần force push:

```bash
git push --force-with-lease
```

Ưu tiên `--force-with-lease` thay `--force` vì nó từ chối ghi đè nếu remote đã có thay đổi bạn chưa biết.

Không force-push lên main hoặc shared branch nếu quy trình không cho phép.

---

## 18. Quy trình feature branch thực tế

```text
main
 ↓ tạo branch
feature/payment
 ↓ commit nhỏ
push
 ↓ Pull Request
review + CI
 ↓ cập nhật theo review
merge/squash
 ↓ deploy
```

Lệnh:

```bash
git switch main
git pull --ff-only
git switch -c feature/payment
git add .
git commit -m "feat: add payment endpoint"
git push -u origin feature/payment
```

Trước khi mở PR:

- Chạy test và lint.
- Kiểm tra diff.
- Loại bỏ debug code và secret.
- Viết mô tả thay đổi, cách test và rủi ro.

---

## 19. Commit message

Một convention phổ biến:

```text
feat: add transfer endpoint
fix: prevent duplicate payment
refactor: extract account repository
test: add concurrent withdrawal tests
docs: expand database guide
chore: update dependencies
```

Message nên nói mục đích, không chỉ nói “update code”.

Commit tốt giúp review, changelog, bisect và rollback dễ hơn.

---

## 20. Git bisect

Dùng binary search để tìm commit gây bug:

```bash
git bisect start
git bisect bad
git bisect good <known-good-commit>
```

Sau mỗi lần Git checkout commit giữa, test rồi đánh dấu:

```bash
git bisect good
# hoặc
git bisect bad
```

Kết thúc:

```bash
git bisect reset
```

Có thể tự động hóa bằng test script.

---

## 21. Câu hỏi phỏng vấn

### Merge khác rebase thế nào?

Merge kết hợp lịch sử và có thể tạo merge commit, không đổi commit cũ. Rebase phát lại commit lên base mới, tạo hash mới và lịch sử tuyến tính hơn. Không rebase shared history nếu chưa thống nhất.

### Revert khác reset?

Revert tạo commit mới đảo thay đổi nên an toàn với branch public. Reset di chuyển branch về commit khác và có thể viết lại lịch sử.

### Fetch khác pull?

Fetch chỉ tải dữ liệu và cập nhật remote-tracking refs. Pull fetch rồi tích hợp vào branch hiện tại bằng merge hoặc rebase.

### Khi nào conflict xảy ra?

Khi Git không thể tự kết hợp thay đổi, thường vì hai nhánh sửa cùng vùng hoặc xóa/sửa cùng file. Cần giải quyết theo ý nghĩa code, không chỉ xóa marker.

### Vì sao không nên commit `.env`?

Nó thường chứa secret và cấu hình môi trường. Git lưu lịch sử nên xóa ở commit sau không đủ; secret đã lộ cần được rotate.

### `--force-with-lease` tốt hơn `--force` ở đâu?

Nó kiểm tra remote ref còn đúng trạng thái bạn biết, giảm nguy cơ ghi đè commit mới của người khác.

---

## 22. Bài tập thực hành

1. Tạo branch, commit ba thay đổi rồi rebase lên main.
2. Tạo conflict có chủ đích và giải quyết.
3. Thử reset soft, mixed và hard trên repository test.
4. Revert một commit đã push.
5. Dùng reflog cứu commit sau reset nhầm.
6. Dùng cherry-pick đưa hotfix sang release branch.
7. Dùng bisect tìm commit làm test thất bại.

## Checklist

- [ ] Hiểu working tree, staging và repository.
- [ ] Giải thích được commit, branch và HEAD.
- [ ] Phân biệt fetch, pull, merge và rebase.
- [ ] Xử lý conflict đúng quy trình.
- [ ] Phân biệt reset, restore và revert.
- [ ] Biết stash, cherry-pick, reflog và bisect.
- [ ] Không commit secret hoặc force-push shared branch tùy tiện.
