package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/trademaster/backend/project-service/internal/models"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *models.Task) error
	FindByID(id uuid.UUID) (*models.Task, error)
	FindByProjectID(projectID uuid.UUID) ([]*models.Task, error)
	FindByAssignedUser(userID uuid.UUID) ([]*models.Task, error)
	Update(task *models.Task) error
	Delete(id uuid.UUID) error
	CompleteTask(id uuid.UUID) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) FindByID(id uuid.UUID) (*models.Task, error) {
	var task models.Task
	err := r.db.Where("id = ?", id).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) FindByProjectID(projectID uuid.UUID) ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.Where("project_id = ?", projectID).
		Order("order ASC, created_at ASC").
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) FindByAssignedUser(userID uuid.UUID) ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.Where("assigned_to = ? AND status != ?", userID, "completed").
		Order("priority DESC, due_date ASC").
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Task{}, "id = ?", id).Error
}

func (r *taskRepository) CompleteTask(id uuid.UUID) error {
	now := gorm.Expr("NOW()")
	return r.db.Model(&models.Task{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "completed",
			"completed_at": now,
		}).Error
}
