package apiserver

import (
	"github.com/aabbuukkaarr8/internal/config"
	"github.com/aabbuukkaarr8/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/zlog"
	"net/http"
)

type APIServer struct {
	config *config.Config
	router *gin.Engine
}

func New(config *config.Config) *APIServer {
	return &APIServer{
		config: config,
		router: gin.Default(),
	}

}
func (s *APIServer) Run() error {
	zlog.Logger.Info().Msg("Starting API server")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) ConfigureRouter(handler *handler.Handler) {

}
