// Copyright 2016 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package prog

import (
	"bytes"
	"testing"
)

func TestAssignSizeRandom(t *testing.T) {
	target, rs, iters := initTest(t)
	ct := target.DefaultChoiceTable()
	for i := 0; i < iters; i++ {
		p := target.Generate(rs, 10, ct)
		data0 := p.Serialize()
		for _, call := range p.Calls {
			target.assignSizesCall(call)
		}
		if data1 := p.Serialize(); !bytes.Equal(data0, data1) {
			t.Fatalf("different lens assigned, initial:\n%s\nnew:\n%s\n", data0, data1)
		}
		p.Mutate(rs, 10, ct, nil, nil)
		p.Serialize()
		for _, call := range p.Calls {
			target.assignSizesCall(call)
		}
	}
}

func TestAssignSize(t *testing.T) {
	// nolint: lll
	TestDeserializeHelper(t, "test", "64", func(target *Target, p *Prog) {
		for _, call := range p.Calls {
			target.assignSizesCall(call)
		}
	}, []DeserializeTest{
		{
			In:  "test$length0(&(0x7f0000000000)={0xff, 0x0})",
			Out: "test$length0(&(0x7f0000000000)={0xff, 0x2})",
		},
		{
			In:  "test$length1(&(0x7f0000001000)={0xff, 0x0})",
			Out: "test$length1(&(0x7f0000001000)={0xff, 0x4})",
		},
		{
			In:  "test$length2(&(0x7f0000001000)={0xff, 0x0})",
			Out: "test$length2(&(0x7f0000001000)={0xff, 0x8})",
		},
		{
			In:  "test$length3(&(0x7f0000005000)={0xff, 0x0, 0x0})",
			Out: "test$length3(&(0x7f0000005000)={0xff, 0x4, 0x2})",
		},
		{
			In:  "test$length4(&(0x7f0000003000)={0x0, 0x0})",
			Out: "test$length4(&(0x7f0000003000)={0x2, 0x2})",
		},
		{
			In:  "test$length5(&(0x7f0000002000)={0xff, 0x0})",
			Out: "test$length5(&(0x7f0000002000)={0xff, 0x4})",
		},
		{
			In:  "test$length6(&(0x7f0000002000)={[0xff, 0xff, 0xff, 0xff], 0x0})",
			Out: "test$length6(&(0x7f0000002000)={[0xff, 0xff, 0xff, 0xff], 0x4})",
		},
		{
			In:  "test$length7(&(0x7f0000003000)={[0xff, 0xff, 0xff, 0xff], 0x0})",
			Out: "test$length7(&(0x7f0000003000)={[0xff, 0xff, 0xff, 0xff], 0x8})",
		},
		{
			In:  "test$length8(&(0x7f000001f000)={0x00, {0xff, 0x0, 0x00, [0xff, 0xff, 0xff]}, [{0xff, 0x0, 0x00, [0xff, 0xff, 0xff]}], 0x00, 0x0, [0xff, 0xff]})",
			Out: "test$length8(&(0x7f000001f000)={0x32, {0xff, 0x1, 0x10, [0xff, 0xff, 0xff]}, [{0xff, 0x1, 0x10, [0xff, 0xff, 0xff]}], 0x10, 0x1, [0xff, 0xff]})",
		},
		{
			In:  "test$length9(&(0x7f000001f000)={&(0x7f0000000000/0x5000)=nil, 0x0000})",
			Out: "test$length9(&(0x7f000001f000)={&(0x7f0000000000/0x5000)=nil, 0x5000})",
		},
		{
			In:  "test$length10(&(0x7f0000000000/0x5000)=nil, 0x0000, 0x0000, 0x0000, 0x0000)",
			Out: "test$length10(&(0x7f0000000000/0x5000)=nil, 0x5000, 0x5000, 0x2800, 0x1400)",
		},
		{
			In:  "test$length11(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, 0x00)",
			Out: "test$length11(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, 0x30)",
		},
		{
			In:  "test$length12(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, 0x00)",
			Out: "test$length12(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, 0x30)",
		},
		{
			In:  "test$length13(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, &(0x7f0000001000)=0x00)",
			Out: "test$length13(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, &(0x7f0000001000)=0x30)",
		},
		{
			In:  "test$length14(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, &(0x7f0000001000)=0x00)",
			Out: "test$length14(&(0x7f0000000000)={0xff, 0xff, [0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]}, &(0x7f0000001000)=0x30)",
		},
		{
			In:  "test_length15(0xff, 0x0)",
			Out: "test_length15(0xff, 0x2)",
		},
		{
			In:  "test$length16(&(0x7f0000000000)={[0x42, 0x42], 0xff, 0xff, 0xff, 0xff, 0xff})",
			Out: "test$length16(&(0x7f0000000000)={[0x42, 0x42], 0x2, 0x10, 0x8, 0x4, 0x2})",
		},
		{
			In:  "test$length17(&(0x7f0000000000)={0x42, 0xff, 0xff, 0xff, 0xff})",
			Out: "test$length17(&(0x7f0000000000)={0x42, 0x8, 0x4, 0x2, 0x1})",
		},
		{
			In:  "test$length18(&(0x7f0000000000)={0x42, 0xff, 0xff, 0xff, 0xff})",
			Out: "test$length18(&(0x7f0000000000)={0x42, 0x8, 0x4, 0x2, 0x1})",
		},
		{
			In:  "test$length19(&(0x7f0000000000)={{0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0xff}, 0xff, 0xff, 0xff})",
			Out: "test$length19(&(0x7f0000000000)={{0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x14}, 0x14, 0x14, 0x5})",
		},
		{
			In:  "test$length20(&(0x7f0000000000)={{{0xff, 0xff, 0xff, 0xff}, 0xff, 0xff, 0xff}, 0xff, 0xff})",
			Out: "test$length20(&(0x7f0000000000)={{{0x4, 0x4, 0x7, 0x9}, 0x7, 0x7, 0x9}, 0x9, 0x9})",
		},
		{
			In:  "test$length21(&(0x7f0000000000)=0x0, 0x0)",
			Out: "test$length21(&(0x7f0000000000), 0x40)",
		},
		{
			In:  "test$length22(&(0x7f0000000000)='12345', 0x0)",
			Out: "test$length22(&(0x7f0000000000)='12345', 0x28)",
		},
		{
			In:  "test$length23(&(0x7f0000000000)={0x1, {0x2, 0x0}})",
			Out: "test$length23(&(0x7f0000000000)={0x1, {0x2, 0x6}})",
		},
		{
			In:  "test$length24(&(0x7f0000000000)={{0x0, {0x0}}, {0x0, {0x0}}})",
			Out: "test$length24(&(0x7f0000000000)={{0x0, {0x8}}, {0x0, {0x10}}})",
		},
		{
			In:  "test$length26(&(0x7f0000000000), 0x0)",
			Out: "test$length26(&(0x7f0000000000), 0x8)",
		},
		{
			In:  "test$length27(&(0x7f0000000000), 0x0)",
			Out: "test$length27(&(0x7f0000000000), 0x2a)",
		},
		{
			In:  "test$length28(&(0x7f0000000000), 0x0)",
			Out: "test$length28(&(0x7f0000000000), 0x2a)",
		},
		{
			In:  "test$length29(&(0x7f0000000000)={'./a\\x00', './b/c\\x00', 0x0, 0x0, 0x0})",
			Out: "test$length29(&(0x7f0000000000)={'./a\\x00', './b/c\\x00', 0xa, 0x14, 0x21})",
		},
		{
			In:  "test$length30(&(0x7f0000000000)={{{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, {'a', 'aaa', 'aaaaa', 'aaaaaa'}, &(0x7f0000000000)={'a', 'aaa', 'aaaaa', 'aaaaaa'}, &(0x7f0000000000)=&(0x7f0000000000)={'a', 'aaa', 'aaaaa', 'aaaaaa'}, 0x0}, 0x0}, 0x0, &(0x7f0000000000)=0x0, 0x0)",
			Out: "test$length30(&(0x7f0000000000)={{{0x0, 0x18, 0x1, 0x3, 0x5, 0x6}, {'a', 'aaa', 'aaaaa', 'aaaaaa'}, &(0x7f0000000000)={'a', 'aaa', 'aaaaa', 'aaaaaa'}, &(0x7f0000000000)=&(0x7f0000000000)={'a', 'aaa', 'aaaaa', 'aaaaaa'}, 0x2}, 0x4}, 0x40, &(0x7f0000000000)=0x18, 0x2)",
		},
		{
			In:  "test$offsetof0(&(0x7f0000000000)={0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})",
			Out: "test$offsetof0(&(0x7f0000000000)={0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x6, 0x8, 0x10, 0x18, 0x18, 0x20})",
		},
		{
			// If len target points into squashed argument, value is not updated.
			In: `
test$length11(&(0x7f0000000000)=ANY=[@ANYBLOB="11"], 0x42)
test$length30(&(0x7f0000000000)=ANY=[@ANYBLOB="11"], 0x42, &(0x7f0000000000)=0x43, 0x44)
`,
		},
		{
			In:  "test$length32(&(0x7f0000000000)={[0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0], {0x0}, &(0x7f0000000040)={0x0}})",
			Out: "test$length32(&(0x7f0000000000)={[0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0], {0x8}, &(0x7f0000000040)={0x8}})",
		},
		{
			In:  "test$length33(&(0x7f0000000000)={[0x0, 0x0, 0x0, 0x0], 0x0})",
			Out: "test$length33(&(0x7f0000000000)={[0x0, 0x0, 0x0, 0x0], 0x4})",
		},
		{
			In:  "test$length34(&(0x7f0000000000)={[0x0, 0x0, 0x0, 0x0], &(0x7f0000000040)=@u1=0x0})",
			Out: "test$length34(&(0x7f0000000000)={[0x0, 0x0, 0x0, 0x0], &(0x7f0000000040)=@u1=0x4})",
		},
	})
}
