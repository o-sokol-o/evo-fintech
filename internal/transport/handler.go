package transport

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/o-sokol-o/evo-fintech/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/o-sokol-o/evo-fintech/internal/domain"
	v1 "github.com/o-sokol-o/evo-fintech/internal/transport/v1"
)

type IServicesEVO interface {
	GetFilteredData(ctx context.Context, input domain.FilterSearchInput) ([]domain.Transaction, error)
	FetchExternTransactions(ctx context.Context, url string) (domain.Status, error)
}
type IServicesRemote interface {
	Get(ctx context.Context, from, to *int) ([]domain.Transaction, error)
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

func (h *Handler) Init(cfg *domain.Config) *gin.Engine {
	router := gin.Default()

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI_v1(router)

	return router
}

func (h *Handler) initAPI_v1(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.servicesEVO, h.servicesRemote)
	api := router.Group("/api/v1")
	{
		handlerV1.Init(api)
	}
}
