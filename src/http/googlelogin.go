package http

import (
	"fmt"
	"net/http"
	"time"
)

func (h *HTTP) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	if !h.IsDisableLogin && h.GLogin.IsEnabled() {
		http.Redirect(w, r, h.GLogin.LoginURL(), 302)
	} else {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<!DOCTYPE html><html><head><script type=\"text/javascript\">opener.location.reload();window.close();</script></head><body></body></html>")
	}
}

func (h *HTTP) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if !h.IsDisableLogin && h.GLogin.IsEnabled() {
		getState := r.FormValue("state")
		getCode := r.FormValue("code")
		if getState != "" && getCode != "" {
			user := h.GLogin.GetUser(getState, getCode)
			if user != "" {
				if !h.DB.IsUserExists(user) {
					h.DB.AddUser(user, "-", "")
				}

				newCookie := http.Cookie{
					Name:    "SID_HKV",
					Value:   h.CreateSession(user),
					Expires: time.Now().Add(24 * time.Hour),
					Path:    "/",
				}

				http.SetCookie(w, &newCookie)

				tokenCookie := http.Cookie{
					Name:    "ACL_TOKEN",
					Value:   h.DB.GetToken(user),
					Expires: time.Now().Add(10 * 365 * 24 * time.Hour),
					Path:    "/",
				}

				http.SetCookie(w, &tokenCookie)
			}
		}
	}

	fmt.Fprintf(w, "<!DOCTYPE html><html><head><script type=\"text/javascript\">opener.location.reload();window.close();</script></head><body></body></html>")
}
