package server

import (
	"github.com/lefuturiste/keyvaluer/pkg/commands"
)

var commandMap = map[string]interface{}{
	"PING":      commands.Ping,
	"GET":       commands.Get,
	"SET":       commands.Set,
	"DEL":       commands.Del,
	"EXISTS":    commands.Exists,
	"APPEND":    commands.Append,
	"INCR":      commands.Incr,
	"INCRBY":    commands.IncrBy,
	"KEYS":      commands.Keys,
	"DBSIZE":    commands.DbSize,
	"FLUSHALL":  commands.FlushAll,
	"QUIT":      commands.Quit,
	"COMMAND":   commands.Command,
	"SADD":      commands.SAdd,
	"SISMEMBER": commands.SIsMember,
	"SMEMBERS":  commands.SMembers,
	"SREM":      commands.SRem,
	"LPOP":      commands.LPop,
	"RPUSH":     commands.RPush,
	"SELECT":    commands.Select,
	"EXPIRE":    commands.Expire,
	"ECHO":      commands.Echo,
}
