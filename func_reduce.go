package flowgraph

import (
)

// FuncReduce reduces a stream of data into a single empty interface.
func FuncReduce(a,x Edge, reducer func(n *Node, datum,collection interface{}) interface{}, freerun bool) Node {

        var reduceFreerunRdy = func (n *Node) bool {
		a := n.Srcs[0]
		x := n.Dsts[0]
		return a.SrcRdy(n) || x.DstRdy(n)
	}

	var reduceFreerunFire = func (n *Node) {
	        if a.SrcRdy(n) {
  		        n.Aux = reducer(n, a.SrcGet(), n.Aux)
		}
		if x.DstRdy(n) {
		        x.DstPut(n.Aux)
		}
	}

	var reduceSteppedFire = func (n *Node) {
	        n.Aux = reducer(n, a.SrcGet(), n.Aux)
	        x.DstPut(n.Aux)
	}

	var reduceRdy func (n *Node) bool = nil
	if freerun {
	        reduceRdy = reduceFreerunRdy
	}

        var reduceFire func (n *Node) = reduceSteppedFire
	if freerun {
	        reduceFire = reduceFreerunFire
	}
	
	node := MakeNode("reduce", []*Edge{&a}, []*Edge{&x}, reduceRdy, reduceFire)
	node.Aux = make([]string, 0)
	return node
}
