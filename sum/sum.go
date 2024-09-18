package sum

import (
	"fmt"
	"io/ioutil"
	"os"
)

// readFile lê um arquivo a partir do caminho e retorna um slice de bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return nil, err
	}
	return data, nil
}

// Sums representa a soma de bytes de um arquivo e o caminho do arquivo
type Sums struct {
	Sum  int
	Path string
}

// Sum calcula a soma de todos os bytes de um arquivo e envia o resultado em um canal
func Sum(filePath string, s chan Sums) (int, error) {
	data, err := readFile(filePath)
	if err != nil {
		return 0, err
	}

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}
	fmt.Printf("%s: %v\n", filePath, _sum)
	s <- Sums{Sum: _sum, Path: filePath}
	return _sum, nil
}

func ReadFiles(s chan Sums) int {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return -1
	}
	nt := 0
	for _, path := range os.Args[1:] {
		go Sum(path, s) // Usando a função Sum do pacote sumPkg
		nt += 1
	}
	return nt
}
