box.cfg{listen='localhost:3301'}

local microservice = require('microservice')
local httpgo = require('httpgo')
local httpgo_monitoring = require('httpgo-monitoring')
local httpgo_auth = require('httpgo-auth')
local httpgo_crud = require('httpgo-crud')
-- local httpgo_lua = require('httpgo-lua')

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
        }
    }
})

local log = require('log')
local fiber = require('fiber')
local metrics = require('metrics')

fiber.create(function()
    while true do
        log.info(metrics.collect())
        fiber.sleep(5)
    end
end)
