package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"testwork/client/queries"
	"testwork/sdkInit"
	"testwork/transformer"
	"testwork/transformer/dbop"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

const (
	cc_name    = "sacc"
	cc_version = "1.8"
)

var databases []*dbop.MySQLDB
var App sdkInit.Application //app

func Tests() {
	block := transformer.Block{
		BlockID:          0,
		BlockHash:        "0x123",
		PreviousHash:     "",
		CreateTime:       "2022-01-03 01:00:00",
		TransactionCount: 2,
	}

	transactions := []transformer.Transaction{
		{
			TransactionHash: "0x456",
			TransactionType: "transfer",
			CreateTime:      "2022-01-01 00:00:00",
			BlockID:         0,
			Sender:          "Alice",
			Receiver:        "Bob",
			Amount:          1000,
			Memo:            "Payment for services",
			SenderBalance:   9000,
			ReceiverBalance: 12000,
		},
		{
			TransactionHash: "0x789",
			TransactionType: "transfer",
			CreateTime:      "2022-01-03 00:00:00",
			BlockID:         0,
			Sender:          "Bob",
			Receiver:        "Charlie",
			Amount:          10000,
			Memo:            "Payment for goods",
			SenderBalance:   2000,
			ReceiverBalance: 11000,
		},
	}

	accounts := []transformer.Account{
		{
			AccountID:     "Alice",
			AccountType:   "individual",
			CreateTime:    "2021-01-03 00:00:00",
			AccountStatus: "active",
			Balance:       9000,
			CreditScore:   700,
			LastUpdate:    "2022-01-01 00:00:00",
		},
		{
			AccountID:     "Bob",
			AccountType:   "individual",
			CreateTime:    "2021-01-03 00:00:00",
			AccountStatus: "active",
			Balance:       2000,
			CreditScore:   700,
			LastUpdate:    "2022-01-03 00:00:00",
		},
		{
			AccountID:     "Charlie",
			AccountType:   "individual",
			CreateTime:    "2021-01-03 00:00:00",
			AccountStatus: "active",
			Balance:       11000,
			CreditScore:   700,
			LastUpdate:    "2022-01-03 00:00:00",
		},
	}

	operations := []transformer.Operation{
		{
			OperationID:     2,
			OperationType:   "C",
			CreateTime:      "2022-01-01 00:00:00",
			TransactionHash: "0x456",
			AccountID:       "Alice",
			Balance:         9000,
			State:           "active",
		},
		{
			OperationID:     3,
			OperationType:   "C",
			CreateTime:      "2022-01-03 00:00:00",
			TransactionHash: "0x789",
			AccountID:       "Bob",
			Balance:         2000,
			State:           "active",
		},
	}
	err := initDB()
	if err != nil {
		fmt.Println("init failed, err:%v\n", err)
		return
	}
	fmt.Println("InsertDataToFabric---------------InsertDataToFabric")
	InsertDataToFabric(block, accounts, transactions, operations)
	fmt.Println("InsertDataToDatabases---------InsertDataToDatabases")
	InsertDataToDatabases(block, accounts, transactions, operations)
	fmt.Println("QueryDataFromDatabases-------QueryDataFromDatabases")
	bs, ts, counts, balances, trss, trass := QueryDataFromDatabases()
	b := pbftbs(bs)
	fmt.Println(b)
	t := pbftts(ts)
	fmt.Println(t)
	fmt.Println(pbftcs(counts))
	fmt.Println(pbftbls(balances))
	fmt.Println(pbfttrs(trss))
	fmt.Println(pbfttras(trass))
	fmt.Println("QueryDataFromFabric--------------QueryDataFromFabric")
	block1, _, transactions1, _ := QueryDataFromFabric()
	fmt.Println("ADS---------------------------------------------ADS")
	adsbs(b, block1)
	adsts(t, transactions1)

}

