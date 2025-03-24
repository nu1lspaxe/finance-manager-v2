curl -X POST https://localhost:8444/v1/records \
  --cacert certs/cert.pem \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "amount": 45.67, "record_type": "expense", "transaction_date": 1677654321, "detail": "Grocery shopping"}'