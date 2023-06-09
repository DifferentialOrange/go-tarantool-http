package = "httpgo-monitoring"
version = "scm-1"
source = {
    url    = 'git+https://github.com/DifferentialOrange/go-tarantool-http.git',
    branch = "master",
}

dependencies = {
    "lua >= 5.1",
    "httpgo",
    "metrics",
}

build = {
    type = 'cmake',
    variables = {
        TARANTOOL_INSTALL_LIBDIR = '$(LIBDIR)',
        TARANTOOL_INSTALL_LUADIR = '$(LUADIR)',
    },
}