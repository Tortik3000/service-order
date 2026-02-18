package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tortik3000/service-order/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"net/http"

	generatedMenu "github.com/Tortik3000/service-order/generated/api/menu"
	generatedOrder "github.com/Tortik3000/service-order/generated/api/order"
	generatedUser "github.com/Tortik3000/service-order/generated/api/user"
	menuHandler "github.com/Tortik3000/service-order/internal/handlers/menu"
	orderHandler "github.com/Tortik3000/service-order/internal/handlers/order"
	userHandler "github.com/Tortik3000/service-order/internal/handlers/user"
	menuRepoImpl "github.com/Tortik3000/service-order/internal/repository/postgres/menu"
	orderRepoImpl "github.com/Tortik3000/service-order/internal/repository/postgres/order"
	userRepoImpl "github.com/Tortik3000/service-order/internal/repository/postgres/user"
	"github.com/Tortik3000/service-order/internal/repository/transactor"
	menuUC "github.com/Tortik3000/service-order/internal/usecase/menu"
	orderUC "github.com/Tortik3000/service-order/internal/usecase/order"
	userUC "github.com/Tortik3000/service-order/internal/usecase/user"
	metricsHandler "github.com/Tortik3000/service-order/pkg/handlers/metrics"
	"github.com/Tortik3000/service-order/pkg/logger"
	"github.com/Tortik3000/service-order/pkg/metrics"
	grpcruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	googleGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:password@postgres:5432/db?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	zapLogger, _ := zap.NewDevelopment()
	appLogger := logger.NewZap(zapLogger)
	db.SetupPostgres(pool, zapLogger)

	if err := pool.Ping(ctx); err != nil {
		appLogger.Fatal("failed to ping database", logger.Error(err))
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		appLogger.Fatal("failed to listen", logger.Error(err))
	}

	txManager := transactor.New(pool)

	userRepo := userRepoImpl.New(txManager)
	menuRepo := menuRepoImpl.New(txManager)
	orderRepo := orderRepoImpl.New(txManager)

	// Usecases
	mUC := menuUC.NewUseCase(menuRepo)
	uUC := userUC.NewUseCase(userRepo)
	oUC := orderUC.NewUseCase(orderRepo, menuRepo, txManager)

	// Handlers
	mH := menuHandler.NewMenuHandler(mUC)
	uH := userHandler.NewUserHandler(uUC)
	oH := orderHandler.NewOrderHandler(oUC)

	s := googleGRPC.NewServer()
	generatedMenu.RegisterMenuServiceServer(s, mH)
	generatedUser.RegisterUserServiceServer(s, uH)
	generatedOrder.RegisterOrderServiceServer(s, oH)

	reflection.Register(s)

	httpServer := &http.Server{
		Addr: ":8081",
	}

	metricsMdw := metrics.New(appLogger)
	mHandler := metricsHandler.New()

	go func() {
		mux := grpcruntime.NewServeMux()
		opts := []googleGRPC.DialOption{googleGRPC.WithTransportCredentials(insecure.NewCredentials())}
		err := generatedMenu.RegisterMenuServiceHandlerFromEndpoint(ctx, mux, "0.0.0.0:50051", opts)
		if err != nil {
			appLogger.Fatal("failed to register menu handler", logger.Error(err))
		}
		err = generatedUser.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "0.0.0.0:50051", opts)
		if err != nil {
			appLogger.Fatal("failed to register user handler", logger.Error(err))
		}
		err = generatedOrder.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, "0.0.0.0:50051", opts)
		if err != nil {
			appLogger.Fatal("failed to register order handler", logger.Error(err))
		}

		// Apply metrics middleware to gateway mux
		httpHandler := metricsMdw.Metrics(mux)

		finalMux := http.NewServeMux()
		finalMux.Handle("/", httpHandler)
		finalMux.HandleFunc("/metrics", mHandler.GetMetrics)

		httpServer.Handler = finalMux
		appLogger.Info("gateway listening at :8081")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("gateway listen error", logger.Error(err))
		}
	}()

	go func() {
		appLogger.Info("grpc server listening at " + lis.Addr().String())
		if err := s.Serve(lis); err != nil {
			appLogger.Fatal("failed to serve", logger.Error(err))
		}
	}()

	<-ctx.Done()
	appLogger.Info("shutting down servers...")

	const gracefulShutdownTimeout = 5 * time.Second
	shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("http server shutdown error", logger.Error(err))
	}
	s.GracefulStop()
	appLogger.Info("servers exited")
}
