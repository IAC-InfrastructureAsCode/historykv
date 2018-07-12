
# HistoryKV

HistoryKV is a simple app for storing history of Consul Key Value.

Using a simple user-password login, it save the person who edited it, and save the last value.


## Installation

Download on [Release Page](https://github.com/tokopedia/historykv/releases).

```
wget -O historykv https://github.com/tokopedia/historykv/releases/download/[version]/historykv-[os]
chmod +x historykv
./historykv
```
And then access http://127.0.0.1:9500/ via your favorite browser. First user is admin, with password admin.

You can add more user in the admin panel.



Note: It is recommended to use ACL System on the Consul. [Read Here](https://www.consul.io/docs/guides/acl.html) for more info.

## Parameter
```
  -admin-password [password-string]
        This will change the password of user "admin" to desired input.

  -config [path/to/file.ini]
        Configuration File that will be used. Or leave it empty if you want use parameter.

  -h
        List all parameter.
```

### Configuration

We have 2 version of configuration, Parameter and Config File.

#### Parameter Configuration
```
  -ip [ipv4-string / ipv6-string]
        IP for HistoryKV to listen. If you have public ip, it is not recommended to use 0.0.0.0.
        Default: 0.0.0.0

  -port [port-number]
        Port for HistoryKV to listen.
        Default: 9500

  -limit [number]
        Limit for History to save on database. If a key have more than 5, it will automatically deleted.
        Default: 5

  -sqlite-path [path/to/file.db]
        Location for SQLite db to write and read.
        Default: ./historykv.db

  -use-mysql [user:password@tcp(ip:port)/dbname]
        Use MySQL instead of SQLite on saving history. This allow multiple instance running at the same time. Input is MySQL DSN, Ex: "[user]:[password]@tcp(192.168.0.1:3306)/dbname"
        Default:

  -use-redis [ip:port]
        Use Redis instead of MemoryTTL on saving Session. This allow multiple instance running at the same time. Input is IP:Port, Ex: 192.168.0.1:6379
        Default:

  -consul-uri [url-string]
        Consul URI that contain API for KV, without trailing slash.
        Default: http://localhost:8500

  -consul-dc [datacenter-string]
        Consul Datacenter. You must define one if you have more than one cluster.
        Default:

  -consul-prefix [prefix/folder/]
        Key Prefix for Consul KV, with trailing slash. This is useful when you want a specific folder to use instead of root folder. Example: folder/folder-2/
        Default:

  -consul-token [token-string]
        Default ACL Token uses for Consul API to get, edit and delete key value.
        Default:

  -google-login-callback-uri [uri]
        This application uri to use Google Login, used for Callback, without trailing slash. Input this if you want to enable Google Login. Ex: http://consul.internal.com/historykv
        Default:

  -google-login-client-id [client-id-string]
        Google Login OAuth 2.0 Credentials Client ID. Input this if you want to enable Google Login.
        Default:

  -google-login-client-secret [client-secret-string]
        Google Login OAuth 2.0 Credentials Client Secret. Input this if you want to enable Google Login.
        Default:

  -google-login-domain string
        Your Google Login E-Mail Domain. Input this if you want to enable Google Login.
        Default: company.com
```

### File Configuration
Create a .ini file anywhere, contains:
```
[Listen]
; IP for HistoryKV to listen.
; If you have public ip, it is not recommended to use 0.0.0.0.
IP = "0.0.0.0"

; Port for HistoryKV to listen.
Port = 9500

[History]{
; Limit for History to save on database.
; If a key have more than 5, it will automatically deleted.
Limit = 5

[DB]
; Location for SQLite db to write and read.
Path = "./historykv.db"

; If not empty, use MySQL instead of SQLite on saving history.
; This allow multiple instance running at the same time.
; Input is MySQL DSN, Ex: "[user]:[password]@tcp(192.168.0.1:3306)/dbname"
MySQL = ""

[Session]
; If not empty, use Redis instead of MemoryTTL on saving Session.
; This allow multiple instance running at the same time.
; Input is IP:Port, Ex: 192.168.0.1:6379
Redis = ""

[Consul]
; Consul URI that contain API for KV, without trailing slash.
URI = "http://localhost:8500"

; Consul Datacenter. You must define one if you have more than one cluster.
Datacenter = ""

; Key Prefix for Consul KV, with trailing slash.
; This is useful when you want a specific folder to use instead of root folder.
; Example: folder/folder-2/
Prefix = ""

;Default ACL Token uses for Consul API to get, edit and delete key value.
Token = ""


[GoogleLogin]
; If CallbackURI, ClientID, ClientSecret and Domain is not empty
; GoogleLogin will be enabled.

; "This application" URI to use Google Login.
; Used for Callback, without trailing slash.
; Input this if you want to enable Google Login. Ex: http://consul.internal.com/historykv
CallbackURI = ""

; Google Login OAuth 2.0 Credentials Client ID.
; Input this if you want to enable Google Login.
ClientID = ""

; Google Login OAuth 2.0 Credentials Client Secret.
; Input this if you want to enable Google Login.
ClientSecret = ""

; Your Google Login E-Mail Domain.
; Input this if you want to enable Google Login.
Domain = "company.com"
```

## Using Google Login

To use Google Login OAuth2, you must specify the CallbackURI (Application URI), Client ID, Client Secret, Domain in the configuration.



To create Client ID and Secret, you must go to [Google Developer Console](https://console.developers.google.com/).

Then select New Project, fill the form as necessary.

Click on Credentials Tab, Create Credentials -> OAuth ID.

Select Web Application, fill name as you wish.



Javascript Origin is URI of your domain.

Callback URI: http://your.consul.domain/glogin/callback

If you use folder on the HistoryKV, you must change it to http://your.consul.domain/path/glogin/callback



Note that Google only allow a valid domain, but allowing private domain (not accessible by public).

For more info about Google OAuth2, [Read Here](https://developers.google.com/identity/protocols/OAuth2).



## Contribution

You can read CONTRIBUTION file for more info.

## Credit

### HTML

SVG Images by [FontAwesome](https://fontawesome.com/license)

JQuery: jquery.com

### GoLang

SQLite: github.com/mattn/go-sqlite3

MySQL: github.com/go-sql-driver/mysql

MemorTTL: github.com/koding/cache

Redis: github.com/xuyu/goredis

## License
[MIT](https://opensource.org/licenses/MIT)