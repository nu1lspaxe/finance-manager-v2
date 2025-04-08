# Users

Fill the values for the `configs/config.json`

```json
{
    "postgres": {
        "link": ""
    },
    "redis": {
        "link": "localhost:6379"
    },
    "ports": {
        "grpc": "9443",
        "http": "8443"
    },
    "certs": {
        "path": {
            "cert": "certs/cert.pem",
            "key": "certs/key.pem"
        }
    },
    "jwt": {
        "secret": ""
    },
    "fintech": {
        "url": "https://localhost:8855"
    },
    "kafka": {
        "addr": "localhost:9094",
        "topic": "bank-accounts"
    }
}
```