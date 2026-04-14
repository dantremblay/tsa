package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	JWK            []byte
	ValidSignature bool
}

type MyCustomClaims struct {
	Admin bool `json:"admin"`
	jwt.RegisteredClaims
}

// GetAudience returns the first audience string, for backward compatibility.
func (c MyCustomClaims) GetFirstAudience() string {
	if len(c.Audience) > 0 {
		return c.Audience[0]
	}
	return ""
}

// GetExpiresAtUnix returns the expiration time as a Unix timestamp.
func (c MyCustomClaims) GetExpiresAtUnix() int64 {
	if c.ExpiresAt != nil {
		return c.ExpiresAt.Unix()
	}
	return 0
}

func New(jwk []byte, valid bool) *Config {
	return &Config{
		JWK:            jwk,
		ValidSignature: valid,
	}
}

func (c *Config) Create(audience, issuer string, admin bool, ttl int) (string, error) {
	now := time.Now()

	expireMinutes := 5
	if admin {
		expireMinutes = 1440
		if ttl > 0 {
			expireMinutes = ttl
		}
	}

	claims := MyCustomClaims{
		admin,
		jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * time.Duration(expireMinutes))),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(c.JWK)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func (c *Config) GetToken(tokenString string) (*jwt.Token, error) {
	opts := []jwt.ParserOption{}
	if !c.ValidSignature {
		opts = append(opts, jwt.WithoutClaimsValidation())
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return c.JWK, nil
	}, opts...)
	if err != nil && c.ValidSignature {
		return nil, err
	}

	return token, nil
}

func (c *Config) GetStandardClaims(tokenString string) (jwt.RegisteredClaims, error) {
	token, err := c.GetToken(tokenString)
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}

	claims := token.Claims.(jwt.MapClaims)

	cl := jwt.RegisteredClaims{}

	if v, ok := claims["aud"]; ok {
		switch a := v.(type) {
		case string:
			cl.Audience = jwt.ClaimStrings{a}
		case []interface{}:
			for _, s := range a {
				if str, ok := s.(string); ok {
					cl.Audience = append(cl.Audience, str)
				}
			}
		}
	}
	if v, ok := claims["exp"]; ok {
		cl.ExpiresAt = jwt.NewNumericDate(time.Unix(int64(v.(float64)), 0))
	}
	if v, ok := claims["jti"]; ok {
		cl.ID = v.(string)
	}
	if v, ok := claims["iat"]; ok {
		cl.IssuedAt = jwt.NewNumericDate(time.Unix(int64(v.(float64)), 0))
	}
	if v, ok := claims["iss"]; ok {
		cl.Issuer = v.(string)
	}
	if v, ok := claims["nbf"]; ok {
		cl.NotBefore = jwt.NewNumericDate(time.Unix(int64(v.(float64)), 0))
	}
	if v, ok := claims["sub"]; ok {
		cl.Subject = v.(string)
	}

	return cl, nil
}

func (c *Config) GetCustomClaims(tokenString string) (MyCustomClaims, error) {
	token, err := c.GetToken(tokenString)
	if err != nil {
		return MyCustomClaims{}, err
	}

	claims := token.Claims.(jwt.MapClaims)

	admin := false
	if v, ok := claims["admin"]; ok {
		admin = v.(bool)
	}

	stdclaims, err := c.GetStandardClaims(tokenString)
	if err != nil {
		return MyCustomClaims{}, err
	}

	mcc := MyCustomClaims{admin, stdclaims}

	return mcc, nil
}
