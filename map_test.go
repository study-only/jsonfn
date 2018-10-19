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

func TestMarshal(t *testing.T) {
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
