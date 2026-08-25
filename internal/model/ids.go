package model

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

var seqCounter atomic.Uint64

// NewID 生成带业务前缀的唯一 ID：前缀-时间戳基36-随机4字节。
func NewID(prefix string) string {
	seq := seqCounter.Add(1)
	buf := make([]byte, 4)
	_, _ = rand.Read(buf)
	return fmt.Sprintf("%s-%s-%04x-%s",
		prefix,
		fmt.Sprintf("%x", time.Now().UnixNano()/1_000_000%0x1000000),
		seq&0xffff,
		hex.EncodeToString(buf),
	)
}

// IDPrefixes 各实体的 ID 前缀约定。
const (
	PrefixSpec   = "spec"
	PrefixBatch  = "batch"
	PrefixMold   = "mold"
	PrefixHeat   = "heat"
	PrefixPour   = "pour"
	PrefixCurve  = "curve"
	PrefixBell   = "bell"
	PrefixDefect = "def"
	PrefixAlert  = "alert"
)
