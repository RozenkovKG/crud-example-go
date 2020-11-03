package repository

import (
	"crud-example-go/src/dto"
	"crud-example-go/src/intreface/db"
	"testing"
)

type testDB struct {
	FindFunc     func(dest interface{}, conds ...interface{}) (tx db.DB)
	FirstFunc    func(dest interface{}, conds ...interface{}) (tx db.DB)
	SaveFunc     func(value interface{}) (tx db.DB)
	DeleteFunc   func(value interface{}, conds ...interface{}) (tx db.DB)
	PreloadFunc  func(query string, args ...interface{}) (tx db.DB)
	GetErrorFunc func() error
}

func (db *testDB) Find(dest interface{}, conds ...interface{}) (tx db.DB) {
	return db.FindFunc(dest, conds...)
}

func (db *testDB) First(dest interface{}, conds ...interface{}) (tx db.DB) {
	return db.FirstFunc(dest, conds...)
}

func (db *testDB) Save(value interface{}) (tx db.DB) {
	return db.SaveFunc(value)
}

func (db *testDB) Delete(value interface{}, conds ...interface{}) (tx db.DB) {
	return db.DeleteFunc(value, conds...)
}

func (db *testDB) Preload(query string, args ...interface{}) (tx db.DB) {
	return db.PreloadFunc(query, args...)
}

func (db *testDB) GetError() error {
	return db.GetErrorFunc()
}

func TestUserRepository_GetUser(t *testing.T) {
	defer createRecover(t)()

	expectedId := 10
	expectedName := "name_name"

	dataBase := testDB{
		GetErrorFunc: createEmptyGetErrorFunc()}
	dataBase.PreloadFunc = createPreloadFunc(&dataBase)
	dataBase.FirstFunc = func(dest interface{}, conds ...interface{}) (tx db.DB) {
		user := dest.(*User)
		idVal := conds[0].(int)
		user.ID = &idVal
		user.Name = expectedName
		return &dataBase
	}
	dataBase.GetErrorFunc = createEmptyGetErrorFunc()
	repository := UserRepository{DB: &dataBase}

	user, err := repository.GetUser(expectedId)
	if err != nil {
		t.Error(`UserRepository.GetUser(`, expectedId, `) => err :`, err)
	} else {
		if expectedId != *user.ID {
			t.Error(`UserRepository.GetUser(`, expectedId, `) => User.ID !=`, expectedId)
		} else {
			if expectedName != user.Name {
				t.Error(`UserRepository.GetUser(`, expectedId, `) => User.Name !=`, expectedName)
			}
		}
	}
}

func TestUserRepository_GetUsers(t *testing.T) {
	defer createRecover(t)()

	dataBase := testDB{
		GetErrorFunc: createEmptyGetErrorFunc()}
	dataBase.PreloadFunc = createPreloadFunc(&dataBase)
	dataBase.FindFunc = func(dest interface{}, conds ...interface{}) (tx db.DB) {

		users := dest.(*[]User)
		idFirst := 1
		idSecond := 2
		*users = append(*users, User{ID: &idFirst, Name: "user1"}, User{ID: &idSecond, Name: "user2"})

		return &dataBase
	}
	repository := UserRepository{DB: &dataBase}

	users, err := repository.GetUsers()
	if err != nil {
		t.Error(`UserRepository.GetUsers() => err :`, err)
		return
	}
	if len(*users) != 2 {
		t.Error(`UserRepository.GetUsers() => err : count of users is not correct`)
	}
}

