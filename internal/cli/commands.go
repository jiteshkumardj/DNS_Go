package cli

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"golang_united_project_2023/dnscatcher/internal/handlers"
)

var logger *slog.Logger

var rootCmd = &cobra.Command{
	Use:   "catcher",
	Short: "DNS Catcher listens to a network device and outputs the DNS data into stdout",
	Args:  validateCatcher,
	Run: func(cmd *cobra.Command, args []string) {
		deviceName := cmd.Flag("device").Value.String()
		output := cmd.Flag("output").Value.String()
		quitTimeout := cmd.Flag("quit-timeout").Value.String()
		to, err := strconv.Atoi(quitTimeout)
		if err != nil {
			logger.Error(err.Error())
		}
		handler := handlers.NewRootHandler(deviceName, output, logger)
		err = handler.Start()
		if err != nil {
			logger.Error(err.Error())
		}
		exiter(*handler, time.Duration(to)*time.Second)
	},
}

func exiter(h handlers.RootHandler, shutDownTime time.Duration) {
	sc := make(chan os.Signal, 1)
	normalQuit := make(chan struct{}, 1)
	ctx, cancelFunc := context.WithTimeout(context.Background(), shutDownTime)
	defer cancelFunc()
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-sc
	go func() {
		h.Stop(ctx)
		normalQuit <- struct{}{}
	}()
	select {
	case <-sc:
		logger.Info("force quit")
		os.Exit(1)
	case <-normalQuit:
		logger.Info("shutdown complete")
		return
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	var deviceName string
	var output string
	var timeoutSecs int64
	rootCmd.Flags().StringVarP(&deviceName, "device", "d", "", "provide a name of a network device to listen on")
	rootCmd.Flags().StringVarP(&output, "output", "o", "stdout", "sets the output")
	rootCmd.Flags().Int64VarP(&timeoutSecs, "quit-timeout", "q", 5, "sets the timeout for quit")

	logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)
}
