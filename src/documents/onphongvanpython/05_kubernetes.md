# 5. Kubernetes cho Backend Developer

Kubernetes là nền tảng orchestration dùng để triển khai, mở rộng, cập nhật và tự phục hồi các ứng dụng container trên một cluster máy chủ.

## 1. Vì sao cần Kubernetes?

Docker giúp đóng gói và chạy một container. Khi hệ thống có nhiều service, nhiều instance và nhiều máy chủ, ta cần giải quyết:

- Container nên chạy trên node nào?
- Container chết thì ai tạo lại?
- Làm sao scale từ 3 lên 20 instance?
- Làm sao cập nhật không downtime?
- Làm sao service tìm được nhau khi IP Pod thay đổi?
- Làm sao quản lý cấu hình, secret và tài nguyên?

Kubernetes giải quyết các vấn đề này bằng mô hình desired state: ta khai báo trạng thái mong muốn, control plane liên tục cố gắng đưa trạng thái thực tế về đúng trạng thái đó.

```text
Internet
   ↓
Ingress / Load Balancer
   ↓
Service
   ↓
Deployment
   ↓
Pod API 1   Pod API 2   Pod API 3
   ↓
PostgreSQL / Redis / RabbitMQ
```

---

## 2. Kiến trúc cluster

### Control plane

Các thành phần chính:

- API Server: cổng giao tiếp trung tâm của cluster.
- Scheduler: chọn node để chạy Pod.
- Controller Manager: chạy các controller để reconcile desired state.
- etcd: lưu trạng thái cluster dạng key-value.

### Worker node

- kubelet: agent trên node, đảm bảo Pod được chạy.
- container runtime: chạy container, ví dụ containerd.
- kube-proxy hoặc data plane tương đương: hỗ trợ network và Service.

Khi tạo Deployment, request đi vào API Server, được lưu trong etcd. Controller tạo ReplicaSet, ReplicaSet tạo Pod, Scheduler gán Pod vào node, kubelet trên node kéo image và chạy container.

---

## 3. Pod

Pod là đơn vị triển khai nhỏ nhất trong Kubernetes. Một Pod có thể chứa một hoặc nhiều container dùng chung network namespace và volume.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: python-api
spec:
  containers:
    - name: api
      image: ghcr.io/example/python-api:1.0.0
      ports:
        - containerPort: 8000
```

Pod thường là tạm thời. Khi bị thay thế, Pod mới có IP mới. Vì vậy không nên kết nối trực tiếp bằng Pod IP; hãy dùng Service.

### Khi nào một Pod có nhiều container?

Ví dụ sidecar:

- Application container.
- Log collector.
- Proxy hoặc service mesh sidecar.

Không nên gom các service độc lập vào cùng Pod chỉ vì muốn deploy cùng nhau.

---

## 4. Deployment và ReplicaSet

Deployment quản lý rollout và số lượng Pod của ứng dụng stateless.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: python-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: python-api
  template:
    metadata:
      labels:
        app: python-api
    spec:
      containers:
        - name: api
          image: ghcr.io/example/python-api:1.0.0
          ports:
            - containerPort: 8000
```

Deployment tạo ReplicaSet. ReplicaSet đảm bảo luôn có đúng số Pod mong muốn.

Nếu một Pod chết:

```text
replicas mong muốn = 3
replicas thực tế = 2
→ ReplicaSet tạo Pod mới
```

Không nên chỉnh ReplicaSet do Deployment quản lý trực tiếp.

---

## 5. Service

Service cung cấp địa chỉ ổn định và load balancing đến các Pod có label phù hợp.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: python-api
spec:
  selector:
    app: python-api
  ports:
    - port: 80
      targetPort: 8000
  type: ClusterIP
```

Các loại phổ biến:

- ClusterIP: truy cập nội bộ cluster.
- NodePort: mở port trên mọi node.
- LoadBalancer: yêu cầu cloud hoặc implementation load balancer.
- ExternalName: ánh xạ DNS đến tên bên ngoài.

Request đến Service được chuyển đến một endpoint sẵn sàng. Pod mới có IP khác nhưng Service vẫn giữ DNS ổn định.

Ứng dụng khác gọi:

```text
http://python-api
```

Trong namespace khác:

```text
python-api.backend.svc.cluster.local
```

---

## 6. Ingress

Ingress định tuyến HTTP/HTTPS từ bên ngoài vào Service.

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: backend-ingress
spec:
  rules:
    - host: api.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: python-api
                port:
                  number: 80
```

Ingress resource chỉ là cấu hình. Cluster cần Ingress Controller như NGINX Ingress, Traefik hoặc implementation cloud.

---

## 7. ConfigMap và Secret

### ConfigMap

