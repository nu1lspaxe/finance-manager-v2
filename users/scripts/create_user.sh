curl -X POST https://localhost:8443/v1/users \
  --cacert certs/cert.pem \
  -d '{"username": "test1", "email": "test1@example.com", "password": "SECpassword123"}'

curl -X POST https://localhost:8443/v1/users \
  --cacert certs/cert.pem \
  -d '{"username": "test2", "email": "test2@example.com", "password": "SECpassword123"}'

# curl -X POST https://localhost:8443/v1/users/1/accounts/51755563204470440037 \
#   --cacert certs/cert.pem