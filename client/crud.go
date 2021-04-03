package main 

import (
	"database/sql"
	"math/big"
	//"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type DAL struct {
	conStr string
	Database *sql.DB
}

func (self *DAL) OpenConnection(database string) {
	db, err := sql.Open("sqlite3", database)	
	if err != nil {
		panic(err)	
	}
	self.Database = db
}

func (self *DAL) CloseConnection() {
	self.Database.Close()
}

func (self *DAL) InsertIntoFriends(friend_name string, shared_key *big.Int) {
	self.OpenConnection("chat.db")
	_, err := self.Database.Exec("insert into friends (friend_name, shared_key) values ($1, $2)",
		friend_name, shared_key.Bytes())
	if err != nil {
		panic(err)	
	}
	defer self.CloseConnection()
}

func (self *DAL) InsertIntoInvites(inviter_name string, public_key *big.Int) {
	self.OpenConnection("chat.db")
	_, err := self.Database.Exec("insert into invites (inviter_name, public_key) values ($1, $2)",
		inviter_name, public_key.Bytes())
	if err != nil {
		panic(err)	
	}
	defer self.CloseConnection()
}

func (self *DAL) SelectAllFrom(table string) *sql.Rows {
	self.OpenConnection("chat.db")
	result, err := self.Database.Query("select * from " + table)
	if err != nil {
		panic(err)	
	}
	defer self.CloseConnection()
	return result
}

func main() {
	var dal DAL;
	//x := new(big.Int).SetInt64(12345)
	//dal.InsertIntoFriends("24121", x)
	dal.SelectAllFrom("friends")
}

