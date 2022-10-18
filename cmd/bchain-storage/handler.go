package main

import (
	"sync"
	"time"
)

type CheckCache struct {
	out        string
	createTime time.Time
}
type FileToken struct {
	space      string
	token      string
	createTime time.Time
}

type HttpHandler struct {
	checkCache   *CheckCache
	checkCacheLk sync.Mutex

	token   map[string]FileToken
	tokenLk sync.Mutex
}

func (h *HttpHandler) gcToken() {
	now := time.Now()
	for key, val := range h.token {
		if now.Sub(val.createTime) > 3*time.Hour {
			delete(h.token, key)
		}
	}
}
func (h *HttpHandler) AddToken(space, sid, token string) {
	h.tokenLk.Lock()
	defer h.tokenLk.Unlock()
	h.token[sid] = FileToken{
		space:      space,
		token:      token,
		createTime: time.Now(),
	}
}
func (h *HttpHandler) DelayToken(sid string) bool {
	h.tokenLk.Lock()
	defer h.tokenLk.Unlock()
	t, ok := h.token[sid]
	if !ok {
		return false
	}
	t.createTime = time.Now()
	h.token[sid] = t
	return true
}

func (h *HttpHandler) DeleteToken(sid string) {
	h.tokenLk.Lock()
	defer h.tokenLk.Unlock()
	delete(h.token, sid)
}

func (h *HttpHandler) VerifyToken(sid, token string) (FileToken, bool) {
	h.tokenLk.Lock()
	defer h.tokenLk.Unlock()
	h.gcToken()

	t, ok := h.token[sid]
	if !ok {
		return FileToken{}, false
	}
	if t.token != token {
		return FileToken{}, false
	}
	return t, true

}
