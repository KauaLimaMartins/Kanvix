package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kanvix/backend/internal/config"
	bootstraphandler "kanvix/backend/internal/handler/http/bootstrap"
	columnhandler "kanvix/backend/internal/handler/http/column"
	labelhandler "kanvix/backend/internal/handler/http/label"
	projecthandler "kanvix/backend/internal/handler/http/project"
	httprouter "kanvix/backend/internal/handler/http/router"
	statssearchhandler "kanvix/backend/internal/handler/http/statssearch"
	subtaskhandler "kanvix/backend/internal/handler/http/subtask"
	taskhandler "kanvix/backend/internal/handler/http/task"
	userhandler "kanvix/backend/internal/handler/http/user"
	workspacehandler "kanvix/backend/internal/handler/http/workspace"
	"kanvix/backend/internal/infra/cache/redis"
	"kanvix/backend/internal/infra/database/postgres"
	dbrepo "kanvix/backend/internal/infra/database/postgres/repository"
	"kanvix/backend/internal/infra/logging"
	"kanvix/backend/internal/usecase/auth/login"
	"kanvix/backend/internal/usecase/auth/logout"
	"kanvix/backend/internal/usecase/auth/me"
	"kanvix/backend/internal/usecase/auth/register"
	"kanvix/backend/internal/usecase/auth/setup"
	bootstrapget "kanvix/backend/internal/usecase/bootstrap/get"
	columncreate "kanvix/backend/internal/usecase/column/create"
	columndelete "kanvix/backend/internal/usecase/column/delete"
	columnlist "kanvix/backend/internal/usecase/column/list"
	columnupdate "kanvix/backend/internal/usecase/column/update"
	labelcreate "kanvix/backend/internal/usecase/label/create"
	labeldelete "kanvix/backend/internal/usecase/label/delete"
	labellist "kanvix/backend/internal/usecase/label/list"
	labelupdate "kanvix/backend/internal/usecase/label/update"
	projectcreate "kanvix/backend/internal/usecase/project/create"
	projectdelete "kanvix/backend/internal/usecase/project/delete"
	projectlist "kanvix/backend/internal/usecase/project/list"
	projectupdate "kanvix/backend/internal/usecase/project/update"
	searchws "kanvix/backend/internal/usecase/search/workspace"
	statsws "kanvix/backend/internal/usecase/stats/workspace"
	subtaskcreate "kanvix/backend/internal/usecase/subtask/create"
	subtaskdelete "kanvix/backend/internal/usecase/subtask/delete"
	subtasklist "kanvix/backend/internal/usecase/subtask/list"
	subtaskupdate "kanvix/backend/internal/usecase/subtask/update"
	taskcreate "kanvix/backend/internal/usecase/task/create"
	taskdelete "kanvix/backend/internal/usecase/task/delete"
	taskget "kanvix/backend/internal/usecase/task/get"
	tasklist "kanvix/backend/internal/usecase/task/list"
	taskmove "kanvix/backend/internal/usecase/task/move"
	taskupdate "kanvix/backend/internal/usecase/task/update"
	usercreate "kanvix/backend/internal/usecase/user/create_in_workspace"
	userdelete "kanvix/backend/internal/usecase/user/delete_from_workspace"
	userlist "kanvix/backend/internal/usecase/user/list_in_workspace"
	userupdate "kanvix/backend/internal/usecase/user/update_in_workspace"
	workspacecreate "kanvix/backend/internal/usecase/workspace/create"
	workspacedelete "kanvix/backend/internal/usecase/workspace/delete"
	workspacelist "kanvix/backend/internal/usecase/workspace/list"
	workspaceupdate "kanvix/backend/internal/usecase/workspace/update"
)

