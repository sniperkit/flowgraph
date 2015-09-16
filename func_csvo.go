package flowgraph

import (
	"encoding/csv"
	"io"
	"os"
)      			


func csvoRdy (n *Node) bool {
	if n.Aux == nil { return false }
	
	a := n.Srcs
	r := n.Aux.(csvState).csvreader
	h := n.Aux.(csvState).header

	if n.Aux== nil { return false }

	if n.Aux.(csvState).record==nil {
		record,err := r.Read()
		if err == io.EOF {
			os.Exit(0)
			return false
		} else {
			check(err)
			n.Aux = csvState{r, h, record}
		}
	}

	record := n.Aux.(csvState).record
	header := n.Aux.(csvState).header

	for i := range a {
		j := find(a[i].Name, header)
		if j>= 0 {
			if !a[i].SrcRdy(n) {
				if record[j]!="*" {
					return false
				} else {
					a[i].NoOut = true
				}
			}
		} else {
			n.LogError("Named output missing from .csv file:  %s\n", a[i].Name)
			os.Exit(1)
		}
	}
	return true
}

func csvoFire (n *Node) {	 
	a := n.Srcs

	record := n.Aux.(csvState).record
	r := n.Aux.(csvState).csvreader
	header := n.Aux.(csvState).header

	l := len(a)
	if l>len(record) { l = len(record) }
	for i:=0; i<l; i++ {
		j := find(a[i].Name, header)
		if record[j]!="*" {
			v := ParseDatum(record[i])
			if !EqualsTest(n, v, a[i].Val) {
				n.LogError("expected=%T(%v) (0x%x), actual=%T(%v) (0x%x)", v, v, v, a[i].Val, a[i].Val, a[i].Val)	
			}
		}
	}

	n.Aux = csvState{csvreader:r, header:header}
	
}

// FuncCSVO reads a vector of input data values from a Reader.
func FuncCSVO(a []Edge, r io.Reader) Node {

	var ap []*Edge
	for i := range a {
		ap = append(ap, &a[i])
	}

	node := MakeNode("csvo", ap, nil, csvoRdy, csvoFire)
	r2 := csv.NewReader(r)

	// save headers
	headers, err := r2.Read()
	check(err)
	node.Aux = csvState{csvreader:r2, header:headers}

	return node
	
}
	
