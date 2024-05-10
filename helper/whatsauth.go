package helper

import (
	"gocroot/model"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsLoginRequest(msg model.IteungMessage, keyword string) bool {
	return strings.Contains(msg.Message, keyword) && msg.From_link
}

func GetUUID(msg model.IteungMessage, keyword string) string {
	return strings.Replace(msg.Message, keyword, "", 1)
}

func HandlerQRLogin(msg model.IteungMessage, WAKeyword string, WAPhoneNumber string, db *mongo.Database, WAAPIQRLogin string) (resp model.Response) {
	dt := &model.WhatsauthRequest{
		Uuid:        GetUUID(msg, WAKeyword),
		Phonenumber: msg.Phone_number,
		Delay:       msg.From_link_delay,
	}
	resp, _ = PostStructWithToken[model.Response]("Token", WAAPIToken(WAPhoneNumber, db), dt, WAAPIQRLogin)
	return
}

func HandlerIncomingMessage(msg model.IteungMessage, WAPhoneNumber string, db *mongo.Database, WAAPIMessage string) (resp model.Response) {
	dt := &model.TextMessage{
		To:       msg.Chat_number,
		IsGroup:  false,
		Messages: GetRandomReplyFromMongo(msg, db),
	}
	if msg.Chat_server == "g.us" { //jika pesan datang dari group maka balas ke group
		dt.IsGroup = true
	}
	if (msg.Phone_number != "628112000279") && (msg.Phone_number != "6283131895000") { //ignore pesan datang dari iteung
		resp, _ = PostStructWithToken[model.Response]("Token", WAAPIToken(WAPhoneNumber, db), dt, WAAPIMessage)
	}
	return
}

func GetRandomReplyFromMongo(msg model.IteungMessage, db *mongo.Database) string {
	rply, _ := GetRandomDoc[model.Reply](db, "reply", 1)
	replymsg := strings.ReplaceAll(rply[0].Message, "#BOTNAME#", msg.Alias_name)
	replymsg = strings.ReplaceAll(replymsg, "\\n", "\n")
	return replymsg
}

func WAAPIToken(phonenumber string, db *mongo.Database) string {
	filter := bson.M{"phonenumber": phonenumber}
	apitoken, _ := GetOneDoc[model.Profile](db, "profile", filter)
	return apitoken.Token
}
