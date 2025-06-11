package sing

import (
	"encoding/base64"

	"github.com/M-logique/vpnparser/pkgs/parser"
	"github.com/M-logique/vpnparser/pkgs/utils"
	"github.com/gogf/gf/v2/encoding/gjson"
)

/*
https://sing-box.sagernet.org/configuration/inbound/hysteria/
{
  "type": "hysteria",
  "tag": "hysteria-in",

  ... // Listen Fields

  "up": "100 Mbps",
  "up_mbps": 100,
  "down": "100 Mbps",
  "down_mbps": 100,
  "obfs": "fuck me till the daylight",

  "users": [
    {
      "name": "sekai",
      "auth": "",
      "auth_str": "password"
    }
  ],

  "recv_window_conn": 0,
  "recv_window_client": 0,
  "max_conn_client": 0,
  "disable_mtu_discovery": false,
  "tls": {}
}
*/

var SingHysteria string = `{
  "type": "hysteria",
  "tag": "hysteria-in",

  "up": "100 Mbps",
  "up_mbps": 100,
  "down": "100 Mbps",
  "down_mbps": 100,
  "obfs": "fuck me till the daylight",

  "users": [],
  "tls": {}
}`

type SHysteriaOut struct {
    RawUri  string
    Parser  *parser.ParserHysteria
    outbound string
}

func (that *SHysteriaOut) Addr() string {
	if that.Parser == nil {
		return ""
	}
	return that.Parser.GetAddr()
}

func (that *SHysteriaOut) Port() int {
	if that.Parser == nil {
		return 0
	}
	return that.Parser.GetPort()
}

func (that *SHysteriaOut) Scheme() string {
	return parser.SchemeHysteria2
}

func (that *SHysteriaOut) GetRawUri() string {
	return that.RawUri
}

func (that *SHysteriaOut) Parse(rawUri string) {
	that.Parser = &parser.ParserHysteria{}
	that.Parser.Parse(rawUri)
}


func (that *SHysteriaOut) getSettings() string {
    if that.Parser.Address == "" || that.Parser.Port == 0 {
        return "{}"
    }

    j := gjson.New(SingHysteria2) 
	j.Set("type", "hysteria")
	j.Set("tag", utils.OutboundTag)
	j.Set("up", that.Parser.Up)
	j.Set("up_mbps", that.Parser.UpMbps)
	j.Set("down", that.Parser.Down)
	j.Set("down_mbps", that.Parser.DownMbps)
	j.Set("obfs", that.Parser.Obfs)
	
	j.Set("users.0.name", "user1")
	authBase64 := base64.StdEncoding.EncodeToString([]byte(that.Parser.Auth))
	j.Set("users.0.auth", authBase64)
	j.Set("users.0.auth_str", that.Parser.Auth)

	return j.MustToJsonString()
}

func (that *SHysteriaOut) GetOutboundStr() string {
    if that.outbound == "" {
		settings := that.getSettings()
		if settings == "{}" {
			return ""
		}
		cnf := gjson.New(settings)
		cnf = PrepareStreamStr(cnf, that.Parser.StreamField)
		that.outbound = cnf.MustToJsonString()
	}
	return that.outbound
}