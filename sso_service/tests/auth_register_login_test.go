package tests

import (
	"testing"

	"github.com/liriquew/social-todo/sso_service/tests/suite"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	ssov1 "github.com/liriquew/todoprotos/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	username := gofakeit.Username()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Username: username,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUid())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Username: username,
		Password: pass,
	})
	require.NoError(t, err)

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(st.JWTSecret), nil
	})

	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, username, claims["username"].(string))
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	username := gofakeit.Username()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Username: username,
		Password: pass,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUid())

	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Username: username,
		Password: pass,
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUid())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		username    string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with Empty Password",
			username:    gofakeit.Username(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Register with Empty Username",
			username:    "",
			password:    randomFakePassword(),
			expectedErr: "username is required",
		},
		{
			name:        "Register with Both Empty",
			username:    "",
			password:    "",
			expectedErr: "username is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Username: tt.username,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)

		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		username    string
		password    string
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			username:    gofakeit.Username(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Login with Empty Username",
			username:    "",
			password:    randomFakePassword(),
			expectedErr: "username is required",
		},
		{
			name:        "Login with Both Empty Username and Password",
			username:    "",
			password:    "",
			expectedErr: "username is required",
		},
		{
			name:        "Login with Non-Matching Password",
			username:    gofakeit.Username(),
			password:    randomFakePassword(),
			expectedErr: "invalid username or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Username: gofakeit.Username(),
				Password: randomFakePassword(),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Username: tt.username,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
