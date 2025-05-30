package segments

import (
	"testing"

	"github.com/jandedobbeleer/oh-my-posh/src/properties"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime/mock"

	"github.com/stretchr/testify/assert"
)

func TestAWSSegment(t *testing.T) {
	cases := []struct {
		Case            string
		ExpectedString  string
		Profile         string
		DefaultProfile  string
		Vault           string
		Region          string
		DefaultRegion   string
		ConfigFile      string
		Template        string
		ExpectedEnabled bool
		DisplayDefault  bool
	}{
		{Case: "enabled with default user", ExpectedString: "default@eu-west", Region: "eu-west", ExpectedEnabled: true, DisplayDefault: true},
		{Case: "disabled with default user", ExpectedString: "default@eu-west", Region: "eu-west", ExpectedEnabled: false, DisplayDefault: false},
		{Case: "disabled", ExpectedString: "", ExpectedEnabled: false},
		{Case: "enabled with default user", ExpectedString: "default@eu-west", Profile: "default", Region: "eu-west", ExpectedEnabled: true, DisplayDefault: true},
		{Case: "enabled with default profile", ExpectedString: "default@eu-west", DefaultProfile: "default", Region: "eu-west", ExpectedEnabled: true, DisplayDefault: true},
		{Case: "disabled with default user", ExpectedString: "default", Profile: "default", Region: "eu-west", ExpectedEnabled: false, DisplayDefault: false},
		{Case: "enabled no region", ExpectedString: "company", ExpectedEnabled: true, Profile: "company"},
		{Case: "enabled with region", ExpectedString: "company@eu-west", ExpectedEnabled: true, Profile: "company", Region: "eu-west", DefaultRegion: "us-west"},
		{Case: "enabled with default region", ExpectedString: "company@us-west", ExpectedEnabled: true, Profile: "company", DefaultRegion: "us-west"},
		{
			Case:            "template: enabled no region",
			ExpectedString:  "profile: company",
			ExpectedEnabled: true,
			Profile:         "company",
			Template:        "profile: {{.Profile}}{{if .Region}} in {{.Region}}{{end}}",
		},
		{
			Case:            "template: enabled with region",
			ExpectedString:  "profile: company in eu-west",
			ExpectedEnabled: true,
			Profile:         "company",
			Region:          "eu-west",
			Template:        "profile: {{.Profile}}{{if .Region}} in {{.Region}}{{end}}",
		},
		{
			Case:            "template: enabled with region alias that has compound cardinal direction",
			ExpectedString:  "profile: company in apne3",
			ExpectedEnabled: true,
			Profile:         "company",
			Region:          "ap-northeast-3",
			Template:        "profile: {{.Profile}}{{if .Region}} in {{.RegionAlias}}{{end}}",
		},
		{Case: "template: invalid", ExpectedString: "{{ .Burp", ExpectedEnabled: true, Profile: "c", Template: "{{ .Burp"},
	}

	for _, tc := range cases {
		env := new(mock.Environment)
		env.On("Getenv", "AWS_VAULT").Return(tc.Vault)
		env.On("Getenv", "AWS_PROFILE").Return(tc.Profile)
		env.On("Getenv", "AWS_DEFAULT_PROFILE").Return(tc.DefaultProfile)
		env.On("Getenv", "AWS_REGION").Return(tc.Region)
		env.On("Getenv", "AWS_DEFAULT_REGION").Return(tc.DefaultRegion)
		env.On("Getenv", "AWS_CONFIG_FILE").Return(tc.ConfigFile)
		env.On("FileContent", "/usr/home/.aws/config").Return("")
		env.On("Home").Return("/usr/home")
		props := properties.Map{
			properties.DisplayDefault: tc.DisplayDefault,
		}
		env.On("Flags").Return(&runtime.Flags{})

		aws := &Aws{}
		aws.Init(props, env)

		if tc.Template == "" {
			tc.Template = aws.Template()
		}

		assert.Equal(t, tc.ExpectedEnabled, aws.Enabled(), tc.Case)
		assert.Equal(t, tc.ExpectedString, renderTemplate(env, tc.Template, aws), tc.Case)
	}
}
