package caydb

import (
	"fmt"
	"testing"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/quad"
)

func BenchmarkWriteQuads(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()

	store := NewCayleyDB("bolt", "cayley.db")
	defer store.Close()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := WriteQuads(store, Entity{ID: quad.IRI(fmt.Sprintf("一号柜:%d", i)), Type: "容器", Location: "主卧"})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCayleyCreate(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()

	store := NewCayleyDB("bolt", "cayley2.db")
	defer store.Close()
	qw := graph.NewWriter(store)
	qw.Close()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := sch.WriteAsQuads(qw, Entity{ID: quad.IRI(fmt.Sprintf("一号柜:%d", i)), Type: "容器", Location: "主卧"})
		if err != nil {
			b.Fatal(err)
		}
	}
}
