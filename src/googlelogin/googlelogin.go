package googlelogin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	historySession "github.com/tokopedia/historykv/src/session"
	historyUtil "github.com/tokopedia/historykv/src/util"
)

type GLogin interface {
	LoginURL() string
	GetUser(state string, code string) string
	GetDomain() string
	IsEnabled() bool
}

type GLoginData struct {
	OAuth       *oauth2.Config
	Session     historySession.Session
	CallbackURI string
	Domain      string
	Enabled     bool
}

func New(client string, secret string, domain string, callbackURI string, session historySession.Session) GLogin {
	g := GLoginData{
		Enabled: false,
	}

	if client != "" && secret != "" && domain != "" && callbackURI != "" && session != nil {
		g.OAuth = &oauth2.Config{
			RedirectURL:  fmt.Sprintf("%s/glogin/callback", callbackURI),
			ClientID:     client,
			ClientSecret: secret,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		}

		g.Session = session
		g.CallbackURI = callbackURI
		g.Domain = domain
		g.Enabled = true
	}

	var gl GLogin
	gl = &g

	return gl
}

func (g *GLoginData) IsEnabled() bool {
	return g.Enabled
}

func (g *GLoginData) GetDomain() string {
	return g.Domain
}

func (g *GLoginData) LoginURL() string {
	if !g.Enabled {
		return ""
	}

	token := fmt.Sprintf("glogin-%s", historyUtil.CreateRandomToken())
	g.Session.Set(token, "1")
	return g.OAuth.AuthCodeURL(token)
}

func (g *GLoginData) GetUser(state string, code string) string {
	if !g.Enabled {
		return ""
	}

	getState := g.Session.Get(state)

	if getState == "1" {
		g.Session.Delete(state)
		getToken, getTokenErr := g.OAuth.Exchange(oauth2.NoContext, code)
		if getTokenErr == nil {
			client := g.OAuth.Client(oauth2.NoContext, getToken)
			resp, respErr := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
			if respErr == nil {
				defer resp.Body.Close()
				bodyBytes, errReadBody := ioutil.ReadAll(resp.Body)
				if errReadBody == nil {
					var infoData map[string]interface{}
					errUnmarshal := json.Unmarshal(bodyBytes, &infoData)
					if errUnmarshal == nil {
						if infoData["email"] != nil {
							email := infoData["email"].(string)
							domainEmail := fmt.Sprintf("@%s", g.Domain)
							if strings.HasSuffix(email, domainEmail) {
								return strings.Replace(email, domainEmail, "", 1)
							}
						}
					}
				}
			}
		}
	}

	return ""
}
