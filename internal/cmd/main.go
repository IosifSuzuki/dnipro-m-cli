package main

import (
	"context"
	"dniprom-cli/internal/client"
	"dniprom-cli/internal/command"
	"dniprom-cli/internal/container"
	"dniprom-cli/internal/model"
	"dniprom-cli/internal/service/recorder"
	"dniprom-cli/pkg/logger"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	conf, err := model.LoadConfig()
	if err != nil {
		fmt.Println("failed to load config", err.Error())
		os.Exit(1)
	}
	ctx := context.Background()
	log := logger.NewLogger(conf.GetLoggerENV())
	cont := container.NewContainer(log, conf)

	dniproClient := client.NewDniproClient(cont)
	recorder, err := recorder.NewRecorder(ctx, cont)
	if err != nil {
		log.Fatal("fail to create recorder", logger.FError(err))
		return
	}

	warrantyCommand := command.NewWarrantyCommand(
		cont,
		dniproClient,
		recorder,
	)

	rootCmd := &cobra.Command{
		Use:   "root",
		Short: "DniproM CLI tool",
	}
	warrantyCmd := &cobra.Command{
		Use:   "warranty",
		Short: "Collect warranty information",
		Long:  "Collect warranty information for products.",
		Run:   warrantyCommand.Run,
	}

	rootCmd.AddCommand(warrantyCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
