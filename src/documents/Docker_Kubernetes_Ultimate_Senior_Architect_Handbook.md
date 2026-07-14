# KIẾN THỨC CONTAINERIZATION & ORCHESTRATION TOÀN DIỆN

_Cẩm nang Chuyên sâu cấp độ Senior & Architect Engineer (Bản Đầy Đủ Nhất)_

## **CHƯƠNG I: BẢN CHẤT HỆ ĐIỀU HÀNH & KHỞI NGUYÊN CỦA CONTAINER**

## **1. Cơ chế cô lập tầng thấp của Linux Kernel**

Một Senior Engineer nhìn nhận container không phải là một thực thể ảo hóa, mà là các tiến trình Linux thông thường bị áp đặt các ranh giới chặt chẽ thông qua các tính năng nguyên bản của nhân Linux:

- **Namespaces (Cô lập góc nhìn):** Quyết định tiến trình nhìn thấy cái gì.
  - `pid` (Process ID): Container có tiến trình riêng bắt đầu từ PID 1, hoàn toàn tách biệt với cây tiến trình của host.
  - `net` (Network): Cấp cấu trúc mạng riêng (card mạng ảo, bảng routing, iptables, cổng port độc lập).
  - `mnt` (Mount): Tạo một cây thư mục (file system root) riêng, không thấy file hệ thống của host trừ khi được mount vào.
  - `ipc` (Inter-Process Communication): Cô lập các cơ chế giao tiếp liên tiến trình như System V IPC, POSIX message queues.
  - `uts` (UNIX Timesharing System): Cho phép container sở hữu hostname và NIS domain name riêng.
  - `user` (User IDs): Map UID/GID bên trong container ra UID/GID khác trên host (ví dụ: root UID 0 trong container map với UID 10001 không có quyền trên host), ngăn ngừa container breakout.
  - `cgroup`: Cô lập góc nhìn đối với chính cấu trúc cgroup.

- **Control Groups - cgroups (Giới hạn tài nguyên):** Quyết định tiến trình được dùng bao nhiêu tài nguyên. Gồm 2 phiên bản chính:
  - `cgroups v1`: Quản lý tài nguyên theo các hệ thống phân cấp tách biệt cho từng loại (cpu, memory, blkio).
  - `cgroups v2`: Hệ thống phân cấp hợp nhất, giải quyết triệt để vấn đề tranh chấp tài nguyên chéo, quản lý tài nguyên Memory OOM hiệu quả hơn và hỗ trợ tốt các tiến trình non-root (Rootless containers).

- **Capabilities:** Chia nhỏ đặc quyền tối thượng của user `root` thành các đơn vị quyền nhỏ hơn (ví dụ: `CAP_NET_ADMIN` để cấu hình mạng, `CAP_SYS_ADMIN` cho các tác vụ quản trị). Mặc định Docker tắt hầu hết các capabilities nguy hiểm để bảo vệ host.

- **Seccomp (Secure Computing Mode):** Bộ lọc các system call (lời gọi hàm hệ thống). Docker áp dụng một profile mặc định chặn hơn 40 system calls nguy hiểm (như `reboot`, `swapon`) để tránh container can thiệp trực tiếp vào kernel của host máy vật lý.

## **2. Sự Phân Rã và Tiến Hóa của Container Runtimes**

Để đảm bảo tính nhất quán toàn ngành, chuẩn OCI (Open Container Initiative) ra đời, chia runtime làm 2 tầng rõ rệt:

- **Low-level Runtime (runc, crun, gVisor, Kata Containers):** Trực tiếp thực thi việc tạo namespaces/cgroups.
  - `runc`: Bản thực thi tiêu chuẩn mặc định bằng Go.
  - `gVisor / Kata Containers`: Các giải pháp bảo mật nâng cao (Sandboxed Containers). gVisor dịch các system call qua một proxy kernel viết bằng Go, trong khi Kata sử dụng các máy ảo siêu nhỏ (MicroVM) để bọc container lại, cô lập tuyệt đối cho môi trường Multi-tenant.

