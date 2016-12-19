package goserver

import "time"

type SessionData struct {
	sid          string                      // unique session id
	timeAccessed time.Time                   // last access time
	value        map[interface{}]interface{} // session value stored inside
}

func (session *SessionData) Set(key, value interface{}) {
	session.value[key] = value
	session.timeAccessed = time.Now()
}

func (session *SessionData) Get(key interface{}) (interface{}, bool) {
	session.timeAccessed = time.Now()
	if v, ok := session.value[key]; ok {
		return v, true
	}
	return nil, false
}

func (session *SessionData) Delete(key interface{}) {
	delete(session.value, key)
	session.timeAccessed = time.Now()
}

func (session *SessionData) SessionID() string {
	return session.sid
}

func (session *SessionData) Age() int64 {
	return time.Now().Unix() - session.timeAccessed.Unix()
}
