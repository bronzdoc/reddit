// Package reddit implements a basic client for the reddit API
package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Describes a reddit item
type Item struct {
	Title    string
	URL      string
	Comments int `json:"num_comments"`
}

func (i Item) String() string {
	com := ""
	switch i.Comments {
	case 0:
		// Nothing
	case 1:
		com = " (1 comment)"
	default:
		com = fmt.Sprintf(" (%d comments)", i.Comments)
	}

	return fmt.Sprintf("%s%s\n%s", i.Title, com, i.URL)
}

type response struct {
	Data struct {
		Children []struct {
			Data Item
		}
	}
}

// Fetches the most recent items posted in a specified subreddit
func Get(reddit string) ([]Item, error) {
	url := fmt.Sprintf("https://reddit.com/r/%s.json", reddit)
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	r := new(response)
	err = json.NewDecoder(res.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	items := make([]Item, len(r.Data.Children))
	for i, child := range r.Data.Children {
		items[i] = child.Data
	}

	return items, nil
}
