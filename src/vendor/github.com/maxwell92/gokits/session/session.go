package session

import (
	"encoding/json"
	"errors"
	mylog "github.com/maxwell92/gokits/log"

	localtime "github.com/maxwell92/gokits/time"

	redis "github.com/garyburd/redigo/redis"
	"github.com/kataras/iris"
	"strings"
	"sync"
)

var log = mylog.Log

const (
	DEFAULT_EXPIRATION = "86400" // 24*60*60s
)

type Session struct {
	SessionId  string   `json:"sessionId"`
	UserId     string   `json:"userId"`
	UserName   string   `json:"userName"`
	OrgId      string   `json:"orgId"`
	DcList     []string `json:"dcList"`
	CreatedAt  string   `json:"CreatedAt"`
	Expiration string   `json:"Expiration"` // expiration in seconds
}

type UserSessionMap struct {
	Mapping map[string]string `json:"mapping"`
}

type USMap interface {
	SetSession(string, string)
	DelSession(string)
	GetSessionId(string) string
}

func (ss *SessionStore) Validate(ctx *iris.Context) {
	orgId := ctx.RequestHeader("OrgId")
	// userId := ctx.RequestHeader("UserId")
	sessionId := ctx.RequestHeader("Authorization")

	log.Infof("sessionId=%s, orgId=%s", sessionId, orgId)

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if !ok || err != nil {
		ctx.Write("{\"code\": 1401, \"message\":\"请重新登录\"}")
	} else {
		ctx.Next()
	}
}

func (ss *SessionStore) ExistsSessionId(sessionId string) (*Session, bool) {
	session, err := ss.Get(sessionId)
	if err == nil && session != nil {
		log.Infof("Session cached: sessionId=%s, userId=%s, orgId=%s", sessionId, session.UserId, session.OrgId)
		return session, true
	}
	log.Errorf("Session is expired!")
	return nil, false
}

func (usm *UserSessionMap) SetSession(userId, sessionId string) {
	if usm.Mapping == nil {
		usm.Mapping = make(map[string]string)
	}

	usm.Mapping[userId] = sessionId
	ss := SessionStoreInstance()
	conn := ss.pool.Get()
	if conn == nil {
		log.Fatalln("The Connection is nil: conn := ss.pool.Get()")
	}
	defer conn.Close()

	data := usm.Encode()
	redis.String(conn.Do("SET", "userSessionMap", data))
}

func (usm *UserSessionMap) DelSession(userId string) {
	ss := SessionStoreInstance()
	conn := ss.pool.Get()
	if conn == nil {
		log.Fatalln("The Connection is nil: conn := ss.pool.Get()")
	}
	defer conn.Close()

	if usm.Mapping != nil {
		delete(usm.Mapping, userId)
	}

	data := usm.Encode()
	redis.String(conn.Do("SET", "userSessionMap", data))
}

func (usm *UserSessionMap) GetSessionId(userId string) string {
	return GetSessionIdByUserId(userId)
}

func GetSessionIdByUserId(userId string) string {
	if userId == "" {
		log.Errorf("UserSessionMap GetSessionIdByUserId Error: userId=%s", userId)
		return ""
	}

	ss := SessionStoreInstance()
	conn := ss.pool.Get()
	if conn == nil {
		log.Fatalln("The Connection is nil: conn := ss.pool.Get()")
	}
	defer conn.Close()

	data, err := redis.String(conn.Do("GET", "userSessionMap"))
	if err != nil {
		log.Errorf("UserSessionMap Get Error")
		return ""
	}
	usm := Decode(data)
	return usm.Mapping[userId]
}

func UserSessionMapInstance() USMap {
	ss := SessionStoreInstance()
	conn := ss.pool.Get()
	if conn == nil {
		log.Fatalln("The Connection is nil: conn := ss.pool.Get()")
	}
	defer conn.Close()

	data, err := redis.String(conn.Do("GET", "userSessionMap"))
	if err != nil {
		log.Errorf("UserSessionMap Get Error: err=%s", err)
	}

	if strings.EqualFold(data, "") {
		return &UserSessionMap{}
	}

	usm := Decode(data)
	return usm
}

func NewUserSessionMap() {
	ss := SessionStoreInstance()
	conn := ss.pool.Get()
	if conn == nil {
		log.Fatalln("The Connection is nil: conn := ss.pool.Get()")
	}
	defer conn.Close()

	data, err := redis.String(conn.Do("GET", "userSessionMap"))
	if err != nil {
		log.Errorf("UserSessionMap Get Error")
		conn.Do("SET", "userSessionMap", "\"\":\"\"")
	}

	if strings.EqualFold(data, "") {
		log.Infof("UserSessionMap Get Empty")
		conn.Do("SET", "userSessionMap", "\"\":\"\"")
	}
	log.Tracef("UserSessionMap Open Success")
}

func FakeUserSessionMap() USMap {
	return nil
}

func (usm *UserSessionMap) Encode() string {
	data, err := json.Marshal(usm)
	if err != nil {
		log.Errorf("UserSessionMap Encode Error: err=%s", err)
		return ""
	}

	return string(data)
}

func Decode(data string) *UserSessionMap {
	usm := new(UserSessionMap)
	usm.Mapping = make(map[string]string)
	err := json.Unmarshal([]byte(data), usm)
	if err != nil {
		log.Errorf("UserSessionMap Decode Error: err=%s, data=%s", err, data)
		return &UserSessionMap{}
	}

	return usm
}

