package sing

import (
	"github.com/M-logique/vpnparser/pkgs/parser"
	"github.com/M-logique/vpnparser/pkgs/utils"
	"github.com/gogf/gf/v2/encoding/gjson"
)

/*
https://sing-box.sagernet.org/zh/configuration/outbound/hysteria2/
{
  "type": "hysteria2",
  "tag": "hy2-out",

  "server": "127.0.0.1",
  "server_port": 1080,
  "server_ports": [
    "2080:3000"
  ],
  "hop_interval": "",
  "up_mbps": 100,
  "down_mbps": 100,
  "obfs": {
    "type": "salamander",
    "password": "cry_me_a_r1ver"
  },
  "password": "goofy_ahh_password",
  "network": "tcp",
  "tls": {},
  "brutal_debug": false,

  ... // Dial Fields
}
*/

var SingHysteria2 string = `{
    "type": "hysteria2",
    "tag": "hy2-out",
    "server": "127.0.0.1",
    "server_port": 1080,
    "obfs": {
        "type": "salamander",
        "password": "cry_me_a_r1ver"
    },
    "password": "goofy_ahh_password",
    "tls": {},
}`

type SHysteria2Out struct {
  RawUri  string
  Parser  *parser.ParserHysteria2
  outbound string
}

func (that *SHysteria2Out) Addr() string {
	if that.Parser == nil {
		return ""
	}
	return that.Parser.GetAddr()
}

func (that *SHysteria2Out) Port() int {
	if that.Parser == nil {
		return 0
	}
	return that.Parser.GetPort()
}

func (that *SHysteria2Out) Scheme() string {
	return parser.SchemeHysteria2
}

func (that *SHysteria2Out) GetRawUri() string {
	return that.RawUri
}

func (that *SHysteria2Out) Parse(rawUri string) {
	that.Parser = &parser.ParserHysteria2{}
	that.Parser.Parse(rawUri)
}


func (that *SHysteria2Out) getSettings() string {
    if that.Parser.Address == "" || that.Parser.Port == 0 {
        return "{}"
    }

    j := gjson.New(SingHysteria2) 
    j.Set("type", "hysteria2")
    j.Set("server", that.Parser.Address)
    j.Set("server_port", that.Parser.Port)
    if that.Parser.UpMbps != 0 && that.Parser.DownMbps != 0 {
        j.Set("up_mbps", that.Parser.UpMbps)
        j.Set("down_mbps", that.Parser.DownMbps)
    }
    j.Set("password", that.Parser.Password)
    if that.Parser.Network != "" {
        j.Set("network", that.Parser.Network)
    }
    
    j.Set("obfs.type", "")
    if that.Parser.ObfPassword != "" && that.Parser.ObfType != "" {
        j.Set("obfs.type", that.Parser.ObfType)
        j.Set("obfs.password", that.Parser.ObfPassword)
    }

    j.Set("tag", utils.OutboundTag)
    return j.MustToJsonString()
}

func (that *SHysteria2Out) GetOutboundStr() string {
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