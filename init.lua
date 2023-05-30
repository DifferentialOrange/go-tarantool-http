box.cfg{listen='localhost:3301'}

local microservice = require('microservice')
local httpgo = require('httpgo')

microservice.run(httpgo, {
    listen = 'localhost:8081',
    routes = {
        ['/hello'] = function(body) return "Hello world!" .. tostring(body) end,
    }
})
