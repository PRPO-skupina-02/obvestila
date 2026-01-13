# Microservice obvestila

Email notification service for the PRPO Cinema App project.

## Env vars

Check out .env.example for example values

| ENV               | Description                                      |
| ----------------- | ------------------------------------------------ |
| LOG_LEVEL         | Log level (DEBUG, INFO, WARN, ERROR)             |
| TZ                | Timezone                                         |
| RABBITMQ_URL      | URL of the rabbitmq service                      |
| RESEND_API_KEY    | API key for Resend service                       |
| RESEND_FROM_EMAIL | Email address from which emails should originate |

## Running

Run the application via

```shell
godotenv go run main.go
```

Regenerate swagger docs via

```shell
make docs
```

Run all application tests via

```shell
make test
```
