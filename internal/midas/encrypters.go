package validators

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"reflect"
	"strings"
)

var AESKEY = "thisis32byteslongpassphrase!!" // 16, 24, or 32 bytes
var rsaPrivateKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQC+XsvuY47mIxX+LGHpVMUI/Zzb+CCrltxMHo2eFh6ALryJkHAF
J/GTAN39TvCuWRArJBlLTecBy722yU0Oq3FedLG+87Ozt+wXfw6z4ljGv9peb5PD
aZIrDfOMV6fR8dTn/QD/YQzxbDR4bS+td60Sp1abrJobS7HDWpKxn4eaCQIDAQAB
AoGBAJlUHuq90022+rLCqoXVWfGY2yssmZu6oWrJvQSUCjLx8bbW1/K/LkpZi3LH
jwDOCUDGDX1inGoM4JuFXQQwv9iPtixwCfLD/JRJQ1omRyJx2ClBsiPdNje1GQMO
Bym9L31ZvTSQjPH5nt+rdoHpGigbH+4RLk+ElLTkBJd/zN4RAkEA9C4ES4XeiKN8
yUbWh57+P3pRAXta3M90LlzJiI+/zH28kASKaD9TqUkhvCoHK4PJnFpR8Op+IV+b
DBm5uiuTLQJBAMeV8wfdbYYoU8HiBJntgOSyT3xizrNnZsp6k+Lnl8DHCe+ybBuG
lmRHoU9SGxp+VhtpbG20Bh8pCbv9YvCxG80CQQCO/1Pslp1YD8ZIaX/BNM9YhV1j
LMZtgeBcNmKf4u9D5m7DOKWFn3BzNyzWcRZ52Vf8hLhwCiOLj93RHE+0Q0iRAkBB
wMXzZmDJ3QlTC7pGV/ep4JDNQuQkOMGlnWKRU4ksSqacYGS7YMi1OuAK+NrTDKIj
n8TIE5Icu/FoDJ+G+mJVAkBkGEvkuxLgtYVbT1Ch5/NcPQB4N6bhAXRAB8Hn+jJD
EJ4XxIfdgn6lQoQRufFWt1EHMIboASYbC1QnsdSfhtVB
-----END RSA PRIVATE KEY-----`
var rsaPublicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC+XsvuY47mIxX+LGHpVMUI/Zzb
+CCrltxMHo2eFh6ALryJkHAFJ/GTAN39TvCuWRArJBlLTecBy722yU0Oq3FedLG+
87Ozt+wXfw6z4ljGv9peb5PDaZIrDfOMV6fR8dTn/QD/YQzxbDR4bS+td60Sp1ab
rJobS7HDWpKxn4eaCQIDAQAB
-----END PUBLIC KEY-----`

func ProcessSecrets(val interface{}, encryptMode bool) error {
	v := reflect.ValueOf(val).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := v.Type().Field(i).Tag.Get("secret")

		if tag == "" {
			continue
		}

		fieldName := v.Type().Field(i).Name
		fieldValue := field.String()

		parts := strings.Split(tag, ":")
		operation := parts[0]
		scheme := "default"
		if len(parts) > 1 {
			scheme = parts[1]
		}

		switch operation {
		case "encrypt":
			if encryptMode {
				encryptedValue, err := applyEncryption(fieldValue, scheme)
				if err != nil {
					return fmt.Errorf("failed to encrypt %s: %v", fieldName, err)
				}
				field.SetString(encryptedValue)
			} else {
				decryptedValue, err := applyDecryption(fieldValue, scheme)
				if err != nil {
					return fmt.Errorf("failed to decrypt %s: %v", fieldName, err)
				}
				field.SetString(decryptedValue)
			}
		case "mask":
			if encryptMode {
				maskedValue, original := mask(fieldValue)
				field.SetString(maskedValue)

				// Store the original value (you may store it in a separate map or a struct for recovery)
				// Here, we store the original value directly on the struct.
				field.Set(reflect.ValueOf(original)) // This is for example, adjust as per the real requirement
			} else {
				unmaskedValue := unmask(fieldValue)
				field.SetString(unmaskedValue)
			}
		}
	}
	return nil
}

func applyEncryption(data, scheme string) (string, error) {
	switch scheme {
	case "AES":
		return aesEncrypt(data, AESKEY)
	case "RSA":
		pubKey, err := parseRSAPublicKey(rsaPublicKeyPEM)
		if err != nil {
			fmt.Println(err)
		}
		return rsaEncrypt(data, pubKey)
	default:
		return encrypt(data)
	}
}

func applyDecryption(data, scheme string) (string, error) {
	switch scheme {
	case "AES":
		return aesDecrypt(data, AESKEY)
	case "RSA":
		privKey, _ := parseRSAPrivateKey(rsaPrivateKeyPEM)
		return rsaDecrypt(data, privKey)
	default:
		return decrypt(data)
	}
}

func padKey(key string, requiredLength int) string {
	if len(key) >= requiredLength {
		return key[:requiredLength] // Truncate if too long
	}
	return key + strings.Repeat("0", requiredLength-len(key)) // Pad with "0"s if too short
}

/*
func hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return string(hash[:])
}
*/

// AES Encrypt
func aesEncrypt(data, key string) (string, error) {
	key = padKey(key, 32) // or hashKey(key) to ensure 32 bytes

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(data))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AES Decrypt
func aesDecrypt(encryptedData, key string) (string, error) {
	key = padKey(key, 32) // or hashKey(key) to ensure 32 bytes

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// RSA Encrypt
func rsaEncrypt(data string, publicKey *rsa.PublicKey) (string, error) {
	if publicKey == nil {
		return "", fmt.Errorf("RSA public key is nil")
	}

	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(data), nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data using RSA: %v", err)
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// RSA Decrypt
func rsaDecrypt(encryptedData string, privateKey *rsa.PrivateKey) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}

// Parse RSA Private Key
func parseRSAPrivateKey(pemEncoded string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemEncoded))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func parseRSAPublicKey(pemEncoded string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemEncoded))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	var pubKey *rsa.PublicKey
	var err error
	if block.Type == "PUBLIC KEY" {
		// PKCS#8 format (-----BEGIN PUBLIC KEY-----)
		pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS#8 public key: %v", err)
		}
		pubKey = pubKeyInterface.(*rsa.PublicKey)
	} else if block.Type == "RSA PUBLIC KEY" {
		// PKCS#1 format (-----BEGIN RSA PUBLIC KEY-----)
		pubKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS#1 public key: %v", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported public key type: %s", block.Type)
	}

	return pubKey, nil
}

func encrypt(data string) (string, error) {
	return "encrypted_" + data, nil
}

func decrypt(data string) (string, error) {
	return strings.TrimPrefix(data, "encrypted_"), nil
}

// Mask and unmask implementation
func mask(data string) (string, string) {
	if len(data) <= 4 {
		return "****", data // Return the original value to restore later
	}
	masked := strings.Repeat("*", len(data)-4) + data[len(data)-4:]
	return masked, data // Store the original value for later recovery
}

func unmask(data string) string {
	if strings.Contains(data, "*") {
		return strings.ReplaceAll(data, "*", "") // Reverse the mask to get the original data
	}
	return data // If not masked, return the data as is
}
