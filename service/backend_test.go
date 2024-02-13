package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMnemonic(t *testing.T) {
	app, _, db := setupService(t)
	defer db.Close()

	cipheredMnemonic, err := app.EncryptMnemonic("mnemonic", "keyphrase")
	assert.NoError(t, err)

	err = app.PutMnemonic(cipheredMnemonic)
	assert.NoError(t, err)

	mnemonic, err := app.GetMnemonic()
	assert.NoError(t, err)
	assert.Equal(t, cipheredMnemonic, mnemonic)

	mnemonicStr, err := app.DecryptMnemonic(cipheredMnemonic, "keyphrase")
	assert.NoError(t, err)
	assert.Equal(t, "mnemonic", mnemonicStr)
}
