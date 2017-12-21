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

type payload struct {
  target, name, message string
}

func (p *payload) Load() payload {
  flag.StringVar( &p.target,
                  "target",
                  "#general",
                  "The target to send messages to.")
  flag.StringVar( &p.name,
                  "name",
                  "Slackr",
                  "The username to send messages as.")
  flag.StringVar( &p.message,
                  "message",
                  "",
                  "The message to send.")
  return *p
}

func main() {
  msg_payload := payload{}
  msg_payload.Load()
  fmt.Println("Payload:")
  fmt.Println("  - target:  ", msg_payload.target)
  fmt.Println("  - name:    ", msg_payload.name)
  fmt.Println("  - message: ", msg_payload.message)
}
