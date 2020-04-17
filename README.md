# go-away

A tiny CLI tool to change your Slack status when you're on webcam. Useful for stopping your girlfriend from walking into the room when you're on a work call

### Setup

1. `make install`

2. [Set up a slack app](https://api.slack.com/authentication/basics) and install it in your slack workspace 

3. Find your "OAuth Access Token" from the Slack app console and set it to a `SLACK_API_TOKEEN` environment variable

### Usage

Your slack API token should be set in the `SLACK_API_TOKEN` env variable

```
$ go-away --help
NAME:
   goaway - ./goaway

USAGE:
   go-away [global options] command [command options] [arguments...]

DESCRIPTION:
   automatically set a custom slack status when you're on webcam

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug               enable debug logging (default: false)
   --status-text value   text to use for slack status (default: "On webcam")
   --status-emoji value  emoji to use for slack status (default: "🎥")
   --help, -h            show help (default: false)

```

### Issues
- Currently only supports OSX
