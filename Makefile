
export GOPATH=$(shell pwd)
export GOBIN=$(shell pwd)/bin
export HOSTNAME=$(shell hostname)

all:
	mkdir -p $(GOBIN)
	go get github.com/stretchr/testify/assert
	go get
	go install

clean:
	rm -rf bin

test:
	go test *.go

test.verbose:
	go test *.go -v

run.testdb:
	-killall mongod
	sleep 4
	rm -rf testdb
	mkdir -p testdb/one
	mkdir -p testdb/two
	mkdir -p testdb/three
	mongod --replSet atoll-test --logpath testdb/mongod.log --fork --dbpath testdb/one --port 26017
	mongod --replSet atoll-test --logpath testdb/mongod.log --fork --dbpath testdb/two --port 26016
	mongod --replSet atoll-test --logpath testdb/mongod.log --fork --dbpath testdb/three --port 26015
	echo "rs.initiate()" | mongo --port 26017
	echo "rs.add(\""$(HOSTNAME)":26016\")" | mongo --port 26017
	echo "rs.addArb(\""$(HOSTNAME)":26015\")" | mongo --port 26017
	
