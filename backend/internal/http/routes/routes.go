package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/config"
	"kanvix/backend/internal/http/handlers"
	"kanvix/backend/internal/http/middleware"
	"kanvix/backend/internal/services"
)

type Deps struct {
	Log        *slog.Logger
	Cfg        config.Config
	Auth       services.AuthService
	AppService services.AppService
}

func NewRouter(d Deps) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(d.Log))
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(d.Cfg.AllowedOrigins))

	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	authHandler := handlers.AuthHandler{Cfg: d.Cfg, Service: d.Auth}
	authmw := middleware.Auth{CookieName: d.Cfg.CookieName, Service: d.Auth}

	api := r.Group("/api")
	{
		api.GET("/auth/setup", authHandler.Setup)
		api.POST("/auth/first-signup", authHandler.FirstSignup)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/auth/me", authHandler.Me)
		api.POST("/auth/logout", authmw.Required(), authHandler.Logout)

		protected := api.Group("")
		protected.Use(authmw.Required())

		bootstrap := handlers.BootstrapHandler{Service: d.AppService}
		protected.GET("/bootstrap", bootstrap.Get)

		workspaces := handlers.WorkspacesHandler{Service: d.AppService}
		projects := handlers.ProjectsHandler{Service: d.AppService}
		columns := handlers.ColumnsHandler{Service: d.AppService}
		tasks := handlers.TasksHandler{Service: d.AppService}
		subtasks := handlers.SubtasksHandler{Service: d.AppService}
		labels := handlers.LabelsHandler{Service: d.AppService}
		users := handlers.UsersHandler{Service: d.AppService}
		statsSearch := handlers.StatsSearchHandler{Service: d.AppService}

		protected.GET("/workspaces", workspaces.List)
		protected.POST("/workspaces", workspaces.Create)
		protected.PATCH("/workspaces/:workspaceId", workspaces.Update)
		protected.DELETE("/workspaces/:workspaceId", workspaces.Delete)

		protected.GET("/workspaces/:workspaceId/projects", projects.ListByWorkspace)
		protected.POST("/workspaces/:workspaceId/projects", projects.Create)
		protected.PATCH("/projects/:projectId", projects.Update)
		protected.DELETE("/projects/:projectId", projects.Delete)

		protected.GET("/projects/:projectId/columns", columns.ListByProject)
		protected.POST("/projects/:projectId/columns", columns.Create)
		protected.PATCH("/columns/:columnId", columns.Update)
		protected.DELETE("/columns/:columnId", columns.Delete)

		protected.GET("/projects/:projectId/tasks", tasks.ListByProject)
		protected.POST("/projects/:projectId/tasks", tasks.Create)
		protected.GET("/tasks/:taskId", tasks.Get)
		protected.PATCH("/tasks/:taskId", tasks.Update)
		protected.DELETE("/tasks/:taskId", tasks.Delete)
		protected.POST("/tasks/:taskId/move", tasks.Move)
		protected.GET("/tasks/:taskId/subtasks", subtasks.ListByTask)
		protected.POST("/tasks/:taskId/subtasks", subtasks.Create)
		protected.PATCH("/subtasks/:subtaskId", subtasks.Patch)
		protected.DELETE("/subtasks/:subtaskId", subtasks.Delete)

		protected.GET("/workspaces/:workspaceId/labels", labels.ListByWorkspace)
		protected.POST("/workspaces/:workspaceId/labels", labels.Create)
		protected.PATCH("/labels/:labelId", labels.Update)
		protected.DELETE("/labels/:labelId", labels.Delete)

		protected.GET("/workspaces/:workspaceId/stats", statsSearch.WorkspaceStats)
		protected.GET("/workspaces/:workspaceId/search", statsSearch.Search)

		protected.GET("/workspaces/:workspaceId/users", users.ListByWorkspace)
		protected.POST("/workspaces/:workspaceId/users", users.CreateInWorkspace)
		protected.PATCH("/workspaces/:workspaceId/users/:userId", users.PatchInWorkspace)
		protected.DELETE("/workspaces/:workspaceId/users/:userId", users.DeleteFromWorkspace)
	}

	return r
}
