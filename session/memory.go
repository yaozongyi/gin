package session

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type MemorySessionData struct {
	ID string
	Data map[string]interface{}
	rwLock sync.RWMutex //读写锁，用于读多写少的情况，读锁可以重复的加，写锁互斥
}

// 管理全局的Session
type MemoryMgr struct {
	Session map[string]SessionData // 存储所有SessionData的切片
	rwLock sync.RWMutex //读写锁，用于读多写少的情况，读锁可以重复的加，写锁互斥
}

//内存版初始化session仓库
func NewMemory() (Mgr) {
	return &MemoryMgr{
		Session: make(map[string]SessionData,1024),
	}
}
func (ms *MemorySessionData) GetId() string {
	return ms.ID
}

func (ms *MemorySessionData) GetKey(key string) (value interface{}, err error) {
 	// 获取锁
	ms.rwLock.Lock()
	defer ms.rwLock.Unlock()
	value, ok := ms.Data[key]
	if !ok {
		err = fmt.Errorf("key无效")
		return
	}
	return
}
func (ms *MemorySessionData) SetKey(key string, value interface{})  {
	// 获取锁
	ms.rwLock.Lock()
	defer ms.rwLock.Unlock()
	ms.Data[key] = value
}
func (ms *MemorySessionData) DelKey(key string)  {
	// 获取锁
	ms.rwLock.Lock()
	defer ms.rwLock.Unlock()
	delete(ms.Data, key)
}
// 被动设置是为了redis
func (ms *MemorySessionData) Save() {
	return
}

func (mgr *MemoryMgr) Init(addr string, option ...string) {
	//这里创建Init方法纯属妥协，其实memory版的并不需要初始化，前面NewMemory已经把活干完了
	//这里只是为了满足接口的定义，因为redis里需要这个方法取去连接数据库
	return
}
//GetSessionData 根据传进来的SessionID找到对应Session
func (mgr *MemoryMgr) GetSessionData(sessionId string) (sd SessionData, err error) {
	//
	mgr.rwLock.RLock()
	defer mgr.rwLock.RUnlock()
	sd , ok := mgr.Session[sessionId]
	if !ok {
		err = fmt.Errorf("无效的session")
		return
	}
	return
}
func (mgr *MemoryMgr) createSession()(sd SessionData){
	// 创建uuid
	uuidObj, _ := uuid.NewUUID()
	// 创建SessionData
	sd = NewMemorySessionData(uuidObj.String())
	// 创建对应关系
	mgr.Session[sd.GetId()] = sd
	return sd
}
func NewMemorySessionData(id string) SessionData {
	return &MemorySessionData{ID: id, Data: make(map[string]interface{},8)}
}