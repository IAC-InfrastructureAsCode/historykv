package http

import (
	"fmt"
	"net/http"
	"strings"
)

func (h *HTTP) PageGetAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if h.GetUserFromRequest(r) != "admin" {
		fmt.Fprintf(w, `<!DOCTYPE html><html><head><meta http-equiv="refresh" content="1; url=./"><script type="text/javascript">location.href="./";</script></head><body></body></html>`)
	} else {
		fmt.Fprint(w, h.GetAdminHTML())
	}
}

func (h *HTTP) PageGetIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if h.GetUserFromRequest(r) == "admin" {
		fmt.Fprintf(w, `<!DOCTYPE html><html><head><meta http-equiv="refresh" content="1; url=./admin"><script type="text/javascript">location.href="./admin";</script></head><body></body></html>`)
	} else {
		getHTML := h.GetIndexHTML()
		additionalScript := ""
		if h.GLogin.IsEnabled() {
			additionalScript += `<script type="text/javascript">enableGLogin = true;</script>`
		}
		getHTML = strings.Replace(getHTML, "<!-- ADDITIONAL_SCRIPT -->", additionalScript, 1)
		fmt.Fprint(w, getHTML)
	}
}
