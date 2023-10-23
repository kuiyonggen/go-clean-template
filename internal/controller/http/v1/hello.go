package v1

import (
        "fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kuiyonggen/go-clean-template/pkg/logger"
        "github.com/kuiyonggen/go-clean-template/config"
)

type helloRoutes struct {
	l logger.Interface
}

func newHelloRoutes(handler *gin.RouterGroup, cfg *config.Config) {
        r := &helloRoutes{cfg.Logger}

	h := handler.Group("/hello")
	{
		h.GET("/say", r.say)
		h.POST("/greeting", r.greeting)
	}
}

// @Summary      Show hello
// @Description  Show hello
// @ID           echo
// @Tags               hello
// @Accept       json
// @Produce      json
// @Success      200  {object}  response
// @Failure      500  {object}  response
// @Router       /hello/say [get]
func (r *helloRoutes) say(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello, Alice@!@")
}

type greetingRequest struct {
	Name      string `json:"name"       binding:"required"  example:"alice"`
}

// @Summary      Greeting
// @Description  Greeting
// @ID           greeting
// @Tags               hello
// @Accept       json
// @Produce      json
// @Param        request  body      greetingRequest  true  "greeting"
// @Success      200      {object}  response
// @Failure      400      {object}  response
// @Failure      500      {object}  response
// @Router       /hello/greeting [post]
func (r *helloRoutes) greeting(c *gin.Context) {
	var request greetingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - greeting")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Hello, %s@!@", request.Name))
}
