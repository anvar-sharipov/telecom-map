package middleware

import (
	"net/http"

	"github.com/anvar-sharipov/telecom-map/internal/utils"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

// func ErrorMiddleware(h AppHandler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		err := h(w, r)
// 		if err == nil {
// 			return
// 		}

// 		if apiErr, ok := err.(*utils.APIError); ok {
// 			utils.WriteJSON(w, apiErr.Status, map[string]string{
// 				"error": apiErr.Message,
// 			})
// 			return
// 		}

// 		// fallback
// 		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
// 			"error": "internal server error",
// 		})
// 	}
// }

func ErrorMiddleware(next func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			if apiErr, ok := err.(*utils.APIError); ok {
				utils.WriteJSON(w, apiErr.Status, map[string]string{
					"error": apiErr.Message,
				})
				return
			}

			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "internal server error",
			})
		}
	}
}
