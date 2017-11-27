package managers

import (
	"context"
	"os"
	"testing"
	"time"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

type base struct {
	ID int64 `json:"id" datastore:"-"`
}

func (b *base) key(ctx context.Context, kind string) *datastore.Key {
	if b.ID == 0 {
		return datastore.NewIncompleteKey(ctx, kind, nil)
	}
	return datastore.NewKey(ctx, kind, "", b.ID, nil)
}

func (b *base) set(id int64) {
	b.ID = id
}

type Organization struct {
	base

	Name string `json:"name"`
}

func (o *Organization) key(ctx context.Context) *datastore.Key {
	return o.base.key(ctx, "Organization")
}

func save(ctx context.Context, src *Organization) error {
	key, err := datastore.Put(ctx, src.key(ctx), src)
	src.set(key.IntID())
	return err
}

func load(ctx context.Context, id int64, src *Organization) error {
	src.set(id)
	return datastore.Get(ctx, src.key(ctx), src)
}

func TestAdd(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	src := &Organization{
		Name: "asdf",
	}

	err = save(ctx, src)
	t.Log(src, err)

	var dst Organization
	err = load(ctx, src.ID, &dst)
	t.Log(dst, err)

	if os.Getenv("WAIT") != "" {
		<-time.After(time.Second * 30)
	}
}
