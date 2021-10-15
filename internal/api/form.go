package api

import (
	"log"
	"net/http"

	"github.com/containrrr/shoutrrr"
	"github.com/dorianim/formrecevr/internal/config"
	"github.com/dorianim/formrecevr/internal/template"
	"github.com/gin-gonic/gin"
)

// ResponseBody is the body of a response
type ResponseBody struct {
	Message string
}

// PostForm registers the route POST /api/v1/accounts
func PostForm(router *gin.RouterGroup) {
	router.POST("/:formID", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")

		formConfig := getForm(c.Param("formID"))
		if formConfig == nil {
			c.JSON(http.StatusNotFound, ResponseBody{Message: "Form not found"})
			return
		}

		var err error

		switch c.ContentType() {
		case "application/x-www-form-urlencoded":
			err = c.Request.ParseForm()
			break
		case "multipart/form-data":
			err = c.Request.ParseMultipartForm(6400000)
			break
		default:
			log.Printf("Unsupported Content-Type: %s", c.ContentType())
			c.JSON(http.StatusBadRequest, ResponseBody{Message: "Unsupported Content-Type"})
			return
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, ResponseBody{Message: "Malformed form data"})
			log.Printf("Error parsing submited form: %v", err)
			return
		}

		log.Printf("Form data: %v", c.Request.Form)

		atLeastOneSuccess := false
		for _, targetConfig := range formConfig.Targets {
			if !targetConfig.Enabled {
				continue
			}
			log.Printf("Processing target %v", targetConfig)

			targetData, err := template.ExecuteTemplateFromFile(targetConfig.Template, c.Request.Form)
			if err != nil {
				log.Printf("Error processing template %s: %v", targetConfig.Template, err)
				continue
			}

			err = shoutrrr.Send(targetConfig.ShoutrrrURL, targetData)
			if err != nil {
				log.Printf("Error sending form to %s: %v", targetConfig.ShoutrrrURL, err)
				continue
			}

			atLeastOneSuccess = true
		}

		if !atLeastOneSuccess {
			c.JSON(http.StatusInternalServerError, ResponseBody{Message: "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, ResponseBody{Message: "Success"})
	})
}

func getForm(formID string) *config.FormConfig {
	config := config.GetConfig()
	val, ex := config.Forms[formID]
	if !ex {
		return nil
	}
	if !val.Enabled {
		return nil
	}
	return val
}
