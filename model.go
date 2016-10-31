package rainsd

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type CBORWriter interface {
	WriteTag(tag int)
	WriteInt(val int)
	WriteBytes(val []byte)
	WriteString(val string)
	WriteBool(val bool)
	WriteNull()
	WriteMapStart(mlen int)
	WriteMapEnd()
	WriteArrayStart(alen int)
	WriteArrayEnd()
	CheckError() error
}

const CBORTagUTC int = 1
const CBORTagRains int = 15309736

// Internal constants and types for CBOR encoding of data model objects
type mapKey uint

const (
	mk_signatures     mapKey = 0 // Signatures on a message or section
	mk_capabilities   mapKey = 1 // Capabilities of server sending message
	mk_token          mapKey = 2 // Token for referring to a data item
	mk_subject_name   mapKey = 3 // Subject name in an assertion
	mk_subject_zone   mapKey = 4 // Zone name in an assertion
	mk_query_name     mapKey = 5 // Qualified subject name in a query
	mk_context        mapKey = 6 // Context of an assertion
	mk_objects        mapKey = 7 // Objects of an assertion
	mk_query_contexts mapKey = 8
	mk_query_types    mapKey = 9  // acceptable object types for query
	mk_query_opts     mapKey = 10 // Set of query options requested
	mk_shard_range    mapKey = 11 // Lexical range of Assertions in Shard
	mk_note_type      mapKey = 21 // Notification type
	mk_note_data      mapKey = 22 // Additional notification data
	mk_content        mapKey = 23 // Content of a message, shard, or zone
)

type sectionKey uint

const (
	sk_assertion    sectionKey = 1  // Assertion (see Section 5.4)
	sk_shard        sectionKey = 2  // Shard (see Section 5.5)
	sk_zone         sectionKey = 3  // Zone (see Section 5.6)
	sk_query        sectionKey = 4  // Query (see Section 5.7)
	sk_notification sectionKey = 23 // Notification (see Section 5.8)
)

type QueryOption uint

const (
	FastOption       QueryOption = 1
	SmallOption      QueryOption = 2
	QuietOption      QueryOption = 3
	SilentOption     QueryOption = 4
	ExpiredOkOption  QueryOption = 5
	TokenTraceOption QueryOption = 6
	NoVerification   QueryOption = 7
)

type ObjectType uint

const (
	NameType        ObjectType = 1  // name associated with subject
	Ip6AddrType     ObjectType = 2  // IPv6 address of subject
	Ip4AddrType     ObjectType = 3  // IPv4 address of subject
	RedirectionType ObjectType = 4  // name of zone authority server
	DelegationType  ObjectType = 5  // public key for zone delgation
	NamesetType     ObjectType = 6  // name set expression for zone
	CertificateType ObjectType = 7  // certificate information for name
	ServiceType     ObjectType = 8  // service information for srvname
	RegistrarType   ObjectType = 9  // registrar information
	RegistrantType  ObjectType = 10 // registrant information
	InfrakeyType    ObjectType = 11 // public key for RAINS infrastructure
)

type AlgorithmType uint

const (
	ECDSA_256 AlgorithmType = 2
	ECDSA_384 AlgorithmType = 3
)

type NotificationType uint

const (
	Heartbeat           NotificationType = 100 // Connection heartbeat
	Capability          NotificationType = 399 // Capability hash not understood
	MalformedMessage    NotificationType = 400 // Malformed message received
	InconsistentMessage NotificationType = 403 // Inconsistent message received
	NoAssertion         NotificationType = 404 // No assertion exists (client protocol only)
	TooLarge            NotificationType = 413 // Message too large
	ServerError         NotificationType = 500 // Unspecified server error
	NotCapable          NotificationType = 501 // Server not capable
	NotAvailable        NotificationType = 504 // No assertion available
)

type Signature struct {
	alg             AlgorithmType
	validFrom       time.Time
	validUntil      time.Time
	revocationToken []byte
	content         []byte
}

type Object interface {
	Emit(w *CBORWriter) error
	Answers(otypes map[ObjectType]bool) bool
}

func (sig *Signature) Emit(w *CBORWriter) error {
	w.WriteArrayStart(5)
	w.WriteInt(sig.alg)
	w.WriteTag(CBORTagUTC)
	w.WriteInt(sig.validFrom.UTC().Unix())
	w.WriteTag(CBORTagUTC)
	w.WriteInt(sig.validUntil.UTC().Unix())
	w.WriteBytes(sig.revocationToken)
	w.WriteBytes(sig.content)
	w.WriteArrayEnd()
	return w.CheckError()
}

