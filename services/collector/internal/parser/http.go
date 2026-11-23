package parser

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func fetchFeed(ctx context.Context, client *http.Client, url string, headers map[string]string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("fetch feed: HTTP %d", resp.StatusCode)
	}
	return resp.Body, nil
}
