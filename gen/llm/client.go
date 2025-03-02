package llm

type Client interface {
	Call(prompt string) (string, error)
	CallAndCheck(prompt string, check func(string) error) (string, error)
}

type MockClient struct {
}

func (client *MockClient) Call(prompt string) (string, error) {
	return `
		func Fuzz_Mock(f *testing.F) {
			f.Fuzz(func(t *testing.T, data []byte) {
				var s *jwt.ClaimStrings
				var data_0 []byte
				fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes, t, fuzzer.Constructors)
				err := fz.Fill2(&s, &data_0)
				if err != nil || s == nil {
					return
				}

				s.UnmarshalJSON(data_0)
			})
		}
	`, nil
}

func (client *MockClient) CallAndCheck(prompt string, check func(string) error) (string, error) {
	result, err := client.Call(prompt)
	if err != nil {
		return "", nil
	}

	if err = check(result); err != nil {
		return "", nil
	}

	return result, nil
}
