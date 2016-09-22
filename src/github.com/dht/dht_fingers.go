package dht

import (
	"encoding/hex"
	//"fmt"
)

const bits int = 9

type DHTFingers struct {
	nodefingerlist [bits]*DHTNode
}

func init_finger_table(node *DHTNode) [bits]*DHTNode {
	var tempFingerList [bits]*DHTNode
	for i := 0; i < bits; i++ {
		//fmt.Println("init for")
		nodeIdBits, _ := hex.DecodeString(node.nodeId)
		calculatedFinger, _ := calcFinger(nodeIdBits, i+1, bits)

		if calculatedFinger == "" {
			calculatedFinger = "00"
		}
		nodeSucc := node.lookup(calculatedFinger)
		tempFingerList[i] = nodeSucc
	}
	return tempFingerList
}
