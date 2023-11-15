package confluence

import (
	"encoding/json"
	"io"
	"net/url"
	"strings"
)

type ContentResult struct {
	Content               Content `json:"content"`
	Title                 string  `json:"title"`
	Excerpt               string  `json:"excerpt"`
	URL                   string  `json:"url"`
	ResultGlobalContainer struct {
		Title      string `json:"title"`
		DisplayURL string `json:"displayUrl"`
	}
	// breadcrumbs is ignored
	EntityType           string `json:"entityType"`
	IconCSSClass         string `json:"iconCssClass"`
	LastModified         string `json:"lastModified"`
	FriendlyLastModified string `json:"friendlyLastModified"`
}

type SearchResults struct {
	ResultPagination
	Results        []ContentResult `json:"results"`
	TotalSize      int             `json:"totalSize"`
	CqlQuery       string          `json:"cqlQuery"`
	SearchDuration int             `json:"SearchDuration"`
	// links are ignored
}

func (w *Wiki) searchEndpoint() (*url.URL, error) {
	return url.ParseRequestURI(w.endPoint.String() + "/search")
}

func (w *Wiki) Search(cql, cqlContext string, expand []string, limit int) (*SearchResults, error) {
	searchEndPoint, err := w.searchEndpoint()
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("expand", strings.Join(expand, ","))
	data.Set("cqlcontext", cqlContext)
	data.Set("cql", cql)
	searchEndPoint.RawQuery = data.Encode()

	res, err := w.client.Get(searchEndPoint.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var results SearchResults

	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}
