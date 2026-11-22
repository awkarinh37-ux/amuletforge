# amuletforge

Penambang jimat untuk [Amulets](https://at.amulet.garden) — puisi pendek yang kelangkaannya
ditentukan oleh sha256-nya sendiri.

Skor sebuah jimat adalah **panjang deret nibble hex `8` terpanjang** pada sha256 teksnya.
Kontrak menolak apa pun di bawah 4; 4 = common, 5 = uncommon, 6 = rare, 7 = epic,
8 = legendary, 9 = mythic, 10+ = beyond mythic.

Alat ini menambang variasi teks sampai skornya cukup langka, lalu mencetak sha256-nya
supaya bisa diperiksa ulang.

## Kecocokan dengan kontrak

Skor dihitung sama persis dengan `Amulet.getScore`, dan **dicocokkan langsung ke kontrak
mainnet** (`0x2A2127753653F6210d26F5B470738bf13B73423d`):

| teks | kontrak | alat ini |
|---|---|---|
| `snap` | 4 | 4 |
| `= 69` | 5 | 5 |
| `at.amulet.garden` | 4 | 4 |
| `If you can't write poems,\nwrite me` | 5 | 5 |
| `kerajinan tangan kecil *4cqu` (hasil tambang) | 6 | 6 |

## Pakai

```bash
go build -o forge .

# periksa teks yang sudah ada
./forge -check "snap"

# tambang jimat dari teks dasar
./forge -text "baris pertama\nbaris kedua" -min 6

# ambil 3 hasil sekaligus, batas 60 detik
./forge -text "..." -min 7 -n 3 -t 60
```

```
dasar    : "kerajinan tangan kecil"
skor awal: 2 (bukan jimat)
---
kerajinan tangan kecil *4cqu
skor 6 (rare)  sha256 0bd427531ad707f88d036c888888cee660e9dffdcca9e79126aeca9a1d58a0ed
percobaan: 991312, hasil: 1
```

## Test

```bash
go test ./...
```

Memuat empat teks nyata dari daftar premine di repo Amulets sebagai patokan.

## Catatan

- Skor tinggi makin mahal secara eksponensial: 4 muncul ~1/16 percobaan, 9 ~1/16⁶.
- `-min 9` bisa jalan lama. Pakai `-t` sebagai batas waktu.
- Token ID on-chain memakai `keccak256` teks, sedangkan alat ini mencetak `sha256`.
  Keduanya berbeda; sha256 dipakai karena itu dasar skornya, bukan untuk menebak id.
- Teks dikirim apa adanya. Spasi/kapital/emoji mengubah hash, jadi jangan mengubahnya
  setelah ditambang — skor akan hilang.
