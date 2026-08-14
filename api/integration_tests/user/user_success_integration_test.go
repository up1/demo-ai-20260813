package user_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin_ValidCredentials_ReturnsToken(t *testing.T) {
	truncateUsers(t)
	userID := seedUser(t, "somkiat", "12345678")

	resp, body := doLogin(t, "somkiat", "12345678")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "success", body.Status)
	assert.Equal(t, "Login successful", body.Message)
	assert.Equal(t, userID, body.Data.UserID)
	assert.Equal(t, "somkiat", body.Data.Username)
	assert.Equal(t, "123 Main St, City, Country", body.Data.Address)
	assert.NotEmpty(t, body.Data.Token)
}

func TestLogin_ValidCredentials_ResetsFailedAttemptsAfterPriorFailure(t *testing.T) {
	truncateUsers(t)
	seedUser(t, "somkiat", "12345678")

	_, failed := doLogin(t, "somkiat", "wrongpassword")
	assert.Equal(t, "error", failed.Status)

	resp, body := doLogin(t, "somkiat", "12345678")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "success", body.Status)
	assert.NotEmpty(t, body.Data.Token)
}
