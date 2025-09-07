package mnemosyne

import (
	"hash/fnv"
	"math/rand"
)

const (
	bucketSize      = 4    // Number of entries per bucket
	numBuckets      = 1024 // Total number of buckets
	fingerprintSize = 2    // Size of fingerprints in bytes
	maxKicks        = 500  // Maximum number of eviction attempts
)

type CuckooFilter struct {
	buckets [][]uint16
}

func NewCuckooFilter() *CuckooFilter {
	buckets := make([][]uint16, numBuckets)
	for i := range buckets {
		buckets[i] = make([]uint16, bucketSize)
	}
	return &CuckooFilter{buckets: buckets}
}

// hashFunc generates a hash for the given data
func hashFunc(data []byte) uint32 {
	h := fnv.New32()
	h.Write(data)
	return h.Sum32()
}

// fingerprint generates a short fingerprint for an item
func fingerprint(data []byte) uint16 {
	h := hashFunc(data)
	return uint16(h) + 1 // Avoid zero fingerprints
}

// index computes the primary bucket index
func index(data []byte) int {
	return int(hashFunc(data) % numBuckets)
}

// altIndex computes the alternate bucket index using the fingerprint
func altIndex(idx int, fp uint16) int {
	h := hashFunc([]byte{byte(fp), byte(fp >> 8)})
	return (idx ^ int(h%numBuckets)) % numBuckets
}

// Insert tries to insert an item into the filter
func (cf *CuckooFilter) Insert(data []byte) bool {
	fp := fingerprint(data)
	i1 := index(data)
	i2 := altIndex(i1, fp)

	// Try inserting into one of the buckets
	if cf.insertIntoBucket(i1, fp) || cf.insertIntoBucket(i2, fp) {
		return true
	}

	// Eviction process
	i := []int{i1, i2}[rand.Intn(2)]
	for n := 0; n < maxKicks; n++ {
		j := rand.Intn(bucketSize)
		evicted := cf.buckets[i][j]
		cf.buckets[i][j] = fp
		fp = evicted
		i = altIndex(i, fp)

		if cf.insertIntoBucket(i, fp) {
			return true
		}
	}
	return false
}

// insertIntoBucket tries to insert a fingerprint into a bucket
func (cf *CuckooFilter) insertIntoBucket(i int, fp uint16) bool {
	for j := range cf.buckets[i] {
		if cf.buckets[i][j] == 0 { // Empty slot
			cf.buckets[i][j] = fp
			return true
		}
	}
	return false
}

// Lookup checks if an item is in the filter
func (cf *CuckooFilter) Lookup(data []byte) bool {
	fp := fingerprint(data)
	i1 := index(data)
	i2 := altIndex(i1, fp)

	// Check both possible locations
	return cf.contains(i1, fp) || cf.contains(i2, fp)
}

// contains checks if a fingerprint is in a bucket
func (cf *CuckooFilter) contains(i int, fp uint16) bool {
	for _, f := range cf.buckets[i] {
		if f == fp {
			return true
		}
	}
	return false
}

// Delete removes an item from the filter
func (cf *CuckooFilter) Delete(data []byte) bool {
	fp := fingerprint(data)
	i1 := index(data)
	i2 := altIndex(i1, fp)

	// Try deleting from either bucket
	return cf.removeFromBucket(i1, fp) || cf.removeFromBucket(i2, fp)
}

// removeFromBucket removes a fingerprint from a bucket
func (cf *CuckooFilter) removeFromBucket(i int, fp uint16) bool {
	for j := range cf.buckets[i] {
		if cf.buckets[i][j] == fp {
			cf.buckets[i][j] = 0
			return true
		}
	}
	return false
}

// func main() {
// 	cf := NewCuckooFilter()

// 	// Test inserting and looking up values
// 	items := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}

// 	for _, item := range items {
// 		cf.Insert(item)
// 	}

// 	for _, item := range items {
// 		fmt.Printf("Lookup %s: %v\n", item, cf.Lookup(item))
// 	}

// 	// Test deletion
// 	fmt.Println("Deleting 'bar'")
// 	cf.Delete([]byte("bar"))
// 	fmt.Printf("Lookup 'bar': %v\n", cf.Lookup([]byte("bar")))
// }
