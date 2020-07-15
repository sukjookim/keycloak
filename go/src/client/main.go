package main

import (
	"fmt"
        "net/http"
	"log"
	"html/template"
	"net/url"
)

var oauth = struct{
  authURL string
  logout string
}{
	authURL: "http://133.186.159.116:8080/auth/realms/learningApp/protocol/openid-connect/auth",
	logout: "http://133.186.159.116:8080/auth/realms/learningApp/protocol/openid-connect/logout",
}

var t = template.Must(template.ParseFiles("template/index.html"))

type AppVar struct{
  AuthCode string
  SessionState string
}

var appVar = AppVar{}

func main() {
  fmt.Println("hello")
  http.HandleFunc("/", home)
  http.HandleFunc("/login", login)
  http.HandleFunc("/logout", logout)
  http.HandleFunc("/authCodeRedirect", authCodeRedirect)
  http.ListenAndServe(":8081", nil)
}

func home( w http.ResponseWriter, r *http.Request){

	t.Execute(w,nil)
}

func login( w http.ResponseWriter, r *http.Request){
	// create a redirect URL for authentication endpoint.
	req, err := http.NewRequest("GET", oauth.authURL, nil)
	if err != nil {
            log.Print(err)
	    return
	}
        // "state=123abc&client_id=billingApp&response_type=code"
        qs := url.Values{}
	qs.Add("state", "123")
	qs.Add("client_id", "billingApp")
	qs.Add("response_type", "code")
        qs.Add("redirect_uri", "http://133.186.159.116:8081/authCodeRedirect")
	req.URL.RawQuery = qs.Encode()
	http.Redirect(w, r, req.URL.String(), http.StatusFound)
}

func authCodeRedirect(w http.ResponseWriter, r *http.Request) {
   // fmt.Printf("Request queries : %v", r.URL.Query())
   // fmt.Fprintf(w, "handling auth code redirection")
   appVar.AuthCode = r.URL.Query().Get("code")
   appVar.SessionState = r.URL.Query().Get("session_state")
   r.URL.RawQuery = ""
   fmt.Printf("Request queries : %+v\n", appVar)
   // t.Execute(w,nil)
   http.Redirect(w, r, "http://133.186.159.116:8081", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
    q := url.Values{}
    q.Add("redirect_uri", "http://133.186.159.116:8081")

    logoutURL, err := url.Parse(oauth.logout)
    if err != nil {
       log.Println(err)
   }
   logoutURL.RawQuery = q.Encode()
    http.Redirect(w, r, logoutURL.String(), http.StatusFound)
}
