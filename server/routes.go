package server

import (
	"FireBaseEx/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	chi.Router
}

func SetupRoutes() *Server {
	router := chi.NewRouter()
	router.Route("/api", func(api chi.Router) {
		api.Post("/image", handlers.UploadImage)

	})
	return &Server{router}
}
func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
