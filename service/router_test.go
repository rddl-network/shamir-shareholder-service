package service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	log "github.com/rddl-network/go-logger"
	"github.com/rddl-network/shamir-shareholder-service/service"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

func setupService(t *testing.T) (app *service.ShamirService, router *gin.Engine, db *leveldb.DB) {
	router = gin.Default()
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		t.Fatal("Error opening in-memory LevelDB: ", err)
	}
	app = service.NewShamirService(router, db, log.GetLogger(log.DEBUG))
	return
}

func TestGetMnemonicRoute(t *testing.T) {
	app, router, db := setupService(t)
	defer db.Close()

	ciphered, err := app.EncryptMnemonic("mnemonic", "keyphrase")
	assert.NoError(t, err)

	err = app.PutMnemonic(ciphered)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/mnemonic", nil)
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resBody service.MnemonicBody
	err = json.Unmarshal(w.Body.Bytes(), &resBody)
	assert.NoError(t, err)
	assert.Equal(t, "mnemonic", resBody.Mnemonic)
}

func TestPostMnemonicRoute(t *testing.T) {
	_, router, db := setupService(t)
	defer db.Close()

	tests := []struct {
		desc    string
		reqBody service.MnemonicBody
		code    int
	}{
		{
			desc: "valid request",
			reqBody: service.MnemonicBody{
				Mnemonic: "mnemonic",
			},
			code: 200,
		},
		{
			desc:    "invalid request",
			reqBody: service.MnemonicBody{},
			code:    400,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			w := httptest.NewRecorder()
			bodyBytes, err := json.Marshal(tc.reqBody)
			assert.NoError(t, err)
			req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/mnemonic", bytes.NewBuffer(bodyBytes))
			assert.NoError(t, err)
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.code, w.Code)
		})
	}
}
