package masyarakatcontrollers

import (
	"Backend_TA/models"
	"Backend_TA/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Show(c *fiber.Ctx) error {
	var user []models.User
	if err := models.DB.Preload("Masyarakat").Joins("JOIN masyarakat ON masyarakat.nik = user.id").Find(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}
	data := make([]fiber.Map, len(user))
	for i, users := range user {
		data[i] = fiber.Map{
			"nik":           users.ID,
			"email":         users.Email,
			"password":      users.Password,
			"id_masyarakat": users.Masyarakat.Idm,
			"nama":          users.Masyarakat.Nama,
			"tempat_lahir":  users.Masyarakat.Tempat_lahir,
			"birthday":      users.Masyarakat.Birthday[0:10],
			"gender":        users.Masyarakat.Gender,
			"no_hp":         users.Masyarakat.No_hp,
			"alamat":        users.Masyarakat.Alamat,
			"createdAt":     users.Masyarakat.CreatedAt.String()[0:10],
		}
	}
	return c.JSON(data)
}

func ShowId(c *fiber.Ctx) error {
	nik := c.Params("nik")

	if nik == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "NIK kosong"})
	}

	tx := models.DB
	var user models.User
	if err := tx.Preload("Masyarakat").Where("user.id = ?", nik).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"msg": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{
		"nik":           user.ID,
		"email":         user.Email,
		"password":      user.Password,
		"id_masyarakat": user.Masyarakat.Idm,
		"nama":          user.Masyarakat.Nama,
		"tempat_lahir":  user.Masyarakat.Tempat_lahir,
		"birthday":      user.Masyarakat.Birthday[0:10],
		"gender":        user.Masyarakat.Gender,
		"no_hp":         user.Masyarakat.No_hp,
		"alamat":        user.Masyarakat.Alamat,
		"createdAt":     user.Masyarakat.CreatedAt.String()[0:10],
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	tx := models.DB
	nik := c.Params("nik")

	if nik == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "NIK required"})
	}

	// if err := tx.Where("id")
	var cekData models.User
	if err := tx.Preload("Masyarakat").Where("user.id = ?", nik).First(&cekData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"msg": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	var user models.User
	user.ID = nik
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	if user.Email == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Email tidak boleh kosong"})
	}
	var masyarakat models.Masyarakat
	masyarakat.NIK = nik
	if err := c.BodyParser(&masyarakat); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	if masyarakat.Nama == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Nama kosong"})
	}

	if masyarakat.No_hp == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Nomor hp kosong"})
	}

	if masyarakat.Tempat_lahir == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Tempat lahir kosong"})
	}

	if masyarakat.Birthday == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Tanggal lahir kosong"})
	}

	if masyarakat.Alamat == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Alamat kosong"})
	}

	if err := models.ValidateMasyarakat(&masyarakat); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg_validate": err.Error()})
	}

	if err := tx.Where("id = ?", nik).Updates(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}
	if err := tx.Where("nik = ?", nik).Updates(&masyarakat).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{"msg": "Profile berhasil di update"})
}

// Cek ke forum
func UpdatePassword(c *fiber.Ctx) error {
	nik := c.Params("nik")

	if nik == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "NIK kosong"})
	}
	var user models.User
	if err := models.DB.Where("id =?", nik).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"msg": "User tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"msg": err.Error()})
	}

	var isValid models.NewPassword
	if err := c.BodyParser(&isValid); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}

	if isValid.Old_pass == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Old password tidak boleh kosong"})
	}

	if isValid.New_pass == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "New password tidak boleh kosong"})
	}

	if isValid.Konf_pass == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Konfirmasi password tidak boleh kosong"})
	}

	if isValid.New_pass != isValid.Konf_pass {
		return c.Status(400).JSON(fiber.Map{"msg": "Password tidak sesuai"})
	}

	if err := models.ValidatePass(&isValid); err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": "Password harus berjumlah minimal 8 karakter"})
	}

	isValid.Old_pass = utils.EncryptHash(isValid.Old_pass)
	if isValid.Old_pass != user.Password {
		return c.Status(400).JSON(fiber.Map{"msg": "Password lama tidak sesuai"})
	}

	user.Password = utils.EncryptHash(isValid.New_pass)
	user.Konf_pass = utils.EncryptHash(isValid.Konf_pass)

	if err := models.DB.Where("id =?", nik).Updates(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": err.Error()})
	}
	return c.JSON(fiber.Map{"msg": "Password berhasil diubah"})
}

func DeleteProfile(c *fiber.Ctx) error {
	nik := c.Params("nik")
	if nik == "" {
		return c.Status(404).JSON(fiber.Map{"msg": "NIK required"})
	}

	var user models.User
	if err := models.DB.Where("id = ?", nik).Delete(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"msg": "User not found"})
		}

		return c.Status(500).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{"msg": "User berhasil dihapus!"})
}
