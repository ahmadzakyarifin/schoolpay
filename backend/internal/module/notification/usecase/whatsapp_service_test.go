package usecase

import (
	"reflect"
	"testing"
)

func TestNewWhatsAppService(t *testing.T) {
	tests := []struct {
		name string
		want WhatsAppService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWhatsAppService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWhatsAppService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_whatsappService_GetStatus(t *testing.T) {
	tests := []struct {
		name    string
		s       *whatsappService
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &whatsappService{}
			got, err := s.GetStatus()
			if (err != nil) != tt.wantErr {
				t.Fatalf("whatsappService.GetStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("whatsappService.GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_whatsappService_StartSession(t *testing.T) {
	tests := []struct {
		name    string
		s       *whatsappService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &whatsappService{}
			if err := s.StartSession(); (err != nil) != tt.wantErr {
				t.Errorf("whatsappService.StartSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_whatsappService_StopSession(t *testing.T) {
	tests := []struct {
		name    string
		s       *whatsappService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &whatsappService{}
			if err := s.StopSession(); (err != nil) != tt.wantErr {
				t.Errorf("whatsappService.StopSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_whatsappService_GetQR(t *testing.T) {
	tests := []struct {
		name    string
		s       *whatsappService
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &whatsappService{}
			got, err := s.GetQR()
			if (err != nil) != tt.wantErr {
				t.Fatalf("whatsappService.GetQR() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("whatsappService.GetQR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_whatsappService_RegisterWebhook(t *testing.T) {
	tests := []struct {
		name    string
		s       *whatsappService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &whatsappService{}
			if err := s.RegisterWebhook(); (err != nil) != tt.wantErr {
				t.Errorf("whatsappService.RegisterWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_whatsappService_GetChatHistory(t *testing.T) {
	type args struct {
		phone string
	}
	tests := []struct {
		name    string
		s       *whatsappService
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &whatsappService{}
			got, err := s.GetChatHistory(tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Fatalf("whatsappService.GetChatHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("whatsappService.GetChatHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_whatsappService_SendChatMessage(t *testing.T) {
	type args struct {
		phone   string
		message string
	}
	tests := []struct {
		name    string
		s       *whatsappService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &whatsappService{}
			if err := s.SendChatMessage(tt.args.phone, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("whatsappService.SendChatMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
