package vote_storage_test

import (
	"math/rand"
	"testing"

	"github.com/dancannon/gorethink"
	"github.com/dchest/uniuri"
)

var (
	session      *gorethink.Session
	key2find     string
	table2search string
)

func init() {
	var err error
	session, err = gorethink.Connect(gorethink.ConnectOpts{
		Address: "127.0.0.1:28015",
	})
	if err != nil {
		panic(err)
	}

	key2find = uniuri.New()

	// Create a new table
	gorethink.Db("test").TableDrop("benchmark_keys_list").Run(session)
	gorethink.Db("test").TableCreate("benchmark_keys_list").Run(session)

	var klist []*KeysList

	// Populate with sample data
	for n := 0; n < 300; n++ {
		keys := rndStringSlice(999)
		keys = randomlyInsert(keys, key2find)

		y := uniuri.New()
		if n == 153 {
			table2search = y
		}

		klist = append(klist, &KeysList{
			ID:    y,
			Voted: keys,
		})
	}

	gorethink.Db("test").Table("benchmark_keys_list").Insert(klist).Run(session)
}

func rndStringSlice(count int) []string {
	var r []string
	for i := 0; i < count; i++ {
		r = append(r, uniuri.New())
	}
	return r
}

func randomlyInsert(s []string, x string) []string {
	i := rand.Intn(len(s) - 1)

	s = append(s[:i], append([]string{x}, s[i:]...)...)

	return s
}

type KeysList struct {
	ID    string   `gorethink:"id"`
	Voted []string `gorethink:"voted"`
}

func BenchmarkContains(b *testing.B) {
	for n := 0; n < b.N; n++ {
		contains, err := gorethink.Db("test").Table("benchmark_keys_list").Get(table2search).Field("voted").Contains(key2find).Run(session)
		if err != nil {
			b.Log(err)
			b.Fail()
		}

		var res bool
		err = contains.One(&res)
		if err != nil {
			b.Log(err)
			b.Fail()
		}
		if !res {
			b.Log("invalid response")
			b.Fail()
		}
	}
}

func BenchmarkAppend(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := gorethink.Db("test").Table("benchmark_keys_list").Get(table2search).Field("voted").Append(uniuri.New()).Run(session)
		if err != nil {
			b.Log(err)
			b.Fail()
		}
	}
}

func BenchmarkDelete(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := gorethink.Db("test").Table("benchmark_keys_list").Get(table2search).Field("voted").DeleteAt(
			gorethink.Expr(gorethink.Db("test").Table("benchmark_keys_list").Get(table2search).Field("voted").IndexesOf(key2find).AtIndex(0)),
		).Run(session)
		if err != nil {
			b.Log(err)
			b.Fail()
		}
	}
}
