package snowflake

import (
	"bytes"
	"math/rand"
	"net"
)

func createNodeId() int64 {
	var nodeId int64

	interfaces, err := net.Interfaces()
	if err != nil {
		// handle error
	}

	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
			// Use the MAC address as the node ID
			nodeId = int64(i.HardwareAddr[5])
			break
		}
	}

	// If no MAC address is found, generate a random node ID
	if nodeId == 0 {
		nodeId = int64(rand.Uint64())
	}

	// Ensure the node ID is within the allowed range
	nodeId = nodeId & maxNodeId

	return nodeId
}
