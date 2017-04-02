package scheduler

import (
	"net/http"
	"time"

	"github.com/caarlos0/watchub/shared/pages"
)

// ScheduleHandler schedules a new check for the logged user
func (s *Scheduler) ScheduleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.session.Get(r, s.config.SessionName)
		id, _ := session.Values["user_id"].(int)
		login, _ := session.Values["user_login"].(string)
		if session.IsNew || id == 0 {
			http.Error(w, "invalid session", http.StatusForbidden)
			return
		}
		if err := s.store.Schedule(int64(id), time.Now()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pages.Render(w, "checking", pages.PageData{
			User: pages.User{
				ID:    id,
				Login: login,
			},
			ClientID: s.config.ClientID,
		})
	}
}
