package db

import "time"

const (
	MaxClientIDKey = "max_client_id_{%d}_%d_%d" //max_client_id_{'slot'}_'connid'_'clientid'
	LastMessageKey = "last_message_{%d}_%d"     //last_message_{'slot'}_'connID'
	TTL7Day        = 7 * 24 * time.Hour         //一周过期时间
)
