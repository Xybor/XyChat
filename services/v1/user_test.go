package v1_test

import (
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/xybor/xychat/helpers"
	"github.com/xybor/xychat/models"
	services "github.com/xybor/xychat/services/v1"
	xyerrors "github.com/xybor/xychat/xyerrors/v1"
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

func CreateUser(username, password, role string) (models.User, error) {
	user := models.User{Username: &username, Password: &password, Role: &role}
	err := models.GetDB().Create(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func Register(id *uint, name, role string) xyerrors.XyError {
	userService := services.CreateUserService(id, true)
	xerr := userService.Register(name, "password", role)
	if xerr.Errno() != 0 {
		return xerr
	}

	return xyerrors.NoError
}

func Remove(id *uint, name, role string) xyerrors.XyError {
	var user models.User
	var err error

	if user, err = CreateUser(name, "password", role); err != nil {
		return xyerrors.ErrorUnknown
	}

	userService := services.CreateUserService(id, true)

	xerr1 := userService.RemoveByUsername(name)
	if xerr1.Errno() == 0 {
		if user, err = CreateUser(name, "password", role); err != nil {
			return xyerrors.ErrorUnknown
		}
	}

	xerr2 := userService.Remove(user.ID)

	if xerr1.Errno() != 0 {
		return xerr1
	}

	if xerr2.Errno() != 0 {
		return xerr2
	}

	return xyerrors.NoError
}

func SelfRemove(role string) xyerrors.XyError {
	var user models.User
	var err error

	if user, err = CreateUser(role+"removeself1", "password", role); err != nil {
		return xyerrors.ErrorUnknown
	}

	userService := services.CreateUserService(&user.ID, true)

	xerr := userService.RemoveByUsername(*user.Username)
	if xerr.Errno() == 0 {
		return xerr
	}

	if user, err = CreateUser(role+"removeself2", "password", role); err != nil {
		return xyerrors.ErrorUnknown
	}

	xerr = userService.Remove(user.ID)
	if xerr.Errno() != 0 {
		return xerr
	}

	return xyerrors.NoError
}

func Select(id *uint, role string) xyerrors.XyError {
	userService := services.CreateUserService(id, true)
	var user models.User

	usn := "select" + role
	if id == nil {
		usn = "none" + usn
	} else {
		usn = strconv.FormatUint(uint64(*id), 10) + usn
	}

	var err error
	if user, err = CreateUser(usn, "password", role); err != nil {
		return xyerrors.ErrorUnknown
	}

	r, xerr := userService.Select(user.ID)
	if xerr.Errno() != 0 {
		return xerr
	}

	if r.ID != user.ID {
		return xyerrors.ErrorUnknown.New("Different result id")
	}

	r, xerr = userService.SelectByName(usn)
	if xerr.Errno() != 0 {
		return xerr
	}

	if r.ID != user.ID {
		return xyerrors.ErrorUnknown.New("Different result id")
	}

	return xyerrors.NoError
}

func SelfSelect(id *uint) xyerrors.XyError {
	userService := services.CreateUserService(id, true)

	r, xerr := userService.SelfSelect()
	if xerr.Errno() != 0 {
		return xerr
	}

	if id == nil || r.ID != *id {
		return xyerrors.ErrorUnknown.New("Different result id")
	}

	return xyerrors.NoError
}

func UpdateInfo(id *uint, username, role string, age *uint, gender *string) xyerrors.XyError {
	var user models.User
	var err error

	if user, err = CreateUser(username, "password", role); err != nil {
		return xyerrors.ErrorUnknown
	}

	userService := services.CreateUserService(id, true)

	if xerr := userService.UpdateInfo(user.ID, age, gender); xerr.Errno() != 0 {
		return xerr
	}

	return xyerrors.NoError
}

func SelfUpdateInfo(id *uint) xyerrors.XyError {
	age := uint(10)
	gender := "female"

	userService := services.CreateUserService(id, true)

	if xerr := userService.UpdateInfo(*id, &age, &gender); xerr.Errno() != 0 {
		return xerr
	}

	return xyerrors.NoError
}

func UpdateRole(id *uint, username, role, newrole string) xyerrors.XyError {
	var user models.User
	var err error

	if user, err = CreateUser(username, "password", role); err != nil {
		return xyerrors.ErrorUnknown
	}

	userService := services.CreateUserService(id, true)

	if xerr := userService.UpdateRole(user.ID, newrole); xerr.Errno() != 0 {
		return xerr
	}

	return xyerrors.NoError
}

func SelfUpdateRole(id *uint) xyerrors.XyError {
	userService := services.CreateUserService(id, true)

	if xerr := userService.UpdateRole(*id, "member"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		if xerr.Errno() == 0 {
			return xyerrors.ErrorUnknown.New("invalid self update to member")
		}
		return xerr
	}

	if xerr := userService.UpdateRole(*id, "mod"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		if xerr.Errno() == 0 {
			return xyerrors.ErrorUnknown.New("invalid self update to mod")
		}
		return xerr
	}

	if xerr := userService.UpdateRole(*id, "admin"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		if xerr.Errno() == 0 {
			return xyerrors.ErrorUnknown.New("invalid self update to admin")
		}
		return xerr
	}

	return xyerrors.NoError
}

func UpdatePassword(id *uint, username, role string) xyerrors.XyError {
	var user models.User
	var err error

	password := "password"

	if user, err = CreateUser(username, password, role); err != nil {
		return xyerrors.ErrorUnknown
	}

	userService := services.CreateUserService(id, true)

	var xerr1, xerr2 xyerrors.XyError

	xerr1 = userService.UpdatePassword(user.ID, &password, "newpassword")
	xerr2 = userService.UpdatePassword(user.ID, nil, "newpassword")

	if xerr1.Errno() != xerr2.Errno() {
		return xyerrors.ErrorUnknown.New("Error from with and without password is different; " +
			xerr1.Error() + ";" + xerr2.Error(),
		)
	}

	return xerr1
}

func SelfUpdatePassword(username, role string) xyerrors.XyError {
	var user models.User
	var err error

	password := "password"

	if user, err = CreateUser(username, password, role); err != nil {
		return xyerrors.ErrorUnknown
	}

	userService := services.CreateUserService(&user.ID, true)

	if xerr := userService.UpdatePassword(user.ID, nil, "newpassword"); xerr.Errno() != xyerrors.ErrorFailedAuthentication.Errno() {
		if xerr.Errno() == 0 {
			return xyerrors.ErrorUnknown.New("self change password without old password")
		}
		return xerr
	}

	if xerr := userService.UpdatePassword(user.ID, &password, "newpassword"); xerr.Errno() != 0 {
		return xerr
	}

	return xyerrors.NoError
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
	var err error

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
	if xerr := Register(nil, "memberbynone", "member"); xerr.Errno() != 0 {
		t.Log("Register member by none: ", xerr)
		t.Fail()
	}

	if xerr := Register(&member.ID, "memberbymember", "member"); xerr.Errno() != 0 {
		t.Log("Register member by member: ", xerr)
		t.Fail()
	}

	if xerr := Register(&mod.ID, "memberbymod", "member"); xerr.Errno() != 0 {
		t.Log("Register member by mod: ", xerr)
		t.Fail()
	}

	if xerr := Register(&admin.ID, "memberbyadmin", "member"); xerr.Errno() != 0 {
		t.Log("Register member by admin: ", xerr)
		t.Fail()
	}
}

func TestRegisterMod(t *testing.T) {
	if xerr := Register(nil, "modbynone", "mod"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Register mod by none: ", xerr)
		t.Fail()
	}

	if xerr := Register(&member.ID, "modbymember", "mod"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Register mod by member: ", xerr)
		t.Fail()
	}

	if xerr := Register(&mod.ID, "modbymod", "mod"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Register mod by mod: ", xerr)
		t.Fail()
	}

	if xerr := Register(&admin.ID, "modbyadmin", "mod"); xerr.Errno() != 0 {
		t.Log("Register mod by admin: ", xerr)
		t.Fail()
	}
}

func TestRegisterAdmin(t *testing.T) {
	if xerr := Register(nil, "adminbynone", "admin"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Register admin by none: ", xerr)
		t.Fail()
	}

	if xerr := Register(&member.ID, "adminbymember", "admin"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Register admin by member: ", xerr)
		t.Fail()
	}

	if xerr := Register(&mod.ID, "adminbymod", "admin"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Register admin by mod: ", xerr)
		t.Fail()
	}

	if xerr := Register(&admin.ID, "adminbyadmin", "admin"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Register admin by admin: ", xerr)
		t.Fail()
	}
}

func TestRegisterUnknownRole(t *testing.T) {
	if xerr := Register(nil, "userunknownrole", "something"); xerr.Errno() != xyerrors.ErrorUnknownInput.Errno() {
		t.Log("Register unknown role: ", xerr)
	}
}

func TestRegisterDuplication(t *testing.T) {
	if xerr := Register(nil, "memberbynonedup", "member"); xerr.Errno() != 0 {
		t.Log("Register member by none dup 1: ", xerr)
		t.Fail()
	}

	if xerr := Register(nil, "memberbynonedup", "member"); xerr.Errno() != xyerrors.ErrorExistedUsername.Errno() {
		t.Log("Register member by none dup 2: ", xerr)
		t.Fail()
	}
}

func TestRemoveMember(t *testing.T) {
	if xerr := Remove(nil, "memberremovebynone", "member"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Remove member by none: ", xerr)
		t.Fail()
	}

	if xerr := Remove(&member.ID, "memberremovebymember", "member"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Remove member by member: ", xerr)
		t.Fail()
	}

	if xerr := Remove(&mod.ID, "memberremovebymod", "member"); xerr.Errno() != 0 {
		t.Log("Remove member by mod: ", xerr)
		t.Fail()
	}

	if xerr := Remove(&admin.ID, "memberremovebyadmin", "member"); xerr.Errno() != 0 {
		t.Log("Remove member by admin: ", xerr)
		t.Fail()
	}
}
func TestRemoveMod(t *testing.T) {
	if xerr := Remove(nil, "modremovebynone", "mod"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Remove mod by none: ", xerr)
		t.Fail()
	}

	if xerr := Remove(&member.ID, "modremovebymember", "mod"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Remove mod by member: ", xerr)
		t.Fail()
	}

	if xerr := Remove(&mod.ID, "modremovebymod", "mod"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Remove mod by mod: ", xerr)
		t.Fail()
	}

	if xerr := Remove(&admin.ID, "modremovebyadmin", "mod"); xerr.Errno() != 0 {
		t.Log("Remove mod by admin: ", xerr)
		t.Fail()
	}
}

func TestRemoveAdmin(t *testing.T) {
	if xerr := Remove(nil, "adminremovebynone", "admin"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Remove admin by none: ", xerr)
		t.Fail()
	}

	if xerr := Remove(&member.ID, "adminremovebymember", "admin"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Remove admin by member: ", xerr)
		t.Fail()
	}

	if xerr := Remove(&mod.ID, "adminremovebymod", "admin"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Remove admin by mod: ", xerr)
		t.Fail()
	}

	if err := Remove(&admin.ID, "adminremovebyadmin", "admin"); err.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Remove admin by admin: ", err)
		t.Fail()
	}
}

func TestRemoveSelf(t *testing.T) {
	if xerr := SelfRemove("member"); xerr.Errno() != 0 {
		t.Log("Member self removes: ", xerr)
		t.Fail()
	}

	if xerr := SelfRemove("mod"); xerr.Errno() != 0 {
		t.Log("Mod self removes: ", xerr)
		t.Fail()
	}

	if xerr := SelfRemove("admin"); xerr.Errno() != 0 {
		t.Log("Admin self removes: ", xerr)
		t.Fail()
	}
}

func TestAuthenticateTrue(t *testing.T) {
	userService := services.CreateUserService(nil, true)

	r, xerr := userService.Authenticate(adminUsn, adminPwd)
	if xerr.Errno() != 0 {
		t.Log("Admin authenticates: ", xerr)
		t.Fail()
	}

	if r.ID != admin.ID {
		t.Log("Admin authenticates: different result's id")
		t.Fail()
	}

	r, xerr = userService.Authenticate(modUsn, modPwd)
	if xerr.Errno() != 0 {
		t.Log("Mod authenticates: ", xerr)
		t.Fail()
	}

	if r.ID != mod.ID {
		t.Log("Mod authenticates: different result's id")
		t.Fail()
	}

	r, xerr = userService.Authenticate(memUsn, memPwd)
	if xerr.Errno() != 0 {
		t.Log("Member authenticates: ", xerr)
		t.Fail()
	}

	if r.ID != member.ID {
		t.Log("Member authenticates: different result's id")
		t.Fail()
	}
}

func TestAuthenticateFalse(t *testing.T) {
	userService := services.CreateUserService(nil, true)

	if _, xerr := userService.Authenticate(adminUsn, modPwd); xerr.Errno() != xyerrors.ErrorFailedAuthentication.Errno() {
		t.Log("Authentication is expected to be fail but something wrongs: ", xerr)
		t.Fail()
	}
}

func TestAuthenticateByIdTrue(t *testing.T) {
	userService := services.CreateUserService(nil, true)

	r, xerr := userService.AuthenticateById(admin.ID, adminPwd)
	if xerr.Errno() != 0 {
		t.Log("Admin authenticates by id: ", xerr)
		t.Fail()
	}

	if r.ID != admin.ID {
		t.Log("Admin authenticates by id: different result's id")
		t.Fail()
	}

	r, xerr = userService.AuthenticateById(mod.ID, modPwd)
	if xerr.Errno() != 0 {
		t.Log("Mod authenticates by id: ", xerr)
		t.Fail()
	}

	if r.ID != mod.ID {
		t.Log("Mod authenticates: different result's id")
		t.Fail()
	}

	r, xerr = userService.AuthenticateById(member.ID, memPwd)
	if xerr.Errno() != 0 {
		t.Log("Member authenticates: ", xerr)
		t.Fail()
	}

	if r.ID != member.ID {
		t.Log("Member authenticates: different result's id")
		t.Fail()
	}
}

func TestAuthenticateByIdFalse(t *testing.T) {
	userService := services.CreateUserService(nil, true)

	if _, xerr := userService.AuthenticateById(admin.ID, modPwd); xerr.Errno() != xyerrors.ErrorFailedAuthentication.Errno() {
		t.Log("AuthenticateById is expected to be fail but something wrongs: ", xerr)
		t.Fail()
	}
}

func TestSelectMember(t *testing.T) {
	if xerr := Select(nil, "member"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("None selects member: ", xerr)
		t.Fail()
	}

	if xerr := Select(&member.ID, "member"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Member selects member: ", xerr)
		t.Fail()
	}

	if xerr := Select(&mod.ID, "member"); xerr.Errno() != 0 {
		t.Log("Mod selects member: ", xerr)
		t.Fail()
	}

	if xerr := Select(&admin.ID, "member"); xerr.Errno() != 0 {
		t.Log("Admin selects member: ", xerr)
		t.Fail()
	}
}

func TestSelectMod(t *testing.T) {
	if xerr := Select(nil, "mod"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("None selects mod: ", xerr)
		t.Fail()
	}

	if xerr := Select(&member.ID, "mod"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Member selects mod: ", xerr)
		t.Fail()
	}

	if xerr := Select(&mod.ID, "mod"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Mod selects mod: ", xerr)
		t.Fail()
	}

	if xerr := Select(&admin.ID, "mod"); xerr.Errno() != 0 {
		t.Log("Admin selects mod: ", xerr)
		t.Fail()
	}
}

func TestSelectAdmin(t *testing.T) {
	if xerr := Select(nil, "admin"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("None selects admin: ", xerr)
		t.Fail()
	}

	if xerr := Select(&member.ID, "admin"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Member selects admin: ", xerr)
		t.Fail()
	}

	if xerr := Select(&mod.ID, "admin"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Mod selects admin: ", xerr)
		t.Fail()
	}

	if xerr := Select(&admin.ID, "admin"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Admin selects admin: ", xerr)
		t.Fail()
	}
}

func TestSelectSelf(t *testing.T) {
	if err := SelfSelect(nil); err.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("None selects self: ", err)
		t.Fail()
	}

	if err := SelfSelect(&member.ID); err.Errno() != 0 {
		t.Log("Member selects self: ", err)
		t.Fail()
	}

	if err := SelfSelect(&mod.ID); err.Errno() != 0 {
		t.Log("Mod selects self: ", err)
		t.Fail()
	}

	if err := SelfSelect(&admin.ID); err.Errno() != 0 {
		t.Log("Admin selects self: ", err)
		t.Fail()
	}
}

func TestUpdateInfoMember(t *testing.T) {
	age := uint(10)
	gender := "male"

	if xerr := UpdateInfo(nil, "noneupdateinfomember", "member", &age, &gender); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("None updates info member: ", xerr)
		t.Fail()
	}

	if xerr := UpdateInfo(&member.ID, "memberupdateinfomember", "member", &age, &gender); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Member updates info member: ", xerr)
		t.Fail()
	}

	if xerr := UpdateInfo(&mod.ID, "modupdateinfomember", "member", &age, &gender); xerr.Errno() != 0 {
		t.Log("Mod updates info member: ", xerr)
		t.Fail()
	}

	if xerr := UpdateInfo(&admin.ID, "adminupdateinfomember", "member", &age, &gender); xerr.Errno() != 0 {
		t.Log("Admin updates info member: ", xerr)
		t.Fail()
	}
}

