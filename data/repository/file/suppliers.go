package file

import (
	"encoding/json"
	"fmt"
	"http/data/model"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type SuppliersRepository struct {
	Directory string `json:"directory"`
}

func NewSuppliersRepository() *SuppliersRepository {
	return &SuppliersRepository{
		Directory: "data/suppliers",
	}
}

func (r *SuppliersRepository) GetAll() ([]*model.Supplier, error) {
	suppliers := make([]*model.Supplier, 0)

	files, err := ioutil.ReadDir(r.Directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filepath := filepath.Join(r.Directory, file.Name())

			data, err := ioutil.ReadFile(filepath)
			if err != nil {
				return nil, err
			}

			var supplier model.Supplier
			err = json.Unmarshal(data, &supplier)
			if err != nil {
				return nil, err
			}

			suppliers = append(suppliers, &supplier)
		}

	}

	return suppliers, nil
}

func (r *SuppliersRepository) GetById(id int) (model.Supplier, error) {
	filepath := filepath.Join(r.Directory, strconv.Itoa(id)+".json")
	var supplier model.Supplier
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return supplier, err
	}

	err = json.Unmarshal(data, &supplier)

	if err != nil {
		return supplier, err
	}

	return supplier, nil
}

func (r *SuppliersRepository) UpdateSupplierById(id int, newSup model.Supplier) error {
	filepath := filepath.Join(r.Directory, strconv.Itoa(id)+".json")

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	newSup.Id = id

	data, err = json.Marshal(newSup)
	if err != nil {
		return err
	}

	err = os.Truncate(filepath, 0)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}

	return nil

}

func (r *SuppliersRepository) RefreshSuppliers(data []*model.Supplier) error {

	if err := os.RemoveAll(r.Directory); err != nil {
		return fmt.Errorf("failed to clean directory: %v", err)
	}

	if err := os.MkdirAll(r.Directory, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	for _, supplier := range data {
		supplierData, err := json.Marshal(supplier)
		if err != nil {
			log.Printf("Error encoding supplier: %v", err)
			return err
		}

		filePath := filepath.Join(r.Directory, strconv.Itoa(supplier.Id)+".json")

		if err := ioutil.WriteFile(filePath, supplierData, 0644); err != nil {
			log.Printf("Error writing supplier file: %v", err)
			return err
		}

		log.Printf("Supplier file created: %s", filePath)
	}

	return nil
}
