package client

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/KirillKhitev/metrics/internal/metrics"
	"github.com/KirillKhitev/metrics/internal/mygrpc"
	pb "github.com/KirillKhitev/metrics/internal/mygrpc/proto"
)

type GRPCClient struct {
	client pb.MetricsClient
	conn   *grpc.ClientConn
}

func NewGRPCClient() (*GRPCClient, error) {
	conn, err := mygrpc.PrepareClientConnection()
	if err != nil {
		return nil, err
	}

	gRPCClient := &GRPCClient{
		client: pb.NewMetricsClient(conn),
		conn:   conn,
	}

	return gRPCClient, nil
}

func (c *GRPCClient) Send(ctx context.Context, metricsData []metrics.Metrics) error {
	ctxt, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	_, err := c.client.UpdatesMetrics(ctxt, &pb.Request{
		Metrics: c.prepareDataForSend(metricsData),
	})

	return err
}

func (c *GRPCClient) prepareDataForSend(data []metrics.Metrics) []*pb.Metrica {
	result := make([]*pb.Metrica, len(data))

	for i, v := range data {
		m := &pb.Metrica{
			Id: v.ID,
		}

		switch v.MType {
		case "counter":
			m.Mtype = pb.Metrica_COUNTER
			m.Delta = *v.Delta
		case "gauge":
			m.Mtype = pb.Metrica_GAUGE
			m.Value = *v.Value
		}

		result[i] = m
	}

	return result
}

func (c *GRPCClient) Close() error {
	c.conn.Close()

	log.Println("Close gRPC-connection")

	return nil
}
