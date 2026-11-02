package main

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
)

type Options struct {
	Version   bool
	KeyPrefix string
	RefPrefix string
	DryRun    bool
	Repo      string
}

func rootCmd() *cobra.Command {
	opts := &Options{}
	cmd := &cobra.Command{
		Use:   "gh delete-actions-caches",
		Short: "Delete GitHub Actions caches by key prefix or ref prefix",
		Example: `  # Delete caches by key prefix
  gh delete-actions-caches -k Linux-build

  # Delete caches by ref prefix
  gh delete-actions-caches -r refs/heads/feature/

  # Match both key prefix and ref prefix (AND)
  gh delete-actions-caches -k Linux-build -r refs/heads/main

  # Preview caches to be deleted (dry run)
  gh delete-actions-caches -k Linux-build -d

  # Target a specific repository
  gh delete-actions-caches -k Linux-build -R owner/repo`,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Version, "version", false, "show version")
	cmd.Flags().StringVarP(&opts.KeyPrefix, "key-prefix", "k", "", "Delete caches whose key starts with this `prefix`")
	cmd.Flags().StringVarP(&opts.RefPrefix, "ref-prefix", "r", "", "Delete caches whose ref starts with this `prefix`")
	cmd.Flags().BoolVarP(&opts.DryRun, "dryrun", "d", false, "List caches to be deleted without actually deleting")
	cmd.Flags().StringVarP(&opts.Repo, "repo", "R", "", "Target `repository` (OWNER/REPO)")

	return cmd
}

func run(opts *Options) error {
	if opts.Version {
		if info, ok := debug.ReadBuildInfo(); ok {
			fmt.Println(info.Main.Version)
			return nil
		}
		return errors.New("could not read build info")
	}

	if opts.KeyPrefix == "" && opts.RefPrefix == "" {
		return errors.New("--key-prefix or --ref-prefix is required")
	}

	var (
		repo repository.Repository
		err  error
	)
	if opts.Repo != "" {
		repo, err = repository.Parse(opts.Repo)
	} else {
		repo, err = repository.Current()
	}
	if err != nil {
		return err
	}

	return deleteCaches(repo, opts)
}
