package router

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/config"
	authhandler "kanvix/backend/internal/handler/http/auth"
	bootstraphandler "kanvix/backend/internal/handler/http/bootstrap"
	columnhandler "kanvix/backend/internal/handler/http/column"
	labelhandler "kanvix/backend/internal/handler/http/label"
	projecthandler "kanvix/backend/internal/handler/http/project"
	statssearchhandler "kanvix/backend/internal/handler/http/statssearch"
	subtaskhandler "kanvix/backend/internal/handler/http/subtask"
	taskhandler "kanvix/backend/internal/handler/http/task"
	userhandler "kanvix/backend/internal/handler/http/user"
	workspacehandler "kanvix/backend/internal/handler/http/workspace"
	"kanvix/backend/internal/middleware"
	"kanvix/backend/internal/usecase/auth/login"
	"kanvix/backend/internal/usecase/auth/logout"
	"kanvix/backend/internal/usecase/auth/me"
	"kanvix/backend/internal/usecase/auth/register"
	"kanvix/backend/internal/usecase/auth/setup"
)

type AuthDeps struct {
	Setup    setup.UseCase
	Register register.UseCase
	Login    login.UseCase
	Me       me.UseCase
	Logout   logout.UseCase
}

type Deps struct {
	Log         *slog.Logger
	Cfg         config.Config
	Auth        AuthDeps
	Workspace   workspacehandler.Handler
	Tasks       taskhandler.Handler
	Bootstrap   bootstraphandler.Handler
	Projects    projecthandler.Handler
	Columns     columnhandler.Handler
	Labels      labelhandler.Handler
	Subtasks    subtaskhandler.Handler
	StatsSearch statssearchhandler.Handler
	Users       userhandler.Handler
}

func New(d Deps) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(d.Log))
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(d.Cfg.AllowedOrigins))

	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	authHandler := authhandler.Handler{
		Cfg:      d.Cfg,
		Setup:    d.Auth.Setup,
		Register: d.Auth.Register,
		Login:    d.Auth.Login,
		Me:       d.Auth.Me,
		Logout:   d.Auth.Logout,
	}
	authmw := middleware.Auth{CookieName: d.Cfg.CookieName, Me: d.Auth.Me}

	api := r.Group("/api")
	{
		api.GET("/auth/setup", authHandler.SetupHandler)
		api.POST("/auth/first-signup", authHandler.FirstSignupHandler)
		api.POST("/auth/login", authHandler.LoginHandler)
		api.GET("/auth/me", authHandler.MeHandler)
		api.POST("/auth/logout", authmw.Required(), authHandler.LogoutHandler)

		protected := api.Group("")
		protected.Use(authmw.Required())

		protected.GET("/bootstrap", d.Bootstrap.GetHandler)

		protected.GET("/workspaces", d.Workspace.ListHandler)
		protected.POST("/workspaces", d.Workspace.CreateHandler)
		protected.PATCH("/workspaces/:workspaceId", d.Workspace.UpdateHandler)
		protected.DELETE("/workspaces/:workspaceId", d.Workspace.DeleteHandler)

		protected.GET("/workspaces/:workspaceId/projects", d.Projects.ListByWorkspaceHandler)
		protected.POST("/workspaces/:workspaceId/projects", d.Projects.CreateHandler)
		protected.PATCH("/projects/:projectId", d.Projects.UpdateHandler)
		protected.DELETE("/projects/:projectId", d.Projects.DeleteHandler)

		protected.GET("/projects/:projectId/columns", d.Columns.ListByProjectHandler)
		protected.POST("/projects/:projectId/columns", d.Columns.CreateHandler)
		protected.PATCH("/columns/:columnId", d.Columns.UpdateHandler)
		protected.DELETE("/columns/:columnId", d.Columns.DeleteHandler)

		protected.GET("/projects/:projectId/tasks", d.Tasks.ListByProjectHandler)
		protected.POST("/projects/:projectId/tasks", d.Tasks.CreateHandler)
		protected.GET("/tasks/:taskId", d.Tasks.GetHandler)
		protected.PATCH("/tasks/:taskId", d.Tasks.UpdateHandler)
		protected.DELETE("/tasks/:taskId", d.Tasks.DeleteHandler)
		protected.POST("/tasks/:taskId/move", d.Tasks.MoveHandler)
		protected.GET("/tasks/:taskId/subtasks", d.Subtasks.ListByTaskHandler)
		protected.POST("/tasks/:taskId/subtasks", d.Subtasks.CreateHandler)
		protected.PATCH("/subtasks/:subtaskId", d.Subtasks.PatchHandler)
		protected.DELETE("/subtasks/:subtaskId", d.Subtasks.DeleteHandler)

		protected.GET("/workspaces/:workspaceId/labels", d.Labels.ListByWorkspaceHandler)
		protected.POST("/workspaces/:workspaceId/labels", d.Labels.CreateHandler)
		protected.PATCH("/labels/:labelId", d.Labels.UpdateHandler)
		protected.DELETE("/labels/:labelId", d.Labels.DeleteHandler)

		protected.GET("/workspaces/:workspaceId/stats", d.StatsSearch.WorkspaceStatsHandler)
		protected.GET("/workspaces/:workspaceId/search", d.StatsSearch.SearchHandler)

		protected.GET("/workspaces/:workspaceId/users", d.Users.ListByWorkspaceHandler)
		protected.POST("/workspaces/:workspaceId/users", d.Users.CreateInWorkspaceHandler)
		protected.PATCH("/workspaces/:workspaceId/users/:userId", d.Users.PatchInWorkspaceHandler)
		protected.DELETE("/workspaces/:workspaceId/users/:userId", d.Users.DeleteFromWorkspaceHandler)
	}

	return r
}

