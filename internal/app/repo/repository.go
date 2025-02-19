package repo

type Repository interface {
	SaveSrcUrlByHashKey(
		hashKey string,
		srcUrl string,
	) (err error)
	GetSrcUrlByHashKey(hashKey string) (string, error)
}
