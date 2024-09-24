package grpc

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"golang_united_project_2023/dnscatcher/internal/common"
)

type grpcout struct {
}

func New() (grpcout, error) {
	return grpcout{}, common.ErrNotImplemented
}

func (s grpcout) MyName() string {
	return "grpc"
}

func (s grpcout) Write(packet gopacket.Packet, dnsLayer *layers.DNS) error {
	return common.ErrNotImplemented
}

func (s grpcout) Close() error {
	return nil
}