func TestUpdateInfoMod(t *testing.T) {
	age := uint(10)
	gender := "male"

	if xerr := UpdateInfo(nil, "noneupdateinfomod", "mod", &age, &gender); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("None updates info mod: ", xerr)
		t.Fail()
	}

	if xerr := UpdateInfo(&member.ID, "memberupdateinfomod", "mod", &age, &gender); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Member updates info mod: ", xerr)
		t.Fail()
	}

	if xerr := UpdateInfo(&mod.ID, "modupdateinfomod", "mod", &age, &gender); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Mod updates info mod: ", xerr)
		t.Fail()
	}

	if xerr := UpdateInfo(&admin.ID, "adminupdateinfomod", "mod", &age, &gender); xerr.Errno() != 0 {
		t.Log("Admin updates info mod: ", xerr)
		t.Fail()
	}
}

func TestUpdateInfoAdmin(t *testing.T) {
	age := uint(10)
	gender := "male"

	if xerr := UpdateInfo(nil, "noneupdateinfoadmin", "admin", &age, &gender); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("None updates info admin: ", xerr)
		t.Fail()
	}

	if xerr := UpdateInfo(&member.ID, "memberupdateinfoadmin", "admin", &age, &gender); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Member updates info admin: ", xerr)
		t.Fail()
	}

	if xerr := UpdateInfo(&mod.ID, "modupdateinfomadmin", "admin", &age, &gender); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Mod updates info admin: ", xerr)
		t.Fail()
	}

	if xerr := UpdateInfo(&admin.ID, "adminupdateinfoadmin", "admin", &age, &gender); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Admin updates info admin: ", xerr)
		t.Fail()
	}
}

