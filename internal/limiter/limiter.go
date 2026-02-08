package limiter

type Limiter interface {
	Allow(ip string) Result
}


type Result struct {
	Allowed		bool
	Limit		int
	Remaining	int
	ResetAt		int64 //unix ts in seconds
}