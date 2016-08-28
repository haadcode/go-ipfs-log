package ipfs

import (
    "log"
    "github.com/ipfs/go-ipfs/core"
    "golang.org/x/net/context"
    path "github.com/ipfs/go-ipfs/path"
    fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
    dag "github.com/ipfs/go-ipfs/merkledag"
    "immutabledb/interface"
)

var dataDirectory = "/tmp/go-ipfs-log-dev"

type ImmutableIPFS struct {
  immutabledb.ImmutableDB
}

func (db ImmutableIPFS) Put(data []byte) string {
  r, err := fsrepo.Open(dataDirectory)
  if err != nil {
    log.Fatal("Can't open data repository at %s: %s", dataDirectory, err)
  }

  cfg := &core.BuildCfg{
    Repo:   r,
    Online: false,
  }

  nd, err := core.NewNode(context.Background(), cfg)
  if err != nil {
    log.Fatal(err)
  }

  obj := &dag.Node{
    Data: data,
  }

  k, err := nd.DAG.Add(obj)
  if err != nil {
    log.Fatal(err)
  }

  return k.B58String()
}

func (db ImmutableIPFS)  Get(key string) []byte {
  r, err := fsrepo.Open(dataDirectory)
  if err != nil {
    log.Fatal("Can't open data repository at %s: %s", dataDirectory, err)
  }

  cfg := &core.BuildCfg{
    Repo:   r,
    Online: false,
  }

  ctx := context.Background()

  nd, err := core.NewNode(ctx, cfg)
  if err != nil {
    log.Fatal(err)
  }

  fpath := path.Path(key)

  object, err := core.Resolve(ctx, nd, fpath)
  if err != nil {
    log.Fatal(err)
  }

  node := &dag.Node{
    Data:  object.Data,
  }

  return node.Data
}
