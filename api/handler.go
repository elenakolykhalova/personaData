package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	database "personaData/db"
	"personaData/models"
)

// Обогащение возрастом
func fetchAge(name string) (int, error) {
	logrus.Debug("Fetching age for name: ", name)
	resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		logrus.Error("Error fetching age: ", err)
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	age := int(result["age"].(float64))
	logrus.Debug("Fetched age: ", age)
	return age, nil
}

// Обогащение полом
func fetchGender(name string) (string, error) {
	logrus.Debug("Fetching gender for name: ", name)
	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		logrus.Error("Error fetching gender: ", err)
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	gender := result["gender"].(string)
	logrus.Debug("Fetched gender: ", gender)
	return gender, nil
}

// Обогащение национальностью
func fetchNationality(name string) (string, error) {
	logrus.Debug("Fetching nationality for name: ", name)
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		logrus.Error("Error fetching nationality: ", err)
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	country := result["country"].([]interface{})[0].(map[string]interface{})
	nationality := country["country_id"].(string)
	logrus.Debug("Fetched nationality: ", nationality)
	return nationality, nil
}

// Добавление человека в таблицу
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	// Берем название таблицы из .env
	tableName := os.Getenv("PG_DATABASE_NAME")
	if tableName == "" {
		http.Error(w, "TABLE_NAME environment variable is not set", http.StatusInternalServerError)
		return
	}

	// Извлечение query-параметров
	query := r.URL.Query()
	name := query.Get("name")
	surname := query.Get("surname")
	patronymic := query.Get("patronymic")

	// Проверка обязательных параметров
	if name == "" || surname == "" {
		http.Error(w, "name and surname are required", http.StatusBadRequest)
		return
	}

	var person models.Person
	person.Name = name
	person.Surname = surname
	person.Patronymic = patronymic

	// Добавление обогащения в структуру person
	var err error
	person.Age, err = fetchAge(person.Name)
	if err != nil {
		logrus.Warn("Could not fetch age for person: ", err)
	}
	person.Gender, err = fetchGender(person.Name)
	if err != nil {
		logrus.Warn("Could not fetch gender for person: ", err)
	}
	person.Nationality, err = fetchNationality(person.Name)
	if err != nil {
		logrus.Warn("Could not fetch nationality for person: ", err)
	}

	// Сохранение в базе данных
	db := database.GetDB()
	_, err = db.Exec(context.Background(),
		fmt.Sprintf("INSERT INTO %s (name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6)", tableName),
		person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality)
	if err != nil {
		logrus.Error("Error inserting person into database: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Info("Person created: ", person)
	json.NewEncoder(w).Encode(&person)
}
