package utils

import "testing"

func TestIsMimoASRExt(t *testing.T) {
	for _, name := range []string{"answer.wav", "ANSWER.MP3"} {
		if !IsMimoASRExt(name) {
			t.Fatalf("IsMimoASRExt(%q) = false", name)
		}
	}
	for _, name := range []string{"answer.webm", "answer.m4a", "answer.ogg", "answer"} {
		if IsMimoASRExt(name) {
			t.Fatalf("IsMimoASRExt(%q) = true", name)
		}
	}
}
