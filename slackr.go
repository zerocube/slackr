/* ==========================================================================

User Parameters™:
  -target <target>    The Slack channel or username to send a message to.
  -name <name>        The Slack name to send messages as.
  -message <message>  The message to send.
  -webhook_url <url>  The Slack webhook URL.
  -verbose            Enables verbose output

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

var git_version string

type App_Options struct {
	Target  string `json:"channel"`
	Name    string `json:"username"`
	Message string `json:"text"`
	Webhook string
	Emoji   string `json:"icon_emoji"`
	Verbose bool
	Version bool
}

func (opt *App_Options) Load() App_Options {
	flag.StringVar(&opt.Target,
		"target",
		"#general",
		"The target to send messages to.")
	flag.StringVar(&opt.Name,
		"name",
		"Slackr",
		"The username to send messages as.")
	flag.StringVar(&opt.Message,
		"message",
		"",
		"The message to send.")
	flag.StringVar(&opt.Webhook,
		"webhook",
		"",
		"The Slack webhook.")
	flag.StringVar(&opt.Emoji,
		"emoji",
		":thinking_face:",
		"The emoji to use.")
	flag.BoolVar(&opt.Verbose,
		"verbose",
		false,
		"Enables verbose output.")
	flag.BoolVar(&opt.Version,
		"version",
		false,
		"Outputs the version, if known.")
	flag.Parse()
	return *opt
}
func (opt *App_Options) OverrideWebhook(url string) App_Options {
	opt.Webhook = url
	return *opt
}

func main() {
	options := App_Options{}
	options.Load()

	//  If all we need is to output the version, then do that.
	if options.Version == true {
		//  Guilty until proven innocent
		slackr_version := ""

		if git_version != "" {
			slackr_version = git_version
		}

		fmt.Println(
			fmt.Sprintf("Slackr Version: %s", slackr_version),
		)
		os.Exit(0)
	}

	//  Check the environment for a WEBHOOK_URL if there wasn't one specified.
	//  The one provided the command line takes precedence, so we only load from
	//  the environment if we need it.
	if options.Webhook == "" {
		env_webhook := os.Getenv("SLACKR_WEBHOOK_URL")
		if env_webhook != "" {
			options.OverrideWebhook(env_webhook)
		} else {
			fmt.Println(
				"Unable to determine webhook URL from command line parameter, or from",
				"environment variable.",
			)
			os.Exit(1)
		}
	}

	if options.Verbose == true {
		fmt.Println("Payload:")
		fmt.Println("  - target:      ", options.Target)
		fmt.Println("  - name:        ", options.Name)
		fmt.Println("  - message:     ", options.Message)
		fmt.Println("  - emoji:       ", options.Emoji)
		fmt.Printf("  - webhook_url:  **********%s\n", string(
			options.Webhook[len(options.Webhook)-4:],
		))
	}

	json_payload, err := json.Marshal(options)
	if options.Verbose == true {
		fmt.Println("JSON payload map:", json_payload)
	}
	if err != nil {
		fmt.Println("An error occurred while trying to create the JSON payload.")
		log.Fatal(err)
		os.Exit(2)
	}
	post_payload := bytes.NewBuffer(json_payload)
	if options.Verbose == true {
		fmt.Println("POST payload map:", post_payload)
	}

	//  Create the web request
	req, err := http.NewRequest("POST", options.Webhook, post_payload)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}
	response_body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("An error occurred while trying to make a POST request to the",
			"webhook URL.")
		if options.Verbose == true {
			log.Fatal(err)
		}
		os.Exit(4)
	}
	fmt.Printf("%s\n", response_body)
}
