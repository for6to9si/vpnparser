package parser

import (
	"fmt"
	"net/url"
	"strconv"
)

type ParserHysteria struct {
	Address        string
	Port           int
	Up 			  string
	UpMbps        int
	Down          string
	DownMbps      int
	Auth          string
	Obfs          string


	*StreamField
}

func (that *ParserHysteria) Parse(rawURI string) {
	if r, err := url.Parse(rawURI); err == nil {
		that.Address = r.Hostname()
		that.Port, _ = strconv.Atoi(r.Port())


		query := r.Query()
		that.Auth = query.Get("auth")
		that.Obfs = query.Get("obfs")
		
		if downMbps, err := strconv.Atoi(query.Get("downmbps")); err == nil {
			that.DownMbps = downMbps
			that.Down     = fmt.Sprintf("%d Mbps", downMbps)
		}

		if upMbps, err := strconv.Atoi(query.Get("upmbps")); err == nil {
			that.UpMbps = upMbps
			that.Up     = fmt.Sprintf("%d Mbps", upMbps)
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

func (that *ParserHysteria) GetAddr() string {
	return that.Address
}

func (that *ParserHysteria) GetPort() int {
	return that.Port
}
