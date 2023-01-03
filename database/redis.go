package database

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
	"time"
	"vsz-web-backend"
	"vsz-web-backend/config"
)

const SessionTime = time.Hour * 24 * 30

var (
	rdb *redis.Client
	ctx = context.Background()
)

func ConnectRedis() error {
	rediscfg := config.Global.Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     rediscfg.Host + ":6379",
		Password: rediscfg.Pass,
		DB:       rediscfg.DB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println("Finished connecting to redis database")
	return nil
}

func StoreSession(session *vsz_web_backend.Session) (uuid.UUID, error) {
	serializedSession, err := json.Marshal(session)
	if err != nil {
		return uuid.UUID{}, err
	}

	sessionuuid, err := uuid.NewRandom()
	if err != nil {
		return uuid.UUID{}, err
	}

	err = rdb.Set(ctx, "session:"+sessionuuid.String(), serializedSession, SessionTime).Err()
	if err != nil {
		return uuid.UUID{}, err
	}

	return sessionuuid, nil
}

func GetSession(session string) (*vsz_web_backend.Session, error) {
	res, err := rdb.Get(ctx, "session:"+session).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var usersession *vsz_web_backend.Session
	err = json.Unmarshal([]byte(res), &usersession)
	if err != nil {
		return nil, err
	}

	return usersession, nil
}

func DeleteSession(session string) error {
	return rdb.Del(ctx, "session:"+session).Err()
}
