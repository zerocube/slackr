/* ==========================================================================

User Parameters™:
  -t <target>   The Slack channel or username to send a message to.
  -n <name>     The Slack name to send messages as.
  -m <message>  The message to send

============================================================================= */

package main

import (
  "flag"
  "fmt"
)

type app_options struct {
  target,
  name,
  message,
  api_token string
  verbose bool
}

func (opt *app_options) Load() app_options {
  flag.StringVar( &opt.target,
                  "target",
                  "#general",
                  "The target to send messages to.")
  flag.StringVar( &opt.name,
                  "name",
                  "Slackr",
                  "The username to send messages as.")
  flag.StringVar( &opt.message,
                  "message",
                  "",
                  "The message to send.")
  flag.StringVar( &opt.api_token,
                  "token",
                  "",
                  "The Slack API token.")
  flag.BoolVar( &opt.verbose,
                "verbose",
                false,
                "Enables verbose output.")
  return *opt
}

func main() {
  options := app_options{}
  options.Load()
  if options.verbose == true {
    fmt.Println("Payload:")
    fmt.Println("  - target:  ", options.target)
    fmt.Println("  - name:    ", options.name)
    fmt.Println("  - message: ", options.message) 
  }
}
