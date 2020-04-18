package server

import (
	"github.com/lefuturiste/keyvaluer/pkg/commands"
)

var commandMap = map[string]interface{}{
	// connexion
	"PING":   commands.Ping,
	"QUIT":   commands.Quit,
	"ECHO":   commands.Echo,
	"SELECT": commands.Select,

	// strings
	"GET":    commands.Get,
	"SET":    commands.Set,
	"APPEND": commands.Append,
	"INCR":   commands.Incr,
	"INCRBY": commands.IncrBy,

	// keys
	"EXISTS": commands.Exists,
	"KEYS":   commands.Keys,
	"DEL":    commands.Del,
	"EXPIRE": commands.Expire,

	// server
	"DBSIZE":   commands.DbSize,
	"FLUSHALL": commands.FlushAll,
	"COMMAND":  commands.Command,

	// sets
	"SADD":      commands.SAdd,
	"SISMEMBER": commands.SIsMember,
	"SMEMBERS":  commands.SMembers,
	"SREM":      commands.SRem,

	// lists
	"LPOP":  commands.LPop,
	"LPUSH": commands.LPush,
	"RPUSH": commands.RPush,
}
