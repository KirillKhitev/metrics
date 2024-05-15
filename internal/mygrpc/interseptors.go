package mygrpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/subnet"
)

func serverIPInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var ipStr string

	if flags.Args.TrustedSubnet == "" {
		return handler(ctx, req)
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		ips := md.Get("X-Real-IP")
		if len(ips) > 0 {
			ipStr = ips[0]
		}
	}

	ip := net.ParseIP(ipStr)

	if ip == nil {
		return nil, status.Error(codes.InvalidArgument, "missing or invalid IP address")
	}

	_, ipNet, err := net.ParseCIDR(flags.Args.TrustedSubnet)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	if !ipNet.Contains(ip) {
		return nil, status.Error(codes.Aborted, "forbidden")
	}

	return handler(ctx, req)
}

func clientIPInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	ip, err := subnet.GetIP()
	if err != nil {
		return err
	}

	md := metadata.New(map[string]string{"X-Real-IP": ip})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return invoker(ctx, method, req, reply, cc, opts...)
}
