package model

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type Response struct {
	ID    *int64 `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type Tasks struct {
	ID string `json:"id,omitempty"`
	Task
}

type ResponseTasks struct {
	Tasks []Tasks `json:"tasks"`
	Error string  `json:"error,omitempty"`
}
