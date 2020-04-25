package main

type Book struct {
	title  string
	author string
	ISBN   string
}

func resetbook(lib *Library) error {
	_, err := lib.db.Exec(`
	drop table if exists book;
`)
	return err
}
func createbook(lib *Library) error {
	_, err := lib.db.Exec(`
	create table if not exists book
(
    title     varchar(50),
    author    varchar(50),
    ISBN      varchar(100),
    bookid int primary key auto_increment,
    available bool default true
);
`)
	return err
}
func addbook(lib *Library) error {
	return nil
}
