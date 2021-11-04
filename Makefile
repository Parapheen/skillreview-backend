

devmigrateup:
	migrate -path db/migrations -database "postgresql://username:password@localhost:5432/skillreview?sslmode=disable" -verbose up

devmigratedown:
	migrate -path db/migrations -database "postgresql://username:password@localhost:5432/skillreview?sslmode=disable" -verbose down


.PHONY: dev_migratedown dev_migrateup