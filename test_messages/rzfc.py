#!/usr/bin/env python3

# RAINS Zonefile Compiler
# Turns RAINS zonefiles into CBOR, and vice versa
# Prototype/testing; to be incorporated into RAINSD

import re
import cbor
from ipaddress import ip_address
from collections import namedtuple

K_SIGNATURES     = 0
K_CAPABILITIES   = 1 
K_TOKEN          = 2 
K_SUBJECT_NAME   = 3 
K_SUBJECT_ZONE   = 4 
K_QUERY_NAME     = 5 
K_CONTEXT        = 6 
K_OBJECTS        = 7 
K_QUERY_CONTEXTS = 8
K_QUERY_TYPES    = 9
K_QUERY_OPTS     = 10
K_SHARD_RANGE    = 11
K_NOTE_TYPE      = 21
K_NOTE_DATA      = 22
K_CONTENT        = 23 

SEC_ASSERTION     = 1 
SEC_SHARD         = 2 
SEC_ZONE          = 3 
SEC_QUERY         = 4 
SEC_NOTIFICATION  = 23

Token = namedtuple("Token", "t", "v")

scanner = re.Scanner([
    (r"\(",             lambda s,t:(Token(t,                      None))),
    (r"\)",             lambda s,t:(Token(t,                      None))),
    (r"\[",             lambda s,t:(Token(t,                      None))),
    (r"\]",             lambda s,t:(Token(t,                      None))),
    (r",",              lambda s,t:(Token(t,                      None))),
    (r":Z:\s+",         lambda s,t:(Token(t,                      None))),
    (r":S:\s+",         lambda s,t:(Token(t,                      None))),
    (r":A:\s+",         lambda s,t:(Token(t,                      None))),
    (r":sig:\s+",       lambda s,t:(Token(t,                      None))),
    (r":ip4:\s+",       lambda s,t:(Token(t,                      None))),
    (r":ip6:\s+",       lambda s,t:(Token(t,                      None))),
    (r":name:\s+",      lambda s,t:(Token(t,                      None))),
    (r":deleg:\s+",     lambda s,t:(Token(t,                      None))),
    (r":redir:\s+",     lambda s,t:(Token(t,                      None))),
    (r":cert:\s+",      lambda s,t:(Token(t,                      None))),
    (r":infra:\s+",     lambda s,t:(Token(t,                      None))),
    (r":regr:\s+",      lambda s,t:(Token(t,                      None))),
    (r":regt:\s+",      lambda s,t:(Token(t,                      None))),
    (r":srv:\s+",       lambda s,t:(Token(t,                      None))),
    (r"\d+\.\d+\.\d+\.\d+",
                        lambda s,t:(Token("VAL_IP4",     ip_address(t)))),
    (r"::",
                        lambda s,t:(Token("VAL_IP6",     ip_address(t)))),
    (r"([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}",
                        lambda s,t:(Token("VAL_IP6",     ip_address(t)))),
    (r"([0-9a-fA-F]{1,4}:){0,7}:([0-9a-fA-F]{1,4}:){0,7}[0-9a-fA-F]{1,4}",
                        lambda s,t:(Token("VAL_IP6",     ip_address(t)))),
    (r"\d+-\d+-\d+T\d+:\d+:\d+",
                        lambda s,t:(Token("VAL_8601",    
                            (datetime.strptime(t,"%Y-%m-%dT%H:%M:%S") - 
                             datetime(1970,0,1)).total_seconds()))),
    (r"[1-9][0-9]*",
                        lambda s,t:(Token("VAL_INT",                 t))),
    (r"[0-9a-fA-F]{2,128}",
                        lambda s,t:(Token("VAL_HEX",                 t))),
    (r":[a-zA-Z][a-zA-Z0-9_]*:\s+",
                        lambda s,t:(Token("TOK_RESERVED",         None))),
    (r"\S+",            lambda s,t:(Token("VAL_STRING",              t))),
    (r"\s+",            lambda s,t:None),
    (r"#\.*\n",         lambda s,t:None)
])


def is_string(tok):
    return tok.t == "VAL_STRING" || tok.t == "VAL_HEX" || tok_t == "VAL_INT"

