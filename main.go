package main

import (
	"net/http"
	"os"

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
		cli.StringFlag{"s3-bucket", "", "S3 bucket to save files to", "AWS_S3_BUCKET"},
		cli.StringFlag{"port", "3000", "port to run on", "PORT"},
	}
	app.Action = Run
	app.Run(os.Args)
}

func Run(c *cli.Context) {
	client := NewOAuthClient(c)

	routes := gin.New()
	routes.GET("/facebook/auth", client.Authorize)
	routes.GET("/facebook/login", func(c *gin.Context) {
		c.Redirect(302, client.AuthCodeURL("beatsforboobs"))
	})
	routes.GET("/", func(c *gin.Context) {
		c.Redirect(302, "http://beatsforboobs.org")
	})

	http.ListenAndServe(":"+c.String("port"), routes)
}
