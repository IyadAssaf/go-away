package main

import (
	"context"
	"fmt"
	"github.com/IyadAssaf/go-away/internal/status"
	"github.com/getlantern/systray"
	"github.com/gorilla/mux"
	"github.com/phayes/freeport"
	"html/template"
	"net/http"
	"os/exec"
)

func setupPreferencesMenu(ctx context.Context, s *status.SlackStatus, triggerCh chan struct{}) *systray.MenuItem {
	port, err := freeport.GetFreePort()
	if err != nil {
		logError(err)
		return nil
	}
	go runConfigServer(ctx, port, s, triggerCh)

	serverUrl := fmt.Sprintf("http://localhost:%d", port)

	prefMenu := systray.AddMenuItem("Preferences", "")
	go func() {
		for {
			<-prefMenu.ClickedCh

			err := exec.Command("open", serverUrl).Run()
			if err != nil {
				logError(err)
			}
		}
	}()

	return prefMenu
}

func runConfigServer(ctx context.Context, port int, s *status.SlackStatus, triggerCh chan struct{}) {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "text/html")

		tpl, err := template.New("preferences").Parse(`
<!doctype html>
<html lang="en">
<head>
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
</head>
<body>
<div class="container">
	<div class="row">
    	<img class="center-block" src="https://github.com/IyadAssaf/go-away/blob/master/assets/logo/1100-with_padding.png?raw=true"/>
	</div>
	<form action="/preferences" method="POST" onsubmit="window.close()">
		<div class="form-group row">
			<label class="col-sm-2 col-form-label" for="fname">Slack API token</label>
			<div class="col-sm-10">
				<input type="password" class="form-control" id="token" name="token" value="{{.Token}}">
			</div>
		</div>
		<div class="form-group row">
			<label class="col-sm-2 col-form-label" for="fname">Emoji</label>
			<div class="col-sm-10">
				<input type="text" class="form-control" id="emoji" name="emoji" value="{{.Emoji}}">
			</div>
		</div>
		<div class="form-group row">
			<label class="col-sm-2 col-form-label" for="fname">Message</label>
			<div class="col-sm-10">
				<input type="text" class="form-control" id="message" name="message" value="{{.Message}}">
			</div>
		</div>
		<button type="submit" class="btn btn-primary">Save</button>
	</form>
</div>
</body>
</html>
`)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		settings := struct {
			Message string
			Emoji   string
			Token   string
		}{
			Message: settingMessage,
			Emoji:   settingEmoji,
			Token:   settingToken,
		}

		w.Header().Set("Content-Type", "text/html")
		tpl.Execute(w, settings)
	})

	r.HandleFunc("/preferences", func(w http.ResponseWriter, r *http.Request) {
		settingToken = r.FormValue("token")
		settingMessage = r.FormValue("message")
		settingEmoji = r.FormValue("emoji")

		s.WithStatusText(settingMessage)
		s.WithStatusEmoji(settingEmoji)
		s.WithSlackToken(settingToken)
		statusMenu.SetTitle(fmt.Sprintf("Status set to %s %s", settingEmoji, settingMessage))

		triggerCh <- struct{}{}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	http.Handle("/", r)

	// TODO make this closable from tokenCh closing
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil); err != nil {
		panic(err)
	}
}
