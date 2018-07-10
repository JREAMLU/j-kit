package redis

// HSCAN hash scan
const HSCAN = `local h_key = ARGV[1]
local h_cursor = ARGV[2]
local h_count = ARGV[3]
local results = redis.call('HSCAN',h_key,h_cursor,'COUNT',h_count)
local list = {}
if #results > 0 then
    table.insert(list,results[1])

    if #results[2] > 0 then
        for i = 1,#results[2] do
            table.insert(list,results[2][i])
        end
    end
end
return list`

// SSCAN set scan
const SSCAN = `local s_key = ARGV[1]
local s_cursor = ARGV[2]
local s_count = ARGV[3]
local results = redis.call('SSCAN',s_key,s_cursor,'COUNT',s_count)
local list = {}
if #results > 0 then
    table.insert(list,results[1])

    if #results[2] > 0 then
        for i = 1,#results[2] do
            table.insert(list,results[2][i])
        end
    end
end
return list`

// HKEYSCAN hash key scan
const HKEYSCAN = `local h_key = ARGV[1]
local h_cursor = ARGV[2]
local h_count = ARGV[3]
local results = redis.call('HSCAN',h_key,h_cursor,'COUNT',h_count)
local list = {}
if #results > 0 then
    table.insert(list,results[1])

    if #results[2] > 0 then
        for i = 1,#results[2],2 do
            table.insert(list,results[2][i])
        end
    end
end
return list`

// SCAN sacn
const SCAN = `local pattern = ARGV[1]
local cursor = ARGV[2]
local count = ARGV[3]
local results = redis.call('SCAN',cursor,'MATCH',pattern,'COUNT',count)
local list = {}
if #results > 0 then
    table.insert(list,results[1])

    if #results[2] > 0 then
        for i = 1,#results[2] do
            table.insert(list,results[2][i])
        end
    end
end
return list`

// ZSCAN zset scan
const ZSCAN = `local h_key = ARGV[1]
local h_cursor = ARGV[2]
local h_count = ARGV[3]
local results = redis.call('ZSCAN',h_key,h_cursor,'COUNT',h_count)
local list = {}
if #results > 0 then
    table.insert(list,results[1])

    if #results[2] > 0 then
        for i = 1,#results[2],2 do
            table.insert(list,results[2][i])
        end
    end
end
return list`
