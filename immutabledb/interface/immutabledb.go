package immutabledb

type ImmutableDB interface {
  Open(path string) ImmutableDB
  Close() error
  Put ([]byte) string
  Get (string) []byte
}
