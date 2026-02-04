package limiter

type Limiter interface {
	Allow(ip string) bool
}