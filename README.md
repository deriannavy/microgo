# Microgo
Microservice

### Create migrations
migrate create -seq -ext sql -dir ./migrate/migrations (\<verb\>_\<resource\>:create_accounts)

### Run migrations
migrate -path=./migrate/migrations -database=<uri> up