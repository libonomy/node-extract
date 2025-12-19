<br />
<p align="center">
  <h3 align="center">Libonomy AI Node Data Extraction</h3>

Extracts data from blockchain nodes in real time for AI training, classification, and simulation pipelines.

<b>Halt Status:</b> POC and Simulators were complete. In-progress for R&D work related to mainnet universal usage.


</p>

## Getting Started

Follow the instructions below to run the project locally.

## Prerequisites

Make sure you have the following installed:

- **GO**

  Install GO from:

  ```sh
  https://go.dev/doc/install
  ```

  ⚠️ **Supported GO version:** `1.14.x`


## Installation

1. **Clone the repository**

2. **Navigate to the project directory**

3. **Install dependencies**

```sh
go mod tidy
```

## To run the server

```sh
go run main.go
```

This repo contains the implementation of the data extraction module of libonomy daemon. The implementation covers the specifications that enables the daemon to extract the data from the node and feed it into the classification agent.


```
Copyright © 2025 — Libonomy
```