package services

import (
	"antibomberman/mego-post/internal/clients"
	"antibomberman/mego-post/internal/models"
	"antibomberman/mego-post/internal/repositories"
	"context"
	storagePb "github.com/antibomberman/mego-protos/gen/go/storage"
	"log"
)

type categoryService struct {
	categoryRepository repositories.CategoryRepository
	storageClient      *clients.StorageClient
}

func NewCategoryService(catRepo repositories.CategoryRepository, storageClient *clients.StorageClient) CategoryService {
	return &categoryService{
		categoryRepository: catRepo,
		storageClient:      storageClient,
	}
}

func (p *categoryService) Find() ([]models.CategoryDetails, error) {

	categories, err := p.categoryRepository.Find()
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		return nil, err
	}

	CategoryDetails := make([]models.CategoryDetails, len(categories))
	for i, category := range categories {
		img := models.File{}

		if category.Icon != "" {
			rsp, err := p.storageClient.GetObjectUrl(context.Background(), &storagePb.GetObjectUrlRequest{
				FileName: category.Icon,
			})
			if err == nil {
				img.FileName = rsp.FileName
				img.ContentType = rsp.ContentType
				img.Url = rsp.Url
			}
		}
		CategoryDetails[i] = models.CategoryDetails{
			Id:   category.Id,
			Name: category.Name,
			Icon: &img,
		}
	}
	return CategoryDetails, nil
}
func (p *categoryService) GetById(id string) (*models.CategoryDetails, error) {
	category, err := p.categoryRepository.GetById(id)
	if err != nil {
		return nil, err
	}
	img := models.File{}
	if category.Icon != "" {
		rsp, err := p.storageClient.GetObjectUrl(context.Background(), &storagePb.GetObjectUrlRequest{
			FileName: category.Icon,
		})
		if err == nil {
			img.FileName = rsp.FileName
			img.ContentType = rsp.ContentType
			img.Url = rsp.Url
		}
	}
	return &models.CategoryDetails{
		Id:   category.Id,
		Name: category.Name,
		Icon: &img,
	}, nil
}

func (p *categoryService) Create(data models.CategoryCreate) (*models.CategoryDetails, error) {
	//upload main image
	if data.Icon != nil && data.Icon.FileName != "" {
		log.Printf("Uploading main image: %s", data.Icon.FileName)

		rsp, err := p.storageClient.PutObject(context.Background(), &storagePb.PutObjectRequest{
			FileName:    data.Icon.FileName,
			ContentType: data.Icon.ContentType,
			Data:        data.Icon.Data,
		})
		if err != nil {
			log.Printf("Error uploading main image: %v", err)
			return nil, err
		}
		log.Println(rsp)
		data.Icon.FileName = rsp.FileName
	}

	catId, err := p.categoryRepository.Create(data.Name, data.Icon.FileName)
	if err != nil {
		return nil, err
	}

	return p.GetById(catId)
}
func (p *categoryService) Update(data models.CategoryUpdate) (*models.CategoryDetails, error) {
	if data.Icon.FileName != "" {
		err := p.DeleteIcon(data.Id)
		if err != nil {
			return nil, err
		}
		rsp, err := p.storageClient.PutObject(context.Background(), &storagePb.PutObjectRequest{
			FileName:    data.Icon.FileName,
			ContentType: data.Icon.ContentType,
			Data:        data.Icon.Data,
		})
		if err != nil {
			return nil, err
		}
		data.Icon.FileName = rsp.FileName
	}
	log.Printf("Updating category: %s", data.Id)
	log.Printf("Icon: %s", data.Icon.FileName)
	log.Printf("Name: %s", data.Name)
	err := p.categoryRepository.Update(data.Id, data.Name, data.Icon.FileName)
	if err != nil {
		return nil, err
	}

	return p.GetById(data.Id)
}

func (p *categoryService) Delete(id string) error {
	err := p.DeleteIcon(id)
	if err != nil {
		return err
	}
	err = p.categoryRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
func (p *categoryService) DeleteIcon(id string) error {
	category, err := p.categoryRepository.GetById(id)
	if err != nil {
		return err
	}
	if category.Icon == "" {
		return nil
	}
	_, err = p.storageClient.DeleteObject(context.Background(), &storagePb.DeleteObjectRequest{
		FileName: category.Icon,
	})
	if err != nil {
		log.Printf("Error deleting storage object %s: %v", category.Icon, err)
		return err
	}
	return nil

}
