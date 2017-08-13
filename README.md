# Simplechain

### Explanation
A small scale blockchain designed as simple as possible. This bare bones implementation of a distributed ledger offers most of the basic functionality and validation rules that are seen in large scale blockchain applications.

### What is a blockchain anyway?
[From Wikipedia](https://en.wikipedia.org/wiki/Blockchain): A blockchain is a distributed database that is used to maintain a continuously growing list of records, called blocks, that are protected from tampering and open to the public.

### What can Simplechain do?
Blocks contained on the Simplechain can only hold one piece of data, a string. I created this application to focus more on the blockchain data structure itself rather than the blocks and the data held within those. This blockchain design uses a simplified Proof of Work system to create new blocks, inspired by ideas from [this article](https://en.bitcoin.it/wiki/Proof_of_work).

### Usage (written in Go 1.8.3)
To compile:
```
$ go build simplechain.go
```
###
To execute:
```
$ ./simplechain
```