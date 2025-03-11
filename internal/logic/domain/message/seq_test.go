package message_test

import (
	"fmt"
	"nbim/internal/logic/domain/message"
	"testing"
)

func Test_GetSeq(t *testing.T) {
	a, _ := message.Seq.GetSeq(1000)
	b, _ := message.Seq.GetSeq(1000)
	c, _ := message.Seq.GetSeq(1000)
	d, _ := message.Seq.GetSeq(1000)
	e, _ := message.Seq.GetSeq(1000)
	fmt.Printf("seqID-[%d]\n", a)
	fmt.Printf("seqID-[%d]\n", b)
	fmt.Printf("seqID-[%d]\n", c)
	fmt.Printf("seqID-[%d]\n", d)
	fmt.Printf("seqID-[%d]\n", e)
}
