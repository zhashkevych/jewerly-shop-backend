package postgres

import (
	"errors"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"testing"
)

func TestAdminRepository_Authorize(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type mockBehavior func(admin jewerly.AdminUser)

	type args struct {
		email    string
		password string
	}
	testTable := []struct {
		name         string
		admin        jewerly.AdminUser
		args         args
		mockBehavior mockBehavior
		shouldFail   bool
	}{
		{
			name:  "OK",
			admin: jewerly.AdminUser{Id: 14, Login: "test", PasswordHash: "qwerty"},
			args: args{
				email:    "test@test.com",
				password: "qwerty",
			},
			mockBehavior: func(admin jewerly.AdminUser) {
				rows := sqlmock.NewRows([]string{"id", "login", "password_hash"}).AddRow(admin.Id, admin.Login, admin.PasswordHash)
				mock.ExpectQuery("SELECT (.+) FROM admin_users").WillReturnRows(rows)
			},
		},
		{
			name:  "Empty Fields",
			admin: jewerly.AdminUser{Id: 14, Login: "test", PasswordHash: "qwerty"},
			args: args{
				email:    "",
				password: "qwerty",
			},
			mockBehavior: func(admin jewerly.AdminUser) {
				mock.ExpectQuery("SELECT (.+) FROM admin_users").WillReturnError(errors.New("no rows"))
			},
			shouldFail: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.admin)

			r := NewAdminRepository(db)

			err := r.Authorize(testCase.args.email, testCase.args.password)
			if testCase.shouldFail {
				assert.Error(t, err)
				t.Skip("OK")
			}

			assert.NoError(t, err)
		})
	}
}
