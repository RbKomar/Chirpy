package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestCheckJWTValidation(t *testing.T) {
	userID_1 := uuid.New()
	userID_2 := uuid.New()

	token_secret_1 := "itsasecret"
	token_secret_2 := "lessofasecret"
	token_string_1, _ := MakeJWT(userID_1, token_secret_1, 10*time.Second)
	token_string_2, _ := MakeJWT(userID_2, token_secret_2, 10*time.Second)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		userID      uuid.UUID
		wantErr     bool
		wantWrongID bool
	}{
		{
			name:        "correct JWT",
			tokenString: token_string_1,
			tokenSecret: token_secret_1,
			userID:      userID_1,
			wantErr:     false,
			wantWrongID: false,
		},
		{
			name:        "incorrect tokens",
			tokenString: token_string_1,
			tokenSecret: token_secret_2,
			userID:      userID_1,
			wantErr:     true,
			wantWrongID: true,
		},
		{
			name:        "incorrect ID",
			tokenString: token_string_2,
			tokenSecret: token_secret_2,
			userID:      userID_1,
			wantErr:     false,
			wantWrongID: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf(`ValidateJWT() err=%v, wantErr %v`, err, tt.wantErr)
			}
			if (id != tt.userID) != tt.wantWrongID {
				t.Errorf(`Got wrong ID from jwt is: %v, want: %v`, id, tt.userID)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPasswordHashing(t *testing.T) {
	password := "pass"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 1)
	want := string(hashed)
	to_test, err_internal := HashPassword(password)

	if err != nil || err_internal != nil {
		t.Errorf(`One of errors happened, err: %v, err in internal func %v`, err, err_internal)
	}
	if to_test == want {
		t.Errorf(`Password was incorrectly hashed, is %s, want %s`, to_test, want)
	}
}
