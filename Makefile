all: .rocks

.rocks: httpgo/httpgo-scm-1.linux-x86_64.rock microservice/microservice-scm-1.all.rock
	tt rocks install microservice/microservice-scm-1.all.rock
	tt rocks install httpgo/httpgo-scm-1.linux-x86_64.rock

microservice/microservice-scm-1.all.rock: microservice/microservice-scm-1.rockspec microservice/init.lua
	cd ./microservice && tt rocks make && tt rocks pack microservice

httpgo/httpgo-scm-1.linux-x86_64.rock: httpgo/httpgo-scm-1.rockspec httpgo/httpgo.go httpgo/init.lua.in
	cd ./httpgo && tt rocks make && tt rocks pack httpgo
