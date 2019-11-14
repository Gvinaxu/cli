package task

import (
	"errors"
	"log"
	"reflect"
	"runtime"
	"sort"
	"time"
)

const MAX_JOB_COUNT = 300

var local_time *time.Location = time.Local

func ChangeLocalTime(localtime *time.Location) {
	local_time = localtime
}

type Job struct {
	funcs map[string]interface{}

	fparams map[string]([]interface{})

	interval uint64

	jobFunc string

	unit string

	atTime string

	lastRun time.Time

	nextRun time.Time

	period time.Duration

	startDay time.Weekday
}

func NewJob(interval uint64) *Job {
	return &Job{
		funcs:    map[string]interface{}{},
		fparams:  map[string]([]interface{}){},
		interval: interval,
		jobFunc:  "",
		unit:     "",
		atTime:   "",
		lastRun:  time.Unix(0, 0),
		nextRun:  time.Unix(0, 0),
		period:   0,
		startDay: time.Sunday,
	}
}

//作业的核心执行
func (j *Job) Do(job_fun interface{}, params ...interface{}) {
	typ := reflect.TypeOf(job_fun)
	if typ.Kind() != reflect.Func {
		panic("作业队列只允许进行函数的调度")
	}
	fname := getFunctionName(job_fun)
	j.funcs[fname] = job_fun
	j.fparams[fname] = params
	j.jobFunc = fname
	j.scheduleNextRun()
}

//运行作业 并且 重新调度
func (j *Job) run() (result []reflect.Value, err error) {
	f := reflect.ValueOf(j.funcs[j.jobFunc])
	params := j.fparams[j.jobFunc]
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of param is not adapted.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	j.lastRun = time.Now()
	j.scheduleNextRun()
	return
}

//重新下次的调度
func (j *Job) scheduleNextRun() {
	if j.lastRun == time.Unix(0, 0) {
		if j.unit == "weeks" {
			i := time.Now().Weekday() - j.startDay
			if i < 0 {
				i = 7 + i
			}
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-int(i), 0, 0, 0, 0, local_time)

		} else {
			j.lastRun = time.Now()
		}
	}

	if j.period != 0 {
		j.nextRun = j.lastRun.Add(j.period * time.Second)
	} else {
		switch j.unit {
		case "minutes":
			j.period = time.Duration(j.interval * 60)
			break
		case "hours":
			j.period = time.Duration(j.interval * 60 * 60)
			break
		case "days":
			j.period = time.Duration(j.interval * 60 * 60 * 24)
			break
		case "weeks":
			j.period = time.Duration(j.interval * 60 * 60 * 24 * 7)
			break
		case "seconds":
			j.period = time.Duration(j.interval)
		}
		j.nextRun = j.lastRun.Add(j.period * time.Second)
	}
}

//在具体时间执行
func (j *Job) At(t string) *Job {
	hour := int((t[0]-'0')*10 + (t[1] - '0'))
	min := int((t[3]-'0')*10 + (t[4] - '0'))
	if hour < 0 || hour > 23 || min < 0 || min > 59 {
		panic("time format error.")
	}
	// time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	mock := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), int(hour), int(min), 0, 0, local_time)

	if j.unit == "days" {
		if time.Now().After(mock) {
			j.lastRun = mock
		} else {
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, hour, min, 0, 0, local_time)
		}
	} else if j.unit == "weeks" {
		if time.Now().After(mock) {
			i := mock.Weekday() - j.startDay
			if i < 0 {
				i = 7 + i
			}
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-int(i), hour, min, 0, 0, local_time)
		} else {
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-7, hour, min, 0, 0, local_time)
		}
	}
	return j
}

func (j *Job) should_run() bool {
	return time.Now().After(j.nextRun)
}

//设置 时间单元为 1 秒
func (j *Job) Second() (job *Job) {
	if j.interval != 1 {
		panic("")
		return
	}
	job = j.Seconds()
	return
}

//设置 时间单元为秒
func (j *Job) Seconds() (job *Job) {
	j.unit = "seconds"
	return j
}

//设置 时间单元为 1 分钟
func (j *Job) Minute() (job *Job) {
	if j.interval != 1 {
		panic("")
		return
	}
	job = j.Minutes()
	return
}

//设置 时间单元为分钟
func (j *Job) Minutes() (job *Job) {
	j.unit = "minutes"
	return j
}

//设置 时间单元为 1 小时
func (j *Job) Hour() (job *Job) {
	if j.interval != 1 {
		panic("")
		return
	}
	job = j.Hours()
	return
}