- **High-level Runtime (containerd, CRI-O):** Quản lý vòng đời cấp cao bao gồm: kéo ảnh (pull image), quản lý lưu trữ tầng overlay, quản lý mạng logic, cung cấp gRPC API cho các thành phần orchestrator gọi xuống.

## **CHƯƠNG II: THIẾT KẾ, TỐI ƯU HÓA & VẬN HÀNH DOCKER TRONG PRODUCTION**

## **1. Cơ chế Lưu trữ: Copy-on-Write & Overlay2**

Docker Image cấu thành từ các Layer chỉ đọc (Read-Only Layers). Khi container khởi chạy, Docker Engine phủ thêm một Layer đọc-ghi mỏng (Writable Layer/Container Layer) lên trên cùng bằng **Overlay2 Storage Driver**.

- **Copy-on-Write (CoW):** Khi một tiến trình trong container muốn sửa đổi một file thuộc Layer chỉ đọc bên dưới, tập tin đó sẽ được sao chép từ tầng dưới lên Writable Layer rồi mới tiến hành chỉnh sửa. Điều này giúp tiết kiệm tài nguyên ổ đĩa tối đa vì nhiều container có thể dùng chung các lớp ảnh gốc.

- **Chiến lược I/O cho Production:** Tuyệt đối không chạy các ứng dụng có I/O nặng (Database như PostgreSQL, MySQL, Redis, Kafka) trực tiếp trên Writable Layer của container vì hiệu năng CoW rất kém. Bắt buộc sử dụng `Named Volumes` hoặc `Bind Mounts` để bypass qua hệ thống overlay, ghi dữ liệu trực tiếp xuống ổ cứng máy host với tốc độ native.

## **2. Kiến trúc Mạng Docker nâng cao**

- **Bridge Network:** Tạo một switch ảo (thường tên là `docker0`). Các container kết nối qua cặp veth (virtual ethernet). Docker sử dụng `iptables` trên host để làm NAT (Network Address Translation) và Port Forwarding khi cấu hình tham số `-p host_port:container_port`.

- **Host Network:** Container dùng chung hoàn toàn network namespace với host. Không có NAT, hiệu năng mạng cao nhất, nhưng dễ xung đột port và mất tính cô lập an toàn.

- **Overlay Network:** Sử dụng giao thức VXLAN (Virtual Extensible LAN) để đóng gói các gói tin tầng 2 vào trong các gói tin tầng 3 UDP, tạo nên một mạng phẳng ảo kết nối các container nằm trên các máy vật lý (máy host) khác nhau.

## **3. Kỹ thuật viết Dockerfile cấp độ Senior**

Một bản thiết kế Dockerfile tối ưu cho môi trường Enterprise cần đáp ứng 3 tiêu chí: Tốc độ build (nhờ Cache), Dung lượng tối thiểu (giảm chi phí lưu trữ/băng thông mạng) và Bảo mật tuyệt đối.

```dockerfile
# STAGE 1: Phát triển và Xây dựng Artifact (Heavyweight Image)
FROM node:20-alpine AS builder
WORKDIR /usr/src/app
# Tận dụng tối đa Layer Cache cho Dependencies
COPY package*.json ./
RUN npm ci --silent
# Copy toàn bộ mã nguồn và thực hiện biên dịch/đóng gói
COPY . .
RUN npm run build

# STAGE 2: Môi trường chạy Production (Minimalist Runtime Image)
FROM node:20-alpine AS runner
WORKDIR /usr/share/app
# Thiết lập biến môi trường tối ưu hiệu năng ứng dụng
ENV NODE_ENV=production
# Chỉ cài đặt các package cần thiết cho runtime, loại bỏ devDependencies
COPY package*.json ./
RUN npm ci --only=production --silent
# Sao chép kết quả đã biên dịch từ Stage trước sang
COPY --from=builder /usr/src/app/dist ./dist
# Tạo User hệ thống không có quyền đặc quyền cao (Non-root) để vận hành app
RUN addgroup -g 1001 -S node_group && adduser -u 1001 -S node_user -G node_group
USER node_user
# Cấu hình tín hiệu dừng chuẩn xác để tránh zombie process
STOPSIGNAL SIGTERM
EXPOSE 3000
CMD ["node", "dist/main.js"]
```

