TCP Fast Open
=============

Go example of TCP Fast Open (TFO) as described by RFC7413 and available in Linux Kernel 3.7 (server and client support).

Check support for TCP Fast Open by checking: `/proc/sys/net/ipv4/tcp_fastopen`, ensure this value is 3 for client and
server support. If necessary, echo 3 to this file, eg:

```bash
# echo 0 > /proc/sys/net/ipv4/tcp_fastopen
# cat /proc/sys/net/ipv4/tcp_fastopen
0
# echo 3 > /proc/sys/net/ipv4/tcp_fastopen
# cat /proc/sys/net/ipv4/tcp_fastopen
3
```

Using standard Linux Kernel system calls, this program shows steps required to configure a server and client to
establish a TCP connection using TFO.

The program simply establishes a connection to itself, using go routines, but doesn't demonstrate a complete client
server architecture with correct tear down.

Note: Go itself provides better handlers for establishing and listening to sockets via the
[net](https://golang.org/pkg/net/) package, but these packages don't currently provide TFO support.

Usage
=====

```bash
Usage of ./tcp-fast-open:
Options:
  -s 127.0.0.1  --server=127.0.0.1  Server to connect to (and listen if listening)
  -p 2222       --port=2222         Port to connect to (and listen to if listening)
  -l            --listen            Create a listening TFO socket
                --help              show usage message
```

Create a listening socket, then connect to it:
```bash
tcp-fast-open -l
```

Connect to remote server:

```bash
tcp-fast-open -s 192.0.2.1 -p 80
```

Once connected, you'll need to use `ip tcp_metrics` (check for `fo_cookie`) or `tcpdump` to determine whether a TFO connection was successful.


Behaviour
=========

TFO's goal is to establish a connection regardless of client, server or middleware support. The system calls provided by
the Kernel therefore abstract these details from the caller. For this reason this program itself can't determine whether:
- the connection was not established using cookies, but have been successfully transferred for consecutive connections
- the connection was establish using cookies
- the connection used an invalid cookie

Therefore, for this program, you'll need to use tcpdump to analyse the packets yourself: `tcpdump -s 0 -XX -nn -i lo port
2222`

When analysing the traces, pay particular attention to the following:

1. SYN packet should not contain data (like most standard SYN packets), but should contain to TFO option.
2. SYN-ACK response should contain a cookie to be cached by the client
3. ACK packet, as usual, contains data (and PUSH bit)

On the second, and consecutive connections, the following behaviour should occur:

1. SYN packet contains data, as well as TFO option and cookie
2. Usual SYN-ACK response
3. ACK packet contains no data
