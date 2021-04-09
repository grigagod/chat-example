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

func (c *Client) GetRequestsList() []*sdb.Request {
	var requests []*sdb.Request
	c.dal.Db.Where("receiver_name = ?", c.username).Find(&requests)
	return requests
}

func (c *Client) GetFriendsList() []*sdb.Friend {
	var friends []*sdb.Friend
	c.dal.Db.Where("owner_name = ?", c.username).Find(&friends)
	return friends
}

func (dal *DAL) InsertIntoMessages(senderName string, receiverName string, message string, timestamp int64) error {
	msg := sdb.NewReceivedMessage(senderName, receiverName, message, timestamp)
	err := dal.Db.Create(&msg).Error
	return err
}

func (c *Client) GetMessagesList(friendname string) []*sdb.Message {
	var messages []*sdb.Message
	c.dal.Db.Where("sender_name IN ? AND receiver_name IN ?", []string{c.username, friendname}, []string{c.username, friendname}).Order("timestamp").Find(&messages)
	return messages
}

func (c *Client) InsertIntoFriends(friendName string, sharedKey *big.Int) error {
	friend := sdb.NewFriend(friendName, sharedKey.Bytes(), c.username)
	err := c.dal.Db.Create(&friend).Error
	return err
}

func (c *Client) InsertIntoRequests(inviter_name string, public_key []byte) error {
	request := sdb.NewRequest(inviter_name, public_key, c.username)
	err := c.dal.Db.Create(&request).Error
	return err
}

func (c *Client) DeleteFromRequests(inviter_name string) error {
	err := c.dal.Db.Where("sender_name = ? AND receiver_name = ?", inviter_name, c.username).Delete(sdb.Request{}).Error
	return err
}
