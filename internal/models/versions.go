package models

import "github.com/google/uuid"

type TenderVersion struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;unique"`
	Version int       `gorm:"not null"`

	Status      int    `gorm:"not null"`
	Name        string `gorm:"type:varchar(50);default:null"`
	Description string `gorm:"type:text;default:null"`
	ServiceType string `gorm:"type:varchar(255);default:null"`

	TenderID uuid.UUID `gorm:"type:uuid;not null"`  // Связь с таблицей tender через ID
	Tender   *Tender   `gorm:"foreignKey:TenderID"` // Связь с тендером через foreign key
}

type BidVersion struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;unique"`
	Version int       `gorm:"not null"`

	Status      int    `gorm:"not null"`
	Name        string `gorm:"type:varchar(50);default:null"`
	Description string `gorm:"type:text;default:null"`
	VotesNumber int    `gorm:"not null;default:0"`

	BidID uuid.UUID `gorm:"type:uuid;not null"` // Связь с таблицей tender через ID
	Bid   *Bid      `gorm:"foreignKey:BidID"`   // Связь с тендером через foreign key
}
