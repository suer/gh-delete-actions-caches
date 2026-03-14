package main

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
)

func deleteCaches(repo repository.Repository, opts *Options) error {
	client, err := api.NewRESTClient(api.ClientOptions{Host: repo.Host})
	if err != nil {
		return err
	}

	caches, err := fetchCaches(client, repo)
	if err != nil {
		return err
	}

	targets := filterCaches(caches, opts)

	if opts.DryRun {
		fmt.Printf("%d cache(s) to be deleted:\n", len(targets))
		for _, cache := range targets {
			fmt.Printf("  key: %s, ref: %s\n", cache.Key, cache.Ref)
		}
		return nil
	}

	for _, cache := range targets {
		url := fmt.Sprintf("repos/%s/%s/actions/caches/%d", repo.Owner, repo.Name, cache.ID)
		if err := client.Delete(url, nil); err != nil {
			return fmt.Errorf("failed to delete cache %d (%s): %w", cache.ID, cache.Key, err)
		}
		fmt.Printf("Deleted: %s (ref: %s)\n", cache.Key, cache.Ref)
	}
	fmt.Printf("%d cache(s) deleted.\n", len(targets))

	return nil
}
