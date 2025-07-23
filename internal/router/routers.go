package router

import (
	"app_aggregator/internal/handlers"
	"app_aggregator/internal/repository"
	"github.com/gin-gonic/gin"
)

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

func (b *Builder) OrganizationRouter() {
	organizationRepository := repository.NewOrganizationRepository(b.repository)
	organizationHandler := handlers.NewOrganizationHandler(organizationRepository)

	api := b.engine.Group("/api/v1")
	{
		api.GET("/organizations", organizationHandler.GetAll)
		api.POST("/organizations", organizationHandler.Create)
		api.PATCH("/admin/organizations/:uuid", organizationHandler.Update)
		api.DELETE("/admin/organizations/:uuid", organizationHandler.Delete)
	}
}
