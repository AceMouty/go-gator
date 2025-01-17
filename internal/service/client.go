package service

import (
	"context"
	"net/http"
)

func FetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	// identify ourselves to the server we are requesting to
	req.Header.Set("User-Agent", "gator")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rssFeed, err := mapResponseToRssFeed(resp)
	if err != nil {
		return nil, err
	}

	return rssFeed, nil
}
