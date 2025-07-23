package router

import (
	"app_aggregator/internal/handlers"
	"app_aggregator/internal/repository"
	"github.com/gin-gonic/gin"
)

const apiV1 = "/api/v1"
const apiAdmin = apiV1 + "/admin"

type Builder struct {
	engine     *gin.Engine
	repository *repository.Repository
}

func NewBuilder(repo *repository.Repository) *Builder {
	return &Builder{
		engine:     gin.Default(),
		repository: repo,
	}
}

func (b *Builder) GetEngine() *gin.Engine {
	return b.engine
}

func (b *Builder) OrganizationRouter() *Builder {
	organizationRepository := repository.NewOrganizationRepository(b.repository)
	organizationHandler := handlers.NewOrganizationHandler(organizationRepository)

	v1 := b.engine.Group(apiV1)
	{
		v1.GET("/organizations", organizationHandler.GetAll)
	}

	admin := b.engine.Group(apiAdmin)
	{
		admin.POST("/organizations", organizationHandler.Create)
		admin.PATCH("/organizations/:uuid", organizationHandler.Update)
		admin.DELETE("/organizations/:uuid", organizationHandler.Delete)
	}
	return b
}

func (b *Builder) LoanApplicationsRouter() *Builder {
	loanApplicationsRepository := repository.NewLoanApplicationsRepository(b.repository)
	loanApplicationsHandler := handlers.NewLoanApplicationsHandler(loanApplicationsRepository)
	api := b.engine.Group("/api/v1")

	{
		api.GET("/loan_applications", loanApplicationsHandler.GetAll)
		api.POST("/loan_applications", loanApplicationsHandler.Create)
	}

	return b
}
