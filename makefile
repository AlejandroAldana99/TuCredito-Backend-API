# migrations
migrations-up:
	migrate -source ./migrations -database $(DATABASE_URL) up

# migrations-down
migrations-down:
	migrate -source ./migrations -database $(DATABASE_URL) down

# create-migration
create-migration:
	migrate create -ext sql -dir ./migrations -seq $(NAME)