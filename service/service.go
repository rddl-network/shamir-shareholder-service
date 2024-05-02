package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/rddl-network/go-utils/logger"
	"github.com/rddl-network/shamir-shareholder-service/config"
	"github.com/syndtr/goleveldb/leveldb"
)

type ShamirService struct {
	router *gin.Engine
	db     *leveldb.DB
	logger log.AppLogger
}

func NewShamirService(router *gin.Engine, db *leveldb.DB, logger log.AppLogger) *ShamirService {
	service := &ShamirService{router: router, db: db, logger: logger}
	service.configureRouter()
	service.registerRoutes()
	return service
}

func (ss *ShamirService) configureTLS(caCertFile []byte) (tlsConfig *tls.Config) {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertFile)

	return &tls.Config{
		ClientCAs:                caCertPool,
		ClientAuth:               tls.RequireAndVerifyClientCert,
		MinVersion:               tls.VersionTLS13,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
}

func (ss *ShamirService) Run() (err error) {
	cfg := config.GetConfig()
	caCertFile, err := os.ReadFile(cfg.CertsPath + "ca.crt")
	if err != nil {
		return err
	}

	tlsConfig := ss.configureTLS(caCertFile)
	server := &http.Server{
		Addr:      fmt.Sprintf("%s:%d", cfg.ServiceHost, cfg.ServicePort),
		TLSConfig: tlsConfig,
		Handler:   ss.router,
	}

	// workaround to listen on tcp4 and not tcp6
	// https://stackoverflow.com/a/38592286
	ln, err := net.Listen("tcp4", server.Addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	return server.ServeTLS(ln, cfg.CertsPath+"server.crt", cfg.CertsPath+"server.key")
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