Dùng cho cấu hình không nhạy cảm.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: api-config
data:
  APP_ENV: production
  LOG_LEVEL: info
```

### Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: api-secret
type: Opaque
stringData:
  DATABASE_URL: postgresql://user:password@postgres:5432/app
```

Secret trong manifest không tự động an toàn. Base64 chỉ là encoding. Production nên kết hợp encryption at rest, RBAC và secret manager như Vault hoặc cloud secret manager.

Inject vào container:

```yaml
envFrom:
  - configMapRef:
      name: api-config
  - secretRef:
      name: api-secret
```

Khi ConfigMap hoặc Secret thay đổi, biến môi trường trong Pod đang chạy không tự cập nhật; thường cần rollout restart.

---

## 8. Liveness, readiness và startup probe

### Readiness probe

Xác định Pod đã sẵn sàng nhận traffic chưa.

```yaml
readinessProbe:
  httpGet:
    path: /health/ready
    port: 8000
  initialDelaySeconds: 5
  periodSeconds: 5
```

Nếu readiness thất bại, Pod vẫn chạy nhưng bị loại khỏi endpoint của Service.

### Liveness probe

Xác định container có bị treo và cần restart không.

```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 8000
  periodSeconds: 10
```

Không nên để liveness phụ thuộc trực tiếp database. Database tạm lỗi có thể khiến toàn bộ API restart liên tục.

### Startup probe

Dành cho ứng dụng khởi động lâu. Trong thời gian startup probe chưa thành công, liveness chưa được áp dụng.

Thiết kế endpoint:

- `/health/live`: process và event loop còn hoạt động.
- `/health/ready`: instance có thể nhận request; có thể kiểm tra dependency quan trọng với timeout ngắn.

---

## 9. Resource requests và limits

```yaml
resources:
  requests:
    cpu: "250m"
    memory: "256Mi"
  limits:
    cpu: "1"
    memory: "512Mi"
```

- Request dùng cho scheduling và đảm bảo tài nguyên tương đối.
- Limit là ngưỡng tối đa.
- Vượt memory limit thường bị `OOMKilled`.
- CPU limit thường gây throttling thay vì kill.

Không đặt request quá thấp chỉ để Pod dễ schedule. Cần đo usage thực tế.

---

## 10. Rolling update và rollback

```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxUnavailable: 0
    maxSurge: 1
```

Luồng rollout:

1. Deployment tạo Pod dùng image mới.
2. Pod mới vượt startup và readiness probe.
3. Service bắt đầu gửi traffic.
4. Pod cũ được giảm dần.

Readiness probe sai có thể đưa traffic vào Pod chưa migration xong hoặc chưa kết nối dependency.

Lệnh thường dùng:

```bash
kubectl rollout status deployment/python-api
kubectl rollout history deployment/python-api
kubectl rollout undo deployment/python-api
```

Database migration cần backward-compatible vì trong rolling update có lúc code cũ và mới chạy đồng thời.

---

## 11. StatefulSet, DaemonSet và Job

### StatefulSet

Phù hợp workload cần identity và storage ổn định, ví dụ một số database hoặc broker.

### DaemonSet

Chạy một Pod trên mỗi node, ví dụ log agent hoặc node monitoring.

### Job

Chạy tác vụ đến khi hoàn thành.

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: database-migration
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: migrate
          image: ghcr.io/example/python-api:1.0.0
          command: ["alembic", "upgrade", "head"]
