# Bank System

Fill values for the `configs/config.json`

```json
{
    "postgres": {
        "connection_string": ""
    },
    "redis": {
        "addr": "",
        "password": "",
        "db": 0
    },
    "server": {
        "port": 8855
    },
    "openrouter": {
        "url": "",
        "api_key": "",
        "model": {
            "llama_3-1_8b": "meta-llama/llama-3.1-8b-instruct:free"
        },
        "message": {
            "deposit": "(Do not give amount) generate a short sentence (as category) to describe a reason that you deposit from bank today",
            "withdraw": "(Do not give amount) generate a short sentence (as category) to describe a reason that you withdraw from bank today"
        },
        "role": {
            "user": "user",
            "system": "system",
            "asistant": "assistant",
            "tool": "tool"
        }
    },
    "certs": {
        "path": {
            "cert": "certs/cert.pem",
            "key": "certs/key.pem"
        }
    }
}
```