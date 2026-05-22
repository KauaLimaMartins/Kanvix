package bootstrap

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/handler/http/httputil"
	"kanvix/backend/internal/usecase/bootstrap/get"
)

type Handler struct {
	Get get.UseCase
}

type meUserDTO struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	AvatarColor string `json:"avatarColor"`
	Role        string `json:"role"`
}

type workspaceDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon,omitempty"`
	Color string `json:"color,omitempty"`
	Role  string `json:"role,omitempty"`
}

type projectDTO struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspaceId"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type columnDTO struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
}

type taskDTO struct {
	ID          string   `json:"id"`
	ProjectID   string   `json:"projectId"`
	ColumnID    string   `json:"columnId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
	DueDate     *string  `json:"dueDate,omitempty"`
	AssigneeID  *string  `json:"assigneeId,omitempty"`
	Order       int      `json:"order"`
	CreatedAt   string   `json:"createdAt"`
}

type labelDTO struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspaceId"`
	Name        string `json:"name"`
	Color       string `json:"color"`
}

type userDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AvatarColor string `json:"avatarColor"`
	Disabled    bool   `json:"disabled"`
}

func (h Handler) GetHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	out, err := h.Get.Execute(c.Request.Context(), userID)
	if err != nil {
		httputil.RespondError(c, err)
		return
	}

	resp := struct {
		User       meUserDTO     `json:"user"`
		Workspaces []workspaceDTO `json:"workspaces"`
		Projects   []projectDTO  `json:"projects"`
		Columns    []columnDTO   `json:"columns"`
		Tasks      []taskDTO     `json:"tasks"`
		Labels     []labelDTO    `json:"labels"`
		Users      []userDTO     `json:"users"`
	}{
		User: meUserDTO{
			ID:          out.User.ID,
			Email:       out.User.Email,
			Name:        out.User.Name,
			AvatarColor: out.User.AvatarColor,
			Role:        out.User.Role,
		},
		Workspaces: make([]workspaceDTO, 0, len(out.Workspaces)),
		Projects:   make([]projectDTO, 0, len(out.Projects)),
		Columns:    make([]columnDTO, 0, len(out.Columns)),
		Tasks:      make([]taskDTO, 0, len(out.Tasks)),
		Labels:     make([]labelDTO, 0, len(out.Labels)),
		Users:      make([]userDTO, 0, len(out.Users)),
	}

	for _, u := range out.Users {
		resp.Users = append(resp.Users, userDTO{ID: u.ID, Name: u.Name, AvatarColor: u.AvatarColor, Disabled: u.Disabled})
	}
	for _, w := range out.Workspaces {
		resp.Workspaces = append(resp.Workspaces, workspaceDTO{
			ID:    w.ID,
			Name:  w.Name,
			Icon:  w.Icon,
			Color: w.Color,
			Role:  out.RoleByWS[w.ID],
		})
	}
	for _, p := range out.Projects {
		resp.Projects = append(resp.Projects, projectDTO{
			ID:          p.ID,
			WorkspaceID: p.WorkspaceID,
			Name:        p.Name,
			Description: p.Description,
		})
	}
	for _, c2 := range out.Columns {
		resp.Columns = append(resp.Columns, columnDTO{ID: c2.ID, ProjectID: c2.ProjectID, Name: c2.Name, Order: c2.Order})
	}
	for _, l := range out.Labels {
		resp.Labels = append(resp.Labels, labelDTO{ID: l.ID, WorkspaceID: l.WorkspaceID, Name: l.Name, Color: l.Color})
	}
	for _, t := range out.Tasks {
		lbls := out.TaskLabels[t.ID]
		if lbls == nil {
			lbls = []string{}
		}
		resp.Tasks = append(resp.Tasks, taskDTO{
			ID:          t.ID,
			ProjectID:   t.ProjectID,
			ColumnID:    t.ColumnID,
			Title:       t.Title,
			Description: t.Description,
			Labels:      lbls,
			DueDate:     t.DueDate,
			AssigneeID:  t.AssigneeID,
			Order:       t.Order,
			CreatedAt:   t.CreatedAt.UTC().Format(time.RFC3339Nano),
		})
	}

	c.JSON(http.StatusOK, resp)
}

