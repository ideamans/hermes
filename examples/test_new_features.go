package main

import "github.com/ideamans/hermes"

type testNewFeatures struct{}

func (t *testNewFeatures) Name() string {
	return "test_new_features"
}

func (t *testNewFeatures) Email() hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: "Test User",
			Intros: []string{
				"This email demonstrates the new features:",
				"- Icon displayed to the left of the product name",
				"- Custom logo width (200px)",
				"- Greeting is hidden (HideGreeting = true)",
				"- Signature is hidden (HideSignature = true)",
			},
			Outros: []string{
				"Notice that neither the greeting nor the signature line is displayed.",
			},
		},
	}
}
