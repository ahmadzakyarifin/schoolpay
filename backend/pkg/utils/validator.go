package utils

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
	"regexp"
)

var (
	nikRegex = regexp.MustCompile(`^[0-9]{16}$`)
)

// ValidatePhoneStruct is a custom validator for go-playground/validator
func ValidatePhoneStruct(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true
	}
	return ValidatePhoneNumber(phone)
}

// ValidateNIKStruct is a custom validator for go-playground/validator
func ValidateNIKStruct(fl validator.FieldLevel) bool {
	nik := fl.Field().String()
	if nik == "" {
		return true
	}
	return nikRegex.MatchString(nik)
}

func GetValidationErrors(err error) map[string][]string {
	messages := make(map[string][]string)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			key := toSnakeCase(e.Field())
			overrides := map[string]string{
				"NISN": "nisn", "NIK": "nik", "NIS": "nis", "RT": "rt", "RW": "rw",
				"JurusanID": "jurusan_id", "ClassID": "class_id", "ParentID": "parent_id",
				"PhoneNumber": "phone_number", "BirthDate": "birth_date", "EntryYear": "entry_year",
			}
			if val, ok := overrides[e.Field()]; ok {
				key = val
			}
			messages[key] = append(messages[key], getErrorMessage(e))
		}
	} else if err != nil {
		errStr := err.Error()
		errStrLower := strings.ToLower(errStr)

		// Detect birth_date errors by content or field name
		if strings.Contains(errStrLower, "tidak valid dalam kalender nyata") ||
			strings.Contains(errStrLower, "birth_date") ||
			strings.Contains(errStrLower, "birthdate") {

			if strings.Contains(errStr, "tidak valid dalam kalender nyata") {
				parts := strings.Split(errStr, ": ")
				msg := parts[len(parts)-1]
				messages["birth_date"] = append(messages["birth_date"], msg)
			} else {
				messages["birth_date"] = append(messages["birth_date"], "Format tanggal lahir tidak valid (gunakan DD/MM/YYYY)")
			}
		} else {
			// Catch generic binding errors (e.g. "cannot unmarshal string into uint")
			messages["_general"] = append(messages["_general"], "Format data tidak valid atau ada data wajib yang kosong.")
		}
	}
	return messages
}

func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func getErrorMessage(e validator.FieldError) string {
	fieldLabels := map[string]string{
		// Common & Auth
		"Name":        "Nama Lengkap",
		"Email":       "Email",
		"PhoneNumber": "Nomor WhatsApp",
		"Password":    "Password",
		"NewPassword": "Password Baru",
		"Token":       "Token Akses",
		"Role":        "Hak Akses",

		// Student Specific
		"NISN":       "NISN",
		"NIS":        "Nomor Induk Siswa (NIS)",
		"NIK":        "NIK (KTP/KK)",
		"BirthPlace": "Tempat Lahir",
		"BirthDate":  "Tanggal Lahir",
		"Religion":   "Agama",
		"Address":    "Alamat Lengkap",
		"EntryYear":  "Tahun Angkatan",
		"ClassID":    "Kelas",
		"JurusanID":  "Jurusan",
		"ParentID":   "Orang Tua/Wali",
		"RT":         "RT",
		"RW":         "RW",
		"ZipCode":    "Kode Pos",
		"Province":   "Provinsi",
		"City":       "Kota/Kabupaten",
		"District":   "Kecamatan",
		"Village":    "Kelurahan/Desa",
		"Gender":     "Jenis Kelamin",
		"Status":     "Status",

		// Academic Specific
		"Code":     "Kode Jurusan",
		"Grade":    "Tingkat/Kelas",
		"Year":     "Tahun Pelajaran",
		"IsActive": "Status Aktif",

		// Finance Specific
		"BillTypeID":    "Jenis Tagihan",
		"TargetType":    "Target Penerima",
		"TargetID":      "ID Target",
		"Amount":        "Nominal Tagihan",
		"DefaultAmount": "Nominal Standar",
		"PaidAmount":    "Jumlah Terbayar",
		"DueDate":       "Tanggal Jatuh Tempo",
		"Method":        "Metode Pembayaran",
		"Reference":     "Referensi Pembayaran",
		"Type":          "Tipe Tagihan",
	}

	fieldLabel, ok := fieldLabels[e.Field()]
	if !ok {
		fieldLabel = e.Field()
	}

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi.", fieldLabel)
	case "email":
		return "Format email tidak valid (contoh: user@gmail.com)."
	case "min":
		return fmt.Sprintf("%s minimal %s karakter.", fieldLabel, e.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter.", fieldLabel, e.Param())
	case "numeric":
		return fmt.Sprintf("%s harus berupa angka.", fieldLabel)
	case "unique":
		return fmt.Sprintf("%s sudah terdaftar di sistem.", fieldLabel)
	case "phone", "custom_phone":
		return "Nomor telepon tidak sesuai format internasional (contoh: 628...).."
	case "custom_nik":
		return "NIK harus terdiri dari 16 digit angka."
	}
	return fmt.Sprintf("Data %s tidak valid.", fieldLabel)
}

func ValidatePhoneNumber(phone string) bool {
	num, err := phonenumbers.Parse("+"+phone, "")
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumber(num)
}

func IsValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
