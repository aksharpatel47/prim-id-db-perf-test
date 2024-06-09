package main

import (
	"aksharpatel47.com/prim-id-perf-test/models"
	"aksharpatel47.com/prim-id-perf-test/repositories"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"math/rand"
	"os"
	"time"
)

func insertData(conn *pgx.Conn) {
	dayRandomIntGenFunc := func() func() int64 {
		d := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		dayCounter := 0

		return func() int64 {
			id := repositories.GenerateRandomIntOfLength18WithDatePrefix(d)
			dayCounter++
			if dayCounter%5000 == 0 {
				d = d.Add(time.Hour * 24)
			}
			return id
		}
	}

	gen := dayRandomIntGenFunc()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	start := time.Now()

	columnsWithId := []string{"id", "first_name", "last_name", "data_1", "data_2", "data_3", "data_4", "data_5"}
	columnsWithoutId := columnsWithId[1:]

	copyCount, err := repositories.InsertPeopleData(conn, repositories.PersonTableNameBigserial, columnsWithoutId, repositories.GeneratePeopleInputData[int64](nil))

	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted", copyCount, "people rows with bigserial prefix")

	bigSerialPeopleIds, err := repositories.GetAllPersonIds[int64](conn, repositories.PersonTableNameBigserial)

	rnd.Shuffle(len(bigSerialPeopleIds), func(i, j int) {
		bigSerialPeopleIds[i], bigSerialPeopleIds[j] = bigSerialPeopleIds[j], bigSerialPeopleIds[i]
	})

	addressCopyCount, err := repositories.InsertAddressData(conn, repositories.AddressTableNameBigserial, repositories.GenerateAddressInputData[int64](bigSerialPeopleIds))

	end := time.Now()

	fmt.Println("Inserted", addressCopyCount, "address rows with bigserial prefix")
	fmt.Println("Time taken", end.Sub(start))

	start = time.Now()

	copyCount, err = repositories.InsertPeopleData(conn, repositories.PersonTableNameUUIDV7, columnsWithId, repositories.GeneratePeopleInputData[string](repositories.GenerateUUIDV7))

	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted", copyCount, "people rows with UUID v7")

	uuidV7PeopleIds, err := repositories.GetAllPersonIds[string](conn, repositories.PersonTableNameUUIDV7)

	rnd.Shuffle(len(uuidV7PeopleIds), func(i, j int) {
		uuidV7PeopleIds[i], uuidV7PeopleIds[j] = uuidV7PeopleIds[j], uuidV7PeopleIds[i]
	})

	addressCopyCount, err = repositories.InsertAddressData(conn, repositories.AddressTableNameUUIDV7, repositories.GenerateAddressInputData[string](uuidV7PeopleIds))

	if err != nil {
		panic(err)
	}

	end = time.Now()

	fmt.Println("Inserted", addressCopyCount, "rows with UUID v7")
	fmt.Println("Time taken", end.Sub(start))

	start = time.Now()

	copyCount, err = repositories.InsertPeopleData(conn, repositories.PersonTableNameRandomInt, columnsWithId, repositories.GeneratePeopleInputData(repositories.GenerateRandomIntOfLength18))

	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted", copyCount, "people rows with random int prefix")

	randomIntPeopleIds, err := repositories.GetAllPersonIds[int64](conn, repositories.PersonTableNameRandomInt)

	rnd.Shuffle(len(randomIntPeopleIds), func(i, j int) {
		randomIntPeopleIds[i], randomIntPeopleIds[j] = randomIntPeopleIds[j], randomIntPeopleIds[i]
	})

	addressCopyCount, err = repositories.InsertAddressData(conn, repositories.AddressTableNameRandomInt, repositories.GenerateAddressInputData[int64](randomIntPeopleIds))

	if err != nil {
		panic(err)
	}

	end = time.Now()

	fmt.Println("Inserted", addressCopyCount, "address rows with random int prefix")
	fmt.Println("Time taken", end.Sub(start))

	start = time.Now()

	copyCount, err = repositories.InsertPeopleData(conn, repositories.PersonTableNameDateRandomInt, columnsWithId, repositories.GeneratePeopleInputData(gen))

	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted", copyCount, "people rows with date prefix")

	dateRandomIntPeopleIds, err := repositories.GetAllPersonIds[int64](conn, repositories.PersonTableNameDateRandomInt)

	rnd.Shuffle(len(dateRandomIntPeopleIds), func(i, j int) {
		dateRandomIntPeopleIds[i], dateRandomIntPeopleIds[j] = dateRandomIntPeopleIds[j], dateRandomIntPeopleIds[i]
	})

	addressCopyCount, err = repositories.InsertAddressData(conn, repositories.AddressTableNameDateRandomInt, repositories.GenerateAddressInputData[int64](dateRandomIntPeopleIds))

	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted", addressCopyCount, "address rows with date prefix")

	end = time.Now()
	fmt.Println("Time taken", end.Sub(start))

	start = time.Now()

	copyCount, err = repositories.InsertPeopleData(conn, repositories.PersonTableNameUUID, columnsWithoutId, repositories.GeneratePeopleInputData[string](nil))

	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted", copyCount, "people rows with UUID")

	uuidPeopleIds, err := repositories.GetAllPersonIds[string](conn, repositories.PersonTableNameUUID)

	rnd.Shuffle(len(uuidPeopleIds), func(i, j int) {
		uuidPeopleIds[i], uuidPeopleIds[j] = uuidPeopleIds[j], uuidPeopleIds[i]
	})

	addressCopyCount, err = repositories.InsertAddressData(conn, repositories.AddressTableNameUUID, repositories.GenerateAddressInputData[string](uuidPeopleIds))

	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted", addressCopyCount, "address rows with UUID")

	end = time.Now()
	fmt.Println("Time taken", end.Sub(start))

}

