package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"project/internal/models"
	"project/internal/schemes"
	"slices"
	"strconv"
	"time"
)

func (a *Application) PostNewBid(c *gin.Context) {
	// "Предложения могут создавать пользователи от имени своей организации."
	var err error
	var request schemes.PostNewBidRequest
	if err = c.ShouldBindBodyWithJSON(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}

	tender, err := a.repo.GetTenderByID(request.TenderID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	if tender == nil {
		a.respondError(c, http.StatusNotFound, "Данного тендера не существует")
		return
	}

	var organization *models.Organization
	var employee *models.Employee
	switch request.AuthorType {
	case "Organization":
		{
			a.respondError(c, http.StatusBadRequest, "Организации запрещено создавать предложение без ответственного лица")
			return
		}
	case "User":
		{
			employee, err = a.repo.GetEmployeeByID(request.AuthorID)
			if err != nil {
				a.respondError(c, http.StatusInternalServerError, "")
				return
			}
			if employee == nil {
				a.respondError(c, http.StatusUnauthorized, "")
				return
			}

			var isOk bool
			organization, isOk = a.getOrganization(c, tender.OrganizationID)
			if !isOk {
				return
			}
			isOk = a.checkUserResponsibleOrganization(c, employee, organization)
			if !isOk {
				return
			}
		}
	default:
		{
			a.respondError(c, http.StatusBadRequest, "")
			return
		}
	}

	var employeeID uuid.UUID
	if employee != nil {
		employeeID = employee.ID
	}
	bid := models.Bid{
		Name:           request.Name,
		Description:    request.Description,
		Status:         models.BidCreated,
		AuthorType:     request.AuthorType,
		Version:        1,
		CreatedAt:      time.Now(),
		EmployeeID:     employeeID,
		Employee:       employee,
		OrganizationID: organization.ID,
		Organization:   organization,
		TenderID:       tender.ID,
		Tender:         tender,
	}

	log.Println(bid)
	err = a.repo.AddBid(&bid)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}

	response := schemes.ToDefaultBidResponse(&bid)
	c.JSON(http.StatusOK, response)
}

// GetMyBids ищу только bids, где пользователь является создателем, потому что так написано в yml
func (a *Application) GetMyBids(c *gin.Context) {
	var request schemes.GetMyBidsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if request.Limit == 0 {
		request.Limit = 5
	}

	employee, isOk := a.getUser(c, request.Username)
	if !isOk {
		return
	}
	bids, err := a.repo.GetUserBids(request.Limit, request.Offset, employee.ID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}

	response := schemes.ArrayToDefaultBidResponses(bids)
	c.JSON(http.StatusOK, response)
}

