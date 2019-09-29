package main

import (
	"container/list"
	"errors"
	"time"
)

const ExpireMinute int = 10

//cacheNode 节点
type cacheNode struct {
	Key, Value interface{}
	CreateTime time.Time
}

//newCacheNode 新建节点
func newCacheNode(k, v interface{}) *cacheNode {
	createTime := time.Now()
	return &cacheNode{k, v, createTime}
}

//LRUCache LRU缓存
type LRUCache struct {
	Capacity   int
	linkedList *list.List
	cacheMap   map[interface{}]*list.Element
	Expire     *float64 //过期时间 time.Minute；nil不过期
}

//NewLruCache 新建 LRUCache
func NewLruCache(cap int, expire *float64) *LRUCache {
	return &LRUCache{Capacity: cap, linkedList: list.New(), cacheMap: make(map[interface{}]*list.Element), Expire: expire}
}

//Size 元素个数
func (lru *LRUCache) Size() int {
	return lru.linkedList.Len()
}

//Set 设置缓存
func (lru *LRUCache) Set(k, v interface{}) error {
	//异常
	if lru.linkedList == nil {
		return errors.New("linkedList未初始化")
	}
	//修改已存在节点
	//mNode 是list.Element结构
	if mNode, ok := lru.cacheMap[k]; ok {
		lru.linkedList.MoveToFront(mNode)
		cacheP := mNode.Value.(*cacheNode)
		cacheP.CreateTime = time.Now()
		cacheP.Value = v
		return nil
	}
	//直接添加新节点 返回值为list.Element结构
	newEle := &list.Element{}
	newEle = lru.linkedList.PushFront(newCacheNode(k, v))
	lru.cacheMap[k] = newEle
	//超过容量弹出最后一个（最少使用）
	if lru.Size() > lru.Capacity {
		lastNode := lru.linkedList.Back()
		//几乎不会发生；容量是0的时候考虑下
		if lastNode == nil {
			return nil
		}
		cNode := lastNode.Value.(*cacheNode)
		delete(lru.cacheMap, cNode.Key)
		lru.linkedList.Remove(lastNode)
	}
	return nil
}

//Get 取值
func (lru *LRUCache) Get(k interface{}) (v interface{}, ret bool, err error) {

	if lru.cacheMap == nil {
		return v, false, errors.New("LRUCache结构体未初始化")
	}

	if ele, ok := lru.cacheMap[k]; ok {
		//设置了过期时间
		if lru.Expire != nil {
			createTime := ele.Value.(*cacheNode).CreateTime
			expire := time.Now().Sub(createTime).Minutes()
			if expire > *lru.Expire {
				return nil, false, nil
			}
		}

		lru.linkedList.MoveToFront(ele)
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
		lru.linkedList.Remove(ele)
		return true
	}
	return false
}
