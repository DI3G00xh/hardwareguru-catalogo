package tests

import (
	"context"
	"net"
	"testing"

	"github.com/DI3G00xh/hardwareguru-catalogo/internal/service"
	grpcTransport "github.com/DI3G00xh/hardwareguru-catalogo/internal/transport/grpc"
	hardwarepb "github.com/DI3G00xh/hardwareguru-catalogo/proto/hardware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func setupGRPCTestServer(t *testing.T) (hardwarepb.CatalogoServiceClient, func()) {
	repo := setupSeededDB(t)
	catalogoSvc := service.NewCatalogoService(repo)
	compatibilidadSvc := service.NewCompatibilidadService(repo)
	handler := grpcTransport.NewHandler(catalogoSvc, compatibilidadSvc)

	lis := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	hardwarepb.RegisterCatalogoServiceServer(grpcServer, handler)

	go func() {
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			t.Errorf("error inesperado en servidor gRPC de prueba: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("error conectando cliente gRPC en memoria: %v", err)
	}

	client := hardwarepb.NewCatalogoServiceClient(conn)

	cleanup := func() {
		conn.Close()
		grpcServer.GracefulStop()
	}

	return client, cleanup
}

func TestGRPC_ListarComponentes_Integracion(t *testing.T) {
	client, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := context.Background()

	// Listar todos los componentes del catálogo oficial
	res, err := client.ListarComponentes(ctx, &hardwarepb.ListarComponentesRequest{})
	if err != nil {
		t.Fatalf("error al listar componentes vía gRPC: %v", err)
	}

	if len(res.GetComponentes()) < 25 {
		t.Fatalf("se esperaban al menos 25 componentes en el catálogo, obtenidos %d", len(res.GetComponentes()))
	}

	// Filtrar por GPUs
	resGPU, err := client.ListarComponentes(ctx, &hardwarepb.ListarComponentesRequest{Tipo: "gpu"})
	if err != nil {
		t.Fatalf("error al filtrar GPUs vía gRPC: %v", err)
	}
	if len(resGPU.GetComponentes()) == 0 {
		t.Fatalf("se esperaba al menos una GPU en el catálogo")
	}
	for _, gpu := range resGPU.GetComponentes() {
		if gpu.GetTipo() != "gpu" {
			t.Errorf("componente no es tipo gpu: %s (%s)", gpu.GetId(), gpu.GetTipo())
		}
	}
}

func TestGRPC_ObtenerSpecs_Integracion(t *testing.T) {
	client, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := context.Background()

	// Componente existente con especificaciones y benchmarks
	specs, err := client.ObtenerSpecs(ctx, &hardwarepb.ObtenerSpecsRequest{Id: "cpu-9600x"})
	if err != nil {
		t.Fatalf("error al obtener specs de cpu-9600x: %v", err)
	}

	if specs.GetNombre() != "Ryzen 5 9600X" {
		t.Errorf("esperaba 'Ryzen 5 9600X', obtenido '%s'", specs.GetNombre())
	}
	if specs.GetSocket() != "AM5" {
		t.Errorf("esperaba socket 'AM5', obtenido '%s'", specs.GetSocket())
	}
	if specs.GetTdp() != 65 {
		t.Errorf("esperaba TDP 65, obtenido %d", specs.GetTdp())
	}
	if specs.GetEspecificacionesJson() == "" || specs.GetEspecificacionesJson() == "{}" {
		t.Errorf("especificaciones_json no debería estar vacío para cpu-9600x")
	}

	// Componente inexistente -> NotFound
	_, err = client.ObtenerSpecs(ctx, &hardwarepb.ObtenerSpecsRequest{Id: "componente-fantasma-999"})
	if status.Code(err) != codes.NotFound {
		t.Errorf("esperaba código NotFound, obtenido %v", status.Code(err))
	}
}

func TestGRPC_VerificarCompatibilidad_Integracion(t *testing.T) {
	client, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := context.Background()

	// Solicitud de build válida a través de gRPC
	reqValida := &hardwarepb.VerificarCompatibilidadRequest{
		Build: &hardwarepb.BuildRequest{
			CpuId:         "cpu-9600x",
			MotherboardId: "mb-b650",
			RamId:         "ram-ddr5-32",
			GpuId:         "gpu-5060ti",
			PsuId:         "psu-650",
			RamCantidad:   2,
		},
	}

	resValida, err := client.VerificarCompatibilidad(ctx, reqValida)
	if err != nil {
		t.Fatalf("error inesperado en llamada RPC VerificarCompatibilidad: %v", err)
	}

	if !resValida.GetCompatible() {
		t.Errorf("la build válida debería retornar compatible=true, errores: %v", resValida.GetErrores())
	}
	if resValida.GetConsumoEstimadoW() != 354 {
		t.Errorf("esperaba consumo de 354W, obtenido %d", resValida.GetConsumoEstimadoW())
	}

	// Solicitud con incompatibilidad de socket vía gRPC
	reqInvalida := &hardwarepb.VerificarCompatibilidadRequest{
		Build: &hardwarepb.BuildRequest{
			CpuId:         "cpu-14600k", // LGA1700
			MotherboardId: "mb-b650",     // AM5
			RamId:         "ram-ddr5-32",
			PsuId:         "psu-650",
		},
	}

	resInvalida, err := client.VerificarCompatibilidad(ctx, reqInvalida)
	if err != nil {
		t.Fatalf("error inesperado en llamada RPC: %v", err)
	}

	if resInvalida.GetCompatible() {
		t.Errorf("la build con choque de sockets debería retornar compatible=false")
	}
	if len(resInvalida.GetErrores()) == 0 {
		t.Errorf("se esperaba al menos 1 error reportado")
	}
}

