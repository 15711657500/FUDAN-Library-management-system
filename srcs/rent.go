package main

func resetrent(lib *Library) error {
	_, err := lib.db.Exec(`
	drop table if exists rent;
`)
	return err
}

func createrent(lib *Library) error {
	_, err := lib.db.Exec(`
	create table if not exists rent
(
    rentdate date,
    duedate date,
    returndate date,
    fine float,
    rentid int primary key auto_increment,
    username varchar(50) references user(username),
    bookid int references book(bookid),
    extend int default 0
);
`)
	return err
}
