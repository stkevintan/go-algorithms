package main

const (
	LeftIndex = iota
	RightIndex
)

type SplayNode struct {
	Child  [2]*SplayNode //子节点
	Parent *SplayNode    //父节点
	Value  int           //当前节点值
	Size   int           //当前树的节点数
	Sum    int           //当前子树的值
	Mark   bool
}

func getSize(node *SplayNode) int {
	if node == nil {
		return 0
	}
	return node.Size
}
func (node *SplayNode) PushUp() {
	node.Size = 1 + getSize(node.Child[LeftIndex]) + getSize(node.Child[RightIndex])
}

func (node *SplayNode) PushDown() {
	if !node.Mark {
		return
	}
	node.Child[LeftIndex], node.Child[RightIndex] = node.Child[RightIndex], node.Child[LeftIndex]
	if node.Child[LeftIndex] != nil {
		node.Child[LeftIndex].Mark = !node.Child[LeftIndex].Mark
	}
	if node.Child[RightIndex] != nil {
		node.Child[RightIndex].Mark = !node.Child[RightIndex].Mark
	}
	node.Mark = false
}

func NewSplayNode(val int) *SplayNode {
	return &SplayNode{Value: val, Size: 1}
}

func switchSide(side int) int {
	return 1 ^ side
}

func (_ *SplayTree) PushUp(node *SplayNode) {
	return
}
func (_ *SplayTree) PushDown(node *SplayNode) {
	return
}

func (node *SplayNode) getSide() int {
	parent := node.Parent
	if parent == nil {
		return LeftIndex
	}
	if parent.Child[RightIndex] == node {
		return RightIndex
	}
	return LeftIndex
}

type SplayTree struct {
	Root *SplayNode
}

//Link  把from节点连接到to节点上, side == {LeftIndex, RightIndex}
func (tree *SplayTree) Link(from *SplayNode, to *SplayNode, pos int) {
	if to != nil {
		to.Child[pos] = from
	}
	if from != nil {
		from.Parent = to
	}
}

func (tree *SplayTree) Rotate(node *SplayNode) {
	if node == tree.Root {
		return
	}
	parent := node.Parent

	parentSide := parent.getSide()
	mySide := node.getSide()
	// 父亲视角： 先处理老，再处理小，最后再处理自己
	// 让儿子直接接替自己
	tree.Link(node, parent.Parent, parentSide)

	//让孙子来代替自己的原来儿子的位置，（注意选择另外一边的孙子，因为另外一边孙子的位置将来会被自己代替
	tree.Link(node.Child[switchSide(mySide)], parent, mySide)

	//认自己儿子做父 自己去代替之前孙子的位置
	tree.Link(parent, node, switchSide(mySide))

	//这时候parent应该在node的下面，先更新子节点
	parent.PushUp()
	//node 本来也需要更新的，但是正常情况下本函数都是一层层往上调用的（详见Splay）所以现在更新有点浪费
	// node.PushUp()
}

// Splay 把from节点伸展到to节点之下
func (tree *SplayTree) Splay(from, to *SplayNode) {
	for from.Parent != to {
		y := from.Parent
		z := y.Parent
		if z == to {
			tree.Rotate(from)
			break
		}

		yPos := y.getSide()
		zPos := z.getSide()

		//方向相同，则先旋转父亲，再旋转自己
		if yPos == zPos {
			tree.Rotate(y)
		} else {
			// 否则旋转自己两次
			tree.Rotate(from)
		}
		tree.Rotate(from)
	}
	if to == nil {
		tree.Root = from
	}
	// 更新最后一个Rotate没有up过的节点
	from.PushUp()
}

func (tree *SplayTree) Find(val int) (*SplayNode, *SplayNode) {
	if tree.Root == nil {
		return nil, nil
	}
	curNode := tree.Root
	var preNode *SplayNode
	for curNode != nil {
		if curNode.Value > val {
			preNode = curNode
			curNode = curNode.Child[LeftIndex]
		} else if curNode.Value < val {
			preNode = curNode
			curNode = curNode.Child[RightIndex]
		} else {
			// 找到了，我们可以把它旋转到根，保证后续查找高效
			tree.Splay(curNode, nil)
			return curNode, preNode
		}
	}
	return nil, preNode
}

