# Makefile for create, up and down migrations
# DATABASE_URL is the connection string like: postgres://postgres:password@localhost:5432/creditdb?sslmode=disable
# NAME is a simple string without spaces and sequence like: init_schema


# migrations
migrations-up:
	migrate -source ./migrations -database $(DATABASE_URL) up

# migrations-down
migrations-down:
	migrate -source ./migrations -database $(DATABASE_URL) down

# create-migration
create-migration:
	migrate create -ext sql -dir ./migrations -seq $(NAME)

# get the current version of the migrations
migrations-version:
	migrate -path ./migrations -database $(DATABASE_URL) version

# force a specific version of the migrations
migrations-force:
	migrate -path ./migrations -database $(DATABASE_URL) force $(VERSION)
