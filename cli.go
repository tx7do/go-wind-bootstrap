package bootstrap

import (
	"github.com/spf13/cobra"
)

// CommandFlags holds the CLI flags consumed by the bootstrap framework.
type CommandFlags struct {
	// Conf is the path to the bootstrap config file.
	// Supported formats: .yaml, .yml, .json, .bin, .pb.
	Conf string
}

// NewCommandFlags returns a CommandFlags with sensible defaults.
func NewCommandFlags() *CommandFlags {
	return &CommandFlags{
		Conf: "config.yaml",
	}
}

// AddFlags binds the flags to the given cobra command.
func (f *CommandFlags) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&f.Conf, "conf", "c", f.Conf, "config file path (yaml/json/bin)")
}

// NewRootCmd creates the root cobra command with the standard bootstrap flags.
// The runE function is called after flags are parsed.
func NewRootCmd(f *CommandFlags, runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the application",
		RunE:  runE,
	}
	f.AddFlags(cmd)
	return cmd
}