func TestUserRepository_SaveUser(t *testing.T) {
	defer createRecover(t)()

	inputUser := dto.User{
		Name: "user_name",
		Tags: []dto.Tag{
			{Name: "tag_name1"},
			{Name: "tag_name2"}}}

	userId := 122
	tagIdFirst := 123
	tagIdSecond := 124
	expectedUser := dto.User(inputUser)
	expectedUser.ID = &userId
	expectedUser.Tags = []dto.Tag{inputUser.Tags[0], inputUser.Tags[1]}
	expectedUser.Tags[0].ID = &tagIdFirst
	expectedUser.Tags[0].UserID = &userId
	expectedUser.Tags[1].ID = &tagIdSecond
	expectedUser.Tags[1].UserID = &userId

	dataBase := testDB{
		GetErrorFunc: createEmptyGetErrorFunc()}
	dataBase.PreloadFunc = createPreloadFunc(&dataBase)
	dataBase.FindFunc = func(dest interface{}, conds ...interface{}) (tx db.DB) {

		tags := dest.(*[]Tag)
		idFirst := 1
		idSecond := 2
		*tags = append(*tags, Tag{ID: &idFirst, Name: "tag1"}, Tag{ID: &idSecond, Name: "tag2"})

		return &dataBase
	}
	dataBase.DeleteFunc = func(value interface{}, conds ...interface{}) (tx db.DB) {
		return &dataBase
	}
	dataBase.SaveFunc = func(value interface{}) (tx db.DB) {
		user := value.(*User)
		user.ID = expectedUser.ID
		user.Tags[0].ID = expectedUser.Tags[0].ID
		user.Tags[0].UserID = expectedUser.ID
		user.Tags[1].ID = expectedUser.Tags[1].ID
		user.Tags[1].UserID = expectedUser.ID
		return &dataBase
	}

	repository := UserRepository{DB: &dataBase}
	outputUser, err := repository.SaveUser(&inputUser)

	if err != nil {
		t.Error(`UserRepository.SaveUsers(user) => err :`, err)
		return
	}
	if *expectedUser.ID != *outputUser.ID {
		t.Error(`UserRepository.SaveUsers(user) => expected ID =`, *expectedUser.ID,
			`, output ID =`, *outputUser.ID)
		return
	}
	if len(expectedUser.Tags) != len(outputUser.Tags) {
		t.Error(`UserRepository.SaveUsers(user) => expected Tags count =`, len(expectedUser.Tags),
			`, output Tags count =`, len(outputUser.Tags))
	}
	if *expectedUser.Tags[0].ID != *outputUser.Tags[0].ID {
		t.Error(`UserRepository.SaveUsers(user) => expected Tag[0].ID =`, *expectedUser.Tags[0].ID,
			`, output Tag[0].ID =`, *outputUser.Tags[0].ID)
	}
	if *expectedUser.Tags[1].ID != *outputUser.Tags[1].ID {
		t.Error(`UserRepository.SaveUsers(user) => expected Tag[1].ID =`, *expectedUser.Tags[1].ID,
			`, output Tag[1].ID =`, *outputUser.Tags[1].ID)
	}
	if *expectedUser.Tags[0].UserID != *outputUser.Tags[0].UserID {
		t.Error(`UserRepository.SaveUsers(user) => expected Tag[0].UserID =`, *expectedUser.Tags[0].UserID,
			`, output Tag[0].UserID =`, *outputUser.Tags[0].UserID)
	}
	if *expectedUser.Tags[1].UserID != *outputUser.Tags[1].UserID {
		t.Error(`UserRepository.SaveUsers(user) => expected Tag[1].UserID =`, *expectedUser.Tags[1].UserID,
			`, output Tag[1].UserID =`, *outputUser.Tags[1].UserID)
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	defer createRecover(t)()

	userId := 3332

	dataBase := testDB{
		GetErrorFunc: createEmptyGetErrorFunc()}
	dataBase.PreloadFunc = createPreloadFunc(&dataBase)
	dataBase.DeleteFunc = func(value interface{}, conds ...interface{}) (tx db.DB) {
		id := conds[0].(int)

		if id != userId {
			t.Error(`UserRepository.SaveUsers(id) => user id =`, userId,
				`input id =`, id)
		}

		return &dataBase
	}

	repository := UserRepository{DB: &dataBase}
	_ = repository.DeleteUser(userId)
}

func createPreloadFunc(dataBase db.DB) func(query string, args ...interface{}) (tx db.DB) {
	return func(query string, args ...interface{}) (tx db.DB) {
		return dataBase
	}
}

func createEmptyGetErrorFunc() func() error {
	return func() error {
		return nil
	}
}

func createRecover(t *testing.T) func() {
	return func() {
		if r := recover(); r != nil {
			t.Error(`Panic:`, r)
		}
	}
}
