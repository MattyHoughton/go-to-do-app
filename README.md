# To-Do List Web Application

This project is a simple web-based To-Do list application implemented in Go. The application allows users to manage a list of tasks stored in memory, with features for creating, reading, updating, and deleting tasks. The UI is styled using CSS to provide a clean and modern user interface.

# Features

Create Task: Users can add a new task to the To-Do list.
Read Tasks: The home page displays a list of all tasks with their current status.
Update Task: Users can update the description and status of an existing task.
Delete Task: Users can remove a task from the list.
API: A basic REST API is provided for programmatic access to the tasks.

# Project Structure

The project directory is structured as follows:

main.go: The main Go application file containing the server logic, handlers, and task management logic.
go.mod: Go modules files for dependency management.
templates/: Directory containing HTML templates for the web pages.
create.html: Template for the task creation page.
index.html: Template for the home page displaying the list of tasks.
update.html: Template for the task update page.
static/css/: Directory containing the CSS file for styling the web pages.
style.css: The CSS file that styles all the web pages.
main_test.go: Contains unit tests for the application.

# Getting Started

# Prerequisites

Go: You need to have Go installed on your machine. You can download it from golang.org.

# Installation

Clone the repository:
git clone https://github.com/MattyHoughton/go-to-do-app
cd go-to-do-app

# Install dependencies:

go mod tidy

# Running the Application

Run the Go server:
go run main.go

Access the web app:
Open your web browser and go to http://localhost:8080.

# Usage

Home Page: Displays a list of all tasks with options to update or delete each task.
Create Page: Accessible via the "Create New Task" button on the home page. Allows adding a new task.
Update Page: Accessible via the "Update" button next to each task. Allows updating the task's description and status.

# API Endpoints

GET /api/tasks: Retrieve the list of tasks in JSON format.
POST /api/task: Create a new task (requires a JSON body).
PUT /api/task: Update an existing task (requires a JSON body).
DELETE /api/task?id={taskID}: Delete a task by its ID.

# CSS Styling

The application uses a simple and consistent CSS file (style.css) to style all pages.
Buttons, tables, and forms are styled to provide a clean and user-friendly interface.

# Testing

To run tests, navigate to the project directory and run:

go test -v
