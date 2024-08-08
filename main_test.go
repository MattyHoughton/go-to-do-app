package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateTaskHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/create", strings.NewReader("task=Test Task"))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("Expected status %v, got %v", http.StatusSeeOther, status)
	}

	if len(tasks) == 0 || tasks[len(tasks)-1].Task != "Test Task" {
		t.Error("Expected task 'Test Task' to be created")
	}
}

func TestReadTasksHandler(t *testing.T) {
	// Add a task to ensure we have data to read
	tasks = []Task{
		{Number: 1, Task: "Test Task", Status: "Not Started"},
	}

	req, err := http.NewRequest("GET", "/api/tasks", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(readTasksHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	var responseTasks []Task
	err = json.NewDecoder(rr.Body).Decode(&responseTasks)
	if err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	if len(responseTasks) == 0 || responseTasks[0].Task != "Test Task" {
		t.Error("Expected task 'Test Task' to be in the response")
	}
}

func TestUpdateTaskHandler(t *testing.T) {
	// Add a task to be updated
	tasks = []Task{
		{Number: 1, Task: "Old Task", Status: "Not Started"},
	}

	task := Task{
		Number: 1,
		Task:   "Updated Task",
		Status: "In Progress",
	}
	body, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	req, err := http.NewRequest("PUT", "/api/task", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(taskAPIHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	if len(tasks) == 0 || tasks[0].Task != "Updated Task" {
		t.Error("Expected task 'Updated Task' to be updated")
	}
}

func TestDeleteTaskHandler(t *testing.T) {
	// Add a task to be deleted
	tasks = []Task{
		{Number: 1, Task: "Task to Delete", Status: "Not Started"},
	}

	req, err := http.NewRequest("DELETE", "/api/task?id=1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(taskAPIHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Expected status %v, got %v", http.StatusNoContent, status)
	}

	if len(tasks) != 0 {
		t.Error("Expected tasks list to be empty after deletion")
	}
}

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	if !strings.Contains(rr.Body.String(), "To-Do List") {
		t.Error("Expected 'To-Do List' to be in the response body")
	}
}

func TestCreatePageHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/create", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	if !strings.Contains(rr.Body.String(), "Create Task") {
		t.Error("Expected 'Create Task' to be in the response body")
	}
}
