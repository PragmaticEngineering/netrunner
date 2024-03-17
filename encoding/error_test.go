package encoding

import (
	"errors"
	"testing"
)

func TestEncodingError_Error(t *testing.T) {
	tests := []struct {
		name           string
		encodingError  *EncodingError
		expectedOutput string
	}{
		{
			name: "Error with additional info",
			encodingError: &EncodingError{
				err:            errors.New("some error"),
				sentinelError:  WriteResponse,
				additionalInfo: []string{"info1", "info2"},
			},
			expectedOutput: "write response error error: some error. AdditionalInformation: [info1 info2]",
		},
		{
			name: "Error without additional info",
			encodingError: &EncodingError{
				err:           errors.New("another error"),
				sentinelError: EncodeJSON,
			},
			expectedOutput: "encode json error error: another error.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := test.encodingError.Error()
			if output != test.expectedOutput {
				t.Errorf("Expected error message '%s', but got '%s'", test.expectedOutput, output)
			}
		})
	}
}

func TestEncodingError_Unwrap(t *testing.T) {
	err := errors.New("wrapped error")
	encodingError := newEncodingError(DecodeJSON, err)

	unwrapped := encodingError.Unwrap()
	if unwrapped != err {
		t.Errorf("Expected unwrapped error to be '%v', but got '%v'", err, unwrapped)
	}
}
