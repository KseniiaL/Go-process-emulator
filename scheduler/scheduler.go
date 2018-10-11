package scheduler

import (
	"fmt"
	"math/rand"
	"github.com/KseniiaL/Go-process-emulator/process"
)

type Scheduler struct {
	RR       []process.Process
	SRTF     []process.Process
	RRfin    []process.Process
	SRTFfin  []process.Process
	curT     uint64
	lastRRid uint
	id       uint64
}

func findMin(s []process.Process) int64 {
	if len(s) == 0 {
		return -1
	}
	var id int64 = 0
	for i := 0; i < len(s); i++ {
		if s[i].ExecTime < s[id].ExecTime {
			id = int64(i)
		}
	}
	return id
}

func (sch *Scheduler) _RR(count int) {
	var s []process.Process
	for i := range sch.RR {
		if uint(i) == sch.lastRRid {
			if sch.RR[i].StartTime == 0 {
				sch.RR[i].StartTime = sch.curT
			}

			sch.RR[i].ExecTime--

			if sch.RR[i].ExecTime > 0 {
				s = append(s, sch.RR[i])
			} else {

				sch.RR[i].EndTime = sch.curT + 1
				sch.RR[i].WaitTime = sch.RR[i].EndTime - sch.RR[i].ActualExec - sch.RR[i].CreateTime
				sch.RR[i].WorkTime = sch.RR[i].EndTime + sch.RR[i].ActualExec - sch.RR[i].StartTime - 1

				sch.RRfin = append(sch.RRfin, sch.RR[i])
			}
		} else {
			s = append(s, sch.RR[i])
		}
	}

	sch.RR = s

	sch.lastRRid++
	if sch.lastRRid >= uint(len(sch.RR)) {
		sch.lastRRid = 0
	}

	// creating new processes
	if rand.Uint64()%100 > 80 && count < 50 {
		num := rand.Uint64()%2 + 1
		for num > 0 {
			sch.RR = append(sch.RR, process.Process{}.CreateProcess(sch.curT, sch.id))
			sch.id++
			num--
		}
		sch.SRTF = append(sch.SRTF, process.Process{}.CreateProcess(sch.curT, sch.id))
		sch.id++
	}
}

func Routine() {
	cnt := 0
	scheduler := Scheduler{}
	for i := 0; i < 50 || len(scheduler.RR) > 0 || len(scheduler.SRTF) > 0; i++ {
		scheduler.curT++
		if cnt%10 <= 8 {
			scheduler._RR(i)
			cnt++
		} else {

			var s []process.Process
			min := findMin(scheduler.SRTF)

			if min < 0 {
				cnt++
				continue
			}

			for scheduler.SRTF[min].ExecTime > 0 {

				scheduler.SRTF[min].ExecTime--

				if scheduler.SRTF[min].StartTime == 0 {
					scheduler.SRTF[min].StartTime = scheduler.curT
				}

				// creating new processes
				if rand.Uint64()%100 > 80 && i < 50 {
					num := rand.Uint64()%2 + 1
					for num > 0 {
						scheduler.RR = append(scheduler.RR, process.Process{}.CreateProcess(scheduler.curT, scheduler.id))
						scheduler.id++
						num--
					}
					scheduler.SRTF = append(scheduler.SRTF, process.Process{}.CreateProcess(scheduler.curT, scheduler.id))
					scheduler.id++
				}

				min = findMin(scheduler.SRTF)
				cnt++
			}

			for i := range scheduler.SRTF {
				if int64(i) != min {
					s = append(s, scheduler.SRTF[i])
				} else {
					scheduler.SRTF[i].EndTime = scheduler.curT + 1
					scheduler.SRTF[i].WaitTime = scheduler.SRTF[i].EndTime - scheduler.SRTF[i].ActualExec - scheduler.SRTF[i].CreateTime
					scheduler.SRTF[i].WorkTime = scheduler.SRTF[i].EndTime + scheduler.SRTF[i].ActualExec - scheduler.SRTF[i].StartTime - 1

					scheduler.SRTFfin = append(scheduler.SRTFfin, scheduler.SRTF[i])
				}
			}

			scheduler.SRTF = s
		}
	}

	fmt.Println("Main processes:")
	print(scheduler.RRfin)
	fmt.Println("\nBack processes:")
	print(scheduler.SRTFfin)
}

func print(c []process.Process) {
	for i := range c {
		fmt.Printf("Id: %3d,    exetime: %3d,    create: %3d,    start: %3d,    end: %3d,    work: %3d,    wait: %3d\n",
			c[i].Id, c[i].ActualExec, c[i].CreateTime, c[i].StartTime, c[i].EndTime, c[i].WorkTime, c[i].WaitTime)
	}
}
