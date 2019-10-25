package main

import "context"

// 基于以上LRUCache、上下文、chanle、多goroutine实现

// nodeGoLRU node
type nodeGoLRU struct {
	Key, Value interface{}
}

// getResult set结果
type getResultNode struct {
	Value interface{} // 对应值
	Ret   bool        // 是否存在
	Err   error       // 错误
}

// GoLRUCache 多goroutine共享LRU缓存
type GoLRUCache struct {
	lruCache    *LRUCache
	setCh       chan nodeGoLRU     // 设置缓存
	setResultCh chan error         // 接收设置结果
	getCh       chan interface{}   // 读取缓存
	getResultCh chan getResultNode // 接收读取结果
}

// NewGoLRUCache New GoLRUCache
// caps 容量
// expire 有效期，分钟。nil不失效
func NewGoLRUCache(caps int, expire *float64) (*GoLRUCache, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	setCh := make(chan nodeGoLRU)
	setResultCh := make(chan error)
	getCh := make(chan interface{})
	getResultCh := make(chan getResultNode)
	goCache := &GoLRUCache{
		lruCache:    NewLruCache(caps, expire),
		getCh:       getCh,
		getResultCh: getResultCh,
		setCh:       setCh,
		setResultCh: setResultCh,
	}
	go func() {
	Loop:
		for {
			select {
			case node := <-goCache.setCh:
				goCache.setHandle(node)
			case k := <-goCache.getCh:
				goCache.getHandle(k)
			case <-ctx.Done():
				break Loop
			}
		}
		close(getCh)
		close(setResultCh)
		close(setCh)
		close(getResultCh)
	}()
	return goCache, cancel
}

// Set 设置
func (g *GoLRUCache) Set(k, v interface{}) error {
	node := nodeGoLRU{Key: k, Value: v}
	g.setCh <- node
	return <-g.setResultCh
}

// Get 取值
func (g *GoLRUCache) Get(k interface{}) (result interface{}, ret bool, err error) {
	g.getCh <- k
	getResult := <-g.getResultCh
	result, ret, err = getResult.Value, getResult.Ret, getResult.Err
	return
}

// Size 缓存大小
func (g *GoLRUCache) Size() int {
	return g.lruCache.Size()
}

// setHandle 设置缓存Handle
func (g *GoLRUCache) setHandle(node nodeGoLRU) {
	k, v := node.Key, node.Value
	g.setResultCh <- g.lruCache.Set(k, v)
}

// getHandle 获取结果
func (g *GoLRUCache) getHandle(k interface{}) {
	v, ret, err := g.lruCache.Get(k)
	g.getResultCh <- getResultNode{Value: v, Ret: ret, Err: err}
}
