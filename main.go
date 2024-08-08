package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"text/template"
)

type Task struct {
	Number int    `json:"number"`
	Task   string `json:"task"`
	Status string `json:"status"`
}

var (
	tasks     []Task
	nextID    int
	mu        sync.Mutex
	templates *template.Template
)

func init() {
	nextID = 1
	templates = template.Must(template.New("").ParseGlob("templates/*.html"))
}

// Handler for the home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	tasks = renumberTasks(tasks)
	templates.ExecuteTemplate(w, "index.html", tasks)
}

// Handler for creating tasks
func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "create.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		taskDesc := r.FormValue("task")
		status := "Not Started"

		mu.Lock()
		task := Task{
			Number: nextID,
			Task:   taskDesc,
			Status: status,
		}
		nextID++
		tasks = append(tasks, task)
		mu.Unlock()

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// Handler for updating tasks
func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		idStr := r.FormValue("id")
		taskDesc := r.FormValue("task")
		status := r.FormValue("status")

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid task number", http.StatusBadRequest)
			return
		}

		mu.Lock()
		for i, task := range tasks {
			if task.Number == id {
				tasks[i].Task = taskDesc
				tasks[i].Status = status
				break
			}
		}
		mu.Unlock()

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid task number", http.StatusBadRequest)
		return
	}

	mu.Lock()
	var task *Task
	for _, t := range tasks {
		if t.Number == id {
			task = &t
			break
		}
	}
	mu.Unlock()

	if task == nil {
		http.NotFound(w, r)
		return
	}

	templates.ExecuteTemplate(w, "update.html", task)
}

// Handler for deleting tasks
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task number", http.StatusBadRequest)
		return
	}

	mu.Lock()
	for i, task := range tasks {
		if task.Number == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}
	tasks = renumberTasks(tasks)
	mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// API Handlers
func taskAPIHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var task Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		task.Number = nextID
		nextID++
		task.Status = "Not Started"

		mu.Lock()
		tasks = append(tasks, task)
		mu.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)

	case http.MethodPut:
		var updatedTask Task
		err := json.NewDecoder(r.Body).Decode(&updatedTask)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		mu.Lock()
		for i, task := range tasks {
			if task.Number == updatedTask.Number {
				tasks[i].Task = updatedTask.Task
				tasks[i].Status = updatedTask.Status
				mu.Unlock()
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(tasks[i])
				return
			}
		}
		mu.Unlock()

		http.NotFound(w, r)

	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid task number", http.StatusBadRequest)
			return
		}

		mu.Lock()
		for i, task := range tasks {
			if task.Number == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				tasks = renumberTasks(tasks)
				mu.Unlock()
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		mu.Unlock()

		http.NotFound(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handler to read all tasks via API
func readTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		mu.Lock()
		defer mu.Unlock()

		tasks = renumberTasks(tasks)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// Renumber tasks to keep the numbering sequential
func renumberTasks(tasks []Task) []Task {
	for i := range tasks {
		tasks[i].Number = i + 1
	}
	return tasks
}

// Set up web UI and API handlers
func setupHandlers() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/api/tasks", readTasksHandler)
	http.HandleFunc("/api/task", taskAPIHandler)

	// Serve static files (CSS, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func main() {
	setupHandlers()

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
