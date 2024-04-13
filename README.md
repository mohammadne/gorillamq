# GorillaMQ

![go-version](https://img.shields.io/badge/Golang-1.21-66ADD8?style=for-the-badge&logo=go)
![build-status](https://img.shields.io/github/actions/workflow/status/gorillamq/gorillamq/test.yaml?logo=github&style=for-the-badge)
![app-version](https://img.shields.io/github/v/tag/gorillamq/gorillamq?sort=semver&style=for-the-badge&logo=github)
![coverage](https://img.shields.io/codecov/c/github/gorillamq/gorillamq?logo=codecov&style=for-the-badge)

The cloud and edge native messaging broker server written in Go

A fast message broker implemented with Golang programming language. You can use GorillaMQ in order to make communication between clients with sending and receiving events.

## TODOs

- implement connection-pool for clients
- improve concurrency (implement worker-pool in broker module)
- secure and insecure connection at the same time (gorillamqs protocol)
- implement basic auth for gorillamqs
- think about deployment and horizontal scaling
- Durability and persistans challanges
