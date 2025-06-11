package parser

import (
	"net/url"
	"strconv"
)


type ParserHysteria2 struct {
	Address        string
	Port           int
	UpMbps        int
	DownMbps      int
	Password      string
	ObfType       string
	ObfPassword   string


	*StreamField
}

var ObfsTypes map[string]struct{} = map[string]struct{}{
	"salamander":                {},
} 

func (that *ParserHysteria2) Parse(rawURI string) {
	if r, err := url.Parse(rawURI); err == nil {
		that.Address = r.Hostname()
		that.Port, _ = strconv.Atoi(r.Port())

		that.Password = r.User.Username()

		query := r.Query()

		if _, ok := ObfsTypes[query.Get("obfs")]; ok {
			that.ObfPassword = query.Get("obfs-password")
			that.ObfType = query.Get("obfs")
		}


		that.StreamField = &StreamField{
			Network:          query.Get("type"),
			StreamSecurity:   query.Get("security"),
			Path:             query.Get("path"),
			Host:             query.Get("host"),
			GRPCServiceName:  query.Get("serviceName"),
			GRPCMultiMode:    query.Get("mode"),
			ServerName:       query.Get("sni"),
			TLSALPN:          query.Get("alpn"),
			Fingerprint:      query.Get("fp"),
			RealityShortId:   query.Get("sid"),
			RealitySpiderX:   query.Get("spx"),
			RealityPublicKey: query.Get("pbk"),
			PacketEncoding:   query.Get("packetEncoding"),
			TCPHeaderType:    query.Get("headerType"),
			TLSAllowInsecure: query.Get("insecure"),
		}
	}
}

func (that *ParserHysteria2) GetAddr() string {
	return that.Address
}

func (that *ParserHysteria2) GetPort() int {
	return that.Port
}

// func (that *ParserHysteria2) Show() {
// 	fmt.Printf("addr: %s, port: %v, password: %s\n",
// 		that.Server,
// 		that.ServerPort,
// 		that.Password)
// }
