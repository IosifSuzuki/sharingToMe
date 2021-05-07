package midlleware

import (
	"IosifSuzuki/sharingToMe/internal/JWT"
	"IosifSuzuki/sharingToMe/internal/defaults"
	"IosifSuzuki/sharingToMe/internal/models"
	"IosifSuzuki/sharingToMe/internal/utility"
	"IosifSuzuki/sharingToMe/pkg/loger"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var time = time.Now()
		ip, _ := utility.GetIP(r)
		loger.PrintInfo("%s: \"%s\" try to get access to page, HTTP URL: \"%s\" header of request: ", time.Format("2006/01/02 15:04:05"), ip, r.URL.Path)
		for key, value := range r.Header {
			fmt.Printf("%s: %s\n", key, strings.Join(value, ","))
		}
		next.ServeHTTP(w, r)
	})
}

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token, ok := r.Header["Authorization"]; ok && token != nil {
			_, err := JWT.ValidateToken(token[0])
			if err != nil {
				w.Header().Set(defaults.ContentType, "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(struct {
					Error models.Error `json:"error"`
				}{
					Error: models.Error{
						Message: err.Error(),
						Code:    http.StatusUnauthorized,
					},
				})
				return
			}
		} else {
			w.Header().Set(defaults.ContentType, "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(struct {
				Error models.Error `json:"error"`
			}{
				Error: models.Error{
					Message: defaults.UnauthorizedUser,
					Code:    http.StatusUnauthorized,
				},
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
