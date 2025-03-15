# Bare Metal TCP Server

TCP chat server, implemented in Go using only the standard library, that implements a custom binary protocol.

This was developed during the [Coderful 2025 Workshop](https://edizione1.coderful.io/workshop/beyond-frameworks-anatomia-comunicazione-client-server-gabriele-santomaggio/) led by [Gabriele Santomaggio](https://github.com/Gsantomaggio)

The protocol specifications can be found in the [README of the original repo](https://github.com/Gsantomaggio/chat/blob/4f9ed460c1fc04d94e22223cbb655de3d05d1378/README.md), and also appended to this README, in case the original repo gets deleted.

## Project Overview

The goal of this project is to create a simple client-server chat application using a custom binary protocol. The server is responsible for:

- Handling client connections
- Processing login and message commands
- Managing user sessions and statuses
- Storing and dispatching messages


## Running the Server

To run the server, navigate to the server directory and run it:

```sh
go run .
```

The server will start listening on port 5555.


# Original README

## Console Chat

### Disclaimer:

This project's scope is only to learn how to write a client-server application from scratch.
The code is demonstrative and focuses more on understanding how the TCP works than production code.
The client applications are meant to test the server quickly.

We provide examples in the following languages:

- Server [golang](./server/go) [Full features]
- Client [golang](./server/go/tcp_client) [Full features]
- Server [rust](./server/rust) [Partial features]
- Server [NodeJS](./server/nodejs) [Partial features]
- Server [Python](./server/python) [Full features]

### How to run the golang stack:

- Server: `/server/go/run/server` and run `go run main.go localhost:5555`
- Client: `/server/go/run/client` and run `go run main.go localhost:5555` ( yes, the client is inside the `server` directory because they share the same codec)
- You can use two different terminals with two different users

### Protocol definition:

The protocol is a binary protocol with the following structure:

### Data types

- `byte` 1 byte
- `uint16` 2 bytes
- `uint32` 4 bytes
- `string` 2 bytes for the length + N bytes for the string
- `[]byte` 2 bytes for the length + N bytes for the string
- `uint64` 8 bytes

### Header

| Name      | Type     |
| --------- | -------- |
| `version` | `byte`   |
| `command` | `uint16` |

### CommandLogin

| Name            | Type     | value(s) | reference         |
| --------------- | -------- | -------- | ----------------- |
| `version`       | `byte`   | 0x01     | `Header::version` |
| `key`           | `uint16` | 0x01     | `Header::command` |
| `correlationId` | `uint32` |          |                   |
| `username`      | `string` |          |                   |

### CommandMessage

| Name            | Type     | value(s) | reference         |
| --------------- | -------- | -------- | ----------------- |
| `version`       | `byte`   | 0x01     | `Header::version` |
| `key`           | `uint16` | 0x02     | `Header::command` |
| `correlationId` | `uint32` |          |                   |
| `message`       | `string` |          |                   |
| `From`          | `string` |          |                   |
| `To`            | `string` |          |                   |
| `Time`          | `uint64` |          |                   |

## Response

All the commands will have a response with the following structure:

| Name            | Type     | value(s) | reference         |
| --------------- | -------- | -------- | ----------------- |
| `version`       | `byte`   | 0x01     | `Header::version` |
| `key`           | `uint16` | 0x03     | `Header::command` |
| `correlationId` | `uint32` |          |                   |
| `code`          | `uint16` |          | `ResponseCodes`   |

### ResponseCodes

| Name                     | value(s) |
| ------------------------ | -------- |
| `OK`                     | 0x01     |
| `ErrorUserNotFound`      | 0x03     |
| `ErrorUserAlreadyLogged` | 0x04     |

## Data (bytes) written on the socket

1. Write the length of whole message (header + command) as a `uint32`
1. Write the header
1. Write the command

### Example with `CommandLogin`

If `user1` wants to login, this will be the structure of data sent.

- `Header`:
  - `version` = 0x01 => 1 byte
  - `command` = 0x01 => 2 bytes
- `CommandLogin`:
  - `correlationId` = 0x01 => 4 bytes
  - `username` = "user1" => length = 2 + 5 = 7

In this case the client will:

- write the length in bytes (as `uint32`) of the whole message: 3 + 11 = 14
  - note: this `uint32` is excluded from the total bytes count
- write the `header` + `message`:

```
0x00 0x00 0x00 0x0E (`uint32`)
0x01 =>  version  (1 byte)
0x00 0x01  => command (`uint16`)
0x00 0x00 0x00 0x01  => correlationId (`uint32`)
0x00 0x05 => username  length  (`uint16`)
0x75 0x73 0x65 0x72 0x31 => username (user1) (5 bytes)
```
- Total bytes written: 14 (body)  + 4 (len of body) = 18
- Send the message
- Read the `Response`

### Read data from the socket

- Read the length of the whole message (`header` + `command`) as a `uint32`
- Read the `header`
- Read the `command key`
- Read the `command` based on the `key`

For Example: `CommandLogin` with username `user1`

- Read the first 4 bytes for the length of the whole message: 14
- Ensure the socket buffer has at least 14 bytes
- Read the header:
  - Read the `version`: 0x01
  - Read the `command`: 0x01
  - Read the `command` based on the `key`: `CommandLogin`
    - Read the `correlationId`: 0x01
    - Read the username: "user1"
  - Process the command
  - Send the `Response`

### CorrelationId

The `correlationId` is a unique identifier for each command sent by the client.
The server will send back the `correlationId` to the client.
So the client can match the response with the command sent.
You must be sure to provide a unique `correlationId` for each command sent.

### Server Side Mandatory Features

- [x] Login (without password. It is enough to send the username)
- [x] Send message and dispatch to the correct user
- [x] Store in memory the users with the status (online/offline)
- [x] Store in memory the off-line messages when the user is not online
- [x] Send the off-line messages when the user logs in
- [x] Check if the user is already logged in
- [x] Check if the destination user exists

### Server Side Nice to have Features

- [ ] Logout
- [ ] Send message to multiple users
- [ ] Send message to all users
- [ ] Command to get the list of users
- [ ] Persist the users and messages in a database