def consume_symbol(ts, sym):
    if ts[0].t != sym:
        raise ValueError("expected "+str(sym)+", got "+str(ts.v))
    return ts[1:]

def section(ts):
    if ts[0].t == ":S:":
        if !is_string(ts[1]) || !isstring(ts[2]) :
            raise ValueError("missing zone and/or context in :Z:")
        z, ts = zone(ts[3:], ts[1].v, ts[2].v)
        return [ SEC_ZONE, z ]
    elif ts[0][0] == ":S:":
        if !is_string(ts[1]) || !isstring(ts[2]) :
            raise ValueError("missing zone and/or context in bare :S:")
        s, ts = shard(ts[3:], ts[1].v, ts[2].v, True)
        return [ SEC_SHARD, s], ts
    elif ts[0][0] == ":A:":
        if !is_string(ts[1]) || !isstring(ts[2]) :
            raise ValueError("missing zone and/or context in bare :A:")
        a, ts = assertion(ts[3:], ts[1].v, ts[2].v, True)
        return [ SEC_ASSERTION, a ]
    else:
        raise ValueError("expected :Z:, :S:, or :A:")

def zone(ts, zone_name, context_name):
    out = { K_ZONE_NAME:    zone_name
            K_CONTEXT:      context_name
            K_CONTENT:      [] }

    ts = consume_symbol(ts, "[")

    # eat content
    while ts[0].t != "]":
        if ts[0].t == ":S:":
            s, ts = shard(ts[1:], zone_name, context_name, False)
            out[K_CONTENT].append(s)
        elif ts[0].t != ":A:":
            a, ts = assertion(ts[1:], zone_name, context_name, False)
            out[K_CONTENT].append(a)
        else:
            raise ValueError("expected :S:, :A:, or ]")
    ts = consume_symbol(ts, "]")

    # and signatures, if present
    out[K_SIGNATURES], ts = signatures(ts)

    return out, ts

def shard(ts, zone_name, context_name, is_section):
    out = { K_ZONE_NAME:    zone_name,
            K_CONTEXT:      context_name,
            K_SIGNATURES:   []
            K_CONTENT:      []
            K_SHARD_RANGE:  [] }

    # check for range
    if ts[0].t == "(":
        ts = ts[1:]
        if is_string(ts[0]):
            out[K_SHARD_RANGE].append(ts[0].v)
            ts = consume_symbol(ts[1:],",")
        elif ts[0].t == ",":
            out[K_SHARD_RANGE].append(None)
            ts = ts[1:]
        else:
            raise ValueError("expected shard range begin or ,")

        if is_string(ts[0]):
            out[K_SHARD_RANGE].append(ts[0].v)
            ts = ts[1:]
        elif (ts[0].t == ")"):
            out[K_SHARD_RANGE].append(None)
        else:
            raise ValueError("expected shard range end or )")
        ts = consume_symbol(ts, ")")

    ts = consume_symbol(ts, "[")

    # eat content
    while ts[0].t != "]":
        if ts[0].t != ":A:":
            a, ts = assertion(ts[1:], zone_name, context_name)
            out[K_CONTENT].append(a)
        else:
            raise ValueError("expected :A:, or ]")
    ts = consume_symbol(ts, "]")    

    # and signatures, if present
    out[K_SIGNATURES], ts = signatures(ts)

    return out, ts

def assertion(ts, zone_name, context_name, is_section):
    pass

def signatures(ts)
    out = []

    # check for signature
    if ts[0].t == "(":
        ts = ts[1:]
        while ts[0].t != ")":
            s, ts = signature(ts)
            out.append(s)
        ts = consume_symbol(")")

    return out, ts

def signature(ts):
    pass

test_zone_1 = """
:Z: example.com . [
    :S: [
        :A: _smtp._tcp [ :srv: mx 25 10 ]
        :A: foobaz [
            :ip4: 192.0.2.33
            :ip6: 2001:db8:cffe:7ea::33
        ]
        :A: quuxnorg [
            :ip4: 192.0.3.33
            :ip6: 2001:db8:cffe:7eb::33
        ]
    ]
]
"""