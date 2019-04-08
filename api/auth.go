package api

import (
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	jwtrequest "github.com/dgrijalva/jwt-go/request"
	"github.com/metiago/zbx1/common/request"

	"github.com/metiago/zbx1/repository"
)

var (
	privKeyPath = os.Getenv("PRI_RSA")
	pubKeyPath  = os.Getenv("PUB_RSA")
	verifyKey   *rsa.PublicKey
	signKey     *rsa.PrivateKey
)

type token struct {
	Hash string `json:"token"`
}

func init() {
	absPath, _ := filepath.Abs(privKeyPath)
	signBytes, err := ioutil.ReadFile(absPath)
	checkOpenRSAFile(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	checkOpenRSAFile(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	checkOpenRSAFile(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	checkOpenRSAFile(err)
}

func checkOpenRSAFile(err error) {
	if err != nil {
		log.Printf("Error open RSA key. %v", err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	var user repository.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		request.Handle400(w, err)
		return
	}

	// FIXME: MOVE TO DI SERVICE
	u, err := repository.FindUserByUsername(user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			request.Handle401(w)
			return
		}
		log.Println(err)
		request.Handle500(w, err)
		return
	}

	// validate user credentials
	ok, err := repository.AuthUser(user.Username, user.Password)
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
	}

	if !ok {
		request.Handle401(w)
		return
	}

	hash := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), jwt.MapClaims{
		"iss": "admin",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp": time.Now().Add(time.Minute * 60).Unix(),
		"uinf": struct {
			ID   int
			Name string
			Role string
		}{u.ID, u.Name, "Member"},
	})

	tokenString, err := hash.SignedString(signKey)
	if err != nil {
		log.Fatal("Error creating token")
	}

	tokenResponse(token{tokenString}, w)
}

func authHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if r.Header.Get("Authorization") == "" {
		request.Handle403(w)
		return
	}

	extractor, err := jwtrequest.HeaderExtractor{"Authorization"}.ExtractToken(r)

	token, err := jwt.Parse(extractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return verifyKey, nil
	})

	if err != nil {
		log.Println(err)
		switch err.(type) {

		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				request.Handle406(w)
				return
			default:
				request.Handle403(w)
				return
			}
		default:
			request.Handle403(w)
			return
		}
	}
	if token.Valid {
		next(w, r)
	} else {
		request.Handle403(w)
	}
}

func tokenResponse(response token, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
