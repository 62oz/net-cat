# net-cat

This project is a recreation of the NetCat command-line utility in a server-client architecture that can establish a TCP connection between a server and multiple clients for a group chat. The project is written in Go and has several features such as a name requirement for the client, the ability to control connections quantity, and the ability to send messages to the chat.

To run the project, use the following command:

go

go run .

You can specify the port number as follows:

go

go run . <port_number>

If there is no port specified, then the default port is set to 8989.

When a client connects, they will be greeted with the Linux logo and prompted to enter their name. The client's name must not be empty to be accepted. The project supports up to 10 connections.

The messages sent by clients to the chat are identified by the time it was sent and the client's name in the following format: [time][client.name]:[client.message]. If a client joins the chat, all previous messages sent to the chat are uploaded to the new client, and the rest of the clients are informed of the new client joining. If a client leaves the chat, the rest of the clients are informed that the client left, but they must not disconnect.

The following packages are allowed for the project: io, log, os, fmt, net, sync, time, bufio, errors, strings, and reflect. It is recommended to have test files for unit testing both the server connection and the client.

Please note that the code must respect good practices, use Go-routines and channels or Mutexes to handle multiple connections, and handle errors from both server-side and client-side.

Have fun using this project to create your own chat group!
