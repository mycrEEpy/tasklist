package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rs/xid"

	"github.com/gorilla/mux"
)

const (
	persistFile = "tasks.json"
)

var (
	listenAddr    = flag.String("listen", ":8080", "listen address")
	persistFolder = flag.String("persist", "/tmp/tasklist", "folder to persist tasks")
	tasks         = make([]task, 0, 4)
)

type task struct {
	ID    string `json:"id"`
	What  string `json:"what"`
	Where string `json:"where"`
	Who   string `json:"who"`
}

func main() {
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/tasks/save", saveTasks).Methods(http.MethodGet)
	r.HandleFunc("/tasks/load", loadTasks).Methods(http.MethodGet)
	r.HandleFunc("/tasks", getTasks).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{id}", getSingleTask).Methods(http.MethodGet)
	r.HandleFunc("/tasks", createTask).Methods(http.MethodPost)
	r.HandleFunc("/tasks/{id}", deleteTask).Methods(http.MethodDelete)
	fmt.Printf("Listening on %s\n", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, r); err != nil {
		fmt.Printf("ListenAndServe error: %s\n", err)
	}
}

func getTasks(res http.ResponseWriter, req *http.Request) {
	jsonBytes, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Fprint(res, err)
		return
	}
	fmt.Fprintf(res, "Tasks: %s", jsonBytes)
}

func getSingleTask(res http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	for _, t := range tasks {
		if t.ID != v["id"] {
			continue
		}
		jsonBytes, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			fmt.Fprint(res, err)
			return
		}
		fmt.Fprintf(res, "Task: %s", jsonBytes)
		return
	}
	fmt.Fprintf(res, "No task found with id %s", v["id"])
}

func createTask(res http.ResponseWriter, req *http.Request) {
	var t task
	d := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err := d.Decode(&t)
	if err != nil {
		fmt.Fprint(res, err)
		return
	}
	if t.ID == "" {
		t.ID = xid.New().String()
	}
	tasks = append(tasks, t)
	fmt.Fprintf(res, "Created task")
}

func deleteTask(res http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	for i, t := range tasks {
		if t.ID != v["id"] {
			continue
		}
		tasks = append(tasks[:i], tasks[i+1:]...)
		fmt.Fprint(res, "Deleted task")
		return
	}
	fmt.Fprintf(res, "No task found with id %s", v["id"])
}

func saveTasks(res http.ResponseWriter, req *http.Request) {
	jsonBytes, err := json.Marshal(tasks)
	if err != nil {
		fmt.Fprint(res, err)
		return
	}
	if _, err := os.Stat(*persistFolder); os.IsNotExist(err) {
		os.Mkdir(*persistFolder, 0755)
	}
	err = ioutil.WriteFile(filepath.Join(*persistFolder, persistFile), jsonBytes, 0644)
	if err != nil {
		fmt.Fprint(res, err)
		return
	}
	fmt.Fprint(res, "Persisted tasks")
}

func loadTasks(res http.ResponseWriter, req *http.Request) {
	jsonBytes, err := ioutil.ReadFile(filepath.Join(*persistFolder, persistFile))
	if err != nil {
		fmt.Fprint(res, err)
		return
	}
	err = json.Unmarshal(jsonBytes, &tasks)
	if err != nil {
		fmt.Fprint(res, err)
		return
	}
	fmt.Fprint(res, "Loaded tasks")
}
