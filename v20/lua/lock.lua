-- 判断能不能拿到锁
local val = redis.call('get', KEYS[1])
if val == false then
-- 锁不存在，那么自己加锁
    return redis.call('set', KEYS[1], ARGV[1], 'PX', ARGV[2])
elseif val == ARGV[1] then
-- 如果锁存在，而且是自己的，就刷新过期时间
    redis.call('expire', KEYS[1], ARGV[2])
    return  "OK"
else
-- 锁存在，但是不是自己的，就返回 ""，表示锁已经被占有
    return ""
end