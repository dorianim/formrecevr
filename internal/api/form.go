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
		case "multipart/form-data":
			err = c.Request.ParseMultipartForm(6400000)
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

		atLeastOneSuccess := false
		for _, targetConfig := range formConfig.Targets {
			if !targetConfig.Enabled {
				continue
			}

			templateData := make(map[string]interface{})
			for k, v := range c.Request.Form {
				if len(v) == 1 {
					templateData[k] = v[0]
				} else {
					templateData[k] = v
				}
			}
			templateData["params"] = targetConfig.Params

			targetData, err := template.ExecuteTemplateFromFile(targetConfig.Template, templateData)
			if err != nil {
				log.Printf("Error processing body template %s: %v", targetConfig.Template, err)
				continue
			}

			shoutrrrURL, err := template.ExecuteTemplateFromString(targetConfig.ShoutrrrURL, templateData)
			if err != nil {
				log.Printf("Error processing URL template %s: %v", targetConfig.ShoutrrrURL, err)
				continue
			}

			err = shoutrrr.Send(shoutrrrURL, targetData)
			if err != nil {
				log.Printf("Error sending form to %s: %v", shoutrrrURL, err)
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
	for _, form := range config.Forms {
		if form.Id == formID && form.Enabled {
			return form
		}
	}
	return nil
}
