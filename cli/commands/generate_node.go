package commands

import (
	"encoding/base64"
	"github.com/alonmuroch/dkg-cli/cli/commands/storage"
	dkg2 "github.com/alonmuroch/dkg-cli/dkg"
	"github.com/alonmuroch/dkg-cli/utils"
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
		dkg2.InitBLS()

		logger := utils.Logger().With(zap.String("cmd", "generate-node"))

		outputPath, err := cmd.Flags().GetString(OutputFolder)
		if err != nil {
			panic(err.Error())
		}

		password, err := cmd.Flags().GetString(PasswordFlag)
		if err != nil {
			panic(err.Error())
		}

		nonce, err := cmd.Flags().GetString(NonceFlag)
		if err != nil {
			panic(err.Error())
		}

		dkgConfig, err := cmd.Flags().GetIntSlice(DKGConfigFlag)
		if err != nil {
			panic(err.Error())
		}
		if len(dkgConfig) != 2 {
			panic("dkg config has to have 2 values")
		}

		operatorID, err := cmd.Flags().GetInt(OperatorIDFlag)
		if err != nil {
			panic(err.Error())
		}

		// decode encryption pub key
		encryptionPKStr, err := cmd.Flags().GetString(OperatorEncryptionKeyFlag)
		if err != nil {
			panic(err.Error())
		}
		byts, err := base64.StdEncoding.DecodeString(encryptionPKStr)
		if err != nil {
			panic(err.Error())
		}

		node := dkg2.NewNode(uint32(operatorID), dkg2.Suite.G1().(kyber.Suite), byts)
		node.Nonce = nonce
		node.N = dkgConfig[0]
		node.T = dkgConfig[1]

		// save secrets
		if err := storage.SaveNodeToDisk(node, outputPath, password); err != nil {
			panic(err.Error())
		}

		// save public node data
		logger.Info("", zap.Any("node", node))
	},
}

func initGenerateNodeCmd() {
	setOperatorIDFlag(generateNodeCommand)
	setOperatorEncryptionKeyFlag(generateNodeCommand)
	setOutputFlag(generateNodeCommand)
	setPasswordFlag(generateNodeCommand)
	setNonceFlag(generateNodeCommand)
	setDKGConfigFlag(generateNodeCommand)
}
