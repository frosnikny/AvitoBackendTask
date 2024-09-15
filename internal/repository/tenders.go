package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"project/internal/models"
)

func (r *Repository) AddTender(tender *models.Tender) error {
	err := r.db.Create(&tender).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SaveTender(tender *models.Tender) error {
	err := r.db.Save(&tender).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetTenderByID(tenderUUID uuid.UUID) (*models.Tender, error) {
	var err error
	var tender *models.Tender
	err = r.db.Table("tenders").
		Where("id = ?", tenderUUID).
		First(&tender).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return tender, nil
}

func (r *Repository) GetFilteredTenders(limit int, offset int, serviceType string) ([]models.Tender, error) {
	var tenders []models.Tender
	query := r.db.Table("tenders").
		Where("status = ?", models.TenderPublished).
		Limit(limit).
		Offset(offset).
		Order("name asc")

	if serviceType != "" {
		query = query.Where("service_type = ?", serviceType)
	}

	if err := query.Find(&tenders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tenders, nil
		}
		return nil, err
	}

	return tenders, nil
}

func (r *Repository) GetUserTenders(limit int, offset int, authorUUID uuid.UUID) ([]models.Tender, error) {
	var err error
	var tenders []models.Tender
	err = r.db.Table("tenders").
		Where("author_id = ?", authorUUID).
		Limit(limit).
		Offset(offset).
		Order("name asc").
		Find(&tenders).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tenders, nil
		}
		return nil, err
	}

	return tenders, nil
}
