package email

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestEmail_EmailBytes(t *testing.T) {
	testTable := []struct {
		name       string
		email      Email
		shouldFail bool
	}{
		{
			name:  "OK",
			email: Email{ToEmail: "test@test.com", ToName: "Test", FromEmail: "zhashkevychmaksim@gmail.com", FromName: "Maksim", Body: "hey yo!", Subject: "HEY!"},
		},
		{
			name:       "Empty To email",
			email:      Email{ToEmail: "", ToName: "Test", FromEmail: "zhashkevychmaksim@gmail.com", FromName: "Maksim", Body: "hey yo!", Subject: "HEY!"},
			shouldFail: true,
		},
		{
			name:       "Invalid To email 1",
			email:      Email{ToEmail: "@gmail.com", ToName: "Test", FromEmail: "zhashkevychmaksim@gmail.com", FromName: "Maksim", Body: "hey yo!", Subject: "HEY!"},
			shouldFail: true,
		},
		{
			name:       "Invalid To email 2",
			email:      Email{ToEmail: "qwe@.com", ToName: "Test", FromEmail: "zhashkevychmaksim@gmail.com", FromName: "Maksim", Body: "hey yo!", Subject: "HEY!"},
			shouldFail: true,
		},
		{
			name:       "Invalid To email 3",
			email:      Email{ToEmail: "qwe@gmail.", ToName: "Test", FromEmail: "zhashkevychmaksim@gmail.com", FromName: "Maksim", Body: "hey yo!", Subject: "HEY!"},
			shouldFail: true,
		},
		{
			name:       "Empty To Name",
			email:      Email{ToEmail: "qwe@gmail.com", ToName: "", FromEmail: "zhashkevychmaksim@gmail.com", FromName: "Maksim", Body: "hey yo!", Subject: "HEY!"},
			shouldFail: true,
		},
		{
			name:       "Empty From Email",
			email:      Email{ToEmail: "qwe@gmail.com", ToName: "Test", FromEmail: "", FromName: "Maksim", Body: "hey yo!", Subject: "HEY!"},
			shouldFail: true,
		},
		{
			name:       "Empty From Name",
			email:      Email{ToEmail: "qwe@gmail.com", ToName: "Test", FromEmail: "zhashkevychmaksim@gmail.com", FromName: "", Body: "hey yo!", Subject: "HEY!"},
			shouldFail: true,
		},
		{
			name:       "Empty Body",
			email:      Email{ToEmail: "qwe@gmail.com", ToName: "Test", FromEmail: "zhashkevychmaksim@gmail.com", FromName: "Maksim", Body: "", Subject: "HEY!"},
			shouldFail: true,
		},
		{
			name:       "Empty Subject",
			email:      Email{ToEmail: "qwe@gmail.com", ToName: "Test", FromEmail: "zhashkevychmaksim@gmail.com", FromName: "Maksim", Body: "", Subject: ""},
			shouldFail: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := testCase.email.EmailBytes()
			if testCase.shouldFail {
				assert.Error(t, err)
				t.Skip("ok")
			}

			assert.NoError(t, err)

			if !strings.Contains(string(got), `Content-Type: text/html; charset="UTF-8"`) {
				t.Fatal("Content-Type header missing")
			}
			if !strings.Contains(string(got), fmt.Sprintf(`To: "%s" <%s>`, testCase.email.ToName, testCase.email.ToEmail)) {
				t.Fatal("To header missing")
			}
			if !strings.Contains(string(got), fmt.Sprintf(`From: "%s" <%s>`, testCase.email.FromName, testCase.email.FromEmail)) {
				t.Fatal("From header missing")
			}
			if !strings.Contains(string(got), fmt.Sprintf(`Subject: %s`, testCase.email.Subject)) {
				t.Fatal("Subject header missing")
			}
			if !strings.Contains(string(got), fmt.Sprintf("\r\n%s", testCase.email.Body)) {
				t.Fatal("Body missing")
			}
		})
	}
}
