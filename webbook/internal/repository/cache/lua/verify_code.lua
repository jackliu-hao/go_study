
-- 获取验证码的键
local key = KEYS[1]
-- 构造验证码计数器的键
local cntKey = key..":cnt"
-- 用户输入的验证码
local expectedCode = ARGV[1]

-- 获取当前验证码的计数器值
local cnt = tonumber(redis.call("get", cntKey))
-- 获取当前的验证码
local code = redis.call("get", key)

-- 检查计数器是否已经耗尽
if cnt == nil or cnt <= 0 then
    -- 验证次数耗尽了
    return -1
end

-- 比较用户输入的验证码和存储的验证码是否一致
if code == expectedCode then
    -- 验证成功，将计数器重置为0
    redis.call("set", cntKey, 0)
    return 0
else
    -- 验证失败，计数器减1
    redis.call("decr", cntKey)
    -- 不相等，用户输错了
    return -2
end
