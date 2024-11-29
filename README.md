# ledis-server

---
- [Overview of Redis](#overview-of-redis)
- [Which commands do we need to implement](#which-commands-do-we-need-to-implement-)
  - [String Commands](#string-commands)
  - [List Commands](#list-commands)
  - [Set Commands](#set-commands)
  - [Data Expiration Commands](#data-expiration-commands)
  - [Snapshot Commands](#snapshot-commands)
- [Why Golang ?](#why-golang-)
- [Initial Ideas](#initial-ideas)
- [Designs](#designs)
- [Things can be improved more](#things-can-be-improved-more)
---
### Overview of Redis
- Redis is a well-known in-memory data store that supports many features like key/value store, set, list, and many other data structures.
- Also Redis has TTL (Time to live) which mean that one key will have a duration of time to live and after that duration, key will automatically expire to users (which may be good if we have stale data and do not want to keep it for so long).
- Two popular in-memory data stores are Redis and Memcached.
- What makes Redis different from Memcached is that Redis uses a single-threaded architecture (i.e., its server runs on a single thread). This makes it easier to manage and simplifies handling concurrency, as there are no race conditions due to the single-threaded model.
- **Redis** also support multiple nodes so make it a highly available cluster.

### Which commands do we need to implement ?
#### String Commands:
- SET key value: Set the string value associated with the specified key, overwriting any existing value.
- GET key: Retrieve the string value stored at the specified key.

#### List Commands:
A List is an ordered collection of strings, where duplicates are allowed.

- `LLEN key`: Return the length of the list stored at the specified key.
- `RPUSH key value1 [value2...]`: Append one or more values to the list. If the list doesn't exist, it is created. Returns the length of the list after the operation.
- `LPOP key`: Remove and return the first item from the list stored at the specified key.
- `RPOP key`: Remove and return the last item from the list stored at the specified key.
- `LRANGE key` start stop: Return a range of elements from the list, inclusive of the start and stop indices. Both start and stop are zero-based non-negative integers.

#### Set Commands:
A Set is an unordered collection of unique string values (duplicates are not allowed).

- `SADD key value1 [value2...]`: Add one or more values to the set stored at the specified key.
- `SCARD key`: Return the number of elements in the set stored at the specified key.
- `SMEMBERS key`: Return an array of all the members of the set stored at the specified key.
- `SREM key value1 [value2...]`: Remove one or more values from the set stored at the specified key.
- `SINTER key1 [key2...]`: Return the intersection of sets stored at the specified keys.

#### Data Expiration Commands:
- `KEYS`: List all available keys in the database.
- `DEL key`: Delete the specified key.
- `FLUSHDB`: Clear all keys from the database.
- `EXPIRE key seconds`: Set a time-to-live (TTL) for the specified key. Returns the number of seconds remaining until the key expires.
- `TTL key`: Retrieve the remaining time-to-live (TTL) of the specified key.

#### Snapshot Commands:
- `SNAPSHOT`: Create a snapshot of the current state of the database.
- `RESTORE`: Restore the database from the last saved snapshot.

### Why Golang ?
- Actually, I am more familiar with Golang, I used it a lot to work with other projects.
- One more reason is that, I have an improvement for this mini project that I have made not a long time before also using Golang. So I think using the same language may make these 2 projects has the consistency.

### Initial Ideas
- At first, after decided to choose **Golang** as the main programming language for this project. What I need to think about is that how to implement these commands. And the **important things** to me is that: ***How can we make sure the correctness of my code and how to extend if some commands need to add more later.***
- Yeah, so our main 2 problems here are about the **Maintenance** and **Correctness** (I will skip about performance here because this project is quite simple to monitor the performance, and I think it will faster than any kind of in-memory storage, because of its **simplification**).
- Actually I have a thinking for the **Maintenance** like this: ***If one feature need to add to the application, if you open some files, and edit them. Then this tends to be the bad design. But if you open some files, and write more code to it, then it tends to be a good design***. This is not absolutely correct or wrong, because it is just my experiences after writing codes. So don't be harsh for this 😁.
- So after thinking for half of a day, I have figure out some designs for this and I think it can be well maintained later on (Refer to the below sections for [Designs](#designs)).
- And for `Data Expiration Commands`, I have several ways to implement this, e.g: We can `lazy-check` when the `key` is queried by users, and can remove it if it is expired. Or we can have some algorithm to periodically run to remove the `expiration keys`. And in this, I have implement these 2 solutions. For the first one, It is easy to understand. But the latter one, how we design algorithm. Luckily, I have read `Redis documentation` a long time ago and know that `Redis` have a probability algorithm to select a small subset of expiration keys and remove them from current memory. 
- For `Snapshot Commands`, if we encounter a `SNAPSHOT` command, we just need to store current in-memory data to `a file`. There is a tricky case for this. Assume that we store the `data` in file name `data.rdb`, but then we call `SNAPSHOT` again, and unluckily, this time the `SNAPSHOT` failed in the middle of time when it writes to `data.rdb`, and our last snapshot for `data.rdb` will be lost and corrupted for now. So I have researched and figured out one solution like this. We will write to temporaty file called `temp_data.rdb`, and when the `write` complete, we can use the `mv temp_data.rdb data.rdb` command to make `temp_data` become the latest `snapshot` for `ledis`. Why it is correct ? because `mv` operation is atomic in most file system architecture. So our `SNAPSHOT` command is also atomic too.
- 
### Designs

### Things can be improved more