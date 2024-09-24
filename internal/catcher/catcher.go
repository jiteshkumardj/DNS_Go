package catcher

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Catcher struct {
	deviceName string
	handle     *pcap.Handle
	cancelFunc context.CancelFunc
	writer     PacketWriter
	logger     *slog.Logger
}

type PacketWriter interface {
	Write(gopacket.Packet, *layers.DNS) error
	MyName() string
	Close() error
}

func New(deviceName string, snapshotLength int32, promiscuous bool, timeout time.Duration, p PacketWriter, logger *slog.Logger) (*Catcher, error) {
	handle, err := pcap.OpenLive(deviceName, snapshotLength, promiscuous, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to open device: %w", err)
	}
	return &Catcher{
		deviceName: deviceName,
		handle:     handle,
		writer:     p,
		logger:     logger,
	}, nil
}

func (c *Catcher) Listen() {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancelFunc = cancel
	c.logger.Info("will send result to", "writer", c.writer.MyName())
	fmt.Printf("Start listening on the device: %s\n", c.deviceName)
	c.logger.Info("start listening on the device", "device", c.deviceName)
	go c.catchPackets(ctx)
}

func (c *Catcher) catchPackets(ctx context.Context) {
	packetSource := gopacket.NewPacketSource(c.handle, c.handle.LinkType())
	for {
		select {
		case packet := <-packetSource.Packets():
			processPacket(c.writer, packet)
		case <-ctx.Done():
			return
		}
	}
}

func processPacket(printer PacketWriter, packet gopacket.Packet) error {
	if layer := packet.Layer(layers.LayerTypeDNS); layer != nil {
		dnsLayer, ok := layer.(*layers.DNS)
		if !ok {
			return nil
		}
		printer.Write(packet, dnsLayer)
	}
	return nil
}

type closerFunc func(*Catcher)

func (c *Catcher) Stop(ctx context.Context) {
	wg := sync.WaitGroup{}
	closers := []closerFunc{
		closeCatcher,
		closeWriter,
	}
	wg.Add(len(closers))
	for _, closer := range closers {
		go func(fnc closerFunc) {
			timeOutedCloser(ctx, fnc, c)
			wg.Done()
		}(closer)
	}
	wg.Wait()
}

func timeOutedCloser(ctx context.Context, closer closerFunc, catcher *Catcher) {
	doneCh := make(chan struct{}, 1)
	go func() {
		closer(catcher)
		doneCh <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		catcher.logger.Info("force quit by timeout")
		return
	case <-doneCh:
		return
	}
}

func closeCatcher(c *Catcher) {
	c.logger.Info("stopping the DNS catcher...")
	c.cancelFunc()
	c.handle.Close()
	c.logger.Info("dns catcher stopped")
}

func closeWriter(c *Catcher) {
	c.logger.Info("stopping the PacketWriter catcher...")
	c.writer.Close()
	c.logger.Info("packet writer catcher stopped")
}
