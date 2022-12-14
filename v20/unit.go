package v20

import "time"

// s6Unit 缓存单元
// 这玩意本地缓存可以直接存在内存里没问题
// 但是如果是 redis 这种实现，那就需要考虑序列化
type s6Unit struct {
	// value 缓存的值
	value any
	// deadline 过期时间
	deadline time.Time
}

// f8CheckDeadline 缓存过期没有
// true=没过期；false=过期
func (p7this *s6Unit) f8CheckDeadline(checkTime time.Time) bool {
	if p7this.deadline.IsZero() {
		// 如果没有设置过期时间，那就不会过期
		return true
	}
	// 否则比较一下校验时间是不是在过期时间之前
	return checkTime.Before(p7this.deadline)
}
