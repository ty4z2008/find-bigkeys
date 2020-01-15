// Copyright ty4z2008
//
//Licensed under the MIT License
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	logger "github.com/ty4z2008/find-bigkeys/utils"
	"strings"
)

var client *redis.Client

//flag names
const (
	redisUrl      = "url"
	redisPassword = "p"
	redisDB       = "db"
)

// config struct
type config struct {
	url      string
	password string
	db       int
}

//redis config variable

var redisCfg config

const (
	version = "0.0.1"
	//http://patorjk.com/software/taag/#p=display&f=Small%20Slant&t=Redis%20tools
	banner = `
   ___         ___       __            __  
  / _ \___ ___/ (_)__   / /____  ___  / /__
 / , _/ -_) _  / (_-<  / __/ _ \/ _ \/ (_-<
/_/|_|\__/\_,_/_/___/  \__/\___/\___/_/___/ %s

	`
)

type redisInfo struct{}

func main() {
	Init()
	redisInfo, err := GetInfo()
	if err != nil {
		logger.Info("Error:", err)
		return
	}
	fmt.Println("#############################\n")
	fmt.Println("System information\n")
	fmt.Println("Used memory:", redisInfo["used_memory_human"])
	fmt.Println("Total memory:", redisInfo["total_system_memory_human"])
	fmt.Println("#############################\n")
	fmt.Println("Start scan redis server\n")

	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = client.Scan(cursor, "*", 100).Result()
		if err != nil {
			panic(err)
		}
		//if cursor is zero,it's don't need goroutine
		if cursor == 0 {
			findBigKeyNormal(keys)
			break
		} else {
			go findBigKeyNormal(keys)
		}
	}
}

//print redis tools slogan
func slogan() {
	fmt.Println(fmt.Sprintf(banner, version))
}

//parse command-line flag
func parseFlag() {
	flag.StringVar(&redisCfg.url, redisUrl, "localhost:6379", "redis server url")
	flag.StringVar(&redisCfg.password, redisPassword, "", "redis password (default is empty)")
	flag.IntVar(&redisCfg.db, redisDB, 0, "redis db (default 0)")
	flag.Usage = func() {
		fmt.Printf("Usage: \n")
		flag.PrintDefaults()
	}
	//After all flags are defined, call
	flag.Parse()
}

//Init connection
func Init() {
	slogan()
	parseFlag()
	client = redis.NewClient(&redis.Options{
		Addr:     redisCfg.url,
		Password: redisCfg.password, // no password set
		DB:       redisCfg.db,       // use default DB
	})

}

//Fetch redis server info
func GetInfo() (map[string]interface{}, error) {
	redisInfo, err := client.Info().Result()
	if err != nil {
		return nil, err
	}
	info := make(map[string]interface{})
	//line by line read info
	scanner := bufio.NewScanner(strings.NewReader(redisInfo))
	for scanner.Scan() {
		infoItem := strings.Split(scanner.Text(), ":")
		if len(infoItem) > 1 {
			key := infoItem[0]
			value := infoItem[1]
			info[key] = value
		}
	}

	return info, nil
}

func findBigKeyNormal(keys []string) error {
	for i := 0; i < len(keys); i++ {
		checkBigKey(keys[i])
	}
	return nil
}
func findBigKeySharding(keys []string) []string {
	return []string{"test"}
}

//Check if it is big key
func checkBigKey(key string) (bool, error) {
	var isBigKey bool = false
	var length int64 = 0
	var keyType string
	keyType, _ = client.Type(key).Result()
	if keyType == "string" {
		length = client.StrLen(key).Val()
	} else if keyType == "hash" {
		length = client.HLen(key).Val()
	} else if keyType == "list" {
		length = client.LLen(key).Val()
	} else if keyType == "set" {
		length = client.SCard(key).Val()
	} else if keyType == "zset" {
		length = client.ZCard(key).Val()
	} else if keyType == "stream" {
		length = client.XLen(key).Val()
	} else {
		logger.Info("can't detect type for key")
	}
	if length > 1024 {
		isBigKey = true
		ttl, _ := client.TTL(key).Result()
		fmt.Printf("%s %s %d %s\n", key, keyType, length, ttl)
	}

	return isBigKey, nil
}
