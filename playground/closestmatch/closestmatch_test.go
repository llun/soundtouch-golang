package closestmatch

import (
	"testing"

	"github.com/schollz/closestmatch"
	levenshtein "github.com/schollz/closestmatch/levenshtein"
)

func TestAccuracyBookletters(t *testing.T) {
	wordsToTest := []string{"Office", "Küche", "Schrank", "Schlafzimmer", "Badezimmer", "Prinzessinen", "Wohnzimmer"}
	lv := levenshtein.New(wordsToTest)
	bagSizes := []int{2}
	cm := closestmatch.New(wordsToTest, bagSizes)
	tests := []struct {
		key  string
		want string
	}{
		{"Office", "Office"},
		{"office", "Office"},
		{"OFFICE", "Office"},
		{"Kueche", "Küche"},
		{"kueche", "Küche"},
		{"Küch", "Küche"},
		{"Kuche", "Küche"},
		{"Badezimme", "Badezimmer"},
		{"Bad", "Badezimmer"},
		{"Bade", "Badezimmer"},
		{"Badezi", "Badezimmer"},
		{"Prinz", "Prinzessinen"},
		{"Princesses", "Prinzessinen"},
		{"Shrank", "Schrank"},
		{"schrank", "Schrank"},
		{"Schrang", "Schrank"},
	}
	var lvFails int
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := lv.Closest(tt.key); got != tt.want {
				t.Errorf("closest(%v) = %v, want %v", lv.Accuracy(), got, tt.want)
				lvFails++
			} else {
				t.Logf("closest(%v) matched", lv.Accuracy())
			}
		})
	}

	var cmFails int
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := cm.Closest(tt.key); got != tt.want {
				t.Errorf("closest() = %v, want %v", got, tt.want)
				cmFails++
			}
		})
	}

	t.Errorf(" lvFails= %v vs cmFails= %v", lvFails, cmFails)

}
