package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Layer-Edge/light-node/node"
	"github.com/Layer-Edge/light-node/utils"
	"github.com/joho/godotenv"
)

func Worker(ctx context.Context, wg *sync.WaitGroup, id int, proxy string) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d is shutting down\n", id)
			return
		default:
			fmt.Printf("Worker %d is running with proxy %s...\n", id, proxy)
			node.CollectSampleAndVerify()
			time.Sleep(5 * time.Second)
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, will try to use wallet.txt")
	}

	// Load all private keys and display their public keys and addresses
	privKeys, err := utils.LoadPrivateKeysFromFile()
	if err != nil {
		log.Fatal("Error loading private keys: ", err)
	}

	pubKeys, err := utils.GetAllCompressedPublicKeys()
	if err != nil {
		log.Fatal("Error getting public keys: ", err)
	}

	addresses, err := utils.GetAllWalletAddresses()
	if err != nil {
		log.Fatal("Error getting wallet addresses: ", err)
	}

	log.Printf("Loaded %d private keys from wallet.txt", len(privKeys))
	for i, pubKey := range pubKeys {
		log.Printf("Key %d - Compressed Public Key: %s, Address: %s", i+1, pubKey, addresses[i])
	}

	// Load proxies
	proxies, err := utils.LoadProxiesFromFile()
	if err != nil {
		log.Println("Warning: Error loading proxies: ", err)
		log.Println("Will run without proxies")
		// Create empty proxies to match the number of private keys
		proxies = make([]string, len(privKeys))
	}

	// Make sure we have enough proxies for all private keys
	// If not enough proxies, reuse them in a round-robin fashion
	if len(proxies) < len(privKeys) {
		log.Printf("Warning: Not enough proxies (%d) for all private keys (%d). Will reuse proxies.", len(proxies), len(privKeys))
		originalProxies := make([]string, len(proxies))
		copy(originalProxies, proxies)

		for i := len(proxies); i < len(privKeys); i++ {
			proxyIndex := i % len(originalProxies)
			proxies = append(proxies, originalProxies[proxyIndex])
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGABRT, syscall.SIGTERM)

	// Start a worker for each private key
	for i := 0; i < len(privKeys); i++ {
		wg.Add(1)
		proxy := ""
		if i < len(proxies) {
			proxy = proxies[i]
		}
		go Worker(ctx, &wg, i+1, proxy)
		// Add a small delay between starting workers to avoid overwhelming the system
		time.Sleep(500 * time.Millisecond)
	}

	<-signalChan
	fmt.Println("\nReceived interrupt signal. Shutting down gracefully...")

	cancel()

	wg.Wait()
	fmt.Println("All workers have shut down. Exiting..")
}
