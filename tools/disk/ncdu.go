package disk

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/douban/sa-tools-go/libs/heapq"
)

type NcduNode struct {
	Name  string
	Asize uint64
	Dsize uint64
	Ino   uint64
}

func TopHugeDirsFromNcdu(topN uint64, raw []byte, parentPath string, maxDepth, depth uint64) (map[string]uint64, uint64, error) {
	var rawData interface{}
	json.Unmarshal(raw, &rawData)
	switch rawData.(type) {
	case map[string]interface{}:
		var node NcduNode
		err := json.Unmarshal(raw, &node)
		if err != nil {
			return nil, 0, fmt.Errorf("unmashal node error")
		}
		if depth > maxDepth {
			return nil, node.Dsize, nil
		}
		path := filepath.Join(parentPath, node.Name)
		return map[string]uint64{path: node.Dsize}, 0, nil
	case []interface{}:
		var nodes []json.RawMessage
		err := json.Unmarshal(raw, &nodes)
		if err != nil {
			return nil, 0, fmt.Errorf("unmashal first node error")
		}
		var node NcduNode
		err = json.Unmarshal(nodes[0], &node)
		if err != nil {
			return nil, 0, fmt.Errorf("unmashal first node error")
		}

		subNodes := nodes[1:]
		path := filepath.Join(parentPath, node.Name)

		var size uint64
		var ret map[string]uint64
		if depth >= maxDepth {
			size = node.Dsize
		} else {
			ret = make(map[string]uint64)
		}

		for _, subnode := range subNodes {
			r, s, err := TopHugeDirsFromNcdu(topN, subnode, path, maxDepth, depth+1)
			if err != nil {
				return nil, 0, err
			}
			if depth >= maxDepth {
				size += s
			} else {
				for k, v := range r {
					ret[k] = v
				}
				ret = heapq.Nlargest(topN, ret)
			}
		}
		if depth == maxDepth {
			return map[string]uint64{path: size}, size, nil
		} else {
			return ret, size, nil
		}
	}
	return nil, 0, fmt.Errorf("ncdu data type wrong: %s", string(raw))
}
