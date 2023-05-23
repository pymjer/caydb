package caydb

import (
	"context"
	"fmt"
	"testing"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/schema"
	"github.com/cayleygraph/quad"
	uuid "github.com/satori/go.uuid"
)

var store *cayley.Handle
var sch *schema.Config

type Entity struct {
	ID       quad.IRI `quad:"@id" json:"id"`
	Type     string   `json:"type"` // required field
	Name     string   `quad:"name,optional" json:"name"`
	Location quad.IRI `quad:"location,optional" json:"location"` // can be empty
	Remark   string   `quad:"remark,optional" json:"remark"`
	Count    int      `quad:"count,optional" json:"count"`
}

func (e *Entity) String() string {
	return fmt.Sprintf("Entity{ID:%s, Type:%s, Name:%s, Location:%s, Remark:%s, Count:%d}", e.ID, e.Type, e.Name, e.Location, e.Remark, e.Count)
}

func TestCayleySchemaCreate(t *testing.T) {
	Entitylist := []Entity{
		{ID: "一号柜", Type: "容器", Location: "主卧"},
		{ID: "二号柜", Type: "容器", Location: "次卧"},
		{ID: "三号柜", Type: "容器", Location: "儿童房"},
		{ID: "指甲刀", Type: "物品", Location: "一号柜"},
		{ID: "夏凉被", Type: "物品", Location: "一号柜", Remark: "带有菊花图案"},
		{ID: "数据线", Type: "物品", Location: "二号柜", Count: 20},
		{ID: "充电宝", Type: "物品", Location: "二号柜", Count: 2},
		{ID: "电脑", Type: "物品", Location: "三号柜", Count: 1},
		{ID: "MP3", Type: "物品", Location: "三号柜", Count: 1},
	}
	qw := graph.NewWriter(store)

	for _, v := range Entitylist {
		id, err := sch.WriteAsQuads(qw, v)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("generated id:", id)
	}
	qw.Close()

	var allEntity []Entity
	sch.LoadTo(context.TODO(), store, &allEntity)
	fmt.Println(allEntity)
}

func TestCayleySchemaQuery(t *testing.T) {
	sch := schema.NewConfig()
	var allEntity []Entity
	sch.LoadTo(context.TODO(), store, &allEntity)
	fmt.Println(allEntity)
}

func TestHelloWorld(t *testing.T) {
	store.AddQuad(quad.Make("phrase of the day", "is of course", "Hello World!", nil))

	// Now we create the path, to get to our data
	p := cayley.StartPath(store, quad.String("phrase of the day")).Out(quad.String("is of course"))

	// Now we iterate over results. Arguments:
	// 1. Optional context used for cancellation.
	// 2. Quad store, but we can omit it because we have already built path with it.
	err := p.Iterate(context.TODO()).EachValue(store, func(value quad.Value) {
		nativeValue := quad.NativeOf(value) // this converts RDF values to normal Go types
		fmt.Println(nativeValue)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCayleyGizmo(t *testing.T) {
	// 查询所有容器
	//g.V().has('<type>','<容器>').all()
	p := cayley.StartPath(store).Has(quad.IRI("type"), quad.String("容器"))
	t.Log("查询所有容器")
	err := p.Iterate(context.TODO()).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value) // this converts RDF values to normal Go types
		fmt.Println(nativeValue)
	})
	if err != nil {
		t.Fatal(err)
	}
	// 查询所有物品
	// g.V().has('<type>','<物品>').all()
	t.Log("查询所有物品")
	p = cayley.StartPath(store).Has(quad.IRI("type"), quad.String("物品"))
	err = p.Iterate(context.TODO()).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value) // this converts RDF values to normal Go types
		fmt.Println(nativeValue)
	})
	if err != nil {
		t.Fatal(err)
	}
	// 查询一号柜的物品
	// g.V("<一号柜>").in().all()
	t.Log("查询一号柜的物品")
	p = cayley.StartPath(store, quad.IRI("一号柜")).In()
	err = p.Iterate(context.TODO()).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value) // this converts RDF values to normal Go types
		fmt.Println(nativeValue)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCayleySchemaGet(t *testing.T) {
	sch := schema.NewConfig()
	var gd Entity
	sch.LoadTo(context.TODO(), store, &gd, quad.IRI("充电宝"))
	fmt.Println(gd)
}

func TestCayleySchemaDelete(t *testing.T) {
	EntityIds := []string{"10001", "10002"}
	sch := schema.NewConfig()
	qw := graph.NewRemover(store)

	var gd Entity
	for _, gid := range EntityIds {
		sch.LoadTo(context.TODO(), store, &gd, quad.IRI(gid))
		id, err := sch.WriteAsQuads(qw, gd)
		if err != nil {
			t.Fatalf("删除Entity： %s 失败: %s/n", qw, err)
		}
		fmt.Println("remove id:", id)
	}
	qw.Close()

	var allEntity []Entity
	sch.LoadTo(context.TODO(), store, &allEntity)
	fmt.Println(allEntity)
}

func TestCreateAndGet(t *testing.T) {
	sch := schema.NewConfig()
	sch.GenerateID = func(_ interface{}) quad.Value {
		return quad.IRI(uuid.NewV4().String())
	}
	Entitylist := []Entity{
		{ID: "10001", Name: "test1-1106", Type: "电器", Location: "5号柜", Count: 1},
		{ID: "10002", Name: "test1", Type: "电器", Location: "5号柜", Count: 1},
	}
	qw := graph.NewWriter(store)
	for _, v := range Entitylist {
		id, err := sch.WriteAsQuads(qw, v)
		if err != nil {
			t.Fail()
		}
		fmt.Println("generated id:", id)
	}
	qw.Close()

	var allEntity []Entity
	sch.LoadTo(context.TODO(), store, &allEntity)
	fmt.Println(allEntity)
}

// 所有测试方法之前会调用
func TestMain(m *testing.M) {
	// 打开数据库连接
	store = NewCayleyDB("bolt", "test.db")
	sch = schema.NewConfig()
	// 运行测试
	m.Run()
	// 关闭数据库连接
	store.Close()
}