type NameObject string

func (name *NameObject) Emit(w *CBORWriter) error {
	w.WriteArrayStart(2)
	w.WriteInt(NameType)
	w.WriteString(string(*name))
	w.WriteArrayEnd()
	return w.CheckError()
}

type IP6AddrObject [16]byte

func (addr6 *IP6AddrObject) Emit(w *CBORWriter) error {
	w.WriteArrayStart(2)
	w.WriteInt(Ip6AddrType)
	w.WriteBytes([]byte(*addr6))
	w.WriteArrayEnd()
	return w.CheckError()
}

type IP4AddrObject [4]byte

func (addr4 *IP4AddrObject) Emit(w *CBORWriter) error {
	w.WriteArrayStart(2)
	w.WriteInt(Ip4AddrType)
	w.WriteBytes([]byte(*addr4))
	w.WriteArrayEnd()
	return w.CheckError()
}

type RedirectionObject string

func (redir *RedirectionObject) Emit(w *CBORWriter) error {
	w.WriteArrayStart(2)
	w.WriteInt(RedirectionType)
	w.WriteString(string(*redir))
	w.WriteArrayEnd()
	return w.CheckError()
}

type DelegationObject struct {
	Alg     AlgorithmType
	Content []byte
}

func (deleg *DelegationObject) Emit(w *CBORWriter) error {
	w.WriteArrayStart(3)
	w.WriteInt(DelegationType)
	w.WriteInt(deleg.Alg)
	w.WriteBytes(deleg.Content)
	w.WriteArrayEnd()
	return w.CheckError()
}

type NamesetObject string

func (nset *NamesetObject) Emit(w *CBORWriter) error {
	w.WriteArrayStart(2)
	w.WriteInt(NamesetType)
	w.WriteString(string(*nset))
	w.WriteArrayEnd()
	return w.CheckError()
}

type CertificateObject struct {
}

func (cert *CertificateObject) Emit(w *CBORWriter) error {
	w.WriteArrayStart(1)
	w.WriteInt(CertificateType)
	w.WriteArrayEnd()
	return w.CheckError()
}

type ServiceObject struct {
	Hostname      string
	TransportPort uint16
	Priority      uint16
}

func (svc *ServiceObject) Emit(w *CBORWriter) error {
	w.WriteArrayStart(4)
	w.WriteInt(ServiceType)
	w.WriteString(svc.Hostname)
	w.WriteInt(svc.TransportPort)
	w.WriteInt(svc.Priority)
	w.WriteArrayEnd()
	return w.CheckError()
}

type RegistrarObject string

func (reg *RegistrarObject) Emit(w *CBORWriter) error {
	w.WriteArrayStart(2)
	w.WriteInt(RegistrarType)
	w.WriteString(string(*reg))
	w.WriteArrayEnd()
	return w.CheckError()
}

type RegistrantObject string

func (reg *RegistrantObject) Emit(w *CBORWriter) error {
	w.WriteArrayStart(2)
	w.WriteInt(RegistrantType)
	w.WriteString(string(*reg))
	w.WriteArrayEnd()
	return w.CheckError()
}

type InfrakeyObject struct {
	Alg     AlgorithmType
	Content []byte
}

func (ik *InfrakeyObject) Emit(w *CBORWriter) error {
	w.WriteArrayStart(3)
	w.WriteInt(InfrakeyType)
	w.WriteInt(ik.Alg)
	w.WriteBytes(ik.Content)
	w.WriteArrayEnd()
	return w.CheckError()
}

type MessageSection interface {
	EmitSection(w *CBORWriter) error
}

type Assertion struct {
	Name       string
	Zone       string
	Context    string
	Objects    []Object
	Signatures []Signature

	parent *AssertionSet
}

func (a *Assertion) Emit(w *CBORWriter, bare bool) error {

	if bare {
		w.WriteMapStart(5)
	} else {
		w.WriteMapStart(3)
	}

	w.WriteInt(mk_signatures)
	w.WriteArrayStart(len(a.Signatures))
	for _, sig := range a.Signatures {
		if sig.Emit(w); err != nil {
			return err
		}
	}
	w.WriteArrayEnd()

	w.WriteInt(mk_subject_name)
	w.WriteString(a.Name)

	if bare {
		w.WriteInt(mk_subject_zone)
		w.WriteString(a.Zone)

		w.WriteInt(mk_context)
		w.WriteString(a.Context)
	}

	w.WriteInt(mk_objects)
	w.WriteArrayStart(len(a.Objects))
	for _, o := range a.Objects {
		if err := o.Emit(w); err != nil {
			return err
		}
	}
	w.WriteArrayEnd()

	w.WriteMapEnd()

	return w.CheckError()
}

