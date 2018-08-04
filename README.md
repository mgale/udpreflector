# udpreflector

udpreflector will listen on a network port for UDP traffic and forward it on to
a remote server.

---


## How to use udpreflector

```
./udpreflector --listenport 0.0.0.0:8889 --destip 127.0.0.1 --destport 8888 -v
```