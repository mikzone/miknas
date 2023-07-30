package cmdexec

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/mikzone/miknas/server/miknas"
	"github.com/panjf2000/ants/v2"
)

const JobStWaiting = "waiting"
const JobStLineUp = "lineUp"
const JobStRunning = "running"
const JobStDone = "done"
const JobStCanceled = "canceled"
const JobStErrStop = "errstop"

type JobItem struct {
	JobId        string
	Uid          string
	Cmd          *exec.Cmd
	Cancel       context.CancelFunc
	NameSpace    string
	RunningState string
	CancelUser   string
	OutTxt       string
	FailTxt      string
}

func (item *JobItem) PackClientDict(detail bool) miknas.H {
	ret := miknas.H{
		"jobId":        item.JobId,
		"uid":          item.Uid,
		"cmd":          item.Cmd.String(),
		"cwd":          item.Cmd.Dir,
		"nameSpace":    item.NameSpace,
		"runningState": item.RunningState,
		"cancelUser":   item.CancelUser,
		"failtxt":      item.FailTxt,
	}
	if detail {
		ret["out"] = item.OutTxt
	}
	return ret
}

func (item *JobItem) SetState(state string) {
	item.RunningState = state
}

func (item *JobItem) CheckInState(states ...string) bool {
	for _, state := range states {
		if state == item.RunningState {
			return true
		}
	}
	return false
}

type JobMgr struct {
	GenedId   int
	Jobs      map[string]*JobItem
	Pool      *ants.PoolWithFunc
	Lineups   []string // those can be submit
	Waits     map[string][]string
	LineupMux sync.Mutex
	WaitMux   sync.Mutex
}

func NewJobMgr() *JobMgr {
	jm := &JobMgr{
		Jobs:  map[string]*JobItem{},
		Waits: map[string][]string{},
	}
	pool, err := ants.NewPoolWithFunc(20, func(i interface{}) {
		item := i.(*JobItem)
		jm.RunJobItem(item)
		jm.TryMaintainNamespace(item.NameSpace, true)
	}, ants.WithNonblocking(true))
	if err != nil {
		panic(miknas.NewFailRet("PoolCreateFail"))
	}
	jm.Pool = pool
	return jm
}

func (jm *JobMgr) NewId() string {
	jm.GenedId += 1
	return fmt.Sprint(jm.GenedId)
}

func (jm *JobMgr) AddToLineUp(jobid string) {
	item, exist := jm.Jobs[jobid]
	if !exist {
		return
	}
	jm.LineupMux.Lock()
	defer jm.LineupMux.Unlock()
	item.RunningState = JobStLineUp
	jm.Lineups = append(jm.Lineups, jobid)
}

func (jm *JobMgr) AddJob(item *JobItem) {
	jobid := jm.NewId()
	item.JobId = jobid
	jm.Jobs[jobid] = item
	namespace := item.NameSpace
	if namespace == "" {
		// 空命名空间的直接提交
		jm.AddToLineUp(jobid)
	} else {
		// 有命名空间的先加到队列，后续管理
		jm.WaitMux.Lock()
		defer jm.WaitMux.Unlock()
		_, exist := jm.Waits[item.NameSpace]
		if !exist {
			jm.Waits[namespace] = []string{jobid}
		} else {
			jm.Waits[namespace] = append(jm.Waits[namespace], jobid)
		}
		jm.TryMaintainNamespace(namespace, false)
	}
}

func (jm *JobMgr) TryStopJob(ch *miknas.ContextHelper, jobid string) error {
	item, exist := jm.Jobs[jobid]
	if !exist {
		return miknas.NewFailRet("jobid(%s) is not exist", jobid)
	}
	if item.CheckInState(JobStCanceled, JobStDone, JobStErrStop) {
		return nil
	}
	uid := ch.GetUserAuth().MustGetUid()
	item.CancelUser = uid
	if item.CheckInState(JobStWaiting, JobStLineUp) {
		item.FailTxt += fmt.Sprintf("[用户 %s 取消了该任务]", uid)
		item.SetState(JobStCanceled)
	} else if item.CheckInState(JobStRunning) {
		item.FailTxt += fmt.Sprintf("[用户 %s 取消了该任务]", uid)
		item.Cancel()
	}
	return nil
}

