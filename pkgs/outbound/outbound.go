package outbound

import (
	"fmt"

	"github.com/for6to9si/vpnparser/pkgs/outbound/sing"
	"github.com/for6to9si/vpnparser/pkgs/outbound/xray"
	"github.com/for6to9si/vpnparser/pkgs/parser"
	"github.com/for6to9si/vpnparser/pkgs/utils"
)

type ClientType string

const (
	XrayCore ClientType = "xray"
	SingBox  ClientType = "sing"
)

func GetOutbound(clientType ClientType, rawUri string) (result IOutbound) {
	scheme := utils.ParseScheme(rawUri)
	switch clientType {
	case XrayCore:
		switch scheme {
		case parser.SchemeVmess:
			result = &xray.VmessOut{RawUri: rawUri}
		case parser.SchemeVless:
			result = &xray.VlessOut{RawUri: rawUri}
		case parser.SchemeTrojan:
			result = &xray.TrojanOut{RawUri: rawUri}
		case parser.SchemeSS:
			result = &xray.ShadowSocksOut{RawUri: rawUri}
		default:
			fmt.Println("unsupported protocol: ", scheme)
			return
		}
	case SingBox:
		switch scheme {
		case parser.SchemeVmess:
			result = &sing.SVmessOut{RawUri: rawUri}
		case parser.SchemeVless:
			result = &sing.SVlessOut{RawUri: rawUri}
		case parser.SchemeTrojan:
			result = &sing.STrojanOut{RawUri: rawUri}
		case parser.SchemeSS:
			result = &sing.SShadowSocksOut{RawUri: rawUri}
		case parser.SchemeSSR:
			result = &sing.SShadowSocksROut{RawUri: rawUri}
		case parser.SchemeWireguard:
			result = &sing.SWireguardOut{RawUri: rawUri}
		case parser.SchemeHy2, parser.SchemeHysteria2:
			result = &sing.SHysteria2Out{RawUri: rawUri}
		case parser.SchemeHysteria:
			result = &sing.SHysteriaOut{RawUri: rawUri}
		default:
			fmt.Println("unsupported protocol: ", scheme)
			return
		}
	default:
		fmt.Println("unsupported client type")
	}
	return
}