func TestUpdateInfoWithNull(t *testing.T) {
	age := uint(10)
	gender := "gay"

	if xerr := UpdateInfo(&admin.ID, "adminupdateinfomemberagenull", "member", nil, &gender); xerr.Errno() != 0 {
		t.Log("Admin updates info member with age null: ", xerr)
		t.Fail()
	}

	if xerr := UpdateInfo(&mod.ID, "modupdateinfomembergendernull", "member", &age, nil); xerr.Errno() != 0 {
		t.Log("Mod updates info member with gender null: ", xerr)
		t.Fail()
	}
}

func TestSelfUpdateInfo(t *testing.T) {
	if xerr := SelfUpdateInfo(&member.ID); xerr.Errno() != 0 {
		t.Log("Member self updates info: ", xerr)
		t.Fail()
	}

	if xerr := SelfUpdateInfo(&mod.ID); xerr.Errno() != 0 {
		t.Log("Mod self updates info: ", xerr)
		t.Fail()
	}

	if xerr := SelfUpdateInfo(&admin.ID); xerr.Errno() != 0 {
		t.Log("Admin self updates info: ", xerr)
		t.Fail()
	}
}

func TestUpdateRole(t *testing.T) {
	if xerr := UpdateRole(nil, "noneupdaterolemembertomod", "member", "mod"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("None updates role member to mod: ", xerr)
		t.Fail()
	}

	if xerr := UpdateRole(&member.ID, "memberupdaterolemembertomod", "member", "mod"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Member updates role member to mod: ", xerr)
		t.Fail()
	}

	if xerr := UpdateRole(&mod.ID, "modupdaterolemembertomod", "member", "mod"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Mod updates role member to mod: ", xerr)
		t.Fail()
	}

	if xerr := UpdateRole(&admin.ID, "adminupdaterolemembertomod", "member", "mod"); xerr.Errno() != 0 {
		t.Log("Admin updates role member to mod: ", xerr)
		t.Fail()
	}

	if xerr := UpdateRole(&mod.ID, "modupdaterolemembertoadmin", "member", "admin"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Mod updates role member to admin: ", xerr)
		t.Fail()
	}

	if xerr := UpdateRole(&admin.ID, "adminupdaterolemembertoadmin", "member", "admin"); xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Admin updates role member to admin: ", xerr)
		t.Fail()
	}
}

