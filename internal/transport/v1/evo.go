package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"

	"github.com/o-sokol-o/evo-fintech/internal/domain"
)

//go:generate mockgen -source=evo.go -destination=mocks/mock.go

type IServicesEVO interface {
	GetFilteredData(ctx context.Context, input domain.FilterSearchInput) ([]domain.Transaction, error)
	FetchExternTransactions(ctx context.Context, url string) (domain.Status, error)
}
type IServicesRemote interface {
	Get(ctx context.Context) ([]domain.Transaction, error)
}

type Handler struct {
	servicesEVO    IServicesEVO
	servicesRemote IServicesRemote
}

func NewHandler(servicesEVO IServicesEVO, servicesRemote IServicesRemote) *Handler {
	return &Handler{
		servicesEVO:    servicesEVO,
		servicesRemote: servicesRemote,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	api.GET("/get_csv_mock_remote_service", h.getSourceFileCSV_as_MockRemoteService)

	api.POST("/download_remote_transactions", h.downloadRemoteTransactionsCSV)

	filtered := api.Group("/filtered")
	{
		filtered.POST("/csv", h.getFilteredFileCSV)
		filtered.POST("/json", h.getFilteredDataJSON)
	}
}

// func getHeadersCSV(myStruct domain.Transaction) string {
// 	var header string
// 	e := reflect.ValueOf(&myStruct).Elem()
// 	if e.NumField() > 1 {
// 		header = e.Type().Field(1).Name
// 	}
// 	for i := 2; i < e.NumField(); i++ {
// 		header = header + "," + e.Type().Field(i).Name
// 	}
// 	return header
// }

func buildCSV(transactions []domain.Transaction) *strings.Reader {

	str, _ := gocsv.MarshalString(&transactions)
	return strings.NewReader(str)

	/*
		var builder strings.Builder

		// get the column names first
		builder.WriteString(getHeadersCSV(transactions[0]) + "\n")

		// for _, trns := range transactions {
		// 	w := fmt.Sprintf("%s,%d\n", trns.Service, trns.TransactionId)
		// 	builder.WriteString(w)
		// }

		enc := struct2csv.New()
		// var rows [][]string

		// get the column names first
		// colhdrs, err := enc.GetColNames(transactions[0])
		// if err != nil {
		// 	// handle error
		// }
		// builder.WriteString(strings.Join(colhdrs, ",") + "\n")

		// get the data from each struct
		for _, v := range transactions {
			row, _ := enc.GetRow(v)
			// if err != nil {
			// 	// handle error
			// }
			// rows = append(rows, row)

			builder.WriteString(strings.Join(row, ",") + "\n")

		}

		return strings.NewReader(builder.String())
	*/
}

// @Summary Test service: Gives a CSV file with initial transactions
// @Tags Mock remote service
// @ID getSourceFileCSV_as_MockRemoteService-csv
// @Success 200
// @Success 204
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /api/v1/get_csv_mock_remote_service/ [get]
func (h *Handler) getSourceFileCSV_as_MockRemoteService(ctx *gin.Context) {
	transactions, err := h.servicesRemote.Get(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	if len(transactions) == 0 {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("no data"))
		return
	}

	CSV := buildCSV(transactions)

	headers := map[string]string{
		"Content-Disposition": `attachment; filename="source.csv"`,
	}

	ctx.DataFromReader(http.StatusOK, -1, "text/html; charset=UTF-8", CSV, headers)
}

// @Summary Request filtered csv file
// @Tags Services
// @ID get-filtered-csv
// @Param   input body domain.FilterSearchInput true " "
// @Success 200
// @Success 204
// @Failure 400   {object} domain.ErrorResponse
// @Failure 500   {object} domain.ErrorResponse
// @Router /api/v1/filtered/csv/ [post]
func (h *Handler) getFilteredFileCSV(ctx *gin.Context) {
	var input domain.FilterSearchInput

	if err := ctx.BindJSON(&input); err != nil {
		logrus.Error(err)
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	v := validator.New()
	if err := v.Struct(input); err != nil {
		logrus.Error(err)
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("data validation error"))
		return
	}

	transactions, err := h.servicesEVO.GetFilteredData(context.Background(), input)
	if err != nil {
		logrus.Error(err)
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("service failure"))
		return
	}

	if len(transactions) == 0 {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	CSV := buildCSV(transactions)

	headers := map[string]string{
		"Content-Disposition": `attachment; filename="source.csv"`,
	}

	ctx.DataFromReader(http.StatusOK, -1, "text/html; charset=UTF-8", CSV, headers)
}

// @Summary Request filtered json
// @Tags Services
// @ID get-filtered-json
// @Param   input body     domain.FilterSearchInput true " "
// @Success 200   {object} []domain.Transaction
// @Success 204
// @Failure 400   {object} domain.ErrorResponse
// @Failure 500   {object} domain.ErrorResponse
// @Router /api/v1/filtered/json/ [post]
func (h *Handler) getFilteredDataJSON(ctx *gin.Context) {
	var input domain.FilterSearchInput

	if err := ctx.BindJSON(&input); err != nil {
		logrus.Error(err)
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	v := validator.New()
	if err := v.Struct(input); err != nil {
		logrus.Error(err)
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("data validation error"))
		return
	}

	transactions, err := h.servicesEVO.GetFilteredData(context.Background(), input)
	if err != nil {
		logrus.Error(err)
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("service failure"))
		return
	}

	if len(transactions) == 0 {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

// Summary Request to download remote transactions: the request is executed fake 10 seconds, at other times it gives a status.

// @Summary Request to download remote transactions
// @Tags Services
// @ID request-download-remote-transactions
// @Param   input body     domain.UrlInput true " "
// @Success 200   {object} domain.StatusResponse
// @Success 204
// @Failure 400   {object} domain.ErrorResponse
// @Failure 500   {object} domain.ErrorResponse
// @Router /api/v1/download_remote_transactions/ [post]
func (h *Handler) downloadRemoteTransactionsCSV(ctx *gin.Context) {
	var input domain.UrlInput

	if err := ctx.BindJSON(&input); err != nil {
		logrus.Error(err)
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	status, err := h.servicesEVO.FetchExternTransactions(context.Background(), *input.URL)
	if err != nil {
		logrus.Error(err)
		newErrorResponse(ctx, http.StatusInternalServerError, errors.New("service failure"))
		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: status})
}
