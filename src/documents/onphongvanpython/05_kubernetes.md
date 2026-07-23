# 5. Kubernetes

## 1. Tổng quan

Kubernetes quản lý container ở quy mô cluster.

Luồng thường gặp:

```text
Client → Ingress → Service → Pods
                         ↑
                    Deployment
```

Deployment không nằm trực tiếp trên đường request; nó quản lý ReplicaSet và Pods.

## 2. Pod

Pod là đơn vị deploy nhỏ nhất trong Kubernetes.

Một Pod có thể chứa một hoặc nhiều container dùng chung:

- Network namespace.
- IP.
- Volume.
- Lifecycle.

Pod thường không được tạo thủ công trong production mà được Deployment quản lý.

## 3. Deployment và ReplicaSet

Deployment khai báo trạng thái mong muốn:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api
  template:
    metadata:
      labels:
        app: api
    spec:
      containers:
        - name: api
          image: my-api:1.0.0
```

ReplicaSet đảm bảo số lượng Pod mong muốn.

## 4. Service

Pod có IP thay đổi. Service cung cấp endpoint ổn định và load balancing đến các Pod theo label selector.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: api-service
spec:
  selector:
    app: api
  ports:
    - port: 80
      targetPort: 8000
```

## 5. Các loại Service

### ClusterIP

Chỉ truy cập trong cluster.

### NodePort

Mở một port trên mỗi node.

### LoadBalancer

Yêu cầu cloud hoặc implementation hỗ trợ load balancer.

## 6. Ingress

Ingress định tuyến HTTP/HTTPS theo host hoặc path.

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: api-ingress
spec:
  rules:
    - host: api.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: api-service
                port:
                  number: 80
```

Ingress cần Ingress Controller như NGINX Ingress Controller.

## 7. ConfigMap và Secret

### ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: api-config
data:
  APP_ENV: production
```

### Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: api-secret
type: Opaque
stringData:
  DB_PASSWORD: super-secret
```

Secret không mặc định được mã hóa an toàn chỉ vì dùng base64. Production nên bật encryption at rest và quản lý quyền RBAC.

## 8. Namespace

Namespace giúp phân tách tài nguyên:

```text
dev
staging
production
```

Không phải ranh giới bảo mật tuyệt đối nếu không kết hợp RBAC và NetworkPolicy.

## 9. PV và PVC

PersistentVolume đại diện storage.

PersistentVolumeClaim là yêu cầu storage của workload.

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

## 10. Liveness và readiness

### Liveness probe

Kiểm tra container còn sống hay bị treo. Nếu fail nhiều lần, kubelet restart container.

### Readiness probe

Kiểm tra Pod đã sẵn sàng nhận traffic chưa. Nếu fail, Pod bị loại khỏi Service endpoints nhưng không nhất thiết restart.

```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 8000
  initialDelaySeconds: 10

readinessProbe:
  httpGet:
    path: /health/ready
    port: 8000
  initialDelaySeconds: 5
```

## 11. Rolling update

```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxUnavailable: 0
    maxSurge: 1
```

Giúp cập nhật dần Pod cũ thành Pod mới.

Rollback:

```bash
kubectl rollout undo deployment/api
```

## 12. Resource requests và limits

```yaml
resources:
  requests:
    cpu: "200m"
    memory: "256Mi"
  limits:
    cpu: "1"
    memory: "512Mi"
```

- Requests dùng cho scheduling.
- Limits giới hạn tài nguyên.
- Vượt memory limit có thể bị OOMKilled.
- CPU limit thường dẫn đến throttling.

## 13. CrashLoopBackOff

CrashLoopBackOff nghĩa là container liên tục khởi động rồi crash.

Kiểm tra:

```bash
kubectl get pods
```

```bash
kubectl describe pod <pod-name>
```

```bash
kubectl logs <pod-name> --previous
```

Nguyên nhân thường gặp:

- Sai command.
- Thiếu env.
- Không kết nối được database.
- Migration lỗi.
- Permission.
- Liveness probe quá gắt.
- Ứng dụng thoát ngay.

## 14. Scale ngang

```bash
kubectl scale deployment api --replicas=5
```

HPA có thể scale theo CPU hoặc metric khác.

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api
  minReplicas: 2
  maxReplicas: 10
```

Ứng dụng cần stateless hoặc đưa session ra Redis/database để scale ngang dễ hơn.

## Câu hỏi phỏng vấn

### Pod và container khác nhau?

Container là tiến trình được cô lập. Pod là đơn vị Kubernetes quản lý, có thể chứa nhiều container chia sẻ network và volume.

### Deployment làm gì?

Quản lý desired state, số replica, rolling update, rollback và ReplicaSet.

### Service tìm Pod thế nào?

Thông qua label selector.

### Pod chết thì sao?

Nếu do Deployment quản lý, ReplicaSet tạo Pod mới để duy trì replica mong muốn.

### Làm sao deploy giảm downtime?

- Nhiều replica.
- Readiness probe.
- Rolling update.
- `maxUnavailable: 0`.
- Graceful shutdown.
- Backward-compatible migration.
- PodDisruptionBudget khi phù hợp.

## Cách trả lời khi kinh nghiệm chưa sâu

> Em đã sử dụng Docker thường xuyên và hiểu các thành phần Kubernetes như Deployment, Service, Ingress, ConfigMap, Secret, probe và rolling update. Kinh nghiệm vận hành production Kubernetes của em chưa sâu, nhưng em có thể đọc manifest, triển khai và debug các lỗi cơ bản bằng kubectl.
