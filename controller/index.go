package controller

import (
	"github.com/ChrisTheBaron/strava-ical/entities"
	"github.com/ChrisTheBaron/strava-ical/utils"
	"github.com/golang/glog"
	"net/http"
)

type Index struct {
	controller
	tr *utils.Template
}

// NewIndex returns a new controller.Index
func NewIndex(con *entities.Config, tr *utils.Template) *Index {
	return &Index{controller{config: con}, tr}
}

func (i *Index) Get(w http.ResponseWriter, r *http.Request) {

	glog.Infoln("#controller.Index.Get")

	data := struct{}{}

	err := i.tr.Render(w, data, "index.tmpl")

	if err != nil {
		glog.Error(err)
	}

}
