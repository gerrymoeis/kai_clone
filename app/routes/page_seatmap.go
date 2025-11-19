package routes

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "gothicforge3/app/templates"
)

func init() {
    RegisterRoute(func(r chi.Router) {
        r.Get("/seatmap", func(w http.ResponseWriter, req *http.Request) {
            w.Header().Set("Content-Type", "text/html; charset=utf-8")
            _ = templates.PageSeatmap().Render(req.Context(), w)
        })
        RegisterURL("/seatmap")
    })
}
