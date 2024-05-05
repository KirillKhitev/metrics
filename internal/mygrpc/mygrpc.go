package mygrpc

import (
	"context"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/logger"
	"github.com/KirillKhitev/metrics/internal/metrics"
	pb "github.com/KirillKhitev/metrics/internal/mygrpc/proto"
	"github.com/KirillKhitev/metrics/internal/storage"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer
	storage storage.Repository
}

func (s *MetricsServer) UpdatesMetrics(ctx context.Context, in *pb.Request) (*pb.UpdatesResponse, error) {
	var response pb.UpdatesResponse

	counters, gauges, err := metrics.GetMetricsFromBytes(in.Data)
	if err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		return &response, err
	}

	if errUpdate := s.storage.UpdateCounters(ctx, counters); errUpdate != nil {
		return &response, status.Error(codes.Aborted, "wrong data")
	}

	if errUpdate := s.storage.UpdateGauges(ctx, gauges); errUpdate != nil {
		return &response, status.Error(codes.Aborted, "wrong data")
	}

	return &response, nil
}

func PrepareServer(appStorage storage.Repository) *grpc.Server {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		serverIPInterceptor,
		serverSignatureInterceptor,
		serverGzipInterceptor,
		serverDecryptInterceptor,
	))

	pb.RegisterMetricsServer(s, &MetricsServer{
		storage: appStorage,
	})

	return s
}

func StartServer(s *grpc.Server) error {
	listen, err := net.Listen("tcp", flags.Args.AddrRunGRPC)
	if err != nil {
		return err
	}

	return s.Serve(listen)
}

func ShutdownServer(s *grpc.Server) {
	s.GracefulStop()

	logger.Log.Info("Shutdown gRPC-server")
}

func PrepareClientConnection() (*grpc.ClientConn, error) {
	return grpc.Dial(flags.ArgsClient.AddrRun,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			clientIPInterceptor,
			clientSignatureInterceptor,
		),
	)
}

func PrepareClient(conn *grpc.ClientConn) pb.MetricsClient {
	return pb.NewMetricsClient(conn)
}
