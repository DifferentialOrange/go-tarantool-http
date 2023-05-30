box.cfg{listen='localhost:3301'}

local json = require('json')

local microservice = require('microservice')
local httpgo = require('httpgo')

microservice.run(httpgo, {
    listen = 'localhost:8081',
    routes = {
        ['/hello'] = function(body)
            local args = json.decode(body)

            if type(args) ~= 'table' then
                error("Unsupported type " .. type(args))
            end

            local _, val = next(args)
            return "Hello world from " .. val
        end,
        ['/sum'] = function(body)
            local args = json.decode(body)
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
            return sum
        end,
    }
})