func (tree *SplayTree) Kth(root *SplayNode, k int) *SplayNode {
	if root == nil {
		return nil
	}
	root.PushDown()

	leftSize := getSize(root.Child[LeftIndex])
	if leftSize+1 == k {
		return root
	}
	if leftSize >= k {
		return tree.Kth(root.Child[LeftIndex], k)
	}
	return tree.Kth(root.Child[RightIndex], k-leftSize-1)
}

func (tree *SplayTree) Flip(l, r int) {
	if l >= r {
		return
	}
	x := tree.Kth(tree.Root, l-1)
	y := tree.Kth(tree.Root, r+1)
	tree.Splay(x, nil)
	tree.Splay(y, x)
	y.Child[LeftIndex].Mark = !y.Child[LeftIndex].Mark
}

//Cut cut l to r, and insert to c
func (tree *SplayTree) Cut(l, r, c int) {
	// 把第l -1个点x移到 root， 把第r + 1个点y移到x下。那么y的左子树就是介于l ~ r的点
	x := tree.Kth(tree.Root, l-1)
	y := tree.Kth(tree.Root, r+1)
	tree.Splay(x, nil)
	tree.Splay(y, x)
	cutPart := y.Child[LeftIndex]
	y.PushUp()
	x.PushUp()
	// 把第c个点移到root，把第c+1个点移到c下。那么c+1的左子树一定没有值，直接把cutPart插进去
	x = tree.Kth(tree.Root, c)
	y = tree.Kth(tree.Root, c+1)
	tree.Splay(x, nil)
	tree.Splay(y, x)
	tree.Link(cutPart, y, LeftIndex)
}

//Insert 如果存在相同值的节点，插入失败。如果没找到，直接插到上一个节点的对应位置上即可
func (tree *SplayTree) Insert(val int) bool {
	node := NewSplayNode(val)
	if tree.Root == nil {
		tree.Root = node
		return true
	}
	curNode, preNode := tree.Find(val)

	if curNode != nil {
		return false
	}

	if preNode.Value > val {
		tree.Link(node, preNode, RightIndex)
	} else {
		tree.Link(node, preNode, LeftIndex)
	}
	tree.Splay(node, nil)
	return true
}

func (tree *SplayTree) Replace(from *SplayNode, to *SplayNode) {
	if tree.Root == from {
		tree.Root = to
	} else if from.getSide() == LeftIndex {
		from.Parent.Child[LeftIndex] = to
	} else {
		from.Parent.Child[RightIndex] = to
	}
	if to != nil {
		to.Parent = from.Parent
	}
}

//Delete 把要删除的节点Splay到root上，然后找出左子树最小的那个Splay上来作为新的root（右子树最小的那个也行）
//或者还有一条思路就是把父节点挂在root，子节点挂在root下面，那么这个节点就一定孤零零地挂在root的另外一个儿子上，直接把这个儿子删掉
func (tree *SplayTree) Delete(val int) bool {
	curNode, _ := tree.Find(val)
	if curNode == nil {
		return false
	}
	// 如果不是根节点，那么把它移动到根节点
	if curNode != tree.Root {
		tree.Splay(curNode, nil)
	}
	if curNode.Child[RightIndex] == nil {
		tree.Replace(curNode, curNode.Child[LeftIndex])
	} else if curNode.Child[LeftIndex] == nil {
		tree.Replace(curNode, curNode.Child[RightIndex])
	} else {
		//Find the Max on the left subTree
		maxNode := curNode.Child[LeftIndex]
		for {
			rightNode := maxNode.Child[RightIndex]
			if rightNode == nil {
				break
			}
			maxNode = rightNode
		}
		if maxNode.Parent != curNode {
			tree.Replace(maxNode, maxNode.Child[LeftIndex])
			maxNode.Child[LeftIndex] = curNode.Child[LeftIndex]
			maxNode.Child[LeftIndex].Parent = maxNode
		}
		baseNode := maxNode.Parent
		tree.Replace(curNode, maxNode)
		maxNode.Child[RightIndex] = curNode.Child[RightIndex]
		maxNode.Child[RightIndex].Parent = maxNode
		tree.Splay(baseNode, nil)
	}
	return true
}

func main() {

}
