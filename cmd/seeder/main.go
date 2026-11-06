package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/config"
	"github.com/DI3G00xh/hardwareguru-catalogo/internal/domain"
	"github.com/DI3G00xh/hardwareguru-catalogo/internal/repository"
	"gorm.io/gorm/clause"
)

func main() {
	filePath := flag.String("file", "db/seeds/components.json", "Ruta al archivo JSON de componentes")
	flag.Parse()

	log.Printf("Iniciando proceso de seeding desde '%s'...", *filePath)

	cfg := config.Load()
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Error al conectar a PostgreSQL: %v", err)
	}

	data, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Error al leer archivo de semillas (%s): %v", *filePath, err)
	}

	var componentes []domain.Componente
	if err := json.Unmarshal(data, &componentes); err != nil {
		log.Fatalf("Error al deserializar JSON de componentes: %v", err)
	}

	if len(componentes) == 0 {
		log.Println("No se encontraron componentes en el archivo de semillas.")
		return
	}

	// Inserción idempotente con clausula OnConflict (upsert)
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&componentes)

	if result.Error != nil {
		log.Fatalf("Error al insertar componentes en la base de datos: %v", result.Error)
	}

	log.Printf("Seeding completado exitosamente: %d componentes procesados/actualizados.", len(componentes))

	conteoPorTipo := make(map[domain.TipoComponente]int)
	piezasClaveEncontradas := make(map[string]bool)
	piezasClave := []string{"cpu-9600x", "cpu-14600k", "gpu-5060ti", "mb-b650", "ram-ddr5-32"}

	for _, c := range componentes {
		conteoPorTipo[c.Tipo]++
		for _, clave := range piezasClave {
			if c.ID == clave {
				piezasClaveEncontradas[clave] = true
			}
		}
	}

	fmt.Println("\n--- Resumen por Tipo de Componente ---")
	for tipo, count := range conteoPorTipo {
		fmt.Printf("  • %-10s: %d\n", tipo, count)
	}

	fmt.Println("\n--- Verificación de Piezas Clave ---")
	for _, clave := range piezasClave {
		if piezasClaveEncontradas[clave] {
			fmt.Printf("  [OK] %s presente en el catálogo\n", clave)
		} else {
			fmt.Printf("  [FALTA] %s no se encontró\n", clave)
		}
	}
}

