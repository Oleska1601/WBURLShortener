package service

type Service struct {
	cache CacheI
	repo  RepoI
}

func New(cache CacheI, repo RepoI) *Service {
	return &Service{
		cache: cache,
		repo:  repo,
	}
}
