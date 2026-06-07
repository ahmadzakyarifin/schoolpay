package helper

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	// Daftarkan TagNameFunc ke engine validator bawaan Gin.
	// Ini membuat e.Field() otomatis mengembalikan nama tag JSON-nya (misal: "birth_date" atau "nisn"),
	// sehingga kita tidak memerlukan konversi toSnakeCase manual dan map overrides.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// Register custom_date validation tag
		_ = v.RegisterValidation("custom_date", func(fl validator.FieldLevel) bool {
			dateStr := fl.Field().String()
			if dateStr == "" {
				return true
			}
			
			// Allow DD/MM/YYYY
			t, err := time.Parse("02/01/2006", dateStr)
			if err == nil && t.Format("02/01/2006") == dateStr {
				return true
			}
			
			// Allow YYYY-MM-DD
			t2, err2 := time.Parse("2006-01-02", dateStr)
			if err2 == nil && t2.Format("2006-01-02") == dateStr {
				return true
			}

			return false
		})
	}
}

// GetValidationErrors menerjemahkan error binding/validation menjadi map error dalam Bahasa Indonesia.
// Format output: map[nama_field_json][]pesan_error
func GetValidationErrors(err error) map[string][]string {
	messages := make(map[string][]string)

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			key := e.Field() // Langsung dapat "phone_number", "birth_date", dll. dari tag JSON
			valueStr := fmt.Sprintf("%v", e.Value())

			// Validasi Kustom Nomor WhatsApp menggunakan library standar Google (libphonenumber)
			if key == "phone_number" && valueStr != "" {
				if !utils.ValidatePhoneNumber(valueStr) {
					messages[key] = append(messages[key], "Nomor telepon tidak sesuai format internasional (contoh: 628...).")
					continue
				}
			}

			// Jika lolos validasi kustom, gunakan pesan error bawaan sesuai jenis tag validator-nya
			messages[key] = append(messages[key], getErrorMessage(e))
		}
	} else {
		// Menangkap error umum parsing/binding JSON (misal: format tipe data salah, dsb.)
		messages["_general"] = append(messages["_general"], "Format data tidak valid atau ada data wajib yang kosong.")
	}
	return messages
}

// getErrorMessage memetakan nama field Go ke Bahasa Indonesia dan mengembalikan pesan error yang ramah user.
func getErrorMessage(e validator.FieldError) string {
	// Kamus label untuk menerjemahkan nama key JSON ke istilah Indonesia yang dimengerti user.
	fieldLabels := map[string]string{
		// Common & Auth
		"name":             "Nama Lengkap",
		"email":            "Email",
		"phone_number":     "Nomor WhatsApp",
		"password":         "Password",
		"new_password":     "Password Baru",
		"confirm_password": "Konfirmasi Password",
		"token":            "Token Akses",
		"role":             "Hak Akses",

		// Student Specific
		"nisn":        "NISN",
		"nis":         "Nomor Induk Siswa (NIS)",
		"nik":         "NIK (KTP/KK)",
		"birth_place": "Tempat Lahir",
		"birth_date":  "Tanggal Lahir",
		"religion":    "Agama",
		"address":     "Alamat Lengkap",
		"entry_year":  "Tahun Angkatan",
		"class_id":    "Kelas",
		"major_id":    "Jurusan",
		"parent_id":   "Orang Tua/Wali",
		"rt":          "RT",
		"rw":          "RW",
		"zip_code":    "Kode Pos",
		"province":    "Provinsi",
		"city":        "Kota/Kabupaten",
		"district":    "Kecamatan",
		"village":     "Kelurahan/Desa",
		"gender":      "Jenis Kelamin",
		"status":      "Status",

		// Academic Specific
		"code":      "Kode Jurusan",
		"grade":     "Tingkat/Kelas",
		"year":      "Tahun Pelajaran",
		"is_active": "Status Aktif",

		// Finance Specific
		"bill_type_id":   "Jenis Tagihan",
		"target_type":    "Target Penerima",
		"target_id":      "ID Target",
		"amount":         "Nominal Tagihan",
		"default_amount": "Nominal Standar",
		"paid_amount":    "Jumlah Terbayar",
		"due_date":       "Tanggal Jatuh Tempo",
		"method":         "Metode Pembayaran",
		"reference":      "Referensi Pembayaran",
		"type":           "Tipe Tagihan",
	}

	// Jika nama field tidak ada di kamus, gunakan nama field mentah-mentah sebagai fallback.
	fieldLabel, ok := fieldLabels[e.Field()]
	if !ok {
		fieldLabel = e.Field()
	}

	// Buat pesan error ramah pengguna berdasarkan tag validator dari go-playground
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi.", fieldLabel)
	case "email":
		return "Format email tidak valid (contoh: user@gmail.com)."
	case "min":
		return fmt.Sprintf("%s minimal %s karakter.", fieldLabel, e.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter.", fieldLabel, e.Param())
	case "len":
		return fmt.Sprintf("%s harus terdiri dari %s karakter.", fieldLabel, e.Param())
	case "numeric":
		return fmt.Sprintf("%s harus berupa angka.", fieldLabel)
	case "unique":
		return fmt.Sprintf("%s sudah terdaftar di sistem.", fieldLabel)
	case "custom_date":
		return fmt.Sprintf("%s harus berupa format tanggal DD/MM/YYYY atau YYYY-MM-DD yang valid.", fieldLabel)
	case "eqfield":
		return fmt.Sprintf("%s harus sama dengan %s.", fieldLabel, e.Param())
	}
	return fmt.Sprintf("Data %s tidak valid.", fieldLabel)
}