func (a *Application) GetBidsList(c *gin.Context) {
	var request schemes.GetBidsListRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if err := c.ShouldBindQuery(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	tenderUUID, err := uuid.Parse(request.URI.TenderID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	if request.Limit == 0 {
		request.Limit = 5
	}

	bids, err := a.repo.GetBidsByTender(request.Limit, request.Offset, tenderUUID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	log.Println(bids)
	if len(bids) == 0 {
		a.respondError(c, http.StatusNotFound, "Тендер или предложение не найдено")
		return
	}

	employee, isOk := a.getUser(c, request.Username)
	if !isOk {
		return
	}
	// Получаем все организации данного юзера
	organizations, err := a.repo.GetEmployeeOrganizations(employee.ID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	if organizations == nil {
		a.respondError(c, http.StatusForbidden, "")
		return
	}

	// Проверяем, подходят ли организации
	var result []models.Bid
	for _, bid := range bids {
		foundedID := slices.IndexFunc(organizations, func(org models.Organization) bool { return org.ID == bid.OrganizationID })
		if foundedID != -1 {
			result = append(result, bid)
		}
	}

	response := schemes.ArrayToDefaultBidResponses(bids)
	c.JSON(http.StatusOK, response)
}

func (a *Application) checkAuthBid(c *gin.Context, employee *models.Employee, bid *models.Bid) bool {
	if bid.EmployeeID == employee.ID {
		return true
	}

	organizations, err := a.repo.GetEmployeeOrganizations(employee.ID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return false
	}
	foundedID := slices.IndexFunc(organizations, func(org models.Organization) bool { return org.ID == bid.OrganizationID })
	if foundedID == -1 {
		a.respondError(c, http.StatusForbidden, "")
		return false
	}

	return true
}

func (a *Application) checkAndGetBid(c *gin.Context, bidUUID uuid.UUID, username string) (*models.Employee, *models.Bid, bool) {
	employee, isOk := a.getUser(c, username)
	if !isOk {
		return nil, nil, false
	}
	bid, err := a.repo.GetBidByID(bidUUID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return nil, nil, false
	}
	if bid == nil {
		a.respondError(c, http.StatusNotFound, "Предложение не найдено")
		return nil, nil, false
	}

	isOk = a.checkAuthBid(c, employee, bid)
	if !isOk {
		return nil, nil, false
	}

	return employee, bid, true
}

func (a *Application) GetBidStatus(c *gin.Context) {
	var request schemes.GetBidStatusRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if err := c.ShouldBindQuery(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	bidUUID, err := uuid.Parse(request.URI.BidID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}

	_, bid, isOk := a.checkAndGetBid(c, bidUUID, request.Username)
	if !isOk {
		return
	}

	c.String(http.StatusOK, bid.StatusToString())
}

func (a *Application) PutBidStatus(c *gin.Context) {
	var request schemes.PutBidStatusRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if err := c.ShouldBindQuery(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	bidUUID, err := uuid.Parse(request.URI.BidID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}

	_, bid, isOk := a.checkAndGetBid(c, bidUUID, request.Username)
	if !isOk {
		return
	}

	// Начало изменения bid
	bidVersion := a.createBidVersion(bid)

	err = bid.StringToStatus(request.Status)
	if err != nil {
		a.respondError(c, http.StatusBadRequest, "Неверное значение параметра status")
		return
	}

	bid.Version++
	err = a.repo.AddBidVersion(&bidVersion)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	// Конец сохранения версии bid

	err = a.repo.SaveBid(bid)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	// Конец изменения bid

	response := schemes.ToDefaultBidResponse(bid)
	c.JSON(http.StatusOK, response)
}

func (a *Application) PatchBid(c *gin.Context) {
	var request schemes.PatchBidRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if err := c.ShouldBindQuery(&request.Query); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if err := c.ShouldBindBodyWithJSON(&request.Body); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	bidUUID, err := uuid.Parse(request.URI.BidID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}

	_, bid, isOk := a.checkAndGetBid(c, bidUUID, request.Query.Username)
	if !isOk {
		return
	}

	// Начало изменения bid
	bidVersion := a.createBidVersion(bid)

	if request.Body.Name != "" {
		bid.Name = request.Body.Name
	}
	if request.Body.Description != "" {
		bid.Description = request.Body.Description
	}

	bid.Version++
	err = a.repo.AddBidVersion(&bidVersion)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	// Конец сохранения версии bid

	if err = a.repo.SaveBid(bid); err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	// Конец изменения bid

	response := schemes.ToDefaultBidResponse(bid)
	c.JSON(http.StatusOK, response)
}

func (a *Application) PutBitSubmitDecision(c *gin.Context) {
	var request schemes.PutBidSubmitDecisionRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if err := c.ShouldBindQuery(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	bidUUID, err := uuid.Parse(request.URI.BidID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	if !models.IsDecisionCorrect(request.Decision) {
		a.respondError(c, http.StatusBadRequest, "Неправильное значение decision")
		return
	}

	_, bid, isOk := a.checkAndGetBid(c, bidUUID, request.Username)
	if !isOk {
		return
	}
	if bid.Status == models.BidCanceled {
		a.respondError(c, http.StatusMethodNotAllowed, "Запрещено отправлять решения на отклоненное предложение")
		return
	}

	// Начало изменения bid
	bidVersion := a.createBidVersion(bid)

	if request.Decision == "Approved" {
		bid.VotesNumber++

		organizationResponsibleCount, err2 := a.repo.CountOrganizationEmployees(bid.OrganizationID)
		if err2 != nil {
			a.respondError(c, http.StatusInternalServerError, "")
			return
		}
		quorum := min(3, organizationResponsibleCount)

		if int64(bid.VotesNumber) >= quorum {
			bid.Status = models.BidPublished
			bid.VotesNumber = 0

			tender, err3 := a.repo.GetTenderByID(bid.TenderID)
			if err3 != nil {
				a.respondError(c, http.StatusInternalServerError, "")
				return
			}
			if tender == nil {
				a.respondError(c, http.StatusNotFound, "Тендер не найден")
				return
			}
			tender.Status = models.TenderClosed
			err3 = a.repo.SaveTender(tender)
			if err3 != nil {
				a.respondError(c, http.StatusInternalServerError, "")
				return
			}
		}
	} else {
		bid.Status = models.BidCanceled
		bid.VotesNumber = 0
	}

	bid.Version++
	err = a.repo.AddBidVersion(&bidVersion)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	// Конец сохранения версии bid

	err = a.repo.SaveBid(bid)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	// Конец изменения bid

	response := schemes.ToDefaultBidResponse(bid)
	c.JSON(http.StatusOK, response)
}

func (a *Application) PutBitFeedback(c *gin.Context) {
	var request schemes.PutBidFeedbackRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if err := c.ShouldBindQuery(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	bidUUID, err := uuid.Parse(request.URI.BidID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}

	employee, bid, isOk := a.checkAndGetBid(c, bidUUID, request.Username)
	if !isOk {
		return
	}

	feedback := models.Feedback{
		Description: request.BidFeedback,
		CreatedAt:   time.Now(),
		BidID:       bid.ID,
		Bid:         bid,
		AuthorID:    employee.ID,
		Author:      employee,
	}

	err = a.repo.AddFeedback(&feedback)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}

	response := schemes.ToDefaultBidResponse(bid)
	c.JSON(http.StatusOK, response)
}

func (a *Application) GetBidReviews(c *gin.Context) {
	var request schemes.GetBidReviewsRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if err := c.ShouldBindQuery(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	tenderUUID, err := uuid.Parse(request.URI.TenderID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	if request.Limit == 0 {
		request.Limit = 5
	}

	requester, isOk := a.getUser(c, request.RequesterUsername)
	if !isOk {
		return
	}
	tender, err := a.repo.GetTenderByID(tenderUUID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	if tender == nil {
		a.respondError(c, http.StatusNotFound, "Тендер не найден")
		return
	}

	isOk = a.checkTenderAccessByUser(c, requester, tender)
	if !isOk {
		return
	}
	author, isOk := a.getUser(c, request.AuthorUsername)
	if !isOk {
		return
	}

	bids, err := a.repo.GetUserTenderBids(author.ID, tenderUUID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}

	var reviews []models.Feedback
	for _, bid := range bids {
		foundReviews, err := a.repo.GetFeedbacksByBid(bid.ID)
		if err != nil {
			a.respondError(c, http.StatusInternalServerError, "")
			return
		}
		reviews = append(reviews, foundReviews...)
	}

	maxLen := min(request.Offset+request.Limit, len(reviews))
	resultReviews := make([]models.Feedback, maxLen)
	for i := request.Offset; i < maxLen; i++ {
		resultReviews[i-request.Offset] = reviews[i]
	}

	response := schemes.ArrayToBidFeedbackResponses(resultReviews)
	c.JSON(http.StatusOK, response)
}

func (a *Application) PutBidRollback(c *gin.Context) {
	var request schemes.PutBidRollbackRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	bidUUID, err := uuid.Parse(request.URI.BidID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	bidVersion, err := strconv.Atoi(request.URI.Version)
	if err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}

	_, bid, isOk := a.checkAndGetBid(c, bidUUID, request.Username)
	if !isOk {
		return
	}
	if bid.Status == models.BidCanceled {
		a.respondError(c, http.StatusMethodNotAllowed, "Запрещено отправлять решения на отклоненное предложение")
		return
	}

	// Начало изменения bid
	bidCreatedVersion := a.createBidVersion(bid)

	isOk = a.rollbackBidToVersion(c, bid, bidVersion)
	if !isOk {
		return
	}

	bid.Version++
	err = a.repo.AddBidVersion(&bidCreatedVersion)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	// Конец сохранения версии bid

	err = a.repo.SaveBid(bid)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	// Конец изменения bid

	response := schemes.ToDefaultBidResponse(bid)
	c.JSON(http.StatusOK, response)
}
