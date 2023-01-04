package main

import (
	"github.com/alonmuroch/dkg-cli/cli/commands"
	"github.com/alonmuroch/dkg-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Logger is the default logger
var Logger = utils.Logger()

func main() {
	rootCommand := &cobra.Command{
		Use:              "dkg",
		Short:            "SSV DKG Tool",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	}

	rootCommand.AddCommand(commands.GenerateNodeCommand)
	rootCommand.AddCommand(commands.DealCommand)

	Logger.Info("Starting SSV DKG...")
	if err := rootCommand.Execute(); err != nil {
		Logger.Fatal("failed to execute root command", zap.Error(err))
	}
}
