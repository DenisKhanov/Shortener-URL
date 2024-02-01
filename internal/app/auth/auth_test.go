package auth

import (
	"github.com/google/uuid"
	"reflect"
	"testing"
)

type args struct {
	tokenString string
}

func TestBuildJWTString(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"ValidToken", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildJWTString()
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildJWTString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false && len(got) == 0 {
				t.Error("BuildJWTString() produced an empty token, but no error was expected")
			}
		})
	}
}

func TestGenerateUniqueID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"GeneratedID"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateUniqueID()
			if reflect.DeepEqual(got, uuid.Nil) {
				t.Error("GenerateUniqueID() returned the zero UUID, which is unexpected")
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ValidToken", args{tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDY4MTk4NzMsIlVzZXJJRCI6ImFkMTMwNjM1LTExNWMtNDhjMy1iYjZmLTJjNmZmNTIxNzA5ZSJ9.UjLVHSwQcLDZuUKIJoo1H3t8flbciC7eipnc-YcMW0w"}, false},
		{"InvalidToken", args{tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ.eyJleHAiOjE3MDY4MTk4NzMsIlVzZXJJRCI6ImFkMTMwNjM1LTExNWMtNDhjMy1iYjZmLTJjNmZmNTIxNzA5ZSJ9.UjLVHSwQcLDZuUKIJoo1H3t8flbciC7eipnc-YcMW0w"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserID(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false && reflect.DeepEqual(got, uuid.Nil) {
				t.Error("GetUserID() returned the zero UUID, which is unexpected")
			}
		})
	}
}

func TestIsValidToken(t *testing.T) {
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"ValidToken", args{tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDY4MTk4NzMsIlVzZXJJRCI6ImFkMTMwNjM1LTExNWMtNDhjMy1iYjZmLTJjNmZmNTIxNzA5ZSJ9.UjLVHSwQcLDZuUKIJoo1H3t8flbciC7eipnc-YcMW0w"}, true},
		{"InvalidToken", args{tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ.eyJleHAiOjE3MDY4MTk4NzMsIlVzZXJJRCI6ImFkMTMwNjM1LTExNWMtNDhjMy1iYjZmLTJjNmZmNTIxNzA5ZSJ9.UjLVHSwQcLDZuUKIJoo1H3t8flbciC7eipnc-YcMW0w"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidToken(tt.args.tokenString); got != tt.want {
				t.Errorf("IsValidToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
