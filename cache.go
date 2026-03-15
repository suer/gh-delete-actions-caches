package main

import (
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
)

type ActionsCache struct {
	ID  int    `json:"id"`
	Key string `json:"key"`
	Ref string `json:"ref"`
}

type listCachesResponse struct {
	TotalCount    int            `json:"total_count"`
	ActionsCaches []ActionsCache `json:"actions_caches"`
}

func fetchCaches(client *api.RESTClient, repo repository.Repository, keyPrefix string) ([]ActionsCache, error) {
	var all []ActionsCache
	page := 1
	for {
		var resp listCachesResponse
		url := fmt.Sprintf("repos/%s/%s/actions/caches?per_page=100&page=%d", repo.Owner, repo.Name, page)
		if keyPrefix != "" {
			url += "&key=" + keyPrefix
		}
		if err := client.Get(url, &resp); err != nil {
			return nil, err
		}
		all = append(all, resp.ActionsCaches...)
		if len(all) >= resp.TotalCount {
			break
		}
		page++
	}
	return all, nil
}

func filterCaches(caches []ActionsCache, refPrefix string) []ActionsCache {
	if refPrefix == "" {
		return caches
	}
	var result []ActionsCache
	for _, cache := range caches {
		if strings.HasPrefix(cache.Ref, refPrefix) {
			result = append(result, cache)
		}
	}
	return result
}

func deleteCaches(repo repository.Repository, opts *Options) error {
	client, err := api.NewRESTClient(api.ClientOptions{Host: repo.Host})
	if err != nil {
		return err
	}

	caches, err := fetchCaches(client, repo, opts.KeyPrefix)
	if err != nil {
		return err
	}

	targets := filterCaches(caches, opts.RefPrefix)

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
