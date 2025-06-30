package task

import (
	Models "ThreeLayerArch/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type TaskServices interface {
	Add_Task(task string) (bool, error)
	View_Task() ([]Models.Tasks, error)
	Get_By_ID(i int) (Models.Tasks, error)
	Update_Task(i int) (bool, error)
	Delete_Task(i int) (bool, error)
}
type Handler struct {
	service TaskServices
}

// New creates a new task handler
func New(service TaskServices) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Addtask(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	msg, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	var reqBody struct {
		T string `json:"task"`
	}

	err = json.Unmarshal(msg, &reqBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	ans, err := h.service.Add_Task(reqBody.T)

	if err != nil {
		log.Printf("Error in HANDLER.AddTask: %v", err)
		log.Printf("%s", err.Error())
		return
	}

	if ans {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Task added"))
		return
	}
}

func (h *Handler) Viewtask(w http.ResponseWriter, r *http.Request) {
	ans, err := h.service.View_Task()
	if err != nil {
		log.Printf("Error in HANDLER.Viewtask: %v", err)
		return
	}
	for _, v := range ans {
		fmt.Fprintf(w, "ID: %d, Task: %s, Completed: %t\n", v.Tid, v.Task, v.Completed)
	}
}

func (h *Handler) Gettask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	var ans Models.Tasks

	ans, err = h.service.Get_By_ID(index)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ID: %d, Task: %s, Completed: %t\n", ans.Tid, ans.Task, ans.Completed)
}

func (h *Handler) Updatetask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	var ans bool
	ans, err = h.service.Update_Task(index)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	if ans {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Task updated"))
	}
}

func (h *Handler) Deletetask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	var ans bool
	ans, err = h.service.Delete_Task(index)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	if ans {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Task deleted"))

	}
}

// added comment for creating pr
