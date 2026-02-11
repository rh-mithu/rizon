package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserID(t *testing.T) {
	userId := uuid.New()
	tests := []struct {
		name    string
		ctx     context.Context
		want    uuid.UUID
		wantErr bool
	}{
		{
			name:    "Found",
			ctx:     context.WithValue(context.Background(), UserIDKey, userId),
			want:    userId,
			wantErr: false,
		},
		{
			name:    "NotFound",
			ctx:     context.Background(),
			want:    uuid.Nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserID(tt.ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	secret := "test-secret"
	userId := uuid.New()

	tests := []struct {
		name           string
		setupRequest   func(req *http.Request)
		expectedStatus int
	}{
		{
			name: "Success",
			setupRequest: func(req *http.Request) {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": userId.String(),
					"exp": time.Now().Add(time.Hour).Unix(),
				})
				tokenString, _ := token.SignedString([]byte(secret))
				req.Header.Set("Authorization", "Bearer "+tokenString)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "NoAuthHeader",
			setupRequest: func(req *http.Request) {
				// No header
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "InvalidFormat",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Token abc")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "InvalidSignature",
			setupRequest: func(req *http.Request) {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": userId.String(),
					"exp": time.Now().Add(time.Hour).Unix(),
				})
				tokenString, _ := token.SignedString([]byte("wrong-secret"))
				req.Header.Set("Authorization", "Bearer "+tokenString)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "ExpiredToken",
			setupRequest: func(req *http.Request) {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": userId.String(),
					"exp": time.Now().Add(-time.Hour).Unix(),
				})
				tokenString, _ := token.SignedString([]byte(secret))
				req.Header.Set("Authorization", "Bearer "+tokenString)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "NoSub",
			setupRequest: func(req *http.Request) {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"exp": time.Now().Add(time.Hour).Unix(),
				})
				tokenString, _ := token.SignedString([]byte(secret))
				req.Header.Set("Authorization", "Bearer "+tokenString)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "InvalidSubFormat",
			setupRequest: func(req *http.Request) {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": "not-a-uuid",
					"exp": time.Now().Add(time.Hour).Unix(),
				})
				tokenString, _ := token.SignedString([]byte(secret))
				req.Header.Set("Authorization", "Bearer "+tokenString)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "WrongSigningMethod",
			setupRequest: func(req *http.Request) {
				token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
					"sub": userId.String(),
					"exp": time.Now().Add(time.Hour).Unix(),
				})
				// SigningMethodNone uses UnsafeAllowNoneSignatureType but SignedString ignores key
				tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
				req.Header.Set("Authorization", "Bearer "+tokenString)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			tt.setupRequest(req)
			rec := httptest.NewRecorder()

			middleware := AuthMiddleware(secret)
			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
