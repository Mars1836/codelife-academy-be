Tài liệu rèn luyện thuật toán & tư duy hệ thống 

## **TUYỂN TẬP 20 BÀI TOÁN DSA THỰC CHIẾN** 

_Tài liệu tổng hợp 20 bài toán chia thành các chủ đề: Mảng, Chuỗi, Cấu trúc dữ liệu tuyến tính, Cây, Đồ thị và Quy hoạch động. Triển khai code chuẩn bằng Go, tập trung vào tư duy tối ưu bộ nhớ và thời gian xử lý._ 

## **Phần 1: Mảng, Chuỗi & Bảng Băm (Arrays, Strings & Hashing)** 

## **1. Contains Duplicate (Kiểm tra phần tử trùng lặp)** 

_Cho mảng số nguyên nums. Trả về true nếu có bất kỳ giá trị nào xuất hiện ít nhất hai lần, ngược lại trả về false._ 

**Ý tưởng:** Sử dụng Hash Map/Set để lưu trữ các phần tử đã duyệt. Kiểm tra sự tồn tại trong _**O(1)**_ . 

```go
func containsDuplicate(nums []int) bool {
    seen := make(map[int]struct{})
    for _, num := range nums {
        if _, ok := seen[num]; ok {
            return true
        }
        seen[num] = struct{}{} // Dùng struct rỗng để tiết kiệm bộ nhớ
    }
    return false
}
```

Time: _**O(N)**_ | Space: _**O(N)**_ 

Trang 1 / 13 

Tài liệu rèn luyện thuật toán & tư duy hệ thống 

## **2. Valid Anagram (Kiểm tra chuỗi đảo chữ)** 

_Cho hai chuỗi s và t, trả về true nếu t là anagram của s (chứa các ký tự giống nhau với cùng số lượng)._ 

**Ý tưởng:** Đếm tần suất xuất hiện của các ký tự. Với bảng chữ cái tiếng Anh, dùng mảng cố định 26 phần tử sẽ tối ưu hơn dùng Hash Map. 

```go
func isAnagram(s string, t string) bool {
    if len(s) != len(t) { return false }
    counts := make([]int, 26)
    for i := 0; i < len(s); i++ {
        counts[s[i]-'a']++
        counts[t[i]-'a']--
    }
    for _, c := range counts {
        if c != 0 { return false }
    }
    return true
}
```

Time: _**O(N)**_ | Space: _**O(1)**_ 

## **3. Product of Array Except Self (Tích các phần tử ngoại trừ chính nó)** 

_Cho mảng nums, trả về mảng answer sao cho answer[i] bằng tích của tất cả các phần tử ngoại trừ nums[i]. (Yêu cầu không dùng phép chia, time O(N))._ 

**Ý tưởng:** Tính tích tích lũy từ trái sang phải (prefix) và từ phải sang trái (suffix). Kết quả tại i là tích của prefix[i-1] và suffix[i+1]. 

```go
func productExceptSelf(nums []int) []int {
    n := len(nums)
    res := make([]int, n)
    prefix := 1
    for i := 0; i < n; i++ {
        res[i] = prefix
        prefix *= nums[i]
    }
    postfix := 1
    for i := n - 1; i >= 0; i-- {
        res[i] *= postfix
        postfix *= nums[i]
    }
    return res
}
```

Time: _**O(N)**_ | Space: _**O(1)**_ (không tính mảng kết quả) 

Trang 2 / 13 

Tài liệu rèn luyện thuật toán & tư duy hệ thống 

## **4. Maximum Subarray (Mảng con có tổng lớn nhất)** 

_Tìm mảng con liên tiếp có tổng lớn nhất trong một mảng số nguyên._ 

