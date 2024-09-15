package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/internal/models"
)

func (a *Application) createTenderVersion(tender *models.Tender) models.TenderVersion {
	tenderVersion := models.TenderVersion{
		Version:     tender.Version,
		Status:      tender.Status,
		Name:        tender.Name,
		Description: tender.Description,
		ServiceType: tender.ServiceType,
		TenderID:    tender.ID,
		Tender:      tender,
	}

	return tenderVersion
}

func (a *Application) rollbackTenderToVersion(c *gin.Context, tender *models.Tender, version int) bool {
	tenderVersion, err := a.repo.GetTenderVersionByID(tender.ID, version)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return false
	}
	if tenderVersion == nil {
		a.respondError(c, http.StatusNotFound, "Версия тендера не найдена")
		return false
	}

	tender.Name = tenderVersion.Name
	tender.ServiceType = tenderVersion.ServiceType
	tender.Description = tenderVersion.Description
	tender.Status = tenderVersion.Status

	return true
}

func (a *Application) createBidVersion(bid *models.Bid) models.BidVersion {
	bidVersion := models.BidVersion{
		Version:     bid.Version,
		Status:      bid.Status,
		Name:        bid.Name,
		Description: bid.Description,
		VotesNumber: bid.VotesNumber,
		BidID:       bid.ID,
		Bid:         bid,
	}

	return bidVersion
}

func (a *Application) rollbackBidToVersion(c *gin.Context, bid *models.Bid, version int) bool {
	bidVersion, err := a.repo.GetBidVersionByID(bid.ID, version)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return false
	}
	if bidVersion == nil {
		a.respondError(c, http.StatusNotFound, "Версия предложения не найдена")
		return false
	}

	bid.Name = bidVersion.Name
	bid.Description = bidVersion.Description
	bid.Status = bidVersion.Status
	bid.VotesNumber = bidVersion.VotesNumber

	return true
}
