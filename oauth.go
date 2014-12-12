package main

import (
	"log"

	"code.google.com/p/goauth2/oauth"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
)

type OAuthClient struct {
	clientId     string
	clientSecret string
	redirectURL  string
	config       *oauth.Config
}

func NewOAuthClient(c *cli.Context) *OAuthClient {
	return &OAuthClient{
		config: &oauth.Config{
			ClientId:     c.String("app-id"),
			ClientSecret: c.String("app-secret"),
			RedirectURL:  c.String("redirect-url"),
			AuthURL:      "https://graph.facebook.com/oauth/authorize",
			TokenURL:     "https://graph.facebook.com/oauth/access_token",
		},
	}
}

func (o *OAuthClient) Authorize(c *gin.Context) {
	form := struct {
		Code string `form:"code"`
	}{}

	if !c.Bind(&form) {
		c.JSON(500, map[string]string{"status": "unable to parse form data"})
		return
	}

	t := &oauth.Transport{Config: o.config}
	tok, err := t.Exchange(form.Code)
	if err != nil {
		c.Fail(500, err)
		return
	}

	go func() {
		err := Export(tok.AccessToken)
		if err != nil {
			log.Println(err)
		}
	}()

	c.JSON(200, tok)
}

func (o *OAuthClient) AuthCodeURL(state string) string {
	return o.config.AuthCodeURL(state)
}
