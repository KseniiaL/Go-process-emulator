package process

import "math/rand"

type Process struct {
	CreateTime uint64
	StartTime  uint64
	EndTime    uint64
	WaitTime   uint64
	ExecTime   uint64
	WorkTime   uint64
	Id         uint64
	ActualExec uint64
}

func (p Process) CreateProcess(curT uint64, id uint64) Process {
	p.CreateTime = curT
	p.Id = id
	p.ExecTime = (rand.Uint64() % 10) + 1
	p.ActualExec = p.ExecTime
	p.WaitTime = 1
	return p
}
