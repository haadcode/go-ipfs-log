package appendonlylog

import (
  "fmt"
  "testing"
  "strings"
  "bytes"
  db "github.com/haadcode/go-ipfs-log/immutabledb/ipfs"
)

var dataDirectory = "/tmp/go-ipfs-log-dev"
var ipfsdb = db.Open(dataDirectory)

var id = "abc"

// func BenchmarkHello(b *testing.B) {
//     for i := 0; i < b.N; i++ {
//         fmt.Sprintf("hello")
//     }
// }

func TestNew(t *testing.T) {
  var log1 = New(id, ipfsdb)

  if log1 == nil {
    t.Errorf("Couldn't create a log")
  }

  if log1.Id != id {
    t.Errorf("Id not set")
  }

  if log1.db == nil {
    t.Errorf("DB not set")
  }

  if len(log1.Items()) != 0 {
    t.Errorf("Items not empty")
  }
}

func TestAdd(t *testing.T) {
  var value = []byte("Hello1")
  var expectedHash = "QmX96xhp6cUB1YE5nqZsmKHbZFiAEderPc3gapGdwAoEod"

  var log1 = New(id, ipfsdb)

  one := log1.Add(value)

  if one == nil {
    t.Errorf("Entry was not added")
  }

  if strings.Compare(one.Key, expectedHash) != 0 {
    t.Errorf("Wrong key: %s", one.Key)
  }

  if bytes.Compare(one.Value, value) != 0 {
    t.Errorf("Wrong key: %s", one.Key)
  }

  if len(one.Next) != 0 {
    t.Errorf("Wrong next reference: %s", one.Next)
  }

  if len(log1.Items()) != 1 {
    t.Errorf("Wrong items count: %i", len(log1.Items()))
  }
}

func ExampleAdd_one() {
  var value = []byte("Hello1")
  var log1 = New(id, ipfsdb)
  one := log1.Add(value)
  fmt.Println(one.Key)
  fmt.Println(string(one.Value))
  fmt.Println(one.Next)
  // Output:
  // QmX96xhp6cUB1YE5nqZsmKHbZFiAEderPc3gapGdwAoEod
  // Hello1
  // []
}

func ExampleAdd_two() {
  var value1 = []byte("Hello1")
  var value2 = []byte("Hello2")
  var log1 = New(id, ipfsdb)
  log1.Add(value1)
  two := log1.Add(value2)
  fmt.Println(two.Key)
  fmt.Println(string(two.Value))
  fmt.Println(two.Next[0].Key)
  // Output:
  // Qme39B2h1QTDYAwCa4gXa6DB6R3TAFaG2Z8HF48U1wkKE6
  // Hello2
  // QmX96xhp6cUB1YE5nqZsmKHbZFiAEderPc3gapGdwAoEod
}

func ExampleAdd_three() {
  var value1 = []byte("Hello1")
  var value2 = []byte("Hello2")
  var value3 = []byte("Hello3")
  var log1 = New(id, ipfsdb)
  log1.Add(value1)
  log1.Add(value2)
  log1.Add(value3)

  items := log1.Items()
  fmt.Println(len(items))
  fmt.Println(string(items[0].Value))
  fmt.Println(string(items[1].Value))
  fmt.Println(string(items[2].Value))
  // Output:
  // 3
  // Hello1
  // Hello2
  // Hello3
}
