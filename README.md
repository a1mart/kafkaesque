# Kafka-esque
Kafka-like message service with integrated Schema Registry and Validator. Defined in Protobuf with automatically generated OpenAPI. Build around [LMAX Disruptor](https://github.com/a1mart/lmax).

## ToDo
- integrate schema registry/validator with messaging service
- Connect (services --> persistent databases and API services... ETL & integrations)... 'sources' and 'sinks'
- Streaming
- Load balancing/timing algorithms to balance event traffic with limited resources
- Security (TLS/SSL, JWT, OAuth, RBAC/ACL... mTLS)
- Observability and Monitoring (via Prometheus/Grafana stack)
- Testing; CI/CD Pipeline
- Ringbuffer per partition

# Commands
```bash
make run

make producer topic=test

make consumer topic=test

make create-topic TOPIC=my_topic STRATEGY=round_robin

make list-topics
```

# Integrated SwaggerUI served on http://localhost:8080/swagger

# Turn proto cmds into makefile command
protoc -I. -Ithird_party \        
  --go_out=. \
  --go-grpc_out=. \
  --grpc-gateway_out=. \
  --grpc-gateway_opt=logtostderr=true \
  --openapiv2_out=./swagger \
  messaging.proto

# Groups
Server <--> Client
Producer --> Consumer
Publisher --> Subscriber (broadcast/fanout)

# Schema validation
- JSON Schema
- Avro
- Protobuf
... custom validator (reflection)

`Ring buffer` circular array with producering writing entries ar sequence index and consumers reading entries at lower sequence index
`Sequence tracking`
- Producer (write) sequence, next slot producer will claim
- Consumer (read) sequence, highest slot producer has fully processed
`Back pressure` 
if the producer's next sequence has 'wrapped around' and would overwrite unconusmed data. Producer must wait (or block or drop messages) until the consumer moves forward


Specify Topic as FCFS (first in consumer group) or PubSub (broadcast to all)
  Specify whether or not to require validation for producers (register schemas on topic and only valid messages may be published)
  Specify connectors (with implied ordering)
    Associating handlers to topics... consumers... work on streams without consuming
        Plug and play algorithms and modules
          Providers
      Generalizing CRUD operations on generic type <T> 
  Streaming consumers (listening for messages conditionally)

Abstracting to support Agentic System

Distinction between processing, memory, and disk
  Scheduling CPU ()
    Hashing memory (caching)
      Structuring disk (sequentiality)
