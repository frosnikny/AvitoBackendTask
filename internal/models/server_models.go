package models

import (
	"time"

	"github.com/google/uuid"
)

// Employee Модель пользователь
type Employee struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username  string    `gorm:"type:varchar(50);unique;not null"`
	FirstName string    `gorm:"type:varchar(50)"`
	LastName  string    `gorm:"type:varchar(50)"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (Employee) TableName() string {
	return "employee"
}

// OrganizationType Тип организации (ENUM)
type OrganizationType string

const (
	OrganizationTypeIE  OrganizationType = "IE"
	OrganizationTypeLLC OrganizationType = "LLC"
	OrganizationTypeJSC OrganizationType = "JSC"
)

// Organization Модель организации
type Organization struct {
	ID          uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string           `gorm:"type:varchar(100);not null"`
	Description string           `gorm:"type:text"`
	Type        OrganizationType `gorm:"type:organization_type"` // ENUM для типа организации
	CreatedAt   time.Time        `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time        `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (Organization) TableName() string {
	return "organization"
}

// OrganizationResponsible Модель для ответственного за организацию
type OrganizationResponsible struct {
	ID             uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrganizationID uuid.UUID     `gorm:"type:uuid;not null"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	UserID         uuid.UUID     `gorm:"type:uuid;not null"`
	Employee       *Employee     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
