---
name: vibe-code
description: Delegasikan tugas coding ke agent CLI lain (claude, codex, muse, opencode) dengan berperan sebagai user yang memberi instruksi. Gunakan hanya saat user mengetik /vibe-code.
argument-hint: gunakan <claude|codex|muse|opencode> kerjakan <tugas atau file.md>
disable-model-invocation: true
allowed-tools: Read, Glob, Grep, Bash(bash .claude/skills/vibe-code/scripts/run-agent.sh:*), Bash(git status:*), Bash(git diff:*), Bash(git log:*)
---

Input dari user: $ARGUMENTS

## Peran kamu

Kamu BUKAN yang menulis kode. Kamu bertindak sebagai **user** yang memberi instruksi
ke agent coding CLI lain. Terjemahkan permintaan menjadi prompt yang jelas, jalankan
agent secara headless (mode yolo), lalu periksa hasilnya seperti user yang teliti.

## Langkah

1. **Tentukan agent dan tugas** dari input. Kalau agent tidak disebut, tanyakan satu
   kali, jangan menebak. Daftar agent dan flag-nya ada di `agents.md`.
2. **Siapkan konteks.** Kalau tugas menunjuk file (misalnya PRD.md), baca dulu untuk
   memahami isinya. Di prompt ke agent, cukup arahkan agent membacanya dari path-nya.
3. **Tulis prompt** ke file sementara, misalnya `/tmp/vibe-prompt.md`. Gunakan bahasa
   user, spesifik, sebutkan path file, hasil yang diharapkan, dan batasan
   (jangan ubah hal di luar tugas).
4. **Cek repo** dengan `git status`. Kalau ada perubahan yang belum di-commit, beri
   tahu user sebelum menjalankan agent yolo.
5. **Jalankan agent** lewat script:
   `bash .claude/skills/vibe-code/scripts/run-agent.sh <agent> /tmp/vibe-prompt.md`
   Tugas panjang boleh dijalankan di background.
6. **Review hasil**: baca output agent, lalu cek `git status` dan `git diff`.
7. **Beri instruksi lanjutan** seperti user kalau hasilnya kurang atau salah
   (maksimal 3 putaran). Pakai opsi lanjutkan session kalau tersedia (lihat `agents.md`).
8. **Laporkan** ke user: yang diminta, agent yang dipakai, file yang berubah, dan
   apa yang masih perlu dicek manual.

## Aturan

- Jangan menulis atau mengedit kode project sendiri. Semua perubahan lewat agent yang dipilih.
- Kalau perintah gagal (flag tidak dikenal, CLI belum terinstal), cek `<cli> --help`
  dan beri tahu user. Jangan ganti agent diam-diam.
- Jangan menjalankan agent di luar direktori project saat ini kecuali user memintanya.
- Ingatkan user bahwa mode yolo mematikan approval dan sandbox, jadi sarankan
  commit atau branch baru sebelum menjalankan.
