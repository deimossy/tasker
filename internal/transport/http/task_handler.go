package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/deimossy/tasker/internal/entity"
	"github.com/deimossy/tasker/internal/service"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/tasks":
		h.listTasks(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/tasks":
		h.createTask(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/tasks/"):
		h.getTask(w, r)
	default:
		http.NotFound(w, r)
	}
}

// ! GET /tasks?status=...
func (h *TaskHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	statusParam := r.URL.Query().Get("status")
	var status *entity.Status
	if statusParam != "" {
		s := entity.Status(statusParam)
		if !s.IsValid() {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}
		status = &s
	}

	tasks, err := h.svc.ListTasks(ctx, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, tasks)
}

// ! GET /tasks/{id}
func (h *TaskHandler) getTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := h.svc.GetTaskByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, task)
}

// ! POST /tasks
func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var task entity.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	created, err := h.svc.CreateTask(ctx, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, created)
}

// ! writeJSON сериализует ответ в JSON
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}