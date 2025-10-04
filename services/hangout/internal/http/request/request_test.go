package request

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type mockValidator struct {
	shouldFail bool
}

func (m *mockValidator) Validate(i interface{}) error {
	if m.shouldFail {
		return errors.New("validation failed")
	}
	v := validator.New()
	return v.Struct(i)
}

type testRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func TestBindAndValidate(t *testing.T) {
	testCases := []struct {
		name           string
		body           string
		contentType    string
		setupValidator func() echo.Validator
		expectError    bool
		expectedName   string
	}{
		{
			name:        "success: valid payload and validation",
			body:        `{"name": "John Doe", "email": "john.doe@example.com"}`,
			contentType: echo.MIMEApplicationJSON,
			setupValidator: func() echo.Validator {
				return &mockValidator{shouldFail: false}
			},
			expectError:  false,
			expectedName: "John Doe",
		},
		{
			name:        "error: binding fails on malformed json",
			body:        `{"name": "John Doe", "email": }`,
			contentType: echo.MIMEApplicationJSON,
			setupValidator: func() echo.Validator {
				return &mockValidator{shouldFail: false}
			},
			expectError: true,
		},
		{
			name:        "error: validation fails",
			body:        `{"name": "John Doe", "email": "not-an-email"}`,
			contentType: echo.MIMEApplicationJSON,
			setupValidator: func() echo.Validator {
				return &CustomValidator{validator: validator.New()}
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = tc.setupValidator()

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, tc.contentType)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			result, err := BindAndValidate[testRequest](c)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tc.expectedName, result.Name)
			}
		})
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
