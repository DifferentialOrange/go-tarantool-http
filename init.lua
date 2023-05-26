box.cfg{listen='localhost:3301'}

local microservice = require('microservice')
local httpgo = require('httpgo')

microservice.run(httpgo, {
    listen = 'localhost:8081'
})
