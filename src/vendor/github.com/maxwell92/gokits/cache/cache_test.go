package cache

/*
import (
	"testing"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"app/backend/common/util/hashmap"
	"strconv"
)

func NewBuntCache() *BuntCache{
	bc := &BuntCache{
		Indexes:  make([]string, 0),
		Hash: hashmap.HashMap {
			Map: make(map[string]bool),
		},
	}
	return bc
}

func Test_BuntCache(*testing.T) {

	// Construct test data
	var d = new(Deploy)
	d.Data = make([]DeployInfo, 0)

	// 1. Create BuntCache
	cache := NewBuntCache()
	// cache := BuntCacheInstance()
	cache.Create("test")

	// 2. Make indexes
	cache.Index([]string{UserName, DcName, DcId, DeploymentName, OrgName})

	// 3. Load && Update
	data, err := ioutil.ReadFile("t.json")
	if err != nil {
		fmt.Printf("ReadFile error: %s\n", err)
		return
	}

	d.Unmarshal(&data)

	for _, value := range d.Data {
		content, _ := json.Marshal(&value)
		id := strconv.Itoa(int(value.DcID))
		key := id + ":" + value.OrgName + ":" + value.DeploymentName
		value := string(content)
		cache.Update(key, value)
	}

	fmt.Println("=======================================================")
	// 4. Search
	results, _ := cache.Search("liyao.miao 办公网")
	for k, value := range *results {
		fmt.Println(k)
		fmt.Printf("Value: %s\n", value)
	}

	// 5. Delete
	// str := `{"userName":"liyao.miao","dcName":"办公网","dcId":1,"deploymentName":"a2","updateTime":"2017-02-2217: 45: 56+0800CST"}`
	str := "26:ops:a2"
	cache.Delete(str)

	// 6. Search
	fmt.Println("\n=======================================================")
	results, _ = cache.Search("liyao.miao 办公网")
	for _, value := range *results {
		fmt.Printf("Value: %s\n", value)
	}

	fmt.Println("\n=======================================================")
	fmt.Println("Search:")
	results, _ = cache.Search("a1 1 办公网 liyao.miao")
	for _, value := range *results {
		fmt.Println(value)
	}
/*
	fmt.Println("Specify:")
	record, err := cache.Specify("1:dev:a1")
	if err != nil {
		fmt.Printf("error=%s\n", err)
	} else {
		fmt.Println(record)
	}

	fmt.Println("List:")
	results, _  = cache.List(1, "dev")
	for k, value := range *results {
		fmt.Println(k)
		fmt.Println(value)
	}
*/
/*
	fmt.Printf("===================================FIND==============================\n\n")
	m1 := make([]Indexer, 1)
	m1[0] = Indexer{"orgName", "dev"}

	rs, _ := cache.Find(m1...)
	for _, value := range *rs {
		fmt.Println(value)
	}

	fmt.Printf("===================================FIND==============================\n\n")
	m2 := make([]Indexer, 2)
	// m2[0] = Indexer{"orgName", "dev"}
	// m2[1] = Indexer{"deploymentName", "a1"}
	m2[0] = Indexer{"dcId", "1"}
	m2[1] = Indexer{"orgName", "ops"}

	rs, _ = cache.Find(m2...)
	for _, value := range *rs {
		fmt.Println(value)
	}

	fmt.Printf("===================================FIND==============================\n\n")

	m3 := make([]Indexer, 1)
	m3[0] = Indexer{"orgName", "ops"}
	//m3[1] = Indexer{"deploymentName", "app-2"}
	//m3[2] = Indexer{"dcId", "1"}

	rs, err = cache.Find(m3...)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range *rs {
		fmt.Println(value)
	}

}
*/
