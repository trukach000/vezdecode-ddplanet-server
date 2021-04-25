package app

import (
	"ddplanet-server/app/repositories/support"
	"ddplanet-server/pkg/httpext"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type CreateSupportRequestInput struct {
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	LastName   string `json:"lastName"`
	Phone      string `json:"phone"`
	Message    string `json:"message"`
}

type CreateSupportRequestOutput struct {
	RequestID int64 `json:"requestId"`
}

// CreateSupportRequest godoc
// @Summary Create new request to support
// @Description Create new request in DB for further checking
// @Param message body CreateSupportRequestInput true "request body"
// @Accept json
// @Produce json
// @Success 200 {object} CreateSupportRequestOutput
// @Success 400 {object} httpext.ErrorResponse
// @Failure 500 {object} httpext.ErrorResponse
// @Router /support/request [post]
func CreateSupportRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input CreateSupportRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logrus.Errorf("failed to decode payload: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to decode payload",
		}, http.StatusBadRequest)
		return
	}

	reqID, err := support.CreateSupportTask(ctx, input.FirstName, input.SecondName, input.LastName, input.Phone, input.Message)
	if err != nil {
		logrus.Errorf("failed to insert new request to DB: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}, http.StatusInternalServerError)
		return
	}

	httpext.JSON(w, CreateSupportRequestOutput{RequestID: reqID})
}
