package main

import (
	"context"

	"github.com/spf13/cobra"

	"mempro/internal/log"
	"mempro/internal/memops"
	"mempro/internal/memops/fileadapter"
)

var initFlagQuiet bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "scaffold $MEMPRO_HOME (idempotent)",
	Long: `Initialize the mempro home directory with the canonical layout
(spec §2.1): directory tree, empty spine.jsonl/symbols.jsonl, seed
directives, README, .gitignore, git init, and an
initial commit. Re-running over an existing home is a near no-op.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		paths, err := resolvePaths()
		if err != nil {
			return err
		}
		ops := fileadapter.NewFileAdapter(paths)
		// The READ-ONLY §9.1 gate variant. `init` must refuse to scaffold
		// over a home written by a newer mempro, but must NOT stamp: the
		// home may not exist yet, and store.Init already owns the
		// write-if-missing of version.toml.
		if _, err := checkHomeFormat(context.Background(), ops); err != nil {
			return err
		}
		opts := memops.InitOptions{
			Quiet: initFlagQuiet,
			Logger: func(format string, args ...any) {
				log.Info(format, args...)
			},
		}
		return ops.Init(context.Background(), opts)
	},
}

func init() {
	initCmd.Flags().BoolVar(&initFlagQuiet, "quiet", false, "suppress per-step output")
	rootCmd.AddCommand(initCmd)
}
