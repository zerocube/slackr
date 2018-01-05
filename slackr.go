/* ==========================================================================

User Parameters™:
  -t <target>   The Slack channel or username to send a message to.
  -n <name>     The Slack name to send messages as.
  -m <message>  The message to send

============================================================================= */

package main

import (
  "bytes"
  "encoding/json"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"
)

type App_Options struct {
  target,
  name,
  message,
  webhook_url string
  verbose bool
}
type Msg_Payload struct {
  Text,
  Channel,
  Username,
  Icon_Emoji string
  Link_Names int
}

func (opt *App_Options) Load() App_Options {
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
  flag.StringVar( &opt.webhook_url,
                  "webhook",
                  "",
                  "The Slack webhook.")
  flag.BoolVar( &opt.verbose,
                "verbose",
                false,
                "Enables verbose output.")
  flag.Parse()
  return *opt
}
func (opt *App_Options) OverrideWebhook(url string) App_Options {
  opt.webhook_url = url
  return *opt
}

func main() {
  options := App_Options{}
  options.Load()
  if options.verbose == true {
    fmt.Println("Payload:")
    fmt.Println("  - target:  ", options.target)
    fmt.Println("  - name:    ", options.name)
    fmt.Println("  - message: ", options.message) 
  }

  //  Check the environment for a WEBHOOK_URL if there wasn't one specified.
  //  The one provided the command line takes precedence, so we only load from
  //  the environment if we need it.
  if options.webhook_url == "" {
    env_webhook_url := os.Getenv("WEBHOOK_URL")
    if env_webhook_url != "" {
      options.OverrideWebhook(env_webhook_url)
    } else {
      fmt.Println(
        "Unable to determine webhook URL from command line parameter, or from ",
        "environment variable.",
      )
    }
  }

  //  Generate the message payload
  hook_payload := &Msg_Payload{
    Text: options.message,
    Channel: options.target,
    Username: options.name,
  }
  json_payload, err := json.Marshal(hook_payload)
  if err != nil {
    fmt.Println("An error occurred while trying to create the JSON payload.")
    log.Fatal(err)
    os.Exit(1)
  }

  //  Create the web request
  req, err := http.NewRequest(
    "POST",
    options.webhook_url,
    bytes.NewBuffer(json_payload) )
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    log.Fatal(err)
    os.Exit(2)
  }
  response_body, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    fmt.Println("An error occurred while trying to make a POST request to the",
                "webhook URL.")
    if options.verbose == true {
      log.Fatal(err)
    }
    os.Exit(10)
  }
  fmt.Printf("%s", response_body)
}
