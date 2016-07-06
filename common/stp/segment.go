package stp

import (
	"sync"
	"sync/atomic"
)

type Segments struct {
	seqid  uint32
	cap    int
	data   []Segment
	wait   sync.WaitGroup
	waited int32
	guard  sync.Mutex
}

func NewSegments(seqid uint32, n int) *Segments {
	return &Segments{
		seqid: seqid,
		cap:   n,
	}
}

func (s *Segments) Ack(seqid uint32) *Segment {
	skip := 0
	s.guard.Lock()
	for i := range s.data {
		if !s.data[i].inited {
			if i == skip {
				skip++
			}
			continue
		}
		if s.data[i].seqid == seqid {
			if i == skip {
				skip++
			}
			for j := i + 1; j < len(s.data); j++ {
				if s.data[j].inited || j != skip {
					break
				}
				skip++
			}

			seg := s.data[i]
			s.data[i].inited = false

			// force skip if it's full
			if skip > 0 {
				s.data = s.data[skip:]
				if atomic.CompareAndSwapInt32(&s.waited, 1, 0) {
					s.wait.Done()
				}
			}
			s.guard.Unlock()
			return &seg
		}
	}
	s.guard.Unlock()
	return nil
}

func (s *Segments) New() *Segment {
	for {
		s.wait.Wait()
		s.guard.Lock()
		if len(s.data) == s.cap {
			if atomic.CompareAndSwapInt32(&s.waited, 0, 1) {
				s.wait.Add(1)
			}
			s.guard.Unlock()
			continue
		}
		s.seqid++
		s.data = append(s.data, Segment{
			seqid: s.seqid,
		})
		seg := &s.data[len(s.data)-1]
		s.guard.Unlock()
		return seg
	}
}

type Segment struct {
	seqid uint32
	cmd   uint32
	size  uint16
	ts    int64
	data  []byte

	inited bool
}

func (s *Segment) init(data []byte) {
	s.inited = true
}
