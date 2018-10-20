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

type Author struct {
	Id        int
	Name      string
	CountryId int
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

func TestMarshalSelectedFields(t *testing.T) {
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

func TestMarshalAllFields(t *testing.T) {
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

func TestMarshalEmbedded(t *testing.T) {
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
