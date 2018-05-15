package gorgo

import (
	"fmt"
	"sort"

	"github.com/OneOfOne/xxhash"
)

func sortMap(currMap JSONDoc) []string {
	var keys []string
	for k := range currMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func calcHash(sql string) string {
	h := xxhash.New64()
	h.Write([]byte(sql))
	r := h.Sum64()

	rp := &r
	s := fmt.Sprintf("%v", rp)
	return s
}
