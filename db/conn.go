package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func ConnectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to " + dbname)

	return db, nil
}

func RunMigrations(db *sql.DB, migrationDir string) {
	// Criar tabela de controle de migrações, se não existir
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatal("Erro ao criar tabela de migrações:", err)
	}

	// Obter e executar as migrações
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		log.Fatalf("Erro ao ler o diretório de migrações: %v", err)
	}

	for _, file := range files {
		migrationName := file.Name()
		if !isMigrationApplied(db, migrationName) {
			applyMigration(db, migrationDir+"/"+migrationName)
		} else {
			fmt.Printf("Migração %s já foi aplicada\n", migrationName)
		}
	}
}

// Verificar se a migração já foi aplicada
func isMigrationApplied(db *sql.DB, migration string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM migrations WHERE name = $1", migration).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}

// Aplicar a migração
func applyMigration(db *sql.DB, migrationPath string) {
	// Ler o arquivo de migração
	sqlBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo %s: %v\n", migrationPath, err)
	}
	sqlStmt := string(sqlBytes)

	// Executar o SQL
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Erro ao aplicar a migração %s: %v\n", migrationPath, err)
	}

	// Registrar que a migração foi aplicada
	_, err = db.Exec("INSERT INTO migrations (name) VALUES ($1)", migrationPath)
	if err != nil {
		log.Fatalf("Erro ao registrar a migração %s: %v\n", migrationPath, err)
	}

	fmt.Printf("Migração %s aplicada com sucesso\n", migrationPath)
}
