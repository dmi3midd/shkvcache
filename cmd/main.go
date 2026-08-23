package main

import "github.com/dmi3midd/shkvcache"

func main() {
	c, err := shkvcache.NewCache[string](8)
	if err != nil {
		return
	}
	c.Set("key1", "val1", 0)
	val, ok := c.Get("key1")
	if ok {
		println(val)
	}

}
