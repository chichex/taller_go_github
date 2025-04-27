package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"parte3/api"
	"parte3/internal/user"
	"testing"
)

func TestIntegrationCreateAndGet(t *testing.T) {
	app := gin.Default()
	api.InitRoutes(app)

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	res := fakeRequest(app, req)

	require.NotNil(t, res)
	require.Equal(t, http.StatusOK, res.Code)
	require.Contains(t, res.Body.String(), "pong")

	req, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{
		"name":"Ayrton",
		"address": "Pringles",
		"nickname": "Chiche"	
	}`))

	res = fakeRequest(app, req)

	require.NotNil(t, res)
	require.Equal(t, http.StatusCreated, res.Code)

	var resUser *user.User
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &resUser))
	require.Equal(t, "Ayrton", resUser.Name)
	require.Equal(t, "Pringles", resUser.Address)
	require.Equal(t, "Chiche", resUser.NickName)
	require.Equal(t, 1, resUser.Version)
	require.NotEmpty(t, resUser.ID)
	require.NotEmpty(t, resUser.CreatedAt)
	require.NotEmpty(t, resUser.UpdatedAt)

	req, _ = http.NewRequest(http.MethodGet, "/users/"+resUser.ID, nil)

	res = fakeRequest(app, req)

	require.NotNil(t, res)
	require.Equal(t, http.StatusOK, res.Code)
}

func fakeRequest(e *gin.Engine, r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)

	return w
}
