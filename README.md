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
  int32 ram_cantidad = 6;
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
├── .github/
│   └── workflows/
│       └── ci.yml                     # Pipeline CI/CD: Pruebas, cobertura (>=60%) y Docker Build
├── cmd/
│   ├── seeder/
│   │   └── main.go                    # CLI de seeding de piezas oficiales (upsert idempotente)
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
│   │   ├── build.go                   # Modelos de solicitud de build y resultado
│   │   ├── componente.go              # Entidad GORM con tipos JSONB y helpers de hardware
│   │   ├── componente_test.go
│   │   ├── errors.go
│   │   └── seeds_test.go              # Validación de integridad del dataset oficial
│   ├── repository/                    # Acceso a datos con GORM y conexión a PostgreSQL
│   │   ├── componente_repository.go
│   │   ├── componente_repository_test.go
│   │   └── postgres.go                # Inicialización de DB, pool de conexiones y AutoMigrate
│   ├── service/                       # Reglas de negocio y compatibilidad de hardware
│   │   ├── catalogo_service.go        # Servicio de consulta de catálogo y especificaciones
│   │   ├── catalogo_service_test.go
│   │   ├── compatibilidad_service.go  # Motor de validación de 4 reglas de hardware
│   │   └── compatibilidad_service_test.go
│   └── transport/
│       └── grpc/                      # Handlers de RPCs y mappers de dominio a Protobuf
│           ├── handler.go             # Implementación de CatalogoServiceServer
│           ├── handler_test.go
│           └── mapper.go              # Conversores entre Domain y Protobuf
├── proto/
│   └── hardware/
│       ├── hardware.pb.go             # Código Go generado por protoc-gen-go
│       ├── hardware.proto             # Contrato Protobuf gRPC
│       └── hardware_grpc.pb.go        # Stubs gRPC generados por protoc-gen-go-grpc
├── tests/
│   ├── catalogo_service_test.go       # Pruebas de integración gRPC end-to-end con bufconn
│   └── compatibilidad_test.go         # Pruebas integrales de reglas de hardware con dataset real
├── .dockerignore
├── .env.example                       # Plantilla de variables de entorno
├── docker-compose.yml                 # Orquestación local (PostgreSQL + Microservicio Catálogo)
├── Dockerfile                         # Build multi-stage optimizado (< 28 MB)
├── go.mod                             # Definición del módulo y dependencias Go
├── go.sum                             # Checksums de dependencias
├── Makefile                           # Comandos de compilación, ejecución, tests y cobertura
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

### Paso 2: Iniciar con Docker Compose

Puedes levantar todo el ecosistema local (PostgreSQL 16 + Microservicio de Catálogo en Go) con un solo comando:

```bash
docker compose up -d --build
```

O si prefieres correr PostgreSQL en contenedor y el binario de Go en tu entorno local:

```bash
# Iniciar solo PostgreSQL
docker compose up -d postgres
```

### Paso 3: Sembrar Componentes Iniciales (Seeder)

Ejecutar el script de seeding para poblar la base de datos con los 28 componentes oficiales (CPUs Ryzen/Intel, GPUs RTX/RX, Placas AM5/LGA1700, RAMs DDR4/DDR5 y Fuentes):

```bash
make seed
# O alternativamente:
go run cmd/seeder/main.go
```

El script utiliza cláusula `OnConflict{UpdateAll: true}`, por lo que es totalmente idempotente y puede ejecutarse repetidas veces sin duplicar registros.

### Paso 4: Iniciar el Servidor gRPC en Desarrollo

Levantar el microservicio localmente:

```bash
make run
# O alternativamente:
go run cmd/server/main.go
```

El servidor ejecutará `AutoMigrate` automáticamente sobre la base de datos y quedará escuchando en el puerto `:50051` con soporte para gRPC reflection y apagado limpio ante señales `SIGINT` o `SIGTERM`.

### Paso 5: Ejecutar Pruebas Unitarias, Cobertura e Integración

El repositorio cuenta con tres suites de pruebas automatizadas:

```bash
# Correr todas las pruebas con reporte de cobertura estándar
make test

# Generar reporte detallado de cobertura por función (requiere >= 60%)
make test-coverage

# Ejecutar únicamente las pruebas de integración gRPC y compatibilidad
make test-integration
```

---

## 5. Reglas del Motor de Compatibilidad

El motor en [`internal/service/compatibilidad_service.go`](internal/service/compatibilidad_service.go) valida cuatro reglas fundamentales antes de dar el visto bueno a una configuración:

1. **Coincidencia Exacta de Socket:** Comprueba que el socket del procesador coincida exactamente con el de la placa madre (`AM5 == AM5`, `LGA1700 == LGA1700`).
2. **Compatibilidad de Memoria RAM:** Verifica que la tecnología de la memoria (`DDR3`, `DDR4` o `DDR5`) sea la misma soportada por la placa base.
3. **Validación de Ranuras de RAM:** Corrobora que la cantidad física de módulos solicitados no exceda las ranuras (`slotsRam`) disponibles en la placa base (ej. placas Micro-ATX compactas con 2 ranuras frente a 4 módulos solicitados).
4. **Cálculo de Consumo y Potencia de Fuente:** Estima el consumo eléctrico del sistema completo aplicando un margen de seguridad del 20%:
   $$\text{Consumo Estimado (W)} = \lceil (\text{TDP}_{\text{cpu}} + \text{TDP}_{\text{gpu}} + 50\text{W}) \times 1.2 \rceil$$
   Verifica que la fuente de poder elegida entregue una potencia mayor o igual al consumo estimado ($\text{Potencia Fuente} \ge \text{Consumo}$).

