package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	rdb "github.com/redis/go-redis/v9"
)

type RedisWrapper struct {
    ctx context.Context
    rdb *rdb.Client
}

func NewRedisWrapper(
    client *rdb.Client,
) *RedisWrapper {
    return &RedisWrapper{ context.Background(), client };
}

func (self *RedisWrapper) GetLoginAttempts(id uint) (int, error) {
    val, err := self.rdb.
        Get(self.ctx, self.loginAttemptsKey(id)).
        Result();

    if err != nil {
        if err == rdb.Nil {
            return 0, nil;
        }
        return 0, err;
    }

    if val == "" {
        return 0, nil;
    }

    valInt, err := strconv.Atoi(val);
    if err != nil {
        return 0, err;
    }

    return valInt, nil;
}

func (self *RedisWrapper) SetLoginAttempts(id uint, attempts int) error {
    return self.rdb.SetEx(
        self.ctx,
        self.loginAttemptsKey(id),
        strconv.Itoa(attempts),
        5 * time.Minute,
    ).Err();
}

func (self *RedisWrapper) GetLoginBlocked(id uint) (bool, error) {
    val, err := self.rdb.Get(self.ctx, self.loginBlockedKey(id)).Result();
    if err != nil {
        if err == rdb.Nil {
            return false, nil;
        }
        return false, err;
    }

    if val == "" {
        return false, nil;
    }

    return val == "true", nil;
}

func (self *RedisWrapper) SetLoginBlocked(id uint, blocked bool) error {
    // FIXME
    blockedStr := "false";
    if blocked == true {
        blockedStr = "true";
    }

    return self.rdb.SetEx(
        self.ctx,
        self.loginBlockedKey(id),
        blockedStr,
        5 * time.Minute,
    ).Err();
}

func (self *RedisWrapper) loginAttemptsKey(id uint) string {
    return fmt.Sprintf("auth:login_attempts:%d", id);
}

func (self *RedisWrapper) loginBlockedKey(id uint) string {
    return fmt.Sprintf("auth:login_blocked:%d", id);
}
