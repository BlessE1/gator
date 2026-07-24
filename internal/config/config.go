package config

import (
 "path/filepath"
 "os"
 "encoding/json"
)

type Config struct {
 DbUrl string `json:"db_url"`
 CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
 configFilePath, err := getConfigFilePath()

 if err != nil {
  return Config{}, err
 }

 data, err := os.ReadFile(configFilePath)
 if err != nil {
  return Config{}, err
 }

 var config Config
 if err := json.Unmarshal(data, &config); err != nil {
  return Config{}, err
 }

 return config, nil
}

func (c *Config) SetUser(username string) error {
 c.CurrentUserName = username
 err := write(*c)
 if err != nil {
  return err
 }

 return nil
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
 home, err := os.UserHomeDir()
 if err != nil {
  return "", err
 }

 return filepath.Join(home, configFileName), nil
}

func write(cfg Config) error {
 data, err := json.Marshal(cfg)
 if err != nil {
  return err
 }

 configFilePath, err := getConfigFilePath()
 if err != nil {
  return err
 }

 err = os.WriteFile(configFilePath, data, 0600)
 if err != nil {
  return err
 }

 return nil
}
