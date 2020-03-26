package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Params struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}

// https://stackoverflow.com/questions/46948050/how-to-read-request-body-twice-in-golang-middleware

func AuthorizationMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		allowedMutations := []string{"CreateOrder", "Login"}
		var params Params

		// if strings.Contains(r.Header.Get("X-Forwarded-Host"), "graphql") {
		// 	fmt.Println("request coming from graphql playground")
		// 	h.ServeHTTP(w, r)
		// 	return
		// }

		// convert request body to bytes then unmarshal the json
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(b, &params); err != nil {
			h.ServeHTTP(w, r)
			return
		}
		// construct a new ReadCloser
		r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		// let slice of allowed mutations pass through
		for _, m := range allowedMutations {
			if strings.Contains(params.OperationName, m) {
				fmt.Printf("got %s mutation, no token necessary \n", params.OperationName)
				h.ServeHTTP(w, r)
				return
			}
		}
		// verify request is a mutation
		if strings.Contains(params.Query, "mutation") {
			fmt.Println("got mutation request")
			// check for valid JWT in authorization header
			// reqToken := r.Header.Get("Authorization")
			// splitToken := strings.Split(reqToken, "Bearer")
			// if len(splitToken) != 2 {
			// 	log.Fatal("Bearer token not in proper format")
			// }
			// reqToken = strings.TrimSpace(splitToken[1])
			// fmt.Printf("token is %s", reqToken)
			// if reqToken != "" {
			// 	h.ServeHTTP(w, r)
			// 	return
			// }
			h.ServeHTTP(w, r)
			return
		}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
