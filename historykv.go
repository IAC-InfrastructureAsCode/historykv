package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	historyConsul "github.com/tokopedia/historykv/src/consul"
	historyDB "github.com/tokopedia/historykv/src/db"
	historyGoogleLogin "github.com/tokopedia/historykv/src/googlelogin"
	historyHTTP "github.com/tokopedia/historykv/src/http"
	historySession "github.com/tokopedia/historykv/src/session"
	historyUtil "github.com/tokopedia/historykv/src/util"
)

func main() {
	log.Println("--------------------------------------------------")
	log.Println("Starting HistoryKV")
	log.Println("--------------------------------------------------")
	// One Time Password
	flagAdminPassword := flag.String("admin-password", "", "New password for user admin. Leave it empty if you do not want to change password.")

	cfg := getConfig()

	// Get Session
	sessionType := "memory"
	sessionAddress := ""
	if cfg.Session.Redis != "" {
		sessionType = "redis"
		sessionAddress = cfg.Session.Redis
	}

	session := historySession.New(sessionType, sessionAddress)

	// Get API Consul
	if cfg.Consul.Prefix != "" && cfg.Consul.Prefix[len(cfg.Consul.Prefix)-1] != '/' {
		cfg.Consul.Prefix = fmt.Sprintf("%s/", cfg.Consul.Prefix)
	}

	consulAPI := historyConsul.New(cfg.Consul.URI, cfg.Consul.Token, cfg.Consul.Datacenter, cfg.Consul.Prefix)

	// Get DB
	dbType := "sqlite"
	dbPath := cfg.DB.Path
	if cfg.DB.MySQL != "" {
		dbType = "mysql"
		dbPath = cfg.DB.MySQL
	}

	db := historyDB.New(dbType, dbPath, cfg.History.Limit, cfg.Consul.Prefix)

	// If reset Password
	if *flagAdminPassword != "" {
		log.Println("--------------------------------------------------")

		passwordMD5 := historyUtil.CreateHashPassword("admin", *flagAdminPassword)
		if db.UpdateUser("admin", passwordMD5) {
			log.Printf("Admin password changed to: %s\n", *flagAdminPassword)
		} else {
			log.Println("Failed to change admin password.")
		}

		log.Println("--------------------------------------------------")
	}

	// Google Login
	glogin := historyGoogleLogin.New(cfg.GoogleLogin.ClientID, cfg.GoogleLogin.ClientSecret, cfg.GoogleLogin.Domain, cfg.GoogleLogin.CallbackURI, session)
	if glogin.IsEnabled() {
		log.Println("--------------------------------------------------")
		log.Println("Google Login is Enabled...")
		log.Printf("Using E-Mail Domain: [user]@%s\n", glogin.GetDomain())
		log.Println("--------------------------------------------------")
	} else {
		log.Println("--------------------------------------------------")
		log.Println("Google Login is Disabled...")
		log.Println("--------------------------------------------------")
	}

	// Init HTTP
	httpServe := historyHTTP.New(session, db, consulAPI, glogin)

	// Serve HTTP
	http.HandleFunc("/glogin/login", httpServe.GoogleLogin)
	http.HandleFunc("/glogin/callback", httpServe.GoogleCallback)
	http.HandleFunc("/ajax/user/login.json", httpServe.Login)
	http.HandleFunc("/ajax/user/logout.json", httpServe.Logout)
	http.HandleFunc("/ajax/user/change-password.json", httpServe.ChangePassword)
	http.HandleFunc("/ajax/user/change-token.json", httpServe.ChangeToken)
	http.HandleFunc("/ajax/kv/history.json", httpServe.GetHistoryID)
	http.HandleFunc("/ajax/kv/update.json", httpServe.CreateUpdateKey)
	http.HandleFunc("/ajax/kv/delete.json", httpServe.DeleteKey)
	http.HandleFunc("/ajax/kv/value.json", httpServe.GetValueKey)
	http.HandleFunc("/ajax/kv/list.json", httpServe.GetList)
	http.HandleFunc("/ajax/admin/list.json", httpServe.AdminGetUserList)
	http.HandleFunc("/ajax/admin/update.json", httpServe.AdminCreateChangeUser)
	http.HandleFunc("/ajax/admin/delete.json", httpServe.AdminDeleteUser)

	// Page
	http.HandleFunc("/", httpServe.PageGetIndex)
	http.HandleFunc("/admin", httpServe.PageGetAdmin)

	serve := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Listen.IP, cfg.Listen.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Listening to %s:%d\n", cfg.Listen.IP, cfg.Listen.Port)
	log.Println("--------------------------------------------------")
	log.Fatal(serve.ListenAndServe())
}
