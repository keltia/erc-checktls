// kpi.go

/*
This file have func for computing the TLS scores
 */
package main

const (
	protocolTable = map[string]float32{
		"SSL2.0": 0.0,
		"SSL3.0": 80.0,
		"TLS1.0": 90.0,
		"TLS1.1": 95.0,
		"TLS1.2": 100.0,
	}
)

// getProtoScore merge proto + version and retrieve the score
func getProtoScore(p *LabsProtocol) (score float32) {
	proto := p.Name + p.Version
	score = protocolTable[proto]
	return
}

// protocolScore compute the relative ProcotolScore KPI
func protocolScore(r *LabsReport) (score int) {
	var (
		lowest, highest float32
	)

	// check for SSL2.0
	protoList := r.Endpoints[0].Details.Protocols

	if len(protoList) == 1 {
		score = getProtoScore(protoList[0])
		return
	}

	for _, p := range protoList {
		score = getProtoScore(p)
		if score < lowest {
			lowest = score
		}
		if score > highest {
			highest = score
		}
	}
	score = (lowest + highest) / 2
	return
}
