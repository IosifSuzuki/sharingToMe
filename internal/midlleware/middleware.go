package midlleware

import (
	"IosifSuzuki/sharingToMe/internal/utility"
	"IosifSuzuki/sharingToMe/pkg/loger"
	"net/http"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteIp, err := utility.GetIP(r)
		if err != nil {
			loger.PrintError(err)
		}
		loger.PrintInfo("Remote Addr: %s Url request -> %s", remoteIp,  r.URL.Path)
		next.ServeHTTP(w, r)
	})
}