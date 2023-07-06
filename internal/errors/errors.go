package errors

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	Type    string `json:"type"`
}

type ErrorType string

const (
	Provider      ErrorType = "ProviderError"
	Configuration ErrorType = "ConfigurationError"
)

func HandleError(c echo.Context, err error, title string, message string, status int, errType string) {
	c.Logger().Error(err)

	res := &ErrorResponse{
		Message: message,
		Title:   title,
		Status:  http.StatusText(status),
		Type:    errType,
	}

	if err := c.JSON(status, res); err != nil {
		c.Logger().Error(err)
	}
}
