package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/trademaster/backend/estimation-service/internal/models"
	"github.com/trademaster/backend/estimation-service/internal/repository"
	"github.com/trademaster/backend/shared/config"
)

type EstimationService interface {
	CreateEstimate(estimate *models.Estimate) error
	GetEstimate(id uuid.UUID) (*models.Estimate, error)
	ListEstimates(contractorID uuid.UUID, limit, offset int) ([]*models.Estimate, error)
	UpdateEstimate(estimate *models.Estimate) error
	DeleteEstimate(id uuid.UUID) error
	UploadBlueprint(estimateID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*models.Blueprint, error)
	AnalyzeBlueprint(blueprintID string) (*models.BlueprintAnalysis, error)
	CalculateEstimate(blueprintAnalysis *models.BlueprintAnalysis, trade string) (*models.Estimate, error)
}

type estimationService struct {
	estimateRepo      repository.EstimateRepository
	blueprintRepo     repository.BlueprintRepository
	materialPriceRepo repository.MaterialPriceRepository
	config            *config.Config
}

func NewEstimationService(
	estimateRepo repository.EstimateRepository,
	blueprintRepo repository.BlueprintRepository,
	materialPriceRepo repository.MaterialPriceRepository,
	config *config.Config,
) EstimationService {
	return &estimationService{
		estimateRepo:      estimateRepo,
		blueprintRepo:     blueprintRepo,
		materialPriceRepo: materialPriceRepo,
		config:            config,
	}
}

func (s *estimationService) CreateEstimate(estimate *models.Estimate) error {
	// Calculate totals
	estimate.LaborTotal = estimate.LaborHours * estimate.LaborRate

	var materialsTotal float64
	for _, item := range estimate.Items {
		item.Total = item.Quantity * item.UnitPrice
		if item.Category == "material" {
			materialsTotal += item.Total
		}
	}
	estimate.MaterialsTotal = materialsTotal
	estimate.Subtotal = estimate.LaborTotal + estimate.MaterialsTotal
	estimate.Tax = estimate.Subtotal * 0.08 // 8% tax
	estimate.Total = estimate.Subtotal + estimate.Tax

	// Set valid until (30 days from now)
	estimate.ValidUntil = time.Now().AddDate(0, 0, 30)

	return s.estimateRepo.Create(estimate)
}

func (s *estimationService) GetEstimate(id uuid.UUID) (*models.Estimate, error) {
	return s.estimateRepo.FindByID(id)
}

func (s *estimationService) ListEstimates(contractorID uuid.UUID, limit, offset int) ([]*models.Estimate, error) {
	return s.estimateRepo.FindByContractorID(contractorID, limit, offset)
}

func (s *estimationService) UpdateEstimate(estimate *models.Estimate) error {
	// Recalculate totals
	estimate.LaborTotal = estimate.LaborHours * estimate.LaborRate

	var materialsTotal float64
	for _, item := range estimate.Items {
		item.Total = item.Quantity * item.UnitPrice
		if item.Category == "material" {
			materialsTotal += item.Total
		}
	}
	estimate.MaterialsTotal = materialsTotal
	estimate.Subtotal = estimate.LaborTotal + estimate.MaterialsTotal
	estimate.Tax = estimate.Subtotal * 0.08
	estimate.Total = estimate.Subtotal + estimate.Tax

	return s.estimateRepo.Update(estimate)
}

func (s *estimationService) DeleteEstimate(id uuid.UUID) error {
	return s.estimateRepo.Delete(id)
}

func (s *estimationService) UploadBlueprint(estimateID uuid.UUID, file multipart.File, header *multipart.FileHeader) (*models.Blueprint, error) {
	// TODO: Upload to S3
	// For now, we'll store metadata in MongoDB

	blueprint := &models.Blueprint{
		ID:         uuid.New().String(),
		EstimateID: estimateID.String(),
		FileName:   header.Filename,
		FileSize:   header.Size,
		MimeType:   header.Header.Get("Content-Type"),
		FileURL:    fmt.Sprintf("s3://blueprints/%s/%s", estimateID.String(), header.Filename),
		UploadedAt: time.Now(),
	}

	if err := s.blueprintRepo.Create(blueprint); err != nil {
		return nil, err
	}

	return blueprint, nil
}

