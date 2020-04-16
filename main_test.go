package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_PORT = "4545"
)

func getClient(password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:" + TEST_PORT,
		Password: password,
		DB:       0,
	})
}

func TestPing(t *testing.T) {
	startServer(TEST_PORT, "")

	time.Sleep(200 * time.Millisecond)

	client := getClient("")
	pong, err := client.Ping().Result()
	assert.Equal(t, "PONG", pong)
	assert.Nil(t, err)
}

func TestAuth(t *testing.T) {
	os.Setenv("REQUIRED_PASS", "root")
	client := getClient("")

	_, err := client.Ping().Result()
	assert.Contains(t, err, "NOAUTH")

	client = getClient("dd")

	pong, err := client.Ping().Result()
	assert.Equal(t, "", pong)
	assert.Equal(t, err.Error(), "ERR invalid password")

	client = getClient("root")
	pong, err = client.Ping().Result()
	assert.Equal(t, "PONG", pong)
	assert.Nil(t, err)

	os.Setenv("REQUIRED_PASS", "")
}

func TestKey(t *testing.T) {

	client := getClient("")
	var data map[string]string = make(map[string]string)
	data["hello"] = "world"
	data["foo"] = "bar"

	jsonStr, _ := json.Marshal(data)
	setCmd := client.Set("key", jsonStr, -1)
	assert.Nil(t, setCmd.Err())

	getCmd := client.Get("key")
	assert.Nil(t, getCmd.Err())
	result, err := getCmd.Result()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, string(jsonStr), result)

	existsResult, _ := client.Exists("key").Result()
	assert.Equal(t, int64(1), existsResult)

	delCmd := client.Del("key")
	assert.Nil(t, delCmd.Err())
	intResult, err := delCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), intResult)

	existsResult, _ = client.Exists("key").Result()
	assert.Equal(t, int64(0), existsResult)

	flushResult, _ := client.FlushAll().Result()
	assert.Equal(t, "OK", flushResult)

	client.Set("data", "yes_this_is_data", -1)
	client.Set("data_1", "yes_this_is_data", -1)
	client.Set("data_2", "yes_this_is_data", -1)
	client.Set("data_3", "yes_this_is_data", -1)

	keysResult, _ := client.Keys("*").Result()
	var keyArray []string = []string{"data", "data_1", "data_2", "data_3"}

	assert.Equal(t, 4, len(keysResult))
	for _, value := range keysResult {
		assert.Contains(t, keyArray, value)
	}

	client.FlushAll().Result()
	client.Set("data", "0", -1)
	client.Incr("data")
	parsedValue, _ := client.Get("data").Result()
	assert.Equal(t, parsedValue, "1")
	client.IncrBy("data", 5)
	parsedValue, _ = client.Get("data").Result()
	assert.Equal(t, parsedValue, "6")
}
