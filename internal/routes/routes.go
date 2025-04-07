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
	auth_helpers "github.com/cxcnxl/go-crud/internal/auth_helpers"
	"github.com/cxcnxl/go-crud/internal/dto"
	"github.com/cxcnxl/go-crud/internal/middleware"
	"github.com/cxcnxl/go-crud/internal/responses"
	"github.com/cxcnxl/go-crud/internal/redis"
)

func NewRouter(db *gorm.DB, redis *redis.RedisWrapper) *MethodHandler {
    mux := http.NewServeMux();
    service := appservice.NewAppService(db, redis);
    methodHandler := NewMethodHandler(mux);

    methodHandler.HandleFunc(
        "GET",
        "/",
        routeIndex(service),
        middleware.UtilMiddleware,
    );
    methodHandler.HandleFunc(
        "POST",
        "/register",
        routeRegister(service),
        middleware.UtilMiddleware,
    );
    methodHandler.HandleFunc(
        "POST",
        "/login",
        routeLogin(service),
        middleware.UtilMiddleware,
    );
    methodHandler.HandleFunc(
        "GET",
        "/me",
        routeMe(service),
        middleware.AuthMiddleware,
    );

    return methodHandler;
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

func routeLogin(service *appservice.AppService) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer r.Body.Close();

        body, err := io.ReadAll(r.Body);
        if err != nil {
            error := responses.NewErrorResponse("Failed to read request body");
            http.Error(w, error.JsonString(), http.StatusBadRequest);
            return;
        }
        
        var data dto.PostLoginDto;
        if err := json.Unmarshal(body, &data); err != nil {
            error := responses.NewErrorResponse("Invalid body");
            http.Error(w, error.JsonString(), http.StatusBadRequest);
            return;
        }

        user, err := service.LoginUser(data);
        if err != nil {
            error := responses.NewErrorResponse(err.Error());
            http.Error(w, error.JsonString(), http.StatusUnauthorized);
            return;
        }

        jwtPayload := map[string]any {
            "id": user.ID,
            "email": user.Email,
            "username": user.Username,
        };

        jwt := auth_helpers.SignJWT(jwtPayload);

        res := responses.NewDataResponse("auth", map[string]string{
            "auth_token": jwt,
        });
        w.Header().Add("Content-Type", "application/json");
        w.WriteHeader(http.StatusCreated);
        w.Write(res.Json());
    });
}

func routeMe(_ *appservice.AppService) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer r.Body.Close();

        auth := r.Context().Value("auth");

        data, _ := json.Marshal(auth);
        w.Write(data);
    });
}

// ---------- Utils -----------

type MethodHandler struct {
    Mux *http.ServeMux
    methods map[string]middleware.MiddlewareSet
}

func NewMethodHandler(mux *http.ServeMux) *MethodHandler {
    return &MethodHandler{
        mux,
        make(map[string]middleware.MiddlewareSet),
    };
}

func (self *MethodHandler) HandleFunc(
    method string,
    path string,
    handler http.HandlerFunc,
    middlewares middleware.MiddlewareSet,
) {
    self.methods[method+path] = middlewares;

    wrapped := middleware.Wrap(handler, middlewares);

    self.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
        if r.Method == method {
            wrapped(w, r);
        } else {
            // TODO: this avoides logger and restore middlewares
            error := responses.NewErrorResponse("Method not allowed");
            w.Header().Add("Content-Type", "application/json");
            w.WriteHeader(http.StatusMethodNotAllowed);
            w.Write(error.Json());
        }
    });
}
