/* This is the test file for both chat.go and broadcast.go
 */

package v1_test

import (
	"testing"

	"github.com/xybor/xychat/models"
	r "github.com/xybor/xychat/representations/v1"
	services "github.com/xybor/xychat/services/v1"
)

var role = "member"
var usn1 = "USN1"
var usn2 = "USN2"
var usn3 = "USN3"
var usn4 = "USN4"

var chatUser1 = models.User{
	Username: &usn1,
	Password: &usn1,
	Role:     &role,
}

var chatUser2 = models.User{
	Username: &usn2,
	Password: &usn2,
	Role:     &role,
}

var chatUser3 = models.User{
	Username: &usn3,
	Password: &usn3,
	Role:     &role,
}

var chatUser4 = models.User{
	Username: &usn4,
	Password: &usn4,
	Role:     &role,
}

func TestSetupDatabase(t *testing.T) {
	TestInitializeDB(t)
	models.GetDB().Create(&chatUser1)
	models.GetDB().Create(&chatUser2)
	models.GetDB().Create(&chatUser3)
	models.GetDB().Create(&chatUser4)

	us1 := services.CreateUserService(&chatUser1.ID, true)
	us2 := services.CreateUserService(&chatUser2.ID, true)
	us3 := services.CreateUserService(&chatUser3.ID, true)
	us4 := services.CreateUserService(&chatUser4.ID, true)

	rs := services.CreateRoomService(nil)

	xerr := rs.Create(&us1, &us2)
	if xerr.Errno() != 0 {
		t.Log(xerr)
	}

	xerr = rs.Create(&us3, &us4)
	if xerr.Errno() != 0 {
		t.Log(xerr)
	}

	xerr = rs.Create(&us1, &us3)
	if xerr.Errno() != 0 {
		t.Log(xerr)
	}
}

func TestChatJoinBroadcast(t *testing.T) {
	us1 := services.CreateUserService(&chatUser1.ID, true)
	us2 := services.CreateUserService(&chatUser2.ID, true)
	us3 := services.CreateUserService(&chatUser3.ID, true)
	us4 := services.CreateUserService(&chatUser4.ID, true)

	rs := services.CreateRoomService(nil)
	rs.Create(&us1, &us2)
	rs.Create(&us3, &us4)

	chatService1, err := services.CreateChatService(us1)
	if err.Errno() != 0 {
		t.Log(err)
		t.Fail()
	}

	chatService2, err := services.CreateChatService(us2)
	if err.Errno() != 0 {
		t.Log(err)
		t.Fail()
	}

	chatService3, err := services.CreateChatService(us3)
	if err.Errno() != 0 {
		t.Log(err)
		t.Fail()
	}

	chatService4, err := services.CreateChatService(us4)
	if err.Errno() != 0 {
		t.Log(err)
		t.Fail()
	}

	rid1 := uint(1)
	roomService1 := services.CreateRoomService(&rid1)

	rid2 := uint(2)
	roomService2 := services.CreateRoomService(&rid2)

	rid3 := uint(3)
	roomService3 := services.CreateRoomService(&rid3)

	broadcastService1 := services.GetBroadcastService(roomService1)
	if broadcastService1 != nil {
		t.Log("Invalid broadcast Service")
		t.FailNow()
	}

	broadcastService3 := services.GetBroadcastService(roomService3)
	if broadcastService3 != nil {
		t.Log("Invalid broadcast Service")
		t.FailNow()
	}

	broadcastService2 := services.GetBroadcastService(roomService2)
	if broadcastService2 != nil {
		t.Log("Invalid broadcast Service")
		t.FailNow()
	}

	chatService1.Online()

	broadcastService1 = services.GetBroadcastService(roomService1)
	if broadcastService1 == nil {
		t.Log("Cannot create the broadcast Service")
		t.FailNow()
	}

	broadcastService3 = services.GetBroadcastService(roomService3)
	if broadcastService3 == nil {
		t.Log("Cannot create the broadcast Service")
		t.FailNow()
	}

	broadcastService2 = services.GetBroadcastService(roomService2)
	if broadcastService2 != nil {
		t.Log("Invalid broadcast Service")
		t.FailNow()
	}

	chatService3.Online()

	broadcastService2 = services.GetBroadcastService(roomService2)
	if broadcastService2 == nil {
		t.Log("Cannot create the broadcast Service")
		t.FailNow()
	}

	broadcastService3 = services.GetBroadcastService(roomService3)
	if broadcastService3 == nil {
		t.Log("Cannot create the broadcast Service")
		t.FailNow()
	}

	chatService4.Online()

	broadcastService2 = services.GetBroadcastService(roomService2)
	if broadcastService2 == nil {
		t.Log("Cannot create the broadcast Service")
		t.FailNow()
	}

	chatService3.Offline()
	chatService4.Offline()

	broadcastService2 = services.GetBroadcastService(roomService2)
	if broadcastService2 != nil {
		t.Log("Cannot free the broadcast Service")
		t.FailNow()
	}

	broadcastService1 = services.GetBroadcastService(roomService1)
	if broadcastService1 == nil {
		t.Log("Cannot create the broadcast Service")
		t.FailNow()
	}

	chatService2.Online()

	broadcastService1 = services.GetBroadcastService(roomService1)
	if broadcastService1 == nil {
		t.Log("Cannot create the broadcast Service")
		t.FailNow()
	}

	chatService1.Offline()

	broadcastService3 = services.GetBroadcastService(roomService3)
	if broadcastService3 != nil {
		t.Log("Cannot free the broadcast Service")
		t.FailNow()
	}

	broadcastService1 = services.GetBroadcastService(roomService1)
	if broadcastService1 == nil {
		t.Log("Cannot create the broadcast Service")
		t.FailNow()
	}

	chatService2.Offline()

	broadcastService1 = services.GetBroadcastService(roomService1)
	if broadcastService1 != nil {
		t.Log("Cannot free the broadcast Service")
		t.FailNow()
	}
}

