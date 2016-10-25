import cbor
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
            ]
            SIGNATURES : []
        }]
    ]
}