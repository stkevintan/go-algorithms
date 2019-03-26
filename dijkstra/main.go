package main

import (
	"container/heap"
	"fmt"
)

type Edge struct {
	v    int
	w    int
	next int
}

type Vex struct {
	dis int
	u   int
}

type PriorityQueue []*Vex

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	if pq[j].dis == -1 {
		return false
	}
	if pq[i].dis == -1 {
		return true
	}

	return pq[i].dis > pq[j].dis
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Vex)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

var pos = 0
var head [22222]int
var edge [111]Edge
var dis [111]int
var visit [111]bool

func addEdge(u, v, w int) {
	edge[pos] = Edge{v, w, head[u]}
	head[u] = pos
	pos++
}

func initial(n int) {
	pos = 0
	for i := 1; i <= n; i++ {
		dis[i] = -1
		visit[i] = false
		head[i] = -1
	}
}

func dijkstra(n int) {
	dis[1] = 0
	pq := make(PriorityQueue, 1)
	pq[0] = &Vex{dis[1], 1}
	heap.Init(&pq)
	for pq.Len() > 0 {
		top := (pq.Pop()).(*Vex)
		u := top.u
		if visit[u] {
			continue
		}
		visit[u] = true
		for i := head[u]; i != -1; i = edge[i].next {
			v := edge[i].v
			if visit[v] {
				continue
			}
			tmp := dis[u] + edge[i].w

			if dis[v] == -1 || tmp < dis[v] {
				dis[v] = tmp
				pq.Push(&Vex{dis[v], v})
			}
		}
	}
}

func main() {
	var m, n int
	for {
		if num, _ := fmt.Scanf("%d%d", &n, &m); num != 2 || n == 0 && m == 0 {
			break
		}
		initial(n)
		for i := 0; i < m; i++ {
			var u, v, w int
			fmt.Scanf("%d%d%d", &u, &v, &w)
			addEdge(u, v, w)
			addEdge(v, u, w)
		}
		dijkstra(n)
		fmt.Println(dis[n])
	}
}
