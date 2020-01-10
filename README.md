# alien-invasion
Alien Invasion Simulation

## Build & Run

To run the `alien-invasion` have a [Golang](https://golang.org/doc/install) installed. If everything is ready, just run the following:

```
$ go run main.go
```
This will run the simulation using all defaults and current Unix time as a random source.

To list all `cli` options ask for help:
```
$ go run main.go -help
Usage of /main:
  -aliens int
        number of aliens invading (default 10)
  -iterations int
        number of iterations (default 10000)
  -world string
        a file used as world map input (default "./data/map1.txt")
```

You can run the specific simulation by below command:

```
$ go run main.go -aliens 4 -iterations 100
```
