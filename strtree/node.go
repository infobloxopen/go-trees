package strtree

import "fmt"

const (
	dirLeft = iota
	dirRight
)

type node struct {
	key   string
	value interface{}

	chld [2]*node
	red  bool
}

func (n *node) dot() string {
	body := ""

	// Iterate all nodes using breadth-first search algorithm.
	i := 0
	queue := []*node{n}
	for len(queue) > 0 {
		n := queue[0]
		body += fmt.Sprintf("N%d %s\n", i, n.dotString())
		if n != nil && (n.chld[0] != nil || n.chld[1] != nil) {
			// Children for current node if any always go to the end of the queue
			// so we can know their indices using current queue length.
			body += fmt.Sprintf("N%d -> { N%d N%d }\n", i, i+len(queue), i+len(queue)+1)
			queue = append(append(queue, n.chld[0]), n.chld[1])
		}

		queue = queue[1:]
		i++
	}

	return body
}

func (n *node) dotString() string {
	if n == nil {
		return "[label=\"nil\" style=filled fontcolor=white fillcolor=black]"
	}

	k := fmt.Sprintf("%q", n.key)
	if n.value != nil {
		v := fmt.Sprintf("%q", fmt.Sprintf("%#v", n.value))
		k = fmt.Sprintf("\"k: \\\"%s\\\" v: \\\"%s\\\"\"", k[1:len(k)-1], v[1:len(v)-1])
	}

	color := "fontcolor=white fillcolor=black"
	if n.red {
		color = "fillcolor=red"
	}

	return fmt.Sprintf("[label=%s style=filled %s]", k, color)
}

func (n *node) insert(key string, value interface{}, compare Compare) *node {
	if n == nil {
		return &node{
			key:   key,
			value: value,
			red:   false}
	}

	// Using fake root to get rid of corner cases with rotation right under the root.
	root := &node{key: "fake", chld: [2]*node{nil, n}, red: false}
	dir := dirLeft

	// Nodes down the path to current node. All these nodes are copies of nodes from tree.
	var (
		// Grandparent's parent.
		gp *node

		// Grandparent.
		g *node

		// Parent.
		p *node

		// Childern.
		c [2]*node
	)

	// Start with fake root.
	n = root

	// As real root is right child of fake root - go to the right from start.
	r := -1

	// Continue until keys are equal.
	for r != 0 {
		parentDir := dir
		dir = dirLeft
		if r < 0 {
			// Go to the right if current node is less then given key.
			dir = dirRight
		}

		// Propagate set of nodes.
		gp = g
		g = p
		p = n
		n = n.chld[dir]

		if n == nil {
			// If no child in the direction we go insert new red node.
			n = &node{
				key:   key,
				value: value,
				red:   true}

			c = [2]*node{nil, nil}
		} else {
			// Make copy of current node or just use copy of child node if it has been made during color flip.
			if n != c[dir] {
				n = n.fullCopy()
			}

			// Color flip case to maintain invariant that the current node is black and has at least one black child.
			if n.chld[dirLeft] != nil && n.chld[dirRight] != nil && n.chld[dirLeft].red && n.chld[dirRight].red {
				n.red = true
				c = [2]*node{
					n.chld[dirLeft].colorCopy(false),
					n.chld[dirRight].colorCopy(false)}
				n.chld = c
			} else {
				c = [2]*node{nil, nil}
			}
		}
		p.chld[dir] = n

		// Fix red violation.
		if n.red && p != nil && p.red {
			// As root is black we can't be here earlier than fake root becomes parent of grandparent.
			grandParentDir := dirLeft
			if gp.chld[dirRight] == g {
				grandParentDir = dirRight
			}

			if n == p.chld[parentDir] {
				// With single rotation if current node goes in the same direction from
				// parent as parent from grandparent.
				gp.chld[grandParentDir] = g.single(parentDir)

				// The rotation changes parent and grandparent so during next iteration
				// grandparent's parent should remain the same. Here we fix grandparent
				// to keep correct gradparent's parent.
				g = gp
			} else {
				// With double rotation if current node goes in the oposite direction.
				gp.chld[grandParentDir] = g.double(parentDir)

				// The rotation puts grandparent and parent as children of current node.
				// The nodes are copied on previous steps so we put them to to children
				// array to prevent additional coping at the next step. Also in the next
				// step grandparent's parent and grandparent iteslf make step back. So we
				// fix parent to keep correct grandparent but there is no information on
				// parent of grandparent's parent to keep corrent grandparent's parent.
				// Luckily after the rotation current node (next parent) becomes black so
				// we can't make red violation on next iteration.
				c = n.chld
				p = gp
			}
		}

		r = compare(n.key, key)
	}

	n = root.chld[dirRight]
	n.red = false
	return n
}

func (n *node) fullCopy() *node {
	return &node{
		key:   n.key,
		value: n.value,
		chld:  n.chld,
		red:   n.red}
}

func (n *node) colorCopy(color bool) *node {
	return &node{
		key:   n.key,
		value: n.value,
		chld:  n.chld,
		red:   color}
}

func (n *node) single(dir int) *node {
	nDir := 1 - dir
	s := n.chld[dir]
	n.chld[dir] = s.chld[nDir]
	s.chld[nDir] = n

	n.red = true
	s.red = false

	return s
}

func (n *node) double(dir int) *node {
	n.chld[dir] = n.chld[dir].single(1 - dir)
	return n.single(dir)
}

func (n *node) get(key string, compare Compare) (interface{}, bool) {
	for n != nil {
		r := compare(n.key, key)

		if r == 0 {
			return n.value, true
		}

		dir := dirLeft
		if r < 0 {
			dir = dirRight
		}

		n = n.chld[dir]
	}

	return nil, false
}

func (n *node) enumerate(ch chan Pair) {
	if n == nil {
		return
	}

	n.chld[dirLeft].enumerate(ch)

	ch <- Pair{Key: n.key, Value: n.value}

	n.chld[dirRight].enumerate(ch)
}
