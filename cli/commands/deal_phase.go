package commands

import (
	dkg2 "github.com/alonmuroch/dkg-cli/dkg"
	"github.com/drand/kyber/share/dkg"
	"github.com/spf13/cobra"
)

var DealCommand = func() *cobra.Command {
	initDealCmd()
	return dealCommand
}()

var dealCommand = &cobra.Command{
	Use:   "deal",
	Short: "Deal shares and save to json file",
	Run: func(cmd *cobra.Command, args []string) {
		//dkg2.InitBLS()
		//
		//logger := utils.Logger().With(zap.String("cmd", "generate-node"))
		//
		//outputPath, err := cmd.Flags().GetString(OutputFolder)
		//if err != nil {
		//	panic(err.Error())
		//}
		//
		//nodesFolder, err := cmd.Flags().GetString(NodesFolder)
		//if err != nil {
		//	panic(err.Error())
		//}
		//
		//password, err := cmd.Flags().GetString(PasswordFlag)
		//if err != nil {
		//	panic(err.Error())
		//}
		//
		//// load nodes
		//nodes, err := storage.LoadNodes(nodesFolder)
		//if err != nil {
		//	logger.Fatal("", zap.Error(err))
		//}
		//
		//// load secrets
		//secrets := storage.Secrets{}
		//if err := secrets.LoadFromDisk(outputPath); err != nil {
		//	logger.Fatal("", zap.Error(err))
		//}
		//
		//// set drand for deal
		//drandNodes := DrandNodes(nodes)
		//for _, node := range nodes {
		//	err := node.SetupDrandWithConfig(&dkg.Config{
		//		Suite:     dkg2.Suite.G1().(dkg.Suite),
		//		NewNodes:  drandNodes,
		//		Threshold: T,
		//		Auth:      TestAuthScheme,
		//		Nonce:     nonce,
		//	})
		//}
	},
}

func DrandNodes(nodes []*dkg2.Node) []dkg.Node {
	var ret []dkg.Node
	for _, n := range nodes {
		ret = append(ret, dkg.Node{
			Index:  n.Index,
			Public: n.Ecies.GetPublicKey(),
		})
	}

	return ret
}

func initDealCmd() {
	setOutputFlag(dealCommand)
	setPasswordFlag(dealCommand)
	setNodesFolderFlag(dealCommand)
}
