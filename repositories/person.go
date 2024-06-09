package repositories

import (
	"aksharpatel47.com/prim-id-perf-test/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"math/rand"
	"strconv"
	"time"
)

type PersonTableName string

const (
	PersonTableNameBigserial     PersonTableName = "person_bigserial"
	PersonTableNameRandomInt     PersonTableName = "person_random_int"
	PersonTableNameDateRandomInt PersonTableName = "person_date_random_int"
	PersonTableNameUUID          PersonTableName = "person_uuid"
	PersonTableNameUUIDV7        PersonTableName = "person_uuidv7"
)

type AddressTableName string

const (
	AddressTableNameBigserial     AddressTableName = "addresses_bigserial"
	AddressTableNameRandomInt     AddressTableName = "addresses_random_int"
	AddressTableNameDateRandomInt AddressTableName = "addresses_date_random_int"
	AddressTableNameUUID          AddressTableName = "addresses_uuid"
	AddressTableNameUUIDV7        AddressTableName = "addresses_uuidv7"
)

func GetAllPersonIds[T models.PersonID](conn *pgx.Conn, table PersonTableName) ([]T, error) {
	ids := make([]T, 0)

	rows, err := conn.Query(context.Background(), fmt.Sprintf("SELECT id FROM %s", table))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id T
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func GetPeopleWithAddresses[T models.PersonID](conn *pgx.Conn, personTable PersonTableName, addressTable AddressTableName, peopleIds []T) ([]models.Person[T], error) {
	people := make([]models.Person[T], 0)

	query := fmt.Sprintf(`
select p.id, p.first_name, p.last_name, p.data_1, p.data_2, p.data_3, p.data_4, p.data_5, jsonb_agg(jsonb_build_object('id', a.id, 'person_id', a.person_id, 'address', a.address, 'city', a.city, 'state', a.state, 'zip', a.zip)) as addresses
from %s p
left join %s a on p.id = a.person_id
where p.id = ANY($1)
group by p.id
`, personTable, addressTable)

	rows, err := conn.Query(context.Background(), query, peopleIds)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var person models.Person[T]
		err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &person.Data1, &person.Data2, &person.Data3, &person.Data4, &person.Data5, &person.Addresses)

		if err != nil {
			return nil, err
		}

		people = append(people, person)
	}

	return people, nil
}

func GenerateRandomIntOfLength18() int64 {
	return rand.Int63n(1e18)
}

func GenerateRandomIntOfLength18WithDatePrefix(date time.Time) int64 {
	intStr := fmt.Sprintf("%s%s%s%d", fmt.Sprintf("%d", date.Year())[2:], fmt.Sprintf("%02d", date.Month()), fmt.Sprintf("%02d", date.Day()), rand.Int63n(1e12))
	intVal, _ := strconv.ParseInt(intStr, 10, 64)

	return intVal
}

func GenerateUUIDV7() string {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return id.String()
}

type IdGenFunc[T models.PersonID] func() T

func GeneratePeopleInputData[T models.PersonID](idGenFunc IdGenFunc[T]) [][]any {
	fmt.Println("Generating input data")
	var c int = 2 * 1e6
	idSet := make(map[T]struct{})
	data := make([][]any, 0)

	// add column names to data

	for i := 0; i < c; i++ {
		if idGenFunc == nil {
			data = append(data, []any{
				"first_name",
				"last_name",
				"data1",
				"data2",
				"data3",
				"data4",
				"data5",
			})
			continue
		}

		id := idGenFunc()
		for _, ok := idSet[id]; ok; {
			fmt.Println("Duplicate id found, regenerating", id)
			id = idGenFunc()
		}
		data = append(data, []any{
			id,
			"first_name",
			"last_name",
			"data1",
			"data2",
			"data3",
			"data4",
			"data5",
		})
	}

	return data
}

func GenerateAddressInputData[T models.PersonID](personIds []T) [][]any {
	fmt.Println("Generating address input data")
	data := make([][]any, 0)

	for _, personId := range personIds {
		data = append(data, []any{
			personId,
			"street 1",
			"city 1",
			"state 1",
			"zip 1",
		})

		data = append(data, []any{
			personId,
			"street 2",
			"city 2",
			"state 2",
			"zip 2",
		})
	}

	return data
}

func InsertPeopleData(conn *pgx.Conn, tableName PersonTableName, columns []string, data [][]any) (int, error) {
	fmt.Println("Inserting data")
	copyCount, err := conn.CopyFrom(context.Background(), pgx.Identifier{string(tableName)}, columns, pgx.CopyFromRows(data))
	return int(copyCount), err
}

func InsertAddressData(conn *pgx.Conn, tableName AddressTableName, data [][]any) (int, error) {
	fmt.Println("Inserting data")
	copyCount, err := conn.CopyFrom(context.Background(), pgx.Identifier{string(tableName)}, []string{"person_id", "address", "city", "state", "zip"}, pgx.CopyFromRows(data))
	return int(copyCount), err
}
