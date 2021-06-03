# AMATERASU

Amaterasu is simple CQRS and Eventsourcing framework in distributed system. It's written in [golang](https://golang.org/), using [kafka](https://kafka.apache.org/) as message queue and [mysql]() database. Amaterasu is 55% sourced from CQRS journey project Microsoft and did from 2012-2015. The remaining 45% is added from blog posts of tiki, uber, wix on medium, youtube...

# Why name amaterasu ?
Amaterasu is considered not only the god of the sun, but also the god (center) of the universe. Like CQRS & Eventsourcing - the heart of a distributed system.

# Features
- [x] Create [command](#Commands) and [event](#Events)
- [x] Command handler and Event handler handle command and event
- [x] CommandBus and EventBus produce  command and event
- [x] Using ORM to connect mysql database
- [x] Save event in eventsourcing to `event_store` table
- [x] Save messages to `undispatched_message` table before produce to kafka
- [x] Sagas (Process manager) to handle messages from different bounded context
- [x] Order message
- [ ] Snapshot `event_store` using memeto pattern
- [ ] Monitor message
- [ ] Idempotent and duplicate message
- [ ] Mysql bindlog to produce message

# CQRS 
    CQRS is a simple pattern that strictly segregates the responsibility of handling command input into an autonomous system from the responsibility of handling side-effect-free query/read access on the same system. Consequently, the decoupling allows for any number of homogeneous or heterogeneous query/read modules to be paired with a command processor and this principle presents a very suitable foundation for event sourcing, eventual-consistency state replication/fan-out and, thus, high-scale read access. In simple terms: You don't service queries via the same module of a service that you process commands through. For REST heads: GET wires to a different thing from what PUT/POST/DELETE wire up to.