package main

import (
	"testing"
)

func TestValidateJSONFromFile_Table(t *testing.T) {
	tests := []struct {
		name string
		file string
		want bool
	}{
		{"Step1Invalid", "../test_files/step1/invalid.json", false},
		{"Step1Invalid2", "../test_files/step1/invalid2.json", false},
		{"Step1Valid", "../test_files/step1/valid.json", true},
		{"Step2Invalid", "../test_files/step2/invalid.json", false},
		{"Step2Invalid2", "../test_files/step2/invalid2.json", false},
		{"Step2Valid", "../test_files/step2/valid.json", true},
		{"Step2Valid2", "../test_files/step2/valid2.json", true},
		{"Step3Invalid", "../test_files/step3/invalid.json", false},
		{"Step3Valid", "../test_files/step3/valid.json", true},
		{"Step4Invalid", "../test_files/step4/invalid.json", false},
		{"Step4Valid", "../test_files/step4/valid.json", true},
		{"Step4Valid2", "../test_files/step4/valid2.json", true},
		{"Step5Valid1", "../test_files/step5/valid.json", true},
		{"Step5Valid2", "../test_files/step5/valid2.json", true},
		{"Step5Valid3", "../test_files/step5/valid3.json", true},
		{"Step5Invalid1", "../test_files/step5/fail1.json", false},
		{"Step5Invalid1", "../test_files/step5/fail1.json", false},
		{"Step5Invalid2", "../test_files/step5/fail2.json", false},
		{"Step5Invalid3", "../test_files/step5/fail3.json", false},
		{"Step5Invalid4", "../test_files/step5/fail4.json", false},
		{"Step5Invalid5", "../test_files/step5/fail5.json", false},
		{"Step5Invalid6", "../test_files/step5/fail6.json", false},
		{"Step5Invalid7", "../test_files/step5/fail7.json", false},
		{"Step5Invalid8", "../test_files/step5/fail8.json", false},
		{"Step5Invalid9", "../test_files/step5/fail9.json", false},
		{"Step5Invalid10", "../test_files/step5/fail10.json", false},
		{"Step5Invalid11", "../test_files/step5/fail11.json", false},
		{"Step5Invalid12", "../test_files/step5/fail12.json", false},
		{"Step5Invalid13", "../test_files/step5/fail13.json", false},
		{"Step5Invalid14", "../test_files/step5/fail14.json", false},
		{"Step5Invalid15", "../test_files/step5/fail15.json", false},
		{"Step5Invalid16", "../test_files/step5/fail16.json", false},
		{"Step5Invalid17", "../test_files/step5/fail17.json", false},
		{"Step5Invalid19", "../test_files/step5/fail19.json", false},
		{"Step5Invalid20", "../test_files/step5/fail20.json", false},
		{"Step5Invalid21", "../test_files/step5/fail21.json", false},
		{"Step5Invalid22", "../test_files/step5/fail22.json", false},
		{"Step5Invalid23", "../test_files/step5/fail23.json", false},
		{"Step5Invalid24", "../test_files/step5/fail24.json", false},
		{"Step5Invalid25", "../test_files/step5/fail25.json", false},
		{"Step5Invalid26", "../test_files/step5/fail26.json", false},
		{"Step5Invalid27", "../test_files/step5/fail27.json", false},
		{"Step5Invalid28", "../test_files/step5/fail28.json", false},
		{"Step5Invalid29", "../test_files/step5/fail29.json", false},
		{"Step5Invalid30", "../test_files/step5/fail30.json", false},
		{"Step5Invalid31", "../test_files/step5/fail31.json", false},
		{"Step5Invalid32", "../test_files/step5/fail32.json", false},
		{"Step5Invalid33", "../test_files/step5/fail33.json", false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			isValid := validateJSONFromFile(tc.file)
			if isValid != tc.want {
				t.Fatalf("unexpected output.\nexpected: %v\nactual: %v", tc.want, isValid)
			}
		})
	}
}
