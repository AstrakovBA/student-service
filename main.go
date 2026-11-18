package main

import (
	"flag"
	"fmt"
	"net"
	"net/http" // стандартный пакет — остаётся "http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	httptransport "example.com/student-service/http" // ← алиас для локального пакета
	"example.com/student-service/implementation"
	"example.com/student-service/pb"
	"example.com/student-service/repository"
	"example.com/student-service/transport"
	grpctransport "example.com/student-service/transport/grpc"
)

func main() {
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	grpcAddr := flag.String("grpc.addr", ":9090", "gRPC listen address")
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "student-service",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	// --- Общие слои (не зависят от транспорта) ---
	repo := repository.NewMemoryRepository()
	svc := implementation.NewService(repo, logger)
	endpoints := transport.MakeEndpoints(svc)

	// --- HTTP-транспорт (используем алиас httptransport) ---
	httpHandler := httptransport.NewService(endpoints, logger)

	// --- gRPC-транспорт ---
	grpcServer := grpctransport.NewGRPCServer(endpoints)
	grpcBaseServer := grpc.NewServer()
	pb.RegisterStudentServiceServer(grpcBaseServer, grpcServer)
	reflection.Register(grpcBaseServer) // для grpcurl list/describe без proto-файла

	errs := make(chan error, 1)

	// Ожидание сигнала ОС (graceful shutdown)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// Запуск gRPC на :9090
	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		level.Error(logger).Log("during", "listen grpc", "err", err)
		os.Exit(1)
	}
	go func() {
		level.Info(logger).Log("transport", "gRPC", "addr", *grpcAddr)
		errs <- grpcBaseServer.Serve(grpcListener)
	}()

	// Запуск HTTP на :8080
	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, httpHandler)
	}()

	level.Error(logger).Log("exit", <-errs)
}
