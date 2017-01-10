# rainsd

This repository contains a reference implementation of a RAINS server,
supporting authority, intermediary, and query service, as well as associated
tools for managing the server. It has the following entry points:

- `rainsd`:   A RAINS server
- `rainsdig`: A command-line RAINS client for query debugging
- `rainspub`: A command-line RAINS authority client for 
              publishing assertions to an authority service
- `rainsfic`: A tool for generating assertions to publish for
              publishing with `rainspub` for testing purposes.

In addition, the `rainslib` library on which the server and tools are built
provides common handling for the CBOR-based RAINS information model. There may
be a `rainsclientlib` as well for client-side name resolution; bonus points
for integrating this with the Go net package in an idiomatic way.

## rainsd architecture and design

The RAINS server itself is made up of several components:

- A server query engine (`engine.go`): This component stores the assertions
  the server knows about, the queries it doesn't yet have answers to, and the
  chains of keys necessary to verify assertions in each zone.
- A message processor (`inbox.go`): This component processes incoming
  messages, demarshaling them, verifying their signatures, and handing them to
  the query engine for further processing.
- A switchboard (`switchboard.go`): The other components of rainsd operate in
  terms of messages associated with a RAINS server identity. The switchboard
  maintains open connections to other RAINS servers using a least-recently used
  cache, reopening connections that have timed out.

In addition, the RAINS server uses the following component provided by `rainslib`:

- A data model implementation (`model.go`): This component defines the core runtime data types,
  handles marshaling and unmarshaling of RAINS messages into and out of CBOR, the parsing and
  generation of RAINS zonefiles, and utilities for signing assertions and verifying signatures.

### query engine design

The query engine is built around three tables: an *assertion cache*, a
*pending queries cache*, and an *authority table*.

The assertion cache stores assertions this instance knows about, indexed by
the fields in a query: context, zone, name, object type. Each entry has an
expiration time, derived from the last-expiring signature on an assertion, to
aid in reaping expired entries. The assertion cache should be able to return
both a value (for internal queries, i.e. when the RAINS server itself needs to
know another RAINS server's address to connect to it) as well as an assertion
(in order to answer queries from peers).

Assertions in the assertion cache are assumed to have valid signatures at the
time of insertion into the assertion cache; it is the responsibility of the
message processor to check for validity (and to handle errors for invalid
signatures).

Note that the assertion cache is a candidate for refactoring into `rainslib`
as it might be useful in the RAINS client implementation, as well.

The pending queries cache stores unexpired queries for which assertions are not yet
available. A query that cannot be immediately answered is placed into the pending queries cache, and checked against incoming assertions until it expires.

The authority table keeps a mapping of zone names to keys, for verifying signatures. *Open Issue*: this needs to be made available to the message processor, or separated out. 

- - - -
work pointer
- - - - 

The query engine has a simple API, with three entry points:

- `assert(assertion)`: add a signed assertion to the assertion cache. Trigger
   any pending queries answered by it. Add keys to the authority table if the
   assertion is a delegation.
- `query(query, callback)`: Run the
   specified callback when the query is answerable. Do so immediately on an assertion cache hit, or after an assertion is available 
- `reap()`: remove expired queries, assertions, and delegations. This is
   probably simply called by a goroutine waiting on a tick channel.

The design of the internal data structures for the query engine is separate
from that of the `rainslib` data model. The `rainslib` data model is optimized
to be close to the wire, and easy to marshal/unmarshal to and from CBOR. The
query engine structures are optimized for fast access given a key (name to
zone, name to contexts, context/zone/name/type(s) to assertions and/or
queries). The query engine structures point either to raw assertions or raw
queries in the `rainslib` data model, as "provenance" for a given question or answer. 

Care must be taken in this design to handle nonexistence proofs based on shards efficiently. 


The following protocol features still need to be supported by the prototype query engine:

- Nonexistence proofs based on shards/zones: how does the query cache know the difference between "no assertion exists" and "I don't have this assertion"? We probably don't want to build this directly into the main assertion cache. Suggested approach: range index data structure that keeps ranges by zone and context, extracted from shards. This structure is queried on assertion cache miss.

#### short assertions and short queries

There is a fair amount of complexity involved in marshaling and unmarshaling CBOR as defined in the RAINS protocol draft (see [the datamodel](#datamodel) for more details). Prototyping will therefore work on "short assertions" and "short queries" instead.

An unsigned short assertion is a UTF-8 string of the form "A valid-from valid-until context zone subject objtype value" where:

- valid-from is an ISO8601 timestamp
- valid-until is an ISO8601 timestamp
- context is the context of the assumption
- zone is the name of the subject zone
- subject is the subject name within the zone
- objtype is one of:
    - ip4 for an IPv4 address; value is parseable address in string form
    - ip6 for an IPv6 address; value is parseable address in string form
    - name for a name; value is name as string
    - deleg for a delegation; value is cipher number, space, delegation key as hex string
    - redir for a redirection; value is authority server name
    - infra for an infrastructure key; value is cipher number, space, key as hex string
    - cert for a certificate; not yet implemented
    - nameset for a nameset; not yet implemented
    - regr for a registrar; value is unformatted string
    - regt for a registrant; value is unformatted string
    - srv for service info; not yet implemented
- value may contain spaces

A signed short assertion is generated and verified over the unsigned short
assertion with a valid key for that assertion's zone. A signed short assertion
has the form:

"S cipher-number signature unsigned-assertion"

A short query has the form:

"Q valid-until context subject objtype"

(Note that unlike RAINS queries, short queries can only have a single context
and object-type. This simplification may carry over into the protocol.)

### data model marshaling and unmarshaling design {#datamodel}

looks like we have to write our own CBOR serialization/deserialization due to two complications:

- RAINS specifies a normalization for serialized CBOR for signing that CBOR libraries may not honor.
- RAINS specifies integer keys for extensible maps for efficiency, and integer keys are hard/impossible to do with structure tags.

One could/should hack an existing CBOR library to provide these two properties.


## rainspub design

A RAINS server cannot be tested unless fed with validly signed assertions. rainspub takes input in the form of RAINS zonefiles (see [the zonefile definition](#zonefiles))

### zonefile Format {#zonefile}

todo: describe the rains zonefile format here. inspired by BIND zonefiles, close to the wire format, and designed to be easily RDP-parseable.
