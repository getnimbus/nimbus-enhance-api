# Nimbus Enhance API

Build for Nodereal Hackathon

## How to use

1. Run `make` to build executable file
2. Add `.env` file

```env
ENV=local
DEBUG=no
MIGRATION=no
GORM_DSN=host=localhost port=5432 user=user password=password dbname=nimbus sslmode=disable TimeZone=UTC
REDIS_ADDRESS=localhost:6379
REDIS_DB=0
```

3. Run executable file

## References:

https://encodeclub.notion.site/NodeReal-b0d916d076984eb5ad16818e0ddf327c

https://dashboard.nodereal.io/my-apis