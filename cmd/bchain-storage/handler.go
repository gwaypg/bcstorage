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
	spaceName  string
	file       string
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
func (h *HttpHandler) AddToken(spaceName, file, token string) {
	h.tokenLk.Lock()
	defer h.tokenLk.Unlock()
	h.token[token] = FileToken{
		spaceName:  spaceName,
		file:       file,
		createTime: time.Now(),
	}
}
func (h *HttpHandler) DelayToken(file string) bool {
	h.tokenLk.Lock()
	defer h.tokenLk.Unlock()
	t, ok := h.token[file]
	if !ok {
		return false
	}
	t.createTime = time.Now()
	h.token[file] = t
	return true
}

func (h *HttpHandler) DeleteToken(token string) {
	h.tokenLk.Lock()
	defer h.tokenLk.Unlock()
	delete(h.token, token)
}

func (h *HttpHandler) VerifyToken(file, token string) (FileToken, bool) {
	h.tokenLk.Lock()
	defer h.tokenLk.Unlock()
	h.gcToken()

	t, ok := h.token[token]
	if !ok {
		return FileToken{}, false
	}
	if t.file != file {
		return FileToken{}, false
	}
	return t, true

}
