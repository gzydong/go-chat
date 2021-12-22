package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Relation struct {
	rds *redis.Client
}

func NewRelation(rds *redis.Client) *Relation {
	return &Relation{rds: rds}
}

func (r *Relation) keyContactRelation(uid, uid2 int) string {
	if uid2 < uid {
		uid, uid2 = uid2, uid
	}

	return fmt.Sprintf("rds:contact:relation:%d_%d", uid, uid2)
}

func (r *Relation) IsContactRelation(ctx context.Context, uid, uid2 int) error {
	return r.rds.Get(ctx, r.keyContactRelation(uid, uid2)).Err()
}

func (r *Relation) SetContactRelation(ctx context.Context, uid, uid2 int) {
	r.rds.SetEX(ctx, r.keyContactRelation(uid, uid2), "1", time.Hour*1)
}

func (r *Relation) DelContactRelation(ctx context.Context, uid, uid2 int) {
	r.rds.Del(ctx, r.keyContactRelation(uid, uid2))
}
