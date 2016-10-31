import cbor
import json
import time 

import nacl.encoding
import nacl.signing 

from ipaddress import ip_address

CONTENT       = 0 
CAPABILITIES  = 1 
SIGNATURES    = 2 
SUBJECT_NAME  = 3 
SUBJECT_ZONE  = 4 
QUERY_NAME    = 5 
CONTEXT       = 6 
OBJECTS       = 7 
TOKEN         = 8 
SHARD_RANGE   = 11
QUERY_TYPES   = 14
NOTE_TYPE     = 17
QUERY_OPTS    = 22
NOTE_DATA     = 23

ASSERTION     = 1 
SHARD         = 2 
ZONE          = 3 
QUERY         = 4 
NOTIFICATION  = 23

NAME          = 1 
IP6_ADDR      = 2 
IP4_ADDR      = 3 
REDIRECTION   = 4 
DELEGATION    = 5 
NAMESET       = 6 
CERT_INFO     = 7 
SERVICE_INFO  = 8 
REGISTRAR     = 9 
REGISTRANT    = 10
INFRAKEY      = 11

ECDSA_256     = 2
ED25519       = 4

RAINS_TAG_CHEAT = b'\xda\x00\xe9\x9b\xa8'

test_message = {
    CONTENT: [
        [ ASSERTION, {
            SUBJECT_NAME : "a",
            SUBJECT_ZONE : "example.com.",
            CONTEXT :      ".",
            OBJECTS : [
                [ IP4_ADDR, ip_address("192.0.2.33").packed ],
                [ IP4_ADDR, ip_address("192.0.3.33").packed ],
                [ IP6_ADDR, ip_address("2001:db8:0:2::33").packed ],
                [ IP6_ADDR, ip_address("2001:db8:0:3::33").packed ]
            ],
            SIGNATURES : []
        }]
    ]
}

def sign_assertions_in_message(message, ttl, zone_keys):

    for contents in message[CONTENT]:
        if contents[0] == ASSERTION:
            assertion = contents[1]
            # save signatures
            saved_signatures = assertion[SIGNATURES]

            # calculate validity
            time0 = int(time.time())
            time1 = time0 + ttl

            # create stub signature
            assertion[SIGNATURES] = [ [ ED25519, time0, time1, 0, None ] ]

            # serialize object to sign
            bytes_to_sign = cbor.dumps(assertion)

            # sign object
            signing_key = zone_keys[assertion[SUBJECT_ZONE]][0]
            signature = signing_key.sign(bytes_to_sign)[:64]

            # add signature back to message
            assertion[SIGNATURES][0][4] = signature

            # add saved signatures
            assertion[SIGNATURES] = saved_signatures + assertion[SIGNATURES]

def generate_example_com_key(filename):

    sk = nacl.signing.SigningKey.generate()
    vk = sk.verify_key

    zk = { "example.com." : [ sk.encode(encoder=nacl.encoding.HexEncoder).decode(), 
                             vk.encode(encoder=nacl.encoding.HexEncoder).decode() ] }

    with open(filename, mode="w") as f:
        json.dump(zk, f)

def load_zone_keys(filename):
    with open(filename) as f:
        rzk = json.load(f)

    zk = {}
    for z in rzk:
        sk = nacl.signing.SigningKey(rzk[z][0], encoder=nacl.encoding.HexEncoder)
        vk = nacl.signing.VerifyKey(rzk[z][1], encoder=nacl.encoding.HexEncoder)
        zk[z] = [sk, vk]

    return zk



