package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantErr  bool
	}{
		{
			name:  "valid REGISTER message",
			input: `{"type":"REGISTER","role":"developer","description":"Revolutionary Go developer"}`,
			expected: &RegisterMessage{
				Type:        MsgTypeRegister,
				Role:        "developer",
				Description: "Revolutionary Go developer",
			},
			wantErr: false,
		},
		{
			name:  "valid YIELD message",
			input: `{"type":"YIELD","to_role":"tester","payload":"Code ready for testing"}`,
			expected: &YieldMessage{
				Type:    MsgTypeYield,
				ToRole:  "tester",
				Payload: "Code ready for testing",
			},
			wantErr: false,
		},
		{
			name:  "valid ACTIVATE message",
			input: `{"type":"ACTIVATE","from_role":"developer","payload":"Begin testing phase"}`,
			expected: &ActivateMessage{
				Type:     MsgTypeActivate,
				FromRole: "developer",
				Payload:  "Begin testing phase",
			},
			wantErr: false,
		},
		{
			name:  "valid ERROR message",
			input: `{"type":"ERROR","message":"Invalid barrel holder"}`,
			expected: &ErrorMessage{
				Type:    MsgTypeError,
				Message: "Invalid barrel holder",
			},
			wantErr: false,
		},
		{
			name:  "valid ACK_REGISTER message",
			input: `{"type":"ACK_REGISTER"}`,
			expected: &AckRegisterMessage{
				Type: MsgTypeAckRegister,
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{"type":"REGISTER","role":}`,
			wantErr: true,
		},
		{
			name:    "unknown message type",
			input:   `{"type":"UNKNOWN"}`,
			wantErr: true,
		},
		{
			name:    "missing type field",
			input:   `{"role":"developer"}`,
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   ``,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseMessage([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err, "ParseMessage() should return an error")
				return
			}

			assert.NoError(t, err, "ParseMessage() should not return an error")
			assert.Equal(t, tt.expected, result, "ParseMessage() result should match expected")
		})
	}
}
