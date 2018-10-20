package jsonfn

import (
	"strconv"
	"testing"
	"time"
)

type Book struct {
	Id        int
	Title     string
	AuthorId  int
	CreatedAt time.Time
}

func (b Book) Author() Author {
	return Author{
		Id:   b.AuthorId,
		Name: "author" + strconv.Itoa(b.AuthorId),
	}
}

func (b Book) Tags() []Tag {
	return []Tag{
		{Id: 1, Title: "tag1"},
		{Id: 2, Title: "tag2"},
	}
}

type Author struct {
	Id        int
	Name      string
	CountryId int
}

type Tag struct {
	Id    int
	Title string
}

func (t Tag) Categories() []Category {
	return []Category{
		{Id: 1, Title: "Asia"},
		{Id: 2, Title: "Europe"},
	}
}

func (a Author) Country() Country {
	return Country{
		Id:   a.CountryId,
		Name: "country" + strconv.Itoa(a.CountryId),
	}
}

type Country struct {
	Id   int
	Name string
}

type Category struct {
	Id    int
	Title string
}

func TestMarshalSelectedFields(t *testing.T) {
	book := Book{
		Id:        1,
		Title:     "Jane Eyre",
		AuthorId:  2,
		CreatedAt: time.Now(),
	}
	jsonStr, _ := Marshal(book, "Id", "Title")
	if string(jsonStr) != `{"Id":1,"Title":"Jane Eyre"}` {
		t.Errorf("unexpected %s", jsonStr)
	}
}

func TestMarshalAllFields(t *testing.T) {
	book := Author{
		Id:        1,
		Name:      "Liam",
		CountryId: 2,
	}

	jsonStr, _ := Marshal(book)
	if string(jsonStr) != `{"CountryId":2,"Id":1,"Name":"Liam"}` {
		t.Errorf("unexpected %s", jsonStr)
	}

	jsonStr, _ = Marshal(book, "*")
	if string(jsonStr) != `{"CountryId":2,"Id":1,"Name":"Liam"}` {
		t.Errorf("unexpected %s", jsonStr)
	}
}

func TestMarshalEmbedded(t *testing.T) {
	book := Book{
		Id:        1,
		Title:     "Jane Eyre",
		AuthorId:  2,
		CreatedAt: time.Now(),
	}
	jsonStr, _ := Marshal(book, "Id", "Title", "Author{Id,Name}")
	if string(jsonStr) != `{"Author":{"Id":2,"Name":"author2"},"Id":1,"Title":"Jane Eyre"}` {
		t.Errorf("unexpected %s", jsonStr)
	}
}

func TestMarshalEmbeddedAllFields(t *testing.T) {
	book := Book{
		Id:        1,
		Title:     "Jane Eyre",
		AuthorId:  2,
		CreatedAt: time.Now(),
	}
	jsonStr, _ := Marshal(book, "Id", "Title", "Author{}")
	if string(jsonStr) != `{"Author":{"CountryId":0,"Id":2,"Name":"author2"},"Id":1,"Title":"Jane Eyre"}` {
		t.Errorf("unexpected %s", jsonStr)
	}

	jsonStr, _ = Marshal(book, "Id", "Title", "Author{*}")
	if string(jsonStr) != `{"Author":{"CountryId":0,"Id":2,"Name":"author2"},"Id":1,"Title":"Jane Eyre"}` {
		t.Errorf("unexpected %s", jsonStr)
	}
}

func TestMarshalEmbeddedSlice(t *testing.T) {
	book := Book{
		Id:        1,
		Title:     "Jane Eyre",
		AuthorId:  2,
		CreatedAt: time.Now(),
	}
	jsonStr, _ := Marshal(book, "Id", "Title", "Tags{*}")
	if string(jsonStr) != `{"Id":1,"Tags":[{"Id":1,"Title":"tag1"},{"Id":2,"Title":"tag2"}],"Title":"Jane Eyre"}` {
		t.Errorf("unexpected %s", jsonStr)
	}
}

func TestMarshalMultiLayerEmbedded(t *testing.T) {
	book := Book{
		Id:        1,
		Title:     "Jane Eyre",
		AuthorId:  2,
		CreatedAt: time.Now(),
	}
	jsonStr, _ := Marshal(book, "Id", "Title", "Author{Id,Name}", "Author:Country{Id,Name}")
	if string(jsonStr) != `{"Author":{"Country":{"Id":0,"Name":"country0"},"Id":2,"Name":"author2"},"Id":1,"Title":"Jane Eyre"}` {
		t.Errorf("unexpected %s", jsonStr)
	}
}

func TestMarshalMultiLayerEmbeddedSlice(t *testing.T) {
	book := Book{
		Id:        1,
		Title:     "Jane Eyre",
		AuthorId:  2,
		CreatedAt: time.Now(),
	}
	jsonStr, _ := Marshal(book, "Id", "Title", "Tags{Title}", "Tags:Categories{Title}")
	if string(jsonStr) != `{"Id":1,"Tags":[{"Categories":[{"Title":"Asia"},{"Title":"Europe"}],"Title":"tag1"},{"Categories":[{"Title":"Asia"},{"Title":"Europe"}],"Title":"tag2"}],"Title":"Jane Eyre"}` {
		t.Errorf("unexpected %s", jsonStr)
	}
}

func TestWontPanic(t *testing.T) {
	book := Book{
		Id:        1,
		Title:     "Jane Eyre",
		AuthorId:  2,
		CreatedAt: time.Now(),
	}

	defer func() {
		if err := recover(); err != nil {
			t.Error("marshal panic")
		}
	}()
	Marshal(book,
		"Id",
		"foo",
		"*",
		"Author{Id,bar}",
		"Author{*}",
		"Author:Country{Id,foo}",
		"Author:Foo{}",
		"Foo:Bar{*}",
	)
}