## **CHƯƠNG III: KIẾN TRÚC CHUYÊN SÂU CỦA HỆ THỐNG KUBERNETES (K8S INTERNALS)**

## **1. Control Plane Components & Cơ chế đồng thuận**

- **API Server (kube-apiserver):** Điểm giao tiếp duy nhất. Nó hoạt động theo mô hình không trạng thái (stateless), thực hiện tuần tự: Xác thực (Authentication) -> Phân quyền (Authorization - RBAC) -> Kiểm soát truy cập (Admission Control - ví dụ Mutating và Validating Webhooks để tự động chỉnh sửa hoặc kiểm tra tính hợp lệ của manifest trước khi ghi nhận).

- **etcd (Bộ não dữ liệu):** Cơ sở dữ liệu phân tán sử dụng thuật toán đồng thuận **Raft**.
  - Mọi dữ liệu ghi vào K8s đều được lưu tại đây. Chỉ có API Server mới có quyền đọc ghi trực tiếp vào etcd.
  - Trong Production, cụm etcd phải có số lượng node lẻ (3, 5, 7) để tránh tình trạng phân rã não bộ (split-brain). Tốc độ ổ đĩa (IOPS) của node chứa etcd quyết định trực tiếp đến độ ổn định của toàn cụm cluster.

- **Kube Scheduler:** Lựa chọn node tối ưu để chạy Pod qua 2 giai đoạn:
  - _Filtering (Predicates):_ Loại bỏ các node không đủ điều kiện (thiếu CPU/RAM, sai NodeSelector, dính Taints).
  - _Scoring (Priorities):_ Chấm điểm các node còn lại dựa trên các tiêu chí (Ví dụ: ưu tiên node đã có sẵn image, node có lượng tài nguyên phân bổ cân bằng). Node cao điểm nhất sẽ được chọn.

- **Kube Controller Manager:** Tập hợp của nhiều luồng kiểm soát (control loops). Mỗi controller liên tục lắng nghe sự thay đổi trạng thái thông qua cơ chế *Informer/Watch API* của API Server để thực thi hành động đưa trạng thái thực tế về đúng trạng thái mong muốn (Desired State).

## **2. Worker Node Components & Cơ chế Mạng Pod-to-Pod**

- **Kubelet:** Giao tiếp với API Server để nhận PodSpec. Sau đó, thông qua **CRI (Container Runtime Interface)** để yêu cầu container runtime khởi tạo container. Kubelet cũng trực tiếp thực hiện kiểm tra sức khỏe ứng dụng (Probes) và báo cáo trạng thái Node về master.

- **Kube-Proxy:** Quản lý mạng dịch vụ (Service). Có 2 chế độ chính:
  - _iptables mode:_ Kube-proxy ghi các luật iptables lên node. Nhược điểm: Khi cụm lên đến hàng ngàn service, iptables duyệt tuần tự từ trên xuống dưới, gây suy giảm hiệu năng mạng nghiêm trọng.
  - _IPVS mode (IP Virtual Server):_ Sử dụng bảng băm (Hash table) tầng kernel của Netfilter. Tốc độ tìm kiếm là O(1), đáp ứng hoàn hảo cho các cụm cluster siêu lớn ở quy mô Enterprise.

## **CHƯƠNG IV: QUẢN TRỊ TÀI NGUYÊN, CHIẾN LƯỢC TRIỂN KHAI & VẬN HÀNH KHỎE MẠNH**

## **1. Mô hình Tài Nguyên Phức Tạp: Requests, Limits & QoS (Quality of Service)**

