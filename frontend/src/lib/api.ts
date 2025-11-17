import axios from 'axios'

const API_URL = process.env.API_URL || 'http://localhost:8080'
const ML_API_URL = process.env.ML_API_URL || 'http://localhost:8000'

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const mlApi = axios.create({
  baseURL: ML_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Auth API
export const authApi = {
  login: (email: string, password: string) =>
    api.post('/api/v1/auth/login', { email, password }),
  register: (email: string, password: string, role: string) =>
    api.post('/api/v1/auth/register', { email, password, role }),
  logout: () => api.post('/api/v1/auth/logout'),
  getCurrentUser: () => api.get('/api/v1/auth/me'),
}

// Estimation API
export const estimationApi = {
  create: (data: any) => api.post('/api/v1/estimates', data),
  list: (params?: any) => api.get('/api/v1/estimates', { params }),
  get: (id: string) => api.get(`/api/v1/estimates/${id}`),
  update: (id: string, data: any) => api.put(`/api/v1/estimates/${id}`, data),
  uploadBlueprint: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return mlApi.post('/api/v1/cv/analyze-blueprint', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
}

// Project API
export const projectApi = {
  create: (data: any) => api.post('/api/v1/projects', data),
  list: (params?: any) => api.get('/api/v1/projects', { params }),
  get: (id: string) => api.get(`/api/v1/projects/${id}`),
  update: (id: string, data: any) => api.put(`/api/v1/projects/${id}`, data),
}

// Training API
export const trainingApi = {
  listCourses: () => api.get('/api/v1/courses'),
  getCourse: (id: string) => api.get(`/api/v1/courses/${id}`),
  enroll: (courseId: string) => api.post('/api/v1/enrollments', { courseId }),
  chatWithCoach: (messages: any[]) =>
    mlApi.post('/api/v1/coach/chat', { messages }),
}

// Marketplace API
export const marketplaceApi = {
  listJobs: (params?: any) => api.get('/api/v1/jobs', { params }),
  getJob: (id: string) => api.get(`/api/v1/jobs/${id}`),
  applyToJob: (id: string, data: any) =>
    api.post(`/api/v1/jobs/${id}/apply`, data),
}
