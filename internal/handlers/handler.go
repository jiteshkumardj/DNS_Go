package handlers

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"golang_united_project_2023/dnscatcher/internal/catcher"
)

const (
	snapshotLen int32         = 1024
	timeout     time.Duration = 30 * time.Second
)

type dnsListener interface {
	Listen()
	Stop(context.Context)
}

type RootHandler struct {
	deviceName string
	outputName outputType
	catcher    dnsListener
	logger     *slog.Logger
}

func NewRootHandler(deviceName, output string, logger *slog.Logger) *RootHandler {
	return &RootHandler{
		deviceName: deviceName,
		outputName: outputType(output),
		logger:     logger,
	}
}

func (h *RootHandler) Start() error {
	out, err := getOutput(h.outputName)
	if err != nil {
		return err
	}
	if strings.TrimSpace(h.deviceName) == "" {
		if err := selectDevicePrompt(h); err != nil {
			return err
		}
	}
	dnsCatcher, err := catcher.New(h.deviceName, snapshotLen, false, timeout, out, h.logger)
	if err != nil {
		return err
	}
	h.catcher = dnsCatcher
	h.catcher.Listen()
	return nil
}

func (h *RootHandler) Stop(ctx context.Context) {
	if h.catcher != nil {
		h.catcher.Stop(ctx)
	}
}
