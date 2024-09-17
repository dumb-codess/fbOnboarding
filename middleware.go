package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.Header.Get("Authorization")
		logger.Info().Msg("Auth")

		if token == "" {
			writeBadResponse(ctx, errors.New("token is missing"), w)
			return
		}

		pToken, err := validateToken(token, cfg.Secret)
		if err != nil {
			logger.Error().Msgf("invalid token!")
			writeBadResponse(ctx, err, w)
			return
		}

		claims, ok := pToken.Claims.(*CustomClaims)
		if !ok && pToken.Valid {
			logger.Error().Msg("customer_id type assertion error")
			writeBadResponse(ctx, errors.New("invalid customer_id type"), w)
			return
		}

		cID := claims.CustomerID

		if ctx.Value("customer_id") == nil {
			logger.Error().Msgf("customer_id not found!")
			ctx = context.WithValue(ctx, "customer_id", int64(cID))
		}
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)

	})
}

func panicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal server error",
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}

		}()

		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info().Msg("CORS")

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "content-type, Authorization")

		next.ServeHTTP(w, r)
	})
}

func validateToken(tokenString, secret string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
