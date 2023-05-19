```bash
go build -tags go_tarantool_ssl_disable
```

```bash
rm -rf *.snap *.xlog
tarantool init.lua
```
