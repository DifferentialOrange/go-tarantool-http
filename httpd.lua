box.cfg{listen='localhost:3301'}

local http_server = require('http.server')

local httpd = http_server.new('localhost', 8081)

httpd:route({path = '/hello'}, function(req)
    local args = req:json()

    if type(args) ~= 'table' then
        error("Unsupported type " .. type(args))
    end

    local _, val = next(args)
    return {
        status = 200,
        body = "Hello world from " .. val,
    }
end)
httpd:route({path = '/sum'}, function(req)
    local args = req:json()
    local sum = 0

    if type(args) ~= 'table' then
        error("Unsupported type " .. type(args))
    end

    for _, v in ipairs(args) do
        if type(v) == 'number' then
            sum = sum + v
        else
            error("Unsupported value type " .. type(v))
        end
    end
    return {
        status = 200,
        body = sum,
    }
end)
httpd:start()
