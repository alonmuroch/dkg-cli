package commands

import "github.com/spf13/cobra"

const (
	OutputFolder = "output"

	OperatorIDFlag            = "operator-id"
	OperatorEncryptionKeyFlag = "encryption-pk"
	PasswordFlag              = "password"
)

func setOperatorIDFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Int(OperatorIDFlag, 1000000, "set SSV operator ID")
	cmd.MarkPersistentFlagRequired(OperatorIDFlag)
}

func setOperatorEncryptionKeyFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(OperatorEncryptionKeyFlag, "", "set SSV operator encryption public key")
	cmd.MarkPersistentFlagRequired(OperatorEncryptionKeyFlag)
}

func setOutputFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(OutputFolder, "", "set output folder")
	cmd.MarkPersistentFlagRequired(OutputFolder)
}

func setPasswordFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(PasswordFlag, "", "set secret encryption password")
	cmd.MarkPersistentFlagRequired(PasswordFlag)
}
