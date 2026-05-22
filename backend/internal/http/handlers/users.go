package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/services"
)

type UsersHandler struct {
	Service services.AppService
}

func (h UsersHandler) ListByWorkspace(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	users, err := h.Service.ListWorkspaceUsers(c.Request.Context(), userID, workspaceID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

type createUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role"`
}

func (h UsersHandler) CreateInWorkspace(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	u, err := h.Service.CreateWorkspaceUser(c.Request.Context(), userID, workspaceID, req.Email, req.Password, req.Name, req.Role)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": u})
}

type updateUserRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

func (h UsersHandler) UpdateRoleInWorkspace(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	targetUserID := c.Param("userId")
	var req updateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.Service.UpdateWorkspaceUserRole(c.Request.Context(), userID, workspaceID, targetUserID, req.Role); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type patchWorkspaceUserRequest struct {
	Name *string `json:"name"`
	Role *string `json:"role"`
}

func (h UsersHandler) PatchInWorkspace(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	targetUserID := c.Param("userId")
	var req patchWorkspaceUserRequest
	if err := c.ShouldBindJSON(&req); err != nil || (req.Name == nil && req.Role == nil) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.Service.UpdateWorkspaceUser(c.Request.Context(), userID, workspaceID, targetUserID, req.Name, req.Role); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type deleteWorkspaceUserRequest struct {
	Action           string  `json:"action"`
	ReassignToUserID *string `json:"reassignToUserId"`
}

func (h UsersHandler) DeleteFromWorkspace(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	targetUserID := c.Param("userId")

	var req deleteWorkspaceUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.Service.DeleteWorkspaceUser(
		c.Request.Context(),
		userID,
		workspaceID,
		targetUserID,
		services.DeleteWorkspaceUserAction(req.Action),
		req.ReassignToUserID,
	); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
