package app

import "testing"

func TestValidateUsernameSyntax(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		wantError     bool
		expectedError string
	}{
		{
			name:          "username valid",
			username:      "user_123",
			wantError:     false,
			expectedError: "",
		},
		{
			name:          "username too short",
			username:      "ab",
			wantError:     true,
			expectedError: ErrorUsernameTooShort.Error(),
		},
		{
			name:          "username too long",
			username:      "areallyreallylongusername",
			wantError:     true,
			expectedError: ErrorUsernameTooLong.Error(),
		},
		{
			name:          "username contains invalid character",
			username:      "invalid username",
			wantError:     true,
			expectedError: ErrorUsernameInvalidChars.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUsernameSyntax(tt.username)

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

func TestValidateEmailSyntax(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		wantError     bool
		expectedError string
	}{
		{
			name:          "email valid",
			email:         "user@example.com",
			wantError:     false,
			expectedError: "",
		},
		{
			name:          "email invalid",
			email:         "not-email",
			wantError:     true,
			expectedError: ErrorEmailInvalid.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmailSyntax(tt.email)

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
