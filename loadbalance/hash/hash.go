// package hash makes a wrapper for the third party package of consistent hash.
package hash

import (
	"github.com/golang/glog"
	"stathat.com/c/consistent"
)

type Hash interface {
	// Add inserts a string element in the consistent hash.
	Add(string)

	// Remove removes an element from the hash.
	Remove(string)

	// Get returns an element close to where key hashes to in the circle.
	Get(key string) string
}

type StathatHash struct {
	c *consistent.Consistent
}

// Add inserts a string element in the consistent hash.
func (h *StathatHash) Add(node string) {
	h.c.Add(node)
}

// Remove removes an element from the hash.
func (h *StathatHash) Remove(node string) {
	h.c.Remove(node)
}

// Get returns an element close to where key hashes to in the circle.
func (h *StathatHash) Get(key string) string {
	node, err := h.c.Get(key)
	if err != nil {
		glog.Errorf("get error key: %s, %v", key, err)
		return ""
	}

	return node
}

// New creates an object which implements interface of Hash.
func NewHash(nodes []string) Hash {
	h := consistent.New()
	h.Set(nodes)
	return &StathatHash{
		c: h,
	}
}
