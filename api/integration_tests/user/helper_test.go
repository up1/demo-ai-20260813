package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type loginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		UserID   int64  `json:"user_id"`
		Username string `json:"username"`
		Token    string `json:"token"`
	} `json:"data"`
}

// doLogin sends a POST /api/login request against the shared test app and
// decodes the JSON response.
func doLogin(t *testing.T, username, password string) (*http.Response, loginResponse) {
	t.Helper()

	body, err := json.Marshal(map[string]string{"username": username, "password": password})
	require.NoError(t, err)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := testApp.Test(req)
	require.NoError(t, err)
	t.Cleanup(func() { _ = resp.Body.Close() })

	var parsed loginResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&parsed))

	return resp, parsed
}
