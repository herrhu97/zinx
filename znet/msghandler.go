package znet

import (
	"fmt"
	"strconv"

	"github.com/herrhu97/zinx/utils"
	"github.com/herrhu97/zinx/ziface"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgId(), " is not FOUND!")
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	_, ok := mh.Apis[msgId]
	if ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}

func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started.")
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// SendMsgToTaskQueue 轮询分配机制
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(),
		" request MsgID=", request.GetMsgId(), " to WorkerID=", workerID)
	mh.TaskQueue[workerID] <- request
}
