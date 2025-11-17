// User and Authentication Types
export interface User {
  id: string
  email: string
  role: 'contractor' | 'client' | 'admin'
  profile: UserProfile
  createdAt: string
  updatedAt: string
}

export interface UserProfile {
  firstName: string
  lastName: string
  company?: string
  trade?: 'electrical' | 'hvac' | 'welding' | 'plumbing' | 'carpentry'
  phone?: string
  avatar?: string
  bio?: string
  location?: Location
  certifications?: Certification[]
  yearsExperience?: number
}

export interface Certification {
  id: string
  name: string
  issuer: string
  issuedDate: string
  expiryDate?: string
  verified: boolean
}

// Location Types
export interface Location {
  address?: string
  city: string
  state: string
  zip?: string
  country?: string
  coordinates?: {
    lat: number
    lng: number
  }
}

// Project Types
export interface Project {
  id: string
  userId: string
  name: string
  description: string
  status: ProjectStatus
  trade: string
  location: Location
  startDate: string
  endDate: string
  actualEndDate?: string
  budget: number
  actualCost?: number
  progress: number
  tasks: Task[]
  team: TeamMember[]
  materials: Material[]
  photos: Photo[]
  createdAt: string
  updatedAt: string
}

export type ProjectStatus = 'planning' | 'active' | 'on_hold' | 'completed' | 'cancelled'

export interface Task {
  id: string
  projectId: string
  title: string
  description?: string
  status: TaskStatus
  priority: 'low' | 'medium' | 'high' | 'critical'
  assignedTo?: string
  dueDate: string
  completedAt?: string
  order: number
}

export type TaskStatus = 'pending' | 'in_progress' | 'review' | 'completed'

export interface TeamMember {
  id: string
  userId: string
  name: string
  role: string
  avatar?: string
  email?: string
  phone?: string
}

export interface Material {
  id: string
  name: string
  category: string
  quantity: number
  unit: string
  unitPrice: number
  totalPrice: number
  supplier?: string
  ordered: boolean
  received: boolean
}

export interface Photo {
  id: string
  url: string
  thumbnail?: string
  description?: string
  uploadedBy: string
  uploadedAt: string
  location?: {
    lat: number
    lng: number
  }
  tags?: string[]
}

// Estimate Types
export interface Estimate {
  id: string
  userId: string
  projectName: string
  trade: string
  status: 'draft' | 'sent' | 'approved' | 'rejected'
  materials: EstimateMaterial[]
  laborHours: number
  laborRate: number
  markup: number
  totalCost: number
  breakdown: CostBreakdown
  notes?: string
  validUntil?: string
  createdAt: string
  updatedAt: string
}

export interface EstimateMaterial {
  name: string
  category: string
  quantity: number
  unit: string
  unitPrice: number
  totalPrice: number
}

export interface CostBreakdown {
  materials: number
  labor: number
  subtotal: number
  markup: number
  tax?: number
  total: number
}

export interface MaterialPrice {
  name: string
  category: string
  unit: string
  price: number
  trade: string
  supplier?: string
  lastUpdated: string
}

// Training Types
export interface Course {
  id: string
  title: string
  description: string
  longDescription?: string
  trade: string
  level: 'beginner' | 'intermediate' | 'advanced'
  duration: number // hours
  price: number
  certification: boolean
  rating: number
  enrolled: number
  thumbnail?: string
  instructor: Instructor
  modules: Module[]
  prerequisites?: string[]
  whatYouLearn?: string[]
  createdAt: string
  updatedAt: string
}

export interface Instructor {
  id: string
  name: string
  title: string
  bio: string
  avatar?: string
  rating: number
  students: number
  courses: number
}

export interface Module {
  id: string
  courseId: string
  title: string
  description?: string
  order: number
  duration: number
  lessons: Lesson[]
}

export interface Lesson {
  id: string
  moduleId: string
  title: string
  description?: string
  type: 'video' | 'ar' | 'quiz' | 'assessment' | 'reading'
  duration: number
  order: number
  completed: boolean
  locked: boolean
  contentUrl?: string
}

export interface Enrollment {
  id: string
  userId: string
  courseId: string
  status: 'active' | 'completed' | 'dropped'
  progress: number
  startedAt: string
  completedAt?: string
  certificateId?: string
}

// Marketplace Types
export interface Job {
  id: string
  clientId: string
  title: string
  description: string
  longDescription?: string
  trade: string
  budget: number
  budgetType: 'fixed' | 'hourly'
  location: Location
  urgency: 'normal' | 'urgent' | 'critical'
  duration: string
  startDate: string
  status: 'open' | 'in_progress' | 'completed' | 'cancelled'
  requirements: string[]
  scopeOfWork?: string[]
  client: Client
  proposals: number
  matchScore?: number
  postedAt: string
  updatedAt: string
}

export interface Client {
  id?: string
  name: string
  rating: number
  jobsPosted: number
  jobsCompleted?: number
  memberSince: string
  verified: boolean
  avatar?: string
}

export interface Proposal {
  id: string
  jobId: string
  contractorId: string
  proposedBudget: number
  timeline: string
  coverLetter: string
  status: ProposalStatus
  submittedAt: string
  viewedByClient: boolean
  messages: number
  acceptedAt?: string
  rejectedAt?: string
  rejectionReason?: string
  interviewDate?: string
}

export type ProposalStatus =
  | 'pending'
  | 'under_review'
  | 'interview_scheduled'
  | 'accepted'
  | 'rejected'

export interface ContractorReputation {
  contractorId: string
  overallRating: number
  totalJobs: number
  completedJobs: number
  completionRate: number
  onTimeRate: number
  responseRate: number
  reviews: Review[]
}

export interface Review {
  id: string
  jobId: string
  clientId: string
  contractorId: string
  rating: number
  comment: string
  categories: {
    quality: number
    communication: number
    timeliness: number
    professionalism: number
  }
  createdAt: string
}

// Blueprint Analysis Types
export interface BlueprintAnalysis {
  detectedComponents: DetectedComponent[]
  areaSqft: number
  suggestedMaterials: SuggestedMaterial[]
  confidenceScore: number
  warnings?: string[]
}

export interface DetectedComponent {
  type: string
  quantity: number
  location: string
  confidence: number
}

export interface SuggestedMaterial {
  name: string
  quantity: number
  unit: string
  estimatedCost: number
  category: string
}

// API Response Types
export interface ApiResponse<T> {
  data?: T
  error?: string
  message?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  limit: number
  offset: number
}

// Form Types
export interface LoginForm {
  email: string
  password: string
}

export interface SignupForm {
  email: string
  password: string
  confirmPassword: string
  firstName: string
  lastName: string
  company?: string
  trade?: string
}

export interface ProjectForm {
  name: string
  description: string
  trade: string
  location: Location
  startDate: string
  endDate: string
  budget: number
}

export interface EstimateForm {
  projectName: string
  trade: string
  materials: EstimateMaterial[]
  laborHours: number
  laborRate: number
  markup: number
  notes?: string
}

export interface ProposalForm {
  jobId: string
  proposedBudget: number
  timeline: string
  coverLetter: string
}
