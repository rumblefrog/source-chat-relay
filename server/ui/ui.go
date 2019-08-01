package ui

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/rumblefrog/source-chat-relay/server/config"
	"github.com/rumblefrog/source-chat-relay/server/relay"

	"github.com/rumblefrog/source-chat-relay/server/entity"

	packr "github.com/gobuffalo/packr/v2"

	"github.com/sirupsen/logrus"
)

var box *packr.Box

type renderData struct {
	Relay    *relay.Relay
	Entities []*renderEntity
}

type renderEntity struct {
	Entity      *entity.Entity
	Highlighted bool
}

func UIListen() {
	box = packr.New("template", "./template/dist")

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/styles.css", styleHandler)

	if config.Config.UI.Port == 0 {
		config.Config.UI.Port = 8080
	}

	logrus.Infof("UI listener started on port %d", config.Config.UI.Port)

	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Config.UI.Port), nil))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	s, err := box.FindString("index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	temp, err := template.New("index").Funcs(template.FuncMap{
		"humanizeChannelString": entity.HumanizeChannelString,
		"byteToMB": func(b int) string {
			return fmt.Sprintf("%.6f", float64(b)/(1024*1024))
		},
	}).Parse(s)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "text/html")

	entities := entity.GetEntities()

	tRenderData := renderData{
		Relay: relay.Instance,
	}

	for _, v := range entities {
		tRenderData.Entities = append(tRenderData.Entities, &renderEntity{
			Entity: v,
		})
	}

	if r.Method == http.MethodPost {
		action := r.FormValue("btn")

		switch action {
		case "update":
			// Since we already retrieved from cache prior (above and on bot ready), this entity should point to the same address
			// as tRenderData and relay instance entities, therefore the template render should reflect the update
			tEntity, err := entity.GetEntity(r.FormValue("id"))

			if err != nil {
				break
			}

			tEntity.ReceiveChannels = entity.ParseDelimitedChannels(r.FormValue("receiveChannels"))
			tEntity.SendChannels = entity.ParseDelimitedChannels(r.FormValue("sendChannels"))
			tEntity.DisabledReceiveTypes = entity.ParseDelimitedChannels(r.FormValue("disabledReceiveTypes"))
			tEntity.DisabledSendTypes = entity.ParseDelimitedChannels(r.FormValue("disabledSendTypes"))

			tEntity.Propagate()
		case "trace":
			for _, v := range tRenderData.Entities {
				sendChannels := r.FormValue("sendChannels")
				receiveChannels := r.FormValue("receiveChannels")

				if len(sendChannels) != 0 && v.Entity.ReceiveIntersectsWith(entity.ParseDelimitedChannels(sendChannels)) {
					v.Highlighted = true
				}

				if len(receiveChannels) != 0 && v.Entity.SendIntersectsWith(entity.ParseDelimitedChannels(receiveChannels)) {
					v.Highlighted = true
				}
			}
		}
	}

	err = temp.Execute(w, tRenderData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func styleHandler(w http.ResponseWriter, r *http.Request) {
	s, err := box.Find("styles.css")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	buf := bytes.NewBuffer(s)

	w.Header().Set("Content-Type", "text/css; charset=utf-8")

	io.Copy(w, buf)
}
