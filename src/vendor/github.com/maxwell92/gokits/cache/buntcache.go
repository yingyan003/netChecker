package cache

import (
	"github.com/maxwell92/gokits/hashmap"
	"github.com/tidwall/buntdb"
	"github.com/tidwall/gjson"
	"strings"
	"sync"
)

type BuntCache struct {
	sync.Mutex
	DB      *buntdb.DB
	Hash    hashmap.HashMap
	Indexes []string
}

var binstance *BuntCache
var bonce sync.Once

func BuntCacheInstance() *BuntCache {
	bonce.Do(func() {
		binstance = new(BuntCache)
		binstance.Indexes = make([]string, 0)
		binstance.Hash.Map = make(map[string]bool)
	})
	return binstance
}

func (c *BuntCache) Init() error {
	return c.Create(BuntCacheName)
}

// Create
func (c *BuntCache) Create(db string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	bunt, err := buntdb.Open(db, false)
	if err != nil {
		return err
	}
	c.DB = bunt
	return nil
}

// Make Indexes
func (c *BuntCache) Index(columns []string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	for _, index := range columns {
		c.Indexes = append(c.Indexes, index)
		err := c.DB.CreateIndex(index, "*", buntdb.IndexJSON(index))
		if err != nil {
			return err
		}
	}
	return nil
}

// Update
func (c *BuntCache) Update(key, value string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	err := c.DB.Update(func(tx *buntdb.Tx) error {
		c.Hash.AddIfNoExist(key)
		_, _, err := tx.Set(key, value, nil)
		return err
	})

	if err != nil {
		return err
	}
	return nil
}

// Delete
func (c *BuntCache) Delete(key string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Hash.Delete(key)
	err := c.DB.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// Search
func (c *BuntCache) Search(keyword string) (*map[string]string, error) {
	results := make(map[string]string)

	// Splite keys from keyword inputed from user
	keys := []string{}
	for _, key := range strings.Split(keyword, Sep) {
		key = strings.Trim(key, " \n")
		if !strings.EqualFold(key, "") {
			keys = append(keys, key)
		}
	}

	// Search in all index
	for _, index := range c.Indexes {
		for _, key := range keys {
			err := c.doSearch(index, key, &results)
			if err != nil {
				return &results, err
			}

		}
	}
	return &results, nil
}

// doSearch
func (c *BuntCache) doSearch(index, k string, results *map[string]string) error {

	err := c.DB.View(func(tx *buntdb.Tx) error {
		tx.Ascend(index, func(key, value string) bool {
			content := gjson.Get(value, index)
			if ok := strings.Contains(content.String(), k); ok {
				_, exists := (*results)[key]
				if !exists {
					(*results)[key] = value
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// Find
func (c *BuntCache) Find(args ...Indexer) (*map[string]string, error) {
	results := make(map[string]string)
	counters := make(map[string]int32)
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	err := c.DB.View(func(tx *buntdb.Tx) error {
		for _, indexer := range args {
			tx.Ascend(indexer.Index, func(key, value string) bool {
				v := gjson.Get(value, indexer.Index)
				// Equal with the given value of the "Index"
				if 0 == strings.Compare(indexer.Key, v.String()) {
					count := counters[key]
					counters[key] = count + 1
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return &results, err
	}

	// Range the Counters
	len := len(args)
	for key, count := range counters {
		if int32(len) == count {
			c.DB.View(func(tx *buntdb.Tx) error {
				value, _ := tx.Get(key)
				results[key] = value
				return nil
			})
		}
	}
	return &results, nil
}

// Watch
func (c *BuntCache) Watch(name string) error {
	return nil
}
