package commands

import "github.com/spf13/cobra"

const (
	OutputFolder = "output"
	NodesFolder  = "node-folder"

	OperatorIDFlag            = "operator-id"
	OperatorEncryptionKeyFlag = "encryption-pk"
	PasswordFlag              = "password"
	NonceFlag                 = "nonce"
	DKGConfigFlag             = "dkg-config"
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

func setNodesFolderFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(NodesFolder, "", "set nodes folder")
	cmd.MarkPersistentFlagRequired(NodesFolder)
}

func setPasswordFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(PasswordFlag, "", "set secret encryption password")
	cmd.MarkPersistentFlagRequired(PasswordFlag)
}

func setNonceFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(NonceFlag, "", "set dkg session nonce")
	cmd.MarkPersistentFlagRequired(NonceFlag)
}

func setDKGConfigFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().IntSlice(DKGConfigFlag, []int{}, "set dkg config params (N,T)")
	cmd.MarkPersistentFlagRequired(DKGConfigFlag)
}
