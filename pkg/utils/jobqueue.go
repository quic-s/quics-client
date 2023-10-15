package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type JType string

const (
	ConnectServer    JType = "CONNECT_SERVER"
	ConnectRoot      JType = "CONNECT_ROOT"
	DisconnectServer JType = "DISCONNECT_SERVER"
	Rescan           JType = "RESCAN"
	ShowStatus       JType = "SHOW_STATUS"
	ChooseOne        JType = "CHOOSE_ONE"
)

type JStatus string

const (
	Ready   JStatus = "READY"
	Working JStatus = "WORKING"
	Done    JStatus = "DONE"
)

type JPriority string

const (
	Urgent JPriority = "URGENT"
	High   JPriority = "HIGH"
	Medium JPriority = "MEDIUM"
	Low    JPriority = "LOW"
)

type Job struct {
	next     *Job
	JobID    string
	Type     JType
	Status   JStatus
	Priority JPriority

	Path string // which task
	Do   func() // allow any function
}

type JobList struct {
	mut        sync.Mutex
	lowHead    *Job
	mediumHead *Job
	highHead   *Job
	urgentHead *Job
}

/* TODO
1. When Transcation gonna be started -> Create Job
2. When Job is created -> Add Job to JobQueue
3. When Job is added to JobQueue -> JobQueue is sorted by Priority
4. When Transcation is started -> Change Job Status to Working
5. When Transaction is done -> Change Job Status to Done
6. When Job Status is Done -> Remove Job from JobQueue
*/

var JList *JobList

func InitJobList() {
	JList = &JobList{}
}

func (JList *JobList) PushJob(jp JPriority, jt JType, path string, do func()) error {
	JList.mut.Lock()
	defer JList.mut.Unlock()
	newNode := newJobNode(jp, jt, path, do)
	switch jp {
	case Urgent:
		if JList.urgentHead == nil {
			JList.urgentHead = newNode
			return nil
		}
		newNode.next = JList.urgentHead
		JList.urgentHead = newNode
	case High:
		if JList.highHead == nil {
			JList.highHead = newNode
			return nil
		}
		newNode.next = JList.highHead
		JList.highHead = newNode
	case Medium:
		if JList.mediumHead == nil {
			JList.mediumHead = newNode
			return nil
		}
		newNode.next = JList.mediumHead
		JList.mediumHead = newNode
	case Low:
		if JList.lowHead == nil {
			JList.lowHead = newNode
			return nil
		}
		newNode.next = JList.lowHead
		JList.lowHead = newNode
	default:
		return fmt.Errorf("Invalid Priority")
	}
	return nil
}

func findLastJobForPop(target *Job, prev *Job) {

	for true {
		if target.next == nil {
			prev.next = nil // delete last node
			break
		}
		if target.next != nil {
			prev = target
			target = target.next

		}
	}
}

func findJobByPathForPop(path string, target *Job, prev *Job) bool {
	for target.next != nil {
		if target.Path == path {
			prev.next = target.next
			return true
		} else {
			prev = target
			target = target.next
		}
	}
	return false
}

func (JList *JobList) PopJob(jp JPriority) (*Job, error) {
	JList.mut.Lock()
	defer JList.mut.Unlock()
	target := &Job{}
	// get the first node of the list
	switch jp {
	case Urgent:
		if JList.urgentHead == nil {
			return nil, fmt.Errorf("no Job")
		}
		target = JList.urgentHead
		prev := target
		findLastJobForPop(target, prev)
		return target, nil
	case High:
		if JList.highHead == nil {
			return nil, fmt.Errorf("no Job")
		}
		target = JList.highHead
		prev := target
		findLastJobForPop(target, prev)
		return target, nil
	case Medium:
		if JList.mediumHead == nil {
			return nil, fmt.Errorf("no Job")
		}
		target = JList.mediumHead
		prev := target
		findLastJobForPop(target, prev)
		return target, nil
	case Low:
		if JList.lowHead == nil {
			return nil, fmt.Errorf("no Job")
		}
		target = JList.lowHead
		prev := target
		findLastJobForPop(target, prev)
		return target, nil
	default:
		return nil, fmt.Errorf("invalid Priority")
	}
}

func (JList *JobList) MoveJobAToB(path string, jp_a JPriority, jp_b JPriority) error {
	JList.mut.Lock()
	defer JList.mut.Unlock()

	target := &Job{}
	prev := &Job{}
	switch jp_a {
	case Urgent:
		target = JList.urgentHead
		prev = target
	case High:
		target = JList.highHead
		prev = target
	case Medium:
		target = JList.mediumHead
		prev = target
	case Low:
		target = JList.lowHead
		prev = target
	default:
		return fmt.Errorf("invalid Priority")
	}
	if !findJobByPathForPop(path, target, prev) {
		return fmt.Errorf("no Job")
	}
	target.Priority = jp_b
	switch jp_b {
	case Urgent:
		if JList.urgentHead == nil {
			JList.urgentHead = target
			return nil
		}
		target.next = JList.urgentHead
		JList.urgentHead = target
	case High:
		if JList.highHead == nil {
			JList.highHead = target
			return nil
		}
		target.next = JList.highHead
		JList.highHead = target
	case Medium:
		if JList.mediumHead == nil {
			JList.mediumHead = target
			return nil
		}
		target.next = JList.mediumHead
		JList.mediumHead = target
	case Low:
		if JList.lowHead == nil {
			JList.lowHead = target
			return nil
		}
		target.next = JList.lowHead
		JList.lowHead = target
	default:
		return fmt.Errorf("invalid Priority")
	}
	return nil

}

func (JList *JobList) PrintJobList() {
	fmt.Println("Urgent")
	for node := JList.urgentHead; node != nil; node = node.next {
		fmt.Println(node.JobID)
	}
	fmt.Println("High")
	for node := JList.highHead; node != nil; node = node.next {
		fmt.Println(node.JobID)
	}
	fmt.Println("Medium")
	for node := JList.mediumHead; node != nil; node = node.next {
		fmt.Println(node.JobID)
	}
	fmt.Println("Low")
	for node := JList.lowHead; node != nil; node = node.next {
		fmt.Println(node.JobID)
	}
}

func newJobNode(jp JPriority, jt JType, path string, do func()) *Job {
	return &Job{
		Priority: jp,
		Type:     jt,
		Status:   Ready,
		JobID:    uuid.New().String(),

		Path: path,
		Do:   do,
	}
}

func (j *Job) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(j); err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func (j *Job) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(j); err != nil {
		panic(err)
	}
}
