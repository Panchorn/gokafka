package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func Key() []byte {
	return []byte("this_is_a_key_to_encrypt_1234567")
}

func Encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// สร้าง nonce ที่มีขนาด 12 ไบต์ (96 บิต)
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// สร้าง AES cipher ในโหมด GCM
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	//fmt.Println(nonce)

	// เข้ารหัสข้อมูล
	ciphertext := aesgcm.Seal(nil, nonce, data, nil)

	// รวม nonce กับ ciphertext เพื่อส่งต่อไป
	return append(nonce, ciphertext...), nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// แยก nonce ออกจาก ciphertext
	nonce := ciphertext[:12]
	ciphertext = ciphertext[12:]

	// สร้าง AES cipher ในโหมด GCM
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// ถอดรหัสข้อมูล
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
