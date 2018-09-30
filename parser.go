package jsonfn

import (
	"strings"
	"regexp"
)

var fieldExp = regexp.MustCompile(`(\*)|([\w]+[:\w]*)({([\w,*]*)})?`)

type node struct {
	Name     string
	Children []*node
}

func (n node) GetFields() []string {
	var fields []string
	for _, child := range n.Children {
		if len(child.Children) == 0 {
			fields = append(fields, child.Name)
		}
	}

	return fields
}

func (n node) IsLeaf() bool {
	return len(n.Children) == 0
}

func (n *node) AddChild(c *node) {
	if child := n.getChild(c.Name); child != nil {
		child.Merge(c)
	} else {
		n.Children = append(n.Children, c)
	}
}

func (n *node) Merge(other *node) {
	for _, oc := range other.Children {
		if nc := n.getChild(oc.Name); nc != nil {
			nc.Merge(oc)
		} else {
			n.Children = append(n.Children, oc)
		}
	}
}

func (n node) getChild(name string) *node {
	for _, c := range n.Children {
		if c.Name == name {
			return c
		}
	}

	return nil
}

func parseFields(fields []string) *node {
	nd := node{}
	for _, field := range fields {
		child := parseField(field)
		nd.AddChild(child)
	}

	return &nd
}

func parseField(field string) *node {
	name, fields := extractField(field)

	if name == "" {
		return nil
	} else if fields == nil {
		return &node{
			Name: name,
		}
	}

	names := strings.Split(name, ":")
	namesLen := len(names)
	leaf := node{
		Name: names[namesLen-1],
	}
	for _, f := range fields {
		leaf.AddChild(&node{
			Name: f,
		})
	}

	nd := &leaf
	for i := namesLen - 2; i >= 0; i-- {
		current := node{
			Name: names[i],
		}
		current.AddChild(nd)
		nd = &current
	}

	return nd
}

func extractField(field string) (names string, fields []string) {
	matches := fieldExp.FindStringSubmatch(field)

	if matches == nil {
		return
	}

	if matches[1] == "*" || matches[3] == "" {
		return field, nil
	}

	names = matches[2]
	parts := strings.Split(matches[4], ",")
	for _, p := range parts {
		if p != "" {
			fields = append(fields, p)
		}
	}

	return
}
