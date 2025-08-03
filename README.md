# reclip-server

Server side of Reclip application.

### Command to generate self-signed certificate for server
```
openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:secp384r1 -days 365 -nodes -keyout key.pem -out cert.pem -subj "/CN=0.0.0.0" -addext "subjectAltName=IP:0.0.0.0"
```

### Command to generate client key
```
openssl rand -base64 64
```