func TestSelfUpdateRole(t *testing.T) {
	if xerr := SelfUpdateInfo(&member.ID); xerr.Errno() != 0 {
		t.Log("Member self updates role: ", xerr)
		t.Fail()
	}

	if xerr := SelfUpdateInfo(&mod.ID); xerr.Errno() != 0 {
		t.Log("Member self updates role: ", xerr)
		t.Fail()
	}

	if xerr := SelfUpdateInfo(&mod.ID); xerr.Errno() != 0 {
		t.Log("Member self updates role: ", xerr)
		t.Fail()
	}
}

func TestUpdatePasswordMember(t *testing.T) {
	xerr := UpdatePassword(nil, "noneupdatepaswordmember", "member")
	if xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("None updates password member: ", xerr)
		t.Fail()
	}

	xerr = UpdatePassword(&member.ID, "memberupdatepaswordmember", "member")
	if xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Member updates password havmember: ", xerr)
		t.Fail()
	}

	xerr = UpdatePassword(&mod.ID, "modupdatepaswordmember", "member")
	if xerr.Errno() != 0 {
		t.Log("Mod updates password member: ", xerr)
		t.Fail()
	}

	xerr = UpdatePassword(&admin.ID, "adminupdatepaswordmember", "member")
	if xerr.Errno() != 0 {
		t.Log("Admin updates password member: ", xerr)
		t.Fail()
	}
}

