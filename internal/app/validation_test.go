package app

import "testing"

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		wantError     bool
		expectedError string
	}{
		{
			name:          "password valid",
			password:      "ValidPass12!",
			wantError:     false,
			expectedError: "",
		},
		{
			name:          "password too short",
			password:      "Vp1!",
			wantError:     true,
			expectedError: "Password must be at least 12 characters long",
		},
		{
			name:          "password missing uppercase",
			password:      "password123!",
			wantError:     true,
			expectedError: "Password must contain at least 1 uppercase letter",
		},
		{
			name:          "password missing lowercase",
			password:      "PASSWORD123!",
			wantError:     true,
			expectedError: "Password must contain at least 1 lowercase letter",
		},
		{
			name:          "password missing number",
			password:      "PasswordABC!",
			wantError:     true,
			expectedError: "Password must contain at least 1 number",
		},
		{
			name:          "password missing special character",
			password:      "Password123ABC",
			wantError:     true,
			expectedError: "Password must contain at least 1 special character",
		},
		{
			name:          "password empty",
			password:      "",
			wantError:     true,
			expectedError: "Password must be at least 12 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := valdiatePassword(tt.password)

			if tt.wantError {
				if err == nil {
					t.Fatal("expected an error; got 'nil'")
				}
				if err.Error() != tt.expectedError {
					t.Errorf("expected error '%s'; got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("did not expect an error; got '%v'", err)
				}
			}
		})
	}
}
