package handler

import (
	"github.com/aabbuukkaarr8/internal/model"
	"github.com/aabbuukkaarr8/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/zlog"
	"net/http"
)

func (h *Handler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		zlog.Logger.Err(err).Msg("upload failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	var req Request
	err = validator.BindJSON(&req, c.Request)
	if err != nil {
		zlog.Logger.Err(err).Msg("upload failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	action := model.Action{
		Name:   req.Action,
		Params: req.Params,
	}
	id, dst, err := h.srv.SaveImage(c.Request.Context(), "original", header.Filename, file, action)

}
