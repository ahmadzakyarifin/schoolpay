package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(env string) {
	zerolog.TimeFieldFormat = time.RFC3339

	if env == "development" || env == "local" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().Msgf("Logger initialized in %s mode", env)
}

/*
Penjelasan Fungsi InitLogger:
  zerolog.TimeFieldFormat = time.RFC3339
   Mengatur format timestamp di log agar standar internasional (contoh: 2026-06-02T14:18:21+07:00).

  Blok Percabangan `if env == "development" || env == "local":
   - Jika aplikasi berjalan di mode pengembangan (development/local):
     - Log yang dicetak di terminal akan menggunakan `ConsoleWriter` (berwarna-warni dan rapi untuk dibaca manusia).
     - Global log level disetel ke `DebugLevel` agar semua log sistem (mulai dari info umum sampai detail debug teknis) ditampilkan.
   - Jika berjalan di mode produksi (production):
     - Global log level disetel ke `InfoLevel`. Ini memfilter log sehingga log kategori 'Debug' yang terlalu detail/sensitif tidak dicetak, guna menghemat resource penyimpanan logs dan menjaga performa aplikasi. Format log juga akan berupa JSON mentah (tidak berwarna) agar mudah diindeks oleh sistem pengumpul log otomatis.

  log.Info().Msgf(...)
   Mencetak log pemberitahuan awal ke layar terminal bahwa Logger sistem sudah berhasil diaktifkan sesuai mode environment yang dipilih di file .env (APP_ENV).
*/
