package rainsd

import (
	"encoding/binary"
	"io"
	"fmt"
	"time"
)

type ByteStreamWriter interface {
	io.Writer
	io.ByteWriter
}

type ByteStreamReader interface {
	io.Reader
	io.ByteReader
}

type CBORWriter interface {
	WriteTag(tag int) error
	WriteInt(val int) error
	WriteBytes(val []byte) error
	WriteString(val string) error
	WriteBool(val bool) error
	WriteNull() error
	WriteMapStart(mlen int) error
	WriteMapEnd() error
	WriteArrayStart(alen int) error
	WriteArrayEnd() error
}

const CBORTagUTC int = 1

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
}

func (sig *Signature) Emit(w *CBORWriter) error {
	var err error
	if err = w.WriteArrayStart(5); err != nil {
		return err
	}
	if err = w.WriteInt(sig.alg); err != nil {
		return err
	}
	if err = w.WriteTag(CBORTagUTC); err != nil {
		return err
	}
	if err = w.WriteInt(sig.validFrom.UTC().Unix()); err != nil {
		return err
	}
	if err = w.WriteTag(CBORTagUTC); err != nil {
		return err
	}
	if err = w.WriteInt(sig.validUntil.UTC().Unix()); err != nil {
		return err
	}
	if err = w.WriteBytes(sig.revocationToken); err != nil {
		return err
	}
	if err = w.WriteBytes(sig.content); err != nil {
		return err
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}
	return nil
}

type NameObject string

func (name *NameObject) Emit(w *CBORWriter) error {
	var err error
	if err = w.WriteArrayStart(2); err != nil {
		return err
	}
	if err = w.Int(NameType); err != nil {
		return err
	}
	if err = w.String(string(*name)); err != nil {
		return err
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}
	return nil
}

type IP6AddrObject [16]byte

func (addr6 *IP6AddrObject) Emit(w *CBORWriter) error {
	var err error
	if err = w.WriteArrayStart(2); err != nil {
		return err
	}
	if err = w.WriteInt(Ip6AddrType); err != nil {
		return err
	}
	if err = w.WriteBytes([]byte(*addr6)); err != nil {
		return err
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}
	return nil
}

type IP4AddrObject [4]byte

func (addr4 *IP4AddrObject) Emit(w *CBORWriter) error {
	var err error
	if err = w.WriteArrayStart(2); err != nil {
		return err
	}
	if err = w.WriteInt(Ip4AddrType); err != nil {
		return err
	}
	if err = w.WriteBytes([]byte(*addr4)); err != nil {
		return err
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}
	return nil
}

type RedirectionObject string

func (redir *RedirectionObject) Emit(w *CBORWriter) error {
	var err error
	if err = w.WriteArrayStart(2); err != nil {
		return err
	}
	if err = w.WriteInt(RedirectionType); err != nil {
		return err
	}
	if err = w.WriteString(string(*redir)); err != nil {
		return err
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}
	return nil
}

type DelegationObject struct {
	Alg     AlgorithmType
	Content []byte
}

func (deleg *DelegationObject) Emit(w *CBORWriter) error {
	var err error
	if err = w.WriteArrayStart(3); err != nil {
		return err
	}
	if err = w.WriteInt(DelegationType); err != nil {
		return err
	}
	if err = w.WriteInt(deleg.Alg); err != nil {
		return err
	}
	if err = w.WriteBytes(deleg.Content); err != nil {
		return err
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}
	return nil
}

type NamesetObject string

func (nset *NamesetObject) Emit(w *CBORWriter) error {
	var err error
	if err = w.WriteArrayStart(2); err != nil {
		return err
	}
	if err = w.WriteInt(NamesetType); err != nil {
		return err
	}
	if err = w.WriteString(string(*nset)); err != nil {
		return err
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}
	return nil
}

type CertificateObject struct {
}

func (cert *CertificateObject) Emit(w *CBORWriter) error {
	var err error
	if err = w.WriteArrayStart(1); err != nil {
		return err
	}
	if err = w.WriteInt(CertificateType); err != nil {
		return err
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}
	return nil
}

type ServiceObject struct {
	Hostname      string
	TransportPort uint16
	Priority      uint16
}

