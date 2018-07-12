package main

import (
	"flag"
	"log"

	"gopkg.in/gcfg.v1"
)

type Config struct {
	Listen struct {
		IP   string
		Port int
	}
	History struct {
		Limit int
	}
	DB struct {
		Path  string
		MySQL string
	}
	Session struct {
		Redis string
	}
	Consul struct {
		Prefix     string
		URI        string
		Token      string
		Datacenter string
	}
	GoogleLogin struct {
		CallbackURI  string
		ClientID     string
		ClientSecret string
		Domain       string
	}
}

func getConfig() *Config {
	// List of Flag
	flagConfigFile := flag.String("config", "", "Path to use config file instead of parameter.")

	// If not using config file
	flagIP := flag.String("ip", "0.0.0.0", "IP for HistoryKV to listen.")
	flagPort := flag.Int("port", 9500, "Port for HistoryKV to listen.")

	flagHistoryLimit := flag.Int("limit", 5, "Limit for History to save.")

	// DB Flag
	flagSQLitePath := flag.String("sqlite-path", "./historykv.db", "Location for SQLite db to write and read.")
	flagUseMySQL := flag.String("use-mysql", "", "Use MySQL instead of SQLite on saving history. This allow multiple instance running at the same time. Input is MySQL DSN, Ex: \"[user]:[password]@tcp(192.168.0.1:3306)/dbname\"")

	// Session Flag
	flagUseRedis := flag.String("use-redis", "", "Use Redis instead of MemoryTTL on saving Session. This allow multiple instance running at the same time. Input is IP:Port, Ex: 192.168.0.1:6379")

	// Consul API Flag
	flagConsulPrefix := flag.String("consul-prefix", "", "Key Prefix for Consul KV, with trailing slash. This is useful when you want a specific folder to use instead of root folder. Example: folder/folder-2/")
	flagConsulURI := flag.String("consul-uri", "http://localhost:8500", "Consul URI that contain API for KV, without trailing slash.")
	flagConsulToken := flag.String("consul-token", "", "ACL Token uses for Consul API to get, edit and delete key value.")
	flagConsulDatacenter := flag.String("consul-dc", "", "Consul Datacenter. You must define one if you have more than one cluster.")

	// Google Login Flag
	flagGoogleCallbackUri := flag.String("google-login-callback-uri", "", "This application uri to use Google Login, used for Callback, without trailing slash. Input this if you want to enable Google Login. Ex: http://consul.internal.com/historykv")
	flagGoogleClient := flag.String("google-login-client-id", "", "Google Login OAuth 2.0 Credentials Client ID. Input this if you want to enable Google Login.")
	flagGoogleSecret := flag.String("google-login-client-secret", "", "Google Login OAuth 2.0 Credentials Client Secret. Input this if you want to enable Google Login.")
	flagGoogleDomain := flag.String("google-login-domain", "company.com", "Your Google Login E-Mail Domain. Input this if you want to enable Google Login.")

	flag.Parse()

	var cfg Config

	if *flagConfigFile != "" {
		readConfigErr := gcfg.ReadFileInto(&cfg, *flagConfigFile)

		if readConfigErr != nil {
			log.Println("Config Error!")
			log.Println(readConfigErr)
			log.Fatalln("Exiting...")
		}

		if cfg.Listen.IP == "" {
			cfg.Listen.IP = "0.0.0.0"
		}

		if cfg.Listen.Port <= 0 {
			cfg.Listen.Port = 9500
		}

		if cfg.History.Limit <= 0 {
			cfg.History.Limit = 5
		}

		if cfg.DB.Path == "" {
			cfg.DB.Path = "./historykv.db"
		}

		if cfg.Consul.URI == "" {
			cfg.Consul.URI = "http://localhost:8500"
		}
	} else {
		cfg.Listen.IP = *flagIP
		cfg.Listen.Port = *flagPort
		cfg.History.Limit = *flagHistoryLimit
		cfg.DB.Path = *flagSQLitePath
		cfg.DB.MySQL = *flagUseMySQL
		cfg.Session.Redis = *flagUseRedis
		cfg.Consul.Prefix = *flagConsulPrefix
		cfg.Consul.URI = *flagConsulURI
		cfg.Consul.Token = *flagConsulToken
		cfg.Consul.Datacenter = *flagConsulDatacenter
		cfg.GoogleLogin.CallbackURI = *flagGoogleCallbackUri
		cfg.GoogleLogin.ClientID = *flagGoogleClient
		cfg.GoogleLogin.ClientSecret = *flagGoogleSecret
		cfg.GoogleLogin.Domain = *flagGoogleDomain
	}

	return &cfg
}
