package jsonfn

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNode_GetChild(t *testing.T) {
	n := node{
		Children: []*node{
			{Name: "id"},
			{
				Name: "user",
				Children: []*node{
					{Name: "name"},
				},
			},
		},
	}

	c := n.getChild("user")
	if c == nil || c.Name != "user" {
		t.Fatalf("unexpected child: %v", c)
	}

	if len(c.Children) != 1 || c.Children[0].Name != "name" {
		t.Errorf("unexpected child: %v", c)
	}
}

func TestNode_Merge(t *testing.T) {
	n1 := node{
		Children: []*node{
			{Name: "id"},
			{
				Name: "user",
				Children: []*node{
					{Name: "name"},
				},
			},
		},
	}
	n2 := node{
		Children: []*node{
			{Name: "title"},
			{
				Name: "user",
				Children: []*node{
					{Name: "name"},
					{
						Name: "country",
						Children: []*node{
							{Name: "continent"},
						},
					},
				},
			},
		},
	}

	n1.Merge(&n2)

	if !compareNode(&node{
		Children: []*node{
			{Name: "id"},
			{
				Name: "user",
				Children: []*node{
					{Name: "name"},
					{
						Name: "country",
						Children: []*node{
							{Name: "continent"},
						},
					},
				},
			},
			{Name: "title"},
		},
	}, &n1) {
		j, _ := json.MarshalIndent(n1, "", "  ")
		t.Errorf("unexpected merge result: %s", j)
	}
}

func TestNode_AddChild(t *testing.T) {
	n := node{
		Children: []*node{
			{Name: "id"},
			{
				Name: "user",
				Children: []*node{
					{Name: "name"},
				},
			},
		},
	}

	c := node{
		Name: "title",
	}
	n.AddChild(&c)
	if child := n.getChild("title"); child == nil || child.Name != "title" {
		t.Errorf("unexpected added child: %v", child)
	}

	user := node{
		Name: "user",
		Children: []*node{
			{
				Name: "country",
				Children: []*node{
					{Name: "Continent"},
				},
			},
		},
	}
	n.AddChild(&user)
	u := n.getChild("user")
	if !compareNode(&node{
		Name: "user",
		Children: []*node{
			{Name: "name"},
			{
				Name: "country",
				Children: []*node{
					{Name: "Continent"},
				},
			},
		},
	}, u) {
		j, _ := json.MarshalIndent(u, "", "  ")
		t.Errorf("unexpected added child: %s", j)
	}
}

func TestExtractField(t *testing.T) {
	if err := testExtractField("*", nil, "*"); err != nil {
		t.Error(err)
	}
	if err := testExtractField("id", nil, "id"); err != nil {
		t.Error(err)
	}
	if err := testExtractField("book", []string{"id", "title"}, "book{id,title}"); err != nil {
		t.Error(err)
	}
	if err := testExtractField("book:user", []string{"id", "name"}, "book:user{id,name}"); err != nil {
		t.Error(err)
	}
}

func TestParseField(t *testing.T) {
	if err := testParseField(node{Name: "*"}, "*"); err != nil {
		t.Error(err)
	}
	if err := testParseField(node{Name: "id"}, "id"); err != nil {
		t.Error(err)
	}
	if err := testParseField(
		node{
			Name: "book",
			Children: []*node{
				{Name: "id"},
				{Name: "title"},
			},
		},
		"book{id,title}"); err != nil {
		t.Error(err)
	}
	if err := testParseField(
		node{
			Name: "book",
			Children: []*node{
				{
					Name: "user",
					Children: []*node{
						{Name: "id"},
						{Name: "name"},
					},
				},
			},
		}, "book:user{id,name}"); err != nil {
		t.Error(err)
	}
}

func TestParseFields(t *testing.T) {
	if err := testParseFields(node{
		Name: "",
		Children: []*node{
			{Name: "id"},
			{Name: "*"},
		},
	}, []string{"id", "*"}); err != nil {
		t.Error(err)
	}

	if err := testParseFields(
		node{
			Name: "",
			Children: []*node{
				{Name: "id"},
				{Name: "title"},
				{
					Name: "user",
					Children: []*node{
						{Name: "id"},
						{Name: "name"},
					},
				},
			},
		},
		[]string{"id", "title", "user{id,name}"}); err != nil {
		t.Error(err)
	}
	if err := testParseFields(
		node{
			Name: "",
			Children: []*node{
				{Name: "id"},
				{Name: "title"},
				{
					Name: "user",
					Children: []*node{
						{Name: "id"},
						{Name: "name"},
						{
							Name: "country",
							Children: []*node{
								{Name: "id"},
								{Name: "name"},
							},
						},
					},
				},
			},
		},
		[]string{"id", "title", "user{id,name}", "user:country{id,name}"}); err != nil {
		t.Error(err)
	}
}

func testExtractField(expectedNames string, expectedFields []string, field string) error {
	names, fields := extractField(field)
	if names != expectedNames {
		return fmt.Errorf("expect %s, but got %s", expectedNames, names)
	}

	l := len(fields)
	el := len(expectedFields)
	if l != el {
		return fmt.Errorf("expect %v, but got %v", expectedFields, fields)
	}

	for i := 0; i < l; i++ {
		if fields[i] != expectedFields[i] {
			fmt.Errorf("expect %v, but got %v", expectedFields, fields)
		}
	}

	return nil
}

func testParseField(expected node, field string) error {
	actual := parseField(field)
	if compareNode(&expected, actual) {
		return nil
	} else {
		eJson, _ := json.MarshalIndent(expected, "", "  ")
		aJson, _ := json.MarshalIndent(actual, "", "  ")
		return fmt.Errorf("expected %s, but got %s", eJson, aJson)
	}
}

func testParseFields(expected node, fields []string) error {
	actual := parseFields(fields)
	if compareNode(&expected, actual) {
		return nil
	} else {
		eJson, _ := json.MarshalIndent(expected, "", "  ")
		aJson, _ := json.MarshalIndent(actual, "", "  ")
		return fmt.Errorf("expected %s, but got %s", eJson, aJson)
	}
}

func compareNode(expected, actual *node) bool {
	if expected.Name != actual.Name {
		return false
	}

	l := len(actual.Children)
	el := len(expected.Children)
	if l != el {
		return false
	}
	for i := 0; i < l; i++ {
		if !compareNode(expected.Children[i], actual.Children[i]) {
			return false
		}
	}

	return true
}