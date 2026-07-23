# 8. Git

## 1. Các lệnh cơ bản

```bash
git clone <repository-url>
```

```bash
git fetch origin
```

```bash
git pull
```

```bash
git push
```

`fetch` chỉ tải thông tin remote. `pull` thường là fetch rồi merge hoặc rebase.

## 2. Branch

```bash
git switch -c feature/login
```

```bash
git switch main
```

```bash
git branch
```

```bash
git branch -d feature/login
```

## 3. Merge

```bash
git switch main
```

```bash
git merge feature/login
```

Merge giữ lịch sử nhánh và có thể tạo merge commit.

## 4. Rebase

```bash
git switch feature/login
```

```bash
git rebase main
```

Rebase viết lại commit lên đầu mới, giúp lịch sử tuyến tính hơn.

Không nên rebase các commit đã chia sẻ rộng nếu việc rewrite history gây ảnh hưởng người khác.

## 5. Conflict

Quy trình:

```bash
git status
```

Sửa file conflict, sau đó:

```bash
git add <file>
```

Nếu merge:

```bash
git commit
```

Nếu rebase:

```bash
git rebase --continue
```

Hủy rebase:

```bash
git rebase --abort
```

## 6. reset và revert

### reset

Di chuyển HEAD. Có thể thay đổi working tree.

```bash
git reset --soft HEAD~1
```

Giữ thay đổi trong staging.

```bash
git reset --mixed HEAD~1
```

Giữ thay đổi trong working tree.

```bash
git reset --hard HEAD~1
```

Xóa thay đổi local, cần rất cẩn thận.

### revert

```bash
git revert <commit>
```

Tạo commit mới đảo ngược commit cũ. An toàn hơn cho shared branch.

## 7. stash

```bash
git stash
```

```bash
git stash list
```

```bash
git stash pop
```

Dùng khi cần tạm cất thay đổi chưa commit.

## 8. Cherry-pick

```bash
git cherry-pick <commit>
```

Áp dụng một commit cụ thể sang branch hiện tại.

## 9. Pull request

Luồng:

```text
Tạo branch
→ commit
→ push
→ mở PR
→ CI chạy
→ review
→ sửa feedback
→ merge
```

PR nên:

- Scope nhỏ.
- Mô tả rõ.
- Có test.
- Nêu rủi ro.
- Không trộn refactor không liên quan.

## 10. GitFlow và trunk-based

### GitFlow

Có nhiều branch dài hạn:

```text
main
develop
feature/*
release/*
hotfix/*
```

Phù hợp quy trình release truyền thống nhưng có thể phức tạp.

### Trunk-based

Nhánh sống ngắn, merge thường xuyên vào main/trunk, dùng feature flag nếu cần.

Phù hợp CI/CD nhanh.

## 11. .gitignore

Ví dụ Python:

```text
.venv/
__pycache__/
*.pyc
.env
.pytest_cache/
```

## 12. .gitattributes

Dùng để cấu hình:

- Line ending.
- Diff.
- Merge behavior.
- Binary file.
- Git LFS.

Ví dụ:

```text
* text=auto
*.sh text eol=lf
*.bat text eol=crlf
```

## Câu hỏi phỏng vấn

### Merge và rebase khác nhau?

Merge kết hợp lịch sử và có thể tạo merge commit. Rebase viết lại commit để lịch sử tuyến tính hơn.

### reset và revert khác nhau?

Reset thay đổi lịch sử local/HEAD. Revert tạo commit mới đảo ngược thay đổi, phù hợp shared branch.

### Đã push commit lỗi lên shared branch?

Dùng `git revert`, không force-push nếu không có thỏa thuận rõ.

### Conflict xảy ra khi nào?

Khi Git không tự quyết định cách kết hợp thay đổi, thường do sửa cùng vùng file hoặc xóa/sửa đồng thời.
