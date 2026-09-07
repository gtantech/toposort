package graph

import (
	"maps"
)

type Graph[V comparable, E any] interface {
	OutgoingVertices(vertex V) func(yield func(V) bool)
	OutgoingVerticesWithEdge(vertex V) func(yield func(V, E) bool)
	IncomingVertices(vertex V) func(yield func(V) bool)
	IncomingVerticesWithEdge(vertex V) func(yield func(V, E) bool)
	Vertices() func(yield func(V) bool)
	AddVertex(v V)
	RemoveVertex(v V)
	AddEdge(value E, origin V, destination V)
	GetEdgeValue(origin V, destination V) (E, bool)
	RemoveEdge(origin V, destination V)
}

var _ Graph[string, string] = (*dag[string, string])(nil)

type dag[V comparable, E any] struct {
	outgoingVertices  map[V]map[V]E //maps the origin vertex, then destination vertex to the edge value
	incomingVerticies map[V]map[V]E //maps the destination vertex, then origin vertex to the edge value
	uniqueVerticies   map[V]struct{}
}

// IncomingVerticesWithEdge implements [Graph].
func (d *dag[V, E]) IncomingVerticesWithEdge(v V) func(yield func(V, E) bool) {
	destinations := d.incomingVerticies[v]

	return maps.All(destinations)
}

// OutgoingVerticesWithEdge implements [Graph].
func (d *dag[V, E]) OutgoingVerticesWithEdge(v V) func(yield func(V, E) bool) {
	destinations := d.outgoingVertices[v]

	return maps.All(destinations)
}

func New[V comparable, E any]() Graph[V, E] {
	return &dag[V, E]{outgoingVertices: make(map[V]map[V]E), incomingVerticies: make(map[V]map[V]E), uniqueVerticies: make(map[V]struct{})}
}

func (d *dag[V, E]) AddVertex(v V) {
	if _, ok := d.uniqueVerticies[v]; !ok {
		d.uniqueVerticies[v] = struct{}{}
	}
}

func (d *dag[V, E]) RemoveVertex(v V) {
	delete(d.uniqueVerticies, v)
	//ensure that if this vertex v is an origin vertex, any edge relating to it is deleted
	for affectedDestionationVertex := range d.outgoingVertices[v] {
		delete(d.incomingVerticies[affectedDestionationVertex], v)
	}
	//ensure that if this vertex v is an destination vertex, any edge relating to it is deleted
	for affectedOriginVertex := range d.incomingVerticies[v] {
		delete(d.outgoingVertices[affectedOriginVertex], v)
	}
	delete(d.incomingVerticies, v)
	delete(d.outgoingVertices, v)
}

func (d *dag[V, E]) RemoveEdge(origin V, destination V) {
	outgoing, ok := d.outgoingVertices[origin]
	if ok {
		//origin vertex exists
		delete(outgoing, destination)
	}
}

func (d *dag[V, E]) GetEdgeValue(origin V, destination V) (E, bool) {
	destinations, ok := d.outgoingVertices[origin]
	var zero E
	if !ok {
		return zero, false
	}

	//origin vertex exists
	value, ok := destinations[destination]

	if !ok {
		return zero, false
	}
	return value, true
}

func (d *dag[V, E]) AddEdge(value E, origin V, destination V) {
	if _, ok := d.outgoingVertices[origin]; !ok {
		d.outgoingVertices[origin] = make(map[V]E)
	}
	d.outgoingVertices[origin][destination] = value
	if _, ok := d.incomingVerticies[destination]; !ok {
		d.incomingVerticies[destination] = make(map[V]E)
	}
	d.incomingVerticies[destination][origin] = value
	d.uniqueVerticies[origin] = struct{}{}
	d.uniqueVerticies[destination] = struct{}{}
}

// OutgoingVertices implements [graph.Graph].
func (d *dag[V, E]) OutgoingVertices(v V) func(yield func(V) bool) {
	destinations := d.outgoingVertices[v]

	return maps.Keys(destinations)
}

// IncomingVertices implements [graph.Graph].
func (d *dag[V, E]) IncomingVertices(v V) func(yield func(V) bool) {
	origins := d.incomingVerticies[v]

	return maps.Keys(origins)
}

// Vertices implements [graph.Graph].
func (d *dag[V, E]) Vertices() func(yield func(V) bool) {
	return maps.Keys(d.uniqueVerticies)
}
