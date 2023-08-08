# czskm-miniweb
Official Czechoslovak Minimarathon website / RTMP authentication

## Usage

The built file shows you help when running without parameters. Recommended to run with pm2

```sh
pm2 start run --name "RTMP Authentication" -- start
```

Your nginx should have rtmp setup similarly to this for the authentication to work:

```
rtmp {
        server {
                listen 1935;
                chunk_size 4096;
                ping 20s;
                ping_timeout 10s;
                notify_method get;

                application live {
                        live on;
                        on_publish http://localhost:8080/auth;
                        on_publish_done http://localhost:8080/disconnect;
                        record off;
                }
        }
}
```
