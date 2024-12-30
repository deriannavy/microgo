# Microgo
Microservice

### Create migrations
migrate create -seq -ext sql -dir ./migrate/migrations create_accou
nts

### Run migrations
migrate -path=./migrate/migrations -database=<uri> up