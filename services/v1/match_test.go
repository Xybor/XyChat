package v1_test

import (
	"log"
	"testing"
	"time"

	services "github.com/xybor/xychat/services/v1"
)

func TestInitializeMatchQueue(t *testing.T) {
	services.InitializeMatchQueue(10 * time.Second)
	TestInitializeDB(t)
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
	id := uint(1)

	userService := services.CreateUserService(&id)
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
	id := uint(1)

	userService1 := services.CreateUserService(&id)
	matchService1, err := services.CreateMatchService(userService1)
	if err != nil {
		t.Log("Create match service: ", err)
		t.Fail()
	} else {
		defer matchService1.Close()
	}

	userService2 := services.CreateUserService(&id)
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

	id1 := uint(1)
	id2 := uint(2)

	userService1 := services.CreateUserService(&id1)
	matchService1, err := services.CreateMatchService(userService1)
	if err != nil {
		t.Log("Create match service: ", err)
		t.FailNow()
	}
	defer matchService1.Close()

	userService2 := services.CreateUserService(&id2)
	matchService2, err := services.CreateMatchService(userService2)
	if err != nil {
		t.Log("Create match service: ", err)
		t.FailNow()
	}
	defer matchService2.Close()

	matchService1.Register()
	matchService2.Register()

	r1 := matchService1.WaitForRoom()
	r2 := matchService2.WaitForRoom()

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
