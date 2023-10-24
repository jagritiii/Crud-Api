package service

import (
	Mongo "awesomeProject/pkg/dataaccess/mongo"
	"awesomeProject/pkg/dataaccess/redis"
	"awesomeProject/pkg/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

func CreateUser(user []models.User) error {
	manager := Mongo.MongoManager()

	idcount, err := manager.Totalcount(context.Background())
	if err != nil {
		log.Error(err.Error())
		return err
	}
	idcount++
	for _, val := range user {
		val.Id = int(idcount)
		_, err = manager.Insert(context.Background(), val)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		idcount++
	}

	return nil
}
func Getalluser() ([]models.User, error) {
	var res []models.User
	filter := bson.M{}
	manager := Mongo.MongoManager()
	cur, err := manager.Findusers(context.Background(), filter)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	for cur.Next(context.Background()) {
		var temp models.User
		err := cur.Decode(&temp)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		res = append(res, temp)
	}
	return res, nil
}
func UpdateUser(user []models.User) error {
	manager := Mongo.MongoManager()
	redimanager := redis.Redismanager()

	for _, val := range user {

		if val.Id == 0 {
			idcount, err := manager.Totalcount(context.Background())
			fmt.Println(idcount)
			if err != nil {
				log.Error(err.Error())
				return err
			}
			idcount++
			val.Id = int(idcount)
		}
		str := strconv.Itoa(val.Id)
		filter := bson.M{"Id": val.Id}
		update := bson.M{"$set": bson.M{"Company": val.Company, "Profile": val.Profile, "Age": val.Age, "Exp": val.Experience}}
		opts := options.Update().SetUpsert(true)
		_, err := manager.Updateone(context.Background(), filter, update, opts)
		asli, err := json.Marshal(&val)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		err = redimanager.Setredis(context.Background(), str, asli, time.Second*10*500).Err()
		if err != nil {
			log.Error(err.Error())
			return err
		}
		//if k.Err() != nil {
		//	log.Error(err.Error())
		//	return err
		//}
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}
func Getauser(id int) (models.User, error, error) {
	var user models.User
	manager := Mongo.MongoManager()
	redismanager := redis.Redismanager()
	str := strconv.Itoa(id)
	ans, err := redismanager.Getredis(context.Background(), str).Result()
	if err == nil {
		fmt.Println("yei value redis sei aaya hai")
		err = json.Unmarshal([]byte(ans), &user)
		if err != nil {
			log.Error(err.Error())
			return user, nil, err
		}
		return user, nil, nil
	} else {
		fmt.Println("yei value humlong ko mila hai mongosei")
		filter := bson.M{"Id": id}

		err := manager.Findone(context.Background(), filter).Decode(&user)
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Error(err)
			return user, errors.New("hai hi nahi dhundhu kaha sei"), nil
		}
		asli, err := json.Marshal(&user)
		err = redismanager.Setredis(context.Background(), str, asli, time.Second*10*500).Err()
		if err != nil {
			log.Error(err.Error())
		}
		return user, nil, nil

	}
}
func Deleteuser(id int) (error, error) {
	manager := Mongo.MongoManager()
	redismanager := redis.Redismanager()
	str := strconv.Itoa(id)
	err := redismanager.Getredis(context.Background(), str).Err()
	if err != nil {
		log.Error(err.Error())
	}
	if err == nil {
		err = redismanager.Deletekey(context.Background(), str).Err()
		if err != nil {
			log.Error(err.Error())
		}
		return err, nil
	}
	filter := bson.M{"Id": id}
	val := manager.Findanddelete(context.Background(), filter)

	if val == nil {
		return err, err
	}
	return nil, nil
}
