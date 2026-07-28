package services

import "github.com/zhitu/server/internal/models"

// ---------- 子资源通用 CRUD ----------

// Educations
func (s *ProfileService) ListEducations(userID uint) ([]models.UserEducation, error) {
	var list []models.UserEducation
	err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&list).Error
	return list, err
}

func (s *ProfileService) CreateEducation(userID uint, m *models.UserEducation) (*models.UserEducation, error) {
	m.ID = 0
	m.UserID = userID
	if err := s.db.Create(m).Error; err != nil {
		return nil, err
	}
	s.touchCompletion(userID)
	return m, nil
}

func (s *ProfileService) UpdateEducation(userID, id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.UserEducation{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

func (s *ProfileService) DeleteEducation(userID, id uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserEducation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

// Works
func (s *ProfileService) ListWorks(userID uint) ([]models.UserWork, error) {
	var list []models.UserWork
	err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&list).Error
	return list, err
}

func (s *ProfileService) CreateWork(userID uint, m *models.UserWork) (*models.UserWork, error) {
	m.ID = 0
	m.UserID = userID
	if err := s.db.Create(m).Error; err != nil {
		return nil, err
	}
	s.touchCompletion(userID)
	return m, nil
}

func (s *ProfileService) UpdateWork(userID, id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.UserWork{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

func (s *ProfileService) DeleteWork(userID, id uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserWork{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

// Projects
func (s *ProfileService) ListProjects(userID uint) ([]models.UserProject, error) {
	var list []models.UserProject
	err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&list).Error
	return list, err
}

func (s *ProfileService) CreateProject(userID uint, m *models.UserProject) (*models.UserProject, error) {
	m.ID = 0
	m.UserID = userID
	if err := s.db.Create(m).Error; err != nil {
		return nil, err
	}
	s.touchCompletion(userID)
	return m, nil
}

func (s *ProfileService) UpdateProject(userID, id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.UserProject{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

func (s *ProfileService) DeleteProject(userID, id uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserProject{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

// Skills
func (s *ProfileService) ListSkills(userID uint) ([]models.UserSkill, error) {
	var list []models.UserSkill
	err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&list).Error
	return list, err
}

func (s *ProfileService) CreateSkill(userID uint, m *models.UserSkill) (*models.UserSkill, error) {
	m.ID = 0
	m.UserID = userID
	if err := s.db.Create(m).Error; err != nil {
		return nil, err
	}
	s.touchCompletion(userID)
	return m, nil
}

func (s *ProfileService) UpdateSkill(userID, id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.UserSkill{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

func (s *ProfileService) DeleteSkill(userID, id uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserSkill{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

// Honors
func (s *ProfileService) ListHonors(userID uint) ([]models.UserHonor, error) {
	var list []models.UserHonor
	err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&list).Error
	return list, err
}

func (s *ProfileService) CreateHonor(userID uint, m *models.UserHonor) (*models.UserHonor, error) {
	m.ID = 0
	m.UserID = userID
	if err := s.db.Create(m).Error; err != nil {
		return nil, err
	}
	s.touchCompletion(userID)
	return m, nil
}

func (s *ProfileService) UpdateHonor(userID, id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.UserHonor{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

func (s *ProfileService) DeleteHonor(userID, id uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserHonor{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

// Practices
func (s *ProfileService) ListPractices(userID uint) ([]models.UserPractice, error) {
	var list []models.UserPractice
	err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&list).Error
	return list, err
}

func (s *ProfileService) CreatePractice(userID uint, m *models.UserPractice) (*models.UserPractice, error) {
	m.ID = 0
	m.UserID = userID
	if err := s.db.Create(m).Error; err != nil {
		return nil, err
	}
	s.touchCompletion(userID)
	return m, nil
}

func (s *ProfileService) UpdatePractice(userID, id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.UserPractice{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

func (s *ProfileService) DeletePractice(userID, id uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserPractice{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSubResourceNotFound
	}
	s.touchCompletion(userID)
	return nil
}

// touchCompletion 重算并更新完成度
func (s *ProfileService) touchCompletion(userID uint) {
	pct := s.calcCompletion(userID)
	s.db.Model(&models.UserProfile{}).Where("user_id = ?", userID).Update("completion_pct", pct)
}
