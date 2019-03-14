//
//Created by xuzhuoxi
//on 2019-03-09.
//@author xuzhuoxi
//
package basis

import (
	"strconv"
)

var (
	MaxTeamMember = 0
	TeamId        = 1000
	TeamName      = "的队伍"
)

var (
	TeamCorpsId   = 1000
	TeamCorpsName = "的军团"
)

func GetTeamId() string {
	defer func() { TeamId++ }()
	return "T_" + strconv.Itoa(TeamId)
}

func GetTeamCorpsId() string {
	defer func() { TeamCorpsId++ }()
	return "TC_" + strconv.Itoa(TeamId)
}

type ITeamControl interface {
	//队长
	Leader() string
	//用户列表
	MemberList() []string
	//检查用户
	ContainMember(memberId string) bool
	//加入用户,进行唯一性检查
	AcceptMember(memberId string) error
	//从组中移除用户
	DropMember(memberId string) error
	//从组中移除用户
	RiseLeader(memberId string) error
	//解散队伍
	DisbandTeam() error
}
