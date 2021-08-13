package face

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
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
func (j *Jwt) Encode(input map[string]interface{}) (string, error) {
	j.claims = input
	j.token.Claims = j.claims
	result, err := j.token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}
	return result, nil
}

//decode
func (j *Jwt) Decode(input string) (map[string]interface{}, error) {
	//parse input string
	token, err := jwt.Parse(input, j.getValidationKey)
	if err != nil {
		return nil, err
	}
	//check header
	if jwt.SigningMethodHS256.Alg() != token.Header["alg"] {
		return nil, errors.New("jwt header not matched")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("decode jwt data failed")
}

//get validate key
func (j *Jwt) getValidationKey(*jwt.Token) (interface{}, error) {
	return []byte(j.secret), nil
}
