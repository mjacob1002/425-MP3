package filesystem;

import (
	//"fmt"
	"sync"
)

type FileStruct struct {
	Sdfsname string
	Rwlock sync.RWMutex
	deleted bool;
	replica int64; // is -1 when I own it, otherwise it's the node that send the replica
}

type FileRWLock struct {
	Mu sync.Mutex
	num_readers int64
	num_writers int64
	read_cv *sync.Cond
	write_cv *sync.Cond
	read_operations_queued int64
	write_operations_queued int64
}

func (s *FileStruct) RLock() {
	s.Rwlock.RLock();
}

func (s *FileStruct) RUnlock() {
	s.Rwlock.RUnlock();
}

func (s *FileStruct) WLock() {
	s.Rwlock.Lock();
}

func (s *FileStruct) WUnlock(){
	s.Rwlock.Unlock()
}

func NewFileStruct(sdfsName string, replica int64) *FileStruct{
	fileObject := FileStruct {
		Sdfsname: sdfsName,
		deleted: false,
		replica: replica,
	}
	// fileObject.Rwlock.Init()
	return &fileObject
}

func (s *FileRWLock) Init() {
	s.read_cv = sync.NewCond(&(s.Mu));
	s.write_cv = sync.NewCond(&(s.Mu))
	s.num_readers = 0
	s.num_writers = 0
	s.read_operations_queued = 0
	s.write_operations_queued = 0
}

func (s *FileRWLock) RLock(){
	s.Mu.Lock() // acquire the lock
	if s.num_readers >= 2 || s.num_writers >= 0 {
		s.read_operations_queued++;
		// sleep at the mutex under read, then undo the number of read_operations queued
		for s.num_readers >= 2 || s.num_writers >= 0{
			s.read_cv.Wait()
		}
		s.read_operations_queued--;
	}
	s.num_readers++;
	s.Mu.Unlock() // can allow another concurrent reader
}

func (s *FileRWLock) RUnlock(){
	s.Mu.Lock();
	s.num_readers--;
	if (s.write_operations_queued > 0) {
		// signal that write conditional variable
		s.write_cv.Signal()
	} else {
		// signal at read condition variable
		s.read_cv.Signal()
	}
	s.Mu.Unlock();
}

func (s *FileRWLock) WLock() {
	s.Mu.Lock();
	if (s.num_readers > 0  || s.num_writers > 0){
		s.write_operations_queued++;
		// sleep at write condition variable
		for s.num_readers > 0 || s.num_writers > 0 {
			s.write_cv.Wait()
		}
		// when you wake up decremenet writes queued up
		s.write_operations_queued--;
	}
	s.num_writers++;
}

func(s *FileRWLock) WUnlock() {
	s.num_writers--;
	if s.read_operations_queued > 0 {
		// signal read condition variable
		s.read_cv.Signal();
	} else {
		// signal write condition variable
		s.write_cv.Signal();
	}
	s.Mu.Unlock();
}