func (a *Assertion) EmitSection(w *CBORWriter) error {
	w.WriteArrayStart(2)
	w.WriteInt(sk_assertion)
	if err := a.Emit(w, true); err != nil {
		return err
	}
	w.WriteArrayEnd()
}

type AssertionSet struct {
	Zone         string
	Context      string
	Assertions   []Assertion
	Signatures   []Signature
	ShardRange   [2]string
	ZoneComplete bool
}

func (as *AssertionSet) EmitSection(w *CBORWriter) error {
	var err error

	w.WriteArrayStart(2)

	if as.ZoneComplete {
		w.WriteInt(sk_zone)
	} else {
		w.WriteInt(sk_shard)
	}

	if !as.ZoneComplete && (len(as.ShardRange[0]) || len(as.ShardRange[1])) {
		w.WriteMapStart(5)
	} else {
		w.WriteMapStart(4)
	}

	w.WriteInt(mk_signatures)
	w.WriteArrayStart(len(as.Signatures))
	for _, sig := range a.Signatures {
		if err := sig.Emit(w); err != nil {
			return err
		}
	}
	w.WriteArrayEnd()

	w.WriteInt(mk_subject_zone)
	w.WriteString(as.Zone)

	w.WriteInt(mk_context)
	w.WriteString(as.Context)

	if !as.ZoneComplete && (len(as.ShardRange[0]) || len(as.ShardRange[1])) {
		w.WriteInt(mk_shard_range)
		w.WriteArrayStart(2)
		w.WriteString(as.ShardRange[0])
		w.WriteString(as.ShardRange[1])
		w.WriteArrayEnd()
	}

	w.WriteInt(mk_content)
	w.WriteArrayStart(len(as.Assertions))
	for _, a := range as.Assertions {
		if err := a.Emit(w, false); err != nil {
			return err
		}
	}
	w.WriteArrayEnd()

	w.WriteMapEnd()

	w.WriteArrayEnd()
	return w.CheckError()
}

type Query struct {
	Name        string
	Contexts    []string
	Token       [16]byte
	ObjectTypes map[ObjectType]bool
	Options     map[QueryOption]bool
}

// FIXME make this actually look at tokens
// FIXME new token type?
// FIXME update rains-protocol to specify 16 bytes for token
func tokenZero([16]byte) bool {
	return false
}

func (q *Query) EmitSection(w *CBORWriter) error {
	var err error
	var mapLength int = 2

	if !tokenZero(q.Token) {
		mapLength++
	}

	if len(q.ObjectTypes) > 0 {
		mapLength++
	}

	if len(q.Options) > 0 {
		mapLength++
	}

	w.WriteArrayStart(2)

	w.WriteInt(sk_query)
	w.WriteMapStart(mapLength)

	if !tokenZero(q.Token) {
		w.WriteInt(mk_token)
		w.WriteBytes(q.Token)
	}

	w.WriteInt(mk_query_name)
	w.WriteString(q.Name)

	w.WriteInt(mk_query_contexts)
	w.WriteArrayStart(len(q.Contexts))
	for _, ctx := range q.Contexts {
		w.WriteString(ctx)
	}
	w.WriteArrayEnd()

	if len(q.ObjectTypes) > 0 {
		w.WriteInt(mk_query_types)
		w.WriteArrayStart(len(q.ObjectTypes))
		for qt, _ := range q.ObjectTypes {
			w.WriteInt(qt)
		}
		w.WriteArrayEnd()
	}

	if len(q.Options) > 0 {
		w.WriteInt(mk_query_types)
		w.WriteArrayStart(len(q.Options))
		for qo, _ := range q.Options {
			w.WriteInt(qo)
		}
		w.WriteArrayEnd()
	}

	w.WriteMapEnd()

	w.WriteArrayEnd()
	return w.CheckError()
}

type Notification struct {
	NoteType NotificationType
	NoteData string
	Token    [16]byte
}

func (q *Notification) EmitSection(w *CBORWriter) error {
	// FIXME write this
	return w.CheckError()
}

type Message struct {
	Assertions   []Assertion
	Shards       []AssertionSet
	Zones        []AssertionSet
	Capabilities []string
	Token        [16]byte
}

func (m *Message) Emit(w *CBORWriter) error {
	// FIXME write this
	return w.CheckError()
}
