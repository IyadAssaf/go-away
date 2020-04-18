package main

import (
	"context"
	"fmt"
	"github.com/IyadAssaf/go-away/internal/status"
	"os/exec"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/zserge/webview"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"github.com/phayes/freeport"
)


var slackTokenOption string

func main() {
	systray.RunWithAppWindow("Preferences", 500, 500, onReady, onExit)
}

func onReady() {
	var err error
	systray.SetTooltip("Go Away")
	systray.SetIcon(cameraOffIconData)

	ctx, cancel := context.WithCancel(context.Background())
	_ = ctx

	statusCameraOnText := systray.AddMenuItem(fmt.Sprintf(`Slack status set to "%s %s"`, status.DefaultStatusText, status.DefaultStatusEmoji), "")
	statusCameraOnText.Hide()
	statusCameraOnText.Disable()

	port, err := freeport.GetFreePort()
	if err != nil {
		//TODO do something better
		log.Fatal(err)
	}

	s := status.NewSlackStatus()

	tokenCh := make(chan string)
	go runConfigServer(ctx, port, tokenCh)

	go func() {
		t := <-tokenCh
		log.Println("Setting token", t)
		s.WithSlackToken(t)
	}()

	prefs := systray.AddMenuItem("Preferences", "")
	go func() {
		<-prefs.ClickedCh

		cmd := exec.Command("open", fmt.Sprintf("http://localhost:%d", port))
		_ = cmd.Run()
	}()
	//
	mQuitOrig := systray.AddMenuItem("Quit", "")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		cancel()
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	isOnCh := make(chan bool)
	go func() {
		for {
			isOn := <-isOnCh
			switch isOn {
			case true:
				systray.SetIcon(cameraOnIconData)
				statusCameraOnText.Show()
			case false:
				systray.SetIcon(cameraOffIconData)
				statusCameraOnText.Hide()
			}
		}
	}()

	// TODO do something better than this. This needs to be triggered when we update the token in preferences
	for {
		err = s.SetStatusWhenWebcamIsBusy(ctx, isOnCh)
		if err != nil {
			_ = beeep.Notify("go-away", err.Error(), "")
		}
	}
}

func onExit() {
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
</html>
`))
	}).Methods("GET")


	r.HandleFunc("/preferences", func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")
		log.Println("got slack token", token)
		tokenCh <- token

		//TODO stop race condition here since this is in a go routine
		slackTokenOption = token
	})

	http.Handle("/", r)

	// TODO make this closable from tokenCh closing
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil); err != nil {
		panic(err)
	}
}

func handleWindow() {
	w := webview.New(webview.Settings{
		Title: "Some title",
		Height: 1000,
		Width: 1000,
		URL: "https://google.com",
	})
	w.Run()
}
