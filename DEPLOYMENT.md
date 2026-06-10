# Deployment Guide

This guide covers deploying the Todo App using Docker and Kubernetes.

## Table of Contents

1. [Docker Compose (Local Development)](#docker-compose-local-development)
2. [Docker (Production Build)](#docker-production-build)
3. [Kubernetes Deployment](#kubernetes-deployment)

---

## Docker Compose (Local Development)

### Quick Start

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Accessing the App

- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html (enabled by default in dev)
- **Database**: localhost:5432 (postgres/securepassword)

### Swagger Availability

Swagger UI is controlled by `APP_ENV` environment variable:
- **development** (default): Swagger UI enabled
- **production**: Swagger UI disabled for security

### Database Migrations

The migrations are automatically applied when PostgreSQL container starts:

```bash
# Manual migration run
docker-compose exec postgres psql -U postgres -d todo_db -f /docker-entrypoint-initdb.d/001_initial_schema.sql
```

### Environment Variables

Edit `.env` or modify `docker-compose.yml` environment section:

```yaml
environment:
  APP_ENV: development        # development or production
  SERVER_PORT: 8080
  DB_TYPE: postgres
  DB_HOST: postgres
  DB_PORT: 5432
  DB_NAME: todo_db
  DB_USER: postgres
  DB_PASSWORD: securepassword
```

---

## Docker (Production Build)

### Build Image

```bash
# Build the container image
docker build -t todo-api:latest .
```

### Environment Configuration

Control features via `APP_ENV` variable:

```bash
# Development (Swagger enabled)
docker run -e APP_ENV=development -p 8080:8080 todo-api:latest

# Production (Swagger disabled)
docker run -e APP_ENV=production -p 8080:8080 todo-api:latest
```

### Docker Image Info

- **Base**: Alpine 3.19 (minimal, ~15MB runtime)
- **Built with**: Go 1.24.1-alpine3.19
- **Swagger UI**: Controlled via `APP_ENV` variable
  - `APP_ENV=development`: Swagger enabled at `/swagger/index.html`
  - `APP_ENV=production`: Swagger disabled (all access to `/swagger/` returns 404)

### Security Verification

```bash
# Scan for vulnerabilities (requires trivy)
trivy image todo-api:latest

# Auto-scan script
bash scan-vulnerabilities.sh
```

### Run Container

```bash
# Production (Distroless)
docker run -p 8080:8080 \
  -e DB_TYPE=postgres \
  -e DB_HOST=postgres.example.com \
  -e DB_PORT=5432 \
  -e DB_NAME=todo_db \
  -e DB_USER=postgres \
  -e DB_PASSWORD=securepassword \
  -e DB_SSLMODE=require \
  todo-api:latest

# Pre-Production with Debugging
docker run -p 8080:8080 \
  -e DB_TYPE=postgres \
  -e DB_HOST=localhost \
  todo-api:debian

# Development with Shell Access
docker run -it -p 8080:8080 \
  -e DB_TYPE=memory \
  todo-api:ubuntu /bin/bash
```

### Docker Features

✅ **Multi-stage Build**: Reduces image size
- Build stage: ~800MB (with Go compiler)
- Runtime: 30-150MB (depends on base image)

✅ **Security**:
- Non-root user execution
- **Zero known vulnerabilities** (distroless)
- Certificate verification
- Minimal attack surface

✅ **Health Check**: Automatic container health monitoring

### Push to Registry

```bash
# Tag image for registry
docker tag todo-api:latest myregistry.azurecr.io/todo-api:latest

# Push to registry
docker push myregistry.azurecr.io/todo-api:latest
```

**See [DOCKERFILE_SECURITY.md](DOCKERFILE_SECURITY.md) for detailed comparison and security analysis.**

---

## Kubernetes Deployment

### Prerequisites

- Kubernetes cluster (v1.20+)
- kubectl configured
- Docker image pushed to registry

### Deployment

```bash
# Apply all resources
kubectl apply -f deployments/k8s-deployment.yaml

# Verify deployment
kubectl get all -n todo-app

# View logs
kubectl logs -n todo-app deployment/todo-api

# Get service endpoint
kubectl get svc -n todo-app todo-api-service
```

### K8s Resources

The deployment includes:

1. **Namespace**: `todo-app` - Isolated namespace
2. **ConfigMap**: Environment configuration
3. **Secret**: Sensitive data (passwords)
4. **PostgreSQL Deployment**: Database with PVC
5. **PostgreSQL Service**: ClusterIP for internal access
6. **Todo API Deployment**: 3 replicas with RollingUpdate strategy
7. **Todo API Service**: LoadBalancer for external access
8. **HorizontalPodAutoscaler**: Auto-scaling (3-10 pods)
9. **PodDisruptionBudget**: HA configuration

### Configuration

Edit before applying:

```yaml
# Change image registry
image: myregistry.azurecr.io/todo-api:v1.0.0

# Adjust replicas
replicas: 5

# Update resource limits
resources:
  limits:
    memory: "512Mi"
    cpu: "500m"
```

### Scaling

```bash
# Manual scaling
kubectl scale deployment todo-api -n todo-app --replicas=5

# Check autoscaler status
kubectl get hpa -n todo-app

# View HPA metrics
kubectl top pods -n todo-app
```

### Accessing the App

```bash
# Get external IP
kubectl get svc -n todo-app todo-api-service

# Access via LoadBalancer
curl http://<EXTERNAL-IP>:80/swagger/index.html

# Port forward (for local testing)
kubectl port-forward -n todo-app svc/todo-api-service 8080:80
```

### Debugging

```bash
# View pod status
kubectl describe pod <pod-name> -n todo-app

# View pod logs
kubectl logs <pod-name> -n todo-app

# Execute command in pod
kubectl exec -it <pod-name> -n todo-app -- /bin/sh

# Check events
kubectl get events -n todo-app
```

### Cleanup

```bash
# Delete deployment and services
kubectl delete namespace todo-app

# Or delete individual resources
kubectl delete -f deployments/k8s-deployment.yaml
```

---

## Environment Variables Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | 8080 | HTTP server port |
| `DB_TYPE` | memory | Database type: `memory` or `postgres` |
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_NAME` | todo_db | Database name |
| `DB_USER` | postgres | Database username |
| `DB_PASSWORD` | - | Database password |
| `DB_SSLMODE` | disable | SSL mode: `disable`, `require`, `verify-ca`, `verify-full` |

---

## Production Checklist

- [ ] Update `DB_PASSWORD` in Secret
- [ ] Change image tag to specific version
- [ ] Update resource limits based on load testing
- [ ] Configure persistent volume for PostgreSQL
- [ ] Setup monitoring/logging (Prometheus, ELK)
- [ ] Enable network policies
- [ ] Setup ingress controller
- [ ] Configure backup strategy for database
- [ ] Review security policies
- [ ] Test failover scenarios

---

## Troubleshooting

### Database Connection Failed

```bash
# Check PostgreSQL is running
docker-compose ps

# Check logs
docker-compose logs postgres

# Test connection
psql -h localhost -U postgres -d todo_db
```

### Container Won't Start

```bash
# Check image build
docker build --no-cache -t todo-api:latest .

# Run with verbose output
docker run -it todo-api:latest ./server

# Check logs
docker logs <container-id>
```

### K8s Pod Pending

```bash
# Check node resources
kubectl top nodes

# Describe pod for events
kubectl describe pod <pod-name> -n todo-app

# Check image pull
kubectl get events -n todo-app
```

