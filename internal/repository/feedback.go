package repository

import (
	"github.com/google/uuid"
	"project/internal/models"
)

func (r *Repository) AddFeedback(feedback *models.Feedback) error {
	err := r.db.Create(&feedback).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetFeedbacksByBid(bidUUID uuid.UUID) ([]models.Feedback, error) {
	var err error
	var feedbacks []models.Feedback
	err = r.db.Table("feedbacks").
		Where("bid_id = ?", bidUUID).
		Find(&feedbacks).Error

	if err != nil {
		return nil, err
	}

	return feedbacks, nil
}
