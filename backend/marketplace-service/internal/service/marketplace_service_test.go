package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trademaster/backend/marketplace-service/internal/models"
)

// MockMarketplaceRepository is a mock implementation of the repository
type MockMarketplaceRepository struct {
	mock.Mock
}

func (m *MockMarketplaceRepository) CreateJob(job *models.Job) error {
	args := m.Called(job)
	return args.Error(0)
}

func (m *MockMarketplaceRepository) GetJobByID(id uuid.UUID) (*models.Job, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Job), args.Error(1)
}

func (m *MockMarketplaceRepository) ListJobs(filters map[string]interface{}) ([]models.Job, error) {
	args := m.Called(filters)
	return args.Get(0).([]models.Job), args.Error(1)
}

func (m *MockMarketplaceRepository) UpdateJob(job *models.Job) error {
	args := m.Called(job)
	return args.Error(0)
}

func (m *MockMarketplaceRepository) CreateProposal(proposal *models.Proposal) error {
	args := m.Called(proposal)
	return args.Error(0)
}

func (m *MockMarketplaceRepository) GetProposalByID(id uuid.UUID) (*models.Proposal, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Proposal), args.Error(1)
}

func (m *MockMarketplaceRepository) ListProposalsByJob(jobID uuid.UUID) ([]models.Proposal, error) {
	args := m.Called(jobID)
	return args.Get(0).([]models.Proposal), args.Error(1)
}

func (m *MockMarketplaceRepository) UpdateProposal(proposal *models.Proposal) error {
	args := m.Called(proposal)
	return args.Error(0)
}

func (m *MockMarketplaceRepository) GetContractorReputation(contractorID uuid.UUID) (*models.ContractorReputation, error) {
	args := m.Called(contractorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ContractorReputation), args.Error(1)
}

func (m *MockMarketplaceRepository) UpdateReputation(reputation *models.ContractorReputation) error {
	args := m.Called(reputation)
	return args.Error(0)
}

func TestCalculateMatchScore(t *testing.T) {
	tests := []struct {
		name           string
		reputation     *models.ContractorReputation
		expectedScore  float64
		expectedError  bool
	}{
		{
			name: "Perfect contractor - 100% match",
			reputation: &models.ContractorReputation{
				OverallRating:  5.0,
				CompletionRate: 1.0,
				OnTimeRate:     1.0,
				ResponseRate:   1.0,
				TotalJobs:      100,
			},
			expectedScore: 100.0,
			expectedError: false,
		},
		{
			name: "Good contractor - ~90% match",
			reputation: &models.ContractorReputation{
				OverallRating:  4.5,
				CompletionRate: 0.95,
				OnTimeRate:     0.90,
				ResponseRate:   0.85,
				TotalJobs:      50,
			},
			expectedScore: 90.25,
			expectedError: false,
		},
		{
			name: "Average contractor - ~70% match",
			reputation: &models.ContractorReputation{
				OverallRating:  3.5,
				CompletionRate: 0.80,
				OnTimeRate:     0.70,
				ResponseRate:   0.60,
				TotalJobs:      20,
			},
			expectedScore: 68.0,
			expectedError: false,
		},
		{
			name: "Poor contractor - ~50% match",
			reputation: &models.ContractorReputation{
				OverallRating:  2.5,
				CompletionRate: 0.60,
				OnTimeRate:     0.50,
				ResponseRate:   0.40,
				TotalJobs:      10,
			},
			expectedScore: 46.0,
			expectedError: false,
		},
		{
			name: "New contractor with good early performance",
			reputation: &models.ContractorReputation{
				OverallRating:  4.8,
				CompletionRate: 1.0,
				OnTimeRate:     1.0,
				ResponseRate:   1.0,
				TotalJobs:      5,
			},
			expectedScore: 93.8,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockMarketplaceRepository)
			service := &marketplaceService{
				repo: mockRepo,
			}

			jobID := uuid.New()
			contractorID := uuid.New()

			// Setup mock expectations
			mockRepo.On("GetContractorReputation", contractorID).Return(tt.reputation, nil)

			// Calculate match score
			score, err := service.CalculateMatchScore(jobID, contractorID)

			// Assertions
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expectedScore, score, 0.1, "Match score should be within 0.1 of expected")
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateJob(t *testing.T) {
	mockRepo := new(MockMarketplaceRepository)
	service := &marketplaceService{
		repo: mockRepo,
	}

	job := &models.Job{
		Title:       "Test Electrical Job",
		Description: "Test job description",
		Trade:       "electrical",
		Budget:      50000,
		ClientID:    uuid.New(),
		Status:      "open",
		Location: models.Location{
			City:  "San Francisco",
			State: "CA",
		},
	}

	mockRepo.On("CreateJob", job).Return(nil)

	err := service.CreateJob(job)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, job.ID)
	mockRepo.AssertExpectations(t)
}

