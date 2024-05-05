package mygrpc

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/mycrypto"
	pb "github.com/KirillKhitev/metrics/internal/mygrpc/proto"
	"github.com/KirillKhitev/metrics/internal/signature"
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

func serverSignatureInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var hash string

	if flags.Args.Key == "" {
		return handler(ctx, req)
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		hashes := md.Get("HashSHA256")
		if len(hashes) > 0 {
			hash = hashes[0]
		}
	}

	if hash == "" {
		return nil, status.Error(codes.InvalidArgument, "missing or invalid hash")
	}

	hashSum := signature.GetHash(req.(*pb.Request).Data, flags.Args.Key)

	if hash != hashSum {
		return nil, status.Error(codes.Aborted, "wrong hash")
	}

	return handler(ctx, req)
}

func serverGzipInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	zr, err := gzip.NewReader(bytes.NewReader(req.(*pb.Request).Data))
	if err != nil {
		return nil, err
	}

	defer zr.Close()

	var b bytes.Buffer

	_, err = b.ReadFrom(zr)
	if err != nil {
		return nil, status.Error(codes.Aborted, "wrong data")
	}

	req.(*pb.Request).Data = b.Bytes()

	return handler(ctx, req)
}

func serverDecryptInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if flags.Args.CryptoKey == "" {
		return handler(ctx, req)
	}

	r := bytes.NewReader(req.(*pb.Request).Data)
	body, _ := io.ReadAll(r)

	bodyDecrypted, err := mycrypto.Decrypt(body, flags.Args.CryptoKey)
	if err != nil {
		return nil, status.Error(codes.Aborted, "error decrypt data")
	}

	req.(*pb.Request).Data = bodyDecrypted

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

func clientSignatureInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	if flags.ArgsClient.Key == "" {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	hashSum := signature.GetHash(req.(*pb.Request).Data, flags.ArgsClient.Key)
	md := metadata.New(map[string]string{"HashSHA256": hashSum})

	mdOld, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return fmt.Errorf("error by get outgoing context")
	}

	md = metadata.Join(mdOld, md)
	ctx = metadata.NewOutgoingContext(ctx, md)

	return invoker(ctx, method, req, reply, cc, opts...)
}
