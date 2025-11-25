package database

import (
	"log"

	"github.com/Dooform/test-data-api/config"
	"github.com/Dooform/test-data-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	dsn := config.GetDSN()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
}

func Migrate() {
	log.Println("Running migrations...")
	migrator := DB.Migrator()

	if !migrator.HasIndex(&models.AdministrativeBoundary{}, "idx_name1") {
		migrator.CreateIndex(&models.AdministrativeBoundary{}, "name1")
	}
	if !migrator.HasIndex(&models.AdministrativeBoundary{}, "idx_name2") {
		migrator.CreateIndex(&models.AdministrativeBoundary{}, "name2")
	}
	if !migrator.HasIndex(&models.AdministrativeBoundary{}, "idx_name3") {
		migrator.CreateIndex(&models.AdministrativeBoundary{}, "name3")
	}

	// Add full-text search support
	AddFullTextSearch()

	// Add trigram index for infix search
	AddTrigramIndex()

	log.Println("Migrations complete.")
}

func AddFullTextSearch() {
	log.Println("Adding full-text search support...")

	// Add search_vector column
	if !DB.Migrator().HasColumn(&models.AdministrativeBoundary{}, "search_vector") {
		log.Println("Creating search_vector column...")
		if err := DB.Exec(`ALTER TABLE administrative_boundaries ADD COLUMN search_vector tsvector`).Error; err != nil {
			log.Println("Error creating search_vector column:", err)
		}
	}

	// Add index on search_vector
	if !DB.Migrator().HasIndex(&models.AdministrativeBoundary{}, "idx_search_vector") {
		log.Println("Creating index on search_vector...")
		if err := DB.Exec(`CREATE INDEX idx_search_vector ON administrative_boundaries USING gin(search_vector)`).Error; err != nil {
			log.Println("Error creating index on search_vector:", err)
		}
	}

	// Create a function to update the search_vector column
	log.Println("Creating or replacing update_search_vector function...")
	if err := DB.Exec(`
	CREATE OR REPLACE FUNCTION update_search_vector() RETURNS trigger AS $$
	BEGIN
		NEW.search_vector := 
			to_tsvector('thai', COALESCE(NEW.name1, '')) ||
			to_tsvector('thai', COALESCE(NEW.name2, '')) ||
			to_tsvector('thai', COALESCE(NEW.name3, '')) ||
			to_tsvector('english', COALESCE(NEW.name_eng1, '')) ||
			to_tsvector('english', COALESCE(NEW.name_eng2, '')) ||
			to_tsvector('english', COALESCE(NEW.name_eng3, ''));
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`).Error; err != nil {
		log.Println("Error creating or replacing update_search_vector function:", err)
	}

	// Create a trigger to update the search_vector on insert or update
	log.Println("Creating trigger...")
	if err := DB.Exec(`
	DO $$
	BEGIN
		IF NOT EXISTS (
			SELECT 1
			FROM pg_trigger
			WHERE tgname = 'tsvectorupdate'
		) THEN
			CREATE TRIGGER tsvectorupdate
			BEFORE INSERT OR UPDATE ON administrative_boundaries
			FOR EACH ROW EXECUTE FUNCTION update_search_vector();
		END IF;
	END;
	$$;
	`).Error; err != nil {
		log.Println("Error creating trigger:", err)
	}

	// Update existing rows
	log.Println("Updating existing rows...")
	if err := DB.Exec(`UPDATE administrative_boundaries SET search_vector = to_tsvector('thai', COALESCE(name1, '')) || to_tsvector('thai', COALESCE(name2, '')) || to_tsvector('thai', COALESCE(name3, '')) || to_tsvector('english', COALESCE(name_eng1, '')) || to_tsvector('english', COALESCE(name_eng2, '')) || to_tsvector('english', COALESCE(name_eng3, ''))`).Error; err != nil {
		log.Println("Error updating existing rows:", err)
	}
	log.Println("Full-text search support added.")
}

func AddTrigramIndex() {
	log.Println("Adding trigram index for infix search...")
	if !DB.Migrator().HasIndex(&models.AdministrativeBoundary{}, "trgm_idx_administrative_boundaries_names") {
		log.Println("Creating trigram index...")
		if err := DB.Exec(`CREATE INDEX trgm_idx_administrative_boundaries_names ON administrative_boundaries USING gin ((name1 || ' ' || name2 || ' ' || name3 || ' ' || name_eng1 || ' ' || name_eng2 || ' ' || name_eng3) gin_trgm_ops)`).Error; err != nil {
			log.Println("Error creating trigram index:", err)
		}
	}
	log.Println("Trigram index support added.")
}