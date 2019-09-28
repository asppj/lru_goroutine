package main

import (
	"container/list"
	"errors"
)

//cacheNode 节点
type cacheNode struct {
	Key, Value interface{}
}

//newCacheNode 新建节点
func (cnode *cacheNode) newCacheNode(k, v interface{}) *cacheNode {
	return &cacheNode{k, v}
}

//LRUCache LRU缓存
type LRUCache struct {
	Capacity int
	dlist    *list.List
	cacheMap map[interface{}]*list.Element
}

//NewLruCache 新建 LRUCache
func NewLruCache(cap int) *LRUCache {
	return &LRUCache{Capacity: cap, dlist: list.New(), cacheMap: make(map[interface{}]*list.Element)}
}

//Size 元素个数
func (lru *LRUCache) Size() int {
	return lru.dlist.Len()
}

//Set 设置缓存
func (lru *LRUCache) Set(k, v interface{}) error {
	//异常
	if lru.dlist == nil {
		return errors.New("dlist未初始化")
	}
	//修改已存在节点
	//mNode 是list.Element结构
	if mNode, ok := lru.cacheMap[k]; ok {
		lru.dlist.MoveToFront(mNode)
		mNode.Value.(*cacheNode).Value = v
		return nil
	}
	//直接添加新节点 返回值为list.Element结构
	newEle := &list.Element{}
	newEle = lru.dlist.PushFront(&cacheNode{k, v})
	lru.cacheMap[k] = newEle
	//超过容量弹出最后一个（最少使用）
	if lru.Size() > lru.Capacity {
		lastNode := lru.dlist.Back()
		//几乎不会发生；容量是0的时候考虑下
		if lastNode == nil {
			return nil
		}
		cNode := lastNode.Value.(*cacheNode)
		delete(lru.cacheMap, cNode.Key)
		lru.dlist.Remove(lastNode)
	}
	return nil
}

//Get 取值
func (lru *LRUCache) Get(k interface{}) (v interface{}, ret bool, err error) {

	if lru.cacheMap == nil {
		return v, false, errors.New("LRUCache结构体未初始化")
	}

	if ele, ok := lru.cacheMap[k]; ok {
		lru.dlist.MoveToFront(ele)
		return ele.Value.(*cacheNode).Value, true, nil
	}
	return v, false, nil
}

//Remove 移除指定元素
func (lru *LRUCache) Remove(k interface{}) bool {

	if lru.cacheMap == nil {
		return false
	}

	if ele, ok := lru.cacheMap[k]; ok {
		cNode := ele.Value.(*cacheNode)
		delete(lru.cacheMap, cNode.Key)
		lru.dlist.Remove(ele)
		return true
	}
	return false
}
