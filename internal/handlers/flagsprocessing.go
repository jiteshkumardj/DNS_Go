package handlers

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang_united_project_2023/dnscatcher/internal/catcher"
	"golang_united_project_2023/dnscatcher/internal/output/grpc"
	"golang_united_project_2023/dnscatcher/internal/output/stdout"
)

func selectDevicePrompt(h *RootHandler) error {
	devices, err := catcher.ListAvailableDevices()
	if err != nil {
		return err
	}

	fmt.Println("Choose a device:")
	for i, device := range devices {
		fmt.Printf("%d: %s\n", i+1, formatDevice(&device))
	}

	reader := bufio.NewReader(os.Stdin)
	v, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	i, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return err
	}

	if i > len(devices) || i < 1 {
		return errors.New("invalid device number")
	}
	h.deviceName = devices[i-1].Name()
	return nil
}

func getOutput(outputName outputType) (catcher.PacketWriter, error) {
	var out catcher.PacketWriter
	var err error
	switch outputName {
	case outputSTDOUT: //others will be set here. Don't forget to validate
		out, err = stdout.New()
	case outputGRPC:
		out, err = grpc.New()
	default:
		out, _ = stdout.New()
	}
	if err != nil {
		return nil, fmt.Errorf("output: %s : %w", outputName, err)
	}
	if out == nil {
		return nil, errNoOutputGiven
	}
	return out, nil
}

func formatDevice(device *catcher.NetworkDevice) string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%s (%s), IPs: ", device.Name(), device.Description()))
	for _, a := range device.Addresses() {
		b.WriteString(a.IP().String())
	}
	return b.String()
}
