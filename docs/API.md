# TRADEMASTER API Documentation

Complete API reference for the TRADEMASTER platform.

## Table of Contents

- [Authentication](#authentication)
- [Projects API](#projects-api)
- [Estimates API](#estimates-api)
- [Training API](#training-api)
- [Marketplace API](#marketplace-api)
- [Blueprint Analysis API](#blueprint-analysis-api)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)

## Base URL

```
Development: http://localhost:3000/api
Production: https://api.trademaster.io
```

## Authentication

All API requests (except login/signup) require authentication using JWT tokens.

### Headers

```
Authorization: Bearer <your-jwt-token>
Content-Type: application/json
```

### POST /api/auth/signup

Register a new user account.

**Request:**

```json
{
  "email": "john@example.com",
  "password": "SecurePassword123!",
  "firstName": "John",
  "lastName": "Smith",
  "company": "Smith Electrical",
  "trade": "electrical"
}
```

**Response (201):**

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john@example.com",
    "role": "contractor",
    "profile": {
      "firstName": "John",
      "lastName": "Smith",
      "company": "Smith Electrical",
      "trade": "electrical"
    }
  },
  "message": "Signup successful"
}
```

**Errors:**
- 400: Invalid input data
- 409: Email already exists

---

### POST /api/auth/login

Authenticate and receive JWT token.

**Request:**

```json
{
  "email": "john@example.com",
  "password": "SecurePassword123!"
}
```

**Response (200):**

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john@example.com",
    "role": "contractor"
  },
  "message": "Login successful"
}
```

**Note:** JWT token is set in httpOnly cookie.

**Errors:**
- 401: Invalid credentials
- 400: Missing email or password

---

### POST /api/auth/logout

Logout and clear authentication token.

**Response (200):**

```json
{
  "message": "Logout successful"
}
```

---

## Projects API

### GET /api/projects

List all projects for authenticated user.

**Query Parameters:**
- `status` (optional): Filter by status (planning, active, completed)
- `trade` (optional): Filter by trade
- `limit` (optional): Number of results (default: 20)
- `offset` (optional): Pagination offset (default: 0)

**Response (200):**

```json
{
  "projects": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Office Building Rewire",
      "description": "Complete electrical rewiring of 5-story office building",
      "status": "active",
      "trade": "electrical",
      "location": {
        "address": "123 Main St",
        "city": "San Francisco",
        "state": "CA",
        "zip": "94105"
      },
      "startDate": "2024-01-15",
      "endDate": "2024-03-30",
      "budget": 125000,
      "progress": 45,
      "teamMembers": 5,
      "createdAt": "2024-01-10T10:00:00Z",
      "updatedAt": "2024-02-15T14:30:00Z"
    }
  ],
  "total": 12,
  "limit": 20,
  "offset": 0
}
```

---

### POST /api/projects

Create a new project.

**Request:**

```json
{
  "name": "Data Center Installation",
  "description": "New data center electrical infrastructure",
  "trade": "electrical",
  "location": {
    "address": "456 Tech Blvd",
    "city": "Austin",
    "state": "TX",
    "zip": "78701"
  },
  "startDate": "2024-03-01",
  "endDate": "2024-05-15",
  "budget": 250000
}
```

**Response (201):**

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "name": "Data Center Installation",
  "status": "planning",
  "createdAt": "2024-02-15T10:00:00Z"
}
```

**Errors:**
- 400: Invalid input data
- 401: Unauthorized

---

### GET /api/projects/:id

Get project details.

**Response (200):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Office Building Rewire",
  "description": "Complete electrical rewiring",
  "status": "active",
  "trade": "electrical",
  "location": {...},
  "startDate": "2024-01-15",
  "endDate": "2024-03-30",
  "budget": 125000,
  "progress": 45,
  "tasks": [
    {
      "id": "task-1",
      "title": "Install main panel",
      "status": "completed",
      "dueDate": "2024-01-20"
    }
  ],
  "team": [
    {
      "id": "user-1",
      "name": "Mike Johnson",
      "role": "lead_electrician"
    }
  ],
  "materials": [...],
  "photos": [...],
  "createdAt": "2024-01-10T10:00:00Z",
  "updatedAt": "2024-02-15T14:30:00Z"
}
```

**Errors:**
- 404: Project not found
- 401: Unauthorized
- 403: Access denied

---

### PUT /api/projects/:id

Update project.

**Request:**

```json
{
  "status": "completed",
  "progress": 100,
  "actualEndDate": "2024-03-25"
}
```

**Response (200):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Project updated successfully"
}
```

---

### DELETE /api/projects/:id

Delete project.

**Response (200):**

```json
{
  "message": "Project deleted successfully"
}
```

---

## Estimates API

### GET /api/estimates

List all estimates.

**Response (200):**

```json
{
  "estimates": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440000",
      "projectName": "Residential Rewire",
      "trade": "electrical",
      "status": "draft",
      "totalCost": 8500,
      "createdAt": "2024-02-10T10:00:00Z"
    }
  ]
}
```

---

### POST /api/estimates

Create new estimate.

**Request:**

```json
{
  "projectName": "Commercial HVAC Install",
  "trade": "hvac",
  "materials": [
    {
      "name": "AC Unit - 5 Ton",
      "quantity": 2,
      "unitPrice": 2500
    },
    {
      "name": "Ductwork - 50ft",
      "quantity": 10,
      "unitPrice": 45
    }
  ],
  "laborHours": 80,
  "laborRate": 95,
  "markup": 0.15
}
```

**Response (201):**

```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "projectName": "Commercial HVAC Install",
  "totalCost": 10092.50,
  "breakdown": {
    "materials": 5450,
    "labor": 7600,
    "subtotal": 13050,
    "markup": 1957.50,
    "total": 15007.50
  },
  "createdAt": "2024-02-15T10:00:00Z"
}
```

---

## Training API

### GET /api/training/courses

List available training courses.

**Query Parameters:**
- `trade` (optional): Filter by trade
- `level` (optional): beginner, intermediate, advanced

**Response (200):**

```json
{
  "courses": [
    {
      "id": "course-1",
      "title": "Data Center Electrical Systems",
      "description": "Master data center electrical installations",
      "trade": "electrical",
      "level": "advanced",
      "duration": 12,
      "price": 299,
      "certification": true,
      "rating": 4.8,
      "enrolled": 245
    }
  ]
}
```

---

### GET /api/training/courses/:id

Get course details.

**Response (200):**

```json
{
  "id": "course-1",
  "title": "Data Center Electrical Systems",
  "description": "...",
  "modules": [
    {
      "id": "module-1",
      "title": "Introduction to Data Centers",
      "order": 1,
      "lessons": [
        {
          "id": "lesson-1",
          "title": "What is a Data Center?",
          "type": "video",
          "duration": 15,
          "completed": false
        }
      ]
    }
  ],
  "instructor": {...},
  "prerequisites": [...],
  "whatYouLearn": [...]
}
```

---

### POST /api/training/enroll/:courseId

Enroll in a course.

**Response (201):**

```json
{
  "enrollmentId": "enroll-1",
  "courseId": "course-1",
  "status": "active",
  "progress": 0,
  "enrolledAt": "2024-02-15T10:00:00Z"
}
```

---

## Marketplace API

### GET /api/marketplace/jobs

List available jobs.

**Query Parameters:**
- `trade` (optional): Filter by trade
- `budget` (optional): under20k, 20k-50k, over50k
- `urgency` (optional): critical, urgent, normal

**Response (200):**

```json
{
  "jobs": [
    {
      "id": "job-1",
      "title": "Data Center Electrical Installation",
      "description": "Need experienced electrician...",
      "trade": "electrical",
      "budget": 45000,
      "budgetType": "fixed",
      "location": {
        "city": "San Francisco",
        "state": "CA"
      },
      "urgency": "urgent",
      "duration": "6 weeks",
      "postedAt": "2024-02-13T10:00:00Z",
      "proposals": 12,
      "matchScore": 95
    }
  ]
}
```

---

### GET /api/marketplace/jobs/:id

Get job details.

**Response (200):**

```json
{
  "id": "job-1",
  "title": "Data Center Electrical Installation",
  "description": "...",
  "longDescription": "...",
  "requirements": [...],
  "scopeOfWork": [...],
  "client": {
    "name": "TechCorp Inc.",
    "rating": 4.8,
    "jobsPosted": 23,
    "verified": true
  },
  "budget": 45000,
  "startDate": "2024-03-01"
}
```

---

### POST /api/marketplace/proposals

Submit job proposal.

**Request:**

```json
{
  "jobId": "job-1",
  "proposedBudget": 45000,
  "timeline": "6 weeks",
  "coverLetter": "With over 8 years of experience..."
}
```

**Response (201):**

```json
{
  "id": "proposal-1",
  "jobId": "job-1",
  "status": "pending",
  "submittedAt": "2024-02-15T10:00:00Z"
}
```

---

### GET /api/marketplace/proposals

Get user's proposals.

**Response (200):**

```json
{
  "proposals": [
    {
      "id": "proposal-1",
      "jobId": "job-1",
      "jobTitle": "Data Center Electrical Installation",
      "status": "under_review",
      "proposedBudget": 45000,
      "timeline": "6 weeks",
      "submittedAt": "2024-02-15T10:00:00Z",
      "viewedByClient": true
    }
  ]
}
```

---

## Blueprint Analysis API

### POST /api/blueprints/analyze

Analyze construction blueprint using AI/ML.

**Request:**

Content-Type: multipart/form-data

```
file: <blueprint-image.pdf or .jpg>
```

**Response (200):**

```json
{
  "detected_components": [
    {
      "type": "electrical_panel",
      "quantity": 2,
      "location": "main_floor"
    },
    {
      "type": "outlet",
      "quantity": 45,
      "location": "throughout"
    }
  ],
  "area_sqft": 5000,
  "suggested_materials": [
    {
      "name": "12/2 NM-B Cable",
      "quantity": 10,
      "unit": "roll",
      "estimatedCost": 899.90
    },
    {
      "name": "15A Circuit Breaker",
      "quantity": 20,
      "unit": "each",
      "estimatedCost": 250.00
    }
  ],
  "confidence_score": 0.92
}
```

**Errors:**
- 400: Invalid file format
- 413: File too large (max 10MB)

---

## Error Handling

All errors follow this format:

```json
{
  "error": "Error message description",
  "code": "ERROR_CODE",
  "details": {
    "field": "Additional context"
  }
}
```

### Common Error Codes

| Code | Status | Description |
|------|--------|-------------|
| UNAUTHORIZED | 401 | Missing or invalid authentication |
| FORBIDDEN | 403 | Insufficient permissions |
| NOT_FOUND | 404 | Resource not found |
| VALIDATION_ERROR | 400 | Invalid input data |
| RATE_LIMIT_EXCEEDED | 429 | Too many requests |
| SERVER_ERROR | 500 | Internal server error |

## Rate Limiting

API rate limits per authenticated user:

- **Standard tier**: 100 requests/minute
- **Premium tier**: 500 requests/minute
- **Enterprise tier**: Unlimited

Rate limit headers:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1708012800
```

