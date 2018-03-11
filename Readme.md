# CuckooC - Cuckoo Cluster
[![Build Status](https://travis-ci.org/vedhavyas/cuckooc.svg?branch=master)](https://travis-ci.org/vedhavyas/cuckooc)
[![GitHub tag](https://img.shields.io/github/tag/vedhavyas/cuckooc.svg)](https://github.com/vedhavyas/cuckooc/tags)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![GitHub issues](https://img.shields.io/github/issues/vedhavyas/cuckooc.svg)](https://github.com/vedhavyas/cuckooc/issues)
![Contributions welcome](https://img.shields.io/badge/contributions-welcome-orange.svg)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/vedhavyas/cuckooc.svg)](https://github.com/vedhavyas/cuckooc/pulls)
[![Website](https://img.shields.io/website-up-down-green-red/http/vedhavyas.com.svg?label=my-website)](https://vedhavyas.com)

Cuckoo Cluster manages all of your [Cuckoo Filters](https://github.com/vedhavyas/cuckoo-filter).
Cuckoo filter, Practically better than Bloom Filter, support adding and removing items dynamically while achieving even higher performance than Bloom filters.
For applications that store many items and target moderately low false positive rates, cuckoo filters have lower space overhead than space-optimized Bloom filters.
It is a compact variant of a cuckoo hash table that stores only fingerprints —a bit string derived from the item using a hash function— for each item inserted, instead of key-value pairs.
The filter is densely filled with fingerprints (e.g., 95% entries occupied), which confers high space efficiency.

Paper can be found [here](https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf)

## Getting Started

Installation assumes that you have Go environment configured.

### Installing

Go get the project with following command

```
go get -u github.com/vedhavyas/cuckooc/cmd/cuckooc/...
```

## Running the tests

Once inside project' folder, simply run `make test` to run the tests.

## Configuration
Cuckoo Cluster requires a configuration file. Example configuration file can be found [here](testdata/config_example.json).

```json
{
  "debug": true,
  "backup_folder": "./testdata/backups",
  "tcp": ":4000",
  "udp": ":5000"
}
```

### Debug
Debug, if enabled, attaches error reason for failed queries.

### Backup Folder
Cuckoo Cluster will backup Filters inside this folder. Skips backup if empty

### TCP
If not empty, Cuckoo Cluster runs a TCP server at given address. Server will not close any connections and expects client
reuse the same connection

### UDP
If not empty, Cuckoo Cluster runs a UDP server at given address.  

## Commands

### Command format

Command format for CuckooC is as follows
```
[Filter-name] [action] [args...]
```

### Actions
#### new
```
>> test new
true
```
Creates a new filter, named test, with default count(4 << 20) and bucket size (8)

```
>> test new 100 8
true
```
Creates a new filter, named test, with count 100 and bucket size 8

#### set
```
>> test set x y z
true true true
```
Sets the values `x`, `y`, & `z` to filter named `test`

#### setu
```
>> test setu a b
true true
```
Sets the missing values `a`, `b`, & `c` to filter named `test`

#### check
```
>> test check x a 1 2
true true false
```
Checks if the given values are set in the filter `test`

#### delete
```
>> test delete x a 1
true true false
```
Deletes the values if present in the filter `test`

#### count
```
>> test count
3
```
Count returns the total items set in filter `test`

#### loadfactor
```
>> test loadfactor
0.03
```
loadfactor returns the load factor of the filter `test`

#### backup
```
>> test backup
true
```
Backup backups the filter to a persistent storage at the location provided in the config.
Fails if the path is not provided in the config

```
>> test backup /backup/some/path
true
``` 
Backups the filter at the provided argument path. Argument path take precedence over path provided in the configuration.

#### reload
```
>> test reload
true
```
Reload reloads the filter from last backup. Filter will be reloaded from backup path provided in the config 

```
>> test reload /backup/some/path
true
```
Reloads the filter from last backup at path given in argument. Argument path takes precedence over config path

#### stop
```
>> test stop
true
```
Stop backups(if provided in the config) and stops the filter `test`

**Backed up and stopped filter is reloaded back by cuckooC if the backup exists in the path provided in Config or by explicitly calling `reload` action**

### Multiple commands
You can send multiple commands in single request with `\n` as a delimiter
```
>>
    test new
    test set a b c
true
true true true
```

## Clients
CuckooC can be used with either TCP or UDP transport. Client should take care of re-using same TCP connection.
There are no CuckooC clients at the moment. But if you have built one, update the Readme and give me a PR.

## TODO
* Command Set
* and many more...

## Built With

* [Go](https://golang.org/)
* [Glide](https://glide.sh/) - Dependency Management
* [Cuckoo Filter](https://github.com/vedhavyas/cuckoo-filter)

## Contributing

PRs, Issues, and Feedback are very welcome and appreciated.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/vedhavyas/cuckoooc/tags).

## Authors

* **Vedhavyas Singareddi** - [Vedhavyas](https://github.com/vedhavyas)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
