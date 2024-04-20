package entity

type VisitStatus int

const (
	VisitStatusUnVisited VisitStatus = iota
	VisitStatusVisited
	VisitStatusUnVisitable VisitStatus = iota
)

type Edge[T any] struct {
	DestinationNodeId int
	Cost              T
}

type Node[T any] struct {
	Id             int
	OriginalNode   LocationDetails
	VisitStatus    VisitStatus
	GeoCoordinates GeoCoordinates
	OverheadCost   T
}

type Graph[T any] struct {
	Nodes []Node[T]
	Adj   [][]Edge[T]
}
