package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"project/internal/models"
)

func (r *Repository) AddBid(bid *models.Bid) error {
	err := r.db.Create(&bid).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SaveBid(bid *models.Bid) error {
	err := r.db.Save(&bid).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetBidByID(bidUUID uuid.UUID) (*models.Bid, error) {
	var err error
	var bid *models.Bid
	err = r.db.Table("bids").
		Where("id = ?", bidUUID).
		First(&bid).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return bid, nil
}

func (r *Repository) GetUserBids(limit int, offset int, authorUUID uuid.UUID) ([]models.Bid, error) {
	var err error
	var bids []models.Bid
	err = r.db.Table("bids").
		Where("employee_id = ?", authorUUID).
		Limit(limit).
		Offset(offset).
		Order("name asc").
		Find(&bids).Error

	if err != nil {
		return nil, err
	}

	return bids, nil
}

func (r *Repository) GetUserTenderBids(authorUUID uuid.UUID, tenderUUID uuid.UUID) ([]models.Bid, error) {
	var err error
	var bids []models.Bid
	err = r.db.Table("bids").
		Where("employee_id = ?", authorUUID).
		Where("tender_id = ?", tenderUUID).
		Order("name asc").
		Find(&bids).Error

	if err != nil {
		return nil, err
	}

	return bids, nil
}

func (r *Repository) GetBidsByTender(limit int, offset int, tenderUUID uuid.UUID) ([]models.Bid, error) {
	var err error
	var bids []models.Bid
	err = r.db.Table("bids").
		Where("tender_id = ?", tenderUUID).
		Limit(limit).
		Offset(offset).
		Order("name asc").
		Find(&bids).Error

	if err != nil {
		return nil, err
	}

	return bids, nil
}
