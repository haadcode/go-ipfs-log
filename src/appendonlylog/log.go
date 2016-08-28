package appendonlylog

import (
  "log"
  "container/list"
  "immutabledb/interface"
)

type AppendOnlyLog struct {
  Id    string
  items *list.List
  db immutabledb.ImmutableDB
}

func New(id string, db immutabledb.ImmutableDB) *AppendOnlyLog {
  return new(AppendOnlyLog).Init(id, db)
}

func (l *AppendOnlyLog) Init(id string, db immutabledb.ImmutableDB) *AppendOnlyLog {
  l.Id = id
  l.db = db
  l.items = list.New()
  return l
}

func (l *AppendOnlyLog) Add(data []byte) *Entry {
  hash := l.db.Put(data)

  var next []*Entry
  n := l.items.Front()

  if n != nil {
    next = []*Entry{n.Value.(*Entry)}
  }

  e := &Entry{
    Key: hash,
    Value: data,
    Next: next,
  }

  l.items.PushFront(e)

  return e
}

func (l *AppendOnlyLog) Join(other *AppendOnlyLog) *AppendOnlyLog {
  items := l.Items()
  s := other.Items()

  for i := len(s) - 1; i >= 0; i-- {
    if (!contains(items, s[i])) {
      l.items.PushFront(s[i])
    }
  }

  return l
}

func contains(s []*Entry, e *Entry) bool {
  for _, a := range s {
    if a.Key == e.Key {
      return true
    }
  }
  return false
}

func (l *AppendOnlyLog) Head() *Entry {
  head := l.items.Front().Value.(*Entry)
  return head
}

func (l *AppendOnlyLog) Items() []*Entry {
  var more []*Entry

  for e := l.items.Front(); e != nil; e = e.Next() {
    entry := e.Value.(*Entry)
    more = append(more, entry)
  }

  return more
}

func (l *AppendOnlyLog) Print() {
  for _, e := range l.Items() {
    log.Println("Entry:", e.Key)

    data := l.db.Get(e.Key)

    log.Println("Data:", string(data))

    if (len(e.Next) > 0) {
      for _, n := range e.Next {
        log.Println("Next:", n.Key)
      }
    }

    log.Println()
  }
}
