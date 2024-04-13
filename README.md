# GorillaMQ

![go-version](https://img.shields.io/badge/Golang-1.21-66ADD8?style=for-the-badge&logo=go)
![app-version](https://img.shields.io/badge/Version-0.1.0-red?style=for-the-badge&logo=github)

The cloud and edge native messaging broker server written in Go

A fast message broker implemented with Golang programming language. You can use GorillaMQ in order to make communication between clients with sending and receiving events.

## TODOs

- implement connection-pool for clients
- improve concurrency (implement worker-pool in broker module)
- secure and insecure connection at the same time (gorillamqs protocol)
- implement basic auth for gorillamqs
- think about deployment and horizontal scaling
- Durability and persistans challanges
