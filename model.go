package rains

import (
	"time"
)

// Internal constants and types for CBOR encoding of data model objects
type mapKey uint

const (
	content      mapKey = 0  // Content of a message, shard, or zone
	capabilities mapKey = 1  // Capabilities of server sending message
	signatures   mapKey = 2  // Signatures on a message or section
	subject_name mapKey = 3  // Subject name in an assertion
	subject_zone mapKey = 4  // Zone name in an assertion
	query_name   mapKey = 5  // Qualified subject name in a query
	context      mapKey = 6  // Context of an assertion
	objects      mapKey = 7  // Objects of an assertion
	token        mapKey = 8  // Token for referring to a data item
	shard_range  mapKey = 11 // Lexical range of Assertions in Shard
	query_types  mapKey = 14 // acceptable object types for query
	note_type    mapKey = 17 // Notification type
	query_opts   mapKey = 22 // Set of query options requested
	note_data    mapKey = 23 // Additional notification data
)

type sectionKey uint

const (
	assertion    sectionKey = 1  // Assertion (see Section 5.4)
	shard        sectionKey = 2  // Shard (see Section 5.5)
	zone         sectionKey = 3  // Zone (see Section 5.6)
	query        sectionKey = 4  // Query (see Section 5.7)
	notification sectionKey = 23 // Notification (see Section 5.8)
)

type objectType uint

const (
	name         objectType = 1  // name associated with subject
	ip6_addr     objectType = 2  // IPv6 address of subject
	ip4_addr     objectType = 3  // IPv4 address of subject
	redirection  objectType = 4  // name of zone authority server
	delegation   objectType = 5  // public key for zone delgation
	nameset      objectType = 6  // name set expression for zone
	cert_info    objectType = 7  // certificate information for name
	service_info objectType = 8  // service information for srvname
	registrar    objectType = 9  // registrar information
	registrant   objectType = 10 // registrant information
	infrakey     objectType = 11 // public key for RAINS infrastructure
)

type algorithmType uint

const (
	ecdsa_256 algorithmType = 2
	ecdsa_384 algorithmType = 3
)

type Signature struct {
	alg             algorithmType
	validFrom       time.Time
	validUntil      time.Time
	revocationToken []byte
	content         []byte
}

// FIXME: is this an interface? go read the go book...
type Object struct {
}

type Assertion struct {
	subjectName string
	subjectZone string
	context     string
	objects     []Object
	signatures  []Signature
}

type AssertionSet struct {
	subjectZone  string
	context      string
	assertions   []Assertion
	signatures   []Signature
	shardTange   []string
	zoneComplete bool
}

type Query struct {
	name        string
	contexts    []string
	objectTypes []objectType
}
