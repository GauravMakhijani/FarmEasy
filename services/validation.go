package services

import (
	"FarmEasy/api"
	"context"
	"errors"
	"net/http"
	"regexp"

	"github.com/dgrijalva/jwt-go"
)

func ValidateFarmerPhone(phone string) (err error) {
	re := regexp.MustCompile(`^[0-9]{10}$`)
	if !re.MatchString(phone) {
		err = errors.New("invalid phone number")
	}
	return
}

func ValidateFarmerEmail(email string) (err error) {
	re := regexp.MustCompile(`^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$`)
	if !re.MatchString(email) {
		err = errors.New("invalid email")
	}
	return
}

func ValidateBookingslots(slots []uint) (err error) {
	if len(slots) == 0 {
		err = errors.New("no slots selected")
	}
	for _, v := range slots {
		if v > 24 || v < 1 {
			err = errors.New("invalid slot selected")
		}
	}
	return
}

func ValidateBookingDate(date string) (err error) {
	re := regexp.MustCompile(`^([0-9]{4})-([0-9]{2})-([0-9]{2})$`)
	if !re.MatchString(date) {
		err = errors.New("invalid date")
	}
	return
}

func ValidateJWT(tokenString string) (farmerId uint, err error) {
	tokenObject, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}
		return secretKey, nil
	})
	if err != nil {
		return
	}
	claims, ok := tokenObject.Claims.(jwt.MapClaims)
	if !ok {
		return
	}
	farmerId = uint(claims["user_id"].(float64))
	return
}

func ValidateUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			api.Response(w, http.StatusUnauthorized, api.Message{Msg: "Authorization header required"})
			return
		}
		farmerId, err := ValidateJWT(authHeader)
		if err != nil {
			api.Response(w, http.StatusUnauthorized, api.Message{Msg: "Token is invalid"})

			return
		}
		ctx := context.WithValue(r.Context(), "token", farmerId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
