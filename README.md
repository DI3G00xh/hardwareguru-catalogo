# Microservicio de Catálogo (`hardwareguru-catalogo`)

Servicio interno de alto rendimiento desarrollado en **Go (1.22+)** para la plataforma **Hardware Guru**. Gestiona el catálogo de componentes de hardware (CPUs, placas madre, memorias RAM, GPUs y fuentes de poder), sus especificaciones técnicas y el motor de validación de compatibilidad de armado de computadoras.

---

## 1. ¿Cómo Funciona este Microservicio?

El microservicio de Catálogo **no expone endpoints REST directamente al cliente web**. Actúa como un servicio interno conectado a PostgreSQL y expone una interfaz **gRPC (HTTP/2 + Protocol Buffers)**.

El **API Gateway / BFF (NestJS)** asume el rol de Backend-For-Frontend intermediario: recibe las peticiones REST desde el frontend en React (`GET /componentes`, `GET /componentes/:id`), gestiona CORS, autenticación y formato de errores estándar (`ApiErrorShape`), y se comunica con este microservicio vía gRPC en el puerto `50051`.

### Topología y Flujo de Comunicación

```text
+-------------------------------------------------------+
|                 Frontend (React)                      |
+---------------------------+---------------------------+
                            |
                            | HTTP / REST (JSON)
                            | GET /componentes
                            | GET /componentes/:id
                            v
+-------------------------------------------------------+
|             API Gateway / BFF (NestJS)                |
|  - Valida esquemas de entrada y CORS                  |
|  - Aplica formato de error estándar (ApiErrorShape)   |
|  - Transforma JSON REST <-> Protobuf gRPC             |
+---------------------------+---------------------------+
                            |
                            | gRPC (HTTP/2 + Protobuf) :50051
                            | ListarComponentes()
                            | ObtenerSpecs()
                            | VerificarCompatibilidad()
                            v
+-------------------------------------------------------+
|       Microservicio de Catálogo (Go 1.22+)            |
|  - Servidor gRPC (`internal/transport/grpc`)          |
|  - Lógica de compatibilidad de hardware               |
|  - Persistencia y migraciones con GORM                |
+---------------------------+---------------------------+
                            |
                            | Conexión TCP nativa :5432
                            v
+-------------------------------------------------------+
|              PostgreSQL (Base de Datos)               |
|  - Tabla `componentes` con soporte JSONB              |
+-------------------------------------------------------+
```

---

## 2. Contratos gRPC (`proto/hardware/hardware.proto`)

La definición de los contratos e interfaces del servicio se encuentra en `proto/hardware/hardware.proto`:

```protobuf
syntax = "proto3";

package hardware;

option go_package = "github.com/DI3G00xh/hardwareguru-catalogo/proto/hardware;hardwarepb";

// Servicio gRPC de Catálogo
service CatalogoService {
  rpc ListarComponentes (ListarComponentesRequest) returns (ListarComponentesResponse);
  rpc ObtenerSpecs (ObtenerSpecsRequest) returns (ComponenteDetalleResponse);
  rpc VerificarCompatibilidad (VerificarCompatibilidadRequest) returns (CompatibilidadResponse);
}

message ListarComponentesRequest {
  string tipo = 1; // "cpu", "placa", "ram", "gpu", "fuente" o vacío para todos
}

message ComponenteResumen {
  string id = 1;
  string tipo = 2;
  string nombre = 3;
  string fabricante = 4;
  int32 tdp = 5;
  string socket = 6;
  string tipo_ram = 7;
  int32 potencia_w = 8;
  string origen = 9;
}

message ListarComponentesResponse {
  repeated ComponenteResumen componentes = 1;
}

message ObtenerSpecsRequest {
  string id = 1;
}

message ComponenteDetalleResponse {
  string id = 1;
  string tipo = 2;
  string nombre = 3;
  string fabricante = 4;
  int32 tdp = 5;
  string socket = 6;
  string tipo_ram = 7;
  int32 potencia_w = 8;
  string origen = 9;
  string especificaciones_json = 10; // JSON serializado como string
  string benchmarks_json = 11;       // JSON serializado como string
}

message BuildRequest {
  string cpu_id = 1;
  string motherboard_id = 2;
  string ram_id = 3;
  string gpu_id = 4;
  string psu_id = 5;
}

message VerificarCompatibilidadRequest {
  BuildRequest build = 1;
}

message CompatibilidadResponse {
  bool compatible = 1;
  repeated string errores = 2;
  int32 consumo_estimado_w = 3;
}
```

