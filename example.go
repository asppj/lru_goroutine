package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"sync"
	"time"
)

func p(id int, i int) string {
	now := time.Now()
	return fmt.Sprintf("ID:%2d第%v次循环：时间： %v--%v", id, i, now.Format("2006/1/2 15:04:05"), now.UnixNano())
}

const Cap int = 50
const RandMax int64 = 100
const Count int = 10000

//命中率约为Cap/RandMax
const GoroutineNum = 200

func main() {
	fmt.Printf("缓存容量%d个；随机范围0-%d；协程个数:%d个；单个协程循环次数：%d次;容量/净含量\n",
		Cap, RandMax, GoroutineNum, Count)
	lru, cFn := NewGoLRUCache(Cap, nil)
	group := &sync.WaitGroup{}
	for i := 0; i < GoroutineNum; i++ {
		group.Add(1)
		go goP(i, lru, group)
	}
	group.Wait()
	cFn()
}

func goP(id int, lru *GoLRUCache, group *sync.WaitGroup) {
	mathed, noMatched := 0, 0
	for i := 1; i < 500; i++ {
		//rand.Seed(time.Now().UnixNano())
		//k := rand.Intn(RandMax)
		ran, _ := rand.Int(rand.Reader, big.NewInt(RandMax))
		ranStr := fmt.Sprintf("%v", ran)
		k, _ := strconv.Atoi(ranStr)
		//fmt.Printf("k=%3d ", k)
		value, ret, err := lru.Get(k)
		if err != nil {
			log.Fatal(errors.New(fmt.Sprintf("发生错误：%v", err)))
		}
		if ret {
			//fmt.Printf("k=%4d使用缓存：%v （%v）\n",k, value, p(id,i))
			mathed++
		} else {
			value = p(id, i)
			if err := lru.Set(k, value); err == nil {
				//fmt.Printf("k=%4d新缓存：%v \n",k, value)
			}
			//fmt.Printf("不能使用缓存：%v \n", value)
			noMatched++

		}
		//time.Sleep(time.Microsecond * 100)
	}
	group.Done()
	fmt.Printf("GoroutineID:%d 命中缓存:%d次； 未命中缓存:%d次;  命中率：%%%v \n", id, mathed, noMatched, mathed*100/(mathed+noMatched))
}
