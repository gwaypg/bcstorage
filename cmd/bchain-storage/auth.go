package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gwaycc/bchain-storage/lib/bcrypt"
	"github.com/gwaylib/errors"
)

const (
	_leveldbPath = "leveldb"
)

// empty of md5 is 'd41d8cd98f00b204e9800998ecf8427e'
const adminUser = "admin"
const adminDefaultPwd = "d41d8cd98f00b204e9800998ecf8427e"

var (
	_userMap = UserMap{
		Auth:  map[string]UserAuth{},
		Space: map[string]UserSpace{},
	}
)

func genPasswd() string {
	token := [16]byte(uuid.New())
	if time.Now().UnixNano()%2 == 0 {
		return fmt.Sprintf("%X", md5.Sum(token[:]))
	}
	return fmt.Sprintf("%x", md5.Sum(token[:]))
}

type UserSpace struct {
	Name    string
	Attr    int32 // TODO
	Size    int64 // TODO
	Used    int64 // TODO
	Private bool  // TODO
}

type UserAuth struct {
	User      string
	Passwd    string
	SpaceName string
}

type UserMap struct {
	lk    sync.Mutex
	Auth  map[string]UserAuth
	Space map[string]UserSpace
}

func (u *UserMap) GetAuth(user string) (UserAuth, bool) {
	u.lk.Lock()
	defer u.lk.Unlock()
	a, ok := u.Auth[user]
	if !ok {
		return UserAuth{}, false
	}
	return a, true
}

func (u *UserMap) AddSpace(space UserSpace) error {
	u.lk.Lock()
	defer u.lk.Unlock()
	u.Space[space.Name] = space
	if err := os.MkdirAll(filepath.Join(_rootPathFlag, space.Name), 0755); err != nil {
		return errors.As(err)
	}
	return nil
}
func (u *UserMap) GetSpace(name string) (UserSpace, bool) {
	u.lk.Lock()
	defer u.lk.Unlock()
	space, ok := u.Space[name]
	if !ok {
		// TODO: fetch from leveldb
		return UserSpace{}, false
	}
	return space, true
}
func (u *UserMap) SpacePath(spaceName string) (string, error) {
	space, ok := u.GetSpace(spaceName)
	if !ok {
		return "", errors.ErrNoData.As(spaceName)
	}
	return filepath.Join(_rootPathFlag, space.Name), nil
}
func (u *UserMap) AddSpaceUsed(name string, val int64) error {
	u.lk.Lock()
	defer u.lk.Unlock()
	space, ok := u.Space[name]
	if !ok {
		return errors.ErrNoData.As(name)
	}
	space.Used += val
	u.Space[name] = space
	return PutLevelDB(fmt.Sprintf(_leveldb_prefix_space, name), &space)
}

func (u *UserMap) UpdateAuth(auth UserAuth) error {
	u.lk.Lock()
	defer u.lk.Unlock()

	u.Auth[auth.User] = auth
	return PutLevelDB(fmt.Sprintf(_leveldb_prefix_user, auth.User), &auth)
}

func initDaemonAuth() {
	userPrefixLen := len("_user.")
	if err := IterLevelDB("_user.", func(key, val []byte) error {
		auth := UserAuth{}
		if err := json.Unmarshal(val, &auth); err != nil {
			return errors.As(err, string(val))
		}
		_userMap.Auth[string(key[userPrefixLen:])] = auth
		return nil
	}); err != nil {
		panic(err)
	}

	spacePrefixLen := len("_space.")
	if err := IterLevelDB("_space.", func(key, val []byte) error {
		space := UserSpace{}
		if err := json.Unmarshal(val, &space); err != nil {
			return errors.As(err, string(val))
		}
		_userMap.Space[string(key[spacePrefixLen:])] = space
		return nil
	}); err != nil {
		panic(err)
	}

	_, ok := _userMap.GetAuth(adminUser)
	if !ok {
		if err := _userMap.UpdateAuth(UserAuth{
			User:   adminUser,
			Passwd: bcrypt.BcryptPwd(adminDefaultPwd),
		}); err != nil {
			panic(err)
		}
	}
}

func validHttpFilePath(spaceName, file string) (string, bool) {
	space, ok := _userMap.GetSpace(spaceName)
	if !ok {
		return "", false
	}
	rootPath := _rootPathFlag
	absPath, err := filepath.Abs(filepath.Join(rootPath, space.Name, file))
	if err != nil {
		return "", false
	}
	if !strings.HasPrefix(absPath, filepath.Join(rootPath, space.Name)) {
		return "", false
	}
	return absPath, true
}

func authWrite(r *http.Request) (FileToken, bool) {
	username, passwd, ok := r.BasicAuth()
	if !ok {
		log.Infof("no BasicAuth:%s", r.RemoteAddr)
		return FileToken{}, false
	}
	fAuth, ok := _handler.VerifyToken(username, passwd)
	if !ok {
		log.Infof("VerifyToken failed:%s", r.RemoteAddr)
		return FileToken{}, false
	}
	file := r.FormValue("file")
	if _, ok := validHttpFilePath(fAuth.spaceName, file); !ok {
		log.Infof("VerifyRW failed:%s", r.RemoteAddr)
		return FileToken{}, false
	}
	return fAuth, true
}

func authAdmin(r *http.Request) (UserAuth, bool) {
	// auth
	username, passwd, ok := r.BasicAuth()
	if !ok {
		log.Infof("auth failed:%s, no auth", r.RemoteAddr)
		return UserAuth{}, false
	}
	auth, ok := _userMap.GetAuth(username)
	if !ok {
		log.Infof("auth user failed:%s,%s,%s", r.RemoteAddr, username, passwd)
		return UserAuth{}, false
	}
	if !bcrypt.BcryptMatch(passwd, auth.Passwd) {
		// TODO: limit the failed
		log.Infof("auth passwd failed:%s,%s,%s", r.RemoteAddr, username, passwd)
		return UserAuth{}, false
	}
	return auth, true
}
