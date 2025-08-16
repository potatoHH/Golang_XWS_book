--发送到key,也就是code:业务:手机号码
local key =KEYS[1]
-- 使用次数,也就是验证次数
local conKey =key..":cnt"
-- 你的验证码
local val =ARGV[1]
-- 验证吗的有效时间是十分钟,600s
local ttl =tonumber(reids.call("ttl",key ))
if ttl =-1 then
--     有人错误操作
    return -2
--  -2 是key, 不存在,ttl<540 是发了一个验证码,已经超过一分钟了
elseif ttl ==-2 or ttl< 540 then
--     后续如果验证码的不同过期时间,要在这里优化
    redis.call("set",key,val)
    redis.call("expire",key,600)
    redis.call("set",conKey,3)
    redis.call("expire",conKey,600)
    return 0
else
--     已经发送了一个验证码,但是还没到一分钟
    return -1


end
