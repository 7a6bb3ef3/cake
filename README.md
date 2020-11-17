### CAKE

PROXY Project which use a custom time-based encryption protocol.

> Experiment only ,low optimization and incomplete features(DNS proxy...) .v2fly/core + Clash is the recommended way.

#### Usage

> Anyway ,If one day they use the whitelist, we are done.

```
[root@ahoy ~]# chmod 755 build.sh && ./build.sh
[root@ahoy ~]# ./cake
2020-08-27 23:39:49     INFO    Listen on 127.0.0.1:1921

[root@ahoy ~]# ./cakecli -proxy=127.0.0.1:1921 -cryptor=chacha
2020-08-27 23:41:22     INFO    Use cryptor chacha
2020-08-27 23:41:22     INFO    Find apnic-latest.txt
2020-08-27 23:41:22     INFO    Socks5 listen on 127.0.0.1:1920
2020-08-27 23:41:22     INFO    HTTP listen on 127.0.0.1:1919
```
