package hypergraph

import (
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/OneOfOne/xxhash"
)

func getHash(arr []int32) uint32 {
	h := xxhash.New32()
    
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})

	in := ""

	for _, j := range arr {
		in += (strconv.Itoa(int(j)) + "|")
	}
	r := strings.NewReader(in)
	io.Copy(h, r)

	return h.Sum32();
}	