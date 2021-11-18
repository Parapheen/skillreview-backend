

devmigrateup:
	migrate -path db/migrations -database "postgresql://username:password@localhost:5432/skillreview?sslmode=disable" -verbose up

devmigratedown:
	migrate -path db/migrations -database "postgresql://username:password@localhost:5432/skillreview?sslmode=disable" -verbose down

dev:
	docker-compose up --build

test:
	docker-compose -f docker-compose.test.yaml up --build --abort-on-container-exit

cleanf:
	docker-compose down --remove-orphan --volumes
	docker volume rm skillreview-backend_database_postgres_test

lint:
	golangci-lint run

format:
	gofmt -w .

retest: cleanf test

.PHONY: dev_migratedown dev_migrateup test cleanf lint format
