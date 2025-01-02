-- 限流对象
local key = KEYS[1]
-- 窗口大小
local window = tonumber(ARGV[1])
-- 阈值
local threshold = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
-- 窗口的起始时间
local min = now - window

-- 移除 Redis 有序集合中窗口起始时间之前的所有元素
redis.call('ZREMRANGEBYSCORE', key, '-inf', min)
-- 计算当前窗口内的元素数量
local cnt = redis.call('ZCOUNT', key, '-inf', '+inf')
-- local cnt = redis.call('ZCOUNT', key, min, '+inf')
-- 如果当前窗口内的元素数量达到或超过阈值，则执行限流
if cnt >= threshold then
    -- 执行限流
    return "true"
else
    -- 把当前时间添加到有序集合中，作为新的元素
    redis.call('ZADD', key, now, now)
    -- 设置有序集合的过期时间，与窗口大小一致
    redis.call('PEXPIRE', key, window)
    -- 不执行限流
    return "false"
end
