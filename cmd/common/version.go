package common

import (
	"fmt"
	"github.com/Aptomi/aptomi/pkg/version"
	"github.com/spf13/cobra"
)

// Version is the instance of cobra command that shows version from version package (injected at build tome)
var Version = &cobra.Command{
	Use:   "version",
	Short: "Print the Aptomi Client version",
	Run: func(cmd *cobra.Command, args []string) {
		info := version.GetBuildInfo()
		fmt.Printf("Aptomi Client version: %s\n       git commit: %s\n       built: %s\n", info.GitVersion, info.GitCommit, info.BuildDate)
	},
}
