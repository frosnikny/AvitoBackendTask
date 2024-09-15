package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"project/internal/models"
	"slices"
)

func (a *Application) respondError(c *gin.Context, code int, message string) {
	if message == "" {
		if code == http.StatusInternalServerError {
			message = "Непредвиденная ошибка со стороны сервера"
		}
		if code == http.StatusBadRequest {
			message = "Неверный формат запроса или его параметров"
		}
		if code == http.StatusUnauthorized {
			message = "Пользователь не существует или некорректен"
		}
		if code == http.StatusForbidden {
			message = "Недостаточно прав для выполнения действия"
		}
	}
	c.JSON(code, gin.H{"reason": message})
}

func (a *Application) getOrganization(c *gin.Context, organizationID uuid.UUID) (*models.Organization, bool) {
	organization, err := a.repo.GetOrganizationByID(organizationID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return nil, false
	}
	if organization == nil {
		a.respondError(c, http.StatusNotFound, "Данной организации не существует")
		return nil, false
	}
	return organization, true
}

func (a *Application) getUser(c *gin.Context, username string) (*models.Employee, bool) {
	employee, err := a.repo.GetEmployeeByUsername(username)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return nil, false
	}
	if employee == nil {
		a.respondError(c, http.StatusUnauthorized, "")
		return nil, false
	}
	return employee, true
}

func (a *Application) checkUserResponsibleOrganization(c *gin.Context, employee *models.Employee, organization *models.Organization) bool {
	employeeOrganizations, err := a.repo.GetEmployeeOrganizations(employee.ID)
	if err != nil {
		a.respondError(c, http.StatusInternalServerError, "")
		return false
	}
	foundedID := slices.IndexFunc(employeeOrganizations, func(org models.Organization) bool { return org.ID == organization.ID })
	if foundedID == -1 {
		a.respondError(c, http.StatusForbidden, "Данный пользователь не соответствует указанной организации")
		return false
	}
	return true
}
