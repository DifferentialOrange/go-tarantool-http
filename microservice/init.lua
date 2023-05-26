local json = require('json')
local popen = require('popen')

local charset = {} -- [0-9a-zA-Z]
for c = 48, 57  do table.insert(charset, string.char(c)) end
for c = 65, 90  do table.insert(charset, string.char(c)) end
for c = 97, 122 do table.insert(charset, string.char(c)) end

local function random_string(length)
    if not length or length <= 0 then return '' end
    math.randomseed(os.clock()^5)
    return random_string(length - 1) .. charset[math.random(1, #charset)]
end

local user = 'microservice_user'
local password = random_string(15)

box.schema.user.create(user, {if_not_exists = true})
box.schema.user.passwd(user, password)
box.schema.user.grant(user,'read,write,execute,create,drop', 'universe', nil, {if_not_exists = true})

local function run(module, cfg)
    local p, err = popen.new({module.binary_path},
        {env = {
            SERVER_LISTEN = box.info.listen,
            SERVER_USER = user,
            SERVER_PASS = password,
            TT_MICROSERVICE_CFG = json.encode(module.config_prepare(cfg)),
        }}
    )

    if err then
        error(err)
    end

    return p
end

return {
	run = run,
}
