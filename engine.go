// Query engine for RAINSd. Defines runtime structures and interfaces for using them.

package rainsd

import (
	"net"
	"time"
)

type AssertionCallback func(assertion string) error

type QueryEngine struct {
	// map of zone caches, keyed by "context zone"
	zones map[string]ZoneCache
	// list of pending query requests, managed as a heap
	pending []QueryRequest
}

func (*QueryEngine) Assert(assertion string) error {

}

func (*QueryEngine) Query(query string, callback AssertionCallback) error {

}

func (*QueryEngine) Reap() {
	// expire pending queries
}

type ZoneCache struct {
	// keyed by "subject objtype"
	assertions map[string][]ObjectValue
	keys       []ZoneKey
}

type ObjectValue interface {
	// Convert object to a string for presentation (and use in a short assertion)
	String() string
	// Determine whether an object value
	ValidAt(time.Time) bool
}

type ZoneKey struct {
	expires time.Time
}

/*
type AssertionCache struct {
    Contexts map[string]ContextCache
    ReapQueue []ExpiryEvent
}

func (ac *AssertionCache) Parse(astr string) error {
    return nil
}

type ContextCache struct {
    Name string
    Zones map[string]ZoneCache

type ZoneCache struct {
    Name string
    Keys []ZoneKey
    Subjects map[string]SubjectCache
}

type ZoneKey struct {

}

type SubjectCache struct {
    Name string
    Assertions []SingularAssertion
}

type SingularAssertion struct {
    Identifier string
    ValidFrom time.Time
    ValidUntil time.Time
    Provenance string
    Object Object
}

type ExpiryEvent struct {
    ValidUntil time.Time
    IdentifierPath []string
}
*/