### Compilar el contrato Protobuf

Para regenerar el código Go a partir del archivo `.proto`:

```bash
make proto
```

O directamente mediante `protoc`:

```bash
protoc --proto_path=proto/hardware \
       --go_out=proto/hardware --go_opt=paths=source_relative \
       --go-grpc_out=proto/hardware --go-grpc_opt=paths=source_relative \
       proto/hardware/hardware.proto
```

---

## 3. Estructura de Directorios

```text
hardwareguru-catalogo/
├── cmd/
│   ├── seeder/
│   │   └── main.go                    # CLI de seeding de 25+ piezas iniciales (upsert idempotente)
│   └── server/
│       └── main.go                    # Entrypoint del microservicio (servidor gRPC y Graceful Shutdown)
├── db/
│   └── seeds/
│       └── components.json            # Dataset oficial con especificaciones y benchmarks
├── internal/
│   ├── config/                        # Carga de variables de entorno y DSN de conexión
│   │   ├── config.go
│   │   └── config_test.go
│   ├── domain/                        # Entidades del modelo de dominio, tipos y errores
│   │   ├── componente.go              # Entidad GORM con tipos JSONB personalizados
│   │   ├── componente_test.go
│   │   ├── errors.go
│   │   └── seeds_test.go              # Validación de integridad del dataset oficial
│   ├── repository/                    # Acceso a datos con GORM y conexión a PostgreSQL
│   │   └── postgres.go                # Inicialización de DB, pool de conexiones y AutoMigrate
│   ├── service/                       # Reglas de negocio y compatibilidad de hardware
│   └── transport/
│       └── grpc/                      # Handlers de los RPCs y mappers de dominio a protobuf
├── proto/
│   └── hardware/
│       └── hardware.proto             # Contrato Protobuf gRPC
├── .env.example                       # Plantilla de variables de entorno
├── docker-compose.yml                 # PostgreSQL 16 para entorno local
├── go.mod                             # Definición del módulo y dependencias Go
├── go.sum                             # Checksums de dependencias
├── Makefile                           # Comandos de compilación, ejecución y tests
└── README.md
```

---

## 4. Inicialización y Puesta en Marcha

### Requisitos Previos

* **Go 1.22+** (o superior)
* **Docker** y **Docker Compose**
* Compilador de Protobuf `protoc` (opcional, solo para regenerar contratos)

### Paso 1: Configurar Variables de Entorno

Crear el archivo `.env` en la raíz del proyecto (basado en `.env.example`):

```bash
cp .env.example .env
```

Contenido por defecto:

```env
PORT=50051
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=hardwareguru_catalogo
DB_SSLMODE=disable
```

### Paso 2: Iniciar Base de Datos Local

Levantar el contenedor de PostgreSQL mediante Docker Compose:

```bash
docker compose up -d
```

Esto iniciará PostgreSQL 16 Alpine en el puerto `5432` con volumen persistente (`postgres_data`).

### Paso 3: Sembrar Componentes Iniciales (Seeder)

Ejecutar el script de seeding para poblar la base de datos con los 28 componentes oficiales (CPUs Ryzen/Intel, GPUs RTX/RX, Placas AM5/LGA1700, RAMs DDR4/DDR5 y Fuentes):

```bash
make seed
# O alternativamente:
go run cmd/seeder/main.go
```

El script utiliza cláusula `OnConflict{UpdateAll: true}`, por lo que es totalmente idempotente y puede ejecutarse repetidas veces sin duplicar registros.

### Paso 4: Iniciar el Servidor gRPC

Levantar el microservicio:

```bash
make run
# O alternativamente:
go run cmd/server/main.go
```

El servidor ejecutará `AutoMigrate` automáticamente sobre la base de datos y quedará escuchando en el puerto `:50051` con soporte para apagado limpio ante señales `SIGINT` o `SIGTERM`.

### Paso 5: Ejecutar Pruebas Unitarias

Para correr las pruebas unitarias y verificar la cobertura:

```bash
make test
# O alternativamente:
go test -v -cover ./...
```
