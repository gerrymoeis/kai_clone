package routes

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
)

func init() {
    RegisterRoute(func(r chi.Router) {
        r.Get("/api/availability", handleAvailabilityAPI)
        RegisterURL("/api/availability")
    })
}

func handleAvailabilityAPI(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    
    // TODO: Implement your API logic here
    response := map[string]interface{}{
        "success": true,
        "message": "Availability API endpoint",
        "method":  "GET",
    }
    
    _ = json.NewEncoder(w).Encode(response)
}
