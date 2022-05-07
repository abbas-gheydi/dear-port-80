
# Dear Port80

## About The Project:

“Dear Port80” is a zero-config TCP proxy server that hides a SSH connection behind a HTTP server!

```

+--------------------------+      +--------------+                +------------+
|      CLIENT REQUEST      |      | Proxy server |                | web server |
|curl http://10.10.10.1:80 | -->  |10.10.10.1:80 | --> (HTTP) --> |            |
|           or             |      |              |                +------------+
|ssh 10.10.10.1 -p 80      |      +-----+--------+
+--------------------------+            |
                                        |
                                        |
                                        |
                                        |                         +------------+
                                        +------> ( SSH ) -------> | ssh server |
                                                                  +------------+



```
It supports two kinds of upstream servers, the first one is a web server like nginx and the other one is a SSH server. It listens on port 8080 ( by default ) and it serves these two protocols on port 8080 at the same time!   
 It sends **all packets** to the HTTP backend server but if it detects that the request is from a SSH client then it proxies traffic to the SSH server.


## How to Use It:

[Download DearPort80](https://github.com/Abbas-gheydi/dear-port-80/releases) and configur it using command line arguments.
```bash
./dearport80 --help
Usage of dearport80:
  -listen string
    	listen Address (default "0.0.0.0:8080")
  -enable_ssh
    	enable ssh proxy (default true)
  -ssh string
    	SSH upstream server address (default "127.0.0.1:22")
  -http string
    	HTTP upstream server address (default "127.0.0.1:80")


```

For example:
```bash

./dearport80 -listen="0.0.0.0:80" -http="10.10.10.1:80" -ssh="127.0.0.1:22"
```

It listens to port 80 and proxies HTTP traffic to the “10.10.10.1:80” and SSH to the port 22 of the local host.
   
      
      
To run it as a service you can use this guide:

https://www.suse.com/support/kb/doc/?id=000019672

  
  
  
  

  

## License

MIT