func adsts(t transformer.Transaction, transactions1 []transformer.Transaction) {
	for _, b := range transactions1 {
		r1, err := json.Marshal(b)
		//fmt.Println(putdata)
		if err != nil {
			fmt.Errorf("Failed to json asset: %s", err)
		}
		r2, err := json.Marshal(t)
		//fmt.Println(putdata)
		if err != nil {
			fmt.Errorf("Failed to json asset: %s", err)
		}
		if string(r1) == string(r2) {
			fmt.Println(t)
			return
		}
	}
	fmt.Println("failed")
}

func adsbs(b transformer.Block, block1 transformer.Block) {
	r1, err := json.Marshal(b)
	if err != nil {
		fmt.Errorf("Failed to json asset: %s", err)
	}
	r2, err := json.Marshal(block1)
	if err != nil {
		fmt.Errorf("Failed to json asset: %s", err)
	}
	if string(r1) == string(r2) {
		fmt.Println(block1)
		return
	}
	fmt.Println("failed")
}

func pbftbs(bs []transformer.Block) transformer.Block {
	num := make(map[transformer.Block]int)
	b := transformer.Block{}
	for _, result := range bs {
		num[result]++
	}

	// 找到出现次数超过三分之一的结果
	threshold := len(bs) / 3
	for _, count := range num {
		if count > threshold {
			return bs[0]
		}
	}
	return b
}

func pbftts(bs []transformer.Transaction) transformer.Transaction {
	num := make(map[transformer.Transaction]int)
	b := transformer.Transaction{}
	for _, result := range bs {
		num[result]++
	}

	// 找到出现次数超过三分之一的结果
	threshold := len(bs) / 3
	for _, count := range num {
		if count > threshold {
			return bs[0]
		}
	}
	return b
}

func pbftcs(bs []int) int {
	num := make(map[int]int)
	b := -1
	for _, result := range bs {
		num[result]++
	}

	// 找到出现次数超过三分之一的结果
	threshold := len(bs) / 3
	for _, count := range num {
		if count > threshold {
			return bs[0]
		}
	}
	return b
}

func pbftbls(bs []float64) float64 {
	num := make(map[float64]int)
	b := 0.0
	for _, result := range bs {
		num[result]++
	}

	// 找到出现次数超过三分之一的结果
	threshold := len(bs) / 3
	for _, count := range num {
		if count > threshold {
			return bs[0]
		}
	}
	return b
}

func pbfttrs(bs [][]transformer.Transaction) []transformer.Transaction {
	var lens []int
	for _, b := range bs {
		l := len(b)
		lens = append(lens, l)
	}
	num := make(map[int]int)
	for _, result := range lens {
		num[result]++
	}

	// 找到出现次数超过三分之一的结果
	threshold := len(bs) / 3
	for _, count := range num {
		if count > threshold {
			nums := make(map[string]int)
			for _, result := range bs {
				r, err := json.Marshal(result)
				//fmt.Println(putdata)
				if err != nil {
					fmt.Errorf("Failed to json asset: %s", err)
				}
				nums[string(r)]++
			}
			for _, counts := range num {
				if counts > threshold {
					return bs[0]
				}
			}
		}
	}
	return nil
}

func pbfttras(bs [][]transformer.Transaction) []transformer.Transaction {
	var lens []int
	for _, b := range bs {
		l := len(b)
		lens = append(lens, l)
	}
	num := make(map[int]int)
	for _, result := range lens {
		num[result]++
	}
	// 找到出现次数超过三分之一的结果
	threshold := len(bs) / 3
	for _, count := range num {
		if count > threshold {
			nums := make(map[string]int)
			for _, result := range bs {
				r, err := json.Marshal(result)
				//fmt.Println(putdata)
				if err != nil {
					fmt.Errorf("Failed to json asset: %s", err)
				}
				nums[string(r)]++
			}
			for _, counts := range num {
				if counts > threshold {
					return bs[0]
				}
			}
		}
	}
	return nil
}

