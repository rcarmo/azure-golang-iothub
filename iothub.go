package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type IoTHub struct {
	HostName             string
	SharedAccessKeyName  string
	SharedAccessKeyValue string
}

func NewIoTHub(conn string) (hub *IoTHub, err error) {
	fields, err := url.ParseQuery(conn)
	if err != nil {
		log.Fatal(err)
	}
	hub = new(IoTHub)
	t := reflect.ValueOf(hub).Elem()
	for k, v := range fields {
		val := t.FieldByName(k)
		val.Set(reflect.ValueOf(v[0]))
	}
	return hub, nil
}

func buildSasToken(hub *IoTHub) string {
	tokenValidSecs := 365 * 24 * 60 * 60
	key, _ := base64.StdEncoding.DecodeString(hub.SharedAccessKeyValue)
	mac := hmac.New(sha256.New, key)
	timestamp := strconv.FormatInt(time.Now().Unix()+int64(tokenValidSecs), 10)
	signature := url.QueryEscape(base64.URLEncoding.EncodeToString(mac.Sum([]byte(strings.ToLower(hub.HostName) + "\n" + timestamp))))
	return fmt.Sprintf("SharedAccessSignature sr=%s&sig=%s&se=%s&skn=%s",
		strings.ToLower(hub.HostName), signature, timestamp, hub.SharedAccessKeyName)
}

func main() {
	hub, _ := NewIoTHub("HostName=foobar.com;SharedAccessKeyName=blabla;SharedAccessKeyValue=puppies")
	fmt.Println(buildSasToken(hub))
}
