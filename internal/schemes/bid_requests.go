package schemes

import "github.com/google/uuid"

type PostNewBidRequest struct {
	Name        string    `json:"name" binding:"required,max=100"`
	Description string    `json:"description" binding:"required,max=500"`
	TenderID    uuid.UUID `json:"tenderId" binding:"required,uuid,max=100"`
	AuthorType  string    `json:"authorType" binding:"required"`
	AuthorID    uuid.UUID `json:"authorId" binding:"required,uuid,max=100"`
}

type GetMyBidsRequest struct {
	Limit    int    `form:"limit" binding:"omitempty,min=1"`
	Offset   int    `form:"offset" binding:"omitempty,min=0"`
	Username string `form:"username" binding:"required"`
}

type GetBidsListRequest struct {
	URI struct {
		// На самом деле это тендер, но в роутах он указан как bidId, чтобы убрать ошибку из-за одинаковых путей
		TenderID string `uri:"bidId" binding:"required,uuid"`
	}
	Username string `form:"username" binding:"required"`
	Limit    int    `form:"limit" binding:"omitempty,min=1"`
	Offset   int    `form:"offset" binding:"omitempty,min=0"`
}

type GetBidStatusRequest struct {
	URI struct {
		BidID string `uri:"bidId" binding:"required,uuid"`
	}
	Username string `form:"username" binding:"required"`
}

type PutBidStatusRequest struct {
	URI struct {
		BidID string `uri:"bidId" binding:"required,uuid"`
	}
	Status   string `form:"status" binding:"required"`
	Username string `form:"username" binding:"required"`
}

type PatchBidRequest struct {
	URI struct {
		BidID string `uri:"bidId" binding:"required,uuid"`
	}
	Query struct {
		Username string `form:"username" binding:"required"`
	}
	Body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
}

type PutBidSubmitDecisionRequest struct {
	URI struct {
		BidID string `uri:"bidId" binding:"required,uuid"`
	}
	Decision string `form:"decision" binding:"required"`
	Username string `form:"username" binding:"required"`
}

type PutBidFeedbackRequest struct {
	URI struct {
		BidID string `uri:"bidId" binding:"required,uuid"`
	}
	BidFeedback string `form:"bidFeedback" binding:"required"`
	Username    string `form:"username" binding:"required"`
}

type GetBidReviewsRequest struct {
	URI struct {
		TenderID string `uri:"bidId" binding:"required,uuid"`
	}
	AuthorUsername    string `form:"authorUsername" binding:"required"`
	RequesterUsername string `form:"requesterUsername" binding:"required"`
	Limit             int    `form:"limit" binding:"omitempty,min=1"`
	Offset            int    `form:"offset" binding:"omitempty,min=0"`
}

type PutBidRollbackRequest struct {
	URI struct {
		BidID   string `uri:"bidId" binding:"required,uuid"`
		Version string `uri:"version" binding:"required"`
	}
	Username string `form:"username" binding:"required"`
}
