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