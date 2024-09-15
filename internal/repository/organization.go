package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"project/internal/models"
)

func (r *Repository) GetOrganizationByID(organizationUUID uuid.UUID) (*models.Organization, error) {
	organization := &models.Organization{}
	err := r.db.Table("organization").
		Where("id = ?", organizationUUID).
		First(organization).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return organization, nil
}

func (r *Repository) CountOrganizationEmployees(organizationUUID uuid.UUID) (int64, error) {
	var result int64
	err := r.db.Table("organization_responsible").
		Where("organization_id = ?", organizationUUID).
		Count(&result).Error
	return result, err
}
