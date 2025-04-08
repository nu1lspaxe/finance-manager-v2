# Bank Records

Fill values for the `configs/config.json`

```json
{
    "postgres": {
        "connection_string": ""
    },
    "ports": {
        "grpc": "9445",
        "http": "8445"
    },
    "certs": {
        "path": {
            "cert": "certs/cert.pem",
            "key": "certs/key.pem",
            "ca": "certs/ca.pem"
        }
    },
    "jwt": {
        "secret": ""
    },
    "fintech": {
        "url": "https://localhost:8855"
    },
    "kafka": {
        "brokers": [
            "localhost:9094"
        ],
        "topic": "bank-accounts",
        "group_id": "bank-sync-group"
    }
}
```