func TestSubmitProposal(t *testing.T) {
	mockRepo := new(MockMarketplaceRepository)
	service := &marketplaceService{
		repo: mockRepo,
	}

	jobID := uuid.New()
	contractorID := uuid.New()

	job := &models.Job{
		ID:          jobID,
		Title:       "Test Job",
		Status:      "open",
		ClientID:    uuid.New(),
		Budget:      50000,
	}

	proposal := &models.Proposal{
		JobID:          jobID,
		ContractorID:   contractorID,
		ProposedBudget: 48000,
		Timeline:       "6 weeks",
		CoverLetter:    "I am interested in this job",
		Status:         "pending",
	}

	mockRepo.On("GetJobByID", jobID).Return(job, nil)
	mockRepo.On("CreateProposal", proposal).Return(nil)

	err := service.SubmitProposal(proposal)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, proposal.ID)
	assert.Equal(t, "pending", proposal.Status)
	mockRepo.AssertExpectations(t)
}

func TestAcceptProposal(t *testing.T) {
	mockRepo := new(MockMarketplaceRepository)
	service := &marketplaceService{
		repo: mockRepo,
	}

	proposalID := uuid.New()
	jobID := uuid.New()

	proposal := &models.Proposal{
		ID:           proposalID,
		JobID:        jobID,
		ContractorID: uuid.New(),
		Status:       "pending",
	}

	job := &models.Job{
		ID:     jobID,
		Status: "open",
	}

	mockRepo.On("GetProposalByID", proposalID).Return(proposal, nil)
	mockRepo.On("GetJobByID", jobID).Return(job, nil)
	mockRepo.On("UpdateProposal", proposal).Return(nil)
	mockRepo.On("UpdateJob", job).Return(nil)

	err := service.AcceptProposal(proposalID)

	assert.NoError(t, err)
	assert.Equal(t, "accepted", proposal.Status)
	assert.Equal(t, "in_progress", job.Status)
	assert.NotNil(t, proposal.AcceptedAt)
	mockRepo.AssertExpectations(t)
}

func TestRejectProposal(t *testing.T) {
	mockRepo := new(MockMarketplaceRepository)
	service := &marketplaceService{
		repo: mockRepo,
	}

	proposalID := uuid.New()
	reason := "Selected another contractor"

	proposal := &models.Proposal{
		ID:     proposalID,
		Status: "pending",
	}

	mockRepo.On("GetProposalByID", proposalID).Return(proposal, nil)
	mockRepo.On("UpdateProposal", proposal).Return(nil)

	err := service.RejectProposal(proposalID, reason)

	assert.NoError(t, err)
	assert.Equal(t, "rejected", proposal.Status)
	assert.Equal(t, reason, proposal.RejectionReason)
	assert.NotNil(t, proposal.RejectedAt)
	mockRepo.AssertExpectations(t)
}

func TestGetMatchingJobs(t *testing.T) {
	mockRepo := new(MockMarketplaceRepository)
	service := &marketplaceService{
		repo: mockRepo,
	}

	contractorID := uuid.New()
	trade := "electrical"

	jobs := []models.Job{
		{
			ID:     uuid.New(),
			Title:  "Job 1",
			Trade:  trade,
			Status: "open",
		},
		{
			ID:     uuid.New(),
			Title:  "Job 2",
			Trade:  trade,
			Status: "open",
		},
	}

	filters := map[string]interface{}{
		"trade":  trade,
		"status": "open",
	}

	mockRepo.On("ListJobs", filters).Return(jobs, nil)

	result, err := service.GetMatchingJobs(contractorID, trade)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestUpdateContractorReputation(t *testing.T) {
	mockRepo := new(MockMarketplaceRepository)
	service := &marketplaceService{
		repo: mockRepo,
	}

	contractorID := uuid.New()
	rating := 4.5

	existingReputation := &models.ContractorReputation{
		ContractorID:   contractorID,
		OverallRating:  4.0,
		TotalJobs:      10,
		CompletedJobs:  9,
		CompletionRate: 0.9,
	}

	mockRepo.On("GetContractorReputation", contractorID).Return(existingReputation, nil)
	mockRepo.On("UpdateReputation", existingReputation).Return(nil)

	err := service.UpdateContractorReputation(contractorID, rating, true, true)

	assert.NoError(t, err)
	assert.Equal(t, 11, existingReputation.TotalJobs)
	assert.Equal(t, 10, existingReputation.CompletedJobs)
	assert.Greater(t, existingReputation.OverallRating, 4.0)
	mockRepo.AssertExpectations(t)
}
