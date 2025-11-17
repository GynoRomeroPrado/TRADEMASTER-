package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/trademaster/backend/training-service/internal/models"
	"github.com/trademaster/backend/training-service/internal/repository"
	"github.com/trademaster/backend/shared/config"
)

type TrainingService interface {
	// Courses
	GetCourse(id string) (*models.Course, error)
	ListCourses(filters map[string]interface{}, limit, offset int) ([]*models.Course, error)
	GetCoursesByTrade(trade string) ([]*models.Course, error)

	// Enrollments
	EnrollInCourse(userID uuid.UUID, courseID string) (*models.Enrollment, error)
	GetEnrollment(id uuid.UUID) (*models.Enrollment, error)
	GetUserEnrollments(userID uuid.UUID) ([]*models.Enrollment, error)
	UpdateProgress(enrollmentID uuid.UUID, lessonID string, progress float64, timeSpent int) error
	CompleteLesson(enrollmentID uuid.UUID, lessonID string, score *float64) error
	GetEnrollmentProgress(enrollmentID uuid.UUID) (*models.Enrollment, error)

	// Certifications
	IssueCertification(userID uuid.UUID, courseID string) (*models.Certification, error)
	GetUserCertifications(userID uuid.UUID) ([]*models.Certification, error)
	VerifyCertification(certificateNumber string) (*models.Certification, error)
	RevokeCertification(id uuid.UUID) error
}

type trainingService struct {
	courseRepo        repository.CourseRepository
	enrollmentRepo    repository.EnrollmentRepository
	certificationRepo repository.CertificationRepository
	config            *config.Config
}

func NewTrainingService(
	courseRepo repository.CourseRepository,
	enrollmentRepo repository.EnrollmentRepository,
	certificationRepo repository.CertificationRepository,
	config *config.Config,
) TrainingService {
	return &trainingService{
		courseRepo:        courseRepo,
		enrollmentRepo:    enrollmentRepo,
		certificationRepo: certificationRepo,
		config:            config,
	}
}

// Course Management

func (s *trainingService) GetCourse(id string) (*models.Course, error) {
	return s.courseRepo.FindByID(id)
}

func (s *trainingService) ListCourses(filters map[string]interface{}, limit, offset int) ([]*models.Course, error) {
	return s.courseRepo.FindAll(filters, limit, offset)
}

func (s *trainingService) GetCoursesByTrade(trade string) ([]*models.Course, error) {
	return s.courseRepo.FindByTrade(trade)
}

// Enrollment Management

func (s *trainingService) EnrollInCourse(userID uuid.UUID, courseID string) (*models.Enrollment, error) {
	// Check if already enrolled
	existing, err := s.enrollmentRepo.FindByUserAndCourse(userID, courseID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil // Already enrolled
	}

	// Verify course exists
	course, err := s.courseRepo.FindByID(courseID)
	if err != nil {
		return nil, err
	}

	// Create enrollment
	enrollment := &models.Enrollment{
		UserID:     userID,
		CourseID:   courseID,
		Status:     "active",
		Progress:   0,
		EnrolledAt: time.Now(),
	}

	// Initialize lesson progress
	for _, module := range course.Modules {
		for _, lesson := range module.Lessons {
			lessonProgress := models.LessonProgress{
				LessonID:  lesson.ID,
				Status:    "not_started",
				Progress:  0,
				TimeSpent: 0,
			}
			enrollment.ProgressData = append(enrollment.ProgressData, lessonProgress)
		}
	}

	err = s.enrollmentRepo.Create(enrollment)
	if err != nil {
		return nil, err
	}

	return enrollment, nil
}

func (s *trainingService) GetEnrollment(id uuid.UUID) (*models.Enrollment, error) {
	return s.enrollmentRepo.FindByID(id)
}

func (s *trainingService) GetUserEnrollments(userID uuid.UUID) ([]*models.Enrollment, error) {
	return s.enrollmentRepo.FindByUserID(userID)
}

