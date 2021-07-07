The boolean expression indexer lib for audience targeting in online advertisement.



- scoring pruning (Done)
- logging monitoring 
- lock for RW indexing
    - is it really necessary? roll-out deployment may be better, not hurting performance

- HA configuration distributed
    - partitioning by geo or other key?
    - wal log, deep storage 
- http, grpc api
- web UI
- plugin: CDC read

