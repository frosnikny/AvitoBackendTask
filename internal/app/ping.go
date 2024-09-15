package app

import "github.com/gin-gonic/gin"

func (a *Application) GetPing(c *gin.Context) {
	c.String(200, "ok")
}