func TestChat(t *testing.T) {
	TestSetupDatabase(t)

	us1 := services.CreateUserService(&chatUser1.ID, true)
	us2 := services.CreateUserService(&chatUser2.ID, true)
	us3 := services.CreateUserService(&chatUser3.ID, true)
	us4 := services.CreateUserService(&chatUser4.ID, true)

	chatService1, xerr := services.CreateChatService(us1)
	if xerr.Errno() != 0 {
		t.Log(xerr)
		t.Fail()
	}

	chatService2, xerr := services.CreateChatService(us2)
	if xerr.Errno() != 0 {
		t.Log(xerr)
		t.Fail()
	}

	chatService3, xerr := services.CreateChatService(us3)
	if xerr.Errno() != 0 {
		t.Log(xerr)
		t.Fail()
	}

	chatService4, xerr := services.CreateChatService(us4)
	if xerr.Errno() != 0 {
		t.Log(xerr)
		t.Fail()
	}

	rid1 := uint(1)
	roomService1 := services.CreateRoomService(&rid1)
	rid2 := uint(2)
	roomService2 := services.CreateRoomService(&rid2)
	rid3 := uint(3)
	roomService3 := services.CreateRoomService(&rid3)

	broadcastService1 := services.GetBroadcastService(roomService1)
	if broadcastService1 != nil {
		t.Log("Invalid broadcast Service")
		t.FailNow()
	}

	broadcastService2 := services.GetBroadcastService(roomService2)
	if broadcastService2 != nil {
		t.Log("Invalid broadcast Service")
		t.FailNow()
	}

	broadcastService3 := services.GetBroadcastService(roomService3)
	if broadcastService3 != nil {
		t.Log("Invalid broadcast Service")
		t.FailNow()
	}

	chatService1.Online()
	chatService2.Online()
	chatService3.Online()
	chatService4.Online()

	chatService1.ChatHandler = func(cmr r.ChatMessageRepresentation) error {
		t.Logf("1 %d->%d: %s", cmr.UserId, cmr.RoomId, cmr.Message)
		return nil
	}

	chatService2.ChatHandler = func(cmr r.ChatMessageRepresentation) error {
		t.Logf("2 %d->%d: %s", cmr.UserId, cmr.RoomId, cmr.Message)
		return nil
	}

	chatService3.ChatHandler = func(cmr r.ChatMessageRepresentation) error {
		t.Logf("3 %d->%d: %s", cmr.UserId, cmr.RoomId, cmr.Message)
		return nil
	}

	chatService4.ChatHandler = func(cmr r.ChatMessageRepresentation) error {
		t.Logf("4 %d->%d: %s", cmr.UserId, cmr.RoomId, cmr.Message)
		return nil
	}

	broadcastService1 = services.GetBroadcastService(roomService1)
	if broadcastService1 == nil {
		t.Log("Cannot create the broadcast Service")
		t.FailNow()
	}

	broadcastService2 = services.GetBroadcastService(roomService2)
	if broadcastService2 == nil {
		t.Log("Cannot create the broadcast Service")
		t.FailNow()
	}

	broadcastService3 = services.GetBroadcastService(roomService3)
	if broadcastService3 == nil {
		t.Log("Cannot create the broadcast Service")
		t.FailNow()
	}

	chatService1.SendTo(1, "msg 1->1")
	chatService1.SendTo(3, "msg 1->3")
	chatService2.SendTo(1, "msg 2->1")
	chatService3.SendTo(2, "msg 3->2")
	chatService3.SendTo(3, "msg 3->3")
	chatService4.SendTo(2, "msg 4->2")

	chatService1.Offline()
	chatService2.Offline()
	chatService3.Offline()
	chatService4.Offline()

	broadcastService1 = services.GetBroadcastService(roomService1)
	if broadcastService1 != nil {
		t.Log("Cannot free the broadcast Service")
		t.FailNow()
	}

	broadcastService2 = services.GetBroadcastService(roomService2)
	if broadcastService2 != nil {
		t.Log("Cannot free the broadcast Service")
		t.FailNow()
	}

	broadcastService3 = services.GetBroadcastService(roomService3)
	if broadcastService3 != nil {
		t.Log("Cannot free the broadcast Service")
		t.FailNow()
	}
}
