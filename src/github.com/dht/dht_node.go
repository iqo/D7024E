package dht

import (
	"encoding/hex"
	"fmt"
)

type Contact struct {
	ip   string
	port string
}

type DHTNode struct {
	nodeId      string
	successor   *DHTNode
	predecessor *DHTNode
	contact     Contact
	fingers     *DHTFingers
	transport   *Transport
}

func makeDHTNode(nodeId *string, ip string, port string) *DHTNode {
	dhtNode := new(DHTNode)
	dhtNode.contact.ip = ip
	dhtNode.contact.port = port

	if nodeId == nil {
		genNodeId := generateNodeId()
		dhtNode.nodeId = genNodeId
	} else {
		dhtNode.nodeId = *nodeId
	}

	dhtNode.successor = nil
	dhtNode.predecessor = nil

	dhtNode.fingers = new(DHTFingers)
	dhtNode.fingers.nodefingerlist = [bits]*DHTNode{}

	return dhtNode
}

func (dhtNode *DHTNode) addToRing(newDHTNode *DHTNode) {
	//join info / stabilize S.5 -S.8
	if dhtNode.predecessor == nil && dhtNode.successor == nil {

		newDHTNode.predecessor = dhtNode

		newDHTNode.successor = dhtNode

		dhtNode.successor = newDHTNode

		dhtNode.predecessor = newDHTNode

		updateFingers(dhtNode)

		//dhtNode.fingers.nodefingerlist = init_finger_table(dhtNode)

		//newDHTNode.fingers.nodefingerlist = init_finger_table(newDHTNode)

	} else if between([]byte(dhtNode.nodeId), []byte(dhtNode.successor.nodeId), []byte(newDHTNode.nodeId)) {

		dhtNode.successor.predecessor = newDHTNode

		newDHTNode.successor = dhtNode.successor

		dhtNode.successor = newDHTNode

		newDHTNode.predecessor = dhtNode

		newDHTNode.fingers.nodefingerlist = init_finger_table(newDHTNode)
		updateFingers(dhtNode)

	} else {
		dhtNode.successor.addToRing(newDHTNode)
	}
	//updateFingers(newDHTNode)
	updateFingers(dhtNode)
}

func (dhtNode *DHTNode) lookup(key string) *DHTNode {
	if between([]byte(dhtNode.nodeId), []byte(dhtNode.successor.nodeId), []byte(key)) {
		if dhtNode.nodeId == key {

			return dhtNode

		} else {

			return dhtNode.successor
		}
	} else {

		return dhtNode.successor.lookup(key)
	}
}

/*func (dhtNode *DHTNode) lookup(key string) *DHTNode {
	if between([]byte(dhtNode.nodeId),[]byte(dhtNode.successor.nodeId), []byte(key)){
		return dhtNode.successor.predecessor
	} else{

		return dhtNode.successor.lookup(key)
	}
}*/

func (dhtNode *DHTNode) acceleratedLookupUsingFingers(key string) *DHTNode {
	// TODO

	for i := len(dhtNode.fingers.nodefingerlist); i > 0; i-- {
		if between([]byte(dhtNode.nodeId), []byte(dhtNode.fingers.nodefingerlist[i-1].nodeId), []byte(key)) {
			fmt.Println(key, "ligger mellan ", dhtNode.nodeId, "och tillhörande finger ", dhtNode.fingers.nodefingerlist[i-1].nodeId)

		} else {
			return dhtNode.fingers.nodefingerlist[i-1].acceleratedLookupUsingFingers(key)
		}
	}
	return dhtNode // XXX This is not correct obviously
}

func (dhtNode *DHTNode) responsible(key string) bool {
	// TODO
	return false
}

func (dhtNode *DHTNode) printRing() {
	//fmt.Println(dhtNode.nodeId)
	for i := dhtNode; i != dhtNode.predecessor; i = i.successor {
		fmt.Println(i.nodeId)
	}
	fmt.Println(dhtNode.predecessor.nodeId)
	// TODO
}

func (dhtNode *DHTNode) testCalcFingers(m int, bits int) {
	idBytes, _ := hex.DecodeString(dhtNode.nodeId)
	fingerHex, _ := calcFinger(idBytes, m, bits)
	fingerSuccessor := dhtNode.lookup(fingerHex)
	fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
	fmt.Println("successor    " + fingerSuccessor.nodeId)

	dist := distance(idBytes, fingerSuccessorBytes, bits)
	fmt.Println("distance     " + dist.String())
}

func (dhtNode *DHTNode) findSuccessor(node *DHTNode) *DHTNode {
	//psudokod S.5 i boken
	predecessorNode := dhtNode.findPredecessor(node)
	return predecessorNode.successor
}

func (dhtNode *DHTNode) findPredecessor(node *DHTNode) *DHTNode {
	//psudokod S.5 i boken
	successorNode := dhtNode
	return successorNode
}

func (dhtNode *DHTNode) printTable() {
	for i := 0; i < len(dhtNode.fingers.nodefingerlist); i++ {
		fmt.Println("Node  ", dhtNode.nodeId, "finger ", i+1, " poitns at ", dhtNode.fingers.nodefingerlist[i])
	}
}

func updateFingers(node *DHTNode) {
	for i := node; i != node.successor; i = i.predecessor {
		i.fingers.nodefingerlist = init_finger_table(i)
	}
}
