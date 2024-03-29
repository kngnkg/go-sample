// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package service

import (
	"context"
	"github.com/kngnkg/go-sample/model"
	"github.com/kngnkg/go-sample/store"
	"sync"
)

// Ensure, that StoreMock does implement Store.
// If this is not the case, regenerate this file with moq.
var _ Store = &StoreMock{}

// StoreMock is a mock implementation of Store.
//
//	func TestSomethingThatUsesStore(t *testing.T) {
//
//		// make and configure a mocked Store
//		mockedStore := &StoreMock{
//			DeleteFunc: func(ctx context.Context, key string) error {
//				panic("mock out the Delete method")
//			},
//			LoadFunc: func(ctx context.Context, key string) (int, error) {
//				panic("mock out the Load method")
//			},
//			SaveFunc: func(ctx context.Context, key string, uid int) error {
//				panic("mock out the Save method")
//			},
//		}
//
//		// use mockedStore in code that requires Store
//		// and then make assertions.
//
//	}
type StoreMock struct {
	// DeleteFunc mocks the Delete method.
	DeleteFunc func(ctx context.Context, key string) error

	// LoadFunc mocks the Load method.
	LoadFunc func(ctx context.Context, key string) (int, error)

	// SaveFunc mocks the Save method.
	SaveFunc func(ctx context.Context, key string, uid int) error

	// calls tracks calls to the methods.
	calls struct {
		// Delete holds details about calls to the Delete method.
		Delete []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Key is the key argument value.
			Key string
		}
		// Load holds details about calls to the Load method.
		Load []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Key is the key argument value.
			Key string
		}
		// Save holds details about calls to the Save method.
		Save []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Key is the key argument value.
			Key string
			// UID is the uid argument value.
			UID int
		}
	}
	lockDelete sync.RWMutex
	lockLoad   sync.RWMutex
	lockSave   sync.RWMutex
}

