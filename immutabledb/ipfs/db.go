package ipfs

import (
    "log"
    "github.com/haadcode/go-ipfs-log/immutabledb/interface"
    "golang.org/x/net/context"
    "github.com/ipfs/go-ipfs/core"
    path "github.com/ipfs/go-ipfs/path"
    repo "github.com/ipfs/go-ipfs/repo"
    fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
    dag "github.com/ipfs/go-ipfs/merkledag"
)

// Trick to make sure ImmutableIPFS implements ImmutableDB
var _ immutabledb.ImmutableDB = ImmutableIPFS{}

type ImmutableIPFS struct {
  Repo repo.Repo
  Node *core.IpfsNode
}

func Open(path string) ImmutableIPFS {
  r, err := fsrepo.Open(path)
  if err != nil {
    log.Fatal("Can't open data repository at %s: %s", path, err)
  }

  cfg := &core.BuildCfg{
    Repo:   r,
    Online: false,
  }

  nd, err := core.NewNode(context.Background(), cfg)
  if err != nil {
    log.Fatal("Can't create IPFS node: %s", err)
  }

  return ImmutableIPFS{
    Repo: r,
    Node: nd,
  }
}

func (db ImmutableIPFS) Close() error {

}

func (db ImmutableIPFS) Put(data []byte) string {
  obj := &dag.Node{
    Data: data,
  }

  k, err := db.Node.DAG.Add(obj)
  if err != nil {
    log.Fatal(err)
  }

  return k.B58String()
}

func (db ImmutableIPFS) Get(key string) []byte {
  ctx := context.Background()
  fpath := path.Path(key)

  object, err := core.Resolve(ctx, db.Node, fpath)
  if err != nil {
    log.Fatal(err)
  }

  node := &dag.Node{
    Data:  object.Data,
  }

  return node.Data
}
