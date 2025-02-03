package internalhttp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"

	"github.com/google/uuid"
)

func (s *Server) main(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("start processing", "request", r.URL.Path)

	_, err := w.Write([]byte("hello world"))
	if err != nil {
		s.log.Error("failed to write response", "error", err)
		return
	}

	s.log.Debug("successfully finished processing", "request", r.URL.Path)
}

//nolint:dupl
func (s *Server) add(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("start processing add", "request", r.URL.Path)

	var requestEvent model.Event

	if err := json.NewDecoder(r.Body).Decode(&requestEvent); err != nil {
		s.log.Error("failed to json.NewDecoder(r.Body).Decode", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := s.service.Add(r.Context(), requestEvent)
	if err != nil {
		s.log.Error("failed to service.Add", "error", err)
		http.Error(w, "failed to add event", http.StatusInternalServerError)

		return
	}

	_, err = w.Write([]byte("event added successfully"))
	if err != nil {
		s.log.Error("failed to write response", "error", err)
		return
	}

	s.log.Debug("successfully finished processing", "request", r.URL.Path)
}

//nolint:dupl
func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("start processing add", "request", r.URL.Path)

	var requestEvent model.Event

	if err := json.NewDecoder(r.Body).Decode(&requestEvent); err != nil {
		s.log.Error("failed to json.NewDecoder(r.Body).Decode", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := s.service.Update(r.Context(), requestEvent)
	if err != nil {
		s.log.Error("failed to service.Update", "error", err)
		http.Error(w, "failed to add event", http.StatusInternalServerError)

		return
	}

	_, err = w.Write([]byte("event updated successfully"))
	if err != nil {
		s.log.Error("failed to write response", "error", err)
		return
	}

	s.log.Debug("successfully finished processing", "request", r.URL.Path)
}

func (s *Server) delete(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("start processing add", "request", r.URL.Path)

	var id uuid.UUID

	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		s.log.Error("failed to json.NewDecoder(r.Body).Decode", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := s.service.Delete(r.Context(), model.EventUUID(id))
	if err != nil {
		s.log.Error("failed to service.Delete", "error", err)
		http.Error(w, "failed to delete event", http.StatusInternalServerError)

		return
	}

	_, err = w.Write([]byte("event deleted successfully"))
	if err != nil {
		s.log.Error("failed to write response", "error", err)
		return
	}

	s.log.Debug("successfully finished processing", "request", r.URL.Path)
}

//nolint:dupl
func (s *Server) getForDay(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("start processing add", "request", r.URL.Path)

	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		s.log.Error("failed to time.Parse", "dateStr", dateStr)
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	events, err := s.service.GetForDay(r.Context(), date)
	if err != nil {
		s.log.Error("failed to service.GetForDay", "error", err)
		http.Error(w, "failed to get day events", http.StatusInternalServerError)

		return
	}

	marshalledData, err := json.Marshal(events)
	if err != nil {
		s.log.Error("failed to json.Marshal", "error", err)
		http.Error(w, "failed to process events", http.StatusInternalServerError)
	}

	_, err = w.Write(marshalledData)
	if err != nil {
		s.log.Error("failed to write response", "error", err)
		return
	}

	s.log.Debug("successfully finished processing", "request", r.URL.Path)
}

//nolint:dupl
func (s *Server) getForWeek(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("start processing add", "request", r.URL.Path)

	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		s.log.Error("failed to time.Parse", "dateStr", dateStr)
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	events, err := s.service.GetForWeek(r.Context(), date)
	if err != nil {
		s.log.Error("failed to service.GetForWeek", "error", err)
		http.Error(w, "failed to get week events", http.StatusInternalServerError)

		return
	}

	marshalledData, err := json.Marshal(events)
	if err != nil {
		s.log.Error("failed to json.Marshal", "error", err)
		http.Error(w, "failed to process events", http.StatusInternalServerError)
	}

	_, err = w.Write(marshalledData)
	if err != nil {
		s.log.Error("failed to write response", "error", err)
		return
	}

	s.log.Debug("successfully finished processing", "request", r.URL.Path)
}

//nolint:dupl
func (s *Server) getForMonth(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("start processing add", "request", r.URL.Path)

	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		s.log.Error("failed to time.Parse", "dateStr", dateStr)
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	events, err := s.service.GetForMonth(r.Context(), date)
	if err != nil {
		s.log.Error("failed to service.GetForMonth", "error", err)
		http.Error(w, "failed to get month events", http.StatusInternalServerError)

		return
	}

	marshalledData, err := json.Marshal(events)
	if err != nil {
		s.log.Error("failed to json.Marshal", "error", err)
		http.Error(w, "failed to process events", http.StatusInternalServerError)
	}

	_, err = w.Write(marshalledData)
	if err != nil {
		s.log.Error("failed to write response", "error", err)
		return
	}

	s.log.Debug("successfully finished processing", "request", r.URL.Path)
}
