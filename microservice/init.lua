local json = require('json')
local fiber = require('fiber')
local popen = require('popen')
local log = require('log')

local function create(module, cfg)
    local full_cfg = module.build_config(cfg)

    return {
        run = function()
            fiber.create(function()
                fiber.name('microservice')

                while true do
                    local p, err = popen.new({module.binary_path},
                        {env = {
                            TT_MICROSERVICE_CFG = json.encode(full_cfg),
                        }}
                    )

                    if err then
                        error(err)
                    end

                    while p.pid ~= nil do
                        p:wait()
                    end

                    log.info("process is no longer alive, info: %s", json.encode(p:info()))
                    p:close()
                end
            end)

        end
    }
end

return {
	create = create,
}