func (svc *ServiceObject) Emit(w *CBORWriter) error {
	var err error
	if err = w.WriteArrayStart(4); err != nil {
		return err
	}
	if err = w.WriteInt(ServiceType); err != nil {
		return err
	}
	if err = w.WriteString(svc.Hostname); err != nil {
		return err
	}
	if err = w.WriteInt(svc.TransportPort); err != nil {
		return err
	}
	if err = w.WriteInt(svc.Priority); err != nil {
		return err
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}
	return nil
}

type RegistrarObject string

func (reg *RegistrarObject) Emit(w *CBORWriter) error {
	var err error
	if err = w.WriteArrayStart(2); err != nil {
		return err
	}
	if err = w.WriteInt(RegistrarType); err != nil {
		return err
	}
	if err = w.WriteString(string(*reg)); err != nil {
		return err
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}
	return nil
}

type RegistrantObject string

func (reg *RegistrantObject) Emit(w *CBORWriter) error {
	var err error
	if err = w.WriteArrayStart(2); err != nil {
		return err
	}
	if err = w.WriteInt(RegistrantType); err != nil {
		return err
	}
	if err = w.WriteString(string(*reg)); err != nil {
		return err
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}
	return nil
}

type InfrakeyObject struct {
	Alg     AlgorithmType
	Content []byte
}

func (ik *InfrakeyObject) Emit(w *CBORWriter) error {
	var err error
	if err = w.WriteArrayStart(3); err != nil {
		return err
	}
	if err = w.WriteInt(InfrakeyType); err != nil {
		return err
	}
	if err = w.WriteInt(ik.Alg); err != nil {
		return err
	}
	if err = w.WriteBytes(ik.Content); err != nil {
		return err
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}
	return nil
}

type MessageSection interface {
	EmitSection(w *CBORWriter) error
}

type Assertion struct {
	Name         string
	Zone         string
	Context      string
	Objects      []Object
	Signatures   []Signature
    
    parent       *AssertionSet
}

func (a *Assertion) Emit(w *CBORWriter, bare bool) error {
	var err error

	if err = w.WriteMapStart(bare ? 5 : 3); err != nil {
		return err
	}

	if err = w.WriteInt(mk_signatures); err != nil {
		return err
	}
	if err = w.WriteArrayStart(len(a.Signatures)); err != nil {
		return err
	}
	for _, sig := range a.Signatures {
		if err = sig.Emit(w); err != nil {
			return err
		}
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}

	if err = w.WriteInt(mk_subject_name); err != nil {
		return err
	}
	if err = w.WriteString(a.Name); err != nil {
		return err
	}

	if bare {
		if err = w.WriteInt(mk_subject_zone); err != nil {
			return err
		}
		if err = w.WriteString(a.Zone); err != nil {
			return err
		}

		if err = w.WriteInt(mk_context); err != nil {
			return err
		}
		if err = w.WriteString(a.Context); err != nil {
			return err
		}
	}

	if err = w.WriteInt(mk_objects); err != nil {
		return err
	}
	if err = w.WriteArrayStart(len(a.Objects)); err != nil {
		return err
	}
	for _, o := range(a.Objects) {
		if err = o.Emit(w); 
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}

	if err = w.WriteMapEnd(); err != nil {
		return err
	}

	return nil
}

func (a *Assertion) EmitSection(w *CBORWriter) error {
	var err error 

	if err = w.WriteArrayStart(2); err != nil {
		return err
	}
	if err = w.WriteInt(sk_assertion); err != nil {
		return err
	}
	if err = a.Emit(w, true); err != nil {
		return err
	}
	if err = w.WriteArrayEnd(); err != nil {
		return err
	}
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

	if err = w.WriteArrayStart(2); err != nil {
		return err
	}

	if as.ZoneComplete {
		sectionKey = sk_zone
	} else {
		sectionKey = sk_shard
	}

	if err = w.WriteInt(sectionKey); err != nil {
		return err
	}

	if err = w.WriteMapStart(); err != nil {
		return err
	}

	// WORK POINTER

	if err = w.WriteMapEnd(); err != nil {
		return err
	}


}


type Query struct {
	Name        string
	Context     string
	ObjectTypes map[ObjectType]bool
}

type Notification struct {
	NoteType NotificationType
	NoteData string
}

type Message struct {
	Assertions   []Assertion
	Shards       []AssertionSet
	Zones        []AssertionSet
	Capabilities []string
	token        string
}
