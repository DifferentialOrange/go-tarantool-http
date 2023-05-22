local ddl = require('ddl')

local function init(opts) -- luacheck: no unused args
    if opts.is_master then
        local ok, err = ddl.set_schema({
            spaces = {
                customers = {
                    engine = 'memtx',
                    is_local = false,
                    temporary = false,
                    format = {
                        {name = 'id', is_nullable = false, type = 'unsigned'},
                        {name = 'bucket_id', is_nullable = false, type = 'unsigned'},
                        {name = 'name', is_nullable = false, type = 'string'},
                    },
                    indexes = {
                        {
                            name = 'id',
                            type = 'TREE',
                            unique = true,
                            parts = {
                                {path = 'id', is_nullable = false, type = 'unsigned'},
                            },
                        },
                        {
                            name = 'bucket_id',
                            type = 'TREE',
                            unique = false,
                            parts = {
                                {path = 'bucket_id', is_nullable = false, type = 'unsigned'},
                            },
                        },
                    },
                    sharding_key = {'id'},
                },
            },
        })

        if not ok then
            error(err)
        end
    end

    return true
end

local function stop()
    return true
end

local function validate_config(conf_new, conf_old) -- luacheck: no unused args
    return true
end

local function apply_config(conf, opts) -- luacheck: no unused args
    -- if opts.is_master then
    -- end

    return true
end

return {
    init = init,
    stop = stop,
    validate_config = validate_config,
    apply_config = apply_config,
    dependencies = {'cartridge.roles.crud-storage'},
}
