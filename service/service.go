package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
)

type ShamirService struct {
	router *gin.Engine
	db     *leveldb.DB
}

func NewShamirService(router *gin.Engine, db *leveldb.DB) *ShamirService {
	service := &ShamirService{router: router, db: db}
	service.configureRouter()
	service.registerRoutes()
	return service
}

func (ss *ShamirService) Run() (err error) {
	return ss.router.Run()
}

func (ss *ShamirService) EncryptMnemonic(mnemonic string, keyPhrase string) (ciphered []byte, err error) {
	gcmInstance, err := ss.getGCMInstance(keyPhrase)
	if err != nil {
		return
	}

	nonce := make([]byte, gcmInstance.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return
	}

	ciphered = gcmInstance.Seal(nonce, nonce, []byte(mnemonic), nil)
	return
}

func (ss *ShamirService) DecryptMnemonic(ciphered []byte, keyPhrase string) (mnemonic string, err error) {
	gcmInstance, err := ss.getGCMInstance(keyPhrase)
	if err != nil {
		return
	}

	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]

	mnemonicBytes, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		return
	}

	return string(mnemonicBytes), err
}

func (ss *ShamirService) getGCMInstance(keyPhrase string) (gcmInstance cipher.AEAD, err error) {
	sha256Hash := sha256.Sum256([]byte(keyPhrase))
	aesBlock, err := aes.NewCipher(sha256Hash[:])
	if err != nil {
		return
	}

	return cipher.NewGCM(aesBlock)
}