func InsertDataToFabric(block transformer.Block, accounts []transformer.Account, transactions []transformer.Transaction, operations []transformer.Operation) error {
	putdata, err := json.Marshal(block)
	//fmt.Println(putdata)
	if err != nil {
		return fmt.Errorf("Failed to json asset: %s", err)
	}
	opera := []string{"set", "Block", string(putdata)}
	_, err = App.Set(opera)
	if err != nil {
		return fmt.Errorf("addu failed, err:%v\n", err)
	}
	putdata, err = json.Marshal(accounts)
	//fmt.Println(putdata)
	if err != nil {
		return fmt.Errorf("Failed to json asset: %s", err)
	}
	opera = []string{"set", "Account", string(putdata)}
	_, err = App.Set(opera)
	if err != nil {
		return fmt.Errorf("addu failed, err:%v\n", err)
	}
	putdata, err = json.Marshal(transactions)
	//fmt.Println(putdata)
	if err != nil {
		return fmt.Errorf("Failed to json asset: %s", err)
	}
	opera = []string{"set", "Transaction", string(putdata)}
	_, err = App.Set(opera)
	if err != nil {
		return fmt.Errorf("addu failed, err:%v\n", err)
	}
	putdata, err = json.Marshal(operations)
	//fmt.Println(putdata)
	if err != nil {
		return fmt.Errorf("Failed to json asset: %s", err)
	}
	opera = []string{"set", "Operation", string(putdata)}
	_, err = App.Set(opera)
	if err != nil {
		return fmt.Errorf("addu failed, err:%v\n", err)
	}
	return nil
}

