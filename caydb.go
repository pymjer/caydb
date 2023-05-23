package caydb

import (
	"log"
	"os"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/kv/bolt"
	_ "github.com/cayleygraph/cayley/graph/sql/mysql"
	"github.com/cayleygraph/cayley/schema"
	"github.com/cayleygraph/quad"
)

var defaultWriter graph.BatchWriter
var defaultSchemaCfg *schema.Config

func NewCayleyDB(backend string, address string) (store *cayley.Handle) {
	log.Println("Try cayley backend connection")
	// 只有第一次可以创建，后面都是打开
	var err error
	switch backend {
	case "mysql":
		store, err = cayley.NewGraph("mysql", address, nil)
		if err != nil {
			log.Fatalf("Failure database connection: %s\n", err.Error())
			os.Exit(0)
		}
	case "bolt":
		if _, err := os.Stat(address); os.IsNotExist(err) {
			err := graph.InitQuadStore("bolt", address, nil)
			if err != nil {
				log.Fatalf("Failure database connection: %s\n", err.Error())
				os.Exit(0)
			}
		}
		store, err = cayley.NewGraph(backend, address, nil)
		if err != nil {
			log.Fatalf("Failure database connection\n")
			os.Exit(0)
		}
	default:
		log.Fatalf("不支持的数据库类型: %s\n", backend)
		os.Exit(0)
	}
	return store
}

func CleanUp() {
	if defaultWriter != nil {
		defaultWriter.Close()
	}
}

func WriteQuads(store *cayley.Handle, o interface{}) (id quad.Value, err error) {
	if defaultWriter == nil {
		defaultWriter = graph.NewWriter(store)
	}
	if defaultSchemaCfg == nil {
		defaultSchemaCfg = schema.NewConfig()
	}
	id, err = defaultSchemaCfg.WriteAsQuads(defaultWriter, o)
	return
}
