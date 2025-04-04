package routes

import (
	"fmt"
	"log/slog"
	"net/http"
)

func NewRouter() *http.ServeMux {
    mux := http.NewServeMux();

    mux.HandleFunc("/", routeIndex);

    return mux;
}

func routeIndex(w http.ResponseWriter, r *http.Request) {
    _, err := fmt.Fprintf(w, "Hello, World!");
    if err != nil {
        slog.Error("Error wrtiting response: " + err.Error());
    }
}
