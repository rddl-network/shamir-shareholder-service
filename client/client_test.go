package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rddl-network/shamir-shareholder-service/client"
	"github.com/stretchr/testify/assert"
)

func TestGetMnemonic(t *testing.T) {
	t.Parallel()
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/mnemonic", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"mnemonic":"word1"}`))
		assert.NoError(t, err)
	}))
	defer mockServer.Close()

	c := client.NewShamirShareholderClient(mockServer.URL, mockServer.Client())
	res, err := c.GetMnemonic(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, "word1", res.Mnemonic)
}

func TestPostMnemonic(t *testing.T) {
	t.Parallel()
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/mnemonic", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	c := client.NewShamirShareholderClient(mockServer.URL, mockServer.Client())
	err := c.PostMnemonic(context.Background(), "someSecret")

	assert.NoError(t, err)
}
