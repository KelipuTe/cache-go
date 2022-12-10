-- 判断能不能拿到锁，拿到的锁是不是自己的
if redis.call("get", KEYS[1]) == ARGV[1]
then
-- 锁存在，而且是自己的，就把自己删了，表示解锁
    return redis.call("del", KEYS[1])
else
-- 锁存在，但是不是自己的，就返回 0，表示锁不存在或者锁已经被占有
    return 0
end