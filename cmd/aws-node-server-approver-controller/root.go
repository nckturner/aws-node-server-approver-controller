package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nckturner/aws-node-server-approver-controller/pkg/config"
	"github.com/nckturner/aws-node-server-approver-controller/pkg/controller"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "k8s.io/kubernetes/pkg/client/metrics/prometheus" // for client metric registration
	_ "k8s.io/kubernetes/pkg/util/reflector/prometheus" // for reflector metric registration
	_ "k8s.io/kubernetes/pkg/util/workqueue/prometheus" // for workqueue metric registration
	_ "k8s.io/kubernetes/pkg/version/prometheus"        // for version metric registration
)

var (
	cfgFile string
	master  string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "Load configuration from `filename`")
	rootCmd.Flags().StringVar(&master, "master", "", "Master address")
	rootCmd.Flags().String("kubeconfig", "", "Path to your kubeconfig")
	viper.BindPFlag("master", rootCmd.Flags().Lookup("master"))
	viper.BindPFlag("kubeconfig", rootCmd.Flags().Lookup("kubeconfig"))
}

func main() {
	// hack to make flag.Parsed return true such that glog is happy
	// about the flags having been parsed
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	/* #nosec */
	_ = fs.Parse([]string{})
	flag.CommandLine = fs

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile == "" {
		return
	}
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Can't read configuration file %q: %v\n", cfgFile, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use: "aws-node-server-approver-controller",
	Long: `The AWS node server approver is a controller responsible
for approving certificate signing requests from AWS EC2 instances for
kubelet server certificates.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := getConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		c, err := controller.NewController(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		c.Run()
	},
}

func getConfig() (*config.AWSNodeServerApproverConfig, error) {
	cfg := config.AWSNodeServerApproverConfig{
		Options: config.AWSNodeServerApproverOptions{
			Master:     viper.GetString("master"),
			Kubeconfig: viper.GetString("kubeconfig"),
		},
	}

	// TODO validation
	return &cfg, nil
}
