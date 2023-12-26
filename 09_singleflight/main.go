/*
* @Author: cc
* @Date:   2023/12/26
* @Description: 使用SingleFlight解决缓存击穿，当大量的请求同时查询一个key时，只执行其中一个请求，其他请求共享结果
 */
package main

import (
	"errors"
	"golang.org/x/sync/singleflight"
	"log"
	"sync"
)

var errorNotExist = errors.New("not exist")
var gsf singleflight.Group

func getData(key string) (string, error) {
	data, err := getDataFromCache(key)
	// 从缓存中未获取到数据,从db中获取
	if err == errorNotExist {
		//data, err = getDataFromDB(key)
		// 使用singleflight
		v, err, _ := gsf.Do(key, func() (interface{}, error) {
			return getDataFromDB(key)
		})
		if err != nil {
			log.Println(err)
			return "", err
		}
		data = v.(string)
	} else if err != nil {
		return "", err
	}
	return data, nil
}

// 模拟从cache中获取值，cache中无该值
func getDataFromCache(key string) (string, error) {
	return "", errorNotExist
}

// 模拟从数据库中获取值
func getDataFromDB(key string) (string, error) {
	log.Printf("get %s from database", key)
	return "data", nil
}

func main() {
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			data, err := getData("key")
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(data)
		}()
	}
	wg.Wait()
}
