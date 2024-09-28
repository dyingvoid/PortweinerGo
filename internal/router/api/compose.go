package api

import (
	"fmt"
	"net/http"
	"start/internal/config"
	"start/internal/infrastructure"

	"github.com/gin-gonic/gin"
)

// DeployStackRequest represents the request body for deploying a stack
// @Description DeployStackRequest represents the request body for deploying a stack
type DeployStackRequest struct {
	// StackName is the name of the stack to deploy
	// @Description StackName is the name of the stack to deploy
	StackName string `json:"stack_name" binding:"required"`

	// Dir is the directory where the compose and .env files are located
	// @Description Dir is the directory where the compose and .env files are located
	Dir string `json:"directory" biniding:"required"`
}

// StopStackRequest represents the request body for stopping a stack
// @Description StopStackRequest represents the request body for stopping a stack
type StopStackRequest struct {
	// StackName is the name of the stack to stop
	// @Description StackName is the name of the stack to stop
	StackName string `json:"stack_name" binding:"required"`
}

// DeployStack godoc
// @Summary Deploy a stack
// @Description Deploy a stack using the provided stack name and directory
// @Tags Stack
// @Accept json
// @Produce json
// @Param request body DeployStackRequest true "DeployStackRequest"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /deploy [post]
func DeployStack(c *gin.Context, cfg *config.TerminalConfig) {
	var req DeployStackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": "Could not bind the body", "error": err.Error()},
		)
		return
	}

	dir := "uploads/" + req.Dir
	if err := infrastructure.ProduceFile(dir+"/.env",
		dir+"/docker-compose.yml", req.StackName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not produce file for deploy.",
			"error":   err.Error(),
		})
		return
	}

	t := infrastructure.NewTerminal(cfg)
	res, err := t.ExecuteCmd(fmt.Sprintf("docker stack deploy --compose-file docker-compose.yml %s", req.StackName), "uploads/"+req.Dir)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message":    "Cound not execute cmd command to deploy stack",
				"stack_name": req.StackName,
				"error":      err.Error(),
			})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": res})
}

// StopStack godoc
// @Summary Stop a stack
// @Description Stop a stack using the provided stack name
// @Tags Stack
// @Accept json
// @Produce json
// @Param request body StopStackRequest true "StopStackRequest"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /stop [post]
func StopStack(c *gin.Context, cfg *config.TerminalConfig) {
	var req StopStackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	t := infrastructure.NewTerminal(cfg)
	if _, err := t.ExecuteCmd(fmt.Sprintf("docker stack rm %s", req.StackName), ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not stop stack", "error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stack stopped(probably)",
	})
}
