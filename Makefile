

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

.PHONY: dev_migratedown dev_migrateup