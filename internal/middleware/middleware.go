package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/cxcnxl/go-crud/internal/auth_helpers"
	"github.com/cxcnxl/go-crud/internal/responses"
)

type Middleware func(n http.HandlerFunc) http.HandlerFunc;
type MiddlewareSet []Middleware;

func (self MiddlewareSet) With(others ...Middleware) MiddlewareSet {
    return append(self, others...);
}

func Wrap (
    handler http.HandlerFunc,
    middlewares MiddlewareSet,
) http.HandlerFunc {
    wrapped := handler;

    for i := len(middlewares) - 1; i >= 0; i-- {
        wrapped = middlewares[i](wrapped)
    }

    return wrapped;
}

var UtilMiddleware = MiddlewareSet{
    RecovererMiddleware,
    LoggerMiddleware,
    JSONResponserMiddleware, // (at least) for now all responses are JSON
};

var AuthMiddleware = append(UtilMiddleware, JWTAutherMiddleware);

// --------- Implementations --------

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

func JWTAutherMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        throwUnauthorized := func() {
            error := responses.NewErrorResponse("unauthorized");
            w.Header().Add("Content-Type", "application/json");
            w.WriteHeader(http.StatusUnauthorized);
            w.Write(error.Json());
        }

        auth := r.Header.Get("Authorization");
        if auth == "" {
            throwUnauthorized();
            return;
        }

        auth = strings.TrimSpace(auth);
        parts := strings.Split(auth, " ");
        if len(parts) < 2 {
            throwUnauthorized();
            return;
        }

        token := parts[1];
        claims, err := auth_helpers.DecodeJWT(token);
        if err != nil && errors.Is(err, auth_helpers.InvalidJwtError{}) {
            throwUnauthorized();
            return;
        }

        ctx:= context.WithValue(r.Context(), "auth", claims);
        r = r.WithContext(ctx);

        next.ServeHTTP(w, r);
    });
}

func JSONResponserMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json");
        
        next.ServeHTTP(w, r)
    })
}

func POSTHandlerMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            error := responses.NewErrorResponse("Only POST method allowed");
            http.Error(w, error.JsonString(), http.StatusMethodNotAllowed);
            return;
        }

        next.ServeHTTP(w, r);
    });
}

func GETHandlerMiddleware (next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "GET" {
            error := responses.NewErrorResponse("Only GET method allowed");
            http.Error(w, error.JsonString(), http.StatusMethodNotAllowed);
            return;
        }

        next.ServeHTTP(w, r);
    });
}

func PUTHandlerMiddleware (next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "PUT" {
            error := responses.NewErrorResponse("Only PUT method allowed");
            http.Error(w, error.JsonString(), http.StatusMethodNotAllowed);
            return;
        }

        next.ServeHTTP(w, r);
    });
}

func DELETEHandlerMiddleware (next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "DELETE" {
            error := responses.NewErrorResponse("Only DELETE method allowed");
            http.Error(w, error.JsonString(), http.StatusMethodNotAllowed);
            return;
        }

        next.ServeHTTP(w, r);
    });
}
