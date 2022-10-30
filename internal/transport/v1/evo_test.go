package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/o-sokol-o/evo-fintech/internal/domain"
	mock_service "github.com/o-sokol-o/evo-fintech/internal/transport/v1/mocks"
)

func getPointerString(s string) *string {
	return &s
}

func getPointerInt(x int) *int {
	return &x
}

func getPointerTime(t string) *time.Time {
	dt, err := time.Parse("2006-01-02T15:04:05.000Z", t)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &dt
}

func TestHandler_downloadRemoteTransactionsCSV(t *testing.T) {

	type mockBehavior func(s *mock_service.MockIServicesEVO, ctx context.Context, url string)

	fmt.Println("----------------  Test download Remote Transactions CSV ---------------")

	testTable := []struct {
		name                 string
		inputBody            string
		inputURL             domain.UrlInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Test case: OK",
			inputBody: `{"url":"http://12345"}`,
			inputURL:  domain.UrlInput{URL: getPointerString("http://12345")},

			mockBehavior: func(s *mock_service.MockIServicesEVO, ctx context.Context, url string) {
				s.EXPECT().FetchExternTransactions(ctx, url).Return(domain.DownloadOk, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"last_download_status":"successfully"}`,
		},
		{
			name:      "Test case: Service Failure",
			inputBody: `{"url":"http://12345"}`,
			inputURL:  domain.UrlInput{URL: getPointerString("http://12345")},

			mockBehavior: func(s *mock_service.MockIServicesEVO, ctx context.Context, url string) {
				s.EXPECT().FetchExternTransactions(ctx, url).Return(domain.DownloadError, errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service failure"}`,
		},
		{
			name:      "Test case: Bad request: Extra comma",
			inputBody: `{"url":"http://12345",}`,
			inputURL:  domain.UrlInput{URL: getPointerString("http://12345")},

			mockBehavior:         func(s *mock_service.MockIServicesEVO, ctx context.Context, url string) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			fmt.Printf("---------------- Start %s ---------------\n", testCase.name)

			//Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			servicesEVO := mock_service.NewMockIServicesEVO(c)
			testCase.mockBehavior(servicesEVO, context.Background(), *testCase.inputURL.URL)

			h := NewHandler(servicesEVO, nil)

			// Init Endpoint
			router := gin.Default()
			router.POST("/api/v1/download_remote_transactions/", h.downloadRemoteTransactionsCSV)

			// Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/download_remote_transactions/", bytes.NewBufferString(testCase.inputBody))

			// Send request
			router.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			body := w.Body.String()
			assert.Equal(t, testCase.expectedResponseBody, body)

			fmt.Printf("---------------- End %s ---------------\n", testCase.name)
		})
	}
}

