package cmd

import (
	"github.com/spf13/cobra"
)

const defaultAddr = "localhost:50051"
const addrUsage = "Address to connect to the maestro server"

var addr string

var assetCmd = &cobra.Command{
	Use:   "asset",
	Short: "maestro-cli asset manages your assets",
}

func init() {
	roodCmd.AddCommand(assetCmd)
	assetCmd.PersistentFlags().StringVar(&addr, "addr", defaultAddr, addrUsage)
}
