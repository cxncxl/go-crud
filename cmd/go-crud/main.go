package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cxcnxl/go-crud/internal/routes"
)

func main() {
    router := routes.NewRouter();

    const port int = 8080;
    addr := fmt.Sprintf(":%d", port);

    slog.Info(fmt.Sprintf("server is running on localhost:%d", port));

    err := http.ListenAndServe(addr, router);
    if err != nil {
        slog.Error("Error starting server: " + err.Error());
        panic(err);
    }
}
