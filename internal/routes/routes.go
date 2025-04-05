package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cxcnxl/go-crud/internal/middleware"
)

func NewRouter() *http.ServeMux {
    mux := http.NewServeMux();

    mux.HandleFunc("/", applyMiddleware(routeIndex, utilMiddleware));

    return mux;
}

func routeIndex(w http.ResponseWriter, r *http.Request) {
    _, err := fmt.Fprintf(w, "Hello, World!");
    if err != nil {
        slog.Error("Error wrtiting response: " + err.Error());
    }
}

func applyMiddleware (
    handler http.HandlerFunc,
    middlewares []middleware.Middleware,
) http.HandlerFunc {
    wrapped := handler;

    for _, m := range middlewares {
        wrapped = m(wrapped);
    }

    return wrapped;
}

var utilMiddleware = []middleware.Middleware{
    middleware.LoggerMiddleware,
    middleware.RecovererMiddleware,
};
