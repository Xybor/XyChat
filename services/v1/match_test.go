package v1_test

import (
	"log"
	"testing"
	"time"

	"github.com/xybor/xychat/models"
	r "github.com/xybor/xychat/representations/v1"
	services "github.com/xybor/xychat/services/v1"
)

var matchUser1 models.User
var matchUser2 models.User

func TestInitializeMatchQueue(t *testing.T) {
	services.InitializeMatchQueue(10 * time.Second)
	TestInitializeDB(t)

	matchUser1, err = CreateUser("USN1", "PWD1", "member")
	if err != nil {
		log.Panicln(err)
	}

	matchUser2, err = CreateUser("USN2", "PWD2", "member")
	if err != nil {
		log.Panicln(err)
	}
}

func TestNoneMatchService(t *testing.T) {
	userService := services.CreateUserService(nil)
	_, err := services.CreateMatchService(userService)
	if err != services.ErrorPermission {
		t.Log("Create none match service")
		t.Fail()
	}
}

func TestMatchServiceJoinQueue(t *testing.T) {
	userService := services.CreateUserService(&matchUser1.ID)
	matchService, err := services.CreateMatchService(userService)
	if err != nil {
		t.Log("Create match service: ", err)
		log.Fatal()
	}

	defer matchService.Close()

	matchService.Register()
	time.Sleep(1 * time.Second)
	_, ok := services.GetMatchQueue().GetQueue()[matchService]

	if !ok {
		t.Log("Register failure")
		t.Fail()
	}

	matchService.Unregister()
	time.Sleep(1 * time.Second)

	_, ok = services.GetMatchQueue().GetQueue()[matchService]

	if ok {
		t.Log("Unregister failure")
		t.Fail()
	}

	if services.GetMatchQueue().GetQueueLen() != 0 {
		t.Log("Non-empty queue")
		t.Fail()
	}
}

func TestDuplicatedMatchService(t *testing.T) {
	userService1 := services.CreateUserService(&matchUser1.ID)
	matchService1, err := services.CreateMatchService(userService1)
	if err != nil {
		t.Log("Create match service: ", err)
		t.Fail()
	} else {
		defer matchService1.Close()
	}

	userService2 := services.CreateUserService(&matchUser1.ID)
	matchService2, err := services.CreateMatchService(userService2)
	if err != services.ErrorDuplicatedConnection {
		t.Log("Create match service: ", err)
		t.Fail()
	}

	if err == nil {
		defer matchService2.Close()
	}
}

func TestTwoMatchServiceJoin(t *testing.T) {
	if services.GetMatchQueue().GetQueueLen() != 0 {
		t.Log("Non-empty queue")
		t.Fail()
	}

	userService1 := services.CreateUserService(&matchUser1.ID)
	matchService1, err := services.CreateMatchService(userService1)
	if err != nil {
		t.Log("Create match service: ", err)
		t.FailNow()
	}
	defer matchService1.Close()

	userService2 := services.CreateUserService(&matchUser2.ID)
	matchService2, err := services.CreateMatchService(userService2)
	if err != nil {
		t.Log("Create match service: ", err)
		t.FailNow()
	}
	defer matchService2.Close()

	matchService1.Register()
	matchService2.Register()

	var r1, r2 r.RoomRepresentation
	matchService1.MatchHandler = func(rr r.RoomRepresentation) { r1 = rr }
	matchService2.MatchHandler = func(rr r.RoomRepresentation) { r2 = rr }

	time.Sleep(11 * time.Second)

	if r1.ID != r2.ID {
		t.Log("Different roomid")
		t.Fail()
	}

	if _, ok := services.GetMatchQueue().GetQueue()[matchService1]; ok {
		t.Log("matchService1 has still existed in matchQueue")
		t.Fail()
	}

	if _, ok := services.GetMatchQueue().GetQueue()[matchService2]; ok {
		t.Log("matchService2 has still existed in matchQueue")
		t.Fail()
	}

	if services.GetMatchQueue().GetQueueLen() != 0 {
		t.Log("Non-empty queue")
		t.Fail()
	}
}
