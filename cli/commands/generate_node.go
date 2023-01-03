package commands

import (
	"encoding/base64"
	dkg2 "github.com/bloxapp/dkg/dkg"
	"github.com/bloxapp/dkg/utils"
	kyber "github.com/drand/kyber/share/dkg"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var GenerateNodeCommand = func() *cobra.Command {
	initGenerateNodeCmd()
	return generateNodeCommand
}()

var generateNodeCommand = &cobra.Command{
	Use:   "generate-node",
	Short: "Generates a session DKG node with a unique public key",
	Run: func(cmd *cobra.Command, args []string) {
		logger := utils.Logger().With(zap.String("cmd", "generate-node"))

		operatorID, err := cmd.Flags().GetInt(OperatorIDFlag)
		if err != nil {
			panic(err.Error())
		}

		encryptionPKStr, err := cmd.Flags().GetString(OperatorEncryptionKey)
		if err != nil {
			panic(err.Error())
		}
		byts, err := base64.StdEncoding.DecodeString(encryptionPKStr)
		if err != nil {
			panic(err.Error())
		}
		encryptionPK, err := dkg2.PemToPublicKey(byts)
		if err != nil {
			panic(err.Error())
		}

		node := dkg2.NewNode(uint32(operatorID), dkg2.Suite.G1().(kyber.Suite), encryptionPK)

		logger.Info("", zap.Any("node", node))
	},
}

func initGenerateNodeCmd() {
	setOperatorIDFlag(generateNodeCommand)
	setOperatorEncryptionKeyFlag(generateNodeCommand)
}
