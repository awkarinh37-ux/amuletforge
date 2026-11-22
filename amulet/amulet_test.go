package amulet

import "testing"

// Teks nyata dari daftar PRELOAD_MINTS di repo Arachnid/amulets.
func TestScore_knownAmulets(t *testing.T) {
	cases := []struct {
		text string
		want int
	}{
		{"Amul888888888888888888888888888", 4},
		{"If you can't write poems,\nwrite me", 5}, // deret '8' panjang
		{"at.amulet.garden", 4},
		{"snap", 4},
		{"= 69", 5},
	}
	for _, c := range cases {
		if got := Score(c.text); got != c.want {
			t.Errorf("Score(%q) = %d, mau %d", c.text, got, c.want)
		}
	}
}

// Deret '8' harus dihitung berurutan nibble, termasuk lintas byte.
func TestScore_runsAcrossBytes(t *testing.T) {
	// cari teks dengan skor >= 4 sungguhan, bukan buatan
	found := 0
	for i := 0; i < 400000 && found == 0; i++ {
		s := "ujian-" + itoa(i)
		if Score(s) >= 4 {
			found = Score(s)
		}
	}
	if found == 0 {
		t.Fatal("tidak menemukan jimat skor >=4 dalam 400k percobaan")
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b []byte
	for i > 0 {
		b = append([]byte{byte('0' + i%10)}, b...)
		i /= 10
	}
	return string(b)
}

func TestRarity_naming(t *testing.T) {
	if RarityOf(4) != "common" {
		t.Error("skor 4 harus common")
	}
	if RarityOf(9) != "mythic" {
		t.Error("skor 9 harus mythic")
	}
	if RarityOf(12) != "beyond mythic" {
		t.Error("skor di atas 9 harus beyond mythic")
	}
	if RarityOf(3) != "bukan jimat" {
		t.Error("skor 3 bukan jimat")
	}
}

func TestNormalize_trimsLineEnds(t *testing.T) {
	got := Normalize("  baris satu   \nbaris dua\t\n\n")
	want := "  baris satu\nbaris dua"
	if got != want {
		t.Errorf("Normalize = %q, mau %q", got, want)
	}
}

func TestHash_stable(t *testing.T) {
	if Hash("snap") != Hash("snap") {
		t.Error("hash harus stabil")
	}
	if Hash("snap") == Hash("snap ") {
		t.Error("spasi akhir harus mengubah hash")
	}
}
