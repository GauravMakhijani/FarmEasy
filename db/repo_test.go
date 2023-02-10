package db

import (
	"FarmEasy/domain"
	"context"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

type DbTestSuite struct {
	suite.Suite
	repo Storer

	mock sqlxmock.Sqlmock
}

func TestDbTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}

func (suite *DbTestSuite) SetupTest() {
	var err error
	var db *sqlx.DB
	db, suite.mock, err = sqlxmock.Newx()
	suite.Require().NoError(err)
	suite.repo = NewPgStore(db)
}

func (suite *DbTestSuite) TearDownSuite() {

	suite.mock.ExpectClose()
}

func (s *DbTestSuite) Test_pgStore_RegisterFarmer() {
	t := s.T()
	type args struct {
		ctx    context.Context
		farmer *domain.FarmerResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args, sqlxmock.Sqlmock)
	}{
		// TODO: Add test cases.
		{
			name: "positiveTest",
			args: args{
				ctx: context.TODO(),
				farmer: &domain.FarmerResponse{
					Id:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@gmail.com",
					Phone:     "1234567890",
					Address:   "1234, abc street, xyz city",
					Password:  "password",
				},
			},
			wantErr: false,
		},
		{
			name: "negativeTest",
			args: args{
				ctx: context.TODO(),
				farmer: &domain.FarmerResponse{
					Id:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@gmail.com",
					Phone:     "1234567890",
					Address:   "1234, abc street, xyz city",
					Password:  "password",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			if tt.wantErr {
				err = errors.New("mocked error")
			} else {
				err = nil
			}
			rows := sqlxmock.NewRows([]string{"id"}).AddRow(1)

			s.mock.ExpectQuery("INSERT INTO farmers").WithArgs(tt.args.farmer.FirstName, tt.args.farmer.LastName, tt.args.farmer.Email, tt.args.farmer.Phone, tt.args.farmer.Address, tt.args.farmer.Password).WillReturnError(err).WillReturnRows(rows)

			if err := s.repo.RegisterFarmer(tt.args.ctx, tt.args.farmer); tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
