package auth

import (
	"context"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/navaz-alani/omegle/pb/go/pb/auth"
)

type AuthService struct {
	auth.UnimplementedAuthServer
	jwtSecret string
  users     map[string]bool
}

type JwtClaims struct {
	jwt.StandardClaims
	Username string
}

func NewAuthService(secret string) (*AuthService, error) {
	if secret == "" {
		return nil, fmt.Errorf("[auth] cannot initialize - empty secret")
	}
	return &AuthService{
		jwtSecret: secret,
    users: make(map[string]bool),
	}, nil
}

func (a *AuthService) GetCert(ctx context.Context, req *auth.Request) (*auth.Cert, error) {
	var username string
	if req.GetRequestedUsername() == "" {
    for {
      if _, ok := a.users[username]; username != "" && ok {
        a.users[username] = true
        break
      } else {
        username = GenerateRandomName()
      }
    }
	} else {
      if _, ok := a.users[username]; ok {
        return nil, fmt.Errorf("[auth] username taken")
      }
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &JwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		Username: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if tokenString, err := token.SignedString(a.jwtSecret); err != nil {
		return nil, err
	} else {
		tsProto, _ := ptypes.TimestampProto(expirationTime)
		return &auth.Cert{
			Jwt:        tokenString,
			Username:   username,
			Expiration: tsProto,
		}, nil
	}

	return nil, nil
}

func (a *AuthService) jwtDecode(jwtToken string) (*jwt.Token, *JwtClaims, error) {
	claims := &JwtClaims{}
	if tkn, err := jwt.ParseWithClaims(jwtToken, claims,
		func(token *jwt.Token) (interface{}, error) {
			return a.jwtSecret, nil
		}); err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, nil, fmt.Errorf("[auth] signature validation fail")
		}
		return nil, nil, fmt.Errorf("[auth] token parse fail")
	} else if !tkn.Valid {
		return nil, nil, fmt.Errorf("[auth] token validation fail")
	} else {
		return tkn, claims, nil
	}
}

func (a *AuthService) VerifCert(ctx context.Context, c *auth.Cert) (*empty.Empty, error) {
	if _, claims, err := a.jwtDecode(c.GetJwt()); err != nil {
		return nil, err
	} else {
		if claims.Username == c.Username {
			return nil, fmt.Errorf("[auth] token validation fail")
		}
		return nil, nil
	}
}

func (a *AuthService) RenewCert(ctx context.Context, c *auth.Cert) (*auth.Cert, error) {
	if _, claims, err := a.jwtDecode(c.GetJwt()); err != nil {
		return nil, err
	} else if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return c, fmt.Errorf("[auth] premature renew request")
	} else {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		if tokenString, err := token.SignedString(a.jwtSecret); err != nil {
			return nil, err
		} else {
			tsProto, _ := ptypes.TimestampProto(time.Now().Add(5 * time.Minute))
			return &auth.Cert{
				Jwt:        tokenString,
				Username:   c.GetUsername(),
				Expiration: tsProto,
			}, nil
		}
	}
	return nil, nil
}
