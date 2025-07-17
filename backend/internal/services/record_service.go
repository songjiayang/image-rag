package services

import (
	"fmt"
	"image-rag-backend/internal/database"
	"image-rag-backend/internal/models"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecordService struct {
	db *gorm.DB
}

func NewRecordService() *RecordService {
	return &RecordService{db: database.DB}
}

func (s *RecordService) GetDB() *gorm.DB {
	return s.db
}

func (s *RecordService) CreateRecord(name, description string) (*models.Record, error) {
	record := &models.Record{
		Name:        name,
		Description: description,
	}

	if err := s.db.Create(record).Error; err != nil {
		return nil, fmt.Errorf("failed to create record: %w", err)
	}

	return record, nil
}

func (s *RecordService) GetRecord(id uint) (*models.Record, error) {
	var record models.Record
	if err := s.db.Preload("Images").First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("record not found")
		}
		return nil, fmt.Errorf("failed to get record: %w", err)
	}
	return &record, nil
}

func (s *RecordService) GetRecords(limit, offset int) ([]models.Record, int64, error) {
	var records []models.Record
	var total int64

	s.db.Model(&models.Record{}).Count(&total)

	if err := s.db.Preload("Images").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get records: %w", err)
	}

	return records, total, nil
}

func (s *RecordService) UpdateRecord(id uint, name, description string) (*models.Record, error) {
	record, err := s.GetRecord(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		record.Name = name
	}
	if description != "" {
		record.Description = description
	}

	if err := s.db.Save(record).Error; err != nil {
		return nil, fmt.Errorf("failed to update record: %w", err)
	}

	return record, nil
}

func (s *RecordService) DeleteRecord(id uint) error {
	result := s.db.Delete(&models.Record{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete record: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

func (s *RecordService) AddImageToRecord(recordID uint, filename string, vectorID string) (*models.Image, error) {
	// Ensure record exists
	_, err := s.GetRecord(recordID)
	if err != nil {
		return nil, err
	}

	image := &models.Image{
		RecordID: recordID,
		Filename: filename,
		Path:     filepath.Join("uploads", filename),
		VectorID: vectorID,
	}

	if err := s.db.Create(image).Error; err != nil {
		return nil, fmt.Errorf("failed to add image: %w", err)
	}

	return image, nil
}

func (s *RecordService) DeleteImage(id uint) error {
	var image models.Image
	if err := s.db.First(&image, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("image not found")
		}
		return fmt.Errorf("failed to get image: %w", err)
	}

	// Delete file from filesystem
	if err := os.Remove(image.Path); err != nil && !os.IsNotExist(err) {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to delete file %s: %v\n", image.Path, err)
	}

	if err := s.db.Delete(&image, id).Error; err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}

func (s *RecordService) GetImage(id uint) (*models.Image, error) {
	var image models.Image
	if err := s.db.First(&image, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("image not found")
		}
		return nil, fmt.Errorf("failed to get image: %w", err)
	}
	return &image, nil
}

func (s *RecordService) GetImagesByRecordID(recordID uint) ([]models.Image, error) {
	var images []models.Image
	if err := s.db.Where("record_id = ?", recordID).Find(&images).Error; err != nil {
		return nil, fmt.Errorf("failed to get images: %w", err)
	}
	return images, nil
}

// ValidateImageFile validates image file extension and size
func ValidateImageFile(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowed[ext] {
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	return nil
}

// EnsureDirectoryExists creates directory if it doesn't exist
func EnsureDirectoryExists(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// DeleteImageByPath deletes image file by path
func (s *RecordService) DeleteImageByPath(path string) error {
	return os.Remove(path)
}

// FileExists checks if file exists
func (s *RecordService) FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// GenerateUniqueFilename generates a unique filename to prevent collisions
func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	name := strings.TrimSuffix(originalFilename, ext)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "..", "_")

	return fmt.Sprintf("%s_%s%s", uuid.New().String()[:8], name, ext)
}
