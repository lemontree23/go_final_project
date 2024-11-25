package handlers

import (
	"fmt"
	"net/http"
	"scheduler/internal/scheduler"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	if nowStr == "" || dateStr == "" || repeatStr == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "invalid 'now' parameter", http.StatusBadRequest)
		return
	}

	nextDate, err := scheduler.NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, nextDate)
}
