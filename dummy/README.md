# Dummy

Dummy is a simple replica implementation that can be used to test the torture testing framework

Dummy replica is a distributed application where a set of Dummy replicas sent and receive random messages

Periodically Dummy replica prints the throughput, average latency and the resource usage.

```build.sh``` in the main directory contains the instructions to run the dummy replica

To run the web front ends

```http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=10100```

```http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=20100```

```http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=30100```

```http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=40100```

```http://localhost:63342/toture-testing-consensus/dummy/run/index.html/?port=50100```