type queryBenchmark struct {
	personTableName  repositories.PersonTableName
	addressTableName repositories.AddressTableName
}

func getPersonIdBatches[T models.PersonID](ids []T, batchSize int) [][]T {
	batches := make([][]T, 0)

	for i := 0; i < len(ids); i += batchSize {
		batch := ids[i:min(i+batchSize, len(ids))]
		batches = append(batches, batch)
	}

	return batches
}

func benchmarkBatchRandomIdQueries(conn *pgx.Conn) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	personIds, err := repositories.GetAllPersonIds[int64](conn, repositories.PersonTableNameRandomInt)

	if err != nil {
		panic(err)
	}

	rnd.Shuffle(len(personIds), func(i, j int) {
		personIds[i], personIds[j] = personIds[j], personIds[i]
	})

	batches := getPersonIdBatches(personIds, 100)

	start := time.Now()
	batchesProcessed := 0

	for _, batch := range batches {
		_, err := repositories.GetPeopleWithAddresses(conn, repositories.PersonTableNameRandomInt, repositories.AddressTableNameRandomInt, batch)
		if err != nil {
			panic(err)
		}

		batchesProcessed += 1
	}

	end := time.Now()

	fmt.Println("Finished processing random int in", end.Sub(start))

	personIds, err = repositories.GetAllPersonIds[int64](conn, repositories.PersonTableNameDateRandomInt)

	if err != nil {
		panic(err)
	}

	rnd.Shuffle(len(personIds), func(i, j int) {
		personIds[i], personIds[j] = personIds[j], personIds[i]
	})

	batches = getPersonIdBatches(personIds, 100)

	start = time.Now()

	batchesProcessed = 0

	for _, batch := range batches {
		_, err := repositories.GetPeopleWithAddresses(conn, repositories.PersonTableNameDateRandomInt, repositories.AddressTableNameDateRandomInt, batch)
		if err != nil {
			panic(err)
		}

		batchesProcessed += 1
	}

	end = time.Now()

	fmt.Println("Finished processing date random int in", end.Sub(start))

	stringPersonIds, err := repositories.GetAllPersonIds[string](conn, repositories.PersonTableNameUUID)

	if err != nil {
		panic(err)
	}

	rnd.Shuffle(len(stringPersonIds), func(i, j int) {
		stringPersonIds[i], stringPersonIds[j] = stringPersonIds[j], stringPersonIds[i]
	})

	stringBatches := getPersonIdBatches(stringPersonIds, 100)

	start = time.Now()

	batchesProcessed = 0

	for _, batch := range stringBatches {
		_, err := repositories.GetPeopleWithAddresses(conn, repositories.PersonTableNameUUID, repositories.AddressTableNameUUID, batch)
		if err != nil {
			panic(err)
		}

		batchesProcessed += 1
	}

	end = time.Now()

	fmt.Println("Finished processing UUID in", end.Sub(start))

	stringPersonIds, err = repositories.GetAllPersonIds[string](conn, repositories.PersonTableNameUUIDV7)

	if err != nil {
		panic(err)
	}

	rnd.Shuffle(len(stringPersonIds), func(i, j int) {
		stringPersonIds[i], stringPersonIds[j] = stringPersonIds[j], stringPersonIds[i]
	})

	stringBatches = getPersonIdBatches(stringPersonIds, 100)

	start = time.Now()
	batchesProcessed = 0

	for _, batch := range stringBatches {
		_, err := repositories.GetPeopleWithAddresses(conn, repositories.PersonTableNameUUIDV7, repositories.AddressTableNameUUIDV7, batch)
		if err != nil {
			panic(err)
		}

		batchesProcessed += 1
	}

	end = time.Now()

	fmt.Println("Finished processing UUID v7 in", end.Sub(start))

	personIds, err = repositories.GetAllPersonIds[int64](conn, repositories.PersonTableNameBigserial)

	if err != nil {
		panic(err)
	}

	rnd.Shuffle(len(personIds), func(i, j int) {
		personIds[i], personIds[j] = personIds[j], personIds[i]
	})

	batches = getPersonIdBatches(personIds, 100)

	start = time.Now()
	batchesProcessed = 0

	for _, batch := range batches {
		_, err := repositories.GetPeopleWithAddresses(conn, repositories.PersonTableNameBigserial, repositories.AddressTableNameBigserial, batch)
		if err != nil {
			panic(err)
		}

		batchesProcessed += 1
	}

	end = time.Now()

	fmt.Println("Finished processing bigserial in", end.Sub(start))
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	defer conn.Close(context.Background())

	insertData(conn)

	benchmarkBatchRandomIdQueries(conn)

	//insertData(conn)

	//randomIds, err := getAllIds[int64](conn, "person_random_int")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//randomIdsWithDatePrefix, err := getAllIds[int64](conn, "person_date_random_int")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//_, err = getPersonWithId[int64](conn, "person_random_int", randomIds[0])
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	//
	//rnd.Shuffle(len(randomIds), func(i, j int) {
	//	randomIds[i], randomIds[j] = randomIds[j], randomIds[i]
	//})
	//
	//rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	//
	//rnd.Shuffle(len(randomIdsWithDatePrefix), func(i, j int) {
	//	randomIdsWithDatePrefix[i], randomIdsWithDatePrefix[j] = randomIdsWithDatePrefix[j], randomIdsWithDatePrefix[i]
	//})
	//
	////singleRandomIdQueryPerformance, err := os.Create("single_random_id_query_performance.txt")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	////singleRandomIdWithDatePrefixQueryPerformance, err := os.Create("single_random_id_with_date_prefix_query_performance.txt")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	////start := time.Now()
	////for _, id := range randomIds {
	////	_, err := getPersonWithId[int64](conn, "person_random_int", id)
	////	if err != nil {
	////		panic(err)
	////	}
	////
	////	//_, err = singleRandomIdQueryPerformance.WriteString(fmt.Sprintf("%d,%s\n", id, executionTime))
	////	//
	////	//if err != nil {
	////	//	panic(err)
	////	//}
	////}
	////end := time.Now()
	////
	////fmt.Println("Single random id query performance", end.Sub(start))
	//
	////start = time.Now()
	////for _, id := range randomIdsWithDatePrefix {
	////	_, err := getPersonWithId[int64](conn, "person_date_random_int", id)
	////	if err != nil {
	////		panic(err)
	////	}
	////
	////	//_, err = singleRandomIdWithDatePrefixQueryPerformance.WriteString(fmt.Sprintf("%d,%s\n", id, executionTime))
	////
	////	//if err != nil {
	////	//	panic(err)
	////	//}
	////}
	////
	////end = time.Now()
	////
	////fmt.Println("Single random id with date prefix query performance", end.Sub(start))
	//
	////multipleRandomIdsQueryPerformance, err := os.Create("multiple_random_ids_query_performance.txt")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//start := time.Now()
	//batchSize := 100
	//
	//for i := 0; i < len(randomIds); i += batchSize {
	//	ids := randomIds[i:min(i+batchSize, len(randomIds))]
	//	_, err := getPeopleWithIds(conn, "person_random_int", ids)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	//_, err = multipleRandomIdsQueryPerformance.WriteString(fmt.Sprintf("%d,%s\n", len(ids), executionTime))
	//	//
	//	//if err != nil {
	//	//	panic(err)
	//	//}
	//}
	//end := time.Now()
	//
	//fmt.Println("Multiple random ids query performance", end.Sub(start))
	//
	////multipleRandomIdsWithDatePrefixQueryPerformance, err := os.Create("multiple_random_ids_with_date_prefix_query_performance.txt")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//start = time.Now()
	//
	//for i := 0; i < len(randomIdsWithDatePrefix); i += batchSize {
	//	ids := randomIdsWithDatePrefix[i:min(i+batchSize, len(randomIdsWithDatePrefix))]
	//	_, err := getPeopleWithIds(conn, "person_date_random_int", ids)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	//_, err = multipleRandomIdsWithDatePrefixQueryPerformance.WriteString(fmt.Sprintf("%d,%s\n", len(ids), executionTime))
	//	//
	//	//if err != nil {
	//	//	panic(err)
	//	//}
	//}
	//
	//end = time.Now()
	//
	//fmt.Println("Multiple random ids with date prefix query performance", end.Sub(start))
}
