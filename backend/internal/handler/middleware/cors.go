package middleware

import "net/http"

// CORS はFrontendからBackend APIへの
// Cross-Origin Requestを許可するMiddleware。
//
// Local Developmentでは、
// Next.js:
//
//	http://localhost:3000
//
// Backend:
//
//	http://localhost:8080
//
// とOriginが異なるため、Browserから直接APIへアクセスする場合は
// CORS Headerが必要になる。
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.Header().Set(
				"Access-Control-Allow-Origin",
				"http://localhost:3000",
			)

			w.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, PATCH, DELETE, OPTIONS",
			)

			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization",
			)

			// BrowserはPOSTの前にOPTIONSで
			// Preflight Requestを送る場合がある。
			if r.Method == http.MethodOptions {
				w.WriteHeader(
					http.StatusNoContent,
				)

				return
			}

			next.ServeHTTP(
				w,
				r,
			)
		},
	)
}
