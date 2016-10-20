# rainsd

This repository contains a reference implementation of a RAINS server,
supporting authority, intermediary, and query service, as well as associated
tools for managing the server. It has the following entry points:

- `rainsd`:    A RAINS server
- `rainsdig`:  A command-line RAINS client for query debugging
- `rainspub`:  A command-line RAINS authority client for 
               publishing assertions to an authority service

In addition, the `rainslib` library on which the server and tools are built provides common handling for the CBOR-based RAINS information model.