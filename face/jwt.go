package face

import (
	"github.com/dgrijalva/jwt-go"
	"log"
)

/*
 * jwt face
 */

//face info
type Jwt struct {
	secret string //secret key string
	token *jwt.Token //jwt token instance
	claims jwt.MapClaims //jwt claims object
}

//construct
func NewJwt(secretKey string) *Jwt {
	this := &Jwt{
		secret:secretKey,
		token:jwt.New(jwt.SigningMethodHS256),
		claims:make(jwt.MapClaims),
	}
	return this
}

//encode
func (j *Jwt) Encode(input map[string]interface{}) string {
	j.claims = input
	j.token.Claims = j.claims
	result, err := j.token.SignedString([]byte(j.secret))
	if err != nil {
		log.Println("Encode jwt failed, error:", err.Error())
		return ""
	}
	return result
}

//decode
func (j *Jwt) Decode(input string) map[string]interface{} {
	//parse input string
	token, err := jwt.Parse(input, j.getValidationKey)
	if err != nil {
		log.Println("Decode ", input, " failed, error:", err.Error())
		return nil
	}
	//check header
	if jwt.SigningMethodHS256.Alg() != token.Header["alg"] {
		log.Println("Header error")
		return nil
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims
	}
	return nil
}

//get validate key
func (j *Jwt) getValidationKey(*jwt.Token) (interface{}, error) {
	return []byte(j.secret), nil
}
