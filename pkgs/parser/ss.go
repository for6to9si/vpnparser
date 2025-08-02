package parser

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/for6to9si/vpnparser/pkgs/utils"
)

var SSMethod map[string]struct{} = map[string]struct{}{
	"2022-blake3-aes-128-gcm":       {},
	"2022-blake3-aes-256-gcm":       {},
	"2022-blake3-chacha20-poly1305": {},
	"none":                          {},
	"aes-128-gcm":                   {},
	"aes-192-gcm":                   {},
	"aes-256-gcm":                   {},
	"chacha20-ietf-poly1305":        {},
	"xchacha20-ietf-poly1305":       {},
	"aes-128-ctr":                   {},
	"aes-192-ctr":                   {},
	"aes-256-ctr":                   {},
	"aes-128-cfb":                   {},
	"aes-192-cfb":                   {},
	"aes-256-cfb":                   {},
	"rc4-md5":                       {},
	"chacha20-ietf":                 {},
	"xchacha20":                     {},
}

/*
shadowsocks: ['plugin', 'obfs', 'obfs-host', 'mode', 'path', 'mux', 'host']
*/

type ParserSS struct {
	Address  string
	Port     int
	Method   string
	Password string

	Host     string
	Mode     string
	Mux      string
	Path     string
	Plugin   string
	OBFS     string
	OBFSHost string

	*StreamField
}

func (that *ParserSS) Parse(rawUri string) {
	rawUri = that.handleSS(rawUri)
	if u, err := url.Parse(rawUri); err == nil {
		that.StreamField = &StreamField{}
		that.Address = u.Hostname()
		that.Port, _ = strconv.Atoi(u.Port())

		that.Method = u.User.Username()
		if that.Method == "rc4" {
			that.Method = "rc4-md5"
		}
		if _, ok := SSMethod[that.Method]; !ok {
			that.Method = "none"
		}
		that.Password, _ = u.User.Password()

		query := u.Query()
		that.Host = query.Get("host")
		that.Mode = query.Get("mode")
		that.Mux = query.Get("mux")
		that.Path = query.Get("path")
		that.Plugin = query.Get("plugin")
		that.OBFS = query.Get("obfs")
		that.OBFSHost = query.Get("obfs-host")
	}
}

func (that *ParserSS) handleSS(rawUri string) string {
	rawUri = strings.ReplaceAll(rawUri, "#ss#\u00261@", "@")

	u, err := url.Parse(rawUri)
	if err != nil {
		fmt.Println(1)
		return rawUri
	}
	if !strings.Contains(rawUri, "@") {
		fmt.Println(2)
		u.Scheme = ""
		u.Fragment = ""

		body := strings.Trim(u.String(), "/")
		fmt.Println(body)

		if decoded, err := base64.StdEncoding.DecodeString(utils.ResoveBase64Padding(body)); err == nil {
			rawUri = strings.Replace(rawUri, body, string(decoded), 1)
			fmt.Println(rawUri)
		} else {
			fmt.Println(err)
		}
	} else {
		username := u.User.Username()

		if decoded, err := base64.StdEncoding.DecodeString(utils.ResoveBase64Padding(username)); err == nil {
			rawUri = strings.Replace(rawUri, username, string(decoded), 1)
		} else {
			fmt.Println(err)
		}
	}

	return rawUri
}

func (that *ParserSS) GetAddr() string {
	return that.Address
}

func (that *ParserSS) GetPort() int {
	return that.Port
}

func (that *ParserSS) Show() {
	fmt.Printf("addr: %s, port: %d, method: %s, password: %s\n",
		that.Address,
		that.Port,
		that.Method,
		that.Password)
}
