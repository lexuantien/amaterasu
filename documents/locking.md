# Locking
- These are stategies to handle contention for resources.
- There are two general models to deal with these situations that you can observe in the wild:
    - Pessimistic. You take a time-bombed lock on an object and reserve it for the client. Ticketmaster operates their web site that way. When they show you a set of seats tickets they are yours until the clock ticks down.
    - Optimistic. There's no lock. Expedia/Orbitz work that way. When they show you a seat, you don't get it until you book it and there is no promise made.
