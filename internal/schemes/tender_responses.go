package schemes

import (
	"github.com/google/uuid"
	"project/internal/models"
	"time"
)

type DefaultTenderResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ServiceType    string    `json:"serviceType"`
	Status         string    `json:"status"`
	OrganizationID uuid.UUID `json:"organizationId"`
	Version        int       `json:"version"`
	CreatedAt      time.Time `json:"createdAt"`
}

func ToDefaultTenderResponse(tender *models.Tender) DefaultTenderResponse {
	return DefaultTenderResponse{
		ID:             tender.ID,
		Name:           tender.Name,
		Description:    tender.Description,
		ServiceType:    tender.ServiceType,
		Status:         tender.StatusToString(),
		OrganizationID: tender.OrganizationID,
		Version:        tender.Version,
		CreatedAt:      tender.CreatedAt,
	}
}

func ArrayToDefaultTenderResponses(tenders []models.Tender) []DefaultTenderResponse {
	result := make([]DefaultTenderResponse, len(tenders))
	for i, tender := range tenders {
		result[i] = ToDefaultTenderResponse(&tender)
	}
	return result
}