**Ý tưởng (Kadane's Algorithm):** Duyệt qua mảng, duy trì tổng cục bộ. Nếu tổng cục bộ nhỏ hơn 0, ta reset nó về 0 vì nó sẽ làm giảm tổng của chuỗi tiếp theo. 

```go
func maxSubArray(nums []int) int {
    maxSum := nums[0]
    currentSum := 0
    for _, n := range nums {
        if currentSum < 0 { currentSum = 0 }
        currentSum += n
        if currentSum > maxSum { maxSum = currentSum }
    }
    return maxSum
}
```

Time: _**O(N)**_ | Space: _**O(1)**_ 

## **Phần 2: Hai Con Trỏ & Cửa Sổ Trượt (Two Pointers & Sliding Window)** 

## **5. Valid Palindrome (Chuỗi đối xứng)** 

_Kiểm tra một chuỗi có phải là chuỗi đối xứng hay không, chỉ xét các ký tự chữ và số, bỏ qua khoảng trắng và ký tự đặc biệt._ 

**Ý tưởng:** Dùng hai con trỏ từ hai đầu chuỗi tiến vào giữa. Bỏ qua các ký tự không hợp lệ. 

```go
import "strings"

func isPalindrome(s string) bool {
    s = strings.ToLower(s)
    l, r := 0, len(s)-1
    for l < r {
        for l < r && !isAlphanumeric(s[l]) { l++ }
        for l < r && !isAlphanumeric(s[r]) { r-- }
        if s[l] != s[r] { return false }
        l++; r--
    }
    return true
}

func isAlphanumeric(c byte) bool {
    return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}
```

Time: _**O(N)**_ | Space: _**O(1)**_ 

Trang 3 / 13 

Tài liệu rèn luyện thuật toán & tư duy hệ thống 

## **6. 3Sum (Tổng 3 số bằng 0)** 

_Tìm tất cả các bộ 3 số trong mảng có tổng bằng 0. Không được chứa bộ 3 trùng lặp._ 

**Ý tưởng:** Sắp xếp mảng trước. Cố định một số thứ nhất, dùng 2 con trỏ tìm 2 số còn lại. Bỏ qua các giá trị trùng lặp khi duyệt. 

```go
import "sort"

func threeSum(nums []int) [][]int {
    sort.Ints(nums)
    res := [][]int{}
    for i := 0; i < len(nums)-2; i++ {
        if i > 0 && nums[i] == nums[i-1] { continue } // Tránh trùng
        l, r := i+1, len(nums)-1
        for l < r {
            sum := nums[i] + nums[l] + nums[r]
            if sum > 0 { r-- } else if sum < 0 { l++ } else {
                res = append(res, []int{nums[i], nums[l], nums[r]})
                l++; r--
                for l < r && nums[l] == nums[l-1] { l++ }
            }
        }
    }
    return res
}
```

Time: _**O(N^2)**_ | Space: _**O(1)**_ hoặc _**O(N)**_ tùy thuật toán sort. 

## **7. Container With Most Water (Tối ưu hóa diện tích hình học)** 

_Cho mảng chiều cao các cột, tìm 2 cột tạo thành một thùng chứa được nhiều nước nhất._ 

**Ý tưởng:** Bài toán tối ưu diện tích _**Area = min(h[L], h[R]) × (R - L)**_ . Dùng 2 con trỏ ở 2 đầu. Luôn di chuyển con trỏ đang trỏ vào cột thấp hơn với hy vọng tìm được cột cao hơn ở phía trong. 

```go
func maxArea(height []int) int {
    l, r := 0, len(height)-1
    maxA := 0
    for l < r {
        area := (r - l) * min(height[l], height[r])
        if area > maxA { maxA = area }
        if height[l] < height[r] { l++ } else { r-- }
    }
    return maxA
}

func min(a, b int) int { if a < b { return a }; return b }
```

Time: _**O(N)**_ | Space: _**O(1)**_ 

Trang 4 / 13 

Tài liệu rèn luyện thuật toán & tư duy hệ thống 

## **8. Minimum Window Substring (Cửa sổ trượt động)** 

_Cho chuỗi s và t. Tìm chuỗi con ngắn nhất trong s chứa tất cả các ký tự của t._ 

**Ý tưởng:** Mở rộng cửa sổ (con trỏ right) đến khi chứa đủ t. Sau đó thu hẹp cửa sổ (con trỏ left) để tìm độ dài ngắn nhất. Sử dụng mảng/map để đếm tần suất. 

```go
func minWindow(s string, t string) string {
    if len(t) == 0 { return "" }
    countT, window := make(map[byte]int), make(map[byte]int)
    for i := range t { countT[t[i]]++ }
    have, need := 0, len(countT)
    res, resLen := []int{-1, -1}, 1<<31-1 // Max int
    l := 0
    for r := range s {
        c := s[r]
        window[c]++
        if countT[c] > 0 && window[c] == countT[c] { have++ }
        for have == need {
            if (r - l + 1) < resLen { res = []int{l, r}; resLen = r - l + 1 }
            window[s[l]]--
            if countT[s[l]] > 0 && window[s[l]] < countT[s[l]] { have-- }
            l++
        }
    }
    if resLen != 1<<31-1 { return s[res[0] : res[1]+1] }
    return ""
}
```

Time: _**O(S + T)**_ | Space: _**O(S + T)**_ 

Trang 5 / 13 

Tài liệu rèn luyện thuật toán & tư duy hệ thống 

## **Phần 3: Stack, Ngăn xếp & Tính toán** 

## **9. Valid Parentheses (Ngoặc hợp lệ)** 

_Kiểm tra chuỗi chứa '()', '{}', '[]' có được đóng mở hợp lệ không._ 

**Ý tưởng:** Dùng stack. Gặp dấu mở thì push vào stack, gặp dấu đóng thì pop ra kiểm tra xem có khớp không. 

```go
func isValid(s string) bool {
    stack := []rune{}
    match := map[rune]rune{')': '(', '}': '{', ']': '['}
    for _, char := range s {
        if val, ok := match[char]; ok {
            if len(stack) > 0 && stack[len(stack)-1] == val {
                stack = stack[:len(stack)-1]
            } else { return false }
        } else { stack = append(stack, char) }
    }
    return len(stack) == 0
}
```

Time: _**O(N)**_ | Space: _**O(N)**_ 

Trang 6 / 13 

Tài liệu rèn luyện thuật toán & tư duy hệ thống 

## **10. Evaluate Reverse Polish Notation (Ký pháp Ba Lan ngược)** 

_Tính toán giá trị của biểu thức toán học dạng hậu tố (RPN)._ 

**Ý tưởng:** Duyệt qua các token. Nếu là số thì đẩy vào stack, nếu là toán tử thì pop 2 số trên cùng ra, thực hiện phép toán và đẩy kết quả lại vào stack. 

```go
import "strconv"

func evalRPN(tokens []string) int {
    stack := []int{}
    for _, t := range tokens {
        if t == "+" || t == "-" || t == "*" || t == "/" {
            b, a := stack[len(stack)-1], stack[len(stack)-2]
            stack = stack[:len(stack)-2]
            switch t {
            case "+": stack = append(stack, a+b)
            case "-": stack = append(stack, a-b)
            case "*": stack = append(stack, a*b)
            case "/": stack = append(stack, a/b)
            }
        } else {
            num, _ := strconv.Atoi(t)
            stack = append(stack, num)
        }
    }
    return stack[0]
}
```

Time: _**O(N)**_ | Space: _**O(N)**_ 

## **Phần 4: Tìm Kiếm Nhị Phân (Binary Search)** 

## **11. Binary Search Cơ Bản** 

_Tìm phần tử trong mảng đã được sắp xếp._ 

```go
func search(nums []int, target int) int {
    l, r := 0, len(nums)-1
    for l <= r {
        m := l + (r-l)/2
        if nums[m] == target { return m }
        if nums[m] < target { l = m + 1 } else { r = m - 1 }
    }
    return -1
}
```

Time: _**O(log N)**_ 

Trang 7 / 13 

Tài liệu rèn luyện thuật toán & tư duy hệ thống 

## **12. Search a 2D Matrix (Tìm kiếm ma trận 2D)** 

_Tìm một giá trị trong ma trận m x n, trong đó mỗi hàng được sắp xếp và phần tử đầu hàng lớn hơn phần tử cuối hàng trước._ 

**Ý tưởng:** Chuyển chỉ số 1D thành tọa độ 2D: _**row = mid / cols**_ , _**col = mid % cols**_ . 

```go
func searchMatrix(matrix [][]int, target int) bool {
    rows, cols := len(matrix), len(matrix[0])
    l, r := 0, rows*cols-1
    for l <= r {
        m := l + (r-l)/2
        val := matrix[m/cols][m%cols]
        if val == target { return true }
        if val < target { l = m + 1 } else { r = m - 1 }
    }
    return false
}
```

Time: _**O(log(M*N))**_ 

## **Phần 5: Danh Sách Liên Kết & Ứng Dụng Thực Tế (Linked Lists & Caching)** 

## **13. Merge Two Sorted Lists** 

_Gộp hai danh sách liên kết đã sắp xếp thành một danh sách liên kết mới cũng được sắp xếp._ 

```go
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
    dummy := &ListNode{}
    curr := dummy
    for l1 != nil && l2 != nil {
        if l1.Val < l2.Val { 
            curr.Next = l1
            l1 = l1.Next 
        } else { 
            curr.Next = l2
            l2 = l2.Next 
        }
        curr = curr.Next
    }
    if l1 != nil { curr.Next = l1 }
    if l2 != nil { curr.Next = l2 }
    return dummy.Next
}
```

Trang 8 / 13 

Tài liệu rèn luyện thuật toán & tư duy hệ thống 

## **14. Linked List Cycle (Phát hiện chu trình)** 

_Kiểm tra xem danh sách liên kết có chứa chu trình hay không._ 

**Ý tưởng:** Thuật toán Rùa và Thỏ (Floyd's). Một con trỏ nhảy 1 bước, một con trỏ nhảy 2 bước. Nếu có chu trình, chúng sẽ gặp nhau. 

```go
func hasCycle(head *ListNode) bool {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast { return true }
    }
    return false
}
```

Trang 9 / 13 

Tài liệu rèn luyện thuật toán & tư duy hệ thống 

## **15. LRU Cache (Hệ thống Caching Cache Least Recently Used)** 

_Thiết kế cấu trúc dữ liệu cho bộ nhớ đệm LRU. Được áp dụng rộng rãi trong thiết kế cơ sở dữ liệu và hạ tầng mạng (như Redis cache policies). Cần thực hiện Get và Put trong thời gian_ _**O(1)** ._ 

**Ý tưởng:** Kết hợp một Bảng băm (để tra cứu nhanh _**O(1)**_ ) và Danh sách liên kết đôi (để đẩy các mục mới/được truy cập lên đầu và loại bỏ mục ở cuối nhanh chóng). 

```go
type Node struct { 
    key, val int
    prev, next *Node 
}

type LRUCache struct {
    capacity int
    cache map[int]*Node
    head, tail *Node
}

func Constructor(capacity int) LRUCache {
    h, t := &Node{}, &Node{}
    h.next, t.prev = t, h
    return LRUCache{capacity, make(map[int]*Node), h, t}
}

func (this *LRUCache) remove(node *Node) {
    node.prev.next, node.next.prev = node.next, node.prev
}

func (this *LRUCache) insert(node *Node) { // Chèn vào ngay sau head
    node.next, node.prev = this.head.next, this.head
    this.head.next.prev, this.head.next = node, node
}

func (this *LRUCache) Get(key int) int {
    if node, ok := this.cache[key]; ok {
        this.remove(node)
        this.insert(node)
        return node.val
    }
    return -1
}

func (this *LRUCache) Put(key int, value int) {
    if node, ok := this.cache[key]; ok {
        this.remove(node)
    }
    newNode := &Node{key: key, val: value}
    this.cache[key] = newNode
    this.insert(newNode)
    if len(this.cache) > this.capacity {
        lru := this.tail.prev
        this.remove(lru)
        delete(this.cache, lru.key)
    }
}
```

Trang 10 / 13 

Tài liệu rèn luyện thuật toán & tư duy hệ thống 

## **Phần 6: Cây (Trees)** 

## **16. Invert Binary Tree (Đảo ngược cây nhị phân)** 

_Đảo ngược các nhánh trái/phải của mọi nút trên cây._ 

```go
func invertTree(root *TreeNode) *TreeNode {
    if root == nil { return nil }
    root.Left, root.Right = root.Right, root.Left
    invertTree(root.Left)
    invertTree(root.Right)
    return root
}
```

## **17. Lowest Common Ancestor of a BST (Tổ tiên chung gần nhất)** 

_Tìm tổ tiên chung sâu nhất của hai nút p và q trên Cây tìm kiếm nhị phân._ 

**Ý tưởng:** Lợi dụng tính chất của BST: Nút gốc nằm giữa p và q chính là tổ tiên chung gần nhất. 

```go
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
    curr := root
    for curr != nil {
        if p.Val > curr.Val && q.Val > curr.Val { 
            curr = curr.Right 
        } else if p.Val < curr.Val && q.Val < curr.Val { 
            curr = curr.Left 
        } else { 
            return curr 
        }
    }
    return nil
}
```

Trang 11 / 13 

Tài liệu rèn luyện thuật toán & tư duy hệ thống 

## **Phần 7: Đồ Thị & Ma Trận (Graphs & Grids)** 

## **18. Number of Islands (Đếm số hòn đảo)** 

_Cho ma trận 2D chỉ chứa '1' (đất) và '0' (nước). Tính số cụm đất liền kề nhau. Bài toán có tính chất tương tự như tìm vùng liên thông trong kiến trúc hạ tầng mạng._ 

**Ý tưởng:** Duyệt qua từng ô. Khi gặp '1', đánh dấu nó là hòn đảo mới và dùng DFS để chìm toàn bộ các phần đất liền kề (đổi thành '0') để không đếm lại. 

```go
func numIslands(grid [][]byte) int {
    if len(grid) == 0 { return 0 }
    count := 0
    var dfs func(r, c int)
    dfs = func(r, c int) {
        if r < 0 || c < 0 || r >= len(grid) || c >= len(grid[0]) || grid[r][c] == '0' { 
            return 
        }
        grid[r][c] = '0'
        dfs(r+1, c); dfs(r-1, c); dfs(r, c+1); dfs(r, c-1)
    }
    for r := range grid {
        for c := range grid[r] {
            if grid[r][c] == '1' { count++; dfs(r, c) }
        }
    }
    return count
}
```

## **19. Max Area of Island (Diện tích đảo lớn nhất)** 

_Giống bài trên nhưng yêu cầu trả về diện tích lớn nhất._ 

```go
func maxAreaOfIsland(grid [][]int) int {
    maxArea := 0
    var dfs func(r, c int) int
    dfs = func(r, c int) int {
        if r < 0 || c < 0 || r >= len(grid) || c >= len(grid[0]) || grid[r][c] == 0 { 
            return 0 
        }
        grid[r][c] = 0 // visited
        return 1 + dfs(r+1, c) + dfs(r-1, c) + dfs(r, c+1) + dfs(r, c-1)
    }
    for r := range grid {
        for c := range grid[r] {
            if grid[r][c] == 1 {
                area := dfs(r, c)
                if area > maxArea { maxArea = area }
            }
        }
    }
    return maxArea
}
```

Trang 12 / 13 

Tài liệu rèn luyện thuật toán & tư duy hệ thống 

## **Phần 8: Quy Hoạch Động Chuyển Dạng Dãy Số (Dynamic Programming & Sequences)** 

## **20. Climbing Stairs (Bài toán dãy số và tổ hợp)** 

_Bạn đang leo cầu thang n bậc. Mỗi bước bạn có thể đi 1 hoặc 2 bậc. Có bao nhiêu cách để lên đỉnh?_ 

**Ý tưởng:** Về mặt toán học, đây là bài toán cấp số cộng dạng Dãy Fibonacci. Số cách để tới bậc _**i**_ bằng tổng số cách tới bậc _**i-1**_ và _**i-2**_ . Thay vì dùng mảng _**O(N)**_ , ta chỉ cần lưu 2 giá trị biến đổi trạng thái để tiết kiệm bộ nhớ. 

```go
func climbStairs(n int) int {
    one, two := 1, 1
    for i := 0; i < n-1; i++ {
        temp := one
        one = one + two
        two = temp
    }
    return one
}
```

Time: _**O(N)**_ | Space: _**O(1)**_ 

Trang 13 / 13 
