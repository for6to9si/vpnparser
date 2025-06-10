package utils

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
)

const (
	OutboundTag string = "PROXY_OUT"
)

func SetJsonObjectByString(key, value string, gJSON *gjson.Json) (newGJSON *gjson.Json) {
	if gJSON == nil {
		return
	}
	tempValue := "OOXXOOXX"
	gJSON.Set(key, tempValue)
	result := strings.ReplaceAll(gJSON.MustToJsonString(), fmt.Sprintf(`"%s"`, tempValue), value)
	return gjson.New(result)
}

func ParseScheme(rawUri string) (scheme string) {
	sp := "://"
	sList := strings.Split(rawUri, sp)
	if len(sList) == 2 {
		scheme = sList[0] + sp
	}
	return
}

func ResoveBase64Padding(s string) string {
	if len(s) % 4 == 0 {
		return s
	}

	paddingNeeded := (4 - (len(s) % 4)) % 4

	for i := 0; i < paddingNeeded; i++ {
		s += "="
	}

	return s
}