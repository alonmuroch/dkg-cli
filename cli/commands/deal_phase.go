package commands

import "github.com/spf13/cobra"

var DealCommand = func() *cobra.Command {
	initDealCmd()
	return dealCommand
}()

var dealCommand = &cobra.Command{
	Use:   "deal",
	Short: "Deal shares and save to json file",
}

func initDealCmd() {

}
