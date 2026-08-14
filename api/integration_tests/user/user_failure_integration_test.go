package user_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin_WrongPassword_ReturnsInvalidCredentialError(t *testing.T) {
	truncateUsers(t)
	seedUser(t, "somkiat", "12345678")

	resp, body := doLogin(t, "somkiat", "wrongpassword")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "error", body.Status)
	assert.Equal(t, "Invalid username or password", body.Message)
}

func TestLogin_UnknownUsername_ReturnsInvalidCredentialError(t *testing.T) {
	truncateUsers(t)

	resp, body := doLogin(t, "no-such-user", "12345678")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "error", body.Status)
	assert.Equal(t, "Invalid username or password", body.Message)
}

func TestLogin_MissingUsername_ReturnsValidationError(t *testing.T) {
	truncateUsers(t)

	resp, body := doLogin(t, "", "12345678")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "error", body.Status)
	assert.Equal(t, "Please enter a valid email address.", body.Message)
}

func TestLogin_ShortPassword_ReturnsValidationError(t *testing.T) {
	truncateUsers(t)
	seedUser(t, "somkiat", "12345678")

	resp, body := doLogin(t, "somkiat", "short")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "error", body.Status)
	assert.Equal(t, "Password must be at least 8 characters long.", body.Message)
}

func TestLogin_ThirdFailedAttempt_LocksAccountFor400(t *testing.T) {
	truncateUsers(t)
	seedUser(t, "somkiat", "12345678")

	for i := 0; i < 2; i++ {
		resp, body := doLogin(t, "somkiat", "wrongpassword")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "Invalid username or password", body.Message)
	}

	resp, body := doLogin(t, "somkiat", "wrongpassword")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "error", body.Status)
	assert.Contains(t, body.Message, "locked")

	lockedResp, lockedBody := doLogin(t, "somkiat", "12345678")
	assert.Equal(t, http.StatusBadRequest, lockedResp.StatusCode)
	assert.Contains(t, lockedBody.Message, "locked")
}