// Delete calls DeleteFunc.
func (mock *StoreMock) Delete(ctx context.Context, key string) error {
	if mock.DeleteFunc == nil {
		panic("StoreMock.DeleteFunc: method is nil but Store.Delete was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Key string
	}{
		Ctx: ctx,
		Key: key,
	}
	mock.lockDelete.Lock()
	mock.calls.Delete = append(mock.calls.Delete, callInfo)
	mock.lockDelete.Unlock()
	return mock.DeleteFunc(ctx, key)
}

// DeleteCalls gets all the calls that were made to Delete.
// Check the length with:
//
//	len(mockedStore.DeleteCalls())
func (mock *StoreMock) DeleteCalls() []struct {
	Ctx context.Context
	Key string
} {
	var calls []struct {
		Ctx context.Context
		Key string
	}
	mock.lockDelete.RLock()
	calls = mock.calls.Delete
	mock.lockDelete.RUnlock()
	return calls
}

// Load calls LoadFunc.
func (mock *StoreMock) Load(ctx context.Context, key string) (int, error) {
	if mock.LoadFunc == nil {
		panic("StoreMock.LoadFunc: method is nil but Store.Load was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Key string
	}{
		Ctx: ctx,
		Key: key,
	}
	mock.lockLoad.Lock()
	mock.calls.Load = append(mock.calls.Load, callInfo)
	mock.lockLoad.Unlock()
	return mock.LoadFunc(ctx, key)
}

// LoadCalls gets all the calls that were made to Load.
// Check the length with:
//
//	len(mockedStore.LoadCalls())
func (mock *StoreMock) LoadCalls() []struct {
	Ctx context.Context
	Key string
} {
	var calls []struct {
		Ctx context.Context
		Key string
	}
	mock.lockLoad.RLock()
	calls = mock.calls.Load
	mock.lockLoad.RUnlock()
	return calls
}

// Save calls SaveFunc.
func (mock *StoreMock) Save(ctx context.Context, key string, uid int) error {
	if mock.SaveFunc == nil {
		panic("StoreMock.SaveFunc: method is nil but Store.Save was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Key string
		UID int
	}{
		Ctx: ctx,
		Key: key,
		UID: uid,
	}
	mock.lockSave.Lock()
	mock.calls.Save = append(mock.calls.Save, callInfo)
	mock.lockSave.Unlock()
	return mock.SaveFunc(ctx, key, uid)
}

// SaveCalls gets all the calls that were made to Save.
// Check the length with:
//
//	len(mockedStore.SaveCalls())
func (mock *StoreMock) SaveCalls() []struct {
	Ctx context.Context
	Key string
	UID int
} {
	var calls []struct {
		Ctx context.Context
		Key string
		UID int
	}
	mock.lockSave.RLock()
	calls = mock.calls.Save
	mock.lockSave.RUnlock()
	return calls
}

// Ensure, that HealthRepositoryMock does implement HealthRepository.
// If this is not the case, regenerate this file with moq.
var _ HealthRepository = &HealthRepositoryMock{}

// HealthRepositoryMock is a mock implementation of HealthRepository.
//
//	func TestSomethingThatUsesHealthRepository(t *testing.T) {
//
//		// make and configure a mocked HealthRepository
//		mockedHealthRepository := &HealthRepositoryMock{
//			PingFunc: func(ctx context.Context, db store.DBConnection) error {
//				panic("mock out the Ping method")
//			},
//		}
//
//		// use mockedHealthRepository in code that requires HealthRepository
//		// and then make assertions.
//
//	}
type HealthRepositoryMock struct {
	// PingFunc mocks the Ping method.
	PingFunc func(ctx context.Context, db store.DBConnection) error

	// calls tracks calls to the methods.
	calls struct {
		// Ping holds details about calls to the Ping method.
		Ping []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Db is the db argument value.
			Db store.DBConnection
		}
	}
	lockPing sync.RWMutex
}

// Ping calls PingFunc.
func (mock *HealthRepositoryMock) Ping(ctx context.Context, db store.DBConnection) error {
	if mock.PingFunc == nil {
		panic("HealthRepositoryMock.PingFunc: method is nil but HealthRepository.Ping was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Db  store.DBConnection
	}{
		Ctx: ctx,
		Db:  db,
	}
	mock.lockPing.Lock()
	mock.calls.Ping = append(mock.calls.Ping, callInfo)
	mock.lockPing.Unlock()
	return mock.PingFunc(ctx, db)
}

// PingCalls gets all the calls that were made to Ping.
// Check the length with:
//
//	len(mockedHealthRepository.PingCalls())
func (mock *HealthRepositoryMock) PingCalls() []struct {
	Ctx context.Context
	Db  store.DBConnection
} {
	var calls []struct {
		Ctx context.Context
		Db  store.DBConnection
	}
	mock.lockPing.RLock()
	calls = mock.calls.Ping
	mock.lockPing.RUnlock()
	return calls
}

// Ensure, that AuthRepositoryMock does implement AuthRepository.
// If this is not the case, regenerate this file with moq.
var _ AuthRepository = &AuthRepositoryMock{}

// AuthRepositoryMock is a mock implementation of AuthRepository.
//
//	func TestSomethingThatUsesAuthRepository(t *testing.T) {
//
//		// make and configure a mocked AuthRepository
//		mockedAuthRepository := &AuthRepositoryMock{
//			GetUserByUserNameFunc: func(ctx context.Context, db store.DBConnection, userName string) (*model.User, error) {
//				panic("mock out the GetUserByUserName method")
//			},
//		}
//
//		// use mockedAuthRepository in code that requires AuthRepository
//		// and then make assertions.
//
//	}
type AuthRepositoryMock struct {
	// GetUserByUserNameFunc mocks the GetUserByUserName method.
	GetUserByUserNameFunc func(ctx context.Context, db store.DBConnection, userName string) (*model.User, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetUserByUserName holds details about calls to the GetUserByUserName method.
		GetUserByUserName []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Db is the db argument value.
			Db store.DBConnection
			// UserName is the userName argument value.
			UserName string
		}
	}
	lockGetUserByUserName sync.RWMutex
}

// GetUserByUserName calls GetUserByUserNameFunc.
func (mock *AuthRepositoryMock) GetUserByUserName(ctx context.Context, db store.DBConnection, userName string) (*model.User, error) {
	if mock.GetUserByUserNameFunc == nil {
		panic("AuthRepositoryMock.GetUserByUserNameFunc: method is nil but AuthRepository.GetUserByUserName was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Db       store.DBConnection
		UserName string
	}{
		Ctx:      ctx,
		Db:       db,
		UserName: userName,
	}
	mock.lockGetUserByUserName.Lock()
	mock.calls.GetUserByUserName = append(mock.calls.GetUserByUserName, callInfo)
	mock.lockGetUserByUserName.Unlock()
	return mock.GetUserByUserNameFunc(ctx, db, userName)
}

// GetUserByUserNameCalls gets all the calls that were made to GetUserByUserName.
// Check the length with:
//
//	len(mockedAuthRepository.GetUserByUserNameCalls())
func (mock *AuthRepositoryMock) GetUserByUserNameCalls() []struct {
	Ctx      context.Context
	Db       store.DBConnection
	UserName string
} {
	var calls []struct {
		Ctx      context.Context
		Db       store.DBConnection
		UserName string
	}
	mock.lockGetUserByUserName.RLock()
	calls = mock.calls.GetUserByUserName
	mock.lockGetUserByUserName.RUnlock()
	return calls
}

// Ensure, that UserRepositoryMock does implement UserRepository.
// If this is not the case, regenerate this file with moq.
var _ UserRepository = &UserRepositoryMock{}

// UserRepositoryMock is a mock implementation of UserRepository.
//
//	func TestSomethingThatUsesUserRepository(t *testing.T) {
//
//		// make and configure a mocked UserRepository
//		mockedUserRepository := &UserRepositoryMock{
//			GetAllUsersFunc: func(ctx context.Context, db store.DBConnection) ([]*model.User, error) {
//				panic("mock out the GetAllUsers method")
//			},
//			GetUserByUserNameFunc: func(ctx context.Context, db store.DBConnection, userName string) (*model.User, error) {
//				panic("mock out the GetUserByUserName method")
//			},
//			RegisterUserFunc: func(ctx context.Context, db store.DBConnection, u *model.User) (*model.User, error) {
//				panic("mock out the RegisterUser method")
//			},
//		}
//
//		// use mockedUserRepository in code that requires UserRepository
//		// and then make assertions.
//
//	}
type UserRepositoryMock struct {
	// GetAllUsersFunc mocks the GetAllUsers method.
	GetAllUsersFunc func(ctx context.Context, db store.DBConnection) ([]*model.User, error)

	// GetUserByUserNameFunc mocks the GetUserByUserName method.
	GetUserByUserNameFunc func(ctx context.Context, db store.DBConnection, userName string) (*model.User, error)

	// RegisterUserFunc mocks the RegisterUser method.
	RegisterUserFunc func(ctx context.Context, db store.DBConnection, u *model.User) (*model.User, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetAllUsers holds details about calls to the GetAllUsers method.
		GetAllUsers []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Db is the db argument value.
			Db store.DBConnection
		}
		// GetUserByUserName holds details about calls to the GetUserByUserName method.
		GetUserByUserName []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Db is the db argument value.
			Db store.DBConnection
			// UserName is the userName argument value.
			UserName string
		}
		// RegisterUser holds details about calls to the RegisterUser method.
		RegisterUser []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Db is the db argument value.
			Db store.DBConnection
			// U is the u argument value.
			U *model.User
		}
	}
	lockGetAllUsers       sync.RWMutex
	lockGetUserByUserName sync.RWMutex
	lockRegisterUser      sync.RWMutex
}

