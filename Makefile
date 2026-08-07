MIGRATE := migrate
MIGRATIONS := internal/storage/migrations
DB ?= Lattice.db
DBURL := sqlite3://$(DB)

.PHONY: db-up db-down db-drop db-status db-new

db-up: ## Apply all pending migrations
	$(MIGRATE) -path $(MIGRATIONS) -database "$(DBURL)" up

db-down: ## Roll back the most recent migration
	$(MIGRATE) -path $(MIGRATIONS) -database "$(DBURL)" down 1

db-status: ## Show current migration version
	$(MIGRATE) -path $(MIGRATIONS) -database "$(DBURL)" version

db-drop: ## Drop all tables
	$(MIGRATE) -path $(MIGRATIONS) -database "$(DBURL)" drop

db-new: ## Create a new migration: make db-new NAME=add_customers
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS) -seq "$(NAME)"