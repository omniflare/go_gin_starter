package utils

import (
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"net/mail"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidatePassword(pass string) error {
	if len(pass) < 8 || len(pass) > 32 {
		return errors.New("password length should be between 8 and 32")
	}
	reUpperCase := regexp.MustCompile(`[A-Z]`)
	reNumber := regexp.MustCompile(`[0-9]`)
	reSpecialChar := regexp.MustCompile(`[^a-zA-Z0-9]`)

	if !reUpperCase.MatchString(pass) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !reNumber.MatchString(pass) {
		return errors.New("password must contain at least one number")
	}
	if !reSpecialChar.MatchString(pass) {
		return errors.New("password must contain at least one special character")
	}
	return nil
}

func MatchHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateOTP() string {
	const otpChars = "1234567890"
	b := make([]byte, 6)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(otpChars))))
		if err != nil {
			log.Fatal("Error while generating otp")
		}
		b[i] = otpChars[n.Int64()]
	}
	return string(b)
}