When rate limit exceeded:

```json
{
  "error": "Rate limit exceeded",
  "code": "RATE_LIMIT_EXCEEDED",
  "retryAfter": 60
}
```

## Pagination

List endpoints support pagination:

**Query Parameters:**
- `limit`: Results per page (default: 20, max: 100)
- `offset`: Number of results to skip (default: 0)

**Response includes:**

```json
{
  "data": [...],
  "total": 150,
  "limit": 20,
  "offset": 40
}
```

## Webhooks

Subscribe to events:

```json
{
  "event": "project.completed",
  "url": "https://your-app.com/webhooks",
  "secret": "webhook-signing-secret"
}
```

**Event Types:**
- `project.created`
- `project.updated`
- `project.completed`
- `proposal.submitted`
- `proposal.accepted`
- `course.completed`
- `certification.issued`

## SDKs

Official SDKs available:

- JavaScript/TypeScript: `npm install @trademaster/sdk`
- Python: `pip install trademaster-sdk`
- Go: `go get github.com/trademaster/sdk-go`

**Example (JavaScript):**

```javascript
import { TradeMaster } from '@trademaster/sdk'

const client = new TradeMaster({
  apiKey: 'your-api-key',
  environment: 'production'
})

const projects = await client.projects.list()
```

## Support

- API Status: https://status.trademaster.io
- Documentation: https://docs.trademaster.io
- Support: api-support@trademaster.io
