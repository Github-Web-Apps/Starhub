package controllers

import (
	"net/http"

	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/datastore"
	"github.com/caarlos0/watchub/shared/dto"
	"github.com/caarlos0/watchub/shared/pages"
	"github.com/gorilla/sessions"
)

// Index ctrl
type Index struct {
	Base
	store datastore.Datastore
}

// NewIndex ctrl
func NewIndex(
	config config.Config,
	session sessions.Store,
	store datastore.Datastore,
) *Index {
	return &Index{
		Base: Base{
			config:  config,
			session: session,
		},
		store: store,
	}
}

// Handler handles /
func (ctrl *Index) Handler(w http.ResponseWriter, r *http.Request) {
	var data = dto.IndexPageData{
		PageData: ctrl.sessionData(w, r),
	}
	if data.User.ID > 0 {
		var err error
		var id = int64(data.User.ID)
		data.StarCount, err = ctrl.store.StarCount(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.FollowerCount, err = ctrl.store.FollowerCount(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.RepositoryCount, err = ctrl.store.RepositoryCount(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	pages.Render(w, "index", data)
}
