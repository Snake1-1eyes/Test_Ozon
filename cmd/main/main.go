// package main

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	shortener "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener"
// 	shortenerDelivery "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener/delivery/http"
// 	inmemoryRepo "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener/repo/inmemory"
// 	postgresRepo "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener/repo/postgres"
// 	shortenerUsecase "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener/usecase"

// 	"github.com/gorilla/mux"
// 	"github.com/jackc/pgx/v4"
// 	"github.com/jackc/pgx/v4/pgxpool"
// 	"github.com/joho/godotenv"
// )

// func init() {
// 	if err := godotenv.Load(); err != nil {
// 		fmt.Println("Error loading .env:", err)
// 	}
// }

// func main() {
// 	var repoInstance shortener.ShortenerRepo
// 	var closeFunc func()

// 	storageType := os.Getenv("STORAGE")
// 	if storageType == "inmemory" {
// 		fmt.Println("Using In-Memory storage")
// 		repoInstance = inmemoryRepo.NewInMemoryRepo()
// 		closeFunc = func() {}
// 	} else {
// 		fmt.Println("Using PostgreSQL storage")
// 		db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
// 		if err != nil {
// 			fmt.Println("Error connecting to database:", err)
// 			return
// 		}
// 		conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
// 		if err != nil {
// 			fmt.Println("Error connecting to database:", err)
// 			db.Close()
// 			return
// 		}
// 		repoInstance = postgresRepo.NewShortenerRepo(db, *conn)
// 		closeFunc = func() {
// 			db.Close()
// 		}
// 	}
// 	defer closeFunc()

// 	// Инициализируем Usecase и HTTP Delivery
// 	uc := shortenerUsecase.NewShortenerUsecase(repoInstance)
// 	hd := shortenerDelivery.NewShortenerHandler(uc)

// 	r := mux.NewRouter().PathPrefix("/api").Subrouter()
// 	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusNotFound)
// 		fmt.Println("not found")
// 	})
// 	http.Handle("/", r)

// 	// Настройка эндпоинтов
// 	shortenerRouter := r.PathPrefix("/v1").Subrouter()
// 	{
// 		shortenerRouter.Handle("/create", http.HandlerFunc(hd.CreateAndSaveShortLink)).
// 			Methods(http.MethodPost, http.MethodOptions)
// 		shortenerRouter.Handle("/get", http.HandlerFunc(hd.GetShortLink)).
// 			Methods(http.MethodGet, http.MethodOptions)
// 	}

// 	// Настройка graceful shutdown
// 	signalCh := make(chan os.Signal, 1)
// 	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

// 	server := &http.Server{
// 		Handler:           r,
// 		Addr:              ":8080",
// 		ReadTimeout:       10 * time.Second,
// 		WriteTimeout:      10 * time.Second,
// 		ReadHeaderTimeout: 10 * time.Second,
// 	}

// 	go func() {
// 		if err := server.ListenAndServe(); err != nil {
// 			fmt.Println("Server Stopped:", err)
// 		}
// 	}()
// 	fmt.Println("Server started on :8080")

// 	<-signalCh
// 	fmt.Println("Shutdown signal received")

// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	if err := server.Shutdown(ctx); err != nil {
// 		fmt.Println("Server shutdown failed:", err)
// 	}
// }

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
	postgresRepo "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener/repo/postgres"
	shortenerUsecase "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener/usecase"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	// Импорт сгенерированного кода от proto
	shortenerpb "github.com/Snake1-1eyes/Test_Ozon/internal/pkg/shortener/delivery/grpc/shortenerpb"
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
		db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Println("Error connecting to database:", err)
			return
		}
		conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Println("Error connecting to database:", err)
			db.Close()
			return
		}
		repoInstance = postgresRepo.NewShortenerRepo(db, *conn)
		closeFunc = func() {
			db.Close()
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

	// Создание и запуск gRPC сервера
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

	// Запуск gRPC сервера
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
