package utils

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

var (
	privateKeys []string
	keysLoaded  bool
	keysMutex   sync.Mutex
)

// LoadPrivateKeysFromFile loads private keys from wallet.txt file
func LoadPrivateKeysFromFile() ([]string, error) {
	keysMutex.Lock()
	defer keysMutex.Unlock()

	if keysLoaded && len(privateKeys) > 0 {
		return privateKeys, nil
	}

	file, err := os.Open("wallet.txt")
	if err != nil {
		// If wallet.txt doesn't exist, try to use the private key from .env
		privKey := GetEnv("PRIVATE_KEY", "")
		if privKey != "" {
			privateKeys = []string{privKey}
			keysLoaded = true
			return privateKeys, nil
		}
		return nil, fmt.Errorf("failed to open wallet.txt and no PRIVATE_KEY in .env: %v", err)
	}
	defer file.Close()

	var keys []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key := strings.TrimSpace(scanner.Text())
		if key != "" {
			keys = append(keys, key)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading wallet.txt: %v", err)
	}

	if len(keys) == 0 {
		// If wallet.txt is empty, try to use the private key from .env
		privKey := GetEnv("PRIVATE_KEY", "")
		if privKey != "" {
			keys = []string{privKey}
		} else {
			return nil, fmt.Errorf("wallet.txt is empty and no PRIVATE_KEY in .env")
		}
	}

	privateKeys = keys
	keysLoaded = true
	return keys, nil
}

// GetRandomPrivateKey returns a random private key from the loaded keys
func GetRandomPrivateKey() (string, error) {
	keys, err := LoadPrivateKeysFromFile()
	if err != nil {
		return "", err
	}

	if len(keys) == 0 {
		return "", fmt.Errorf("no private keys available")
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(keys))
	return keys[randomIndex], nil
}

func GetCompressedPublicKey() (string, error) {
	privKey, err := GetRandomPrivateKey()
	if err != nil {
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	// Get the public key
	publicKey := privateKey.Public().(*ecdsa.PublicKey)

	// Serialize the public key in compressed format
	compressedPubKey := secp256k1.CompressPubkey(publicKey.X, publicKey.Y)

	// Convert to hex string
	return hex.EncodeToString(compressedPubKey), nil
}

// GetAllCompressedPublicKeys returns all compressed public keys from the loaded private keys
func GetAllCompressedPublicKeys() ([]string, error) {
	keys, err := LoadPrivateKeysFromFile()
	if err != nil {
		return nil, err
	}

	var pubKeys []string
	for _, privKey := range keys {
		privateKey, err := crypto.HexToECDSA(privKey)
		if err != nil {
			log.Printf("Invalid private key: %v", err)
			continue
		}

		publicKey := privateKey.Public().(*ecdsa.PublicKey)
		compressedPubKey := secp256k1.CompressPubkey(publicKey.X, publicKey.Y)
		pubKeys = append(pubKeys, hex.EncodeToString(compressedPubKey))
	}

	if len(pubKeys) == 0 {
		return nil, fmt.Errorf("no valid public keys generated")
	}

	return pubKeys, nil
}

func GetWalletAddress() (*string, error) {
	privKey, err := GetRandomPrivateKey()
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.PublicKey

	walletAddress := crypto.PubkeyToAddress(publicKey).Hex()
	return &walletAddress, nil
}

// GetAllWalletAddresses returns all wallet addresses from the loaded private keys
func GetAllWalletAddresses() ([]string, error) {
	keys, err := LoadPrivateKeysFromFile()
	if err != nil {
		return nil, err
	}

	var addresses []string
	for _, privKey := range keys {
		privateKey, err := crypto.HexToECDSA(privKey)
		if err != nil {
			log.Printf("Invalid private key: %v", err)
			continue
		}

		publicKey := privateKey.PublicKey
		walletAddress := crypto.PubkeyToAddress(publicKey).Hex()
		addresses = append(addresses, walletAddress)
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("no valid wallet addresses generated")
	}

	return addresses, nil
}

func SignMessage(message string) (*string, error) {
	privKey, err := GetRandomPrivateKey()
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)

	data := []byte(prefix)
	hash := crypto.Keccak256Hash(data)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}
	signature[64] += 27

	hexSign := hexutil.Encode(signature)
	return &hexSign, nil
}

// SignMessageWithSpecificKey signs a message with a specific private key
func SignMessageWithSpecificKey(message string, privKeyIndex int) (*string, error) {
	keys, err := LoadPrivateKeysFromFile()
	if err != nil {
		return nil, err
	}

	if privKeyIndex < 0 || privKeyIndex >= len(keys) {
		return nil, fmt.Errorf("private key index out of range")
	}

	privKey := keys[privKeyIndex]
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)

	data := []byte(prefix)
	hash := crypto.Keccak256Hash(data)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}
	signature[64] += 27

	hexSign := hexutil.Encode(signature)
	return &hexSign, nil
}

func VerifyMessage(sign string, message string, expectedAddress string) error {
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	messageHash := crypto.Keccak256Hash([]byte(prefix))
	signature, err := hexutil.Decode(sign)
	if err != nil {
		log.Printf("invalid signature")
		return err
	}
	if len(signature) != 65 {
		log.Println("invalid signature length")
		return fmt.Errorf("invalid signature length")
	}

	signature[64] -= 27

	pubKey, err := crypto.SigToPub(messageHash.Bytes(), signature)
	if err != nil {
		log.Println("Failed to recover public key:", err)
		return fmt.Errorf("Failed to recover public key:", err)
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()

	log.Printf("recovered address: %s, expected address: %s", recoveredAddress, expectedAddress)
	return nil
}
