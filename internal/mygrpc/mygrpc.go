package mygrpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/reflection"
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

	counters, gauges := sortMetrics(in.Metrics)

	if errUpdate := s.storage.UpdateCounters(ctx, counters); errUpdate != nil {
		return &response, status.Error(codes.Aborted, "wrong data")
	}

	if errUpdate := s.storage.UpdateGauges(ctx, gauges); errUpdate != nil {
		return &response, status.Error(codes.Aborted, "wrong data")
	}

	return &response, nil
}

func PrepareServer(appStorage storage.Repository) *grpc.Server {
	s := grpc.NewServer(grpc.UnaryInterceptor(serverIPInterceptor))

	pb.RegisterMetricsServer(s, &MetricsServer{
		storage: appStorage,
	})

	reflection.Register(s)

	return s
}

func sortMetrics(request []*pb.Metrica) (counters, gauges []metrics.Metrics) {
	for _, metrica := range request {
		switch metrica.Mtype {
		case pb.Metrica_COUNTER:
			v := metrica.GetDelta()
			m := metrics.Metrics{
				ID:    metrica.GetId(),
				MType: "counter",
				Delta: &v,
			}
			counters = append(counters, m)
		case pb.Metrica_GAUGE:
			v := metrica.GetValue()
			m := metrics.Metrics{
				ID:    metrica.GetId(),
				MType: "gauge",
				Value: &v,
			}
			gauges = append(gauges, m)
		}
	}

	return counters, gauges
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
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		grpc.WithUnaryInterceptor(clientIPInterceptor),
	)
}
