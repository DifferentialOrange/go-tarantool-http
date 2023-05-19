local fiber = require('fiber')
local popen = require('popen')

local listen = 'localhost:3301'

box.cfg{listen = listen}

local function get_server_config()
	return {somedata = 'data'}
end

rawset(_G, 'get_server_config', get_server_config)

box.schema.user.grant('guest', 'read,write,execute,create,alter,drop,usage,session', 'universe')

local p = popen.new({'go-tarantool-http'}, {
    env = {
        LISTEN_ADDR = listen,
        SERVER_USER = 'guest',
        SERVER_PASS = '',
    }
})

fiber.sleep(10)
