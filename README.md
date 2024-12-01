# Ledis Server

---
- [Overview of Redis](#overview-of-redis)
- [Installation](#installation)
- [Which commands do we need to implement](#which-commands-do-we-need-to-implement-)
  - [String Commands](#string-commands)
  - [List Commands](#list-commands)
  - [Set Commands](#set-commands)
  - [Data Expiration Commands](#data-expiration-commands)
  - [Snapshot Commands](#snapshot-commands)
- [Why Golang ?](#why-golang-)
- [Initial Ideas](#initial-ideas)
  - [Some first thoughts about how to maintain](#some-first-thoughts-about-maintain)
  - [Expiration checking and strategy](#expiration-checking-and-strategy)
  - [Snapshot edge case](#snapshot-edge-case)
  - [Select the right strategy to implement each value data structure](#select-the-right-strategy-to-implement-each-value-data-structure)
- [Designs](#designs)
- [Things can be improved more](#things-can-be-improved-more)
  - [Transaction](#transactions)
  - [Snapshot](#snapshot)
---
### Overview of Redis
- Redis is a well-known in-memory data store that supports many features like key/value store, set, list, and many other data structures.
- Also Redis has TTL (Time to live) which mean that one key will have a duration of time to live and after that duration, key will automatically expire to users (which may be good if we have stale data and do not want to keep it for so long).
- Two popular in-memory data stores are Redis and Memcached.
- What makes Redis different from Memcached is that Redis uses a single-threaded architecture (i.e., its server runs on a single thread). This makes it easier to manage and simplifies handling concurrency, as there are no race conditions due to the single-threaded model.
- **Redis** also support multiple nodes so make it a highly available cluster.

### Installation
- Ensure that `Go` has installed on your machine, if not please refer to this [official document](https://go.dev/doc/install) to install it.
- Run the server by run the `make run_server` command in your terminal, if `make` has not installed yet u can install it, or just simply run the server directly using `go run main.go`
- You can run the tests as well `make test` or `go test -v ./...`.
- The server will start on the port `6379` same as `Redis`

### Which commands do we need to implement ?
#### String Commands:
- `SET key value`: Set the string value associated with the specified key, overwriting any existing value.
- `GET key:` Retrieve the string value stored at the specified key.

#### List Commands:
A List is an ordered collection of strings, where duplicates are allowed.

- `LLEN key`: Return the length of the list stored at the specified key.
- `RPUSH key value1 [value2...]`: Append one or more values to the list. If the list doesn't exist, it is created. Returns the length of the list after the operation.
- `LPOP key`: Remove and return the first item from the list stored at the specified key.
- `RPOP key`: Remove and return the last item from the list stored at the specified key.
- `LRANGE key start stop`: Return a range of elements from the list, inclusive of the start and stop indices. Both start and stop are zero-based non-negative integers.

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

#### Some first thoughts about maintain.
- At first, after decided to choose **Golang** as the main programming language for this project. What I need to think about is that how to implement these commands. And the **important things** to me is that: ***How can we make sure the correctness of my code and how to extend if some commands need to add more later.***
- Yeah, so our main 2 problems here are about the **Maintenance** and **Correctness** (I will skip about performance here because this project is quite simple to monitor the performance, and I think it will faster than any kind of in-memory storage, because of its **simplification**).
- Actually I have a thinking for the **Maintenance** like this: ***If one feature need to add to the application, if you open some files, and edit them. Then this tends to be the bad design. But if you open some files, and write more code to it, then it tends to be a good design***. This is not absolutely correct or wrong, because it is just my experiences after writing codes. So don't be harsh for this 😁.
- So after thinking for half of a day, I have figure out some designs for this and I think it can be well maintained later on (Refer to the below sections for [Designs](#designs)).

#### Expiration checking and strategy
- And for `Data Expiration Commands`, I have several ways to implement this.
- We can `lazy-check` when the `key` is queried by users, and can remove it if it is expired. 
- Or we can have some algorithm to periodically run to remove the `expiration keys`. 
- And in this, I have implement these 2 solutions. For the first one, It is easy to understand. 
- The latter one, I have several options:
  - Periodically run another thread to clean up the expiration keys. But the drawback is it will lock the current thread for a long time. Time complexity will be `O(N) for N is number of keys`. But it will work well if **large portion** of keys expire likely at the same time.
  - Maintain a `sorted set` (Balanced Binary Tree underlying) of all the `(expirationTime, key)` in `increasing orders`. Everytime we have a `new key with TTL set`, we can put it in this `sorted set`. And periodically we can get the smallest `expirationTime` out of the set in `O(logN) with N is the number of current keys in the set`. So it will reduce time complexity to `O(M * log(N)) with M is the number of expired keys, but it can go up to O(N * log(N)) if all the keys are expired`. But overall this solution may work well if we have **small portion** of expiration keys at the same time.
  - `Redis` has the random algorithm to select a small portion of keys and check to expire them. But currently what I am using is a map, so it's hard to select random keys from the map (can only iterate through, but this will not a random solution). 
    - I have think storing those keys in list, so we can select random indexes and expire them. But it may not work well because remove the keys from the list cost `O(N)` which we can achieve by an easier solution by iterating through the map.
    - Using a `SortedSet` is not a good idea, because we can know exactly which keys will expire, so we do not need to pick randomly here.
    - Using the `LinkList` here will solve the problem of removing element, but will be hard to get a specific key at one given index (have to iterate through), so it is still `O(N)` overall. 
    - There is one another data structure that store the `keys` along with `expiration time`, and we can remove one `key` in `O(key length)`, and we can **walk** through the children nodes (we can compute the probability of each children nodes that we will walk down). Total approximate complexity is `O(number of size * average length of keys)`
  - Above 4 solutions have either advantages and disadvantages and trade-off between time complexity and maybe memory complexity, and we should choose the one that matches application needs. But for easier in this application, I will choose the first one.
#### Snapshot edge case
- For `Snapshot Commands`, if we encounter a `SNAPSHOT` command, we just need to store current in-memory data to `a file`. There is a tricky case for this. Assume that we store the `data` in file name `data.rdb`, but then we call `SNAPSHOT` again, and unluckily, this time the `SNAPSHOT` failed in the middle of time when it writes to `data.rdb`, and our last snapshot for `data.rdb` will be lost and corrupted for now. So I have researched and figured out one solution like this. We will write to temporaty file called `temp_data.rdb`, and when the `write` complete, we can use the `mv temp_data.rdb data.rdb` command to make `temp_data` become the latest `snapshot` for `ledis`. Why it is correct ? because `mv` operation is atomic in most file system architecture. So our `SNAPSHOT` command is also atomic too.


#### Select the right strategy to implement each value data structure
- For the `Data Structure commands (String, List, Set)`. We can use a map with `key` as `String` and value can be `any`, but I will implement an interface for the `value` to make it extendable. 
- The `List data structure` is quite easy, we can use the builtin `list.List` of `Golang`. 
- The `Set data structure` is that we can use a map for this (to check if a key has existed or not). We can perform operations easily using map. But for easy, I will use the lib for this case.
- The `String`, we can use builtin `string`.

#### Testing
- For the ***Correctness***. Yes we need to write tests to make sure it can be more correct.

#### Keep coding
- So that's it, I will start implementing now and will figure out any problems and fix them.
### Designs
- I have designed several interfaces for this project:
```go
redis/interface.go


type Item interface {
	Value() any
	Type() int
}
```
- So as I said before `Item` is aa interface that will act as a value, any new data structures supported will implement this interface. For example I have a `List` concrete implementation like this:
```go
type ListType struct {
	List *list.List
}

func NewListType() redis.Item {
	return &ListType{
		List: list.New(),
	}
}

func (s *ListType) Value() any {
	return s.List
}

func (s *ListType) Type() int {
	return utils.ListType
}
```
- Actually I dont want to make `ListType` and `List` **publicly accessible**. But for making snapshot, I need it to be public so the process of making snapshot can touch the `List` inside. The same for other data structures.
- And in the `ListType` we may implement the `commands` that we support like this:

```go
func (s *ListType) LLen() int {
	return s.List.Len()
}

func (s *ListType) LPush(values ...*string) int {
	for _, v := range values {
		s.List.PushFront(*v)
	}
	return s.LLen()
}

func (s *ListType) RPush(values ...*string) int {
	for _, v := range values {
		s.List.PushBack(*v)
	}
	return s.LLen()
}
... There more
```

- Our main part is `Redis`, and this is its interface:
```go

type Redis interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()

	// Set key equal to value
	Set(key string, value Item)

	// Get just get but not check expiration
	Get(key string) (Item, bool)

	Delete(key string)

	// Expired check if key has expired or not
	Expired(key string) bool

	// GetOrExpired  lazy key expiration when get it
	GetOrExpired(key string) (Item, bool)

	// Keys return list of keys available (not expired or do not have ttl set)
	Keys() []string

	FlushDB()

	Expire(key string, ttlInSeconds int) error
	TTL(key string) (int, error)
	Gets(keys ...string) []Item

	LoadSnapshot() error
	MakeSnapshot() error
}
``` 

- Its concrete implementation must support all of the functions in this interface. For example:

```go
type redis struct {
	data          map[string]Item
	expirationKey map[string]time.Time
	ttl           map[string]time.Duration
	mu            *deadlock.RWMutex
}

func (s *redis) Lock() {
	s.mu.Lock()
}

func (s *redis) Unlock() {
	s.mu.Unlock()
}

func (s *redis) RLock() {
	s.mu.RLock()
}

func (s *redis) RUnlock() {
	s.mu.RUnlock()
}

func (s *redis) Set(key string, value Item) {
	s.data[key] = value
}

func (s *redis) Get(key string) (Item, bool) {
	value, exist := s.data[key]
	return value, exist
}

... There more
```
- As u can see, our `redis` implementation will store several maps data `data`, `expirationKey` and `ttl`.
- `data` map will store the value.
- `expirationKey` map will store the point in time when it will be expired.
- `ttl` map will store the ttl value of current `key`.

- And for `commands`, I have designed for it an interface look like this:
```go

type ICommandHandler interface {
	CommandName() string
	Execute(args ...string) (any, error)
}

```

- While `ICommandHandler` is an abstraction for those commands that will implement it. It should implement the `CommandName` which return a command unique name like `GET` `SET` `LRANGE` or sth. And the `Execute` function will receive many arguments for that corresponding command and execute it to make change with the `Redis` interface. Example of a `command` that implement the `ICommandHandler` looks like this:

```go

type getCmd struct {
	rds redis.Redis
}

func NewGetCmd(rds redis.Redis) redis.ICommandHandler {
	return &getCmd{rds: rds}
}

func (cmd *getCmd) CommandName() string {
	return "GET"
}

// Execute GET KEY
func (cmd *getCmd) Execute(args ...string) (any, error) {
	if len(args) != 1 {
		return nil, utils.ErrArgsLengthNotMatch
	}

	cmd.rds.Lock()
	defer cmd.rds.Unlock()

	item, exist := cmd.rds.GetOrExpired(args[0])

	if !exist {
		return nil, nil
	}

	if item.Type() != utils.StringType {
		return nil, utils.ErrTypeMismatch(utils.StringType, item.Type())
	}

	v := *item.Value().(*string)
	return v, nil
}
```
- It will keep the `Redis` interface inside, and when its `Execute` function will make change to the `Redis` interface through the exported functions.
- Finally the `ICommandHandler` will be registed to the `ICommandManager`, and the `ICommandManager` interface will map them with their `CommandName` and `Handler` for execution.

```go

type ICommandManager interface {
	Register(handler ICommandHandler) ICommandManager
	Execute(command string, args ...string) (any, error)
}

type commandManager struct {
    commandHandlerMapper map[string]redis.ICommandHandler
    rds                  redis.Redis
}

func (cm *commandManager) Register(handler redis.ICommandHandler) redis.ICommandManager {
    if _, ok := cm.commandHandlerMapper[handler.CommandName()]; ok {
        logging.GetLogger().Fatalln(utils.ErrCommandRegisteredDuplicate(handler.CommandName()))
    }
    cm.commandHandlerMapper[handler.CommandName()] = handler
    return cm
}

func (cm *commandManager) Execute(command string, args ...string) (any, error) {
    handler, ok := cm.commandHandlerMapper[command]
    if !ok {
        return nil, utils.ErrCommandDoesNotExist
    }
    return handler.Execute(args...)
}

```

- The concrete `commandManager` will implement these 2 functions `Register` and `Execute`. 
- `Register` just map the name of the `ICommandHandler` to itself.
- `Execute` will receive `command` and then find the corresponding `ICommandHandler` in the map and execute it with `args`.
- And when we want to have a new `command`, we just add the new `command` to the `Register` function in the custom constructor for `ICommandManager` looks like this:

```go

func NewCommandManager(rds redis.Redis) redis.ICommandManager {
	commandManager := &commandManager{
		commandHandlerMapper: make(map[string]redis.ICommandHandler),
	}
	createCommand := func(newFunc func(redis.Redis) redis.ICommandHandler) redis.ICommandHandler {
		return newFunc(rds)
	}

	commandManager.
		Register(createCommand(string_commands.NewGetCmd)).
		Register(createCommand(string_commands.NewSetCmd)).
		...
```


### Things can be improved more
#### Transactions
- `Redis` does not support real transaction, what it is supporting is just running a sequence of operations and ensure that no operation outside transaction can occur in the middle of transaction (it is easy because `Redis` is one thread). And if one operation fails in the middle of transaction, then it will continue execute other remain operations, without any rolling back.
- And I have done a key-value store that supports transaction (something similar with relational database) in [this project](https://github.com/leductoan3082004/in-memory-storage-engine). But this version only supports key/value store with value is an any type.
- It supports MVCC for a key (which may create and manage several versions).
- Applied repeatable read. Multithreading.
- Some basic operations that a transaction should have `BeginTransaction` `Abort` `Commit, ...`

#### Snapshot
- Although I just stored all the data on memory into the disk, and it's quick and easy to store.
- But if we develop more, we need to think about maybe our data structure will change a lot, and the last snapshot we stored may not fit with the new implementation.
- So we can develop a new strategy to store it. Maybe we will not store the whole memory into disk. Instead we will store all the users' events into one file (maybe several files to prevent loading for so long). With this event file, we can easily recover our snapshot from scratch right ? just go through from beginning of the file and to the end, apply all the events into memory. And bump, we have our snapshot recovered.
- But the main disadvantage of this is that it will take us so long to recover right ?
- To fix this, we may catch the snapshot at some points of time, and we can just apply our events from that point of time after. So it's much quicker.
- And in our code. We may handle Backward compatibility to make our application wont crash when apply the events.