// Package amulet menghitung skor jimat persis seperti kontrak Amulet milik
// Arachnid (at.amulet.garden): skor = panjang deret nibble hex '8' terpanjang
// pada sha256 teks jimat. Ambang mint di kontrak adalah 4.
package amulet

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// Rarity memetakan skor ke tingkat kelangkaan, sesuai komentar di kontrak
// Amulet.sol ("4: common ... 9: mythic, 10+: beyond mythic").
var Rarity = map[int]string{
	0: "bukan jimat", 1: "bukan jimat", 2: "bukan jimat", 3: "bukan jimat",
	4: "common", 5: "uncommon", 6: "rare", 7: "epic", 8: "legendary",
	9: "mythic", 10: "beyond mythic", 11: "beyond mythic", 12: "beyond mythic",
}

// MintThreshold adalah skor minimal yang diterima kontrak.
const MintThreshold = 4

// Score meniru Amulet.getScore. Kontrak memperlakukan sha256 sebagai uint256
// big-endian lalu `hash >>= 4`, artinya pembacaan mulai dari byte TERAKHIR,
// nibble rendah dulu. Contoh nyata: "snap" -> ...88889d42... => skor 4.
func Score(text string) int {
	sum := sha256.Sum256([]byte(text))
	maxRun, run := 0, 0
	for i := len(sum) - 1; i >= 0; i-- {
		b := sum[i]
		for shift := 0; shift < 8; shift += 4 {
			if (b>>shift)&0xF == 8 {
				run++
				if run > maxRun {
					maxRun = run
				}
			} else {
				run = 0
			}
		}
	}
	return maxRun
}

// Hash mengembalikan sha256 heksadesimal dari teks jimat.
func Hash(text string) string {
	sum := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", sum[:])
}

// RarityOf menamai skor; skor di atas 9 tetap 'beyond mythic'.
func RarityOf(score int) string {
	if score > 9 {
		return "beyond mythic"
	}
	if r, ok := Rarity[score]; ok {
		return r
	}
	return "bukan jimat"
}

// TokenID meniru skrip deploy: keccak256 tidak tersedia di pustaka standar,
// jadi id di sini adalah sha256 teks — cukup untuk penamaan lokal, dan
// dicatat berbeda dari id on-chain.
func TokenID(text string) string { return Hash(text) }

// Normalize merapikan jimat: buang spasi di ujung tiap baris, sisakan
// susunan baris apa adanya supaya sha256 tetap bisa direproduksi.
func Normalize(text string) string {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}
	return strings.Trim(strings.Join(lines, "\n"), "\n")
}
