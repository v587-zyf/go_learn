package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"go.uber.org/zap"
	"kernel/log"
	"net/url"
	"sort"
	"strings"
)

type param struct {
	Key   string
	Value string
}

// 按键的ASCII值排序
type byKey []param

func (a byKey) Len() int           { return len(a) }
func (a byKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

func TgParseData(data string) (url.Values, error) {
	return url.ParseQuery(data)
}

func TgGetCheckData(urlQuery url.Values) string {
	var params []param
	for k, vs := range urlQuery {
		if k == "hash" {
			continue
		}
		params = append(params, param{Key: k, Value: vs[0]})
	}
	sort.Sort(byKey(params))

	var sortedQuery strings.Builder
	for _, p := range params {
		if sortedQuery.Len() > 0 {
			sortedQuery.WriteRune('\n')
		}
		sortedQuery.WriteString(p.Key)
		sortedQuery.WriteRune('=')
		sortedQuery.WriteString(p.Value)
	}

	//fmt.Println("Sorted Query String:", sortedQuery.String())
	return sortedQuery.String()
}

func TgGetHmacSha256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	hash := h.Sum(nil)

	//HmacUtils (HmacAlgorithms.HMAC SHA 256, botTokenData).hmac (botToken)

	//fmt.Println("HMAC-SHA256 Hash:", hex.EncodeToString(hash))
	//fmt.Println(base64.StdEncoding.EncodeToString(hash))
	//return hex.EncodeToString(hash)
	return hash
}

func TgCheck(initData, loginToken string) (tgDate url.Values, res bool) {
	tgDate, err := TgParseData(initData)
	if err != nil {
		log.Error("utils.TgParseData", zap.Error(err), zap.String("initData", initData))
		return
	}
	dataCheckString := TgGetCheckData(tgDate)
	botTokenData := "WebAppData"
	secret := TgGetHmacSha256([]byte(botTokenData), []byte(loginToken))
	hash := hex.EncodeToString(TgGetHmacSha256([]byte(secret), []byte(dataCheckString)))
	if hash != tgDate.Get("hash") {
		log.Error("hash not true", zap.String("makeHash", hash), zap.String("hash", tgDate.Get("hash")))
		return
	}

	res = true

	return
}
