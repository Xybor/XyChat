package v1_test

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/xybor/xychat/helpers"
	"github.com/xybor/xychat/models"
	services "github.com/xybor/xychat/services/v1"
)

var adminUsn = "admin"
var adminPwd = "admin"
var adminRole = "admin"
var admin models.User

var modUsn = "mod"
var modPwd = "mod"
var modRole = "mod"
var mod models.User

var memUsn = "member"
var memPwd = "member"
var memRole = "member"
var member models.User

var err error

func CreateUser(username, password, role string) (models.User, error) {
	user := models.User{Username: &username, Password: &password, Role: &role}
	err = models.GetDB().Create(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func Register(id *uint, name, role string) error {
	userService := services.CreateUserService(id)
	err = userService.Register(name, "password", role)
	if err != nil {
		return err
	}

	return nil
}

func Remove(id *uint, name, role string) error {
	var user models.User
	if user, err = CreateUser(name, "password", role); err != nil {
		return err
	}

	userService := services.CreateUserService(id)

	err1 := userService.RemoveByUsername(name)
	if err1 == nil {
		if user, err = CreateUser(name, "password", role); err != nil {
			return err
		}
	}

	err2 := userService.Remove(user.ID)

	if err1 != nil {
		return err1
	}

	if err2 != nil {
		return err2
	}

	return nil
}

func SelfRemove(role string) error {
	var user models.User
	if user, err = CreateUser(role+"removeself1", "password", role); err != nil {
		return err
	}

	userService := services.CreateUserService(&user.ID)

	err := userService.RemoveByUsername(*user.Username)
	if err == nil {
		return err
	}

	if user, err = CreateUser(role+"removeself2", "password", role); err != nil {
		return err
	}

	err = userService.Remove(user.ID)
	if err != nil {
		return err
	}

	return nil
}

func Select(id *uint, role string) error {
	userService := services.CreateUserService(id)
	var user models.User

	usn := "select" + role
	if id == nil {
		usn = "none" + usn
	} else {
		usn = strconv.FormatUint(uint64(*id), 10) + usn
	}

	if user, err = CreateUser(usn, "password", role); err != nil {
		return err
	}

	r, err := userService.Select(user.ID)
	if err != nil {
		return err
	}

	if r.ID != user.ID {
		return errors.New("Different result id")
	}

	r, err = userService.SelectByName(usn)
	if err != nil {
		return err
	}

	if r.ID != user.ID {
		return errors.New("Different result id")
	}

	return nil
}

func SelfSelect(id *uint) error {
	userService := services.CreateUserService(id)

	r, err := userService.SelfSelect()
	if err != nil {
		return err
	}

	if id == nil || r.ID != *id {
		return errors.New("Different result id")
	}

	return nil
}

func UpdateInfo(id *uint, username, role string, age *uint, gender *string) error {
	var user models.User

	if user, err = CreateUser(username, "password", role); err != nil {
		return err
	}

	userService := services.CreateUserService(id)

	if err = userService.UpdateInfo(user.ID, age, gender); err != nil {
		return err
	}

	return nil
}

func SelfUpdateInfo(id *uint) error {
	age := uint(10)
	gender := "female"

	userService := services.CreateUserService(id)

	if err = userService.UpdateInfo(*id, &age, &gender); err != nil {
		return err
	}

	return nil
}

func UpdateRole(id *uint, username, role, newrole string) error {
	var user models.User

	if user, err = CreateUser(username, "password", role); err != nil {
		return err
	}

	userService := services.CreateUserService(id)

	if err = userService.UpdateRole(user.ID, newrole); err != nil {
		return err
	}

	return nil
}

func SelfUpdateRole(id *uint) error {
	userService := services.CreateUserService(id)

	if err = userService.UpdateRole(*id, "member"); err != services.ErrorPermission {
		if err == nil {
			return errors.New("invalid self update to member")
		}
		return err
	}

	if err = userService.UpdateRole(*id, "mod"); err != services.ErrorPermission {
		if err == nil {
			return errors.New("invalid self update to mod")
		}
		return err
	}

	if err = userService.UpdateRole(*id, "admin"); err != services.ErrorPermission {
		if err == nil {
			return errors.New("invalid self update to admin")
		}
		return err
	}

	return nil
}

func UpdatePassword(id *uint, username, role string) error {
	var user models.User
	password := "password"

	if user, err = CreateUser(username, password, role); err != nil {
		return err
	}

	userService := services.CreateUserService(id)

	var err1, err2 error

	err1 = userService.UpdatePassword(user.ID, &password, "newpassword")
	err2 = userService.UpdatePassword(user.ID, nil, "newpassword")

	if err1 != err2 {
		return errors.New("Error from with and without password is different; " +
			err1.Error() + ";" + err2.Error(),
		)
	}

	return err1
}

func SelfUpdatePassword(username, role string) error {
	var user models.User
	password := "password"

	if user, err = CreateUser(username, password, role); err != nil {
		return err
	}

	userService := services.CreateUserService(&user.ID)

	if err = userService.UpdatePassword(user.ID, nil, "newpassword"); err != services.ErrorInvalidOldPassword {
		if err == nil {
			return errors.New("self change password without old password")
		}
		return err
	}

	if err = userService.UpdatePassword(user.ID, &password, "newpassword"); err != nil {
		return err
	}

	return nil
}

func TestInitializeDB(t *testing.T) {
	dbname := helpers.MustReadEnv("POSTGRES_DBNAME")
	if !strings.Contains(dbname, "Test") {
		log.Fatal("Please use a test database for testing")
	}

	models.InitializeDB()
	models.CreateTables(true)
}

func TestCreateSeedingUsers(t *testing.T) {
	if admin, err = CreateUser(adminUsn, adminPwd, adminRole); err != nil {
		log.Fatal("Register admin user: ", err)
	}

	if mod, err = CreateUser(modUsn, modPwd, modRole); err != nil {
		log.Fatal("Register mod user: ", err)
	}

	if member, err = CreateUser(memUsn, memPwd, memRole); err != nil {
		log.Fatal("Register member user: ", err)
	}
}

func TestRegisterMember(t *testing.T) {
	if err := Register(nil, "memberbynone", "member"); err != nil {
		t.Log("Register member by none: ", err)
		t.Fail()
	}

	if err := Register(&member.ID, "memberbymember", "member"); err != nil {
		t.Log("Register member by member: ", err)
		t.Fail()
	}

	if err := Register(&mod.ID, "memberbymod", "member"); err != nil {
		t.Log("Register member by mod: ", err)
		t.Fail()
	}

	if err := Register(&admin.ID, "memberbyadmin", "member"); err != nil {
		t.Log("Register member by admin: ", err)
		t.Fail()
	}
}

func TestRegisterMod(t *testing.T) {
	if err := Register(nil, "modbynone", "mod"); err != services.ErrorPermission {
		t.Log("Register mod by none: ", err)
		t.Fail()
	}

	if err := Register(&member.ID, "modbymember", "mod"); err != services.ErrorPermission {
		t.Log("Register mod by member: ", err)
		t.Fail()
	}

	if err := Register(&mod.ID, "modbymod", "mod"); err != services.ErrorPermission {
		t.Log("Register mod by mod: ", err)
		t.Fail()
	}

	if err := Register(&admin.ID, "modbyadmin", "mod"); err != nil {
		t.Log("Register mod by admin: ", err)
		t.Fail()
	}
}

func TestRegisterAdmin(t *testing.T) {
	if err := Register(nil, "adminbynone", "admin"); err != services.ErrorPermission {
		t.Log("Register admin by none: ", err)
		t.Fail()
	}

	if err := Register(&member.ID, "adminbymember", "admin"); err != services.ErrorPermission {
		t.Log("Register admin by member: ", err)
		t.Fail()
	}

	if err := Register(&mod.ID, "adminbymod", "admin"); err != services.ErrorPermission {
		t.Log("Register admin by mod: ", err)
		t.Fail()
	}

	if err := Register(&admin.ID, "adminbyadmin", "admin"); err != services.ErrorPermission {
		t.Log("Register admin by admin: ", err)
		t.Fail()
	}
}

func TestRegisterUnknownRole(t *testing.T) {
	if err := Register(nil, "userunknownrole", "something"); err != services.ErrorUnknownRole {
		t.Log("Register unknown role: ", err)
	}
}

func TestRegisterDuplication(t *testing.T) {
	if err := Register(nil, "memberbynonedup", "member"); err != nil {
		t.Log("Register member by none dup 1: ", err)
		t.Fail()
	}

	if err := Register(nil, "memberbynonedup", "member"); err != services.ErrorExistedUsername {
		t.Log("Register member by none dup 2: ", err)
		t.Fail()
	}
}

func TestRemoveMember(t *testing.T) {
	if err := Remove(nil, "memberremovebynone", "member"); err != services.ErrorPermission {
		t.Log("Remove member by none: ", err)
		t.Fail()
	}

	if err := Remove(&member.ID, "memberremovebymember", "member"); err != services.ErrorPermission {
		t.Log("Remove member by member: ", err)
		t.Fail()
	}

	if err := Remove(&mod.ID, "memberremovebymod", "member"); err != nil {
		t.Log("Remove member by mod: ", err)
		t.Fail()
	}

	if err := Remove(&admin.ID, "memberremovebyadmin", "member"); err != nil {
		t.Log("Remove member by admin: ", err)
		t.Fail()
	}
}
func TestRemoveMod(t *testing.T) {
	if err := Remove(nil, "modremovebynone", "mod"); err != services.ErrorPermission {
		t.Log("Remove mod by none: ", err)
		t.Fail()
	}

	if err := Remove(&member.ID, "modremovebymember", "mod"); err != services.ErrorPermission {
		t.Log("Remove mod by member: ", err)
		t.Fail()
	}

	if err := Remove(&mod.ID, "modremovebymod", "mod"); err != services.ErrorPermission {
		t.Log("Remove mod by mod: ", err)
		t.Fail()
	}

	if err := Remove(&admin.ID, "modremovebyadmin", "mod"); err != nil {
		t.Log("Remove mod by admin: ", err)
		t.Fail()
	}
}

func TestRemoveAdmin(t *testing.T) {
	if err := Remove(nil, "adminremovebynone", "admin"); err != services.ErrorPermission {
		t.Log("Remove admin by none: ", err)
		t.Fail()
	}

	if err := Remove(&member.ID, "adminremovebymember", "admin"); err != services.ErrorPermission {
		t.Log("Remove admin by member: ", err)
		t.Fail()
	}

	if err := Remove(&mod.ID, "adminremovebymod", "admin"); err != services.ErrorPermission {
		t.Log("Remove admin by mod: ", err)
		t.Fail()
	}

	if err := Remove(&admin.ID, "adminremovebyadmin", "admin"); err != services.ErrorPermission {
		t.Log("Remove admin by admin: ", err)
		t.Fail()
	}
}

func TestRemoveSelf(t *testing.T) {
	if err = SelfRemove("member"); err != nil {
		t.Log("Member self removes: ", err)
		t.Fail()
	}

	if err = SelfRemove("mod"); err != nil {
		t.Log("Mod self removes: ", err)
		t.Fail()
	}

	if err = SelfRemove("admin"); err != nil {
		t.Log("Admin self removes: ", err)
		t.Fail()
	}
}

func TestAuthenticateTrue(t *testing.T) {
	userService := services.CreateUserService(nil)

	r, err := userService.Authenticate(adminUsn, adminPwd)
	if err != nil {
		t.Log("Admin authenticates: ", err)
		t.Fail()
	}

	if r.ID != admin.ID {
		t.Log("Admin authenticates: different result's id")
		t.Fail()
	}

	r, err = userService.Authenticate(modUsn, modPwd)
	if err != nil {
		t.Log("Mod authenticates: ", err)
		t.Fail()
	}

	if r.ID != mod.ID {
		t.Log("Mod authenticates: different result's id")
		t.Fail()
	}

	r, err = userService.Authenticate(memUsn, memPwd)
	if err != nil {
		t.Log("Member authenticates: ", err)
		t.Fail()
	}

	if r.ID != member.ID {
		t.Log("Member authenticates: different result's id")
		t.Fail()
	}
}

func TestAuthenticateFalse(t *testing.T) {
	userService := services.CreateUserService(nil)

	if _, err = userService.Authenticate(adminUsn, modPwd); err != services.ErrorFailedAuthentication {
		t.Log("Authentication is expected to be fail but something wrongs: ", err)
		t.Fail()
	}
}

func TestAuthenticateByIdTrue(t *testing.T) {
	userService := services.CreateUserService(nil)

	r, err := userService.AuthenticateById(admin.ID, adminPwd)
	if err != nil {
		t.Log("Admin authenticates by id: ", err)
		t.Fail()
	}

	if r.ID != admin.ID {
		t.Log("Admin authenticates by id: different result's id")
		t.Fail()
	}

	r, err = userService.AuthenticateById(mod.ID, modPwd)
	if err != nil {
		t.Log("Mod authenticates by id: ", err)
		t.Fail()
	}

	if r.ID != mod.ID {
		t.Log("Mod authenticates: different result's id")
		t.Fail()
	}

	r, err = userService.AuthenticateById(member.ID, memPwd)
	if err != nil {
		t.Log("Member authenticates: ", err)
		t.Fail()
	}

	if r.ID != member.ID {
		t.Log("Member authenticates: different result's id")
		t.Fail()
	}
}

func TestAuthenticateByIdFalse(t *testing.T) {
	userService := services.CreateUserService(nil)

	if _, err = userService.AuthenticateById(admin.ID, modPwd); err != services.ErrorFailedAuthentication {
		t.Log("AuthenticateById is expected to be fail but something wrongs: ", err)
		t.Fail()
	}
}

func TestSelectMember(t *testing.T) {
	if err = Select(nil, "member"); err != services.ErrorPermission {
		t.Log("None selects member: ", err)
		t.Fail()
	}

	if err = Select(&member.ID, "member"); err != services.ErrorPermission {
		t.Log("Member selects member: ", err)
		t.Fail()
	}

	if err = Select(&mod.ID, "member"); err != nil {
		t.Log("Mod selects member: ", err)
		t.Fail()
	}

	if err = Select(&admin.ID, "member"); err != nil {
		t.Log("Admin selects member: ", err)
		t.Fail()
	}
}

func TestSelectMod(t *testing.T) {
	if err = Select(nil, "mod"); err != services.ErrorPermission {
		t.Log("None selects mod: ", err)
		t.Fail()
	}

	if err = Select(&member.ID, "mod"); err != services.ErrorPermission {
		t.Log("Member selects mod: ", err)
		t.Fail()
	}

	if err = Select(&mod.ID, "mod"); err != services.ErrorPermission {
		t.Log("Mod selects mod: ", err)
		t.Fail()
	}

	if err = Select(&admin.ID, "mod"); err != nil {
		t.Log("Admin selects mod: ", err)
		t.Fail()
	}
}

func TestSelectAdmin(t *testing.T) {
	if err = Select(nil, "admin"); err != services.ErrorPermission {
		t.Log("None selects admin: ", err)
		t.Fail()
	}

	if err = Select(&member.ID, "admin"); err != services.ErrorPermission {
		t.Log("Member selects admin: ", err)
		t.Fail()
	}

	if err = Select(&mod.ID, "admin"); err != services.ErrorPermission {
		t.Log("Mod selects admin: ", err)
		t.Fail()
	}

	if err = Select(&admin.ID, "admin"); err != services.ErrorPermission {
		t.Log("Admin selects admin: ", err)
		t.Fail()
	}
}

func TestSelectSelf(t *testing.T) {
	if err := SelfSelect(nil); err != services.ErrorPermission {
		t.Log("None selects self: ", err)
		t.Fail()
	}

	if err := SelfSelect(&member.ID); err != nil {
		t.Log("Member selects self: ", err)
		t.Fail()
	}

	if err := SelfSelect(&mod.ID); err != nil {
		t.Log("Mod selects self: ", err)
		t.Fail()
	}

	if err := SelfSelect(&admin.ID); err != nil {
		t.Log("Admin selects self: ", err)
		t.Fail()
	}
}

func TestUpdateInfoMember(t *testing.T) {
	age := uint(10)
	gender := "male"

	if err = UpdateInfo(nil, "noneupdateinfomember", "member", &age, &gender); err != services.ErrorPermission {
		t.Log("None updates info member: ", err)
		t.Fail()
	}

	if err = UpdateInfo(&member.ID, "memberupdateinfomember", "member", &age, &gender); err != services.ErrorPermission {
		t.Log("Member updates info member: ", err)
		t.Fail()
	}

	if err = UpdateInfo(&mod.ID, "modupdateinfomember", "member", &age, &gender); err != nil {
		t.Log("Mod updates info member: ", err)
		t.Fail()
	}

	if err = UpdateInfo(&admin.ID, "adminupdateinfomember", "member", &age, &gender); err != nil {
		t.Log("Admin updates info member: ", err)
		t.Fail()
	}
}

func TestUpdateInfoMod(t *testing.T) {
	age := uint(10)
	gender := "male"

	if err = UpdateInfo(nil, "noneupdateinfomod", "mod", &age, &gender); err != services.ErrorPermission {
		t.Log("None updates info mod: ", err)
		t.Fail()
	}

	if err = UpdateInfo(&member.ID, "memberupdateinfomod", "mod", &age, &gender); err != services.ErrorPermission {
		t.Log("Member updates info mod: ", err)
		t.Fail()
	}

	if err = UpdateInfo(&mod.ID, "modupdateinfomod", "mod", &age, &gender); err != services.ErrorPermission {
		t.Log("Mod updates info mod: ", err)
		t.Fail()
	}

	if err = UpdateInfo(&admin.ID, "adminupdateinfomod", "mod", &age, &gender); err != nil {
		t.Log("Admin updates info mod: ", err)
		t.Fail()
	}
}

func TestUpdateInfoAdmin(t *testing.T) {
	age := uint(10)
	gender := "male"

	if err = UpdateInfo(nil, "noneupdateinfoadmin", "admin", &age, &gender); err != services.ErrorPermission {
		t.Log("None updates info admin: ", err)
		t.Fail()
	}

	if err = UpdateInfo(&member.ID, "memberupdateinfoadmin", "admin", &age, &gender); err != services.ErrorPermission {
		t.Log("Member updates info admin: ", err)
		t.Fail()
	}

	if err = UpdateInfo(&mod.ID, "modupdateinfomadmin", "admin", &age, &gender); err != services.ErrorPermission {
		t.Log("Mod updates info admin: ", err)
		t.Fail()
	}

	if err = UpdateInfo(&admin.ID, "adminupdateinfoadmin", "admin", &age, &gender); err != services.ErrorPermission {
		t.Log("Admin updates info admin: ", err)
		t.Fail()
	}
}

func TestUpdateInfoWithNull(t *testing.T) {
	age := uint(10)
	gender := "gay"

	if err = UpdateInfo(&admin.ID, "adminupdateinfomemberagenull", "member", nil, &gender); err != nil {
		t.Log("Admin updates info member with age null: ", err)
		t.Fail()
	}

	if err = UpdateInfo(&mod.ID, "modupdateinfomembergendernull", "member", &age, nil); err != nil {
		t.Log("Mod updates info member with gender null: ", err)
		t.Fail()
	}
}

func TestSelfUpdateInfo(t *testing.T) {
	if err = SelfUpdateInfo(&member.ID); err != nil {
		t.Log("Member self updates info: ", err)
		t.Fail()
	}

	if err = SelfUpdateInfo(&mod.ID); err != nil {
		t.Log("Mod self updates info: ", err)
		t.Fail()
	}

	if err = SelfUpdateInfo(&admin.ID); err != nil {
		t.Log("Admin self updates info: ", err)
		t.Fail()
	}
}

func TestUpdateRole(t *testing.T) {
	if err = UpdateRole(nil, "noneupdaterolemembertomod", "member", "mod"); err != services.ErrorPermission {
		t.Log("None updates role member to mod: ", err)
		t.Fail()
	}

	if err = UpdateRole(&member.ID, "memberupdaterolemembertomod", "member", "mod"); err != services.ErrorPermission {
		t.Log("Member updates role member to mod: ", err)
		t.Fail()
	}

	if err = UpdateRole(&mod.ID, "modupdaterolemembertomod", "member", "mod"); err != services.ErrorPermission {
		t.Log("Mod updates role member to mod: ", err)
		t.Fail()
	}

	if err = UpdateRole(&admin.ID, "adminupdaterolemembertomod", "member", "mod"); err != nil {
		t.Log("Admin updates role member to mod: ", err)
		t.Fail()
	}

	if err = UpdateRole(&mod.ID, "modupdaterolemembertoadmin", "member", "admin"); err != services.ErrorPermission {
		t.Log("Mod updates role member to admin: ", err)
		t.Fail()
	}

	if err = UpdateRole(&admin.ID, "adminupdaterolemembertoadmin", "member", "admin"); err != services.ErrorPermission {
		t.Log("Admin updates role member to admin: ", err)
		t.Fail()
	}
}

func TestSelfUpdateRole(t *testing.T) {
	if err = SelfUpdateInfo(&member.ID); err != nil {
		t.Log("Member self updates role: ", err)
		t.Fail()
	}

	if err = SelfUpdateInfo(&mod.ID); err != nil {
		t.Log("Member self updates role: ", err)
		t.Fail()
	}

	if err = SelfUpdateInfo(&mod.ID); err != nil {
		t.Log("Member self updates role: ", err)
		t.Fail()
	}
}

func TestUpdatePasswordMember(t *testing.T) {
	err = UpdatePassword(nil, "noneupdatepaswordmember", "member")
	if err != services.ErrorPermission {
		t.Log("None updates password member: ", err)
		t.Fail()
	}

	err = UpdatePassword(&member.ID, "memberupdatepaswordmember", "member")
	if err != services.ErrorPermission {
		t.Log("Member updates password havmember: ", err)
		t.Fail()
	}

	err = UpdatePassword(&mod.ID, "modupdatepaswordmember", "member")
	if err != nil {
		t.Log("Mod updates password member: ", err)
		t.Fail()
	}

	err = UpdatePassword(&admin.ID, "adminupdatepaswordmember", "member")
	if err != nil {
		t.Log("Admin updates password member: ", err)
		t.Fail()
	}
}

func TestUpdatePasswordMod(t *testing.T) {
	err = UpdatePassword(nil, "noneupdatepaswordmod", "mod")
	if err != services.ErrorPermission {
		t.Log("None updates password mod: ", err)
		t.Fail()
	}

	err = UpdatePassword(&member.ID, "memberupdatepaswordmod", "mod")
	if err != services.ErrorPermission {
		t.Log("Member updates password mod: ", err)
		t.Fail()
	}

	err = UpdatePassword(&mod.ID, "modupdatepaswordmod", "mod")
	if err != services.ErrorPermission {
		t.Log("Mod updates password mod: ", err)
		t.Fail()
	}

	err = UpdatePassword(&admin.ID, "adminupdatepaswordmod", "mod")
	if err != nil {
		t.Log("Admin updates password mod: ", err)
		t.Fail()
	}
}

func TestUpdatePasswordAdmin(t *testing.T) {
	err = UpdatePassword(nil, "noneupdatepaswordadmin", "admin")
	if err != services.ErrorPermission {
		t.Log("None updates password admin: ", err)
		t.Fail()
	}

	err = UpdatePassword(&member.ID, "memberupdatepaswordadmin", "admin")
	if err != services.ErrorPermission {
		t.Log("Member updates password admin: ", err)
		t.Fail()
	}

	err = UpdatePassword(&mod.ID, "modupdatepaswordadmin", "admin")
	if err != services.ErrorPermission {
		t.Log("Mod updates password admin: ", err)
		t.Fail()
	}

	err = UpdatePassword(&admin.ID, "adminupdatepaswordadmin", "admin")
	if err != services.ErrorPermission {
		t.Log("Admin updates password admin: ", err)
		t.Fail()
	}
}

func TestSelfUpdatePassword(t *testing.T) {
	err = SelfUpdatePassword("memberselfupdatepassword", "member")
	if err != nil {
		t.Log("Member self updates password: ", err)
		t.Fail()
	}

	err = SelfUpdatePassword("modselfupdatepassword", "mod")
	if err != nil {
		t.Log("Mod self updates password: ", err)
		t.Fail()
	}

	err = SelfUpdatePassword("Adminselfupdatepassword", "admin")
	if err != nil {
		t.Log("Admin self updates password: ", err)
		t.Fail()
	}
}
