package squeue

import (
	"reflect"
	"time"

	"github.com/yyliziqiu/slib/sfile"
	"github.com/yyliziqiu/slib/ssnap"
)

// SnapSave 保存队列数据快照
func (q *Queue) SnapSave() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	return ssnap.Save(q.path, q.copyList())
}

// SnapLoad 加载队列数据快照
func (q *Queue) SnapLoad(item any) error {
	exist, err := sfile.Exist(q.path)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	lst := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(item)), 0, 0)
	lsp := reflect.New(lst.Type())

	err = ssnap.Load(q.path, lsp.Interface())
	if err != nil {
		return err
	}

	size := lsp.Elem().Len()
	data := lsp.Elem().Slice(0, size)

	var list []any
	for i := 0; i < size; i++ {
		list = append(list, data.Index(i).Interface())
	}

	q.reset(list)

	return nil
}

// SnapDuplicate 保存队列数据快照副本
func (q *Queue) SnapDuplicate(d time.Duration) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	return ssnap.Duplicate(q.path, q.copyList(), d)
}
