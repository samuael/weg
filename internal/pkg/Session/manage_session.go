package session

import (
	"fmt"
	"net/http"
	"time"

	entity "github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("This is my Secret Key For Encryption ")
var cookiename = "snrs"

// Cookiehandler representing the Cookie Handler Interface
type Cookiehandler struct {
}

// NewCookieHandler representing New Cookie thing
func NewCookieHandler() *Cookiehandler {
	return &Cookiehandler{}
}

// SaveSession to save the Session in the User  Session Header
func (sessh *Cookiehandler) SaveSession(writer http.ResponseWriter, session *entity.Session, host string) bool {
	// Declare the expiration time of the token
	Succesfull := false
	expirationTime := time.Now().Add(24 * time.Hour)
	session.StandardClaims = jwt.StandardClaims{
		// In JWT, the expiry time is expressed as unix milliseconds
		ExpiresAt: expirationTime.Unix(),
		// HttpOnly:  true,
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, session)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return Succesfull
	}
	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	cookie := http.Cookie{
		Name:    entity.SessionName,
		Value:   tokenString,
		Expires: expirationTime,
		// Domain:   host,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(writer, &cookie)
	Succesfull = true
	return Succesfull
}

// DeleteSession representing del
func (sessh *Cookiehandler) DeleteSession(writer http.ResponseWriter, request *http.Request) bool {
	session := entity.Session{}
	expirationTime := time.Now().Add( -2400 * time.Hour)
	session.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, session)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return false
	}
	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	cookie := http.Cookie{
		Name:    entity.SessionName,
		Value:   tokenString,
		Expires: expirationTime,
		// Domain:   host,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(writer, &cookie)
	return true
}

// SaveLang to save the Session in the User  Session Header
func (sessh *Cookiehandler) SaveLang(writer http.ResponseWriter, lang string, host string) bool {
	// Declare the expiration time of the token
	expirationTime := time.Now().Add(2400 * time.Hour)
	// Declare the token with the algorithm used for signing, and the claims
	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	if lang =="" {
		return false 
	}
	cookie := http.Cookie{
		Name:    entity.LanguageSessionName,
		Value:   lang,
		Expires: expirationTime,
		// Domain:   host,
		SameSite: http.SameSiteDefaultMode,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(writer, &cookie)
	return true
}

// DeleteLang representing del
func (sessh *Cookiehandler) DeleteLang(writer http.ResponseWriter, request *http.Request) bool {
	cookie := http.Cookie{
		Name: entity.LanguageSessionName,
		// Domain:  ":8080",
		Path:    "/",
		Expires: time.Now().Add(-10 * time.Second),
	}
	http.SetCookie(writer, &cookie)
	return true
}

// GetSession returns a Session Struct Having the Data of the User
func (sessh *Cookiehandler) GetSession(request *http.Request) *entity.Session {
	cookie, err := request.Cookie(entity.SessionName)
	defer recover()
	if err != nil {
		return nil
	}
	tknStr := cookie.Value
	session := &entity.Session{}
	tkn, err := jwt.ParseWithClaims(tknStr, session, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil
		}
		return nil
	}
	if !tkn.Valid {
		return nil 
	}
	return session
}

// GetLang returns a Session Struct Having the Data of the User
func (sessh *Cookiehandler) GetLang(request *http.Request) string {
	cookie, err := request.Cookie(entity.LanguageSessionName )
	defer recover()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	tknStr := cookie.Value
	return tknStr
}

// RandomToken random token Generator for CSRF and related technologies
func (sessh *Cookiehandler) RandomToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString

}

// ValidateForm representing the Form Value
func (sessh *Cookiehandler) ValidateForm(tokenstring string) bool {
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return false
	}
	return true
}