func (s *trainingService) UpdateProgress(enrollmentID uuid.UUID, lessonID string, progress float64, timeSpent int) error {
	enrollment, err := s.enrollmentRepo.FindByID(enrollmentID)
	if err != nil {
		return err
	}

	// Update lesson progress
	found := false
	for i := range enrollment.ProgressData {
		if enrollment.ProgressData[i].LessonID == lessonID {
			enrollment.ProgressData[i].Progress = progress
			enrollment.ProgressData[i].TimeSpent += timeSpent
			if progress >= 100 {
				enrollment.ProgressData[i].Status = "completed"
				now := time.Now()
				enrollment.ProgressData[i].CompletedAt = &now
			} else if progress > 0 {
				enrollment.ProgressData[i].Status = "in_progress"
			}
			found = true
			break
		}
	}

	if !found {
		return errors.New("lesson not found in enrollment")
	}

	// Calculate overall progress
	totalLessons := len(enrollment.ProgressData)
	completedLessons := 0
	for _, lp := range enrollment.ProgressData {
		if lp.Status == "completed" {
			completedLessons++
		}
	}
	enrollment.Progress = float64(completedLessons) / float64(totalLessons) * 100

	// Update last accessed
	now := time.Now()
	enrollment.LastAccessedAt = &now

	// Check if course completed
	if enrollment.Progress >= 100 && enrollment.Status != "completed" {
		enrollment.Status = "completed"
		enrollment.CompletedAt = &now
	}

	return s.enrollmentRepo.Update(enrollment)
}

func (s *trainingService) CompleteLesson(enrollmentID uuid.UUID, lessonID string, score *float64) error {
	enrollment, err := s.enrollmentRepo.FindByID(enrollmentID)
	if err != nil {
		return err
	}

	// Mark lesson as completed
	for i := range enrollment.ProgressData {
		if enrollment.ProgressData[i].LessonID == lessonID {
			enrollment.ProgressData[i].Status = "completed"
			enrollment.ProgressData[i].Progress = 100
			enrollment.ProgressData[i].Score = score
			now := time.Now()
			enrollment.ProgressData[i].CompletedAt = &now
			break
		}
	}

	// Calculate overall progress
	totalLessons := len(enrollment.ProgressData)
	completedLessons := 0
	for _, lp := range enrollment.ProgressData {
		if lp.Status == "completed" {
			completedLessons++
		}
	}
	enrollment.Progress = float64(completedLessons) / float64(totalLessons) * 100

	// Check if course completed
	if enrollment.Progress >= 100 && enrollment.Status != "completed" {
		enrollment.Status = "completed"
		now := time.Now()
		enrollment.CompletedAt = &now
	}

	return s.enrollmentRepo.Update(enrollment)
}

func (s *trainingService) GetEnrollmentProgress(enrollmentID uuid.UUID) (*models.Enrollment, error) {
	return s.enrollmentRepo.FindByID(enrollmentID)
}

// Certification Management

func (s *trainingService) IssueCertification(userID uuid.UUID, courseID string) (*models.Certification, error) {
	// Verify course is completed
	enrollment, err := s.enrollmentRepo.FindByUserAndCourse(userID, courseID)
	if err != nil {
		return nil, err
	}
	if enrollment == nil {
		return nil, errors.New("user not enrolled in course")
	}
	if enrollment.Status != "completed" {
		return nil, errors.New("course not completed")
	}

	// Get course details
	course, err := s.courseRepo.FindByID(courseID)
	if err != nil {
		return nil, err
	}

	if course.Certification == nil || !course.Certification.Enabled {
		return nil, errors.New("course does not offer certification")
	}

	// Check if already certified
	certifications, _ := s.certificationRepo.FindByUserID(userID)
	for _, cert := range certifications {
		if cert.CourseID == courseID && cert.Status == "active" {
			return cert, nil // Already certified
		}
	}

	// Create certification
	certification := &models.Certification{
		UserID:            userID,
		CourseID:          courseID,
		CertificateType:   course.Certification.Name,
		IssuedAt:          time.Now(),
		Status:            "active",
		VerificationURL:   "", // Will be set after creation
		BlockchainTxHash:  "", // TODO: Implement blockchain verification
	}

	err = s.certificationRepo.Create(certification)
	if err != nil {
		return nil, err
	}

	// Set verification URL
	certification.VerificationURL = s.generateVerificationURL(certification.CertificateNumber)
	s.certificationRepo.Update(certification)

	// TODO: Record on blockchain for verification
	// s.recordOnBlockchain(certification)

	return certification, nil
}

func (s *trainingService) GetUserCertifications(userID uuid.UUID) ([]*models.Certification, error) {
	return s.certificationRepo.FindByUserID(userID)
}

func (s *trainingService) VerifyCertification(certificateNumber string) (*models.Certification, error) {
	return s.certificationRepo.FindByCertificateNumber(certificateNumber)
}

func (s *trainingService) RevokeCertification(id uuid.UUID) error {
	return s.certificationRepo.Revoke(id)
}

// Helper functions

func (s *trainingService) generateVerificationURL(certificateNumber string) string {
	baseURL := s.config.Environment
	return baseURL + "/verify/" + certificateNumber
}
