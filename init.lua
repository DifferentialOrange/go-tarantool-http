box.cfg{listen='localhost:3301'}

local microservice = require('microservice')
local httpgo = require('httpgo')
-- local httpgo_auth = require('httpgo-auth')
local httpgo_crud = require('httpgo-crud')
-- local httpgo_lua = require('httpgo-lua')

microservice.run(httpgo, {
    listen = 'localhost:8081',
    pipeline = {
        -- {
        --     plugin = httpgo_auth.plugin,
        --     cfg = {
        --         token = '9June2023',
        --     }
        -- },
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
