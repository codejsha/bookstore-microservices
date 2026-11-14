package support

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type cachedProfile struct {
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func TestCacheClientSet_structValue_storesJsonWithTtl(t *testing.T) {
	c, valkey := newFakeValkeyClient(t)

	err := c.Set(context.Background(), "profile:1", cachedProfile{Email: "a@x.com", Age: 30}, 90*time.Second)
	if err != nil {
		t.Fatalf("set: %v", err)
	}

	stored, ok := valkey.value("profile:1")
	if !ok || stored != `{"email":"a@x.com","age":30}` {
		t.Fatalf("stored = %q (%v), want the json document", stored, ok)
	}
	if got := valkey.ttl("profile:1"); got != 90*time.Second {
		t.Errorf("ttl = %v, want 90s", got)
	}
}

func TestCacheClientSet_unmarshalableValue_reportsErrorWithoutWriting(t *testing.T) {
	c, valkey := newFakeValkeyClient(t)

	err := c.Set(context.Background(), "profile:1", make(chan int), time.Minute)

	if err == nil {
		t.Fatal("err = nil, want a marshal error")
	}
	if writes := valkey.commandsNamed("SET"); len(writes) != 0 {
		t.Errorf("writes = %v, want none", writes)
	}
}

func TestCacheClientGet_storedJson_decodesIntoDest(t *testing.T) {
	c, valkey := newFakeValkeyClient(t)
	valkey.put("profile:1", `{"email":"a@x.com","age":30}`)

	var got cachedProfile
	if err := c.Get(context.Background(), "profile:1", &got); err != nil {
		t.Fatalf("get: %v", err)
	}

	if want := (cachedProfile{Email: "a@x.com", Age: 30}); got != want {
		t.Fatalf("profile = %+v, want %+v", got, want)
	}
}

func TestCacheClientGet_missingKey_reportsRedisNil(t *testing.T) {
	c, _ := newFakeValkeyClient(t)

	var got cachedProfile
	err := c.Get(context.Background(), "profile:missing", &got)

	if !errors.Is(err, redis.Nil) {
		t.Fatalf("err = %v, want redis.Nil", err)
	}
	if got != (cachedProfile{}) {
		t.Errorf("profile = %+v, want dest untouched", got)
	}
}

func TestCacheClientGet_storedValueNotJson_reportsDecodeError(t *testing.T) {
	c, valkey := newFakeValkeyClient(t)
	valkey.put("profile:1", "not-json")

	var got cachedProfile
	err := c.Get(context.Background(), "profile:1", &got)

	if err == nil || errors.Is(err, redis.Nil) {
		t.Fatalf("err = %v, want a decode error distinct from a miss", err)
	}
}

func TestCacheClientMGet_someKeysMissing_keepsPositionsWithNilGaps(t *testing.T) {
	c, valkey := newFakeValkeyClient(t)
	valkey.put("k1", "v1")
	valkey.put("k3", "v3")

	got, err := c.MGet(context.Background(), "k1", "k2", "k3")
	if err != nil {
		t.Fatalf("mget: %v", err)
	}

	if want := []any{"v1", nil, "v3"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("values = %#v, want %#v", got, want)
	}
}

func TestCacheClientMGet_noKeys_skipsTheRoundTrip(t *testing.T) {
	c, valkey := newFakeValkeyClient(t)

	got, err := c.MGet(context.Background())
	if err != nil {
		t.Fatalf("mget: %v", err)
	}

	if got != nil {
		t.Errorf("values = %#v, want nil", got)
	}
	if calls := valkey.commandsNamed("MGET"); len(calls) != 0 {
		t.Errorf("mget commands = %v, want none", calls)
	}
}

func TestCacheClientDel_severalKeys_removesEachOfThem(t *testing.T) {
	c, valkey := newFakeValkeyClient(t)
	valkey.put("k1", "v1")
	valkey.put("k2", "v2")
	valkey.put("k3", "v3")

	if err := c.Del(context.Background(), "k1", "k2", "absent"); err != nil {
		t.Fatalf("del: %v", err)
	}

	for _, key := range []string{"k1", "k2"} {
		if _, ok := valkey.value(key); ok {
			t.Errorf("%s still stored, want it deleted", key)
		}
	}
	if _, ok := valkey.value("k3"); !ok {
		t.Errorf("k3 deleted, want it kept")
	}
}

func TestCacheClientExists_keyPresence_reportsWhetherStored(t *testing.T) {
	c, valkey := newFakeValkeyClient(t)
	valkey.put("present", "v")

	for key, want := range map[string]bool{"present": true, "absent": false} {
		got, err := c.Exists(context.Background(), key)
		if err != nil {
			t.Fatalf("exists %s: %v", key, err)
		}
		if got != want {
			t.Errorf("exists(%s) = %v, want %v", key, got, want)
		}
	}
}

func TestCacheClientSetOperations_membersAddedAndRemoved_trackMembership(t *testing.T) {
	c, valkey := newFakeValkeyClient(t)
	ctx := context.Background()

	if err := c.SAdd(ctx, "index", "a", "b", "c"); err != nil {
		t.Fatalf("sadd: %v", err)
	}
	if err := c.SAdd(ctx, "index", "b"); err != nil {
		t.Fatalf("sadd duplicate: %v", err)
	}
	if err := c.SRem(ctx, "index", "a", "absent"); err != nil {
		t.Fatalf("srem: %v", err)
	}
	got, err := c.SMembers(ctx, "index")
	if err != nil {
		t.Fatalf("smembers: %v", err)
	}

	if want := []string{"b", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("members = %v, want %v", got, want)
	}
	if want := []string{"b", "c"}; !reflect.DeepEqual(valkey.members("index"), want) {
		t.Fatalf("stored members = %v, want %v", valkey.members("index"), want)
	}
}

func TestCacheClientSMembers_unknownKey_emptyList(t *testing.T) {
	c, _ := newFakeValkeyClient(t)

	got, err := c.SMembers(context.Background(), "index")
	if err != nil {
		t.Fatalf("smembers: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("members = %v, want none", got)
	}
}

func TestCacheClientTryLock_lockFree_storesOwnerTokenUnderPrefixedKeyWithTtl(t *testing.T) {
	c, valkey := newFakeValkeyClient(t)

	token, err := c.TryLock(context.Background(), "refresh", 10*time.Second)
	if err != nil {
		t.Fatalf("try lock: %v", err)
	}

	if token == "" {
		t.Fatal("token = \"\", want an owner token for an acquired lock")
	}
	stored, ok := valkey.value("lock:refresh")
	if !ok || stored != token {
		t.Fatalf("lock:refresh = %q (%v), want the owner token %q", stored, ok, token)
	}
	if got := valkey.ttl("lock:refresh"); got != 10*time.Second {
		t.Errorf("ttl = %v, want 10s so a crashed holder cannot keep the lock", got)
	}
	if _, ok := valkey.value("refresh"); ok {
		t.Errorf("unprefixed key written, want only lock:refresh")
	}
}

func TestCacheClientTryLock_lockHeld_emptyTokenAndHolderKept(t *testing.T) {
	c, valkey := newFakeValkeyClient(t)
	holder, err := c.TryLock(context.Background(), "refresh", 10*time.Second)
	if err != nil {
		t.Fatalf("first try lock: %v", err)
	}

	token, err := c.TryLock(context.Background(), "refresh", 10*time.Second)
	if err != nil {
		t.Fatalf("second try lock: %v", err)
	}

	if token != "" {
		t.Fatalf("token = %q, want empty for a contended lock", token)
	}
	if stored, _ := valkey.value("lock:refresh"); stored != holder {
		t.Errorf("lock owner = %q, want the first holder %q", stored, holder)
	}
}

func TestCacheClientUnlock_ownerToken_releasesTheLock(t *testing.T) {
	c, valkey := newFakeValkeyClient(t)
	token, err := c.TryLock(context.Background(), "refresh", 10*time.Second)
	if err != nil {
		t.Fatalf("try lock: %v", err)
	}

	if err := c.Unlock(context.Background(), "refresh", token); err != nil {
		t.Fatalf("unlock: %v", err)
	}

	if _, ok := valkey.value("lock:refresh"); ok {
		t.Fatal("lock:refresh still stored, want it released")
	}
	again, err := c.TryLock(context.Background(), "refresh", 10*time.Second)
	if err != nil || again == "" {
		t.Fatalf("relock = (%q, %v), want the lock acquirable again", again, err)
	}
}

func TestCacheClientUnlock_foreignToken_sendsCompareAndDeleteAndKeepsTheLock(t *testing.T) {
	c, valkey := newFakeValkeyClient(t)
	holder, err := c.TryLock(context.Background(), "refresh", 10*time.Second)
	if err != nil {
		t.Fatalf("try lock: %v", err)
	}

	if err := c.Unlock(context.Background(), "refresh", "someone-else"); err != nil {
		t.Fatalf("unlock: %v", err)
	}

	if stored, _ := valkey.value("lock:refresh"); stored != holder {
		t.Fatalf("lock owner = %q, want the holder %q kept", stored, holder)
	}
	evals := valkey.commandsNamed("EVAL")
	if len(evals) != 1 {
		t.Fatalf("eval commands = %v, want 1", evals)
	}
	script, keyAndToken := evals[0][1], evals[0][2:]
	if !strings.Contains(script, "redis.call('GET', KEYS[1]) == ARGV[1]") || !strings.Contains(script, "redis.call('DEL', KEYS[1])") {
		t.Errorf("script = %q, want an owner-checked delete", script)
	}
	if want := []string{"1", "lock:refresh", "someone-else"}; !reflect.DeepEqual(keyAndToken, want) {
		t.Errorf("script args = %v, want %v", keyAndToken, want)
	}
}

func TestCacheClient_clientClosed_reportsErrorOnEveryOperation(t *testing.T) {
	ctx := context.Background()
	calls := map[string]func(c *CacheClient) error{
		"get": func(c *CacheClient) error {
			var dest cachedProfile
			return c.Get(ctx, "k", &dest)
		},
		"mget": func(c *CacheClient) error {
			_, err := c.MGet(ctx, "k")
			return err
		},
		"set": func(c *CacheClient) error { return c.Set(ctx, "k", "v", time.Minute) },
		"del": func(c *CacheClient) error { return c.Del(ctx, "k") },
		"exists": func(c *CacheClient) error {
			_, err := c.Exists(ctx, "k")
			return err
		},
		"sadd": func(c *CacheClient) error { return c.SAdd(ctx, "k", "m") },
		"srem": func(c *CacheClient) error { return c.SRem(ctx, "k", "m") },
		"smembers": func(c *CacheClient) error {
			_, err := c.SMembers(ctx, "k")
			return err
		},
		"tryLock": func(c *CacheClient) error {
			_, err := c.TryLock(ctx, "k", time.Second)
			return err
		},
		"unlock": func(c *CacheClient) error { return c.Unlock(ctx, "k", "token") },
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			c, _ := newFakeValkeyClient(t)
			if err := c.Close(); err != nil {
				t.Fatalf("close: %v", err)
			}

			err := call(c)

			if !errors.Is(err, redis.ErrClosed) {
				t.Fatalf("err = %v, want redis.ErrClosed", err)
			}
		})
	}
}
