# TRADEMASTER Deployment Guide

Complete guide for deploying TRADEMASTER to development, staging, and production environments.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Setup](#environment-setup)
- [Development Deployment](#development-deployment)
- [Production Deployment](#production-deployment)
- [Database Migration](#database-migration)
- [Monitoring & Logging](#monitoring--logging)
- [Scaling](#scaling)
- [Disaster Recovery](#disaster-recovery)

## Prerequisites

### Required Software

- Docker 20.10+
- Docker Compose 2.0+
- kubectl 1.24+ (for Kubernetes deployments)
- AWS CLI (for AWS deployments)
- PostgreSQL 14+
- MongoDB 5.0+
- Redis 7.0+

### Cloud Resources

**Production Environment:**
- AWS RDS PostgreSQL instance (db.r5.xlarge minimum)
- MongoDB Atlas M30 cluster
- AWS ElastiCache Redis cluster
- AWS S3 bucket for file storage
- AWS EKS cluster (3 nodes minimum)
- CloudFront CDN
- Route53 DNS

## Environment Setup

### 1. Clone Repository

```bash
git clone https://github.com/your-org/trademaster.git
cd trademaster
```

### 2. Environment Variables

Create environment files for each service:

**Backend Services (.env):**

```bash
# Database
DATABASE_URL=postgresql://user:password@localhost:5432/trademaster
MONGODB_URI=mongodb://localhost:27017/trademaster
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-super-secret-key-change-in-production
JWT_EXPIRY=7d

# Services URLs
AUTH_SERVICE_URL=http://localhost:8001
ESTIMATION_SERVICE_URL=http://localhost:8002
PROJECT_SERVICE_URL=http://localhost:8003
TRAINING_SERVICE_URL=http://localhost:8004
MARKETPLACE_SERVICE_URL=http://localhost:8005

# ML Services
ML_CV_SERVICE_URL=http://localhost:9001
ML_COST_SERVICE_URL=http://localhost:9002

# Storage
S3_BUCKET=trademaster-uploads
AWS_REGION=us-west-2

# Email (SendGrid)
SENDGRID_API_KEY=your-sendgrid-key
FROM_EMAIL=noreply@trademaster.io

# Stripe
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...
```

**Frontend (.env.local):**

```bash
NEXT_PUBLIC_API_URL=http://localhost:3000/api
AUTH_SERVICE_URL=http://localhost:8001
ESTIMATION_SERVICE_URL=http://localhost:8002
PROJECT_SERVICE_URL=http://localhost:8003
TRAINING_SERVICE_URL=http://localhost:8004
MARKETPLACE_SERVICE_URL=http://localhost:8005
ML_CV_SERVICE_URL=http://localhost:9001
```

**Mobile (.env):**

```bash
API_URL=http://localhost:3000/api
GOOGLE_MAPS_API_KEY=your-google-maps-key
```

## Development Deployment

### Using Docker Compose

1. **Start all services:**

```bash
make dev-up
```

This starts:
- PostgreSQL (port 5432)
- MongoDB (port 27017)
- Redis (port 6379)
- Elasticsearch (port 9200)
- RabbitMQ (port 5672, 15672)
- MinIO (port 9000)
- All backend microservices
- Frontend (port 3000)

2. **Run database migrations:**

```bash
make migrate-up
```

3. **Seed database with demo data:**

```bash
make seed
```

4. **Access services:**

- Frontend: http://localhost:3000
- Auth Service: http://localhost:8001
- Estimation Service: http://localhost:8002
- Project Service: http://localhost:8003
- Training Service: http://localhost:8004
- Marketplace Service: http://localhost:8005
- RabbitMQ Management: http://localhost:15672
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3001

### Manual Setup (without Docker)

**Backend Services:**

```bash
# Terminal 1 - Auth Service
cd backend/auth-service
go run cmd/main.go

# Terminal 2 - Estimation Service
cd backend/estimation-service
go run cmd/main.go

# Terminal 3 - Project Service
cd backend/project-service
go run cmd/main.go

# ... repeat for other services
```

**ML Services:**

```bash
# Terminal - CV Service
cd ml-services/computer-vision
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python app.py
```

**Frontend:**

```bash
cd frontend
npm install
npm run dev
```

## Production Deployment

### AWS Infrastructure Setup

**1. Create Infrastructure with Terraform:**

```bash
cd infrastructure/terraform
terraform init
terraform plan
terraform apply
```

This creates:
- VPC with public/private subnets
- EKS cluster
- RDS PostgreSQL instance
- ElastiCache Redis cluster
- S3 buckets
- Load balancers
- Security groups

**2. Configure kubectl:**

```bash
aws eks update-kubeconfig --name trademaster-prod --region us-west-2
```

### Kubernetes Deployment

**1. Create namespace:**

```bash
kubectl create namespace trademaster-prod
```

**2. Apply secrets:**

```bash
kubectl create secret generic app-secrets \
  --from-env-file=.env.production \
  --namespace=trademaster-prod
```

**3. Deploy services:**

```bash
# Backend services
kubectl apply -f k8s/backend/ --namespace=trademaster-prod

# ML services
kubectl apply -f k8s/ml-services/ --namespace=trademaster-prod

# Frontend
kubectl apply -f k8s/frontend/ --namespace=trademaster-prod

# Ingress
kubectl apply -f k8s/ingress/ --namespace=trademaster-prod
```

**4. Verify deployment:**

```bash
kubectl get pods --namespace=trademaster-prod
kubectl get services --namespace=trademaster-prod
```

### Using Deployment Scripts

**Development:**

```bash
./scripts/deploy-dev.sh
```

**Production:**

```bash
./scripts/deploy-prod.sh
```

The production script includes:
- Pre-deployment checks
- Database backup
- Rolling deployment
- Health checks
- Rollback on failure
- Slack notifications

## Database Migration

### PostgreSQL Migrations

**Create migration:**

```bash
make migrate-create name=add_new_table
```

**Apply migrations:**

```bash
# Development
make migrate-up

# Production
DATABASE_URL=$PROD_DATABASE_URL make migrate-up
```

**Rollback:**

```bash
make migrate-down
```

### MongoDB Schema Updates

MongoDB is schema-less, but application-level schema validation is defined in:

```
backend/training-service/internal/models/
```

## Monitoring & Logging

### Prometheus Metrics

Access metrics at: `https://metrics.trademaster.io`

**Key metrics:**
- Request rate per service
- Response time (p50, p95, p99)
- Error rate
- Database connection pool usage
- Cache hit rate

### Grafana Dashboards

Access dashboards at: `https://grafana.trademaster.io`

**Pre-configured dashboards:**
- Service Health Overview
- API Performance
- Database Performance
- User Activity
- Business Metrics

### Application Logs

**View logs:**

```bash
# All pods
kubectl logs -f -l app=auth-service --namespace=trademaster-prod

# Specific pod
kubectl logs -f pod-name --namespace=trademaster-prod

# Stream to CloudWatch
aws logs tail /aws/eks/trademaster-prod --follow
```

### Error Tracking

Sentry integration for error tracking:

```bash
SENTRY_DSN=your-sentry-dsn
```

## Scaling

### Horizontal Pod Autoscaling

**Configure HPA:**

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: auth-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: auth-service
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### Database Scaling

**Read Replicas:**

```bash
# Create read replica
aws rds create-db-instance-read-replica \
  --db-instance-identifier trademaster-replica-1 \
  --source-db-instance-identifier trademaster-primary
```

**Connection Pooling:**

Use PgBouncer for PostgreSQL connection pooling:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: pgbouncer-config
data:
  pgbouncer.ini: |
    [databases]
    trademaster = host=rds-endpoint port=5432 dbname=trademaster

    [pgbouncer]
    pool_mode = transaction
    max_client_conn = 1000
    default_pool_size = 25
```

## Disaster Recovery

### Backup Strategy

**Database Backups:**

- Automated daily snapshots (7-day retention)
- Point-in-time recovery enabled
- Cross-region backup replication

**Backup script:**

```bash
#!/bin/bash
# backup-databases.sh

# PostgreSQL
pg_dump -h $DB_HOST -U $DB_USER -d trademaster > backup_$(date +%Y%m%d).sql
aws s3 cp backup_$(date +%Y%m%d).sql s3://trademaster-backups/postgres/

# MongoDB
mongodump --uri=$MONGODB_URI --out=mongo_backup_$(date +%Y%m%d)
tar -czf mongo_backup_$(date +%Y%m%d).tar.gz mongo_backup_$(date +%Y%m%d)
aws s3 cp mongo_backup_$(date +%Y%m%d).tar.gz s3://trademaster-backups/mongodb/
```

**Automated backups (cron):**

```bash
0 2 * * * /scripts/backup-databases.sh
```

### Restore Procedure

**PostgreSQL:**

```bash
# Download backup
aws s3 cp s3://trademaster-backups/postgres/backup_20240215.sql .

# Restore
psql -h $DB_HOST -U $DB_USER -d trademaster < backup_20240215.sql
```

**MongoDB:**

```bash
# Download and extract
aws s3 cp s3://trademaster-backups/mongodb/mongo_backup_20240215.tar.gz .
tar -xzf mongo_backup_20240215.tar.gz

# Restore
mongorestore --uri=$MONGODB_URI mongo_backup_20240215
```

### Rollback Deployment

```bash
# Kubernetes rollback
kubectl rollout undo deployment/auth-service --namespace=trademaster-prod

# Check rollout status
kubectl rollout status deployment/auth-service --namespace=trademaster-prod
```

## SSL/TLS Configuration

**Using cert-manager:**

```bash
# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.12.0/cert-manager.yaml

# Create issuer
kubectl apply -f k8s/cert-issuer.yaml
```

## Performance Optimization

### CDN Configuration

CloudFront distribution for static assets:

```bash
# Invalidate cache
aws cloudfront create-invalidation \
  --distribution-id DISTRIBUTION_ID \
  --paths "/*"
```

### Database Indexing

**PostgreSQL indexes:**

```sql
-- Add indexes for common queries
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_jobs_trade ON jobs(trade);
CREATE INDEX idx_jobs_location ON jobs(location);
```

### Caching Strategy

- Redis for session storage
- Application-level caching for frequently accessed data
- CDN caching for static assets
- Browser caching headers

## Security Checklist

- [ ] Environment variables secured (AWS Secrets Manager)
- [ ] SSL/TLS certificates installed
- [ ] Database encryption at rest enabled
- [ ] VPC security groups configured
- [ ] IAM roles with least privilege
- [ ] Regular security updates applied
- [ ] Rate limiting configured
- [ ] CORS properly configured
- [ ] SQL injection prevention
- [ ] XSS protection enabled

## Support

For deployment issues:
- Slack: #trademaster-ops
- Email: ops@trademaster.io
- On-call: PagerDuty rotation
