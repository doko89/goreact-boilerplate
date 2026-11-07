---
name: vibe-code
description: Delegasikan tugas coding ke agent CLI lain (claude, codex, muse, opencode) dengan berperan sebagai user yang memberi instruksi. Gunakan hanya saat user mengetik /vibe-code.
argument-hint: gunakan <claude|codex|muse|opencode> kerjakan <tugas atau file.md> [tanpa commit]
disable-model-invocation: true
allowed-tools: Read, Glob, Grep, Bash(bash .agents/skills/vibe-code/scripts/run-agent.sh:*), Bash(git status:*), Bash(git diff:*), Bash(git log:*), Bash(git show:*), Bash(git branch:*)
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
3. **Cek repo** dengan `git status`, `git branch --show-current`, dan
   `git log --oneline -5`. Catat perubahan yang sudah ada sebelum agent berjalan
   dan gaya pesan commit repo ini. Kalau ada perubahan yang belum di-commit, beri
   tahu user sebelum lanjut.
4. **Tulis prompt** ke file sementara, misalnya `/tmp/vibe-prompt.md`. Gunakan bahasa
   user, spesifik, sebutkan path file, hasil yang diharapkan, dan batasan.
   Kecuali user menulis "tanpa commit", akhiri prompt dengan instruksi commit:

   > Setelah pekerjaan selesai dan test/build (jika ada) lolos, commit hasil kerjamu:
   > - Stage hanya file yang kamu ubah untuk tugas ini. Jangan memakai `git add -A`
   >   atau `git add .`, dan jangan ikut meng-commit perubahan lain yang sudah ada
   >   sebelum kamu mulai.
   > - Tulis pesan commit yang singkat dan jelas, mengikuti gaya pesan commit di
   >   riwayat repo ini.
   > - Jangan push, jangan ubah branch, jangan rewrite history (tanpa amend,
   >   rebase, atau force).
   > - Kalau test/build gagal, jangan commit. Jelaskan apa yang gagal.
   > - Setelah commit, tampilkan hasil `git log -1 --stat`.
5. **Jalankan agent** lewat script:
   `bash .agents/skills/vibe-code/scripts/run-agent.sh <agent> /tmp/vibe-prompt.md`
   Tugas panjang boleh dijalankan di background.
6. **Review hasil**: baca output agent, lalu verifikasi sendiri dengan
   `git status`, `git log -3 --stat`, dan `git show --stat HEAD`.
   Pastikan: commit benar-benar dibuat, isinya hanya file yang relevan, pesannya
   masuk akal, dan tidak ada perubahan tersisa yang seharusnya ikut ter-commit.
7. **Beri instruksi lanjutan** seperti user kalau hasilnya kurang atau salah
   (maksimal 3 putaran), termasuk kalau commit belum dibuat atau isinya salah.
   Pakai opsi lanjutkan session kalau tersedia (lihat `agents.md`).
8. **Laporkan** ke user: yang diminta, agent yang dipakai, hash dan pesan commit,
   file yang berubah, dan apa yang masih perlu dicek manual.

## Aturan

- Jangan menulis, mengedit, atau meng-commit kode sendiri. Semua perubahan dan
  commit dilakukan oleh agent yang dipilih.
- Jangan push ke remote kecuali user memintanya secara eksplisit.
- Kalau perintah gagal (flag tidak dikenal, CLI belum terinstal), cek `<cli> --help`
  dan beri tahu user. Jangan ganti agent diam-diam.
- Jangan menjalankan agent di luar direktori project saat ini kecuali user memintanya.
- Kalau repo ada di branch utama (main/master), sarankan membuat branch baru dulu.
- Ingatkan user bahwa mode yolo mematikan approval dan sandbox.
