package main

import (
	"context"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/gorilla/mux"
	"github.com/phayes/freeport"
	"log"
	"net/http"
	"os/exec"
)

func setupPreferencesMenu(ctx context.Context, tokenCh chan string) *systray.MenuItem {
	port, err := freeport.GetFreePort()
	if err != nil {
		logError(err)
		return nil
	}
	go runConfigServer(ctx, port, tokenCh)

	serverUrl := fmt.Sprintf("http://localhost:%d", port)

	prefMenu := systray.AddMenuItem("Preferences", "")
	go func() {
		<-prefMenu.ClickedCh

		//systray.ShowAppWindow(serverUrl)

		err := exec.Command("open", serverUrl).Run()
		if err != nil {
			logError(err)
		}
	}()

	return prefMenu
}

func runConfigServer(ctx context.Context, port int, tokenCh chan string) {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<body>

<form action="/preferences">
  <label for="fname">Slack API token</label><br>
  <input type="text" id="token" name="token" value=""><br>
  <input type="submit" value="Submit">
</form> 
</body>
</html>s
`))
	})

	r.HandleFunc("/preferences", func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")
		log.Println("got slack token", token)
		tokenCh <- token
	})

	http.Handle("/", r)

	go func() {
		// TODO make this closable from tokenCh closing
		if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil); err != nil {
			panic(err)
		}
	}()
}