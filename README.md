# ztcache
Go implementation of "zero time cache pattern"

    import "github.com/hnw/ztcache"

Package ztcache は zero time cache pattern の実装です。
同じ結果を返す処理が並行に呼び出された場合に、1回の処理にまとめて相乗りさせるような仕組みを提供します。

## Usage

#### type ZTCache

```go
type ZTCache struct {
}
```

ZTCache は動作中の処理を管理する構造体です。

#### func  New

```go
func New() *ZTCache
```
New は ZTCache構造体のコンストラクタです。

#### func (*ZTCache) Get

```go
func (c *ZTCache) Get(key string, f func() string) (interface{}, error)
```
Get は処理結果を取得するためのラッパー関数です。
同じkeyで同時に呼び出された処理がある場合に並行に走らせず、前の処理の終了を待った上で必要なら相乗りして実行結果を複数呼び出し間で共有します。
