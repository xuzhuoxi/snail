//
//Created by xuzhuoxi
//on 2019-02-10.
//@author xuzhuoxi
//
package imodule

import (
	"sync"
	"time"
)

//统计间隔
const DefaultStatsInterval = int64(5 * time.Minute)

type IServiceState interface {
	GetName() string
	GetWeight() float64
}

type IServiceStateDetail interface {
	//运行时间
	GetPassNano() int64

	//当前统计的服务权重(连接数*统计时间/统计响应时间)
	//越大代表压力越大
	StatsWeight() float64
	//统计时间段的请求密度(次数/秒)
	StatsReqDensity() int
	//统计时间段的响应密度(统计响应时间/统计时间)
	StatsRespDensity() float64
}

func NewServiceState(name string, weight float64) *ServiceState {
	return &ServiceState{Name: name, Weight: weight}
}

func NewServiceStateDetail(name string, statsInterval int64) *ServiceStateDetail {
	return &ServiceStateDetail{Name: name, StatsInterval: statsInterval}
}

//------------------------------

type ServiceState struct {
	//名称
	Name string
	//压力
	Weight float64
}

func (ss *ServiceState) GetName() string {
	return ss.Name
}

func (ss *ServiceState) GetWeight() float64 {
	return ss.Weight
}

//------------------------------

type ServiceStateDetail struct {
	Name string
	//启动时间戳(纳秒)
	StartTimestamp int64
	//最大连接数
	MaxLinkCount uint64
	//总请求数
	TotalReqCount int64
	//总响应时间
	TotalRespTime int64
	//最大响应时间(纳秒)
	MaxRespTime int64

	//连接数
	LinkCount uint64

	//统计开始时间戳(纳秒)
	StatsTimestamp int64
	//统计请求数
	StatsReqCount int64
	//统计响应时间(纳称)
	StatsRespUnixNano int64
	//统计间隔
	StatsInterval int64

	lock sync.RWMutex
}

//启动时间
func (s ServiceStateDetail) GetPassNano() int64 {
	return s.getStatsPass()
}

//当前统计的服务权重(连接数*统计时间/统计响应时间)
//越大代表压力越大
func (s ServiceStateDetail) StatsWeight() float64 {
	if 0 == s.StatsRespUnixNano {
		return float64(s.LinkCount)
	} else {
		pass := s.getStatsPass()
		return float64(s.LinkCount) * float64(pass) / float64(s.StatsRespUnixNano)
	}
}

//统计时间段的请求密度(次数/秒)
func (s ServiceStateDetail) StatsReqDensity() int {
	pass := s.getStatsPass()
	return int(int64(time.Second) * s.StatsReqCount / pass)
}

//统计时间段的响应密度(统计响应时间/统计时间)
func (s ServiceStateDetail) StatsRespDensity() float64 {
	pass := time.Now().UnixNano() - s.StartTimestamp
	return float64(s.StatsRespUnixNano) / float64(pass)
}

//-------------------------------

//启动
func (s *ServiceStateDetail) Start() {
	s.lock.Lock()
	defer s.lock.Unlock()
	now := time.Now().UnixNano()
	s.StartTimestamp = now
	s.StatsTimestamp = now
}

//增加一个连接
func (s *ServiceStateDetail) AddLinkCount() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.LinkCount++
	if s.LinkCount > s.MaxLinkCount { //更新最大连接数
		s.MaxLinkCount = s.LinkCount
	}
}

//减少一个连接
func (s *ServiceStateDetail) RemoveLinkCount() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.LinkCount--
}

//----------------------

//增加一个请求
func (s *ServiceStateDetail) AddReqCount() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.statsFull() {
		s.statsReset()
	}
	s.TotalReqCount++
	s.StatsReqCount++
}

//增加响应时间量
func (s *ServiceStateDetail) AddRespUnixNano(unixNano int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.StatsRespUnixNano += unixNano
	s.TotalRespTime += unixNano
	if unixNano > s.MaxRespTime { //更新最大响应时间量
		s.MaxRespTime = unixNano
	}
}

//重新统计
func (s *ServiceStateDetail) ReStats() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.statsReset()
}

//----------------------

//重置统计数据
func (s *ServiceStateDetail) statsReset() {
	s.StatsReqCount = 0
	s.StatsTimestamp = time.Now().UnixNano()
	s.StatsRespUnixNano = 0
}

func (s ServiceStateDetail) getStatsPass() int64 {
	return time.Now().UnixNano() - s.StartTimestamp
}

func (s ServiceStateDetail) statsFull() bool {
	return s.getStatsPass() >= s.StatsInterval
}
