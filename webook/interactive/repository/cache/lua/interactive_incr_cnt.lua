local key =KEYS[1]
-- 对应到的就是hincrby 中的field
local cntKey=ARGV[1]
local delta=tonumber(ARGV[2])
local exists=redis.call("EXISTS",key)
if exists==1 then
--     delta +1 or -1
    redis.call("HINCRBY",key,cntKey,delta")
--     说明自增成功了
    return 1
else
--     自增不成功
    return 0
end
