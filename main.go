package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/crypto/sha3"
)

var operatorList = []string{}

type Holder struct {
	account solana.PublicKey
	status  int
}

func generateHolderSeed(resourceID int, prefix []byte) string {
	// Преобразование resourceID в байты
	aid := big.NewInt(int64(resourceID)).Bytes()
	// Формируем базу для сидов
	seedBase := append(prefix, aid...)
	// Вычисляем keccak hash и обрезаем до 32 символов
	hash := sha3.NewLegacyKeccak256()
	hash.Write(seedBase)
	keccakHash := hash.Sum(nil)
	return hex.EncodeToString(keccakHash)[:32]
}

func generateHolderAddress(operator string, resourceID int, program string) (solana.PublicKey, error) {
	programKey, err := solana.PublicKeyFromBase58(program)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("invalid program address: %v", err)
	}
	operatorKey, err := solana.PublicKeyFromBase58(operator)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("invalid operator address: %v", err)
	}
	prefix := []byte("holder-")
	seed := generateHolderSeed(resourceID, prefix)
	holderAddress, err := solana.CreateWithSeed(operatorKey, seed, programKey)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to create address with seed: %v", err)
	}
	//log.Println("Holder addresses:", operator, resourceID, holderAddress)
	return holderAddress, nil
}

func fetchAccountData(client *rpc.Client, account solana.PublicKey, dataSize *uint64) ([]byte, error) {
	offset := uint64(0)
	accountInfo, err := client.GetAccountInfoWithOpts(context.TODO(), account, &rpc.GetAccountInfoOpts{
		Encoding: solana.EncodingBase64,
		DataSlice: &rpc.DataSlice{
			Offset: &offset,
			Length: dataSize,
		},
	})
	if err != nil {
		//log.Println("Failed to fetch account data: ", err)
		return nil, err
	}
	if accountInfo == nil || accountInfo.Value == nil || accountInfo.Value.Data == nil {
		return nil, fmt.Errorf("account %s not found", account)
	}

	data := accountInfo.Value.Data.GetBinary()
	return data, nil
}

func checkHolderStatus(holders <-chan solana.PublicKey, resChan chan<- Holder, client *rpc.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	dataSize := uint64(8)
	for {
		select {
		case account, ok := <-holders:
			if !ok {
				return
			}
			acc := Holder{account, 0}
			data, err := fetchAccountData(client, account, &dataSize)
			if err != nil {
				log.Printf("Failed to fetch data for account %s: %v", account, err)
				resChan <- acc
				continue
			}
			if len(data) > 0 {
				//if data[0] != 32 && data[0] != 52 && data[0] != 25 {
				//	log.Println("Unmatched data value: ", data[0], " for account: ", account)
				//}
				//switch data[0] {
				//case 32:
				//	acc.status = "finalized"
				//case 52:
				//	acc.status = "clean"
				//case 25:
				//	acc.status = "in-use"
				//	log.Println("In use account: ", account)
				//default:
				//	acc.status = "unmatched"
				//}
				acc.status = int(data[0])
			}
			resChan <- acc
		}
	}
}

func readOperatorsFromFile(filename string) ([]string, error) {
	path := filepath.Join("keys", filename)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла %s: %v", filename, err)
	}
	defer file.Close()

	var operators []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if addr := strings.TrimSpace(scanner.Text()); addr != "" {
			operators = append(operators, addr)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error read file  %s: %v", filename, err)
	}

	return operators, nil
}

func writeHolderToFile(holder Holder) error {
	// Создаём папку holders, если не существует
	dir := "holders"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("ошибка создания папки %s: %v", dir, err)
	}

	fileName := filepath.Join(dir, fmt.Sprintf("%d.txt", holder.status))
	if holder.status == 0 {
		fileName = filepath.Join(dir, "notexist.txt")
	}

	// Открываем файл для дозаписи
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла %s: %v", fileName, err)
	}
	defer f.Close()

	// Записываем адрес с переводом строки
	if _, err := f.WriteString(holder.account.String() + "\n"); err != nil {
		return fmt.Errorf("ошибка записи в файл %s: %v", fileName, err)
	}

	return nil
}