func TestHandler_getFilteredDataJSON(t *testing.T) {

	type mockBehavior func(s *mock_service.MockIServicesEVO, ctx context.Context, input domain.FilterSearchInput)

	fmt.Println("----------------  Test getFilteredDataJSON ---------------")

	testTrns := domain.Transaction{
		TransactionId:      18,
		RequestId:          20190,
		TerminalId:         3523,
		PartnerObjectId:    1111,
		AmountTotal:        120,
		AmountOriginal:     120,
		CommissionPS:       0.08,
		CommissionClient:   0,
		CommissionProvider: -0.24,
		DateInput:          *getPointerTime("2022-08-23T11:58:16.000Z"),
		DatePost:           *getPointerTime("2022-08-23T14:58:16.000Z"),
		Status:             "accepted",
		PaymentType:        "cash",
		PaymentNumber:      "PS16698375",
		ServiceId:          14150,
		Service:            "Поповнення карток",
		PayeeId:            15933855,
		PayeeName:          "privat",
		PayeeBankMfo:       271768,
		PayeeBankAccount:   "UA713620688819353",
		PaymentNarrative:   "Перерахування коштів згідно договору про надання послуг А11/27123 від 19.11.2020 р.",
	}

	testTrnsStr := `{"transaction_id":18,"request_id":20190,"terminal_id":3523,"partner_object_id":1111,"amount_total":120,"amount_original":120,"commission_ps":0.08,"commission_client":0,"commission_provider":-0.24,"date_input":"2022-08-23T11:58:16Z","date_post":"2022-08-23T14:58:16Z","status":"accepted","payment_type":"cash","payment_number":"PS16698375","service_id":14150,"service":"Поповнення карток","payee_id":15933855,"payee_name":"privat","payee_bnank_mfo":271768,"payee_bnank_account":"UA713620688819353","payment_narrative":"Перерахування коштів згідно договору про надання послуг А11/27123 від 19.11.2020 р."}`

	testTable := []struct {
		name                 string
		inputBody            string
		inputFilter          domain.FilterSearchInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Test case: OK",
			inputBody: `{
						"payment_narrative": "договору про надання послуг А11/27123",
						"payment_type": "cash",
						"period": {
							"from": "2022-08-23T11:56:00.000Z",
							"to": "2022-08-24T00:00:00.000Z"
						  },
						"status": "accepted",
						"terminal_id": [3521,3522,3523,3524,3525,3526,3527,3528,3529],
						"transaction_id": 18
					  }`,

			inputFilter: domain.FilterSearchInput{
				TransactionId: getPointerInt(18),
				TerminalId:    []int{3521, 3522, 3523, 3524, 3525, 3526, 3527, 3528, 3529},
				Status:        getPointerString("accepted"),
				PaymentType:   getPointerString("cash"),
				Period: &domain.Period{ // наприклад: from 2022-08-12, to 2022-09-01 повинен повернути всі транзакції за вказаний період
					From: getPointerTime("2022-08-23T11:56:00.000Z"),
					To:   getPointerTime("2022-08-24T00:00:00.000Z"),
				},
				PaymentNarrative: getPointerString("договору про надання послуг А11/27123"), // частково вказаному
			},
			mockBehavior: func(s *mock_service.MockIServicesEVO, ctx context.Context, input domain.FilterSearchInput) {
				s.EXPECT().GetFilteredData(ctx, input).Return(
					[]domain.Transaction{
						testTrns,
						testTrns,
						testTrns,
					}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "[" + testTrnsStr + "," + testTrnsStr + "," + testTrnsStr + "]",
		},
		{
			name: "Test case: OK = no content",
			inputBody: `{
						"payment_narrative": "договору про надання послуг А11/27123",
						"payment_type": "cash",
						"period": {
							"from": "2022-08-23T11:56:00.000Z",
							"to": "2022-08-24T00:00:00.000Z"
						  },
						"status": "accepted",
						"terminal_id": [3521,3522,3523,3524,3525,3526,3527,3528,3529],
						"transaction_id": 18
					  }`,

			inputFilter: domain.FilterSearchInput{
				TransactionId: getPointerInt(18),
				TerminalId:    []int{3521, 3522, 3523, 3524, 3525, 3526, 3527, 3528, 3529},
				Status:        getPointerString("accepted"),
				PaymentType:   getPointerString("cash"),
				Period: &domain.Period{ // наприклад: from 2022-08-12, to 2022-09-01 повинен повернути всі транзакції за вказаний період
					From: getPointerTime("2022-08-23T11:56:00.000Z"),
					To:   getPointerTime("2022-08-24T00:00:00.000Z"),
				},
				PaymentNarrative: getPointerString("договору про надання послуг А11/27123"), // частково вказаному
			},
			mockBehavior: func(s *mock_service.MockIServicesEVO, ctx context.Context, input domain.FilterSearchInput) {
				s.EXPECT().GetFilteredData(ctx, input).Return(nil, nil)
			},
			expectedStatusCode:   204,
			expectedResponseBody: "",
		},
		{
			name: "Test case: Service Failure",
			inputBody: `{
						"payment_narrative": "договору про надання послуг А11/27123",
						"payment_type": "cash",
						"period": {
							"from": "2022-08-23T11:56:00.000Z",
							"to": "2022-08-24T00:00:00.000Z"
						  },
						"status": "accepted",
						"terminal_id": [3521,3522,3523,3524,3525,3526,3527,3528,3529],
						"transaction_id": 18
					  }`,

			inputFilter: domain.FilterSearchInput{
				TransactionId: getPointerInt(18),
				TerminalId:    []int{3521, 3522, 3523, 3524, 3525, 3526, 3527, 3528, 3529},
				Status:        getPointerString("accepted"),
				PaymentType:   getPointerString("cash"),
				Period: &domain.Period{ // наприклад: from 2022-08-12, to 2022-09-01 повинен повернути всі транзакції за вказаний період
					From: getPointerTime("2022-08-23T11:56:00.000Z"),
					To:   getPointerTime("2022-08-24T00:00:00.000Z"),
				},
				PaymentNarrative: getPointerString("договору про надання послуг А11/27123"), // частково вказаному
			},
			mockBehavior: func(s *mock_service.MockIServicesEVO, ctx context.Context, input domain.FilterSearchInput) {
				s.EXPECT().GetFilteredData(ctx, input).Return(nil, errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service failure"}`,
		},
		{
			name: "Test case: Bad request: Period.To absent",
			inputBody: `{
					"payment_type": "cash",
					"period": {
						"from": "2022-08-23T11:56:00.000Z"
					  },
					"status": "accepted",
				  }`,

			mockBehavior:         func(s *mock_service.MockIServicesEVO, ctx context.Context, input domain.FilterSearchInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name: "Test case: Data validation error: Type error",
			inputBody: `{
				"payment_type": "ca  sh"
			  }`,

			mockBehavior:         func(s *mock_service.MockIServicesEVO, ctx context.Context, input domain.FilterSearchInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"data validation error"}`,
		},
		{
			name: "Test case: Bad request: Extra comma",
			inputBody: `{
				"status": "accepted",
			  }`,

			mockBehavior:         func(s *mock_service.MockIServicesEVO, ctx context.Context, input domain.FilterSearchInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			fmt.Printf("---------------- Start %s ---------------\n", testCase.name)

			//Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			servicesEVO := mock_service.NewMockIServicesEVO(c)
			testCase.mockBehavior(servicesEVO, context.Background(), testCase.inputFilter)

			h := NewHandler(servicesEVO, nil)

			// Init Endpoint
			router := gin.Default()
			router.POST("/api/v1/filtered/json/", h.getFilteredDataJSON)

			// Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/filtered/json/", bytes.NewBufferString(testCase.inputBody))

			// Send request
			router.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			body := w.Body.String()
			assert.Equal(t, testCase.expectedResponseBody, body)

			fmt.Printf("---------------- End %s ---------------\n", testCase.name)
		})
	}
}

func TestHandler_getFilteredFileCSV(t *testing.T) {

	type mockBehavior func(s *mock_service.MockIServicesEVO, ctx context.Context, input domain.FilterSearchInput)

	fmt.Println("----------------  Test getFilteredDataJSON ---------------")
	/*
		testTrns := domain.Transaction{
			TransactionId:      18,
			RequestId:          20190,
			TerminalId:         3523,
			PartnerObjectId:    1111,
			AmountTotal:        120,
			AmountOriginal:     120,
			CommissionPS:       0.08,
			CommissionClient:   0,
			CommissionProvider: -0.24,
			DateInput:          *getPointerTime("2022-08-23T11:58:16.000Z"),
			DatePost:           *getPointerTime("2022-08-23T14:58:16.000Z"),
			Status:             "accepted",
			PaymentType:        "cash",
			PaymentNumber:      "PS16698375",
			ServiceId:          14150,
			Service:            "Поповнення карток",
			PayeeId:            15933855,
			PayeeName:          "privat",
			PayeeBankMfo:       271768,
			PayeeBankAccount:   "UA713620688819353",
			PaymentNarrative:   "Перерахування коштів згідно договору про надання послуг А11/27123 від 19.11.2020 р.",
		}

		testTrnsStr := `{"transaction_id":18,"request_id":20190,"terminal_id":3523,"partner_object_id":1111,"amount_total":120,"amount_original":120,"commission_ps":0.08,"commission_client":0,"commission_provider":-0.24,"date_input":"2022-08-23T11:58:16Z","date_post":"2022-08-23T14:58:16Z","status":"accepted","payment_type":"cash","payment_number":"PS16698375","service_id":14150,"service":"Поповнення карток","payee_id":15933855,"payee_name":"privat","payee_bnank_mfo":271768,"payee_bnank_account":"UA713620688819353","payment_narrative":"Перерахування коштів згідно договору про надання послуг А11/27123 від 19.11.2020 р."}`
	*/
	testTable := []struct {
		name                 string
		inputBody            string
		inputFilter          domain.FilterSearchInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		// {
		// 	name: "Test case: OK",
		// 	inputBody: `{
		// 				"payment_narrative": "договору про надання послуг А11/27123",
		// 				"payment_type": "cash",
		// 				"period": {
		// 					"from": "2022-08-23T11:56:00.000Z",
		// 					"to": "2022-08-24T00:00:00.000Z"
		// 				  },
		// 				"status": "accepted",
		// 				"terminal_id": [3521,3522,3523,3524,3525,3526,3527,3528,3529],
		// 				"transaction_id": 18
		// 			  }`,

		// 	inputFilter: domain.FilterSearchInput{
		// 		TransactionId: getPointerInt(18),
		// 		TerminalId:    []int{3521, 3522, 3523, 3524, 3525, 3526, 3527, 3528, 3529},
		// 		Status:        getPointerString("accepted"),
		// 		PaymentType:   getPointerString("cash"),
		// 		Period: &domain.Period{ // наприклад: from 2022-08-12, to 2022-09-01 повинен повернути всі транзакції за вказаний період
		// 			From: getPointerTime("2022-08-23T11:56:00.000Z"),
		// 			To:   getPointerTime("2022-08-24T00:00:00.000Z"),
		// 		},
		// 		PaymentNarrative: getPointerString("договору про надання послуг А11/27123"), // частково вказаному
		// 	},
		// 	mockBehavior: func(s *mock_service.MockIServicesEVO, ctx context.Context, input domain.FilterSearchInput) {
		// 		s.EXPECT().GetFilteredData(ctx, input).Return(
		// 			[]domain.Transaction{
		// 				testTrns,
		// 				testTrns,
		// 				testTrns,
		// 			}, nil)
		// 	},
		// 	expectedStatusCode:   200,
		// 	expectedResponseBody: "[" + testTrnsStr + "," + testTrnsStr + "," + testTrnsStr + "]",
		// },

		{
			name: "Test case: OK = no content",
			inputBody: `{
						"payment_narrative": "договору про надання послуг А11/27123",
						"payment_type": "cash",
						"period": {
							"from": "2022-08-23T11:56:00.000Z",
							"to": "2022-08-24T00:00:00.000Z"
						  },
						"status": "accepted",
						"terminal_id": [3521,3522,3523,3524,3525,3526,3527,3528,3529],
						"transaction_id": 18
					  }`,

			inputFilter: domain.FilterSearchInput{
				TransactionId: getPointerInt(18),
				TerminalId:    []int{3521, 3522, 3523, 3524, 3525, 3526, 3527, 3528, 3529},
				Status:        getPointerString("accepted"),
				PaymentType:   getPointerString("cash"),
				Period: &domain.Period{ // наприклад: from 2022-08-12, to 2022-09-01 повинен повернути всі транзакції за вказаний період
					From: getPointerTime("2022-08-23T11:56:00.000Z"),
					To:   getPointerTime("2022-08-24T00:00:00.000Z"),
				},
				PaymentNarrative: getPointerString("договору про надання послуг А11/27123"), // частково вказаному
			},
			mockBehavior: func(s *mock_service.MockIServicesEVO, ctx context.Context, input domain.FilterSearchInput) {
				s.EXPECT().GetFilteredData(ctx, input).Return(nil, nil)
			},
			expectedStatusCode:   204,
			expectedResponseBody: "",
		},
		{
			name: "Test case: Service Failure",
			inputBody: `{
						"payment_narrative": "договору про надання послуг А11/27123",
						"payment_type": "cash",
						"period": {
							"from": "2022-08-23T11:56:00.000Z",
							"to": "2022-08-24T00:00:00.000Z"
						  },
						"status": "accepted",
						"terminal_id": [3521,3522,3523,3524,3525,3526,3527,3528,3529],
						"transaction_id": 18
					  }`,

			inputFilter: domain.FilterSearchInput{
				TransactionId: getPointerInt(18),
				TerminalId:    []int{3521, 3522, 3523, 3524, 3525, 3526, 3527, 3528, 3529},
				Status:        getPointerString("accepted"),
				PaymentType:   getPointerString("cash"),
				Period: &domain.Period{ // наприклад: from 2022-08-12, to 2022-09-01 повинен повернути всі транзакції за вказаний період
					From: getPointerTime("2022-08-23T11:56:00.000Z"),
					To:   getPointerTime("2022-08-24T00:00:00.000Z"),
				},
				PaymentNarrative: getPointerString("договору про надання послуг А11/27123"), // частково вказаному
			},
			mockBehavior: func(s *mock_service.MockIServicesEVO, ctx context.Context, input domain.FilterSearchInput) {
				s.EXPECT().GetFilteredData(ctx, input).Return(nil, errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"service failure"}`,
		},
		{
			name: "Test case: Bad request: Period.To absent",
			inputBody: `{
					"payment_type": "cash",
					"period": {
						"from": "2022-08-23T11:56:00.000Z"
					  },
					"status": "accepted",
				  }`,

			mockBehavior:         func(s *mock_service.MockIServicesEVO, ctx context.Context, input domain.FilterSearchInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name: "Test case: Data validation error: Type error",
			inputBody: `{
				"payment_type": "ca  sh"
			  }`,

			mockBehavior:         func(s *mock_service.MockIServicesEVO, ctx context.Context, input domain.FilterSearchInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"data validation error"}`,
		},
		{
			name: "Test case: Bad request: Extra comma",
			inputBody: `{
				"status": "accepted",
			  }`,

			mockBehavior:         func(s *mock_service.MockIServicesEVO, ctx context.Context, input domain.FilterSearchInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			fmt.Printf("---------------- Start %s ---------------\n", testCase.name)

			//Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			servicesEVO := mock_service.NewMockIServicesEVO(c)
			testCase.mockBehavior(servicesEVO, context.Background(), testCase.inputFilter)

			h := NewHandler(servicesEVO, nil)

			// Init Endpoint
			router := gin.Default()
			router.POST("/api/v1/filtered/csv/", h.getFilteredFileCSV)

			// Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/filtered/csv/", bytes.NewBufferString(testCase.inputBody))

			// Send request
			router.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			body := w.Body.String()
			assert.Equal(t, testCase.expectedResponseBody, body)

			fmt.Printf("---------------- End %s ---------------\n", testCase.name)
		})
	}
}
