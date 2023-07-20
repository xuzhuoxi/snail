// Package imodule
// Created by xuzhuoxi
// on 2019-02-10.
// @author xuzhuoxi
//
package imodule

import (
	"runtime"
	"sync"
	"time"
)

type ISockState interface {
	GetSockName() string
	GetSockSockConnections() uint64
	GetSockWeight() float64
}

type ISockStateDetail interface {
	// GetPassNano 运行时间
	GetPassNano() int64

	// StatsWeight 当前统计的服务权重(连接数*统计时间/统计响应时间)
	// 越大代表压力越大
	StatsWeight() float64

	// RespCoefficient 响应系数(响应总时间 / (统计总时间 * 逻辑cpu数)),
	// 注意：结果正常设置下为[0,1]
	RespCoefficient() float64
	//平均响应时间(响应总时间/响应次数)
	RespAvgTime() float64
	//请求密度(次数/秒)
	ReqDensityTime() int

	//区间响应系数(区间响应总时间 / (区间统计总时间 * 逻辑cpu数)),
	// 注意：结果正常设置下为[0,1]
	StatsRespCoefficient() float64
	//区间平均响应时间(响应总时间/响应次数)
	StatsRespAvgTime() float64
	//区间时间请求密度(次数/秒)
	StatsReqDensityTime() int
}

func NewSockStateDetail(name string, statsInterval int64) *SockStateDetail {
	return &SockStateDetail{SockName: name, StatsInterval: statsInterval}
}

//------------------------------

//Sock的拥有者信息
type SockOwner struct {
	//平台id
	PlatformId string
	//模块id
	ModuleId string
	//模块类型名称
	ModuleName ModuleName
}

type SockState struct {
	//名称
	SockName string
	//连接数
	SockConnections uint64
	//压力
	SockWeight float64
	////响应系数(响应总时间 / (统计总时间 * 逻辑cpu数)),
	////注意：结果正常设置下为[0,1]
	//RespCoefficient float64
	////平均响应时间(响应总时间/响应次数)
	//RespAvgTime float64
	////请求密度(次数/秒)
	//ReqDensityTime int
}

func (ss *SockState) GetSockName() string {
	return ss.SockName
}

func (ss *SockState) GetSockSockConnections() uint64 {
	return ss.SockConnections
}

func (ss *SockState) GetSockWeight() float64 {
	return ss.SockWeight
}

//------------------------------

type SockStateDetail struct {
	SockName string
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
func (s SockStateDetail) GetPassNano() int64 {
	return time.Now().UnixNano() - s.StartTimestamp
}

//当前统计的服务权重(连接数 + 统计响应时间 / 统计时间 )
//越大代表压力越大
func (s SockStateDetail) StatsWeight() float64 {
	if 0 == s.StatsRespUnixNano {
		return 0
	} else {
		return s.StatsRespCoefficient()
	}
}

//响应系数(响应总时间/统计总时间),
// 注意：结果正常设置下为[0,1]
func (s SockStateDetail) RespCoefficient() float64 {
	return float64(s.TotalRespTime) / (float64(s.GetPassNano()) * float64(runtime.NumCPU()))
}

//平均响应时间(响应总时间/响应次数)
func (s SockStateDetail) RespAvgTime() float64 {
	return float64(s.TotalRespTime) / float64(s.TotalReqCount)
}

//时间请求密度(次数/秒)
func (s SockStateDetail) ReqDensityTime() int {
	pass := s.GetPassNano()
	return int(int64(time.Second) * s.TotalReqCount / pass)
}

//区间响应系数(响应总时间/统计总时间),
// 注意：结果正常设置下为[0,1]
func (s SockStateDetail) StatsRespCoefficient() float64 {
	return float64(s.StatsRespUnixNano) / (float64(s.getStatsPass()) * float64(runtime.NumCPU()))
}

//区间平均响应时间(响应总时间/响应次数)
func (s SockStateDetail) StatsRespAvgTime() float64 {
	return float64(s.StatsRespUnixNano) / float64(s.StatsReqCount)
}

//区间时间请求密度(次数/秒)
func (s SockStateDetail) StatsReqDensityTime() int {
	pass := s.getStatsPass()
	return int(int64(time.Second) * s.StatsReqCount / pass)
}

//-------------------------------

//启动
func (s *SockStateDetail) Start() {
	s.lock.Lock()
	defer s.lock.Unlock()
	now := time.Now().UnixNano()
	s.StartTimestamp = now
	s.StatsTimestamp = now
}

//增加一个连接
func (s *SockStateDetail) AddLinkCount() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.LinkCount++
	if s.LinkCount > s.MaxLinkCount { //更新最大连接数
		s.MaxLinkCount = s.LinkCount
	}
}

//减少一个连接
func (s *SockStateDetail) RemoveLinkCount() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.LinkCount--
}

//----------------------

//增加一个请求
func (s *SockStateDetail) AddReqCount() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.statsFull() {
		s.statsReset()
	}
	s.TotalReqCount++
	s.StatsReqCount++
}

//增加响应时间量
func (s *SockStateDetail) AddRespUnixNano(unixNano int64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.StatsRespUnixNano += unixNano
	s.TotalRespTime += unixNano
	if unixNano > s.MaxRespTime { //更新最大响应时间量
		s.MaxRespTime = unixNano
	}
}

//重新统计
func (s *SockStateDetail) ReStats() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.statsReset()
}

//----------------------

//重置统计数据
func (s *SockStateDetail) statsReset() {
	s.StatsReqCount = 0
	s.StatsTimestamp = time.Now().UnixNano()
	s.StatsRespUnixNano = 0
}

func (s SockStateDetail) getStatsPass() int64 {
	return time.Now().UnixNano() - s.StatsTimestamp
}

func (s SockStateDetail) statsFull() bool {
	return s.getStatsPass() >= s.StatsInterval
}
