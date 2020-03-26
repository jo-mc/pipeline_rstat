package main

import (
	"fmt"
	"log"
	"os"

	"redo_pipeline_rstat/cmdr"
	

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "uppercase",
	Short: "Transform the input to uppercase letters",
	Long: `Simple demo of the usage of linux pipes
Transform the input (pipe of file) to uppercase letters`,
	RunE: func(cmd *cobra.Command, args []string) error {
		print = logNoop
		if cmdr.Flags.Verbose {
			print = logOut
		}
		return cmdr.RunCommand()
	},
}

// var flags struct {
// 	filepath string
// 	verbose  bool
// }

var flagsName = struct {
	file, fileShort       string
	verbose, verboseShort string
}{
	"file", "f",
	"verbose", "v",
}

var print func(s string)

func main() {
	rootCmd.Flags().StringVarP(
		&cmdr.Flags.Filepath,
		flagsName.file,
		flagsName.fileShort,
		"", "path to the file")
	rootCmd.PersistentFlags().BoolVarP(
		&cmdr.Flags.Verbose,
		flagsName.verbose,
		flagsName.verboseShort,
		false, "log verbose output")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func logNoop(s string) {}

func logOut(s string) {
	log.Println(s)
}
