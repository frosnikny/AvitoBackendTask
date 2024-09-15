package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"project/internal/models"
)

func (r *Repository) AddTenderVersion(tenderVersion *models.TenderVersion) error {
	err := r.db.Create(&tenderVersion).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetTenderVersionByID(tenderVersionUUID uuid.UUID, version int) (*models.TenderVersion, error) {
	var err error
	var tenderVersion *models.TenderVersion
	err = r.db.Table("tender_versions").
		Where("version = ?", version).
		Where("tender_id = ?", tenderVersionUUID).
		First(&tenderVersion).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return tenderVersion, nil
}
