package rains

import (
	"github.com/ugorji/go/codec"
	"time"
)

// Internal constants and types for CBOR encoding of data model objects
type mapKey uint

const (
	mk_content      mapKey = 0  // Content of a message, shard, or zone
	mk_capabilities mapKey = 1  // Capabilities of server sending message
	mk_signatures   mapKey = 2  // Signatures on a message or section
	mk_subject_name mapKey = 3  // Subject name in an assertion
	mk_subject_zone mapKey = 4  // Zone name in an assertion
	mk_query_name   mapKey = 5  // Qualified subject name in a query
	mk_context      mapKey = 6  // Context of an assertion
	mk_objects      mapKey = 7  // Objects of an assertion
	mk_token        mapKey = 8  // Token for referring to a data item
	mk_shard_range  mapKey = 11 // Lexical range of Assertions in Shard
	mk_query_types  mapKey = 14 // acceptable object types for query
	mk_note_type    mapKey = 17 // Notification type
	mk_query_opts   mapKey = 22 // Set of query options requested
	mk_note_data    mapKey = 23 // Additional notification data
)

type sectionKey uint

const (
	sk_assertion    sectionKey = 1  // Assertion (see Section 5.4)
	sk_shard        sectionKey = 2  // Shard (see Section 5.5)
	sk_zone         sectionKey = 3  // Zone (see Section 5.6)
	sk_query        sectionKey = 4  // Query (see Section 5.7)
	sk_notification sectionKey = 23 // Notification (see Section 5.8)
)

type objectType uint

const (
	ot_name         objectType = 1  // name associated with subject
	ot_ip6_addr     objectType = 2  // IPv6 address of subject
	ot_ip4_addr     objectType = 3  // IPv4 address of subject
	ot_redirection  objectType = 4  // name of zone authority server
	ot_delegation   objectType = 5  // public key for zone delgation
	ot_nameset      objectType = 6  // name set expression for zone
	ot_cert_info    objectType = 7  // certificate information for name
	ot_service_info objectType = 8  // service information for srvname
	ot_registrar    objectType = 9  // registrar information
	ot_registrant   objectType = 10 // registrant information
	ot_infrakey     objectType = 11 // public key for RAINS infrastructure
)

type algorithmType uint

const (
	ECDSA_256 algorithmType = 2
	ECDSA_384 algorithmType = 3
)

type notificationType uint

const (
	Heartbeat           notificationType = 100 // Connection heartbeat
	Capability          notificationType = 399 // Capability hash not understood
	MalformedMessage    notificationType = 400 // Malformed message received
	InconsistentMessage notificationType = 403 // Inconsistent message received
	NoAssertion         notificationType = 404 // No assertion exists (client protocol only)
	TooLarge            notificationType = 413 // Message too large
	ServerError         notificationType = 500 // Unspecified server error
	NotCapable          notificationType = 501 // Server not capable
	NotAvailable        notificationType = 504 // No assertion available
)

type Signature struct {
	alg             algorithmType
	validFrom       time.Time
	validUntil      time.Time
	revocationToken []byte
	content         []byte
}

type NameObject string

type IP6AddrObject [16]byte

type IP4AddrObject [4]byte

type RedirectionObject string

type DelegationObject struct {
	Alg     algorithmType
	Content []byte
}

type NamesetObject string

type CertificateObject struct {
}

type ServiceObject struct {
	Hostname      string
	TransportPort uint16
	Priority      uint16
}

type RegistrarObject string

type RegistrantObject string

type InfrakeyObject struct {
	Alg     algorithmType
	Content []byte
}

type Assertion struct {
	Name         string
	Zone         string
	Context      string
	Names        []NameObject
	IP6Addrs     []IP6AddrObject
	IP4Addrs     []IP4AddrObject
	Redirections []RedirectionObject
	Delegations  []DelegationObject
	Namesets     []NamesetObject
	Registrars   []RegistrarObject
	Registrants  []RegistrantObject
	Certificate  []CertificateObject
	Infrakeys    []InfrakeyObject
	Signatures   []Signature
	parent       *AssertionSet
}

type AssertionSet struct {
	Zone         string
	Context      string
	Assertions   []Assertion
	Signatures   []Signature
	ShardRange   [2]string
	ZoneComplete bool
}

type Query struct {
	Name        string
	Context     string
	ObjectTypes []objectType
}

type Notification struct {
	NoteType notificationType
	NoteData string
}

type Message struct {
	Assertions   []Assertion
	Shards       []AssertionSet
	Zones        []AssertionSet
	Capabilities []string
	token        string
}
