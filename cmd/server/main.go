package main

import (
    "fmt"
    "log"
    "net"
    "os"

    "github.com/joho/godotenv"
    "google.golang.org/grpc"
)

func main() {
    _ = godotenv.Load()

    port := os.Getenv("PORT")
    if port == "" {
        port = "50051"
    }

    listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
    if err != nil {
        log.Fatalf("Error al escuchar en puerto %s: %v", port, err)
    }

    grpcServer := grpc.NewServer()

    log.Printf("Servidor gRPC de Catálogo escuchando en el puerto :%s", port)
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatalf("Error al levantar el servidor gRPC: %v", err)
    }
}
