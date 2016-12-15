// Query engine for RAINSd. Defines runtime structures and interfaces for using them.

package rainsd


import (
    "net"
    "time"
)

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