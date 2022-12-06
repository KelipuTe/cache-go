package v20

import "time"

// 缓存单元
type s6Unit struct {
	// 值
	value any
	// 过期时间
	deadline time.Time
}

// true=没过期；false=过期
func (p7this *s6Unit) f8CheckDeadline(t4time time.Time) bool {
	if p7this.deadline.IsZero() {
		// 如果没有设置过期时间
		return true
	}
	// 否则比较一下过期时间和当前时间
	return p7this.deadline.Before(t4time)
}