```

### CronJob

Chạy Job theo lịch, ví dụ báo cáo hằng ngày hoặc cleanup.

---

## 12. PersistentVolume và PersistentVolumeClaim

Pod filesystem là tạm thời. Dữ liệu cần tồn tại phải dùng volume phù hợp.

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: app-data
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

Ứng dụng stateless nên lưu file lên object storage thay vì phụ thuộc local disk của Pod.

---

## 13. Horizontal Pod Autoscaler

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: python-api
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: python-api
  minReplicas: 3
  maxReplicas: 20
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

CPU không phải lúc nào cũng phản ánh tải. Worker queue có thể scale theo queue length; API I/O-bound có thể cần metric request rate hoặc latency.

HPA cần resource requests để tính utilization chính xác.

---

## 14. Namespace, RBAC và ServiceAccount

Namespace phân tách logic tài nguyên, ví dụ `dev`, `staging`, `production` hoặc theo team.

RBAC kiểm soát ai hoặc workload nào được phép làm gì.

Nguyên tắc:

- Least privilege.
- Không cấp `cluster-admin` cho application.
- Mỗi workload quan trọng có ServiceAccount riêng.
- Hạn chế quyền đọc Secret.

---

## 15. Debug sự cố

### Pod Pending

Kiểm tra:

```bash
kubectl describe pod <pod-name>
```

Nguyên nhân thường gặp:

- Không đủ CPU/RAM.
- PVC chưa bind.
- Node selector hoặc affinity không khớp.
- Taint chưa có toleration.

### ImagePullBackOff

- Sai image/tag.
- Registry private nhưng thiếu imagePullSecret.
- Node không kết nối được registry.

### CrashLoopBackOff

```bash
kubectl logs <pod-name>
kubectl logs <pod-name> --previous
kubectl describe pod <pod-name>
```

Nguyên nhân:

- Application crash.
- Sai environment variable.
- Migration thất bại.
- Liveness probe quá gắt.
- Permission hoặc filesystem read-only.

### Service không truy cập được Pod

Kiểm tra label và endpoint:

```bash
kubectl get pods --show-labels
kubectl get service python-api
kubectl get endpoints python-api
```

Nếu endpoint rỗng, selector không khớp hoặc Pod chưa ready.

### Debug network

```bash
kubectl exec -it <pod-name> -- sh
kubectl run debug --rm -it --image=curlimages/curl -- sh
```

Kiểm tra DNS, port, NetworkPolicy và service name.

---

## 16. Best practices cho Python Backend

- Chạy process dưới non-root user.
- Image có tag bất biến; tránh `latest` ở production.
- Xử lý SIGTERM và graceful shutdown.
- Có readiness, liveness và startup probe phù hợp.
- Đặt request/limit dựa trên đo lường.
- Log ra stdout/stderr.
- Không lưu state quan trọng trong filesystem của Pod.
- Migration backward-compatible.
- Dùng PodDisruptionBudget khi cần đảm bảo availability.
- Có nhiều replica trải trên node/zone khi phù hợp.
- Không đưa secret vào image hoặc Git.

Python API cần dừng nhận request, hoàn thành request đang chạy và đóng connection pool khi nhận SIGTERM.

---

## 17. Câu hỏi phỏng vấn

### Pod khác container thế nào?

Container là process được cô lập. Pod là đơn vị Kubernetes quản lý, có thể chứa một hoặc nhiều container chia sẻ network và volume.

### Deployment khác StatefulSet?

Deployment phù hợp workload stateless, Pod có thể thay thế tự do. StatefulSet cung cấp identity, thứ tự và storage ổn định cho workload stateful.

### Service hoạt động thế nào khi Pod đổi IP?

Service chọn Pod bằng label selector. Endpoint controller cập nhật danh sách Pod IP, còn DNS và ClusterIP của Service vẫn ổn định.

### Readiness khác liveness?

Readiness quyết định Pod có nhận traffic không. Liveness quyết định container có cần restart không. Readiness fail không nhất thiết restart Pod.

### Vì sao ứng dụng có 3 replica vẫn downtime khi deploy?

Có thể do readiness sai, `maxUnavailable` quá cao, migration phá tương thích, tất cả Pod cùng phụ thuộc một dependency lỗi hoặc Pod bị đặt trên cùng node.

### Secret có được mã hóa không?

Secret thường chỉ base64 trong YAML. Mức bảo vệ phụ thuộc encryption at rest, RBAC, etcd và secret manager.

### Cách trả lời khi kinh nghiệm Kubernetes chưa sâu

> Em đã triển khai ứng dụng Docker và hiểu các thành phần Kubernetes gồm Pod, Deployment, Service, Ingress, ConfigMap, Secret, probe và rolling update. Em có thể đọc manifest, deploy và debug các lỗi cơ bản bằng `kubectl`. Kinh nghiệm vận hành cluster production của em chưa sâu nên em sẽ nói rõ phạm vi đã thực hiện thay vì khẳng định quá mức.

---

## 18. Bài tập thực hành

1. Viết Deployment ba replica cho FastAPI.
2. Tạo Service và Ingress cho domain `api.local`.
3. Thêm readiness, liveness và startup probe.
4. Giới hạn CPU/RAM và mô phỏng OOMKilled.
5. Deploy image lỗi rồi dùng `kubectl logs --previous` để debug.
6. Tạo ConfigMap, Secret và inject vào container.
7. Thực hiện rolling update và rollback.
8. Tạo Job chạy Alembic migration.

## Checklist

- [ ] Giải thích được desired state và reconciliation.
- [ ] Phân biệt Pod, Deployment, ReplicaSet và StatefulSet.
- [ ] Hiểu Service và Ingress.
- [ ] Phân biệt readiness, liveness và startup probe.
- [ ] Hiểu request, limit và OOMKilled.
- [ ] Debug được Pending, ImagePullBackOff và CrashLoopBackOff.
- [ ] Hiểu rolling update và migration tương thích ngược.
- [ ] Biết ConfigMap, Secret, RBAC và ServiceAccount.
