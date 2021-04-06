package main

import (
	"log"
	"math/big"
	"os"

	//"fmt"
	"github.com/grigagod/chat-example/sdb"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

type DAL struct {
	DSN string
	Db  *gorm.DB
}

func createDAL(dns string, logger *log.Logger) *DAL {
	var dal = new(DAL)
	dal.DSN = dns
	db, err := sdb.CreateConnection(dal.DSN, logger)
	if err != nil {
		os.Create(dal.DSN)
		db, err := sdb.CreateConnection(dal.DSN, logger)
		if err != nil {
			panic(err)
		}
		dal.Db = db
		return dal
	}
	dal.Db = db
	return dal
}

func (dal *DAL) InsertIntoUsers(username string, privateKey *big.Int) error {
	user := sdb.NewUser(username, privateKey.Bytes())
	err := dal.Db.Create(&user).Error
	return err
}

func (dal *DAL) GetUser(username string) (*sdb.User, error) {
	var user sdb.User
	if err := dal.Db.First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (dal *DAL) GetRequestsList(username string) []*sdb.Request {
	var requests []*sdb.Request
	dal.Db.Where("receiver_name = ?", username).Find(&requests)
	return requests
}

func (dal *DAL) GetFriendsList(username string) []*sdb.Friend {
	var friends []*sdb.Friend
	dal.Db.Where("owner_name = ?", username).Find(&friends)
	return friends
}

func (dal *DAL) InsertIntoMessages(senderName string, receiverName string, message string, timestamp int64) error {
	msg := sdb.NewReceivedMessage(senderName, receiverName, message, timestamp)
	err := dal.Db.Create(&msg).Error
	return err
}

func (dal *DAL) GetMessagesList(userName string, friendName string) []*sdb.Message {
	var messages []*sdb.Message
	dal.Db.Where("sender_name IN ? AND receiver_name IN ?", userName, friendName).Find(&messages)
	return messages
}

func (dal *DAL) InsertIntoFriends(friendName string, sharedKey *big.Int, username string) error {
	friend := sdb.NewFriend(friendName, sharedKey.Bytes(), username)
	err := dal.Db.Create(&friend).Error
	return err
}

func (dal *DAL) InsertIntoRequests(inviter_name string, public_key []byte, username string) error {
	request := sdb.NewRequest(inviter_name, public_key, username)
	err := dal.Db.Create(&request).Error
	return err
}

func (dal *DAL) DeleteFromRequests(inviter_name string, username string) error {
	err := dal.Db.Where("sender_name = ? AND receiver_name = ?", inviter_name, username).Delete(sdb.Request{}).Error
	return err
}
