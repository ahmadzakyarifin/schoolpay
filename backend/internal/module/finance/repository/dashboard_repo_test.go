package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/uptrace/bun"
)

func TestNewDashboardRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want DashboardRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDashboardRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDashboardRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dashboardRepo_CountNewUsersByPeriod(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx   context.Context
		start *time.Time
		end   *time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &dashboardRepo{
				db: tt.fields.db,
			}
			got, err := r.CountNewUsersByPeriod(tt.args.ctx, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Fatalf("dashboardRepo.CountNewUsersByPeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("dashboardRepo.CountNewUsersByPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}
