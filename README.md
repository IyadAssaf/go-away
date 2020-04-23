<p align="center">-<img alt="go-away" src="assets/logo/1000-no_padding.png" width="750"/></p>


A tiny CLI tool and taskbar app to change your Slack status when you're on webcam. 
Useful for stopping peoples from walking into the room when you're on a work call

### Setup

1. [Download and install the most recent version of Go](https://golang.org/dl). 
 
2. [Set up a slack app](https://api.slack.com/authentication/basics) and install it in your slack workspace
 

### GoAway.app

- Run `make install-app` and copy your app to your Applications folder

- Start the app, you'll see the app in your task bar

![go-away](assets/screenshot.png)
  
- Select "Preferences" from the dropdown menu and enter your "OAuth Access Token" from the Slack app console. You can also edit your slack emoji and status message here.

![prefs](assets/prefs.png)

### CLI `go-away` 

3. Find your "OAuth Access Token" from the Slack app console and set it to an environment variable called `SLACK_API_TOKEN`

- Run `make install-cli` to install it to install the CLI tool

```
$ go-away --help
NAME:
   go-away - ./go-away

USAGE:
   go-away [global options] command [command options] [arguments...]

DESCRIPTION:
   update slack with a status when you're on webcam

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug               enable debug logging (default: false)
   --status-text value   text to use for slack status (default: "On webcam")
   --status-emoji value  emoji to use for slack status (default: "🎥")
   --refresh-rate value  number of seconds to refresh webcam status (default: 5)
   --help, -h            show help (default: false)
```
