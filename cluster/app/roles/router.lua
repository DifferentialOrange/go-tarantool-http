local popen = require('popen')

local process = nil

local function init(opts) -- luacheck: no unused args
    local p, err = popen.new({'go-tarantool-http'}, {
        env = {
            LISTEN_ADDR = box.info.listen,
            SERVER_USER = 'admin',
            SERVER_PASS = 'secret-cluster-cookie',
        }
    })

    assert(p ~= nil, err)
    assert(err == nil, err)
    process = p

    return true
end

local function stop()
    return true
end

local function validate_config(conf_new, conf_old) -- luacheck: no unused args
    return true
end

local function apply_config(conf, opts) -- luacheck: no unused args
    -- if opts.is_master then
    -- end

    return true
end

return {
    init = init,
    stop = stop,
    validate_config = validate_config,
    apply_config = apply_config,
    dependencies = {'cartridge.roles.crud-router'},
    process = process,
}
