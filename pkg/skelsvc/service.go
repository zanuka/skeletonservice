package skelsvc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var signingKey = []byte("secret")
var TokenAuth = jwtauth.New("HS256", signingKey, nil)

// Service describe auth service.
type Service interface {
	Health() bool
	Login(string, string) (string, error)
}

type JWTPayload struct {
	Type   string `json:"type,omitempty"`
	UserId string `json:"userId,omitempty"`
}

// AuthService implementation of the Service interface.
type AuthService struct{}

// Health implementation of the Service
func (AuthService) Health() bool {
	return true
}

// Login implementation of the Service.
func (AuthService) Login(username string, password string) (string, error) {

	// Query the db to get the user
	sqlStatement := `select person.id, person.hashed_password from person
		LEFT JOIN person_email on person_email.person_id = person.id
		WHERE person_email.email = '%v'`

	sqlStatement = fmt.Sprintf(sqlStatement, username)

	rows, err := authconfig.DBconn.Query(sqlStatement)
	defer rows.Close()
	if err != nil {
		return "", &autherror.QueryFailure
	}

	userColVal := processRows(rows)

	// did we find a user with that username/email
	if userColVal.FindValueForCol("id").ToString() == "" {
		return "", &autherror.InvalidCredentials
	}

	// check the password against the hash one we have
	hashUserPassword := userColVal.FindValueForCol("hashed_password").ToString()
	hashPass := password + authconfig.AppConfig.HashPepper

	if !checkPasswordHash(hashPass, hashUserPassword) {
		return "", &autherror.InvalidCredentials
	}

	// everything checked out make the token
	payload := JWTPayload{Type: "squishy", UserId: userColVal.FindValueForCol("id").ToString()}
	token := createToken(payload)
	return token, nil

}

func createToken(payload JWTPayload) string {

	//set the default claims
	claims := jwt.MapClaims{}

	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiryIn(claims, 24*time.Hour)

	// add our payload
	jsonValue, _ := json.Marshal(payload)
	claims["payload"] = string(jsonValue)

	_, tokenString, _ := TokenAuth.Encode(claims)

	return tokenString

}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 1)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func processRows(rows *sql.Rows) ColVal {

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtr := make([]interface{}, count)

	for rows.Next() {

		for i, _ := range columns {
			valuePtr[i] = &values[i]
		}

		rows.Scan(valuePtr...)
	}

	cv := ColVal{Columns: columns, Values: values}
	return cv

}

func interfaceToString(value interface{}) string {

	var v string

	if value != nil {

		b, ok := value.([]byte)

		if ok {
			v = string(b)
		} else {
			v = value.(string)
		}
	}

	return v
}
