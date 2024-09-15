package schemes

import (
	"github.com/google/uuid"
	"project/internal/models"
	"time"
)

type DefaultBidResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AuthorType  string    `json:"authorType"`
	AuthorID    uuid.UUID `json:"authorId"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
}

func ToDefaultBidResponse(bid *models.Bid) DefaultBidResponse {
	var authorID uuid.UUID
	if bid.AuthorType == "Organization" {
		authorID = bid.OrganizationID
	} else {
		authorID = bid.EmployeeID
	}
	return DefaultBidResponse{
		ID:          bid.ID,
		Name:        bid.Name,
		Description: bid.Description,
		Status:      bid.StatusToString(),
		AuthorType:  bid.AuthorType,
		AuthorID:    authorID,
		Version:     bid.Version,
		CreatedAt:   bid.CreatedAt,
	}
}

func ArrayToDefaultBidResponses(bids []models.Bid) []DefaultBidResponse {
	result := make([]DefaultBidResponse, len(bids))
	for i, bid := range bids {
		result[i] = ToDefaultBidResponse(&bid)
	}
	return result
}

type BidFeedbackResponse struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

func ToBidFeedbackResponse(feedback *models.Feedback) BidFeedbackResponse {
	return BidFeedbackResponse{
		ID:          feedback.ID,
		Description: feedback.Description,
		CreatedAt:   feedback.CreatedAt,
	}
}

func ArrayToBidFeedbackResponses(feedbacks []models.Feedback) []BidFeedbackResponse {
	result := make([]BidFeedbackResponse, len(feedbacks))
	for i, feedback := range feedbacks {
		result[i] = ToBidFeedbackResponse(&feedback)
	}
	return result
}
