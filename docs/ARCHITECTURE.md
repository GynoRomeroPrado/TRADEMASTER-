# TRADEMASTER - System Architecture

## Table of Contents
1. [Overview](#overview)
2. [System Components](#system-components)
3. [Data Models](#data-models)
4. [API Design](#api-design)
5. [ML/AI Pipeline](#mlai-pipeline)
6. [Security](#security)
7. [Scalability](#scalability)

## Overview

TRADEMASTER follows a microservices architecture with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                     API Gateway (Kong)                       │
│                    JWT Auth, Rate Limiting                   │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
┌───────▼────────┐   ┌────────▼────────┐   ┌───────▼────────┐
│  Frontend App  │   │   Mobile App    │   │  Admin Panel   │
│   (Next.js)    │   │ (React Native)  │   │   (Next.js)    │
└────────────────┘   └─────────────────┘   └────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
┌───────▼────────┐   ┌────────▼────────┐   ┌───────▼────────┐
│  Estimation    │   │    Project      │   │    Training    │
│   Service      │   │    Service      │   │    Service     │
│     (Go)       │   │     (Go)        │   │     (Go)       │
└───────┬────────┘   └────────┬────────┘   └───────┬────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
┌───────▼────────┐   ┌────────▼────────┐   ┌───────▼────────┐
│  Marketplace   │   │   Auth Service  │   │  ML Services   │
│   Service      │   │     (Go)        │   │   (Python)     │
│     (Go)       │   │                 │   │  TensorFlow    │
└───────┬────────┘   └────────┬────────┘   └───────┬────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
┌───────▼────────┐   ┌────────▼────────┐   ┌───────▼────────┐
│   PostgreSQL   │   │    MongoDB      │   │ Elasticsearch  │
│  (Relational)  │   │   (Documents)   │   │  (Search)      │
└────────────────┘   └─────────────────┘   └────────────────┘
```

## System Components

### 1. Frontend (Next.js 14)
**Location**: `/frontend`

**Responsibilities**:
- Server-side rendering for SEO
- Client-side interactivity
- Real-time updates via WebSockets
- File uploads (blueprints, photos)

**Key Technologies**:
- Next.js 14 (App Router)
- TypeScript
- Tailwind CSS + shadcn/ui
- Zustand (state management)
- React Query (data fetching)
- Socket.io-client (real-time)

**Routes**:
```
/                       → Landing page
/dashboard              → Contractor dashboard
/projects               → Project list
/projects/[id]          → Project details
/estimates              → Estimates list
/estimates/new          → Create estimate
/training               → Training hub
/training/[course]      → Course viewer
/marketplace            → Jobs marketplace
/settings               → User settings
```

### 2. Mobile App (React Native)
**Location**: `/mobile`

**Responsibilities**:
- Field data collection
- AR training modules
- Photo capture with geolocation
- Offline-first functionality

**Key Technologies**:
- React Native + Expo
- ARCore / ARKit
- React Native Maps
- Expo Camera
- AsyncStorage (offline)
- React Navigation

**Screens**:
```
HomeScreen              → Dashboard
ProjectListScreen       → Active projects
ProjectDetailScreen     → Project details + map
CameraScreen            → Progress photos
ARTrainingScreen        → AR lessons
MarketplaceScreen       → Job search
ProfileScreen           → User profile
```

### 3. Backend Microservices (Go)

#### 3.1 Estimation Service
**Location**: `/backend/estimation-service`

**Responsibilities**:
- Parse PDF blueprints (via ML service)
- Calculate material quantities
- Predict costs with ML
- Generate estimates

**Endpoints**:
```
POST   /api/v1/estimates              → Create estimate
GET    /api/v1/estimates              → List estimates
GET    /api/v1/estimates/:id          → Get estimate
PUT    /api/v1/estimates/:id          → Update estimate
POST   /api/v1/estimates/:id/convert  → Convert to project
POST   /api/v1/blueprints/upload      → Upload blueprint
POST   /api/v1/blueprints/analyze     → Analyze blueprint
GET    /api/v1/materials/prices       → Get material prices
```

**Database Tables**:
- `estimates` (PostgreSQL)
- `estimate_items` (PostgreSQL)
- `blueprints` (MongoDB - with GridFS)
- `material_prices` (PostgreSQL)

#### 3.2 Project Service
**Location**: `/backend/project-service`

**Responsibilities**:
- Project CRUD operations
- Task scheduling
- Progress tracking
- Team assignments
- Route optimization

**Endpoints**:
```
POST   /api/v1/projects               → Create project
GET    /api/v1/projects               → List projects
GET    /api/v1/projects/:id           → Get project
PUT    /api/v1/projects/:id           → Update project
DELETE /api/v1/projects/:id           → Delete project
POST   /api/v1/projects/:id/tasks     → Create task
PUT    /api/v1/tasks/:id              → Update task
POST   /api/v1/tasks/:id/complete     → Complete task
GET    /api/v1/projects/:id/schedule  → Get schedule
POST   /api/v1/projects/:id/photos    → Upload photo
GET    /api/v1/projects/:id/timeline  → Get timeline
```

**Database Tables**:
- `projects` (PostgreSQL)
- `tasks` (PostgreSQL)
- `project_photos` (PostgreSQL - metadata, S3 - files)
- `schedules` (PostgreSQL)
- `assignments` (PostgreSQL)

#### 3.3 Training Service
**Location**: `/backend/training-service`

**Responsibilities**:
- Course management
- Progress tracking
- Certification issuance
- AI coaching integration

**Endpoints**:
```
GET    /api/v1/courses                → List courses
GET    /api/v1/courses/:id            → Get course
POST   /api/v1/enrollments            → Enroll in course
GET    /api/v1/enrollments            → My enrollments
PUT    /api/v1/progress/:id           → Update progress
POST   /api/v1/certifications         → Issue certification
GET    /api/v1/certifications         → My certifications
POST   /api/v1/coaching/chat          → AI coaching chat
GET    /api/v1/coaching/feedback      → Get AI feedback
```

**Database Tables**:
- `courses` (MongoDB)
- `lessons` (MongoDB)
- `enrollments` (PostgreSQL)
- `progress` (PostgreSQL)
- `certifications` (PostgreSQL)
- `coaching_sessions` (MongoDB)

#### 3.4 Marketplace Service
**Location**: `/backend/marketplace-service`

**Responsibilities**:
- Job posting
- Contractor matching
- Bidding system
- Payment processing
- Reputation management

**Endpoints**:
```
GET    /api/v1/jobs                   → List jobs
POST   /api/v1/jobs                   → Post job
GET    /api/v1/jobs/:id               → Get job details
POST   /api/v1/jobs/:id/apply         → Apply to job
GET    /api/v1/applications           → My applications
PUT    /api/v1/applications/:id       → Update application
POST   /api/v1/jobs/:id/award         → Award job
GET    /api/v1/contractors/search     → Search contractors
POST   /api/v1/reviews                → Leave review
GET    /api/v1/reputation/:id         → Get reputation
```

**Database Tables**:
- `jobs` (PostgreSQL + Elasticsearch)
- `applications` (PostgreSQL)
- `contracts` (PostgreSQL)
- `reviews` (PostgreSQL)
- `reputation_scores` (PostgreSQL)
- `payments` (PostgreSQL)

#### 3.5 Auth Service
**Location**: `/backend/auth-service`

**Responsibilities**:
- User authentication
- OAuth integration
- JWT management
- Role-based access control

**Endpoints**:
```
POST   /api/v1/auth/register          → Register user
POST   /api/v1/auth/login             → Login
POST   /api/v1/auth/refresh           → Refresh token
POST   /api/v1/auth/logout            → Logout
GET    /api/v1/auth/me                → Current user
POST   /api/v1/auth/oauth/google      → Google OAuth
POST   /api/v1/auth/password/reset    → Password reset
```

**Database Tables**:
- `users` (PostgreSQL)
- `user_profiles` (PostgreSQL)
- `user_roles` (PostgreSQL)
- `sessions` (Redis)
- `refresh_tokens` (PostgreSQL)

### 4. ML/AI Services (Python)

#### 4.1 Cost Prediction Service
**Location**: `/ml-services/cost-prediction`

**Models**:
- LSTM for time-series forecasting
- XGBoost for project duration
- Random Forest for risk assessment

**Endpoints**:
```
POST   /api/v1/predict/material-costs → Predict material costs
POST   /api/v1/predict/duration       → Predict project duration
POST   /api/v1/predict/risk           → Assess risk
POST   /api/v1/anomaly/detect         → Detect cost anomalies
```

#### 4.2 Computer Vision Service
**Location**: `/ml-services/cv-blueprint-analysis`

**Models**:
- YOLO v8 for object detection
- OCR for text extraction
- Custom CNN for room recognition

**Endpoints**:
```
POST   /api/v1/cv/analyze-blueprint   → Analyze blueprint
POST   /api/v1/cv/extract-text        → Extract text/dimensions
POST   /api/v1/cv/detect-components   → Detect electrical components
```

#### 4.3 Training Coach Service
**Location**: `/ml-services/training-coach`

**Models**:
- GPT-4 for conversational coaching
- Computer vision for movement analysis
- Speech-to-text for voice commands

**Endpoints**:
```
POST   /api/v1/coach/chat             → Chat with AI coach
POST   /api/v1/coach/analyze-movement → Analyze technique
POST   /api/v1/coach/voice-command    → Process voice command
POST   /api/v1/coach/feedback         → Get personalized feedback
```

#### 4.4 Inventory Optimization Service
**Location**: `/ml-services/inventory-optimization`

**Models**:
- Reinforcement learning for stock levels
- Prophet for demand forecasting

**Endpoints**:
```
POST   /api/v1/inventory/optimize     → Optimize stock levels
POST   /api/v1/inventory/forecast     → Forecast demand
POST   /api/v1/inventory/reorder      → Get reorder recommendations
```

## Data Models

### User
```json
{
  "id": "uuid",
  "email": "string",
  "password_hash": "string",
  "role": "contractor | worker | admin",
  "profile": {
    "first_name": "string",
    "last_name": "string",
    "phone": "string",
    "company": "string",
    "trade": "electrical | hvac | welding",
    "license_number": "string",
    "certifications": ["string"],
    "location": {
      "lat": "float",
      "lng": "float",
      "address": "string"
    }
  },
  "subscription": {
    "plan": "free | pro | premium | enterprise",
    "status": "active | canceled | past_due",
    "expires_at": "timestamp"
  },
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Project
```json
{
  "id": "uuid",
  "name": "string",
  "description": "string",
  "owner_id": "uuid",
  "client": {
    "name": "string",
    "email": "string",
    "phone": "string",
    "address": "string"
  },
  "location": {
    "lat": "float",
    "lng": "float",
    "address": "string"
  },
  "status": "draft | scheduled | in_progress | completed | canceled",
  "type": "electrical | hvac | welding",
  "budget": "float",
  "actual_cost": "float",
  "estimated_hours": "float",
  "actual_hours": "float",
  "start_date": "date",
  "end_date": "date",
  "team": ["uuid"],
  "tasks": ["uuid"],
  "materials": [{
    "item_id": "uuid",
    "name": "string",
    "quantity": "float",
    "unit": "string",
    "cost": "float"
  }],
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Estimate
```json
{
  "id": "uuid",
  "project_name": "string",
  "client": {
    "name": "string",
    "email": "string",
    "phone": "string"
  },
  "contractor_id": "uuid",
  "status": "draft | sent | accepted | rejected | expired",
  "type": "electrical | hvac | welding",
  "blueprint_id": "uuid",
  "labor": {
    "hours": "float",
    "rate": "float",
    "total": "float"
  },
  "materials": [{
    "name": "string",
    "quantity": "float",
    "unit": "string",
    "unit_price": "float",
    "total": "float"
  }],
  "subtotal": "float",
  "tax": "float",
  "total": "float",
  "confidence_score": "float",
  "valid_until": "date",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Course
```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "trade": "electrical | hvac | welding",
  "level": "beginner | intermediate | advanced",
  "duration_hours": "float",
  "format": "video | ar | vr | mixed",
  "certification": {
    "enabled": "boolean",
    "name": "string",
    "issuer": "string"
  },
  "modules": [{
    "id": "uuid",
    "title": "string",
    "lessons": [{
      "id": "uuid",
      "title": "string",
      "type": "video | ar | quiz | assessment",
      "content_url": "string",
      "duration_minutes": "float"
    }]
  }],
  "prerequisites": ["uuid"],
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Job
```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "client_id": "uuid",
  "trade": "electrical | hvac | welding",
  "location": {
    "lat": "float",
    "lng": "float",
    "address": "string"
  },
  "budget_range": {
    "min": "float",
    "max": "float"
  },
  "start_date": "date",
  "duration_days": "int",
  "required_certifications": ["string"],
  "status": "open | in_review | awarded | in_progress | completed | canceled",
  "applications_count": "int",
  "awarded_to": "uuid",
  "created_at": "timestamp",
  "expires_at": "timestamp"
}
```

## API Design

### Authentication
All API requests (except auth endpoints) require JWT token:

```
Authorization: Bearer <jwt_token>
```

### Error Handling
Standardized error responses:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "field": "email",
      "issue": "Invalid email format"
    },
    "request_id": "uuid"
  }
}
```

### Pagination
List endpoints support cursor-based pagination:

```
GET /api/v1/projects?limit=20&cursor=abc123

Response:
{
  "data": [...],
  "pagination": {
    "next_cursor": "xyz789",
    "has_more": true
  }
}
```

### Rate Limiting
- Free tier: 100 requests/15 minutes
- Pro tier: 1000 requests/15 minutes
- Enterprise: Custom limits

Headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1234567890
```

## ML/AI Pipeline

### 1. Blueprint Analysis Pipeline
```
PDF Upload → Image Preprocessing → YOLO Detection → OCR →
Entity Extraction → Quantity Calculation → Cost Prediction
```

### 2. Cost Prediction Pipeline
```
Historical Data → Feature Engineering → LSTM Training →
Model Serving → Real-time Prediction → Confidence Score
```

### 3. Training Coach Pipeline
```
User Query → GPT-4 → Context Retrieval → Response Generation →
Voice Synthesis → Delivery
```

### Model Deployment
- TensorFlow Serving for LSTM models
- TorchServe for YOLO models
- OpenAI API for GPT-4
- Model versioning with MLflow
- A/B testing framework

## Security

### 1. Authentication & Authorization
- JWT with short expiry (15 min)
- Refresh tokens (30 days)
- Role-based access control (RBAC)
- OAuth 2.0 for third-party login

### 2. Data Protection
- Encryption at rest (AES-256)
- Encryption in transit (TLS 1.3)
- PII encryption in database
- Secure file storage (S3 with encryption)

### 3. API Security
- Rate limiting
- Request validation (JSON Schema)
- SQL injection prevention (parameterized queries)
- XSS prevention (input sanitization)
- CSRF protection (SameSite cookies)

### 4. Compliance
- GDPR compliance
- SOC 2 Type II (planned)
- PCI DSS for payments (via Stripe)

## Scalability

### Horizontal Scaling
- Stateless microservices
- Load balancing with AWS ALB
- Auto-scaling based on CPU/memory

### Database Scaling
- Read replicas for PostgreSQL
- Sharding strategy for high-volume tables
- MongoDB replica sets
- Redis cluster

### Caching Strategy
- Redis for session data
- CDN for static assets (CloudFront)
- API response caching (15-60 seconds)
- Browser caching for images

### Performance Targets
- API latency: p95 < 200ms
- Page load: < 2 seconds
- Mobile app start: < 1 second
- ML inference: < 500ms
- Uptime: 99.9%

### Monitoring
- Prometheus + Grafana for metrics
- ELK Stack for logs
- Sentry for error tracking
- Datadog for APM
- Custom business metrics dashboard

---

**Last Updated**: 2025-01-17
