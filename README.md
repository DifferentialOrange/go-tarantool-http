## Build
```bash
cd ./microservice && tt rocks make && tt rocks pack microservice; cd ..
cd ./httpgo && tt rocks make && tt rocks pack httpgo; cd ..
tt rocks install ./microservice/microservice-scm-1.all.rock ./httpgo/httpgo-scm-1.linux-x86_64.rock
tt run init.lua
```

## Performance test
```bash
k6 run load.js --duration 30s --vus 100
```
