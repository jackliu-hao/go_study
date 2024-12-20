-- 获取传入的key，用于标识验证码
-- phone_code:login:152xxxx
local key = KEYS[1]
-- 构造计数器的key
-- phone_code:login:152xxxx:cnt
local cntKey = key..":cnt"
-- 你准备的存储的验证码
local val = ARGV[1]

-- 获取key的剩余生存时间
local ttl = tonumber(redis.call("ttl", key))
-- 检查key的生存时间，决定是否发送验证码
if ttl == -1 then
    --    key 存在，但是没有过期时间
    return -2
elseif ttl == -2 or ttl < 540 then
    --    可以发验证码 , 540 = 600 -60 ，超过一分钟了 ， 并且key不存在，key不存在的时候，ttl == -2
    redis.call("set", key, val)
    -- 600 秒
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 发送太频繁
    return -1
end
