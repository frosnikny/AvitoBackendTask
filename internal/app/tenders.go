package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"project/internal/models"
	"project/internal/schemes"
	"strconv"
	"time"
)

func (a *Application) checkTenderAccessByUser(c *gin.Context, employee *models.Employee, tender *models.Tender) bool {
	organization, err := a.repo.GetOrganizationByID(tender.OrganizationID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return false
	}
	if organization == nil {
		a.respondError(c, http.StatusNotFound, "Данной организации не существует")
		return false
	}

	return a.checkUserResponsibleOrganization(c, employee, organization)
}

func (a *Application) GetTenders(c *gin.Context) {
	var request schemes.GetTendersRequest
	err := c.ShouldBind(&request)
	if err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if request.ServiceType != "" && !models.IsServiceTypeCorrect(request.ServiceType) {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if request.Limit == 0 {
		request.Limit = 5
	}

	tenders, err := a.repo.GetFilteredTenders(request.Limit, request.Offset, request.ServiceType)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}

	response := schemes.ArrayToDefaultTenderResponses(tenders)
	c.JSON(http.StatusOK, response)
}

func (a *Application) PostNewTender(c *gin.Context) {
	var request schemes.PostNewTenderRequest
	if err := c.ShouldBind(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if !models.IsServiceTypeCorrect(request.ServiceType) {
		a.respondError(c, http.StatusBadRequest, "Неправильно указан тип тендера")
		return
	}

	// Получаем организацию и автора тендера
	organization, isOk := a.getOrganization(c, request.OrganizationID)
	if !isOk {
		return
	}
	employee, err := a.repo.GetEmployeeByUsername(request.CreatorUsername)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	if employee == nil {
		a.respondError(c, http.StatusUnauthorized, "")
		return
	}

	// Проверяем, принадлежит ли пользователь организации
	if isOk := a.checkUserResponsibleOrganization(c, employee, organization); !isOk {
		return
	}

	tender := models.Tender{
		Name:           request.Name,
		Description:    request.Description,
		Status:         models.TenderCreated,
		ServiceType:    request.ServiceType,
		Version:        1,
		CreatedAt:      time.Now(),
		OrganizationID: request.OrganizationID,
		Organization:   organization,
		AuthorID:       employee.ID,
		Author:         employee,
	}

	err = a.repo.AddTender(&tender)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}

	response := schemes.ToDefaultTenderResponse(&tender)
	c.JSON(http.StatusOK, response)
}

func (a *Application) GetMyTenders(c *gin.Context) {
	var request schemes.GetMyTendersRequest
	if err := c.ShouldBind(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if request.Limit == 0 {
		request.Limit = 5
	}

	// Ищем пользователя с данным username
	employee, err := a.repo.GetEmployeeByUsername(request.Username)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	if employee == nil {
		a.respondError(c, http.StatusUnauthorized, "")
		return
	}

	// Выбираем тендеры данного пользователя
	tenders, err := a.repo.GetUserTenders(request.Limit, request.Offset, employee.ID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}

	response := schemes.ArrayToDefaultTenderResponses(tenders)
	c.JSON(http.StatusOK, response)
}

func (a *Application) GetTenderStatus(c *gin.Context) {
	var request schemes.GetTenderStatusRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	tenderUUID, err := uuid.Parse(request.URI.TenderID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
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

	if tender.Status != models.TenderPublished {
		if request.Username == "" {
			a.respondError(c, http.StatusForbidden, "")
			return
		}
		employee, isOk := a.getUser(c, request.Username)
		if !isOk {
			return
		}
		if isOk = a.checkTenderAccessByUser(c, employee, tender); !isOk {
			return
		}
	}

	c.String(http.StatusOK, tender.StatusToString())
}

func (a *Application) PutTenderStatus(c *gin.Context) {
	var request schemes.EditTenderStatusRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	tenderUUID, err := uuid.Parse(request.URI.TenderID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
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

	employee, isOk := a.getUser(c, request.Username)
	if !isOk {
		return
	}
	if isOk = a.checkTenderAccessByUser(c, employee, tender); !isOk {
		return
	}

	// Начало изменения tender
	tenderCreatedVersion := a.createTenderVersion(tender)

	err = tender.StringToStatus(request.Status)
	if err != nil {
		a.respondError(c, http.StatusBadRequest, "Неверное значение параметра status")
		return
	}

	tender.Version++
	err = a.repo.AddTenderVersion(&tenderCreatedVersion)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	// Конец сохранения версии tender

	err = a.repo.SaveTender(tender)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	// Конец изменения tender

	response := schemes.ToDefaultTenderResponse(tender)
	c.JSON(http.StatusOK, response)
}

func (a *Application) PatchTender(c *gin.Context) {
	var request schemes.EditTenderRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if err := c.ShouldBindQuery(&request.Query); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	tenderUUID, err := uuid.Parse(request.URI.TenderID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	if request.ServiceType != "" {
		if !models.IsServiceTypeCorrect(request.ServiceType) {
			a.respondError(c, http.StatusBadRequest, "Неправильно указан тип тендера")
			return
		}
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

	employee, isOk := a.getUser(c, request.Query.Username)
	if !isOk {
		return
	}
	if isOk = a.checkTenderAccessByUser(c, employee, tender); !isOk {
		return
	}

	// Начало изменения tender
	tenderCreatedVersion := a.createTenderVersion(tender)

	if request.Name != "" {
		tender.Name = request.Name
	}
	if request.ServiceType != "" {
		tender.ServiceType = request.ServiceType
	}
	if request.Description != "" {
		tender.Description = request.Description
	}

	tender.Version++
	err = a.repo.AddTenderVersion(&tenderCreatedVersion)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	// Конец сохранения версии tender

	if err = a.repo.SaveTender(tender); err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	// Конец изменения tender

	response := schemes.ToDefaultTenderResponse(tender)
	c.JSON(http.StatusOK, response)
}

func (a *Application) PutTenderRollback(c *gin.Context) {
	var request schemes.PutTenderRollbackRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		a.respondError(c, http.StatusBadRequest, "")
		return
	}
	tenderUUID, err := uuid.Parse(request.URI.TenderID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	tenderVersion, err := strconv.Atoi(request.URI.Version)
	if err != nil {
		a.respondError(c, http.StatusBadRequest, "")
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

	employee, isOk := a.getUser(c, request.Username)
	if !isOk {
		return
	}
	if isOk = a.checkTenderAccessByUser(c, employee, tender); !isOk {
		return
	}

	// Начало изменения tender
	tenderCreatedVersion := a.createTenderVersion(tender)

	isOk = a.rollbackTenderToVersion(c, tender, tenderVersion)
	if !isOk {
		return
	}

	tender.Version++
	err = a.repo.AddTenderVersion(&tenderCreatedVersion)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	// Конец сохранения версии tender

	err = a.repo.SaveTender(tender)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return
	}
	// Конец изменения tender

	response := schemes.ToDefaultTenderResponse(tender)
	c.JSON(http.StatusOK, response)
}
