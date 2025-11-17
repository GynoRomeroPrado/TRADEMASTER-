package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trademaster/backend/estimation-service/internal/models"
)

// MockEstimationRepository is a mock implementation of the repository
type MockEstimationRepository struct {
	mock.Mock
}

func (m *MockEstimationRepository) CreateEstimate(estimate *models.Estimate) error {
	args := m.Called(estimate)
	return args.Error(0)
}

func (m *MockEstimationRepository) GetEstimateByID(id uuid.UUID) (*models.Estimate, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Estimate), args.Error(1)
}

func (m *MockEstimationRepository) ListEstimatesByUser(userID uuid.UUID) ([]models.Estimate, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Estimate), args.Error(1)
}

func (m *MockEstimationRepository) UpdateEstimate(estimate *models.Estimate) error {
	args := m.Called(estimate)
	return args.Error(0)
}

func (m *MockEstimationRepository) DeleteEstimate(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEstimationRepository) GetMaterialPrice(name, trade string) (*models.MaterialPrice, error) {
	args := m.Called(name, trade)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MaterialPrice), args.Error(1)
}

func (m *MockEstimationRepository) ListMaterialPrices(trade string) ([]models.MaterialPrice, error) {
	args := m.Called(trade)
	return args.Get(0).([]models.MaterialPrice), args.Error(1)
}

func TestCalculateEstimateTotal(t *testing.T) {
	tests := []struct {
		name           string
		estimate       *models.Estimate
		expectedTotal  float64
	}{
		{
			name: "Simple estimate with materials and labor",
			estimate: &models.Estimate{
				Materials: []models.Material{
					{Name: "Wire", Quantity: 10, UnitPrice: 25.0},
					{Name: "Outlet", Quantity: 5, UnitPrice: 15.0},
				},
				LaborHours: 8,
				LaborRate:  75.0,
			},
			expectedTotal: 925.0, // (10*25 + 5*15) + (8*75) = 325 + 600 = 925
		},
		{
			name: "Estimate with markup",
			estimate: &models.Estimate{
				Materials: []models.Material{
					{Name: "Cable", Quantity: 5, UnitPrice: 100.0},
				},
				LaborHours: 10,
				LaborRate:  80.0,
				Markup:     0.15, // 15% markup
			},
			expectedTotal: 1495.0, // ((5*100) + (10*80)) * 1.15 = 1300 * 1.15 = 1495
		},
		{
			name: "Estimate with no materials",
			estimate: &models.Estimate{
				Materials:  []models.Material{},
				LaborHours: 5,
				LaborRate:  100.0,
			},
			expectedTotal: 500.0, // 5*100 = 500
		},
		{
			name: "Estimate with no labor",
			estimate: &models.Estimate{
				Materials: []models.Material{
					{Name: "Equipment", Quantity: 1, UnitPrice: 5000.0},
				},
				LaborHours: 0,
				LaborRate:  0,
			},
			expectedTotal: 5000.0, // 1*5000 = 5000
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockEstimationRepository)
			service := &estimationService{
				repo: mockRepo,
			}

			total := service.calculateTotal(tt.estimate)

			assert.InDelta(t, tt.expectedTotal, total, 0.01, "Total should match expected value")
		})
	}
}

func TestCreateEstimate(t *testing.T) {
	mockRepo := new(MockEstimationRepository)
	service := &estimationService{
		repo: mockRepo,
	}

	estimate := &models.Estimate{
		UserID:      uuid.New(),
		ProjectName: "Office Rewiring",
		Trade:       "electrical",
		Materials: []models.Material{
			{Name: "Wire", Quantity: 100, UnitPrice: 2.50},
		},
		LaborHours: 40,
		LaborRate:  75.0,
	}

	mockRepo.On("CreateEstimate", estimate).Return(nil)

	err := service.CreateEstimate(estimate)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, estimate.ID)
	assert.Greater(t, estimate.TotalCost, 0.0)
	mockRepo.AssertExpectations(t)
}

func TestGetEstimateByID(t *testing.T) {
	mockRepo := new(MockEstimationRepository)
	service := &estimationService{
		repo: mockRepo,
	}

	estimateID := uuid.New()
	expectedEstimate := &models.Estimate{
		ID:          estimateID,
		ProjectName: "Test Project",
		TotalCost:   1000.0,
	}

	mockRepo.On("GetEstimateByID", estimateID).Return(expectedEstimate, nil)

	result, err := service.GetEstimateByID(estimateID)

	assert.NoError(t, err)
	assert.Equal(t, estimateID, result.ID)
	assert.Equal(t, "Test Project", result.ProjectName)
	mockRepo.AssertExpectations(t)
}

