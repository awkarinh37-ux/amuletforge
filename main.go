// Command amuletforge menambang jimat: variasi teks sampai skor sha256 cukup
// langka. Skor dihitung sama persis dengan kontrak Amulet (sudah dicocokkan
// terhadap getScore on-chain di 0x2A21...423d).
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"amuletforge/amulet"
)

func main() {
	base := flag.String("text", "", "teks dasar jimat (baris baru dipisah \\n)")
	target := flag.Int("min", amulet.MintThreshold, "skor minimal yang dicari (kontrak menolak di bawah 4)")
	limit := flag.Int("n", 1, "berapa hasil yang dikumpulkan")
	workers := flag.Int("w", 0, "jumlah pekerja (0 = jumlah CPU)")
	check := flag.String("check", "", "hanya cetak skor teks ini lalu keluar")
	maxSec := flag.Int("t", 30, "batas waktu detik")
	flag.Parse()

	if *check != "" {
		s := strings.ReplaceAll(*check, `\n`, "\n")
		fmt.Printf("skor     : %d (%s)\n", amulet.Score(s), amulet.RarityOf(amulet.Score(s)))
		fmt.Printf("sha256   : %s\n", amulet.Hash(s))
		if amulet.Score(s) < amulet.MintThreshold {
			fmt.Println("status   : ditolak kontrak (skor < 4)")
		} else {
			fmt.Println("status   : bisa di-mint")
		}
		return
	}
	if *base == "" {
		fmt.Fprintln(os.Stderr, "butuh -text (atau pakai -check)")
		os.Exit(2)
	}

	text := strings.ReplaceAll(*base, `\n`, "\n")
	startScore := amulet.Score(text)
	fmt.Printf("dasar    : %q\nskor awal: %d (%s)\n", text, startScore, amulet.RarityOf(startScore))
	if startScore >= *target {
		fmt.Println("sudah memenuhi target, tidak perlu ditambang.")
		print(text, startScore)
		return
	}

	if *workers <= 0 {
		*workers = 8
	}
	var done atomic.Bool
	var mu sync.Mutex
	var got []string
	var tries atomic.Int64
	deadline := time.Now().Add(time.Duration(*maxSec) * time.Second)

	var wg sync.WaitGroup
	for w := 0; w < *workers; w++ {
		wg.Add(1)
		go func(seed int64) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano() + seed))
			vocab := []string{"✦", "✶", "✸", "·", "°", "*", "amul", "8", "888", "✦✦", "✯", "★"}
			for !done.Load() && time.Now().Before(deadline) {
				tries.Add(1)
				cand := text + " " + vocab[rng.Intn(len(vocab))] + suffix(rng)
				if amulet.Score(cand) >= *target {
					mu.Lock()
					got = append(got, cand)
					if len(got) >= *limit {
						end := true
						mu.Unlock()
						done.Store(true)
						_ = end
						return
					}
					mu.Unlock()
				}
			}
		}(int64(w))
	}
	wg.Wait()

	for _, g := range got {
		print(g, amulet.Score(g))
	}
	fmt.Printf("percobaan: %d, hasil: %d\n", tries.Load(), len(got))
	if len(got) == 0 {
		os.Exit(1)
	}
}

func suffix(rng *rand.Rand) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	n := 1 + rng.Intn(4)
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rng.Intn(len(chars))]
	}
	return string(b)
}

func print(text string, score int) {
	fmt.Println("---")
	fmt.Printf("%s\n", text)
	fmt.Printf("skor %d (%s)  sha256 %s\n", score, amulet.RarityOf(score), amulet.Hash(text))
}
