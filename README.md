## Build
```bash
make .rocks
tt run init.lua
```

## Performance test
```bash
k6 run load.js --duration 30s --vus 100
```
