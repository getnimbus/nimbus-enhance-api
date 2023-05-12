# Go Nimeth

## How to use

1. Run `make` to build executable file
2. Add `.env` file

```env
ENV=local
DEBUG=no
MIGRATION=no
GORM_DSN=host=localhost port=5432 user=user password=password dbname=nimbus sslmode=disable TimeZone=UTC
```

3. Run executable file