package api_test

import (
	"crypto/tls"
	"logger/config"
	"net/http"
	"testing"
)

func TestMain(m *testing.M) {

}

func initialize() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	config.ListenPort = "127.0.0.1:7121"
	config.DBAddress = "127.0.0.1:4002"
	config.AuthAddress = ""
	config.Version = 0.01
	config.Environment = "test"
}
