postgres:
	docker-compose up -d

createdb:
	docker-compose exec pgdb psql -U db_user -c 'CREATE DATABASE users'
	docker-compose exec pgdb psql -U db_user -c 'CREATE DATABASE services'
	docker-compose exec pgdb psql -U db_user -c 'CREATE DATABASE record'

migrateup:
	migrate -path db/migration/users -database "postgresql://db_user:pwd123@localhost:54320/users?sslmode=disable" -verbose up
	migrate -path db/migration/services -database "postgresql://db_user:pwd123@localhost:54320/services?sslmode=disable" -verbose up
	migrate -path db/migration/record -database "postgresql://db_user:pwd123@localhost:54320/record?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration/users -database "postgresql://db_user:pwd123@localhost:54320/users?sslmode=disable" -verbose down
	migrate -path db/migration/services -database "postgresql://db_user:pwd123@localhost:54320/services?sslmode=disable" -verbose down
	migrate -path db/migration/record -database "postgresql://db_user:pwd123@localhost:54320/record?sslmode=disable" -verbose down