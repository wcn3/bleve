//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package searchers

import (
	"testing"

	"github.com/blevesearch/bleve/search"
)

func TestMatchNoneSearch(t *testing.T) {

	twoDocIndexReader := twoDocIndex.Reader()
	defer twoDocIndexReader.Close()

	noneSearcher, err := NewMatchNoneSearcher(twoDocIndexReader)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		searcher search.Searcher
		results  []*search.DocumentMatch
	}{
		{
			searcher: noneSearcher,
			results:  []*search.DocumentMatch{},
		},
	}

	for testIndex, test := range tests {
		defer test.searcher.Close()

		next, err := test.searcher.Next()
		i := 0
		for err == nil && next != nil {
			if i < len(test.results) {
				if next.ID != test.results[i].ID {
					t.Errorf("expected result %d to have id %s got %s for test %d", i, test.results[i].ID, next.ID, testIndex)
				}
				if !scoresCloseEnough(next.Score, test.results[i].Score) {
					t.Errorf("expected result %d to have score %v got  %v for test %d", i, test.results[i].Score, next.Score, testIndex)
					t.Logf("scoring explanation: %s", next.Expl)
				}
			}
			next, err = test.searcher.Next()
			i++
		}
		if err != nil {
			t.Fatalf("error iterating searcher: %v for test %d", err, testIndex)
		}
		if len(test.results) != i {
			t.Errorf("expected %d results got %d for test %d", len(test.results), i, testIndex)
		}
	}
}
