package flowgraph

import (
)

func arbit_rdy (n *Node) bool {
	return (n.Srcs[0].Rdy || n.Srcs[1].Rdy) && n.Dsts[0].Rdy
}

// Arbiter goroutine
func FuncArbit(a, b, x Edge) {

	node := NewNode("arbit", []*Edge{&a, &b}, []*Edge{&x}, arbit_rdy)

	a_last := false

	for {
		node.Tracef("a.Rdy,b.Rdy %v,%v  x.Rdy %v\n", a.Rdy, b.Rdy, x.Rdy);

		if node.Rdy() {
			node.Tracef("writing x.Data  and either a.Ack or b.Ack\n")
			if(a.Rdy && !b.Rdy || a.Rdy && !a_last) {
				a_last = true
				x.Val = a.Val
				node.TraceVals()
				a.Rdy = false
				a.Ack <- true
				node.Tracef("done writing x.Data and a.Ack\n")
			} else if (b.Rdy) {
				a_last = false
				x.Val = b.Val
				node.TraceVals()
				b.Rdy = false
				b.Ack <- true
				node.Tracef("done writing x.Data and b.Ack\n")
			}
			x.Data <- x.Val
			x.Rdy = false
		}

		node.Select()

	}

}