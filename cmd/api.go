package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/a1mart/kafkaesque/internal/generated/messaging"
	"github.com/a1mart/kafkaesque/internal/server"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type APIService struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// List of services with their names and corresponding Swagger JSON file names
var services = []APIService{
	{Name: "Messaging", URL: "messaging.swagger.json"},
	{Name: "Huginn", URL: "simple.swagger.json"},
	{Name: "Stripe", URL: "stripe.swagger.json"},
}

func swaggerServicesHandler(w http.ResponseWriter, r *http.Request) {
	// baseURL := "http://localhost:8080/swagger/"
	// for i := range services {
	// 	services[i].URL = baseURL + services[i].URL
	// }

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(services); err != nil {
		http.Error(w, "Failed to encode services", http.StatusInternalServerError)
	}
}

func runGRPCServer(network string, addr string) *grpc.Server {
	// Set up the gRPC server
	lis, err := net.Listen(network, addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Create server instance with topic tracking
	srv := server.NewServer(8, 3, 5)

	// grpcServer := grpc.NewServer(
	// 	grpc.ChainUnaryInterceptor(
	// 		grpcmiddleware.LoggingUnaryInterceptor,
	// 		grpcmiddleware.AuthUnaryInterceptor,
	// 	),
	// 	grpc.ChainStreamInterceptor(
	// 		grpcmiddleware.LoggingStreamInterceptor,
	// 		grpcmiddleware.AuthStreamInterceptor,
	// 	),
	// )

	// Register the service implementation (use &server{} to pass the instance)
	messaging.RegisterMessagingServiceServer(s, srv)
	messaging.RegisterAdminServiceServer(s, srv)
	// Register reflection service on gRPC server
	reflection.Register(s)

	// Start the server
	log.Printf("gRPC server listening on %s\n", addr)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	return s
}

func runHTTPServer(grpcAddr string, httpAddr string) *http.Server {
	// Set up the HTTP gateway
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	ctx := context.Background()
	err := messaging.RegisterMessagingServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register MLService: %v", err)
	}

	err = messaging.RegisterAdminServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register AuthService: %v", err)
	}

	// Start the HTTP server with the REST API
	log.Printf("HTTP server listening on %s", httpAddr)
	// Create a multiplexer for HTTP routes
	mainMux := http.NewServeMux()
	mainMux.Handle("/", mux) // Register gRPC Gateway routes

	// Dynamically serve Swagger JSON files
	for _, service := range services {
		servicePath := service.URL // Example: "mlservice.json"
		mainMux.HandleFunc("/swagger/"+servicePath, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join("./pkg/swagger", servicePath))
		})
	}

	// Register dynamic service endpoint
	mainMux.HandleFunc("/swagger/services", swaggerServicesHandler)

	// Serve Swagger UI
	mainMux.Handle("/swagger/", http.StripPrefix("/swagger", http.FileServer(http.Dir("./pkg/swagger/swagger-ui"))))

	server := &http.Server{Addr: httpAddr, Handler: mainMux}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve HTTP server: %v", err)
		}
	}()
	return server
}

func main() {
	grpcAddr := "localhost:50051"
	httpAddr := "localhost:8080"

	grpcServer := runGRPCServer("tcp", grpcAddr)
	httpServer := runHTTPServer(grpcAddr, httpAddr)

	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received shutdown signal: %v", sig)

	// Create a timeout context for cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the HTTP server
	log.Println("Shutting down HTTP server...")
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down HTTP server: %v", err)
	}

	// Shutdown the gRPC server
	log.Println("Shutting down gRPC server...")
	grpcServer.Stop() // More reliable immediate shutdown

	log.Println("Servers shut down successfully")
}
