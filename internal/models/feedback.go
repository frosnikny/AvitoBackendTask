package models

import (
	"github.com/google/uuid"
	"time"
)

// Feedback Модель для заявки на тендер
type Feedback struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;unique"`
	Description string    `gorm:"type:text;;not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`

	BidID    uuid.UUID `gorm:"type:uuid;not null"`  // Связь с таблицей Bid через ID
	Bid      *Bid      `gorm:"foreignKey:BidID"`    // Связь с предложением через foreign key
	AuthorID uuid.UUID `gorm:"type:uuid;not null"`  // Связь с таблицей Employee через ID
	Author   *Employee `gorm:"foreignKey:AuthorID"` // Связь с автором через foreign key
}
