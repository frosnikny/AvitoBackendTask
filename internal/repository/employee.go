package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"project/internal/models"
)

func (r *Repository) GetEmployeeByID(employeeUUID uuid.UUID) (*models.Employee, error) {
	employee := &models.Employee{}
	err := r.db.Table("employee").
		Where("id = ?", employeeUUID).
		First(employee).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return employee, nil
}

func (r *Repository) GetEmployeeByUsername(username string) (*models.Employee, error) {
	employee := &models.Employee{}
	err := r.db.Table("employee").
		Where("username = ?", username).
		First(employee).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return employee, nil
}

func (r *Repository) GetEmployeeOrganizations(employeeUUID uuid.UUID) ([]models.Organization, error) {
	var organizations []models.Organization
	err := r.db.Table("organization").
		Joins("JOIN organization_responsible ON organization.id = organization_responsible.organization_id").
		Where("organization_responsible.user_id = ?", employeeUUID).
		Find(&organizations).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return organizations, nil
}