func NewSession(token, userId, userName, orgId string) *Session {
	return &Session{
		SessionId:  token,
		UserId:     userId,
		UserName:   userName,
		OrgId:      orgId,
		CreatedAt:  localtime.NewLocalTime().String(),
		Expiration: DEFAULT_EXPIRATION,
	}
}

func (s *Session) DecodeJson(data string) error {
	err := json.Unmarshal([]byte(data), s)

	if err != nil {
		log.Errorf("Session DecodeJson Error: err=%s", err)
		return err
	}

	return nil
}

func (s *Session) EncodeJson() (string, error) {
	data, err := json.Marshal(s)
	if err != nil {
		log.Errorf("Session EncodeJson Error: err=%s", err)
		return "", err
	}
	return string(data), nil
}

var instance *SessionStore

var once sync.Once

type SessionStore struct {
	pool *redis.Pool
}

func SessionStoreInstance() *SessionStore {
	log.Tracef("SessionStoreInstance Success")
	return instance
}

func NewSessionStore(p *redis.Pool) *SessionStore {
	once.Do(func() {
		instance = &SessionStore{
			pool: p,
		}
		log.Tracef("SessionStore Pool=%p", instance.pool)
	})
	return instance
}

func FakeSessionStore() *SessionStore {
	once.Do(func() {

	})
	return instance
}

func (ss *SessionStore) ValidateOrgId(sessionIdClient string, OrgIdClient string) (bool, error) {
	session, err := ss.Get(sessionIdClient)
	if err != nil {
		log.Errorf("Get session from sessionIdClient error: sessionIdClient: %s, err=%s", sessionIdClient, err)
		return false, err
	}

	// sessionId invalid
	if session == nil && err == nil {
		return false, errors.New("Validate sessionIdClient failed: invalid sessionIdClient")
	}

	if session.OrgId == OrgIdClient {
		return true, nil
	} else {
		return false, errors.New("Validate sessionId failed: OrgId doesn't match")
	}

}

func (ss *SessionStore) ValidateUserId(sessionIdClient string, UserIdClient string) (bool, error) {
	session, err := ss.Get(sessionIdClient)
	if err != nil {
		log.Errorf("Get session from sessionIdClient error: sessionIdClient=%s, err=%s", sessionIdClient, err)
		return false, err
	}

	// sessionId invalid
	if session == nil && err == nil {
		return false, errors.New("Validate sessionIdClient failed: invalid sessionIdClient")
	}
	if session.UserId == UserIdClient {
		return true, nil
	} else {
		return false, errors.New("Validate sessionId failed: UserId doesn't match")
	}
}

func (ss *SessionStore) Get(sessionId string) (*Session, error) {
	conn := ss.pool.Get()
	if conn == nil {
		log.Fatalln("The Connection is nil: conn := ss.pool.Get()")
		return nil, errors.New("The Connection is nil: conn := ss.pool.Get()")
	}

	defer conn.Close()

	// If exists
	exists, err := ss.Exist(sessionId)
	if err != nil {
		log.Fatalf("SessionStore exist error: sessionId=%s, err=%s", sessionId, err)
		return nil, err
	}

	// not exists
	if !exists {
		log.Warnf("The Session not exists: sessionId=%s", sessionId)
		return nil, nil
	}

	// exists
	session := &Session{}

	data, err := redis.Bytes(conn.Do("GET", sessionId))
	if err != nil {
		log.Fatalf("Redis Get error: sessionId=%s, err=%s", sessionId, err)
		return nil, err
	}

	err = json.Unmarshal(data, session)

	if err != nil {
		log.Fatalf("Json unmashal failed: data=%s, err=%s", string(data), err)
		return nil, err
	}

	return session, nil
}

func (ss *SessionStore) Set(session *Session) error {
	conn := ss.pool.Get()

	if conn == nil {
		log.Fatalf("The Connection is nil: conn := ss.pool.Get()")
		return errors.New("The Connection is nil: conn := ss.pool.Get()")
	}

	defer conn.Close()

	data, err := json.Marshal(session)

	if err != nil {
		log.Fatalf("Json marshal err: err=%s", err)
		return err
	}

	_, err = conn.Do("SET", session.SessionId, data, "EX", session.Expiration)

	if err != nil {
		log.Errorf("Redis set error: sessionId=%s, err=%s", session.SessionId, err)
		return err
	}

	return nil
}

func (ss *SessionStore) Delete(sessionId string) error {

	conn := ss.pool.Get()

	if conn == nil {
		log.Fatalln("The Connection is nil: conn := ss.pool.Get()")
		return errors.New("The Connection is nil: conn := ss.pool.Get()")
	}

	defer conn.Close()

	_, err := conn.Do("DEL", sessionId)
	if err != nil {
		log.Errorf("Redis delete error: sessionId=%s, err=%s", sessionId, err)
		return err
	}

	return nil
}

func (ss *SessionStore) Exist(sessionId string) (bool, error) {
	conn := ss.pool.Get()

	if conn == nil {
		log.Fatalln("The Connection is nil: conn := ss.pool.Get()")
		return false, errors.New("The Connection is nil: conn := ss.pool.Get()")
	}

	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", sessionId))

	if err != nil {
		log.Fatalf("Redis Bool error: sessionId=%s, err=%s\n", sessionId, err)
		return false, err
	}

	return exists, nil
}
