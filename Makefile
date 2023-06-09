all: .rocks

.rocks: microservice/microservice-scm-1.all.rock \
	httpgo/httpgo-scm-1.linux-x86_64.rock \
	httpgo-crud/httpgo-crud-scm-1.linux-x86_64.rock \
	httpgo-auth/httpgo-auth-scm-1.linux-x86_64.rock \
	httpgo-monitoring/httpgo-monitoring-scm-1.linux-x86_64.rock \
	httpgo-lua/httpgo-lua-scm-1.linux-x86_64.rock
	tt rocks install microservice/microservice-scm-1.all.rock
	tt rocks install httpgo/httpgo-scm-1.linux-x86_64.rock
	tt rocks install httpgo-crud/httpgo-crud-scm-1.linux-x86_64.rock
	tt rocks install httpgo-auth/httpgo-auth-scm-1.linux-x86_64.rock
	tt rocks install httpgo-monitoring/httpgo-monitoring-scm-1.linux-x86_64.rock
	tt rocks install httpgo-lua/httpgo-lua-scm-1.linux-x86_64.rock

microservice/microservice-scm-1.all.rock: microservice/microservice-scm-1.rockspec microservice/CMakeLists.txt microservice/init.lua
	cd ./microservice && tt rocks make && tt rocks pack microservice

httpgo/httpgo-scm-1.linux-x86_64.rock: httpgo/httpgo-scm-1.rockspec httpgo/CMakeLists.txt httpgo/httpgo.go httpgo/init.lua.in
	cd ./httpgo && tt rocks make && tt rocks pack httpgo

httpgo-crud/httpgo-crud-scm-1.linux-x86_64.rock: httpgo-crud/httpgo-crud-scm-1.rockspec \
												 httpgo-crud/CMakeLists.txt httpgo-crud/httpgo-crud.go \
												 httpgo-crud/init.lua.in \
												 httpgo/httpgo-scm-1.linux-x86_64.rock
	cd ./httpgo-crud && tt rocks install ../httpgo/httpgo-scm-1.linux-x86_64.rock && tt rocks make && tt rocks pack httpgo-crud

httpgo-auth/httpgo-auth-scm-1.linux-x86_64.rock: httpgo-auth/httpgo-auth-scm-1.rockspec \
												 httpgo-auth/CMakeLists.txt httpgo-auth/httpgo-auth.go \
												 httpgo-auth/init.lua.in \
												 httpgo/httpgo-scm-1.linux-x86_64.rock
	cd ./httpgo-auth && tt rocks install ../httpgo/httpgo-scm-1.linux-x86_64.rock && tt rocks make && tt rocks pack httpgo-auth

httpgo-monitoring/httpgo-monitoring-scm-1.linux-x86_64.rock: httpgo-monitoring/httpgo-monitoring-scm-1.rockspec \
												 httpgo-monitoring/CMakeLists.txt httpgo-monitoring/httpgo-monitoring.go \
												 httpgo-monitoring/init.lua.in \
												 httpgo/httpgo-scm-1.linux-x86_64.rock
	cd ./httpgo-monitoring && tt rocks install ../httpgo/httpgo-scm-1.linux-x86_64.rock && tt rocks make && tt rocks pack httpgo-monitoring

httpgo-lua/httpgo-lua-scm-1.linux-x86_64.rock: httpgo-lua/httpgo-lua-scm-1.rockspec \
												 httpgo-lua/CMakeLists.txt httpgo-lua/httpgo-lua.go \
												 httpgo-lua/init.lua.in \
												 httpgo/httpgo-scm-1.linux-x86_64.rock
	cd ./httpgo-lua && tt rocks install ../httpgo/httpgo-scm-1.linux-x86_64.rock && tt rocks make && tt rocks pack httpgo-lua
