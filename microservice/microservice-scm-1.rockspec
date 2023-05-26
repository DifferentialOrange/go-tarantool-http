package = "microservice"
version = "scm-1"
source = {
    url    = 'git+https://github.com/DifferentialOrange/go-tarantool-http.git',
    branch = "master",
}

dependencies = {
    "lua >= 5.1"
}

build = {
    type = 'cmake',
    variables = {
        TARANTOOL_INSTALL_LUADIR = '$(LUADIR)',
    },
}