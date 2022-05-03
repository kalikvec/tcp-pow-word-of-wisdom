# "Word of Wisdom" TCP-server with protection from DDOS based on Proof of Work

## 1. Problem statement
Design and implement “Word of Wisdom” tcp server.  
• TCP server should be protected from DDOS attacks with the [Proof of Work](https://en.wikipedia.org/wiki/Proof_of_work),, the challenge-response protocol should be used.  
• The choice of the POW algorithm should be explained.  
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.  
• Docker file should be provided both for the server and for the client that solves the POW challenge

## 2. How to run.

### 2.1 Requirements
+ [Go 1.18](https://go.dev/dl/) installed (if you want to run server or client without Docker)
+ [Docker](https://docs.docker.com/engine/install/) installed (to run docker-compose)

### 2.2 Start server and client by docker-compose:
```
make start
```

### 2.3 Start only server:
```
make start-server
```

### 2.4 Start only client:
```
make start-client
```

## 3. Alternatives and algorithm selection.

Hashcash's become the algorithm of choice as it has the following 
Pros:
+ ease of implementation and detailed description.
+ computational efficiency of validation on server side
+ dynamic complexity adjustment capabilities depending on client power, so clients may be forced to calculate
different leading zeros count.

Cons:
+ need to choose carefully challenge complexity depending on client computational power.

Alternatives:

[Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree)
+ more difficult to implement comparing with hashcash/
+ less efficient from the server side - server needs more power to check the challenge.

