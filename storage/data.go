package storage

import (
	"encoding/json"
	"os"
	"sync"
)

var ConfigMgr = NewConfigManager("config.json")

type Config struct {
	ServerRoles map[string]string `json:"server_roles"`
}

type ConfigManager struct {
	FilePath string
	config   Config
	mu       sync.RWMutex
}

func NewConfigManager(path string) *ConfigManager {
	return &ConfigManager{
		FilePath: path,
		config:   Config{ServerRoles: make(map[string]string)},
	}
}

func (m *ConfigManager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	file, err := os.ReadFile(m.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // ファイルがない場合は空のまま進む
		}
		return err
	}

	return json.Unmarshal(file, &m.config)
}

func (m *ConfigManager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.FilePath, data, 0644)
}

func (m *ConfigManager) SetRole(roleID, serverID string) string {
	m.mu.Lock()
	m.config.ServerRoles[roleID] = serverID
	m.mu.Unlock()
	if err := m.Save(); err != nil {
		return "設定の保存に失敗しました: " + err.Error()
	}
	return "ロールID " + roleID + " をサーバーID " + serverID + " に設定しました。"
}

func (m *ConfigManager) GetServerID(roleID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, ok := m.config.ServerRoles[roleID]
	if !ok {
		return []string{}
	}
	return []string{id}
}
