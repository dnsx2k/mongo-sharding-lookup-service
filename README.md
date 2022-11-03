## Lookup Service  

### Overview

Piece of "Hot Shard" project which is a custom MongoDB sharding solution.

Lookup Service is responsible for storing shard keys and their location.

### API

#### GET v1/lookup/:key
Will fetch key:location pair

#### POST v1/lookup
Will add new entry with given location for many keys

