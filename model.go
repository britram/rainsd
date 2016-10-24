package rains

import (
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

type NameObject string

func (this *NameObject) sliceContents() []interface{} {
	out := make([]interface{}, 2)
	out[0] = ot_name
	out[1] = *this
	return out
}

type IP6AddrObject [16]byte

func (this *IP6AddrObject) sliceContents() []interface{} {
	out := make([]interface{}, 2)
	out[0] = ot_ip6_addr
	out[1] = *this
	return out
}

type IP4AddrObject [4]byte

func (this *IP4AddrObject) sliceContents() []interface{} {
	out := make([]interface{}, 2)
	out[0] = ot_ip4_addr
	out[1] = *this
	return out
}

type RedirectionObject string

func (this *RedirectionObject) sliceContents() []interface{} {
	out := make([]interface{}, 2)
	out[0] = ot_redirection
	out[1] = *this
	return out
}

type DelegationObject struct {
	Alg     algorithmType
	Content []byte
}

func (this *DelegationObject) sliceContents() []interface{} {
	out := make([]interface{}, 3)
	out[0] = ot_delegation
	out[1] = this.Alg
	out[2] = this.Content
	return out
}

type NamesetObject string

func (this *NamesetObject) sliceContents() []interface{} {
	out := make([]interface{}, 2)
	out[0] = ot_nameset
	out[1] = *this
	return out
}

type CertificateObject struct {
}

type ServiceObject struct {
	Hostname      string
	TransportPort uint16
	Priority      uint16
}

func (this *ServiceObject) sliceContents() []interface{} {
	out := make([]interface{}, 4)
	out[0] = ot_service_info
	out[1] = this.Hostname
	out[2] = this.TransportPort
	out[3] = this.Priority
	return out
}

type RegistrarObject string

func (this *RegistrarObject) sliceContents() []interface{} {
	out := make([]interface{}, 4)
	out[0] = ot_registrar
	out[1] = *this
	return out
}

type RegistrantObject string

func (this *RegistrantObject) sliceContents() []interface{} {
	out := make([]interface{}, 4)
	out[0] = ot_registrar
	out[1] = *this
	return out
}

type InfrakeyObject struct {
	Alg     algorithmType
	Content []byte
}

func (this *InfrakeyObject) sliceContents() []interface{} {
	out := make([]interface{}, 3)
	out[0] = ot_infrakey
	out[1] = this.Alg
	out[2] = this.Content
	return out
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
}

type AssertionSet struct {
	Zone         string
	Context      string
	Assertions   []Assertion
	Signatures   []Signature
	ShardRange   []string
	ZoneComplete bool
}

type Query struct {
	Name        string
	Contexts    []string
	ObjectTypes []objectType
}

type Notification struct {
}
