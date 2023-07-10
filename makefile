run_server:
	nodemon --exec go run main.go --signal SIGTERM

migrate_up:
	migrate -path database/migrations/ -database "postgres://cyicsnej:pkzAVWH--U1AcE4IMQb3rfWvyE-gYK22@arjuna.db.elephantsql.com/cyicsnej?sslmode=disable" -verbose up

migrate_down:
	migrate -path database/migrations/ -database "postgres://cyicsnej:pkzAVWH--U1AcE4IMQb3rfWvyE-gYK22@arjuna.db.elephantsql.com/cyicsnej?sslmode=disable" -verbose down

gen_migrate:
	migrate create -ext sql -dir database/migrations -seq $(name)