//设置 时间单元为小时
func (j *Job) Hours() (job *Job) {
	j.unit = "hours"
	return j
}

//设置 时间单元为 1 天
func (j *Job) Day() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	job = j.Days()
	return
}

//设置 时间单元为天
func (j *Job) Days() *Job {
	j.unit = "days"
	return j
}

//Set the units as weeks
func (j *Job) Weeks() *Job {
	j.unit = "weeks"
	return j
}

//设置 时间为星期一
func (j *Job) Monday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 1
	job = j.Weeks()
	return
}

//设置 时间为星期二
func (j *Job) Tuesday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 2
	job = j.Weeks()
	return
}

//设置 时间为星期三
func (j *Job) Wednesday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 3
	job = j.Weeks()
	return
}

//设置 时间为星期四
func (j *Job) Thursday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 4
	job = j.Weeks()
	return
}

//设置 时间为星期五
func (j *Job) Friday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 5
	job = j.Weeks()
	return
}

//设置 时间为星期六
func (j *Job) Saturday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 6
	job = j.Weeks()
	return
}

//设置 时间为星期日
func (j *Job) Sunday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 0
	job = j.Weeks()
	return
}

//作业调度类
type Scheduler struct {
	//存储这具体的作业
	jobs [MAX_JOB_COUNT]*Job

	//作业的数目
	size int
}

//构造调度器实例
func NewScheduler() *Scheduler {
	return &Scheduler{
		jobs: [MAX_JOB_COUNT]*Job{},
		size: 0,
	}
}

//作业调度类继承sort.Interface{}接口, 方便对作业进行排序

func (s *Scheduler) Len() int {
	return s.size
}

func (s *Scheduler) Swap(i, j int) {
	s.jobs[i], s.jobs[j] = s.jobs[j], s.jobs[i]
}

//next_run从小到大
func (s *Scheduler) Less(i, j int) bool {
	return s.jobs[j].nextRun.After(s.jobs[i].nextRun)
}

//设计 Every 功能, 规划一个调度周期处理 , 返回 Job 实例
func (s *Scheduler) Every(interval uint64) *Job {
	job := NewJob(interval)
	s.jobs[s.size] = job
	s.size++
	return job
}

func getFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf((fn)).Pointer()).Name()
}

func (s *Scheduler) getRunnableJobs() (running_jobs [MAX_JOB_COUNT]*Job, n int) {
	runnable_jobs := [MAX_JOB_COUNT]*Job{}
	n = 0
	sort.Sort(s)
	for i := 0; i < s.size; i++ {
		//判断作业是否应该要允许
		if s.jobs[i].should_run() {
			runnable_jobs[n] = s.jobs[i]
			n++
		} else {
			break
		}
	}
	return runnable_jobs, n
}

//从调度器中删除特点的作业
func (s *Scheduler) Remove(j interface{}) {
	i := 0
	for ; i < s.size; i++ {
		//找到目标执行方法对应的作业索引
		if s.jobs[i].jobFunc == getFunctionName(j) {
			break
		}
	}
	//元素左移
	for index := (i + 1); index < s.size; index++ {
		s.jobs[i] = s.jobs[index]
		i++
	}
	//调度器当前size调整
	s.size = s.size - 1
}

//删除调度器中的所有作业
func (s *Scheduler) Clear() {
	for i := 0; i < s.size; i++ {
		s.jobs[i] = nil
	}
	s.size = 0
}

//执行已经准备好的作业
func (s *Scheduler) RunReadyed() {
	runnable_jobs, n := s.getRunnableJobs()
	if n != 0 {
		for i := 0; i < n; i++ {
			_, err := runnable_jobs[i].run()
			if err != nil {
				log.Println(err)
			}
		}
	}
}

//从调度器中开启所有的作业 (这些作业为那些有计划的作业)
func (s *Scheduler) Start() {
	for {
		s.RunReadyed()
	}
}

//从调度器中运行所有的作业
func (s *Scheduler) RunAll() {
	for i := 0; i < s.size; i++ {
		s.jobs[i].run()
	}
}

//从调度器中延迟规定时间后运行所有的作业
func (s *Scheduler) RunAllwithDelay(d int) {
	for i := 0; i < s.size; i++ {
		s.jobs[i].run()
		time.Sleep(time.Duration(d))
	}
}

//获取作业的下次调度时间

func (s *Scheduler) NextRun() (*Job, time.Time) {
	if s.size <= 0 {
		return nil, time.Now()

	}
	sort.Sort(s)
	return s.jobs[0], s.jobs[0].nextRun
}
