// +build !amd64

package dltree

func treeRawGet(t *Tree, key string, size int) (interface{}, bool) {
	if t == nil {
		return nil, false
	}

	return t.root.get(key, size)
}
