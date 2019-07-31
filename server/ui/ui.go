package ui

import (
	"bytes"
	"html/template"
	"io"
	"net/http"

	"github.com/rumblefrog/source-chat-relay/server/entity"

	packr "github.com/gobuffalo/packr/v2"

	"github.com/sirupsen/logrus"
)

var box *packr.Box

func UIListen() {
	box = packr.New("template", "./template/dist")

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/styles.css", styleHandler)

	logrus.Info("UI listener started")

	logrus.Fatal(http.ListenAndServe(":8080", nil))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	s, err := box.FindString("index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	temp, err := template.New("index").Funcs(template.FuncMap{
		"humanizeChannelString": entity.HumanizeChannelString,
	}).Parse(s)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "text/html")

	err = temp.Execute(w, entity.GetEntities())

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
