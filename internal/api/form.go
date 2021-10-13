package api

import (
	"log"
	"net/http"

	"github.com/containrrr/shoutrrr"
	"github.com/dorianim/formrecevr/internal/config"
	"github.com/gin-gonic/gin"
)

// POST /api/v1/accounts
func PostForm(router *gin.RouterGroup) {
	router.POST("/:formID", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")

		formConfig := getForm(c.Param("formID"))
		if formConfig == nil {
			c.JSON(http.StatusNotFound, map[string]string{"message": "Form not found"})
			return
		}

		atLeastOneSuccess := false
		for _, target := range formConfig.Targets {
			err := shoutrrr.Send(target, "Hello world (or slack channel) !")
			if err != nil {
				log.Printf("Error sending form to %s: %v", target, err)
			} else {
				atLeastOneSuccess = true
			}
		}

		if atLeastOneSuccess {
			c.JSON(http.StatusOK, map[string]string{"message": "success"})
		} else {
			c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
		}
	})
}

func getForm(formID string) *config.FormConfig {
	config := config.GetConfig()
	if val, ex := config.Forms[formID]; ex {
		return &val
	} else {
		return nil
	}
}
