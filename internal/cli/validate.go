package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"golang_united_project_2023/dnscatcher/internal/catcher"
	"golang_united_project_2023/dnscatcher/internal/handlers"
)

func validateCatcher(cmd *cobra.Command, args []string) error {
	deviceName := cmd.Flag("device").Value.String()
	outputName := cmd.Flag("output").Value.String()
	quitTO := cmd.Flag("quit-timeout").Value.String()
	if err := validateCatcherOutput(outputName); err != nil {
		return fmt.Errorf("output name validation falied: %w", err)
	}
	if err := validateCatcherDeviceName(deviceName); err != nil {
		return fmt.Errorf("device name validation falied: %w", err)
	}
	if err := validateCatcherQuitTimeout(quitTO); err != nil {
		return fmt.Errorf("quit timeout validation falied: %w", err)
	}
	return nil
}

func validateCatcherQuitTimeout(qto string) error {
	errMsg := fmt.Errorf("wrong value: %s. Expecting int seconds > 0", qto)
	to, err := strconv.Atoi(qto)
	if err != nil || to < 1 {
		return errMsg
	}
	return nil
}

func validateCatcherOutput(outputName string) error {
	if !handlers.IsKnownOUT(outputName) {
		return fmt.Errorf("unknown output %s. Allowed values: %s", outputName,
			strings.Join(handlers.ListKnownOutputs(), ","))
	}
	return nil
}

func validateCatcherDeviceName(deviceName string) error {
	if deviceName == "" {
		return nil
	}
	devices, err := catcher.ListAvailableDevices()
	if err != nil {
		return err
	}
	names := make([]string, len(devices))
	for i, v := range devices {
		if deviceName == v.Name() {
			return nil
		}
		names[i] = v.Name()
	}
	return fmt.Errorf("unknown device: %s. Allowed values:\n %s", deviceName,
		strings.Join(names, "\n"))
}