func InsertDataToDatabases(block transformer.Block, accounts []transformer.Account, transactions []transformer.Transaction, operations []transformer.Operation) error {
	databases := []string{"test1", "test2", "test3", "test4"}
	for _, dbname := range databases {
		db, err := sql.Open("mysql", "root:abc123@tcp(172.22.84.3:3306)/"+dbname)
		if err != nil {
			fmt.Println(err)
		}
		defer db.Close()

		err = transformer.InsertData(db, block, transactions, accounts, operations)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func QueryDataFromFabric() (transformer.Block, []transformer.Account, []transformer.Transaction, []transformer.Operation) {
	var block transformer.Block
	var accounts []transformer.Account
	var transactions []transformer.Transaction
	var operations []transformer.Operation
	opera := []string{"get", "Block"}
	value, err := App.Get(opera)
	if err != nil {
		fmt.Errorf("read failed, err:%v\n", err)
	}
	err = json.Unmarshal([]byte(value), &block)
	if err != nil {
		fmt.Println(">> Unmarshal error: ", err)
	}
	opera = []string{"get", "Account"}
	value, err = App.Get(opera)
	if err != nil {
		fmt.Errorf("read failed, err:%v\n", err)
	}
	err = json.Unmarshal([]byte(value), &accounts)
	if err != nil {
		fmt.Println(">> Unmarshal error: ", err)
	}
	opera = []string{"get", "Transaction"}
	value, err = App.Get(opera)
	if err != nil {
		fmt.Errorf("read failed, err:%v\n", err)
	}
	err = json.Unmarshal([]byte(value), &transactions)
	if err != nil {
		fmt.Println(">> Unmarshal error: ", err)
	}
	opera = []string{"get", "Operation"}
	value, err = App.Get(opera)
	if err != nil {
		fmt.Errorf("read failed, err:%v\n", err)
	}
	err = json.Unmarshal([]byte(value), &operations)
	if err != nil {
		fmt.Println(">> Unmarshal error: ", err)
	}
	return block, accounts, transactions, operations
}

func QueryDataFromDatabases() ([]transformer.Block, []transformer.Transaction, []int, []float64, [][]transformer.Transaction, [][]transformer.Transaction) {
	// 执行查询操作并获取结果
	var b transformer.Block
	var bs []transformer.Block
	var t transformer.Transaction
	var ts []transformer.Transaction
	var count int
	var counts []int
	var balance float64
	var balances []float64
	var tr transformer.Transaction

	var trss [][]transformer.Transaction
	var tra transformer.Transaction

	var trass [][]transformer.Transaction
	//fmt.Println(databases)
	for _, db := range databases {
		var trs []transformer.Transaction
		var tras []transformer.Transaction
		// 1.通过前一个区块的hash来查询单个区块的信息
		blockRows, err := queries.QueryBlockByPreviousHash(db, "")
		if err != nil {
			fmt.Println("Error querying block: ", err)
		}
		for blockRows.Next() {
			err = blockRows.Scan(&b.BlockID, &b.BlockHash, &b.PreviousHash, &b.CreateTime, &b.TransactionCount)
		}
		bs = append(bs, b)

		// 2.通过交易发起方账户查询单个交易的信息
		transactionRows, err := queries.QueryTransactionBySender(db, "Alice")
		if err != nil {
			fmt.Println("Error querying transaction: ", err)
		}
		for transactionRows.Next() {
			err = transactionRows.Scan(&t.TransactionHash, &t.TransactionType, &t.CreateTime, &t.BlockID, &t.Sender, &t.Receiver, &t.Amount, &t.Memo, &t.SenderBalance, &t.ReceiverBalance)
		}
		ts = append(ts, t)

		// 3.查询创建时间在2022年1月1日至2022年12月31日之间的交易数量 (Range)
		countRows, err := queries.QueryTransactionCountIn2022(db)
		if err != nil {
			fmt.Println("Error querying transaction count: ", err)
		}

		for countRows.Next() {
			err = countRows.Scan(&count)
		}
		counts = append(counts, count)

		// 4.查询由某账户发起的交易金额大于等于10000的总交易金额
		amountRows, err := queries.QueryTotalTransactionAmountByAccount(db, "Bob")
		if err != nil {
			fmt.Println("Error querying total transaction amount: ", err)
		}
		for amountRows.Next() {
			err = amountRows.Scan(&balance)
		}
		balances = append(balances, balance)

		// 5.查询指定账户提出的所有交易信息
		AccountTransactionRows, err := queries.QueryAllTransactionsByAccount(db, "Alice")
		if err != nil {
			fmt.Println("Error querying accountTransaction: ", err)
		}
		for AccountTransactionRows.Next() {
			err = AccountTransactionRows.Scan(&tr.TransactionHash, &tr.TransactionType, &tr.CreateTime, &tr.BlockID, &tr.Sender, &tr.Receiver, &tr.Amount, &tr.Memo, &tr.SenderBalance, &tr.ReceiverBalance)
			trs = append(trs, tr)
		}
		trss = append(trss, trs)

		// 6.查询在2022年1月1日至2022年12月31日之间，指定用户提出的某个交易类型的全部交易信息
		AccountTypeTransactionRows, err := queries.QueryTransactionsByAccountAndTypeIn2022(db, "Alice", "transfer")
		if err != nil {
			fmt.Println("Error querying accountTypeTransaction: ", err)
		}
		for AccountTypeTransactionRows.Next() {
			err = AccountTypeTransactionRows.Scan(&tra.TransactionHash, &tra.TransactionType, &tra.CreateTime, &tra.BlockID, &tra.Sender, &tra.Receiver, &tra.Amount, &tra.Memo, &tra.SenderBalance, &tra.ReceiverBalance)
			tras = append(tras, tra)
		}
		trass = append(trass, tras)
	}
	return bs, ts, counts, balances, trss, trass

}

//common
//初始化DB
func initDB() (err error) {
	// 配置数据库连接信息
	dbConfigs := []dbop.DBConfig{
		{
			User:     "root",
			Password: "abc123",
			Host:     "172.22.84.3",
			Port:     "3306",
			DBName:   "test1",
		},
		{
			User:     "root",
			Password: "abc123",
			Host:     "172.22.84.3",
			Port:     "3306",
			DBName:   "test2",
		},
		{
			User:     "root",
			Password: "abc123",
			Host:     "172.22.84.3",
			Port:     "3306",
			DBName:   "test3",
		},
		{
			User:     "root",
			Password: "abc123",
			Host:     "172.22.84.3",
			Port:     "3306",
			DBName:   "test4",
		},
	}

	// 连接到四个数据库
	databases = make([]*dbop.MySQLDB, 0)
	for _, cfg := range dbConfigs {
		db, err := dbop.NewMySQLDB(cfg)
		if err != nil {
			fmt.Println("Error connecting to database: ", err)
		}
		databases = append(databases, db)
	}
	return nil
}

func main() {
	fmt.Println("go")
	//org信息
	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/home/zhang/dqfabric/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/home/zhang/dqfabric/fixtures/channel-artifacts/Org2MSPanchors.tx",
		},
	}
	// 初始化info
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    "/home/zhang/dqfabric/fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    "/home/zhang/dqfabric/chaincode/go/sacc",
		ChaincodeVersion: cc_version,
	}
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> Sdk set error ", err)
		os.Exit(-1)
	}
	if err := sdkInit.CreateChannel(&info); err != nil {
		fmt.Println(">> Create channel error: ", err)
		os.Exit(-1)
	}
	if err := sdkInit.JoinChannel(&info); err != nil {
		fmt.Println(">> join channel error: ", err)
		os.Exit(-1)
	}
	//chaincode operation
	packageID, err := sdkInit.InstallCC(&info)
	if err != nil {
		fmt.Println(">> install chaincode error: ", err)
		os.Exit(-1)
	}
	//apprrove
	if err := sdkInit.ApproveLifecycle(&info, 1, packageID); err != nil {
		fmt.Println(">> approve chaincode error: ", err)
		os.Exit(-1)
	}
	//init chaincode
	if err := sdkInit.InitCC(&info, false, sdk); err != nil {
		fmt.Println(">> init chaincode error: ", err)
		os.Exit(-1)
	}
	fmt.Println(">> 通过链码外部服务设置链码状态......")
	if err := info.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk); err != nil {
		fmt.Println(">> InitService error: ", err)
		os.Exit(-1)
	}
	App = sdkInit.Application{
		SdkEnvInfo: &info,
	}
	fmt.Println(">> 设置链码状态完成")

	Tests()

}

