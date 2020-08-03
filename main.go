package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var secretKey = []byte("abcdefgh")

type tokenResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type errorResponse struct {
	Token   string `json:"token"`
	Message error  `json:"message"`
}

func login(w http.ResponseWriter, r *http.Request) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  "123asd",
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(tokenResponse{Token: "", Message: "something went wrong!"})
		return
	}
	// w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse{Token: tokenString, Message: "login successful"})
	return
}

func authRoute(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Auth route")
}

func logout(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("logout route")
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		w.Header().Add("content-type","application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func authMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if(tokenString == "") {
			http.Error(w, "Invalid Authorization", http.StatusForbidden)
			return
		}

		authToken := strings.Split(tokenString," ")[1]

		fmt.Println(authToken);
		if authToken == "" {
			http.Error(w, "Invalid Token", http.StatusForbidden)
			return
		}

		// Parse takes the token string and a function for looking up the key. The latter is especially
		// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
		// head of the token to identify which key to use, but the parsed token (head and claims) is provided
		// to the callback, providing flexibility.
		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return secretKey, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["id"], claims["exp"])
		} else {
			fmt.Println(err)
			http.Error(w, "token expired", http.StatusForbidden)
			return
		}

		
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()

	 router.Use(middleware)
	 router.Use(mux.CORSMethodMiddleware(router))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("Hello from API ")
	})

	router.HandleFunc("/api/login", login).Methods(http.MethodPost,http.MethodOptions)

	authRouter := router.PathPrefix("/auth").Subrouter()
	 authRouter.Use(authMiddlware)
//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTYzODQyNTYsImlkIjoiMTIzYXNkIn0.ESmhVuBRDR6Fz1E3PjXt5SriuG2LKvJUqx8zh1en8gc
	authRouter.HandleFunc("/api/auth", authRoute).Methods(http.MethodPost,http.MethodOptions)
	authRouter.HandleFunc("/api/logout", logout)

	log.Fatal(http.ListenAndServe(":8000", router))
}
