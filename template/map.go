package template

import (
	"fmt"

	"github.com/flosch/pongo2"
)

// Map - Template map.
type Map map[string]MapNode

// MapNode - Template map node.
type MapNode interface {
	Execute(Context) (interface{}, error)
}

// ParseMap - Parses map from context map.
func ParseMap(m Context) (res Map, err error) {
	res = make(Map)
	for key, value := range m {
		res[key], err = parseNode(value)
		if err != nil {
			return
		}
	}
	return
}

func parseNode(value interface{}) (MapNode, error) {
	switch v := value.(type) {
	case string:
		t, err := pongo2.FromString(v)
		if err != nil {
			return nil, err
		}
		return &mapNodeTemplate{template: t}, nil
	case map[string]string:
		res := make(Map)
		for key, mvalue := range v {
			t, err := pongo2.FromString(mvalue)
			if err != nil {
				return nil, err
			}
			res[key] = &mapNodeTemplate{template: t}
		}
		return &mapNodeMap{nodes: res}, nil
	case map[string]interface{}:
		res, err := ParseMap(Context(v))
		if err != nil {
			return nil, err
		}
		return &mapNodeMap{nodes: res}, nil
	case Context:
		res, err := ParseMap(v)
		if err != nil {
			return nil, err
		}
		return &mapNodeMap{nodes: res}, nil
	case map[interface{}]interface{}:
		res := make(Map)
		for k, mvalue := range v {
			key, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("Invalid key %#v in templates map. All keys must be strings.", k)
			}
			n, err := parseNode(mvalue)
			if err != nil {
				return nil, err
			}
			res[key] = n
		}
		return &mapNodeMap{nodes: res}, nil
	case []interface{}:
		res := new(mapNodeSlice)
		for _, svalue := range v {
			n, err := parseNode(svalue)
			if err != nil {
				return nil, err
			}
			res.nodes = append(res.nodes, n)
		}
		return res, nil
	}
	return &mapNodeInterface{value: value}, nil
}

// Execute - Executes a map of templates and/or values.
func (nodes Map) Execute(ctx Context) (res Context, err error) {
	res = make(Context)
	for key, value := range nodes {
		res[key], err = value.Execute(ctx)
		if err != nil {
			return
		}
	}
	return
}

// ParseAndMerge - Parses a map of templates and merges into current map.
// If `t == nil` it may return a new map.
func (nodes Map) ParseAndMerge(input Context) (Map, error) {
	if len(input) == 0 {
		return nodes, nil
	}
	extra, err := ParseMap(input)
	if err != nil {
		return nil, err
	}
	for key, value := range extra {
		nodes[key] = value
	}
	return nodes, nil
}

type mapNodeMap struct {
	nodes Map
}

func (node *mapNodeMap) Execute(ctx Context) (interface{}, error) {
	res, err := node.nodes.Execute(ctx)
	return res, err
}

type mapNodeInterface struct {
	value interface{}
}

func (node *mapNodeInterface) Execute(_ Context) (interface{}, error) {
	return node.value, nil
}

type mapNodeSlice struct {
	nodes []MapNode
}

func (node *mapNodeSlice) Execute(ctx Context) (interface{}, error) {
	var res []interface{}
	for _, n := range node.nodes {
		v, err := n.Execute(ctx)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}

type mapNodeTemplate struct {
	template *pongo2.Template
}

func (node *mapNodeTemplate) Execute(ctx Context) (interface{}, error) {
	res, err := node.template.Execute(pongo2.Context(ctx))
	return res, err
}
