package stp

import (
	"testing"

	"github.com/chzyer/logex"
	"github.com/chzyer/test"
)

func TestSegment(t *testing.T) {
	defer test.New(t)

	sgs := NewSegments(0, 3)
	sgs.New().init(nil)
	test.Equal(len(sgs.data), 1)
	sgs.New().init(nil)
	test.Equal(len(sgs.data), 2)
	sgs.New().init(nil)
	test.Equal(len(sgs.data), 3)
	logex.Struct(sgs.data)
	ret := sgs.Ack(3)
	test.NotNil(ret)
	test.Equal(int(ret.seqid), 3)
	test.Equal(len(sgs.data), 3)
	{
		ret := sgs.Ack(2)
		test.NotNil(ret)
		test.Equal(int(ret.seqid), 2)
		test.Equal(len(sgs.data), 3)
	}

	{
		ret := sgs.Ack(1)
		test.NotNil(ret)
		test.Equal(int(ret.seqid), 1)
		test.Equal(len(sgs.data), 0)
	}

}
