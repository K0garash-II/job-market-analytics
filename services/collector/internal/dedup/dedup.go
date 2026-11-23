package dedup

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Deduplicator struct{ client *redis.Client }

func New(client *redis.Client) *Deduplicator { return &Deduplicator{client} }
func fingerprint(source, id string) string {
	return fmt.Sprintf("dedup:%x", sha256.Sum256([]byte(source+":"+id)))
}
func (d *Deduplicator) Seen(ctx context.Context, source, id string) (bool, error) {
	n, err := d.client.Exists(ctx, fingerprint(source, id)).Result()
	return n > 0, err
}
func (d *Deduplicator) Mark(ctx context.Context, source, id string) error {
	return d.client.Set(ctx, fingerprint(source, id), 1, 24*time.Hour).Err()
}
