package http_api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"task-tracker-service/internal/config"
	"task-tracker-service/internal/controller/http_api/dto"
	"task-tracker-service/internal/domain/service"
)

type Handler struct {
	service *service.DayToDoService
	config  *config.Config
	logger  *logrus.Logger
	router  *mux.Router
}

func NewHandler(service *service.DayToDoService, config *config.Config, logger *logrus.Logger, router *mux.Router) *Handler {
	return &Handler{
		service: service,
		config:  config,
		logger:  logger,
		router:  router,
	}
}

func (h *Handler) InitRouts() error {
	h.router.HandleFunc("/tasks", h.addTask).Methods("POST")
	h.router.HandleFunc("/tasks/{date}/{task}", h.deleteTask).Methods("DELETE")
	h.router.HandleFunc("/tasks/{date}", h.deleteDateTasks).Methods("DELETE")
	h.router.HandleFunc("/tasks/{date}", h.getDateTasks).Methods("GET")
	h.router.HandleFunc("/tasks", h.getAllTasks).Methods("GET")
	h.router.HandleFunc("/tasks/{date}/{task}", h.updateTaskByDate).Methods("PUT")
	return nil
}

func (h *Handler) initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func (h *Handler) addTask(writer http.ResponseWriter, request *http.Request) {
	h.initHeaders(writer)
	h.logger.Info("Post Task POST /tasks")
	var dateTask dto.DayToDo
	err := json.NewDecoder(request.Body).Decode(&dateTask)
	if err != nil {
		h.logger.Info("Invalid json received from client")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if !h.validDate(dateTask.Date, writer) {
		return
	}
	err = h.service.AddTask(dateTask.Date, dateTask.Task)
	if err != nil {
		h.logger.Info("Troubles while creating new task:", err)
		writer.WriteHeader(http.StatusNotImplemented)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

func (h *Handler) deleteTask(writer http.ResponseWriter, request *http.Request) {
	h.initHeaders(writer)
	h.logger.Info("Delete task DELETE /tasks/{date}/{task}")
	var dateTask dto.DayToDo
	var isParse bool
	dateTask.Date, isParse = mux.Vars(request)["date"]
	if !isParse {
		h.logger.Info("Invalid json received from client")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	dateTask.Task, isParse = mux.Vars(request)["task"]
	if !isParse {
		h.logger.Info("Invalid json received from client")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if !h.validDate(dateTask.Date, writer) {
		return
	}
	err := h.service.DeleteTask(dateTask.Date, dateTask.Task)
	if err != nil {
		h.logger.Info("Troubles with delete task by date:", err)
		writer.WriteHeader(http.StatusNotImplemented)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteDateTasks(writer http.ResponseWriter, request *http.Request) {
	h.initHeaders(writer)
	h.logger.Info("Delete all tasks by date DELETE /tasks/{date}")
	date, isParse := mux.Vars(request)["date"]
	if !isParse {
		h.logger.Info("Invalid json received from client")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if !h.validDate(date, writer) {
		return
	}
	_, err := h.service.DeleteAllTasks(date)
	if err != nil {
		h.logger.Info("Troubles with delete all tasks by date:", err)
		writer.WriteHeader(http.StatusNotImplemented)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getDateTasks(writer http.ResponseWriter, request *http.Request) {
	h.initHeaders(writer)
	h.logger.Info("Get all tasks by date GET /tasks/{date}")
	date, isParse := mux.Vars(request)["date"]
	if !isParse {
		h.logger.Info("Invalid json received from client")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if !h.validDate(date, writer) {
		return
	}
	taskDate, err := h.service.FindAllTaskDate(date)
	if err != nil {
		h.logger.Info("Troubles with find tasks by date:", err)
		writer.WriteHeader(http.StatusNotImplemented)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(taskDate)
}

func (h *Handler) getAllTasks(writer http.ResponseWriter, request *http.Request) {
	h.initHeaders(writer)
	h.logger.Info("Get all tasks GET /tasks")

	tasks, err := h.service.FindAllTasks()
	if err != nil {
		h.logger.Info("Troubles with find tasks", err)
		writer.WriteHeader(http.StatusNotImplemented)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(tasks)
}

func (h *Handler) updateTaskByDate(writer http.ResponseWriter, request *http.Request) {
	h.initHeaders(writer)
	h.logger.Info("Update tasks by date PUT /tasks/{date}/{task}")
	var dayTask dto.DayToDo
	var isParse bool
	dayTask.Date, isParse = mux.Vars(request)["date"]
	if !isParse {
		h.logger.Info("Invalid json received from client")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	dayTask.Task, isParse = mux.Vars(request)["task"]
	if !isParse {
		h.logger.Info("Invalid json received from client")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if !h.validDate(dayTask.Date, writer) {
		return
	}
	var newTask dto.DayToDo
	err := json.NewDecoder(request.Body).Decode(&newTask)
	if err != nil {
		h.logger.Info("Invalid json received from client")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.service.UpdateTaskByDate(dayTask.Date, dayTask.Task, newTask.Task)
	if err != nil {
		h.logger.Info("Troubles with update task", err)
		writer.WriteHeader(http.StatusNotImplemented)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

func (h *Handler) validDate(date string, writer http.ResponseWriter) bool {
	isValidDate, err := regexp.MatchString(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`, date)
	if err != nil {
		h.logger.Info("Invalid json received from client")
		writer.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !isValidDate {
		h.logger.Info("Invalid date format received from client")
		writer.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}
