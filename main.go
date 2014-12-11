package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
)

func main() {
	app := cli.NewApp()
	app.Name = "facebook-albums"
	app.Usage = "extracts facebook album metadata from facebook via the graph api"
	app.Version = "0.1"
	app.Author = "Matt Ho"
	app.Flags = []cli.Flag{
		cli.StringFlag{"app-id", "", "Facebook app client id", "FACEBOOK_APP_ID"},
		cli.StringFlag{"app-secret", "", "Facebook app client secret", "FACEBOOK_APP_SECRET"},
		cli.StringFlag{"redirect-url", "http://localhost:3000/facebook/auth", "OAuth callback endpoint", "FACEBOOK_OAUTH_URL"},
	}
	app.Action = Run
	app.Run(os.Args)
}

func Server(client *OAuthClient) {
	routes := gin.New()
	routes.GET("/facebook/auth", client.Authorize)

	http.ListenAndServe(":3000", routes)
}

func Run(c *cli.Context) {
	client := NewOAuthClient(c)

	// start the web server
	go Server(client)

	// wait for the server to start
	<-time.After(250 * time.Millisecond)

	cmd := exec.Command("open", client.AuthCodeURL("facebook-albums"))
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln(err)
	}

	// wait for completion
	<-client.Quit
}
