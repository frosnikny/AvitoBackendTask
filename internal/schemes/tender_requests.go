package schemes

import "github.com/google/uuid"

type PostNewTenderRequest struct {
	Name            string    `json:"name" binding:"required,max=100"`
	Description     string    `json:"description"  binding:"required,max=500"`
	ServiceType     string    `json:"serviceType"  binding:"required"`
	OrganizationID  uuid.UUID `json:"organizationId" binding:"required,uuid,max=100"`
	CreatorUsername string    `json:"creatorUsername" binding:"required"`
}

type GetTendersRequest struct {
	Limit       int    `form:"limit" binding:"omitempty,min=1"`
	Offset      int    `form:"offset" binding:"omitempty,min=0"`
	ServiceType string `form:"service_type" binding:"omitempty"`
}

type GetMyTendersRequest struct {
	Limit    int    `form:"limit" binding:"omitempty,min=1"`
	Offset   int    `form:"offset" binding:"omitempty,min=0"`
	Username string `form:"username" binding:"required"`
}

type GetTenderStatusRequest struct {
	URI struct {
		// Не можем из URI взять сразу uid.UUID
		TenderID string `uri:"tenderId" binding:"required,uuid"`
	}
	Username string `form:"username" binding:"omitempty"`
}

// EditTenderStatusRequest Все-таки менять статус тендера может только ответственный
type EditTenderStatusRequest struct {
	URI struct {
		TenderID string `uri:"tenderId" binding:"required,uuid"`
	}
	Status   string `form:"status" binding:"required"`
	Username string `form:"username" binding:"required"`
}

type EditTenderRequest struct {
	URI struct {
		TenderID string `uri:"tenderId" binding:"required,uuid"`
	}
	Query struct {
		Username string `form:"username" binding:"required"`
	}
	Name        string `json:"name" binding:"omitempty,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	ServiceType string `json:"serviceType" binding:"omitempty"`
}

type PutTenderRollbackRequest struct {
	URI struct {
		TenderID string `uri:"tenderId" binding:"required,uuid"`
		Version  string `uri:"version" binding:"required"`
	}
	Username string `form:"username" binding:"required"`
}
