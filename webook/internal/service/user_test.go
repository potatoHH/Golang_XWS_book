package service

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/internal/repository"
	"context"
	"reflect"
	"testing"
)

func TestUserService_Login(t *testing.T) {
	type fields struct {
		repo repository.UserRepository
	}
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &UserService{
				repo: tt.fields.repo,
			}
			got, err := svc.Login(tt.args.ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}
