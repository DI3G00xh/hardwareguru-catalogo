package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/config"
	"github.com/DI3G00xh/hardwareguru-catalogo/internal/repository"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	// Conectar a la base de datos y ejecutar AutoMigrate
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Printf("ADVERTENCIA: No se pudo conectar a PostgreSQL (%v). Verifique que el servicio esté arriba.", err)
	} else {
		log.Println("Conexión a PostgreSQL establecida y AutoMigrate completado exitosamente.")
		_ = db
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("Error al escuchar en puerto %s: %v", cfg.Port, err)
	}

	grpcServer := grpc.NewServer()

	// Canal para apagado elegante (Graceful Shutdown)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Servidor gRPC de Catálogo escuchando en el puerto :%s", cfg.Port)
		if err := grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("Error al levantar el servidor gRPC: %v", err)
		}
	}()

	<-stopChan
	log.Println("Señal de apagado recibida. Deteniendo servidor gRPC de manera ordenada...")
	grpcServer.GracefulStop()
	log.Println("Servidor gRPC detenido correctamente.")
}