func main() {
	log := logging.New()
	slog.SetDefault(log)

	cfg := config.Load()
	ctx := context.Background()

	db, err := postgres.Open(ctx, log, cfg.DBPath)
	if err != nil {
		log.Error("db open failed", "error", err)
		os.Exit(1)
	}
	if err := postgres.AutoMigrate(ctx, db.Gorm); err != nil {
		log.Error("db migrate failed", "error", err)
		os.Exit(1)
	}
	if err := postgres.EnsureOwnerMemberships(ctx, db.Gorm); err != nil {
		log.Error("db memberships failed", "error", err)
		os.Exit(1)
	}

	redisClient := redis.New(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err := redisClient.Ping(ctx); err != nil {
		log.Error("redis unavailable", "error", err)
		os.Exit(1)
	}

	usersRepo := dbrepo.Store{DB: db.Gorm}
	workspacesRepo := dbrepo.Workspaces{DB: db.Gorm}
	workspaceMembersRepo := dbrepo.WorkspaceMembers{DB: db.Gorm}
	projectsRepo := dbrepo.Projects{DB: db.Gorm}
	columnsRepo := dbrepo.Columns{DB: db.Gorm}
	tasksRepo := dbrepo.Tasks{DB: db.Gorm}
	labelsRepo := dbrepo.Labels{DB: db.Gorm}
	taskLabelsRepo := dbrepo.TaskLabels{DB: db.Gorm}
	subtasksRepo := dbrepo.Subtasks{DB: db.Gorm}
	sessions := redis.SessionStore{Client: redisClient.Raw}

	authSetupUC := setup.UseCase{Users: usersRepo}
	authRegisterUC := register.UseCase{
		Users:      usersRepo,
		Workspaces: workspacesRepo,
		Members:    workspaceMembersRepo,
		Sessions:   sessions,
		SessionTTL: cfg.SessionTTL,
	}
	authLoginUC := login.UseCase{Users: usersRepo, Sessions: sessions, SessionTTL: cfg.SessionTTL}
	authMeUC := me.UseCase{Users: usersRepo, Sessions: sessions}
	authLogoutUC := logout.UseCase{Sessions: sessions}

	workspaceListUC := workspacelist.UseCase{Memberships: workspaceMembersRepo}
	workspaceCreateUC := workspacecreate.UseCase{Users: usersRepo, Workspaces: workspacesRepo, Members: workspaceMembersRepo}
	workspaceUpdateUC := workspaceupdate.UseCase{Workspaces: workspacesRepo, Members: workspaceMembersRepo}
	workspaceDeleteUC := workspacedelete.UseCase{Workspaces: workspacesRepo, Members: workspaceMembersRepo}

	taskListUC := tasklist.UseCase{
		Projects:   projectsRepo,
		Tasks:      tasksRepo,
		TaskLabels: taskLabelsRepo,
		Members:    workspaceMembersRepo,
		Workspaces: workspacesRepo,
	}
	taskGetUC := taskget.UseCase{
		Projects:   projectsRepo,
		Tasks:      tasksRepo,
		TaskLabels: taskLabelsRepo,
		Members:    workspaceMembersRepo,
		Workspaces: workspacesRepo,
	}
	taskCreateUC := taskcreate.UseCase{
		Projects:   projectsRepo,
		Tasks:      tasksRepo,
		Members:    workspaceMembersRepo,
		Workspaces: workspacesRepo,
	}
	taskUpdateUC := taskupdate.UseCase{
		Projects:   projectsRepo,
		Tasks:      tasksRepo,
		TaskLabels: taskLabelsRepo,
		Members:    workspaceMembersRepo,
		Workspaces: workspacesRepo,
	}
	taskDeleteUC := taskdelete.UseCase{
		Projects:   projectsRepo,
		Tasks:      tasksRepo,
		Members:    workspaceMembersRepo,
		Workspaces: workspacesRepo,
	}
	taskMoveUC := taskmove.UseCase{
		Projects:   projectsRepo,
		Tasks:      tasksRepo,
		Members:    workspaceMembersRepo,
		Workspaces: workspacesRepo,
	}

	workspaceHTTP := workspacehandler.Handler{
		List:   workspaceListUC,
		Create: workspaceCreateUC,
		Update: workspaceUpdateUC,
		Delete: workspaceDeleteUC,
	}
	taskHTTP := taskhandler.Handler{
		List:   taskListUC,
		Get:    taskGetUC,
		Create: taskCreateUC,
		Update: taskUpdateUC,
		Delete: taskDeleteUC,
		Move:   taskMoveUC,
	}

	projectListUC := projectlist.UseCase{Projects: projectsRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	projectCreateUC := projectcreate.UseCase{Projects: projectsRepo, Columns: columnsRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	projectUpdateUC := projectupdate.UseCase{Projects: projectsRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	projectDeleteUC := projectdelete.UseCase{Projects: projectsRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	projectHTTP := projecthandler.Handler{List: projectListUC, Create: projectCreateUC, Update: projectUpdateUC, Delete: projectDeleteUC}

	columnListUC := columnlist.UseCase{Projects: projectsRepo, Columns: columnsRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	columnCreateUC := columncreate.UseCase{Projects: projectsRepo, Columns: columnsRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	columnUpdateUC := columnupdate.UseCase{Projects: projectsRepo, Columns: columnsRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	columnDeleteUC := columndelete.UseCase{Projects: projectsRepo, Columns: columnsRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	columnHTTP := columnhandler.Handler{List: columnListUC, Create: columnCreateUC, Update: columnUpdateUC, Delete: columnDeleteUC}

	labelListUC := labellist.UseCase{Labels: labelsRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	labelCreateUC := labelcreate.UseCase{Labels: labelsRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	labelUpdateUC := labelupdate.UseCase{Labels: labelsRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	labelDeleteUC := labeldelete.UseCase{Labels: labelsRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	labelHTTP := labelhandler.Handler{List: labelListUC, Create: labelCreateUC, Update: labelUpdateUC, Delete: labelDeleteUC}

	subtaskListUC := subtasklist.UseCase{Projects: projectsRepo, Tasks: tasksRepo, Subtasks: subtasksRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	subtaskCreateUC := subtaskcreate.UseCase{Projects: projectsRepo, Tasks: tasksRepo, Subtasks: subtasksRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	subtaskUpdateUC := subtaskupdate.UseCase{Projects: projectsRepo, Tasks: tasksRepo, Subtasks: subtasksRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	subtaskDeleteUC := subtaskdelete.UseCase{Projects: projectsRepo, Tasks: tasksRepo, Subtasks: subtasksRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	subtaskHTTP := subtaskhandler.Handler{List: subtaskListUC, Create: subtaskCreateUC, Update: subtaskUpdateUC, Delete: subtaskDeleteUC}

	statsRepo := dbrepo.StatsSearch{DB: db.Gorm}
	cache := redis.Cache{Client: redisClient.Raw}
	statsUC := statsws.UseCase{Stats: statsRepo, Cache: cache, CacheTTL: cfg.CacheTTL, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	searchUC := searchws.UseCase{Search: statsRepo, Cache: cache, CacheTTL: cfg.CacheTTL, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	statsSearchHTTP := statssearchhandler.Handler{Stats: statsUC, Search: searchUC}

	userListUC := userlist.UseCase{Users: usersRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	userCreateUC := usercreate.UseCase{Users: usersRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo, Memberships: workspaceMembersRepo}
	userUpdateUC := userupdate.UseCase{Memberships: workspaceMembersRepo, Users: usersRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	userDeleteUC := userdelete.UseCase{Memberships: workspaceMembersRepo, Users: usersRepo, Tasks: tasksRepo, Members: workspaceMembersRepo, Workspaces: workspacesRepo}
	userHTTP := userhandler.Handler{List: userListUC, Create: userCreateUC, Update: userUpdateUC, Delete: userDeleteUC}

	bootstrapUC := bootstrapget.UseCase{
		Users:      usersRepo,
		Members:    workspaceMembersRepo,
		Workspaces: workspacesRepo,
		Projects:   projectsRepo,
		Columns:    columnsRepo,
		Tasks:      tasksRepo,
		Labels:     labelsRepo,
		TaskLabels: taskLabelsRepo,
	}
	bootstrapHTTP := bootstraphandler.Handler{Get: bootstrapUC}

	r := httprouter.New(httprouter.Deps{
		Log: log,
		Cfg: cfg,
		Auth: httprouter.AuthDeps{
			Setup:    authSetupUC,
			Register: authRegisterUC,
			Login:    authLoginUC,
			Me:       authMeUC,
			Logout:   authLogoutUC,
		},
		Workspace:   workspaceHTTP,
		Tasks:       taskHTTP,
		Bootstrap:   bootstrapHTTP,
		Projects:    projectHTTP,
		Columns:     columnHTTP,
		Labels:      labelHTTP,
		Subtasks:    subtaskHTTP,
		StatsSearch: statsSearchHTTP,
		Users:       userHTTP,
	})

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("server listening", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
		}
	}()

	<-stop
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = redisClient.Raw.Close()
	_ = db.SQL.Close()
	_ = srv.Shutdown(shutdownCtx)
}
