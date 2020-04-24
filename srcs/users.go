package main

func resetusers(lib *Library) error {
	_, err := lib.db.Exec(`
drop table if exists user;
`)
	return err
}
func createusers(lib *Library) error {
	_, err := lib.db.Exec(`
	create table if not exists user
(
    username varchar(50) primary key,
    password varchar(50),
    permission varchar(10) default "default"
);
`)
	return err
}
