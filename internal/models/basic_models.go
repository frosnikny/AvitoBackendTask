package models

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

const (
	TenderCreated int = iota
	TenderPublished
	TenderClosed
)

const (
	TenderCreatedStr   = "Created"
	TenderPublishedStr = "Published"
	TenderClosedStr    = "Closed"
)

// Tender Модель для тендера
type Tender struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;unique"`
	Name        string    `gorm:"type:varchar(50);default:null"`
	Description string    `gorm:"type:text;default:null"`
	Status      int       `gorm:"not null"`
	ServiceType string    `gorm:"type:varchar(255);default:null"`
	Version     int       `gorm:"not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`

	OrganizationID uuid.UUID     `gorm:"type:uuid;not null"`        // Связь с таблицей Organization через ID
	Organization   *Organization `gorm:"foreignKey:OrganizationID"` // Связь с организацией через foreign key
	AuthorID       uuid.UUID     `gorm:"type:uuid;not null"`        // Связь с таблицей Employee через ID
	Author         *Employee     `gorm:"foreignKey:AuthorID"`       // Связь с автором через foreign key
}

func (t *Tender) StatusToString() string {
	switch t.Status {
	case TenderCreated:
		return TenderCreatedStr
	case TenderPublished:
		return TenderPublishedStr
	case TenderClosed:
		return TenderClosedStr
	}
	return ""
}

func (t *Tender) StringToStatus(status string) error {
	switch status {
	case TenderCreatedStr:
		t.Status = TenderCreated
	case TenderPublishedStr:
		t.Status = TenderPublished
	case TenderClosedStr:
		t.Status = TenderClosed
	default:
		return errors.New("invalid value")
	}
	return nil
}

func IsServiceTypeCorrect(serviceType string) bool {
	CorrectServiceTypes := map[string]bool{"Construction": true, "Delivery": true, "Manufacture": true}
	return CorrectServiceTypes[serviceType]
}

const (
	BidCreated int = iota
	BidPublished
	BidCanceled
)

const (
	BidCreatedStr   = "Created"
	BidPublishedStr = "Published"
	BidCanceledStr  = "Canceled"
)

// Bid Модель для заявки на тендер
type Bid struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;unique"`
	Name        string    `gorm:"type:varchar(50);default:null"`
	Description string    `gorm:"type:text;default:null"`
	Status      int       `gorm:"not null"`
	AuthorType  string    `gorm:"type:varchar(255);default:null"`
	Version     int       `gorm:"not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	VotesNumber int       `gorm:"not null;default:0"`

	EmployeeID     uuid.UUID     `gorm:"type:uuid;default:null"`    // Связь с таблицей Employee через ID
	Employee       *Employee     `gorm:"foreignKey:EmployeeID"`     // Связь с автором через foreign key
	OrganizationID uuid.UUID     `gorm:"type:uuid;not null"`        // Связь с таблицей Organization через ID
	Organization   *Organization `gorm:"foreignKey:OrganizationID"` // Связь с организацией через foreign key
	TenderID       uuid.UUID     `gorm:"type:uuid;not null"`        // Связь с таблицей tender через ID
	Tender         *Tender       `gorm:"foreignKey:TenderID"`       // Связь с тендером через foreign key
}

func IsDecisionCorrect(serviceType string) bool {
	CorrectServiceTypes := map[string]bool{"Approved": true, "Rejected": true}
	return CorrectServiceTypes[serviceType]
}

func (b *Bid) StatusToString() string {
	switch b.Status {
	case BidCreated:
		return BidCreatedStr
	case BidPublished:
		return BidPublishedStr
	case BidCanceled:
		return BidCanceledStr
	}
	return ""
}

func (b *Bid) StringToStatus(status string) error {
	switch status {
	case BidCreatedStr:
		b.Status = BidCreated
	case BidPublishedStr:
		b.Status = BidPublished
	case BidCanceledStr:
		b.Status = BidCanceled
	default:
		return errors.New("invalid value")
	}
	return nil
}
