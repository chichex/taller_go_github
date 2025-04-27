package user

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_Create_Simple(t *testing.T) {
	s := NewService(NewLocalStorage(), nil)

	input := &User{
		Name:     "Ayrton",
		Address:  "Pringles",
		NickName: "Chiche",
	}

	err := s.Create(input)

	require.Nil(t, err)
	require.NotEmpty(t, input.ID)
	require.NotEmpty(t, input.CreatedAt)
	require.NotEmpty(t, input.UpdatedAt)
	require.Equal(t, 1, input.Version)

	s = NewService(&mockStorage{
		mockSet: func(user *User) error {
			return errors.New("fake error trying to set user")
		},
	}, nil)

	err = s.Create(input)
	require.NotNil(t, err)
	require.EqualError(t, err, "fake error trying to set user")
}

func TestService_Create(t *testing.T) {
	type fields struct {
		storage Storage
	}

	type args struct {
		user *User
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  func(t *testing.T, err error)
		wantUser func(t *testing.T, user *User)
	}{
		{
			name: "error",
			fields: fields{
				storage: &mockStorage{
					mockSet: func(user *User) error {
						return errors.New("fake error trying to set user")
					},
				},
			},
			args: args{
				user: &User{},
			},
			wantErr: func(t *testing.T, err error) {
				require.NotNil(t, err)
				require.EqualError(t, err, "fake error trying to set user")
			},
			wantUser: nil,
		},
		{
			name: "success",
			fields: fields{
				storage: NewLocalStorage(),
			},
			args: args{
				user: &User{
					Name:     "Ayrton",
					Address:  "Pringles",
					NickName: "Chiche",
				},
			},
			wantErr: func(t *testing.T, err error) {
				require.Nil(t, err)
			},
			wantUser: func(t *testing.T, input *User) {
				require.NotEmpty(t, input.ID)
				require.NotEmpty(t, input.CreatedAt)
				require.NotEmpty(t, input.UpdatedAt)
				require.Equal(t, 1, input.Version)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
			}

			err := s.Create(tt.args.user)
			if tt.wantErr != nil {
				tt.wantErr(t, err)
			}

			if tt.wantUser != nil {
				tt.wantUser(t, tt.args.user)
			}
		})
	}
}

type mockStorage struct {
	mockSet    func(user *User) error
	mockRead   func(id string) (*User, error)
	mockDelete func(id string) error
}

func (m *mockStorage) Set(user *User) error {
	return m.mockSet(user)
}

func (m *mockStorage) Read(id string) (*User, error) {
	return m.mockRead(id)
}

func (m *mockStorage) Delete(id string) error {
	return m.mockDelete(id)
}
