package main

import (
	"bytes"
	"fmt"
	"sort"
	"text/tabwriter"

	"golang.org/x/tools/benchmark/parse"
)

func Compare(a, b *Benchmark) ([]byte, error) {
	sa, err := parse.ParseSet(bytes.NewReader([]byte(a.Content)))
	if err != nil {
		return nil, fmt.Errorf("cannot parse a: %s", err)
	}
	sb, err := parse.ParseSet(bytes.NewReader([]byte(b.Content)))
	if err != nil {
		return nil, fmt.Errorf("cannot parse b: %s", err)
	}
	cmps, _ := Correlate(sa, sb)

	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 0, 0, 5, ' ', 0)

	sort.Sort(ByParseOrder(cmps))

	for _, cmp := range cmps {
		if !cmp.Measured(parse.NsPerOp) {
			continue
		}
		delta := cmp.DeltaNsPerOp()
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			cmp.Name(),
			formatNs(cmp.Before.NsPerOp),
			formatNs(cmp.After.NsPerOp),
			delta.Percent())
	}

	tw.Flush()
	return buf.Bytes(), nil
}
