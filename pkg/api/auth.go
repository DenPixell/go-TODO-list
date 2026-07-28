package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type signInRequest struct {
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req signInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err.Error())
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if req.Password != pass {
		writeJson(w, signInResponse{Error: "Неверный пароль"})
		return
	}

	claims := jwt.MapClaims{
		"hash": passwordHash(pass),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(pass))
	if err != nil {
		writeError(w, err.Error())
		return
	}

	writeJson(w, signInResponse{Token: signed})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtStr string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtStr = cookie.Value
			}

			valid := false
			token, err := jwt.Parse(jwtStr, func(t *jwt.Token) (any, error) {
				return []byte(pass), nil
			})
			if err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if hash, ok := claims["hash"].(string); ok && hash == passwordHash(pass) {
						valid = true
					}
				}
			}

			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}
