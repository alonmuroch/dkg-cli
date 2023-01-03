package commands

import "github.com/spf13/cobra"

const (
	OperatorIDFlag        = "operator-id"
	OperatorEncryptionKey = "encryption-pk"
)

func setOperatorIDFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Int(OperatorIDFlag, 1000000, "set SSV operator ID")
	cmd.MarkPersistentFlagRequired(OperatorIDFlag)
}

func setOperatorEncryptionKeyFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(OperatorEncryptionKey, "", "set SSV operator encryption public key")
	cmd.MarkPersistentFlagRequired(OperatorEncryptionKey)
}
