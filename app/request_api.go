package app

import (
	"ddplanet-server/app/repositories/support"
	"ddplanet-server/app/supporting"
	"ddplanet-server/pkg/httpext"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

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

// GetSupportRequests godoc
// @Summary Get list of support requests
// @Description Get list of support requests according to filter
// @Param limit query int false "limit of requests"
// @Param status query []string false "status to filer in" collectionFormat(multi)
// @Param tsCreatedFrom query int false "unix timestamp time of creation"
// @Param tsCreatedTo query int false "unix timestamp time of creation"
// @Param search query string false "search by phone or ID"
// @Produce json
// @Success 200 {object} []supporting.SupportRequest
// @Success 400 {object} httpext.ErrorResponse
// @Failure 500 {object} httpext.ErrorResponse
// @Router /support/requests [get]
func GetSupportRequests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := int64(100)
	tsCreatedFrom := int64(0)
	tsCreatedTo := int64(0)
	availableStatuses := make([]string, 0)

	qvalues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		logrus.Warningf("can't parse query in GetSupportRequests: %s", err)
	}

	limitStr := qvalues.Get("limit")
	if limitStr != "" {
		var err error

		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			logrus.Warningf("can't parse limit query param: %s", err)
		}
	}

	tsCreatedFromStr := qvalues.Get("tsCreatedFrom")
	if tsCreatedFromStr != "" {
		var err error
		tsCreatedFrom, err = strconv.ParseInt(tsCreatedFromStr, 10, 64)
		if err != nil {
			logrus.Warningf("can't parse tsCreatedFrom query param: %s", err)
		}
	}

	tsCreatedToStr := qvalues.Get("tsCreatedTo")
	if tsCreatedToStr != "" {
		var err error
		tsCreatedTo, err = strconv.ParseInt(tsCreatedToStr, 10, 64)
		if err != nil {
			logrus.Warningf("can't parse tsCreatedTo query param: %s", err)
		}
	}

	if availableStatusesFromQuery, exist := qvalues["status"]; exist {
		availableStatuses = availableStatusesFromQuery
	}

	availableStatusesForFilter := make([]supporting.SupportRequestStatus, 0)
	for _, s := range availableStatuses {
		vs := supporting.TryToConvertToValidSupportRequestStatus(s)
		if vs != "" {
			availableStatusesForFilter = append(availableStatusesForFilter, vs)
		}
	}

	search := qvalues.Get("search")

	requests, err := support.GetRequests(
		ctx,
		availableStatusesForFilter,
		tsCreatedFrom,
		tsCreatedTo,
		limit,
		search,
	)
	if err != nil {
		logrus.Errorf("failed to get requests from database: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}, http.StatusInternalServerError)
		return
	}

	httpext.JSON(w, requests)
}