func TestUpdatePasswordMod(t *testing.T) {
	xerr := UpdatePassword(nil, "noneupdatepaswordmod", "mod")
	if xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("None updates password mod: ", xerr)
		t.Fail()
	}

	xerr = UpdatePassword(&member.ID, "memberupdatepaswordmod", "mod")
	if xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Member updates password mod: ", xerr)
		t.Fail()
	}

	xerr = UpdatePassword(&mod.ID, "modupdatepaswordmod", "mod")
	if xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Mod updates password mod: ", xerr)
		t.Fail()
	}

	xerr = UpdatePassword(&admin.ID, "adminupdatepaswordmod", "mod")
	if xerr.Errno() != 0 {
		t.Log("Admin updates password mod: ", xerr)
		t.Fail()
	}
}

func TestUpdatePasswordAdmin(t *testing.T) {
	xerr := UpdatePassword(nil, "noneupdatepaswordadmin", "admin")
	if xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("None updates password admin: ", xerr)
		t.Fail()
	}

	xerr = UpdatePassword(&member.ID, "memberupdatepaswordadmin", "admin")
	if xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Member updates password admin: ", xerr)
		t.Fail()
	}

	xerr = UpdatePassword(&mod.ID, "modupdatepaswordadmin", "admin")
	if xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Mod updates password admin: ", xerr)
		t.Fail()
	}

	xerr = UpdatePassword(&admin.ID, "adminupdatepaswordadmin", "admin")
	if xerr.Errno() != xyerrors.ErrorPermission.Errno() {
		t.Log("Admin updates password admin: ", xerr)
		t.Fail()
	}
}

func TestSelfUpdatePassword(t *testing.T) {
	xerr := SelfUpdatePassword("memberselfupdatepassword", "member")
	if xerr.Errno() != 0 {
		t.Log("Member self updates password: ", xerr)
		t.Fail()
	}

	xerr = SelfUpdatePassword("modselfupdatepassword", "mod")
	if xerr.Errno() != 0 {
		t.Log("Mod self updates password: ", xerr)
		t.Fail()
	}

	xerr = SelfUpdatePassword("Adminselfupdatepassword", "admin")
	if xerr.Errno() != 0 {
		t.Log("Admin self updates password: ", xerr)
		t.Fail()
	}
}
