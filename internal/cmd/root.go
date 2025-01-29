package cmd

import (
	"coffy/internal/coffy"
	"github.com/spf13/cobra"
	"os"
)

var cfgPath string

var callBack func(config *coffy.Config)

var serverCmd = &cobra.Command{
	Use:   "coffy-server",
	Short: "Coffy's backend application server",
	Long: `As part of the coffy application suite, coffy-server handles all the application logic and data persistence.
coffy-server is the heart of the coffy-suite and comes with a HTTP REST interface to interact with it.`,
	Run: run,
}

func init() {
	serverCmd.Flags().StringVarP(&cfgPath, "config", "c", "./coffy.yaml", "path to coffy-server.yaml")
}

func run(cmd *cobra.Command, args []string) {
	f, err := os.Open(cfgPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	config, err := coffy.ParseFile(f)
	if err != nil {
		panic(err)
	}
	callBack(config)
}

func Execute(callback func(cfg *coffy.Config)) {
	if callback != nil {
		callBack = callback
		if err := serverCmd.Execute(); err != nil {
			panic(err)
		}
	} else {
		panic("missing callback function")
	}
}
