package mocks

import (
	"context"
	"strings"

	protomfx "github.com/MainfluxLabs/mainflux/pkg/proto"
	"github.com/MainfluxLabs/mainflux/users"
	"google.golang.org/grpc"
)

var _ protomfx.UsersServiceClient = (*usersServiceClientMock)(nil)

type usersServiceClientMock struct {
	usersByID     map[string]users.User
	usersByEmails map[string]users.User
}

func NewUsersService(usersByID map[string]users.User, usersByEmails map[string]users.User) protomfx.UsersServiceClient {
	return &usersServiceClientMock{usersByID: usersByID, usersByEmails: usersByEmails}
}

func (svc *usersServiceClientMock) GetUsersByIDs(_ context.Context, req *protomfx.UsersByIDsReq, _ ...grpc.CallOption) (*protomfx.UsersRes, error) {
	var users []*protomfx.User
	for _, id := range req.Ids {
		if user, ok := svc.usersByID[id]; ok {
			if req.Email != "" && !strings.Contains(user.Email, req.Email) {
				continue
			}
			users = append(users, &protomfx.User{Id: user.ID, Email: user.Email})
		}
	}

	return &protomfx.UsersRes{Users: users}, nil
}

func (svc *usersServiceClientMock) GetUsersByEmails(_ context.Context, req *protomfx.UsersByEmailsReq, _ ...grpc.CallOption) (*protomfx.UsersRes, error) {
	var users []*protomfx.User
	for _, email := range req.Emails {
		if user, ok := svc.usersByEmails[email]; ok {
			users = append(users, &protomfx.User{Id: user.ID, Email: user.Email})
		}
	}

	return &protomfx.UsersRes{Users: users}, nil
}