////	//cfg := transformer.DBConfig{
////	//	User:     "root",
////	//	Password: "abc123",
////	//	Host:     "172.22.84.3",
////	//	Port:     "3306",
////	//	DBName:   "test1",
////	//}
////	//
////	//db, err := transformer.NewMySQLDB(cfg)
////	//if err != nil {
////	//	fmt.Println("Error connecting to database: ", err)
////	//	return
////	//}
////	//var books Books
////	//books.Book_id = 2
////	//books.Title = "fdsfs"
////	//books.Subject = "34234"
////	//books.Author = "fdsfds"
////	//data, err := json.Marshal(books)
////	//if err != nil {
////	//	fmt.Println("序列号失败", err)
////	//}
////	//fmt.Println(string(data))
////	//// Insert operation
////	//err = db.Insert("INSERT INTO test_table VALUES ('John', 'Doe')")
////	//if err != nil {
////	//	fmt.Println("Error inserting into database: ", err)
////	//	return
////	//}
////	//
////	//// Select operation
////	//rows, err := db.Select("SELECT * FROM test_table")
////	//if err != nil {
////	//	fmt.Println("Error selecting from database: ", err)
////	//	return
////	//}
////	//defer rows.Close()
////	//
////	//// Handle selected rows here...[]interface{}
////	//
////	//// Update operation
////	//err = db.Update("UPDATE test_table SET name='Jane' WHERE surname='Doe'")
////	//if err != nil {
////	//	fmt.Println("Error updating database: ", err)
////	//	return
////	//}
////	//
////	//// Delete operation
////	//err = db.Delete("DELETE FROM test_table WHERE name='Jane'")
////	//if err != nil {
////	//	fmt.Println("Error deleting from database: ", err)
////	//	return
////	//}
////
////	databases := []string{"test1", "test2", "test3", "test4"}
////
////	block := transformer.Block{
////		BlockID:          0,
////		BlockHash:        "0x123",
////		PreviousHash:     "",
////		CreateTime:       time.Date(2022, 1, 3, 1, 0, 0, 0, time.UTC),
////		TransactionCount: 2,
////	}
////
////	transactions := []transformer.Transaction{
////		{
////			TransactionHash: "0x456",
////			TransactionType: "transfer",
////			CreateTime:      time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
////			BlockID:         0,
////			Sender:          "Alice",
////			Receiver:        "Bob",
////			Amount:          1000,
////			Memo:            "Payment for services",
////			SenderBalance:   9000,
////			ReceiverBalance: 12000,
////		},
////		{
////			TransactionHash: "0x789",
////			TransactionType: "transfer",
////			CreateTime:      time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
////			BlockID:         0,
////			Sender:          "Bob",
////			Receiver:        "Charlie",
////			Amount:          10000,
////			Memo:            "Payment for goods",
////			SenderBalance:   2000,
////			ReceiverBalance: 11000,
////		},
////	}
////
////	accounts := []transformer.Account{
////		{
////			AccountID:     "Alice",
////			AccountType:   "individual",
////			CreateTime:    time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
////			AccountStatus: "active",
////			Balance:       9000,
////			CreditScore:   700,
////			LastUpdate:    time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
////		},
////		{
////			AccountID:     "Bob",
////			AccountType:   "individual",
////			CreateTime:    time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
////			AccountStatus: "active",
////			Balance:       2000,
////			CreditScore:   700,
////			LastUpdate:    time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
////		},
////		{
////			AccountID:     "Charlie",
////			AccountType:   "individual",
////			CreateTime:    time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
////			AccountStatus: "active",
////			Balance:       11000,
////			CreditScore:   700,
////			LastUpdate:    time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
////		},
////	}
////
////	operations := []transformer.Operation{
////		{
////			OperationID:     1,
////			OperationType:   "C",
////			CreateTime:      time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
////			TransactionHash: "0x456",
////			AccountID:       "Alice",
////			Balance:         9000,
////			State:           "active",
////		},
////		{
////			OperationID:     2,
////			OperationType:   "C",
////			CreateTime:      time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
////			TransactionHash: "0x789",
////			AccountID:       "Bob",
////			Balance:         2000,
////			State:           "active",
////		},
////	}
////
////	for _, dbname := range databases {
////		db, err := sql.Open("mysql", "root:abc123@tcp(172.22.84.3:3306)/"+dbname)
////		if err != nil {
////			log.Fatal(err)
////		}
////		defer db.Close()
////
////		err = transformer.InsertData(db, block, transactions, accounts, operations)
////		if err != nil {
////			log.Fatal(err)
////		}
////	}
////}
//
//func main() {
//	// 配置数据库连接信息
//	dbConfigs := []dbop.DBConfig{
//		{
//			User:     "root",
//			Password: "abc123",
//			Host:     "192.168.43.69",
//			Port:     "3306",
//			DBName:   "test1",
//		},
//		{
//			User:     "root",
//			Password: "abc123",
//			Host:     "192.168.43.69",
//			Port:     "3306",
//			DBName:   "test2",
//		},
//		{
//			User:     "root",
//			Password: "abc123",
//			Host:     "192.168.43.69",
//			Port:     "3306",
//			DBName:   "test3",
//		},
//		{
//			User:     "root",
//			Password: "abc123",
//			Host:     "192.168.43.69",
//			Port:     "3306",
//			DBName:   "test4",
//		},
//	}
//
//	// 连接到四个数据库
//	databases := make([]*dbop.MySQLDB, 0)
//	for _, cfg := range dbConfigs {
//		db, err := dbop.NewMySQLDB(cfg)
//		if err != nil {
//			fmt.Println("Error connecting to database: ", err)
//			return
//		}
//		defer db.Close()
//		databases = append(databases, db)
//	}
//
//	// 执行查询操作并获取结果
//	blockResults := []interface{}{}
//	transactionResults := []interface{}{}
//	countResults := []interface{}{}
//	amountResults := []interface{}{}
//	var b transformer.Block
//	var bs []transformer.Block
//	for _, db := range databases {
//		// 1.通过前一个区块的hash来查询单个区块的信息
//		blockRows, err := queries.QueryBlockByPreviousHash(db, "")
//		if err != nil {
//			fmt.Println("Error querying block: ", err)
//			return
//		}
//		for blockRows.Next() {
//			err = blockRows.Scan(&b.BlockID, &b.BlockHash, &b.PreviousHash, &b.CreateTime, &b.TransactionCount)
//		}
//		bs = append(bs, b)
//		//for k,v := range bs{
//		//	fmt.Println(k,v)
//		//}
//
//		blockResults = append(blockResults, blockRows)
//
//		// 2.通过交易发起方账户查询单个交易的信息
//		transactionRows, err := queries.QueryTransactionBySender(db, "Alice")
//		if err != nil {
//			fmt.Println("Error querying transaction: ", err)
//			return
//		}
//		transactionResults = append(transactionResults, transactionRows)
//
//		// 3.查询创建时间在2022年1月1日至2022年12月31日之间的交易数量
//		countRows, err := queries.QueryTransactionCountIn2022(db)
//		if err != nil {
//			fmt.Println("Error querying transaction count: ", err)
//			return
//		}
//		countResults = append(countResults, countRows)
//
//		// 4.查询由某账户发起的交易金额大于等于10000的总交易金额
//		amountRows, err := queries.QueryTotalTransactionAmountByAccount(db, "Bob")
//		if err != nil {
//			fmt.Println("Error querying total transaction amount: ", err)
//			return
//		}
//		amountResults = append(amountResults, amountRows)
//	}
//	fmt.Println(bs)
//
//	if pbft.VerifyBlock(bs) {
//		fmt.Println("Verification succeeded.")
//	} else {
//		fmt.Println("Verification failed.")
//	}
//	if pbft.VerifyResults(transactionResults) {
//		fmt.Println("Verification succeeded.")
//	} else {
//		fmt.Println("Verification failed.")
//	}
//	if pbft.VerifyResults(countResults) {
//		fmt.Println("Verification succeeded.")
//	} else {
//		fmt.Println("Verification failed.")
//	}
//	if pbft.VerifyResults(amountResults) {
//		fmt.Println("Verification succeeded.")
//	} else {
//		fmt.Println("Verification failed.")
//	}
//
//	//// 输出countResults的查询结果
//	//for i := 0; i < len(countResults); i++ {
//	//	fmt.Println("Results for count query from database", i+1)
//	//	rows := countResults[i]
//	//	for rows.Next() {
//	//		var count int
//	//		err := rows.Scan(&count)
//	//		if err != nil {
//	//			fmt.Println("Error reading count row: ", err)
//	//			return
//	//		}
//	//		fmt.Println("Count:", count)
//	//	}
//	//	fmt.Println("------")
//	//}
//	//
//	//// 服务器端验证
//	//if pbft.ValidateServer(blockResults...) {
//	//	fmt.Println("Server validation successful")
//	//} else {
//	//	fmt.Println("Server validation failed")
//	//}
//	//if pbft.ValidateServer(transactionResults...) {
//	//	fmt.Println("Server validation successful")
//	//} else {
//	//	fmt.Println("Server validation failed")
//	//}
//	//if pbft.ValidateServer(countResults...) {
//	//	fmt.Println("Server validation successful")
//	//} else {
//	//	fmt.Println("Server validation failed")
//	//}
//	//if pbft.ValidateServer(amountResults...) {
//	//	fmt.Println("Server validation successful")
//	//} else {
//	//	fmt.Println("Server validation failed")
//	//}
//
//	//// 客户端验证
//	//expectedData := "expected_data" // 替换为您的预期数据
//	//
//	//for _, rows := range countResults {
//	//	valid, err := ads.ValidateClient(rows, expectedData)
//	//	if err != nil {
//	//		fmt.Println("Error validating client data: ", err)
//	//		return
//	//	}
//	//
//	//	if valid {
//	//		fmt.Println("Client validation successful")
//	//	} else {
//	//		fmt.Println("Client validation failed")
//	//	}
//	//}
//
//}
