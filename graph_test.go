package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestModuleGraph_Parse(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "full",
			args: args{
				bytes.NewReader([]byte(`github.com/poloxue/testmod golang.org/x/text@v0.3.2
github.com/poloxue/testmod rsc.io/quote/v3@v3.1.0
github.com/poloxue/testmod rsc.io/sampler@v1.3.1
golang.org/x/text@v0.3.2 golang.org/x/tools@v0.0.0-20180917221912-90fa682c2a6e
rsc.io/quote/v3@v3.1.0 rsc.io/sampler@v1.3.0
rsc.io/sampler@v1.3.1 golang.org/x/text@v0.0.0-20170915032832-14c0d48ead0c
rsc.io/sampler@v1.3.0 golang.org/x/text@v0.0.0-20170915032832-14c0d48ead0c`))},
			want: []byte(`digraph {
node [shape=box];
1 [label="github.com/poloxue/testmod"];
2 [label="golang.org/x/text@v0.3.2"];
3 [label="rsc.io/quote/v3@v3.1.0"];
4 [label="rsc.io/sampler@v1.3.1"];
5 [label="golang.org/x/tools@v0.0.0-20180917221912-90fa682c2a6e"];
6 [label="rsc.io/sampler@v1.3.0"];
7 [label="golang.org/x/text@v0.0.0-20170915032832-14c0d48ead0c"]
1 -> 2;
1 -> 3;
1 -> 4;
2 -> 5;
3 -> 6;
4 -> 7;
6 -> 7;
}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moduleGraph := NewModuleGraph(tt.args.reader)
			moduleGraph.Parse()
			for k, v := range moduleGraph.Mods {
				fmt.Println(v, k)
			}

			for k, v := range moduleGraph.Dependencies {
				fmt.Println(k)
				fmt.Println(v)
				fmt.Println()
			}
		})
	}
}