func TestUpdateEstimate(t *testing.T) {
	mockRepo := new(MockEstimationRepository)
	service := &estimationService{
		repo: mockRepo,
	}

	estimate := &models.Estimate{
		ID:          uuid.New(),
		ProjectName: "Updated Project",
		Materials: []models.Material{
			{Name: "Updated Material", Quantity: 10, UnitPrice: 50.0},
		},
		LaborHours: 20,
		LaborRate:  80.0,
	}

	mockRepo.On("GetEstimateByID", estimate.ID).Return(estimate, nil)
	mockRepo.On("UpdateEstimate", estimate).Return(nil)

	err := service.UpdateEstimate(estimate)

	assert.NoError(t, err)
	assert.Greater(t, estimate.TotalCost, 0.0)
	mockRepo.AssertExpectations(t)
}

func TestDeleteEstimate(t *testing.T) {
	mockRepo := new(MockEstimationRepository)
	service := &estimationService{
		repo: mockRepo,
	}

	estimateID := uuid.New()

	mockRepo.On("DeleteEstimate", estimateID).Return(nil)

	err := service.DeleteEstimate(estimateID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetMaterialSuggestions(t *testing.T) {
	mockRepo := new(MockEstimationRepository)
	service := &estimationService{
		repo: mockRepo,
	}

	trade := "electrical"
	materials := []models.MaterialPrice{
		{
			Name:     "12/2 NM-B Cable",
			Category: "wiring",
			Unit:     "roll",
			Price:    89.99,
			Trade:    trade,
		},
		{
			Name:     "15A Circuit Breaker",
			Category: "breakers",
			Unit:     "each",
			Price:    12.50,
			Trade:    trade,
		},
	}

	mockRepo.On("ListMaterialPrices", trade).Return(materials, nil)

	result, err := service.GetMaterialSuggestions(trade)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "12/2 NM-B Cable", result[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestCalculateLaborCost(t *testing.T) {
	tests := []struct {
		name          string
		hours         float64
		rate          float64
		expectedCost  float64
	}{
		{
			name:         "Standard 8-hour day",
			hours:        8,
			rate:         75.0,
			expectedCost: 600.0,
		},
		{
			name:         "Half day",
			hours:        4,
			rate:         100.0,
			expectedCost: 400.0,
		},
		{
			name:         "Full week",
			hours:        40,
			rate:         80.0,
			expectedCost: 3200.0,
		},
		{
			name:         "Decimal hours",
			hours:        6.5,
			rate:         90.0,
			expectedCost: 585.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := tt.hours * tt.rate
			assert.InDelta(t, tt.expectedCost, cost, 0.01)
		})
	}
}

func TestCalculateMaterialsCost(t *testing.T) {
	tests := []struct {
		name          string
		materials     []models.Material
		expectedCost  float64
	}{
		{
			name: "Single material",
			materials: []models.Material{
				{Name: "Wire", Quantity: 10, UnitPrice: 25.0},
			},
			expectedCost: 250.0,
		},
		{
			name: "Multiple materials",
			materials: []models.Material{
				{Name: "Wire", Quantity: 10, UnitPrice: 25.0},
				{Name: "Outlet", Quantity: 5, UnitPrice: 15.0},
				{Name: "Switch", Quantity: 3, UnitPrice: 10.0},
			},
			expectedCost: 355.0, // 250 + 75 + 30
		},
		{
			name:         "No materials",
			materials:    []models.Material{},
			expectedCost: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total := 0.0
			for _, m := range tt.materials {
				total += m.Quantity * m.UnitPrice
			}
			assert.InDelta(t, tt.expectedCost, total, 0.01)
		})
	}
}

func TestApplyMarkup(t *testing.T) {
	tests := []struct {
		name          string
		subtotal      float64
		markup        float64
		expectedTotal float64
	}{
		{
			name:          "15% markup",
			subtotal:      1000.0,
			markup:        0.15,
			expectedTotal: 1150.0,
		},
		{
			name:          "20% markup",
			subtotal:      5000.0,
			markup:        0.20,
			expectedTotal: 6000.0,
		},
		{
			name:          "No markup",
			subtotal:      2500.0,
			markup:        0.0,
			expectedTotal: 2500.0,
		},
		{
			name:          "10% markup",
			subtotal:      750.0,
			markup:        0.10,
			expectedTotal: 825.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total := tt.subtotal * (1 + tt.markup)
			assert.InDelta(t, tt.expectedTotal, total, 0.01)
		})
	}
}
