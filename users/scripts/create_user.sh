curl -X POST https://localhost:8443/v1/users \
  --cacert certs/cert.pem \
  -d '{"username": "test1", "email": "test1@example.com", "password": "securepassword123"}'