package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"logger/config"
	"net/http"
	"regexp"
)

func Start(collector http.HandlerFunc) error {
	mux := http.NewServeMux()

	mux.Handle("/logs", verify(handleRequest(testHandle)))

	// Post logs
	mux.Handle("/logs", verify(handleRequest(collector)))

	return http.ListenAndServe(config.ListenPort, mux)
}

// Need handlers to post, delete, and get logs
// Need handlers to handle authentication

func cors(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", config.CorsAllow)
	rw.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	rw.Header().Set("Access-Control-Allow-Headers", "X-AUTH-TOKEN, X-USER-TOKEN")
	rw.WriteHeader(200)
	rw.Write([]byte("success!\n"))
}

func handleRequest(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", config.CorsAllow)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "X-AUTH-TOKEN, X-USER-TOKEN")

		fn(w, r)

		getStatus := func(w http.ResponseWriter) string {
			r, _ := regexp.Compile("status:([0-9]*)")
			return r.FindStringSubmatch(fmt.Sprintf("%+v", w))[1]
		}

		getWrote := func(w http.ResponseWriter) string {
			r, _ := regexp.Compile("written:([0-9]*)")
			return r.FindStringSubmatch(fmt.Sprintf("%+v", w))[1]
		}

		log.Print(`%s - [%s] %s %s %s(%s) - "User-Agent: %s"`,
			r.RemoteAddr, r.Proto, r.Method, r.RequestURI,
			getStatus(w), getWrote(w), // %s(%s)
			r.Header.Get("User-Agent"))
	}
}

// verify that the token is allowed throught the authenticator
func verify(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		key := req.Header.Get("X-USER-TOKEN")
		// allow browsers to authenticate/fetch logs
		if key == "" {
			query := req.URL.Query()
			key = query.Get("X-USER-TOKEN")
			if key == "" {
				key = query.Get("x-user-token")
			}
		}
		// if !authenticator.Valid(key) {
		// 	rw.WriteHeader(401)
		// 	return
		// }
		fn(rw, req)
	}
}

func parse(r *http.Request, v interface{}) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}

func testHandle(w http.ResponseWriter, r *http.Request) {

}
