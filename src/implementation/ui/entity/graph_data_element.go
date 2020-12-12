package entity

import "fmt"

type GraphDataElements []GraphDataElement

func (g GraphDataElements) MinSize() (min int) {
	for _, v := range g {
		if min == 0 {
			min = v.SizeBytes
		}
		if v.SizeBytes < min {
			min = v.SizeBytes
		}
	}
	return
}

func (g GraphDataElements) MaxSize() (max int) {
	for _, v := range g {
		if max == 0 {
			max = v.SizeBytes
		}
		if v.SizeBytes > max {
			max = v.SizeBytes
		}
	}
	return
}

type GraphDataElement struct {
	SourceHost    string
	SourceService string
	TargetHost    string
	TargetService string
	PacketsCnt    int
	SizeBytes     int
}

func (g GraphDataElement) SourceServiceID() string {
	return fmt.Sprintf("%s-%s", g.SourceHost, g.SourceService)
}

func (g GraphDataElement) TargetServiceID() string {
	return fmt.Sprintf("%s-%s", g.TargetHost, g.TargetService)
}

func (g GraphDataElement) EdgeId() string {
	return fmt.Sprintf("%s-%s-%s-%s", g.SourceHost, g.SourceService, g.TargetHost, g.TargetService)
}
