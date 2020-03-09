# tunnel-helper

a small tool to open ssh tunnel.

## Usage

### how to install

```sh
go install github.com/yeqown/infrastructure/cmd/tunnel-helper
```

### how to use cli

```sh
➜  tunnel-helper git:(master) ✗ tunnel-helper -h                     
Usage of tunnel-helper:
  -c string
        specified config file to load, default is ./.tunnelrc in json format
  -i    create an default format config in current folder
  -p string
        pattern to match with specified tunnel to open (default ".*")
```

## Config file format
tunnelrc file format:
```json
{
    "ssh": {
        "host": "host",
        "user": "username",
        "secret": "password",
        "privateKeyFile": "path/to/.ssh/id_rsa",
        "port": 22
    }, // default ssh config
    "tunnels": [
        {
            "ident": "mongo",
            "ssh": {
                "host": "host",
                "user": "username",
                "secret": "password",
                "privateKeyFile": "path/to/.ssh/id_rsa",
                "port": 22
            },
            "localPort": 27017,             // local port
            "remoteHost": "192.168.3.34",   // remote server host
            "remotePort": 27017             // remote server port
        },
         {
            "ident": "redis",
            "ssh": null, // the global ssh config will be used
            "localPort": 6379,
            "remoteHost": "192.168.3.58",
            "remotePort": 6370
        }
    ]
}
```

## Output

```log
➜  tunnel-helper git:(master) ✗ tunnel-helper -c ./.tunnelrc -p redis
2020/03/09 13:21:53 main.go:183: [WARN] tunnel ident=mongo, not matched with pattern=redis, so skipped
2020/03/09 13:21:53 main.go:230: [INFO] 1 tunnel starting, current: 1
2020/03/09 13:22:14 ssh_tunnel.go:129: [INFO] tunnel=(localhost:6379) accepted connection
2020/03/09 13:22:14 ssh_tunnel.go:141: [INFO] tunnel=(localhost:6379) connected to server=111.231.85.95:22 (1 of 2)
2020/03/09 13:22:15 ssh_tunnel.go:147: [INFO] tunnel=(localhost:6379) connected to remote=192.168.3.58:6370 (2 of 2)
```