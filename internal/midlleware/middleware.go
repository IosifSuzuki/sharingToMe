package midlleware

import (
	"IosifSuzuki/sharingToMe/internal/utility"
	"IosifSuzuki/sharingToMe/pkg/loger"
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