K8s phân loại Pod vào 3 nhóm chất lượng dịch vụ (QoS) dựa trên cách khai báo tài nguyên, quyết định thứ tự bị "hiến tế" khi Node bị cạn kiệt tài nguyên vật lý:

| **Lớp QoS** | **Điều kiện cấu hình** | **Độ ưu tiên & Hành vi khi Node cạn kiệt tài nguyên** |
|---|---|---|
| **Guaranteed** | Tất cả các container trong Pod đều có `Requests == Limits` cho cả CPU và Memory. | Độ ưu tiên cao nhất. Chỉ bị kill khi hệ thống cạn kiệt tài nguyên nghiêm trọng và không còn Pod nào khác để giải phóng. |
| **Burstable** | Có cấu hình Requests và Limits nhưng chúng không bằng nhau, hoặc chỉ cấu hình một trong hai. | Độ ưu tiên trung bình. Có thể dùng vượt mức request lên tới mức limit nếu node còn trống tài nguyên. |
| **BestEffort** | Hoàn toàn không khai báo bất kỳ chỉ số Requests hay Limits nào. | Độ ưu tiên thấp nhất. Sẽ bị K8s tiêu diệt (Evicted/OOMKilled) ngay lập tức khi Node có dấu hiệu quá tải tài nguyên. |

## **2. Cơ chế Tự động Co Giãn (Autoscaling)**

- **HPA (Horizontal Pod Autoscaler):** Tăng/giảm số lượng bản sao (Pod replicas) dựa trên chỉ số CPU, Memory tiêu thụ hoặc Custom Metrics (như số lượng requests/giây lấy từ Prometheus).

- **VPA (Vertical Pod Autoscaler):** Tự động tính toán và điều chỉnh lại mức định lượng Requests/Limits CPU và Memory của Pod. Lưu ý: VPA mặc định sẽ restart lại Pod để áp dụng cấu hình tài nguyên mới.

- **Cluster Autoscaler:** Tự động tương tác với hạ tầng Cloud (AWS ASG, GCP Node Pools) để thêm Node vật lý mới vào cụm khi thấy có nhiều Pod rơi vào trạng thái `Pending` do toàn bộ cluster không còn Node nào đủ tài nguyên đáp ứng.

## **3. Thiết kế Health Checks Chuẩn Chỉ**

Sai lầm phổ biến là cấu hình Liveness Probe trỏ chung vào một API endpoint kiểm tra Database với Readiness Probe. Khi DB sập, Liveness Probe fail khiến K8s restart đồng loạt hàng loạt Pod, tạo ra hiệu ứng sập dây chuyền (Cascading Failure).

- **Startup Probe:** Dành cho ứng dụng mất nhiều thời gian khởi chạy (ví dụ nạp cấu hình nặng ban đầu). Nó block hoàn toàn liveness và readiness probe hoạt động cho tới khi hoàn tất.

- **Liveness Probe:** Chỉ kiểm tra xem tiến trình ứng dụng nội tại có bị treo cứng hay deadlock không. Trỏ vào endpoint gọn nhẹ trả về HTTP 200 nội bộ (ví dụ: `/healthz`).

- **Readiness Probe:** Kiểm tra xem ứng dụng đã sẵn sàng xử lý traffic chưa (đã kết nối thành công tới DB, Redis, Kafka chưa). Nếu fail, tạm thời ngắt Pod khỏi Service để không làm lỗi request của người dùng cuối.

## **CHƯƠNG V: MẠNG NÂNG CAO, LƯU TRỮ BỀN VỮNG & BẢO MẬT ENTERPRISE**

## **1. Ingress, Ingress Controller & Service Mesh**

Service nội cụm chỉ cấp IP nội bộ. Để đưa traffic từ internet vào, ta sử dụng cấu trúc tầng:

- **Ingress Resource:** Bản khai báo các luật định tuyến (Routing rules) dạng HTTP/HTTPS (ví dụ: domain `api.domain.com` trỏ vào Service A).

