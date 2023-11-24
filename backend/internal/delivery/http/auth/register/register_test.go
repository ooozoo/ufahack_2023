package register

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

	"ufahack_2023/internal/delivery/http/auth/register/mocks"
	"ufahack_2023/internal/lib/logger/handlers/slogdiscard"
)

func TestRegisterHandler(t *testing.T) {
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
			statusCode: http.StatusCreated,
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

			userRegisterMock := mocks.NewUserRegister(t)

			if len(tc.respErrors) == 0 || tc.mockError != nil {
				userRegisterMock.
					On("Register", mock.Anything, tc.username, tc.password).
					Return(uuid.Nil, tc.mockError).
					Once()
			}

			handler := New(slogdiscard.NewDiscardLogger(), userRegisterMock)

			input := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, tc.username, tc.password)

			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte(input)))
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
