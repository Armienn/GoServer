package goserver

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type SessionManager struct {
	cookieName  string     //private cookiename
	lock        sync.Mutex // protects session
	maxlifetime int64
	Sessions    map[string]Session
}

func NewSessionManager(cookieName string, maxlifetime int64) *SessionManager {
	return &SessionManager{cookieName: cookieName, maxlifetime: maxlifetime, Sessions: make(map[string]Session)}
}

type Session interface {
	Set(key, value interface{})              //set session value
	Get(key interface{}) (interface{}, bool) //get session value
	Delete(key interface{})                  //delete session value
	SessionID() string                       //back current sessionID
	Age() int64
}

func newSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *SessionManager) SessionStart(w http.ResponseWriter, r *http.Request) Session {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := newSessionId()
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
		http.SetCookie(w, &cookie)
		return manager.sessionInit(sid)
	}
	sid, _ := url.QueryUnescape(cookie.Value)
	if session, ok := manager.Sessions[sid]; ok {
		return session
	}
	return manager.sessionInit(sid)
}

func (manager *SessionManager) sessionInit(sid string) Session {
	session := &SessionData{sid: sid, timeAccessed: time.Now(), value: make(map[interface{}]interface{}, 0)}
	manager.Sessions[sid] = session
	return session
}

func (manager *SessionManager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	manager.lock.Lock()
	defer manager.lock.Unlock()
	sid, _ := url.QueryUnescape(cookie.Value)
	if _, ok := manager.Sessions[sid]; ok {
		delete(manager.Sessions, sid)
	}
	cookie = &http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: time.Now(), MaxAge: -1}
	http.SetCookie(w, cookie)
}

func (manager *SessionManager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	for key, value := range manager.Sessions {
		if value.Age() > manager.maxlifetime {
			delete(manager.Sessions, key)
		}
	}
	time.AfterFunc(time.Duration(manager.maxlifetime), func() { manager.GC() })
}
