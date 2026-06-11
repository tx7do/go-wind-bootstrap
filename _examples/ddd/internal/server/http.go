// Package server 实现接口层（HTTP handler）。
//
// server 层负责 HTTP 请求/响应的编解码和路由，将请求委托给 service 层处理。
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/tx7do/go-wind-bootstrap/_examples/ddd/internal/domain"
	"github.com/tx7do/go-wind-bootstrap/_examples/ddd/internal/service"
	httpAdapter "github.com/tx7do/go-wind-bootstrap/transport/http"
)

// RegisterHTTPRoutes 注册用户相关的 HTTP 路由。
func RegisterHTTPRoutes(srv *httpAdapter.Server, svc *service.UserService) {
	srv.POST("/users", handleCreate(svc))
	srv.GET("/users", handleList(svc))
	srv.GET("/users/{id}", handleGet(svc))
	srv.PUT("/users/{id}", handleUpdate(svc))
	srv.DELETE("/users/{id}", handleDelete(svc))
}

// --- Handlers ---

func handleCreate(svc *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		user, err := svc.CreateUser(&req)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, user)
	}
}

func handleGet(svc *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		user, err := svc.GetUser(id)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, user)
	}
}

func handleList(svc *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		users, err := svc.ListUsers(offset, limit)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, users)
	}
}

func handleUpdate(svc *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		var req domain.UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		user, err := svc.UpdateUser(id, &req)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, user)
	}
}

func handleDelete(svc *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := pathID(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := svc.DeleteUser(id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// --- Helpers ---

func pathID(r *http.Request) (uint64, error) {
	parts := strings.Split(r.URL.Path, "/")
	raw := parts[len(parts)-1]
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %s", raw)
	}
	return id, nil
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
