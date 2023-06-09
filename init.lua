box.cfg{listen='localhost:3301'}

-- local log = require('log')
-- local fiber = require('fiber')
local json = require('json')
-- local metrics = require('metrics')
local json_exporter = require('metrics.plugins.json')

local microservice = require('microservice')
local httpgo = require('httpgo')
local httpgo_monitoring = require('httpgo-monitoring')
local httpgo_auth = require('httpgo-auth')
local httpgo_crud = require('httpgo-crud')
local httpgo_lua = require('httpgo-lua')

rawset(_G, 'my_procedure', function(_)
    return "Select everything"
end)

microservice.run(httpgo, {
    listen = 'localhost:8081',
    pipeline = {
        {
            plugin = httpgo_monitoring.plugin,
            cfg = {
                collector_name = 'tnt_httpgo_latency',
            }
        },
        {
            plugin = httpgo_auth.plugin,
            cfg = {
                token = '9June2023',
            }
        },
        {
            plugin = httpgo_crud.plugin,
            cfg = {
               routes = {
                    {
                        path = '/hello',
                        method = 'GET',
                        handler = httpgo_crud.handlers.Hello,
                    }
                }
            }
        },
        {
            plugin = httpgo_lua.plugin,
            cfg = {
               routes = {
                    {
                        path = '/metrics',
                        method = 'GET',
                        handler = httpgo_lua.inplace_handler(json_exporter.export),
                    },
                    {
                        path = '/sum',
                        method = 'GET',
                        handler = httpgo_lua.inplace_handler(function(body)
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
                        end),
                    },
                    {
                        path = '/my_procedure',
                        method = 'GET',
                        handler = 'my_procedure',
                    },
                }
            }
        }
    }
})
