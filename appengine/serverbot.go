package serverbot

import (
	"fmt"
	"html/template"
	"net/http"
)

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/run", run)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, botSettingsHTML)
}

func run(w http.ResponseWriter, r *http.Request) {
	err := botRunningTemplate.Execute(w, r.FormValue("user"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
      {{.}} is running successfully
    <form action="/" method="post">
      <div><input type="submit" value="Stop"></div>
    </form>
  </body>
</html>
`
