package graph

import (
	"config"
	"log"
)

const (
	// Set a buffer size for all channels between nodes
	CHANNEL_SIZE = 10
)

var (
	// channels is a map that stores all channels between the nodes
	// The name of a channel is "input_node_name-output_node_name"
	channels = map[string](chan Message){}

	// Send a message to the source channel to trigger the processing
	sourceChannel = make(chan Message)

	// Block on the drain channel to wait for the drain node
	drainChannel = make(chan Message)
)

// Graph is a map that contains all nodes
type Graph struct {
	Nodes map[string]*Node
}

func (g *Graph) InitAllNode() {
	for _, node := range g.Nodes {
		log.Printf("Create goroutine for node (%s)\n", node.Name)
		// Create a goroutine
		go node.Run()
	}
}

func (g *Graph) Start() {
	sourceChannel <- Message{1}
}

func (g *Graph) WaitEnd() {
	<-drainChannel
	log.Println("Graph processing ends")
}

type Message struct {
	Quantity int
}

// ConstructGraph constructs a new graph processing engine
// based on a graph configuration
func ConstructGraph(graphConfig *config.GraphConfig) *Graph {
	// Create a new empty graph
	graph := &Graph{
		Nodes: make(map[string]*Node),
	}

	// Parse the topology to construct the graph
	for nodeName, nodeConfig := range graphConfig.Topology {
		// Check whether this node has been created

		node, ok := graph.Nodes[nodeName]
		if !ok {
			// If not, create a new one
			node = &Node{
				Name:     nodeName,
				IsSource: false,
				IsDrain:  false,
				Inputs:   make(map[string]int),
				Outputs:  make(map[string]int),
			}
			graph.Nodes[nodeName] = node
		}

		//Adding dynamically read paths

		//If node is source, set the boolean flag and configure only outputs
		//If node is drain, set the boolean flag and configure only inputs
		//Otherwise, configure no flags but both inputs and ouputs
		
		switch node.Name {
		case "source":
			node.IsSource = true
			node.Outputs = nodeConfig.Outputs

		case "drain":
			node.IsDrain = true
			node.Inputs = nodeConfig.Inputs

		default:
			node.Inputs = nodeConfig.Inputs
			node.Outputs = nodeConfig.Outputs
		}
		for ipNodeName := range node.Inputs {
			channels[ipNodeName+"-"+node.Name] = make(chan Message)
		}
		for opNodeName := range node.Outputs {
			channels[node.Name+"-"+opNodeName] = make(chan Message)
		}
	}

	return graph
}
