#!/usr/bin/env python3

# RAINS Zonefile Compiler
# Turns RAINS zonefiles into CBOR, and vice versa
# Prototype/testing; to be incorporated into RAINSD

import re
import cbor
from ipaddress import ip_address

scanner = re.Scanner([
    (r"\(",         lambda s,t:("TOK_(",                None)),
    (r"\)",         lambda s,t:("TOK_)",                None)),
    (r"\[",         lambda s,t:("TOK_[",                None)),
    (r"\]",         lambda s,t:("TOK_]",                None)),
    (r",",          lambda s,t:("TOK_,",                None)),
    (r":Z:",        lambda s,t:("TOK_ZONE",             None)),
    (r":S:",        lambda s,t:("TOK_SHARD",            None)),
    (r":A:",        lambda s,t:("TOK_ASSERTION",        None)),
    (r":sig:",      lambda s,t:("TOK_SIGNATURE",        None)),
    (r":ip4:",      lambda s,t:("TOK_IP4_ADDR",         None)),
    (r":ip6:",      lambda s,t:("TOK_IP6_ADDR",         None)),
    (r":name:",     lambda s,t:("TOK_NAME",             None)),
    (r":deleg:",    lambda s,t:("TOK_DELEGATION",       None)),
    (r":redir:",    lambda s,t:("TOK_REDIRECTION",      None)),
    (r":cert:",     lambda s,t:("TOK_CERTIFICATE",      None)),
    (r":infra:",    lambda s,t:("TOK_INFRAKEY",         None)),
    (r":regr:",     lambda s,t:("TOK_REGISTRAR",        None)),
    (r":regt:",     lambda s,t:("TOK_REGISTRANT",       None)),
    (r":srv:",      lambda s,t:("TOK_SERVICE",          None)),
    (r":[a-zA-Z][a-zA-Z0-9_]*:",
                    lambda s,t:("TOK_RESERVED",          None)),
    (r"\d+\.\d+\.\d+\.\d+",
                    lambda s,t:("VAL_IP4",     ip_address(t))),
    (r"::",
                    lambda s,t:("VAL_IP6",     ip_address(t))),
    (r"([0-9a-fA-F]{1,4}:){8}",
                    lambda s,t:("VAL_IP6",     ip_address(t))),
    (r"([0-9a-fA-F]{1,4}:){0,7}::[0-9a-fA-F]{1,4}(:[0-9a-fA-F]{1,4}){0,7}",
                    lambda s,t:("VAL_IP6",     ip_address(t))),
    (r"\d+-\d+-\d+T\d+:\d+:\d+",
                    lambda s,t:("VAL_8601",    
                        (datetime.strptime(t,"%Y-%m-%dT%H:%M:%S") - 
                         datetime(1970,0,1)).total_seconds())),
    (r"[0-9a-fA-F]{2,128}",
                    lambda s,t:("VAL_HEX",     t)),
    (r"\S+",        lambda s,t:("VAL_STRING",  t))    
    (r"\s+",        lambda s,t:None),
    (r"#\.*\n",     lambda s,t:None)
])

test_zone_1 = """
:Z: example.com . [
    :S: ( , ) [
        :A: _smtp._tcp [ :srv: mx 25 10 ]
        :A: aaa [
            :ip4: 192.0.2.33
            :ip6: 2001:db8:cffe:7ea::33
        ]
        :A: aab [
            :ip4: 192.0.2.33
            :ip6: 2001:db8:cffe:7ea::33
        ]
    ]
    :S: (aaa_)

]
"""