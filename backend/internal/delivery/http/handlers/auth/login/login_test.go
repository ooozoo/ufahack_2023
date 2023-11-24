package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"ufahack_2023/internal/delivery/http/handlers/auth/login/mocks"
	"ufahack_2023/internal/domain"
	"ufahack_2023/internal/lib/logger/handlers/slogdiscard"
)

func TestLoginHandler(t *testing.T) {
	cases := []struct {
		name       string
		username   string
		password   string
		respErrors []string
		mockError  error
		statusCode int
	}{
		{
			name:       "Success",
			username:   "admin",
			password:   "admin",
			statusCode: http.StatusOK,
		},
		{
			name:       "Empty username",
			username:   "",
			password:   "admin",
			respErrors: []string{"field username is a required field"},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Empty password",
			username:   "admin",
			password:   "",
			respErrors: []string{"field password is a required field"},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userLoginer := mocks.NewUserLoginer(t)

			if len(tc.respErrors) == 0 || tc.mockError != nil {
				userLoginer.
					On("Login", mock.Anything, tc.username, tc.password).
					Return(&domain.User{ID: uuid.Nil}, "", tc.mockError).
					Once()
			}

			handler := New(slogdiscard.NewDiscardLogger(), userLoginer)

			input := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, tc.username, tc.password)

			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.statusCode, rr.Code)

			body := rr.Body.String()

			var resp Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respErrors, resp.Errors)
		})
	}
}
