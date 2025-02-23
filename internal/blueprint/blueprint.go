package blueprint

import (
	"fmt"
)

type RedisStruct struct {
	Network  string
	Addr     string
	Username string
	Password string
	DB       int
}

// return true if external redis is not configured
func (r *RedisStruct) IsNone() bool {
	return r.Network == ""
}

// built redis url from struct
func (r *RedisStruct) String() string {
	return fmt.Sprintf("%s://%s:%s@%s/%v", r.Network, r.Username, r.Password, r.Addr, r.DB)
}
