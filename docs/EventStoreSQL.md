# Event Store SQL

## Race condition

SQL works well for Event storage, but not Event replay.  There is no default way to ensure Event order.  More importantly, there is a race condition in writing and reading Events.

Let's say we have three processes, two writers: `W1` and `W2`, and one reader: `R`.

`W1` and `W2` are both writing Events to the log.  `W1` writes `Event #101` and and `W2` writes `Event #102`.  The SQL server generates the Event positions `#101` and `#102` serially, and so the events are well ordered.

However, it is possible for `W2` to finish writing before `W1`.  Therefore, for a brief time, `Event #102` exists, but `Event #101` does not.  If `R` reads during this time, it will see a gap in the Event log.  Furthermore, if `W1` fails, the process will complete, and `Event #101` may never exist, and the gap will be permanent.  Thus, `R` cannot assume the gap will ever be filled.

`R` has two options:

1. Re-run the query repeatedly for some time window, waiting for the gap to be filled
2. Move on, never to return

However, `R` may not be able to detect the gap.  If our Event log is partitioned, contiguous Events may be split into different partitions.  Thus, the only option is 2, which is unacceptable for Event replays.

### Solution

We need to establish a **High Water Mark**, which is some maximum Event log position that is safe to read.  Any Event before that position is confirmed to be free of any race conditions.

### PostgreSQL solution

PostgreSQL provides this via a unique **Transaction ID**.  Every transaction is automatically provided a monotonically increasing identifier.  Then, when each transaction is committed, an internal **High Water Mark** is updated.  If we store this ID on each Event, we can simply read any Event before the current **High Water Mark**.

We can read transaction information with the following SQL functions

| Function                                  | Information                               |
| ----------------------------------------- | ----------------------------------------- |
| `pg_current_xact_id()`                    | Transaction ID of the current transaction |
| `pg_snapshot_xmin(pg_current_snapshot())` | Transaction ID of the **High Water Mark** |

## References

* https://www.postgresql.org/docs/current/functions-info.html#FUNCTIONS-INFO-SNAPSHOT
* https://event-driven.io/en/ordering_in_postgres_outbox/
* https://github.com/eugene-khyst/postgresql-event-sourcing#4-7-1
