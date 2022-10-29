package services

import (
	"context"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/o-sokol-o/evo-fintech/internal/domain"
	"github.com/o-sokol-o/evo-fintech/pkg/restclient"
)

type ServiceEVO struct {
	repo IRepoEVO
}

func NewEvoServices(repo IRepoEVO) *ServiceEVO {
	return &ServiceEVO{repo: repo}
}

func (s *ServiceEVO) GetFilteredData(ctx context.Context, input domain.FilterSearchInput) ([]domain.Transaction, error) {
	return s.repo.GetFilteredData(ctx, input)
}

func (s *ServiceEVO) FetchExternTransactions(ctx context.Context, url string) (domain.Status, error) {
	// We request a list of transactions from an external service via REST
	transactions, err := getTransactionsRemoteURL(url)
	if err != nil || len(transactions) == 0 {
		return domain.DownloadError, err
	}

	err = s.repo.InsertTransactions(ctx, transactions)
	if err != nil {
		return domain.DownloadError, err
	}

	/*
		// Local store products
		pr, err := s.GetProducts(ctx)
		fmt.Printf("Local store products: %d\n", len(pr))
		storeProd := make(map[string]domain.Transaction)
		if len(pr) > 0 {
			for _, prd := range pr {
				storeProd[prd.Name] = prd
			}
		}

		// Product separation
		newProd := make([]domain.Transaction, 0)
		updateProd := make([]domain.Transaction, 0)
		for i := range transactions {
			prd, ok := storeProd[transactions[i].Name]
			if !ok {
				newProd = append(newProd, transactions[i])
			} else if prd.Cost != transactions[i].Cost {
				transactions[i].ChangeCount = prd.Cost + 1
				updateProd = append(updateProd, transactions[i])
			}
		}

		// Insert products
		if len(newProd) > 0 {
			fmt.Printf("Insert products: ")
			err = s.repo.InsertProducts(ctx, newProd)
			fmt.Printf("%d\n\n", len(newProd))
		}

		// Update products
		if len(updateProd) > 0 {
			fmt.Printf("Update products: ")
			err = s.repo.UpdateProducts(ctx, updateProd)
			fmt.Printf("%d\n\n", len(updateProd))
		}
	*/

	return domain.DownloadOk, nil
}

func getTransactionsRemoteURL(url string) ([]domain.Transaction, error) {
	restClient, err := restclient.NewClient(time.Second * 10)
	if err != nil {
		return nil, err
	}

	// Запрашиваем CSV у внешнего сервиса по REST
	in_csv, err := restClient.Get(url)
	if err != nil {
		return nil, err
	}

	var transactions []domain.Transaction
	// UnmarshalBytes parses the CSV from the bytes in the interface.
	gocsv.UnmarshalBytes(in_csv, &transactions)
	if err != nil {
		return nil, err
	}

	/*
		// example to read uploaded CSV file
		type csvUploadInput struct {
			CsvFile *multipart.FileHeader `form:"file" binding:"required"`
		}
		var input csvUploadInput
		if err := c.ShouldBind(&input); err != nil {
			// handle error
		}
		f, err := input.CsvFile.Open()
		if err != nil {
			// handle error
		}
		defer f.Close()
		fileBytes, err := ioutil.ReadAll(f)
		if err != nil {
			// handle error
		}
		var employee []domain.Transaction
		// UnmarshalBytes parses the CSV from the bytes in the interface.
		gocsv.UnmarshalBytes(fileBytes, &employee)
	*/

	/*
		// Получаем список продуктов
		reader := csv.NewReader(strings.NewReader(string(in_csv)))
		reader.Comma = ','
		transactions := make([]domain.Transaction, 0)
		for {
			lines, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			if len(lines) != 2 {
				continue
			}

			cost, err := strconv.Atoi(lines[1])
			if err != nil {
				cost = 0
			}
			if lines[0] != "" && cost != 0 {
				transactions = append(transactions, domain.Transaction{
					// Name:        lines[0],
					// Cost:        int32(cost),
					// UpdateAt:    time.Now(),
					// ChangeCount: 0,
				})
			}
		}
	*/
	return transactions, nil
}
