package bootstrap

import (
	"database/sql"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/redis/go-redis/v9"

	"LostAndFound/internal/adapters/postgres"
	cache "LostAndFound/internal/adapters/redis"
	s3storage "LostAndFound/internal/adapters/s3"
	sc "LostAndFound/internal/config/storage_config"
	"LostAndFound/internal/domain/repository"
)

type Deps struct {
	UserRepo  repository.UserRepo
	CardRepo  repository.CardRepo
	CacheRepo repository.CacheRepo
	FileStore repository.FileStorage
}

func Init(pg *sql.DB, rd *redis.Client, s3c *s3.S3, cfg sc.S3Config) *Deps {
	return &Deps{
		UserRepo:  postgres.NewUserRepo(pg),
		CardRepo:  postgres.NewCardRepo(pg),
		CacheRepo: cache.NewCacheRepo(rd),
		FileStore: s3storage.NewFileStorage(s3c, cfg),
	}
}
