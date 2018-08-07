package tomlssm

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/BurntSushi/toml"
)

// decode... is Test function for func Decode()
// because use SSMmock.
func decode(in string, v interface{}) (toml.MetaData, error) {
	out, err := toml.Decode(in, v)
	if err != nil {
		return toml.MetaData{}, err
	}

	mock := &mockSSMClient{}
	d := newTestssmDecrypter(mock)

	// override is a method in tomlssm package.
	d.override(v)
	return out, nil
}

// mockSSMClient... stores SSM interface for mock
type mockSSMClient struct {
	ssmiface.SSMAPI
}

// newTestssmDecrypter... returns a new ssmDecrypter for mock.
func newTestssmDecrypter(mock ssmiface.SSMAPI) *ssmDecrypter {
	return &ssmDecrypter{
		svc: mock,
	}
}

// GetParameter... returns "decrypted" that is Decrypted SSM parameter.
func (m *mockSSMClient) GetParameter(i *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	parameter := &ssm.Parameter{
		Value: aws.String("decrypted"),
	}

	return &ssm.GetParameterOutput{
		Parameter: parameter,
	}, nil
}

// TestSSMDecode... testing SSM decode
func TestSSMDecode(t *testing.T) {
	cases := []struct{
		value string
		expected interface{}
	}{
		{
			"value = \"a\"",
			"a",

		},
		{
			"value = \"ssm://encrypt_parameter\"",
			"decrypted",

		},
	}

	type Data struct{
		Value string
	}
	for _, c := range cases {
		var d Data
		if _, err := decode(c.value, &d); err != nil {
			t.Fatalf("failed unmarshal: %s", err)
		}
		if d.Value != c.expected {
			t.Errorf("want %s got %s", c.expected, d.Value)
		}
	}
}