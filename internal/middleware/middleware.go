package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

type Middleware func(n http.HandlerFunc) http.HandlerFunc;

func RecovererMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                slog.Error("Error handling http request:\n", err);
                debug.PrintStack();

                w.Header().Set("Connection", "close");
                w.WriteHeader(http.StatusInternalServerError);
                w.Write([]byte("Internal Server Error\n"));
            }
        }();

        next.ServeHTTP(w, r);
    });
}

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf(
            "[%v] -- %v -- (%v) %v %v\n",
            time.Now(),
            r.RemoteAddr,
            r.Header.Get("User-Agent"),
            r.Method,
            r.URL.Path,
        );

        next.ServeHTTP(w, r);
    });
}
