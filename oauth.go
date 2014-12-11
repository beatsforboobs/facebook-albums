package main

import (
	"code.google.com/p/goauth2/oauth"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
)

type OAuthClient struct {
	clientId     string
	clientSecret string
	redirectURL  string
	Quit         chan struct{}
}

func NewOAuthClient(c *cli.Context) *OAuthClient {
	return &OAuthClient{
		clientId:     c.String("app-id"),
		clientSecret: c.String("app-secret"),
		redirectURL:  c.String("redirect-url"),
		Quit:         make(chan struct{}),
	}
}

func (o *OAuthClient) oauthConfig() *oauth.Config {
	return &oauth.Config{
		ClientId:     o.clientId,
		ClientSecret: o.clientSecret,
		AuthURL:      "https://graph.facebook.com/oauth/authorize",
		TokenURL:     "https://graph.facebook.com/oauth/access_token",
		RedirectURL:  o.redirectURL,
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

	t := &oauth.Transport{Config: o.oauthConfig()}
	tok, err := t.Exchange(form.Code)
	if err != nil {
		c.Fail(500, err)
		return
	}

	c.JSON(200, tok)
}

func (o *OAuthClient) AuthCodeURL(state string) string {
	return o.oauthConfig().AuthCodeURL(state)
}
