package usecase

type Usecase struct {
	cache CacheInterface
	repo  RepoInterface
}

func New(cache CacheInterface, repo RepoInterface) *Usecase {
	return &Usecase{
		cache: cache,
		repo:  repo,
	}
}
