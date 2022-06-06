package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/octo/watersort"
)

type server struct {
	tmpl *template.Template
}

func newServer() *server {
	t, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	return &server{
		tmpl: t,
	}
}

func stateURL(s watersort.State) (string, error) {
	stateParam, err := s.MarshalText()
	if err != nil {
		return "", err
	}

	values := make(url.Values)
	values.Set("state", string(stateParam))

	return "/state?" + values.Encode(), nil
}

func (s server) GenerateStateHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	state := watersort.RandomState(10, 4)

	url, err := stateURL(state)
	if err != nil {
		return err
	}

	http.Redirect(w, req, url, http.StatusFound)
	return nil
}

func (s server) StateHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	stateParam := req.FormValue("state")
	if stateParam == "" {
		return httpError{
			msg:  "the required 'state' parameter is missing",
			code: http.StatusBadRequest,
		}
	}

	var state watersort.State
	if err := state.UnmarshalText([]byte(stateParam)); err != nil {
		return httpError{
			msg:  "failed to parse the 'state' parameter",
			code: http.StatusBadRequest,
		}

	}

	solved := state.Solved()

	var (
		step      watersort.Step
		nextState watersort.State
	)
	if !solved {
		steps, err := watersort.FindSolution(state)
		if err != nil {
			return err
		}
		step = steps[0]

		nextState = state.Clone()
		if err := nextState.Pour(step.From, step.To); err != nil {
			return err
		}
	}

	nextURL, err := stateURL(nextState)
	if err != nil {
		return err
	}

	data := struct {
		State   watersort.State
		Step    watersort.Step
		NextURL string
		Solved  bool
	}{
		State:   state.Clone(),
		Step:    step,
		NextURL: nextURL,
		Solved:  solved,
	}

	return s.tmpl.ExecuteTemplate(w, "state_show.html", data)
}
