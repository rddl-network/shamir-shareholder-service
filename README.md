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
The service can be executed via the following go command without having it previously built:
```bash
go run cmd/shamir-shareholder-service/main.go
```

## Configuration
The service needs to be configured via the ```./app.toml``` file or environment variables. The defaults are
```
db-path = './data'
service-host = 'localhost'
service-port = 8080
key-phrase='keyPhrase'
certs-path='./certs/'
```