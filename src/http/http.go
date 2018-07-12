package http

import (
	"log"
	"net/http"

	historyConsul "github.com/tokopedia/historykv/src/consul"
	historyDB "github.com/tokopedia/historykv/src/db"
	historyGoogleLogin "github.com/tokopedia/historykv/src/googlelogin"
	historySession "github.com/tokopedia/historykv/src/session"
	historyUtil "github.com/tokopedia/historykv/src/util"
)

type HTTP struct {
	Session        historySession.Session
	DB             historyDB.DB
	API            historyConsul.Consul
	GLogin         historyGoogleLogin.GLogin
	IsDisableLogin bool
}

func New(session historySession.Session, db historyDB.DB, api historyConsul.Consul, glogin historyGoogleLogin.GLogin, disableLogin bool) *HTTP {
	var h HTTP

	if session == nil {
		log.Fatalln("Session is Nil. Failed to generate?")
	}

	if db == nil {
		log.Fatalln("History DB is Nil. Are you sure you're connecting to right db?")
	}

	if api == nil {
		log.Fatalln("Consul API is Nil. Are you sure you're connecting to right consul http api?")
	}

	if glogin == nil {
		log.Fatalln("Google Login Data is Nil. Not sure why.")
	}

	h.Session = session
	h.DB = db
	h.API = api
	h.GLogin = glogin
	h.IsDisableLogin = disableLogin

	return &h
}

func (h *HTTP) GetUserFromRequest(r *http.Request) string {
	if h.IsDisableLogin {
		return "anonymous"
	}

	cookieSID, errCookieSID := r.Cookie("SID_HKV")
	if errCookieSID == nil {
		sid := cookieSID.Value
		if sid != "" {
			return h.Session.Get(sid)
		}
	}

	return ""
}

func (h *HTTP) GetTokenFromRequest(r *http.Request) string {
	aclToken, errACLToken := r.Cookie("ACL_TOKEN")
	if errACLToken == nil {
		return aclToken.Value
	}

	return ""
}

func (h *HTTP) CreateSession(user string) string {
	sessionID := historyUtil.CreateSessionToken()
	isSet := h.Session.Set(sessionID, user)

	if isSet {
		return sessionID
	}

	return ""
}
