package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"newstrix/internal/api/handler"
	"newstrix/internal/search"
)

type Router struct {
	r *chi.Mux
}

func SetupRouter(engine *search.SearchEngine) *Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/search", func(r chi.Router) {
		sh := handler.NewSearchHandler(engine)
		r.Get("/semantic", sh.SemanticSearch)
		//r.Get("/date", sh.SearchByDate)
		//r.Get("/{id}", sh.GetByID)
	})

	return &Router{r: r}
}

func (router *Router) Run(addr string) error {
	return http.ListenAndServe(addr, router.r)
}
