package graph

import "log"

// Node is a concurrent processing unit
type Node struct {
	Name     string
	IsSource bool           // is it a source node?
	IsDrain  bool           // is it a drain node?
	Inputs   map[string]int // the inputs information map[node_name](data quantity)
	Outputs  map[string]int // the outputs information map[node_name](data quantity)

}

func (n *Node) Run() {

	//Common messages for each node - print when initialized and print when starting to process data
	log.Printf("Node (%s): Initiated\n", n.Name)
	log.Printf("Node (%s): ----- Start processing data -----\n", n.Name)

	switch n.Name {
	case "source":

		//trigger source channel
		<-sourceChannel
		if n.IsSource == true {
			for opNodeName := range n.Outputs {
				var Qty = n.Outputs[opNodeName]
				channels[n.Name+"-"+opNodeName] <- Message{Qty}
				log.Printf("Node (%s): Send <%d> to (%s)\n", n.Name, Qty, opNodeName)
			}
		}
	default:
		//Recieve dat
		for ipNodeName := range n.Inputs {
			RcdQty := <-channels[ipNodeName+"-"+n.Name]
			log.Printf("Node (%s): Recieved <%d> from (%s)\n", n.Name, RcdQty, ipNodeName)
		}

		for opNodeName := range n.Outputs {
			var Qty = n.Outputs[opNodeName]
			channels[n.Name+"-"+opNodeName] <- Message{Qty}
			log.Printf("Node (%s): Send <%d> to (%s)\n", n.Name, Qty, opNodeName)
		}

	case "drain":

		if n.IsDrain == true {
			for ipNodeName := range n.Inputs {
				RcdQty := <-channels[ipNodeName+"-"+n.Name]
				log.Printf("Node (%s): Recieved <%d> from (%s)\n", n.Name, RcdQty, ipNodeName)
			}
		}
		//trigger drain when done receiving
		drainChannel <- Message{1}

	}
}

/*
	To be completed:
	1. For source node, block on source channel and wait for triggering;
	2. For drain node, once it finishes, trigger the drain channel;
	3. All nodes need to block on input channels, and to start processing once it collects enough inputs,
	  	then sends messages to output channels;
	4. Print information such as:
		log.Printf("Node (%s): Receive <%d> from (%s)\n", n.Name, msg.Quantity, inputNodeName)
		log.Printf("Node (%s): ----- Start processing data ----- \n", n.Name)
		log.Printf("Node (%s): Send <%d> to (%s)\n", n.Name, quantity, outputNodeName)
*/
