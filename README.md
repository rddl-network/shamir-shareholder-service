# shamir-shareholder-service
This service serves the purpose of encrypting and storing a mnemonic that is part of a Shamir Secret Split. To ensure secure communication it utilizes mutual TLS. It offers two routes:

- GET `/mnemonic`
- POST `/mnemonic`

Clients will need a certificate that is signed by the same Root CA. Unauthorized requests will be blocked.

**Curl example:**
```bash
# GET /mnemonic
curl --cert client.crt --key client.key --cacert ca.crt https://localhost:8080/mnemonic

# POST /mnemonic
curl --cert client.crt --key client.key --cacert ca.crt -d "{\"mnemonic\":\"mnemonic phrase\"}" https://localhost:8080/mnemonic
```

## Execution

This service relies on TLS certificates.
For production use you NEED to create your own.
However for convenience you'll find TLS certificates in this repo.
Just change `certs-path` to `'./example/certs/'` (see next section).

The service can be executed via the following go command without having it previously built:
```bash
go run cmd/shamir-shareholder-service/main.go
```

## Configuration
The service needs to be configured via the ```./app.toml``` file or environment variables. The defaults are
```
certs-path = './certs/'
db-path = './data'
key-phrase = 'keyphrase'
log-level = 'error'
service-host = 'localhost'
service-port = 8080
```