func (jm *JobMgr) TryMaintainNamespace(namespace string, needlock bool) {
	if namespace == "" {
		return
	}
	if needlock {
		jm.WaitMux.Lock()
		defer jm.WaitMux.Unlock()
	}
	waitlist, exist := jm.Waits[namespace]
	if !exist {
		return
	}
	for len(waitlist) > 0 {
		jobid := waitlist[0]
		item, exist := jm.Jobs[jobid]
		if !exist {
			waitlist = waitlist[1:]
			jm.Waits[namespace] = waitlist
			continue
		}
		if item.CheckInState(JobStDone, JobStCanceled, JobStErrStop) {
			waitlist = waitlist[1:]
			jm.Waits[namespace] = waitlist
			continue
		} else if item.CheckInState(JobStLineUp, JobStRunning) {
			// 有进行中的，不管
			return
		} else if item.CheckInState(JobStWaiting) {
			jm.AddToLineUp(jobid)
			return
		}
	}
	// namespace 下的消耗完了的话,维护一下lineUps的
	// jm.TryMaintainLineUps()
}

func (jm *JobMgr) TryMaintainLineUps() {
	if jm.Pool.Free() == 0 {
		return
	}
	jm.LineupMux.Lock()
	defer jm.LineupMux.Unlock()
	for len(jm.Lineups) > 0 {
		jobid := jm.Lineups[0]
		item, exist := jm.Jobs[jobid]
		if !exist {
			jm.Lineups = jm.Lineups[1:]
			continue
		}
		if item.CheckInState(JobStCanceled) {
			jm.Lineups = jm.Lineups[1:]
			continue
		}
		if !item.CheckInState(JobStLineUp) {
			fmt.Printf("TryMaintainLineUps Failed by invalid state(%v) of jobid(%s)", item.RunningState, jobid)
			jm.Lineups = jm.Lineups[1:]
			continue
		}
		err := jm.Pool.Invoke(item)
		if err == nil {
			jm.Lineups = jm.Lineups[1:]
			return
		}
	}
}

func (jm *JobMgr) RunJobItem(item *JobItem) {
	c := item.Cmd
	stdout, err := c.StdoutPipe()
	if err != nil {
		item.SetState(JobStErrStop)
		item.FailTxt += fmt.Sprintf("程序获取输出管道失败: %q", err)
		return
	}
	// 命令的错误输出和标准输出都连接到同一个管道
	c.Stderr = c.Stdout
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			item.OutTxt += readString
			if err != nil {
				if err != io.EOF {
					item.OutTxt += fmt.Sprintf("\n读取运行结果出现错误: %q", err)
				}
				return
			}
		}
	}(&wg)
	err = c.Start()
	if err != nil {
		item.FailTxt += fmt.Sprintf("程序Start失败: %q", err)
		item.SetState(JobStErrStop)
		wg.Wait()
		return
	}
	item.SetState(JobStRunning)
	wg.Wait()
	waitErr := c.Wait()
	if err != nil {
		item.FailTxt += fmt.Sprintf("Wait发生错误: %v", waitErr)
		item.SetState(JobStErrStop)
	} else if item.CancelUser != "" {
		item.SetState(JobStCanceled)
	} else {
		item.SetState(JobStDone)
	}
}

func GetJobMgr(ch *miknas.ContextHelper) *JobMgr {
	ext := ch.GetApp().GetExt(ExtId).(*MikNasExt)
	return ext.JmInst
}

func NewJob(name string, arg ...string) *JobItem {
	ctx, cancel := context.WithCancel(context.Background())
	c := exec.CommandContext(ctx, name, arg...)
	item := &JobItem{
		Cmd:          c,
		Cancel:       cancel,
		RunningState: JobStWaiting,
	}
	return item
}

func SubmitJob(ch *miknas.ContextHelper, item *JobItem) {
	userauth := ch.GetUserAuth()
	uid := userauth.MustGetUid()
	item.Uid = uid
	jobmgr := GetJobMgr(ch)
	jobmgr.AddJob(item)
}
