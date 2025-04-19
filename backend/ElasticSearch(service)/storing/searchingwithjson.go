package storing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type SearchRequest struct {
	Query BoolQuery `json:"query"`
}

// BoolQuery represents the boolean query structure
type BoolQuery struct {
	Bool BoolCondition `json:"bool"`
}

// BoolCondition holds the must and should clauses
type BoolCondition struct {
	Must   []MultiMatchQuery  `json:"must"`
	Should []MatchPhraseQuery `json:"should"`
}

// MultiMatchQuery for searching across multiple fields
type MultiMatchQuery struct {
	MultiMatch struct {
		Query     string   `json:"query"`
		Fields    []string `json:"fields"`
		Type      string   `json:"type"`
		Fuzziness string   `json:"fuzziness,omitempty"`
	} `json:"multi_match"`
}



// MatchPhraseQuery for exact phrase matching
type MatchPhraseQuery struct {
	MatchPhrase map[string]struct {
		Query string `json:"query"`
	} `json:"match_phrase"`
}

// Highlighting structure
type Highlight struct {
	Fields map[string]HighlightField `json:"fields"`
}

type HighlightField struct {
	PreTags  []string `json:"pre_tags"`
	PostTags []string `json:"post_tags"`
}

func SearchingJson(query string) ([]map[string]interface{}, error) {
	searchBody := map[string]interface{}{
		"query": BoolQuery{
			Bool: BoolCondition{
				Must: []MultiMatchQuery{
					{
						MultiMatch: struct {
							Query     string   `json:"query"`
							Fields    []string `json:"fields"`
							Type      string   `json:"type"`
							Fuzziness string   `json:"fuzziness,omitempty"`
						}{
							Query:     query,
							Fields:    []string{"name", "name.edge_ngram"},
							Type:      "most_fields",
							Fuzziness: "auto",
						},
					},
				},
				Should: []MatchPhraseQuery{
					{
						MatchPhrase: map[string]struct {
							Query string `json:"query"`
						}{
							"name": {Query: query},
						},
					},
				},
			},
		},
		"highlight": Highlight{
			Fields: map[string]HighlightField{
				// "name": {
				// 	PreTags:  []string{"<b>"},
				// 	PostTags: []string{"</b>"},
				// },
				"name.edge_ngram": {
					PreTags:  []string{"<b>"},
					PostTags: []string{"</b>"},
				},
			},
		},
	}

	// Convert to JSON
	jsonBody, err := json.Marshal(searchBody)
	if err != nil {
		log.Printf("Error marshaling search query: %v\n", err)
		return nil, err
	}

	// Send request to Elasticsearch
	req := bytes.NewReader(jsonBody)
	res, err := ES.Search(
		ES.Search.WithContext(context.Background()),
		ES.Search.WithIndex("products"),
		ES.Search.WithBody(req),
		ES.Search.WithPretty(),
	)

	if err != nil {
		log.Printf("Error executing search query: %v\n", err)
		return nil, err
	}
	defer res.Body.Close()

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return nil, err
	}

	// Convert response to a map
	var rawResult map[string]any
	err = json.Unmarshal(body, &rawResult)
	if err != nil {
		log.Printf("Error unmarshaling response: %v\n", err)
		return nil, err
	}

	hits, ok := rawResult["hits"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("INVALID STRUCTURE %t", ok)
	}

	hitArray, ok := hits["hits"].([]any)
	if !ok {
		return nil, fmt.Errorf("INVALID HITS ARRAY STRUCTURE %t", ok)
	}
	var formattedResults []map[string]any

	for _, h := range hitArray {
		hitMap, _ := h.(map[string]any)
		formattedResult := map[string]any{
			"_source":   hitMap["_source"],
			"highlight": hitMap["highlight"],
		}
		formattedResults = append(formattedResults, formattedResult)
	}

	return formattedResults, nil

}
