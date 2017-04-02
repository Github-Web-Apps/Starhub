package scheduler

import (
	"net/http"
	"time"
)

// ScheduleHandler schedules a new check for the logged user
func (s *Scheduler) ScheduleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.session.Get(r, s.config.SessionName)
		id, _ := session.Values["user_id"].(int)
		if session.IsNew || id == 0 {
			http.Error(w, "not logged in", http.StatusForbidden)
			return
		}
		if err := s.store.Schedule(int64(id), time.Now()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/scheduled", http.StatusTemporaryRedirect)
	}
}
