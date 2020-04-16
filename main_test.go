package main

import (
	"encoding/json"
	"testing"

	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
)

func getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:4545",
		Password: "",
		DB:       0,
	})
}

func TestPing(t *testing.T) {
	client := getClient()
	pong, err := client.Ping().Result()
	assert.Equal(t, "PONG", pong)
	assert.Nil(t, err)
}

func TestKey(t *testing.T) {
	client := getClient()
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
	delCmd := client.Del("key")
	assert.Nil(t, delCmd.Err())
	intResult, err := delCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), intResult)
}
