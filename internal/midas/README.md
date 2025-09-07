# Validation Library
*Everything it touches turns to gold*

Define metadata for attributes using struct tags

Link to other mediums (Disk... SQL, NoSQL, Objects)

Schema Support
    Avro
    JSON Schema
    SQL

Schema registry
    Index schemas, noting version

DTO and context based subsets
    View, Request, Response



CLIENT ---(serialized)---> 
    UNMARSHALL ---(structured)---> 
    VALIDATE ---(sanitized/cleaned and validated)-->
    ENCRYPT ---(secured)--> 
    STORE ---(persisted)--> done

SERVER ---(request)-->
    ROLECHECK ---(authenticated)--->
    RETRIEVE ---(structured)--->
    DECRYPT ---(sometimes/rarely)--->
    MARSHALL ---(serialized)---> done