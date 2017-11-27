package managers

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/fvbock/trie"
	"google.golang.org/api/iterator"
)

type Ngram struct {
	Keyword  string
	Category string
	Count    int
}

type keywordSearcher struct {
	lookup map[string]string
	trie   *trie.Trie
}

func (k keywordSearcher) Search(query string, count int) (*Result, error) {
	prefix := FormatProductKey(query, "_")

	keywords := make([]Keyword, 0)
	total := 0
	var prev *trie.MemberInfo
	for _, member := range k.trie.PrefixMembers(prefix) {
		if total == count {
			break
		}
		if prev != nil {
			if s, isPresent := k.lookup[member.Value]; isPresent {
				if !strings.HasPrefix(member.Value, prev.Value) || member.Count != prev.Count {
					keywords = append(keywords, Keyword{
						Value: s,
						Count: member.Count,
					})
					total++
				}
			}
		}

		prev = member
	}

	result := &Result{
		Query:    prefix,
		Keywords: keywords,
	}
	return result, nil
}

func KeywordSearcher(bucketName, path string) (Searcher, error) {
	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	q := &storage.Query{
		Prefix: path,
	}

	//output/20171113T1953Z/keywords-00004-of-00007.gz

	var files []string

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)
	current := 0
	it := bucket.Objects(ctx, q)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		values := strings.Split(attrs.Name, "/")
		if len(values) == 3 {
			date, _ := strconv.Atoi(values[1])
			if current < date {
				files = make([]string, 0)
				current = date
			}
			files = append(files, attrs.Name)
		}
	}

	lookup := make(map[string]string)
	trie := trie.NewTrie()

	for _, path := range files {
		obj := bucket.Object(path)
		rdr, err := obj.ReadCompressed(true).NewReader(ctx)
		if err != nil {
			return nil, err
		}
		defer rdr.Close()

		gzrdr, err := gzip.NewReader(rdr)
		if err != nil {
			return nil, err
		}
		defer gzrdr.Close()

		scanner := bufio.NewScanner(gzrdr)
		for scanner.Scan() {
			var ngram Ngram
			if err = json.NewDecoder(bytes.NewReader(scanner.Bytes())).Decode(&ngram); err == nil && ngram.Category == "" {
				key := FormatProductKey(ngram.Keyword, "_")
				lookup[key] = ngram.Keyword
				for i := 0; i < ngram.Count; i++ {
					trie.Add(key)
				}
			}
		}
	}

	return keywordSearcher{lookup, trie}, nil
}
