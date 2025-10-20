package handler

type Handler struct {
	srv Service
}

func NewHandler(srv Service) *Handler {
	return &Handler{
		srv: srv,
	}
}
