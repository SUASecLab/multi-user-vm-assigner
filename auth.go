package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func decodeWebsiteToken(tokenString string) (bool, map[string]interface{}) {
	secret := []byte(websiteKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: failed to decode token: unexpected signing method: %v",
				time.Now().Format(time.Stamp), token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		log.Printf("%s: failed to parse token: %s", time.Now().Format(time.Stamp), err)
		return false, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, claims
	} else {
		log.Printf("%s Invalid token: %s", time.Now().Format(time.Stamp), err)
	}

	return false, nil
}

func generateJitsiToken(claims map[string]interface{}) string {
	name, ok := claims["name"].(string)
	if !ok {
		name = "Me"
	}

	moderator, ok := claims["moderator"].(bool)
	if !ok {
		moderator = false
	}

	currentTime := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"context": map[string]map[string]string{
			"user": {
				"name": name,
			},
		},
		"nbf":       currentTime.Unix(),
		"aud":       "jitsi",
		"iss":       jitsiIssuer,
		"room":      "*",
		"moderator": moderator,
		"iat":       currentTime.Unix(),
		"exp":       currentTime.Add(time.Duration(time.Hour * 3)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jitsiKey))
	if err != nil {
		log.Printf("Could not generate Jitsi token: %s\n", err)
	}

	return tokenString
}