- **Ingress Controller (NGINX Ingress, Traefik, HAProxy):** Một Pod chạy trong cụm, đóng vai trò là Reverse Proxy thu nhận cấu hình từ các Ingress Resource để thực thi việc định tuyến traffic thực tế.

- **Service Mesh (Istio, Linkerd):** Quản lý giao tiếp nâng cao dạng Service-to-Service (East-West traffic). Cài đặt một Sidecar Proxy (Envoy) đứng cạnh mỗi Pod để thực thi: Mã hóa toàn bộ traffic bằng **mTLS**, Phân tách traffic (Canary deployment), Rate Limiting, Circuit Breaking và Thu thập chỉ số quan sát trực quan (Observability) một cách tự động.

## **2. Hệ thống Lưu trữ trạng thái trong K8s (Stateful Workloads)**

Đối với ứng dụng cần lưu dữ liệu bền vững, K8s cung cấp mô hình phân tách vai trò:

- **StorageClass (SC):** Do Admin định nghĩa trước các loại ổ đĩa vật lý (ví dụ: SSD tốc độ cao trên Cloud hoặc hệ thống Ceph, NFS ở On-premise).

- **PersistentVolumeClaim (PVC):** Do Developer viết để yêu cầu một lượng không gian lưu trữ và chế độ truy cập (như `ReadWriteOnce` - chỉ 1 node mount được, hoặc `ReadWriteMany` - nhiều node cùng đồng thời mount được).

- **PersistentVolume (PV):** Thực thể lưu trữ vật lý được hệ thống tự động khởi tạo (Dynamic Provisioning) dựa trên yêu cầu từ PVC và StorageClass để gắn trực tiếp vào Pod.

- **StatefulSet:** Controller chuyên dụng cho các ứng dụng Stateful (Database). Khác với Deployment, các Pod trong StatefulSet có định danh cố định (ví dụ: `mysql-0`, `mysql-1`), có thứ tự khởi chạy/sập tuần tự rõ ràng và mỗi Pod được gắn chặt với một PV độc lập không thay đổi ngay cả khi Pod bị restart sang Node khác.

## **3. Chiến lược Bảo mật Toàn diện (Hardening K8s Cluster)**

- **RBAC (Role-Based Access Control):** Không dùng chung quyền. Chia rõ `Role` (trong 1 namespace) và `ClusterRole` (toàn cụm), sau đó gán cho đối tượng bằng `RoleBinding` hoặc `ClusterRoleBinding`.

- **Pod Security Standards (PSS) & Admission Controllers (Kyverno / OPA Gatekeeper):** Thay thế cho PSP cũ. Áp đặt các rule bắt buộc: Không cho phép chạy container bằng quyền root (`MustRunAsNonRoot`), Cấm mount trực tiếp root filesystem của host (`readOnlyRootFilesystem`), và ngăn chặn leo thang đặc quyền.

- **Network Policies (Tường lửa nội bộ):** Cấu hình mô hình Zero-Trust. Mặc định cô lập toàn bộ traffic (Default Deny), chỉ mở các kết nối được chỉ định tường minh thông qua nhãn (Labels Selector).

---

**TƯ DUY ARCHITECT TRONG VẬN HÀNH PRODUCTION:** Hãy chuyển dịch toàn bộ mô hình quản trị sang **GITOPS** sử dụng ArgoCD. Loại bỏ hoàn toàn quyền can thiệp thủ công vào cluster. Mọi cấu hình từ hạ tầng, mạng đến ứng dụng đều phải lưu trữ dưới dạng mã nguồn mã hóa trên Git. Kết hợp với bộ công cụ giám sát **LGTM STACK (LOKI, GRAFANA, TEMPO, MIMIR)** hoặc **PROMETHEUS + DATADOG** để có cái nhìn toàn diện từ hạ tầng cho tới log chi tiết của từng tiến trình.