func (s *estimationService) AnalyzeBlueprint(blueprintID string) (*models.BlueprintAnalysis, error) {
	// Get blueprint
	blueprint, err := s.blueprintRepo.FindByID(blueprintID)
	if err != nil {
		return nil, err
	}

	// Call ML service for analysis
	analysis, err := s.callBlueprintAnalysisML(blueprint.FileURL)
	if err != nil {
		return nil, err
	}

	// Update blueprint with analysis
	blueprint.Analysis = analysis
	now := time.Now()
	blueprint.AnalyzedAt = &now

	if err := s.blueprintRepo.Update(blueprint); err != nil {
		return nil, err
	}

	return analysis, nil
}

func (s *estimationService) CalculateEstimate(blueprintAnalysis *models.BlueprintAnalysis, trade string) (*models.Estimate, error) {
	// Call ML service for cost prediction
	prediction, err := s.callCostPredictionML(blueprintAnalysis, trade)
	if err != nil {
		return nil, err
	}

	// Build estimate from prediction
	estimate := &models.Estimate{
		Type:            trade,
		LaborHours:      prediction.LaborHours,
		LaborRate:       prediction.LaborRate,
		ConfidenceScore: prediction.Confidence,
		Items:           []models.EstimateItem{},
	}

	// Add materials from detected items
	for _, item := range blueprintAnalysis.DetectedItems {
		// Get price from database
		price, err := s.materialPriceRepo.FindByName(item.Label, trade)
		if err != nil {
			// Use default price if not found
			price = &models.MaterialPrice{
				Name:  item.Label,
				Unit:  "each",
				Price: 10.0,
			}
		}

		estimateItem := models.EstimateItem{
			Name:      item.Label,
			Quantity:  float64(item.Quantity),
			Unit:      price.Unit,
			UnitPrice: price.Price,
			Category:  "material",
		}
		estimate.Items = append(estimate.Items, estimateItem)
	}

	return estimate, nil
}

// ML Service Calls

type BlueprintAnalysisRequest struct {
	FileURL string `json:"file_url"`
}

func (s *estimationService) callBlueprintAnalysisML(fileURL string) (*models.BlueprintAnalysis, error) {
	// Call Python ML service
	mlServiceURL := "http://localhost:8000/api/v1/cv/analyze-blueprint"

	reqBody := BlueprintAnalysisRequest{FileURL: fileURL}
	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post(mlServiceURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call ML service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ML service error: %s", string(body))
	}

	var analysis models.BlueprintAnalysis
	if err := json.NewDecoder(resp.Body).Decode(&analysis); err != nil {
		return nil, err
	}

	return &analysis, nil
}

type CostPredictionRequest struct {
	BlueprintAnalysis *models.BlueprintAnalysis `json:"blueprint_analysis"`
	Trade             string                    `json:"trade"`
}

type CostPredictionResponse struct {
	LaborHours float64 `json:"labor_hours"`
	LaborRate  float64 `json:"labor_rate"`
	Confidence float64 `json:"confidence"`
}

func (s *estimationService) callCostPredictionML(analysis *models.BlueprintAnalysis, trade string) (*CostPredictionResponse, error) {
	// Call Python ML service
	mlServiceURL := "http://localhost:8000/api/v1/predict/duration"

	reqBody := CostPredictionRequest{
		BlueprintAnalysis: analysis,
		Trade:             trade,
	}
	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post(mlServiceURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call ML service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("ML service error")
	}

	var prediction CostPredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
		return nil, err
	}

	return &prediction, nil
}
