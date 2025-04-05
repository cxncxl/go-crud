package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"gorm.io/gorm"

	"github.com/cxcnxl/go-crud/internal/app_service"
	"github.com/cxcnxl/go-crud/internal/dto"
	"github.com/cxcnxl/go-crud/internal/middleware"
	"github.com/cxcnxl/go-crud/internal/responses"
)

func NewRouter(db *gorm.DB) *http.ServeMux {
    mux := http.NewServeMux();
    service := appservice.NewAppService(db);

    mux.HandleFunc(
        "/",
        applyMiddleware(routeIndex(service), utilMiddleware),
    );
    mux.HandleFunc(
        "/register",
        applyMiddleware(routeRegister(service), utilMiddleware),
    );
    mux.HandleFunc(
        "/me",
        applyMiddleware(routeMe(service), authMiddleware),
    );

    return mux;
}

func routeIndex(_ *appservice.AppService) http.HandlerFunc {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        _, err := fmt.Fprintf(w, "Hello, World!");
        if err != nil {
            slog.Error("Error wrtiting response: " + err.Error());
        }
    });
}

func routeRegister(service *appservice.AppService) http.HandlerFunc {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        defer r.Body.Close();

        if r.Method != "POST" {
            error := responses.NewErrorResponse("Only POST method allowed");
            http.Error(w, error.JsonString(), http.StatusMethodNotAllowed);
            return;
        }

        body, err := io.ReadAll(r.Body);
        if err != nil {
            error := responses.NewErrorResponse("Failed to read request body");
            http.Error(w, error.JsonString(), http.StatusBadRequest);
            return;
        }

        var data dto.CreateUserDto;
        if err := json.Unmarshal(body, &data); err != nil {
            error := responses.NewErrorResponse("Invalid payload. Expected JSON");
            http.Error(w, error.JsonString(), http.StatusBadRequest);
            return;
        }

        user, err := service.CreateUser(data);
        if err != nil {
            if errors.Is(err, appservice.DuplicateUserEmailError{}) ||
                errors.Is(err, appservice.DuplicateUserUsernameError{}) {
                error := responses.NewErrorResponse(err.Error());
                http.Error(w, error.JsonString(), http.StatusConflict);
                return;
            }

            slog.Error(err.Error());
            error := responses.NewErrorResponse("Failed to create user");
            http.Error(w, error.JsonString(), http.StatusInternalServerError);
            return;
        }

        response := responses.NewDataResponse("user", user);
        w.Header().Add("Content-Type", "application/json");
        w.Write(response.Json());
    });
}

func routeMe(service *appservice.AppService) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer r.Body.Close();

        auth := r.Context().Value("auth");

        data, _ := json.Marshal(auth);
        w.Write(data);
    });
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
    middleware.RecovererMiddleware,
    middleware.LoggerMiddleware,
};

var authMiddleware = append(utilMiddleware, middleware.JWTAutherMiddleware);
