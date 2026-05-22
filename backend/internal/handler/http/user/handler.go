package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/handler/http/httputil"
	createuc "kanvix/backend/internal/usecase/user/create_in_workspace"
	deleteuc "kanvix/backend/internal/usecase/user/delete_from_workspace"
	listuc "kanvix/backend/internal/usecase/user/list_in_workspace"
	updateuc "kanvix/backend/internal/usecase/user/update_in_workspace"
)

type Handler struct {
	List   listuc.UseCase
	Create createuc.UseCase
	Update updateuc.UseCase
	Delete deleteuc.UseCase
}

type userDetailDTO struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	AvatarColor string `json:"avatarColor"`
	Role        string `json:"role"`
	Disabled    bool   `json:"disabled"`
}

func toDTO(u entity.WorkspaceUser) userDetailDTO {
	return userDetailDTO{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		AvatarColor: u.AvatarColor,
		Role:        u.Role,
		Disabled:    u.Disabled,
	}
}

func (h Handler) ListByWorkspaceHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	out, err := h.List.Execute(c.Request.Context(), listuc.In{UserID: userID, WorkspaceID: workspaceID})
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	users := make([]userDetailDTO, 0, len(out.Users))
	for _, u := range out.Users {
		users = append(users, toDTO(u))
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

type createRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role"`
}

func (h Handler) CreateInWorkspaceHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	out, err := h.Create.Execute(c.Request.Context(), createuc.In{
		UserID:      userID,
		WorkspaceID: workspaceID,
		Email:       req.Email,
		Password:    req.Password,
		Name:        req.Name,
		Role:        req.Role,
	})
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": toDTO(out.User)})
}

type patchRequest struct {
	Name *string `json:"name"`
	Role *string `json:"role"`
}

func (h Handler) PatchInWorkspaceHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	targetUserID := c.Param("userId")
	var req patchRequest
	if err := c.ShouldBindJSON(&req); err != nil || (req.Name == nil && req.Role == nil) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.Update.Execute(c.Request.Context(), updateuc.In{
		UserID:       userID,
		WorkspaceID:  workspaceID,
		TargetUserID: targetUserID,
		Name:         req.Name,
		Role:         req.Role,
	}); err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type deleteRequest struct {
	Action           string  `json:"action"`
	ReassignToUserID *string `json:"reassignToUserId"`
}

func (h Handler) DeleteFromWorkspaceHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	targetUserID := c.Param("userId")
	var req deleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.Delete.Execute(c.Request.Context(), deleteuc.In{
		UserID:          userID,
		WorkspaceID:     workspaceID,
		TargetUserID:    targetUserID,
		Action:          deleteuc.Action(req.Action),
		ReassignToUserID: req.ReassignToUserID,
	}); err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

