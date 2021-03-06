package model

import (
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"time"
)

type Topic struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	UName      string        `bson:"uname"`
	Title      string        `bson:"title"`
	Body       template.HTML `bson:"body"`
	Tag        int           `bson:"tag"`
	CreateDate time.Time     `bson:"createdate"`
	Status     int           `bson:"status"`
	LastUpdate time.Time     `bson:"lastupdate"`
	Read       int           `bson:"read"`
	Comment    int           `bson:"comment"`
}

func (dao *Dao) InserTopic(topic *Topic) (string, error) {
	collection := dao.session.DB(DBNAME).C(T_TOPIC)
	//set the time
	topic.CreateDate = time.Now()
	topic.LastUpdate = time.Now()
	topic.Id = bson.NewObjectId()
	topic.Read = 0
	topic.Comment = 0
	_, err := collection.Upsert(bson.M{"_id": topic.Id}, topic)
	if err != nil {
		revel.WARN.Printf("Unable to save topic: %v error %v", topic, err)
	}
	return topic.Id.Hex(), err
}
func (dao *Dao) UpdateTopic(topic *Topic) (string, error) {
	collection := dao.session.DB(DBNAME).C(T_TOPIC)
	//set the time
	_, err := collection.Upsert(bson.M{"_id": topic.Id}, topic)
	if err != nil {
		revel.WARN.Printf("Unable to update topic: %v error %v", topic, err)
	}
	return topic.Id.Hex(), err
}
func (dao *Dao) EditTopic(topic *Topic) (string, error) {
	collection := dao.session.DB(DBNAME).C(T_TOPIC)
	//set the time
	topic.LastUpdate = time.Now()
	_, err := collection.Upsert(bson.M{"_id": topic.Id}, topic)
	if err != nil {
		revel.WARN.Printf("Unable to update topic: %v error %v", topic, err)
		return "", err
	}
	return topic.Id.Hex(), err
}
func (dao *Dao) FindTopicById(id string) *Topic {
	collection := dao.session.DB(DBNAME).C(T_TOPIC)
	topic := new(Topic)
	query := collection.Find(bson.M{"_id": bson.ObjectIdHex(id)})
	query.One(&topic)
	revel.INFO.Printf("The FindTopicById(%s) result: %v", id, topic)
	return topic
}
func (dao *Dao) GetTopics(page int) ([]Topic, int, int) {
	collection := dao.session.DB(DBNAME).C(T_TOPIC)
	topics := []Topic{}
	if page < 1 {
		page = 1
	}
	count, err := collection.Count()
	if err != nil {
		revel.WARN.Printf("Count topic:  error %v", err)
		count = 0
	}
	totalPage := count % PAGESIZE
	if totalPage == 0 {
		totalPage = count / PAGESIZE
	} else {
		totalPage = count/PAGESIZE + 1
	}
	if page > totalPage {
		page = totalPage
	}
	query := collection.Find(bson.M{"status": 1}).Sort("-lastupdate").Limit(PAGESIZE).Skip((page - 1) * PAGESIZE)
	query.All(&topics)
	return topics, page, totalPage
}
func (dao *Dao) GetTopicsWithTag(page, tag int) ([]Topic, int, int) {
	collection := dao.session.DB(DBNAME).C(T_TOPIC)
	topics := []Topic{}
	if page < 1 {
		page = 1
	}
	count, err := collection.Count()
	if err != nil {
		revel.WARN.Printf("Count topic:  error %v", err)
		count = 0
	}
	totalPage := count % PAGESIZE
	if totalPage == 0 {
		totalPage = count / PAGESIZE
	} else {
		totalPage = count/PAGESIZE + 1
	}
	if page > totalPage {
		page = totalPage
	}
	query := collection.Find(bson.M{"tag": tag, "status": 1}).Sort("-lastupdate").Limit(PAGESIZE).Skip((page - 1) * PAGESIZE)
	query.All(&topics)
	return topics, page, totalPage
}

func (dao *Dao) GetTopicsWithTagForMob(page, tag int) []Topic {
	topics := []Topic{}
	if page < 1 {
		return topics
	}
	collection := dao.session.DB(DBNAME).C(T_TOPIC)
	count, err := collection.Count()
	if err != nil {
		revel.WARN.Printf("Count topic:  error %v", err)
		count = 0
	}
	totalPage := count % PAGESIZE
	if totalPage == 0 {
		totalPage = count / PAGESIZE
	} else {
		totalPage = count/PAGESIZE + 1
	}
	if page > totalPage {
		return topics
	}
	query := collection.Find(bson.M{"tag": tag, "status": 1}).Sort("-lastupdate").Limit(PAGESIZE).Skip((page - 1) * PAGESIZE)
	query.All(&topics)
	return topics
}

func (dao *Dao) GetGoodTopics(tag, num int) []Topic {
	collection := dao.session.DB(DBNAME).C(T_TOPIC)
	topics := []Topic{}
	if tag == 10 {
		query := collection.Find(bson.M{"tag": tag, "status": 1}).Sort("-lastupdate").Limit(num)
		query.All(&topics)
	} else {
		query := collection.Find(bson.M{"status": 1}).Sort("-comment", "-read", "-lastupdate").Limit(num)
		query.All(&topics)
	}
	return topics
}
func (dao *Dao) GetTopicsByUserName(name string, limit int) []Topic {
	collection := dao.session.DB(DBNAME).C(T_TOPIC)
	topics := []Topic{}
	query := collection.Find(bson.M{"status": 1, "uname": name}).Sort("-lastupdate").Limit(limit)
	query.All(&topics)
	return topics
}
