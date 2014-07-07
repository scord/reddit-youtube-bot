package serverbot

import (
	"fmt"
	"github.com/scord/reddit-youtube-bot/bot"
	"html/template"
	"net/http"
)

const apiKey = ""

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/run", run)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, botSettingsHTML)
}

func run(w http.ResponseWriter, r *http.Request) {

	err := botRunningTemplate.Execute(w, "RUNNING")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = bot.Run(r.FormValue("chan"), r.FormValue("sub"), r.FormValue("user"), r.FormValue("pass"), apiKey)

	if err != nil {
		err = botRunningTemplate.Execute(w, "FAIL")
		return
	}

}

var botRunningTemplate = template.Must(template.New("run").Parse(botRunningHTML))

const botSettingsHTML = `
<html>
  <body>
    <form action="/run" method="post">
      Bot username:
      <div><textarea name="user" rows="3" cols="60"></textarea></div>
      Bot password:
      <div><textarea name="pass" rows="3" cols="60"></textarea></div>
      Bot subreddit:
      <div><textarea name="sub" rows="3" cols="60"></textarea></div>
      YouTube channel:
      <div><textarea name="chan" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Run"></div>
    </form>
  </body>
</html>
`

const botRunningHTML = `
<html>
  <body>
      {{.}}
    <form action="/" method="post">
      <div><input type="submit" value="Stop"></div>
    </form>
  </body>
</html>
`
