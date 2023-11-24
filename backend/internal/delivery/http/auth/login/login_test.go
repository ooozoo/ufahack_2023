package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"ufahack_2023/internal/delivery/http/auth/login/mocks"
	"ufahack_2023/internal/lib/logger/handlers/slogdiscard"
)

func TestLoginHandler(t *testing.T) {
	cases := []struct {
		name       string
		username   string
		password   string
		respErrors []string
		mockError  error
	}{
		{
			name:     "Success",
			username: "admin",
			password: "admin",
		},
		{
			name:       "Empty username",
			username:   "",
			password:   "admin",
			respErrors: []string{"field username is a required field"},
		},
		{
			name:       "Empty password",
			username:   "admin",
			password:   "",
			respErrors: []string{"field password is a required field"},
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
					Return("", tc.mockError).
					Once()
			}

			handler := New(slogdiscard.NewDiscardLogger(), userLoginer)

			input := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, tc.username, tc.password)

			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respErrors, resp.Errors)
		})
	}
}
