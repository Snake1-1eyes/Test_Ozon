package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	shortener "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener"
	grpcDelivery "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener/delivery/grpc"
	shortenerDelivery "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener/delivery/http"
	inmemoryRepo "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener/repo/inmemory"
	"github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener/repo/postgres"
	shortenerUsecase "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener/usecase"

	shortenerpb "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener/delivery/grpc/shortenerpb"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env:", err)
	}
}

func main() {
	var repoInstance shortener.ShortenerRepo
	var closeFunc func()

	storageType := os.Getenv("STORAGE")
	if storageType == "inmemory" {
		fmt.Println("Using In-Memory storage")
		repoInstance = inmemoryRepo.NewInMemoryRepo()
		closeFunc = func() {}
	} else {
		fmt.Println("Using PostgreSQL storage")

		// Создаем пул соединений
		pool, err := postgres.NewPostgresPool(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Println("Error creating connection pool:", err)
			return
		}

		repoInstance = postgres.NewShortenerRepo(pool)
		closeFunc = func() {
			pool.Close()
		}
	}
	defer closeFunc()

	// Инициализация Usecase и HTTP Delivery
	uc := shortenerUsecase.NewShortenerUsecase(repoInstance)
	hd := shortenerDelivery.NewShortenerHandler(uc)

	// Настройка маршрутов HTTP-сервера
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("not found")
	})
	http.Handle("/", r)
	shortenerRouter := r.PathPrefix("/v1").Subrouter()
	{
		shortenerRouter.Handle("/create", http.HandlerFunc(hd.CreateAndSaveShortLink)).
			Methods(http.MethodPost, http.MethodOptions)
		shortenerRouter.Handle("/get", http.HandlerFunc(hd.GetShortLink)).
			Methods(http.MethodGet, http.MethodOptions)
	}

	// Создание HTTP сервера
	httpServer := &http.Server{
		Handler:           r,
		Addr:              ":8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Создание и запуск GRPC сервера
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		fmt.Println("Failed to start gRPC listener:", err)
		return
	}
	grpcServer := grpc.NewServer()
	grpcSvc := grpcDelivery.NewServer(uc)
	shortenerpb.RegisterShortenerServiceServer(grpcServer, grpcSvc)

	// Канал для сигналов завершения
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Запуск HTTP сервера
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("HTTP Server Stopped:", err)
		}
	}()
	fmt.Println("HTTP server started on :8080")

	// Запуск GRPC сервера
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Println("gRPC Server Stopped:", err)
		}
	}()
	fmt.Println("gRPC server started on :9090")

	// graceful shutdown
	<-signalCh
	fmt.Println("Shutdown signal received")

	httpCtx, httpCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer httpCancel()

	if err := httpServer.Shutdown(httpCtx); err != nil {
		fmt.Println("HTTP server shutdown failed:", err)
	}

	// Грейсфул завершение gRPC сервера
	grpcServer.GracefulStop()
	fmt.Println("Servers gracefully stopped")
}