// GetAllUsers calls GetAllUsersFunc.
func (mock *UserRepositoryMock) GetAllUsers(ctx context.Context, db store.DBConnection) ([]*model.User, error) {
	if mock.GetAllUsersFunc == nil {
		panic("UserRepositoryMock.GetAllUsersFunc: method is nil but UserRepository.GetAllUsers was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Db  store.DBConnection
	}{
		Ctx: ctx,
		Db:  db,
	}
	mock.lockGetAllUsers.Lock()
	mock.calls.GetAllUsers = append(mock.calls.GetAllUsers, callInfo)
	mock.lockGetAllUsers.Unlock()
	return mock.GetAllUsersFunc(ctx, db)
}

// GetAllUsersCalls gets all the calls that were made to GetAllUsers.
// Check the length with:
//
//	len(mockedUserRepository.GetAllUsersCalls())
func (mock *UserRepositoryMock) GetAllUsersCalls() []struct {
	Ctx context.Context
	Db  store.DBConnection
} {
	var calls []struct {
		Ctx context.Context
		Db  store.DBConnection
	}
	mock.lockGetAllUsers.RLock()
	calls = mock.calls.GetAllUsers
	mock.lockGetAllUsers.RUnlock()
	return calls
}

// GetUserByUserName calls GetUserByUserNameFunc.
func (mock *UserRepositoryMock) GetUserByUserName(ctx context.Context, db store.DBConnection, userName string) (*model.User, error) {
	if mock.GetUserByUserNameFunc == nil {
		panic("UserRepositoryMock.GetUserByUserNameFunc: method is nil but UserRepository.GetUserByUserName was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Db       store.DBConnection
		UserName string
	}{
		Ctx:      ctx,
		Db:       db,
		UserName: userName,
	}
	mock.lockGetUserByUserName.Lock()
	mock.calls.GetUserByUserName = append(mock.calls.GetUserByUserName, callInfo)
	mock.lockGetUserByUserName.Unlock()
	return mock.GetUserByUserNameFunc(ctx, db, userName)
}

// GetUserByUserNameCalls gets all the calls that were made to GetUserByUserName.
// Check the length with:
//
//	len(mockedUserRepository.GetUserByUserNameCalls())
func (mock *UserRepositoryMock) GetUserByUserNameCalls() []struct {
	Ctx      context.Context
	Db       store.DBConnection
	UserName string
} {
	var calls []struct {
		Ctx      context.Context
		Db       store.DBConnection
		UserName string
	}
	mock.lockGetUserByUserName.RLock()
	calls = mock.calls.GetUserByUserName
	mock.lockGetUserByUserName.RUnlock()
	return calls
}

// RegisterUser calls RegisterUserFunc.
func (mock *UserRepositoryMock) RegisterUser(ctx context.Context, db store.DBConnection, u *model.User) (*model.User, error) {
	if mock.RegisterUserFunc == nil {
		panic("UserRepositoryMock.RegisterUserFunc: method is nil but UserRepository.RegisterUser was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Db  store.DBConnection
		U   *model.User
	}{
		Ctx: ctx,
		Db:  db,
		U:   u,
	}
	mock.lockRegisterUser.Lock()
	mock.calls.RegisterUser = append(mock.calls.RegisterUser, callInfo)
	mock.lockRegisterUser.Unlock()
	return mock.RegisterUserFunc(ctx, db, u)
}

// RegisterUserCalls gets all the calls that were made to RegisterUser.
// Check the length with:
//
//	len(mockedUserRepository.RegisterUserCalls())
func (mock *UserRepositoryMock) RegisterUserCalls() []struct {
	Ctx context.Context
	Db  store.DBConnection
	U   *model.User
} {
	var calls []struct {
		Ctx context.Context
		Db  store.DBConnection
		U   *model.User
	}
	mock.lockRegisterUser.RLock()
	calls = mock.calls.RegisterUser
	mock.lockRegisterUser.RUnlock()
	return calls
}