func cleanHoldersDirectory() error {
	// Создаём папку holders, если не существует
	dir := "holders"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("ошибка создания папки %s: %v", dir, err)
	}

	// Читаем содержимое папки
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("ошибка чтения папки %s: %v", dir, err)
	}

	// Удаляем все файлы
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("ошибка удаления файла %s: %v", path, err)
		}
	}

	return nil
}

func main() {
	var (
		rpcURL    string
		holderNum int
		workerNum int
		program   string
		keyFiles  string
	)

	flag.StringVar(&rpcURL, "rpc", "", "RPC URL for Solana node")
	flag.IntVar(&holderNum, "holders", 4, "Number of holder accounts per operator (1-32)")
	flag.IntVar(&workerNum, "workers", 16, "Number of parallel workers")
	flag.StringVar(&program, "program", "NeonVMyRX5GbCrsAHnUwx1nYYoJAtskU1bWUo6JGNyG", "Program address")
	flag.StringVar(&keyFiles, "key_files", "p2p.txt,everstake.txt", "Comma-separated list of key files in 'keys' directory")
	flag.Parse()

	if rpcURL == "" || program == "" || keyFiles == "" {
		log.Fatal("RPC URL, program address and key_files must be specified")
	}

	if err := cleanHoldersDirectory(); err != nil {
		log.Fatal("Ошибка очистки папки holders:", err)
	}

	files := strings.Split(keyFiles, ",")
	for _, file := range files {
		operators, err := readOperatorsFromFile(strings.TrimSpace(file))
		if err != nil {
			log.Fatal(err)
		}
		operatorList = append(operatorList, operators...)
	}

	if len(operatorList) == 0 {
		log.Fatal("Не найдено операторов в указанных файлах")
	}

	if holderNum < 1 || holderNum > 256 {
		log.Fatal("holders must be between 1 and 256")
	}

	client := rpc.New(rpcURL)
	var holderAccounts []solana.PublicKey
	var operatorHolders = make(map[string][]solana.PublicKey)

	for _, operator := range operatorList {
		operatorHolders[operator] = make([]solana.PublicKey, holderNum+1)
		for i := 0; i <= holderNum; i++ {
			holderAddress, err := generateHolderAddress(operator, i, program)
			if err != nil {
				log.Fatalf("Error generating holder address: %v", err)
			}
			operatorHolders[operator][i] = holderAddress
			holderAccounts = append(holderAccounts, holderAddress)
		}
	}

	fi, err := os.Create("holder_accounts.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := fi.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	w := bufio.NewWriter(fi)
	for operator, account := range operatorHolders {
		_, err := w.WriteString("Operator: " + operator + "\n")
		if err != nil {
			log.Fatal(err)
		}
		for _, holder := range account {
			_, err := w.WriteString(holder.String() + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}

	}

	resChan := make(chan Holder, len(holderAccounts))

	wg := &sync.WaitGroup{}
	wg.Add(workerNum)

	accountsChan := make(chan solana.PublicKey, len(holderAccounts))
	for i := 0; i < len(holderAccounts); i++ {
		accountsChan <- holderAccounts[i]
	}
	close(accountsChan)

	for i := 0; i < workerNum; i++ {
		go checkHolderStatus(accountsChan, resChan, client, wg)
	}
	wg.Wait()
	close(resChan)

	var statuses = make(map[string]int)
	statuses["finalized"] = 0
	statuses["unmatched"] = 0
	statuses["clean"] = 0
	statuses["in-use"] = 0
	statuses["notexist"] = 0

	for res := range resChan {
		if err := writeHolderToFile(res); err != nil {
			log.Printf("Ошибка записи в файл: %v", err)
		}

		switch res.status {
		case 32:
			statuses["finalized"]++
		case 52:
			statuses["clean"]++
		case 25:
			statuses["in-use"]++
		case 0:
			statuses["notexist"]++
		default:
			statuses["unmatched"]++
		}
	}

	fmt.Printf("Total holder accounts: %d\n", len(holderAccounts))
	fmt.Printf("Finalized accounts (32): %d\n", statuses["finalized"])
	fmt.Printf("Clean accounts (52): %d\n", statuses["clean"])
	fmt.Printf("In use accounts (25): %d\n", statuses["in-use"])
	fmt.Printf("Unmatched accounts: %d\n", statuses["unmatched"])
	fmt.Printf("Not exist accounts: %d\n", statuses["notexist"])
}
