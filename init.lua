box.cfg{listen='localhost:3301'}

local microservice = require('microservice')
local httpgo = require('httpgo')
local httpgo_crud = require('httpgo-crud')

microservice.run(httpgo, {
    listen = 'localhost:8081',
    routes =  {
        {
            path = '/hello',
            handler = {
                type = 'plugin',
                lib = httpgo_crud.lib,
                func = httpgo_crud.Hello,
            },
        },
    }
})
