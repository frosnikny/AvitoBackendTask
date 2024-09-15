package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"project/internal/models"
)

func (r *Repository) AddBidVersion(bidVersion *models.BidVersion) error {
	err := r.db.Create(&bidVersion).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetBidVersionByID(bidVersionUUID uuid.UUID, version int) (*models.BidVersion, error) {
	var err error
	var bidVersion *models.BidVersion
	err = r.db.Table("bid_versions").
		Where("version = ?", version).
		Where("bid_id = ?", bidVersionUUID).
		First(&bidVersion).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return bidVersion, nil
}
