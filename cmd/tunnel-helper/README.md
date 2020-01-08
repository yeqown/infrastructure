# tunnel-helper

a small tool to open ssh tunnel.

```sh
tunnel-helper -c path/to/tunnerlrc
```

tunnelrc file format:
```json
{
    "ssh": {
        "host": "host",
        "user": "username",
        "secret": "password",
        "privateKeyFile": "path/to/privateKey.perm",
        "port": 22
    }, // default ssh config
    "tunnels": [
        {
            "ssh": null,                    // if current ssh is null, the default ssh config will be used
            "localPort": 27017,             // local port
            "remoteHost": "192.168.3.34",   // remote server host
            "remotePort": 27017             // remote server port
        },
         {
            "ssh": null,        
            "localPort": 6379,
            "remoteHost": "192.168.3.58",
            "remotePort": 6370
        }
    ]
}
```

## example

```log
➜  tunnel-helper git:(master) ✗ ./app
2020/01/08 14:29:16 main.go:103: [INFO] 1 tunnel starting, current: 1
2020/01/08 14:29:16 main.go:103: [INFO] 1 tunnel starting, current: 2
2020/01/08 14:29:21 ssh_tunnel.go:113: [INFO] tunnel=(localhost:6379) accepted connection
2020/01/08 14:29:21 ssh_tunnel.go:125: [INFO] tunnel=(localhost:6379) connected to server=111.231.85.95:22 (1 of 2)
2020/01/08 14:29:21 ssh_tunnel.go:131: [INFO] tunnel=(localhost:6379) connected to remote=192.168.3.58:6370 (2 of 2)
```