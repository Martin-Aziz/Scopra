use crate::error::{CliError, Result};
use serde::{Deserialize, Serialize};
use std::fs;
use std::path::{Path, PathBuf};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AppConfig {
    pub gateway_url: String,
    pub access_token: Option<String>,
}

impl Default for AppConfig {
    fn default() -> Self {
        Self {
            gateway_url: "http://localhost:8080".to_string(),
            access_token: None,
        }
    }
}

impl AppConfig {
    pub fn default_path() -> Result<PathBuf> {
        let home = dirs::home_dir()
            .ok_or_else(|| CliError::Config("unable to determine home directory".to_string()))?;
        Ok(home.join(".nexus").join("config.toml"))
    }

    pub fn load_or_default() -> Result<Self> {
        let path = Self::default_path()?;
        if !path.exists() {
            return Ok(Self::default());
        }
        Self::load_from_path(&path)
    }

    pub fn load_from_path(path: &Path) -> Result<Self> {
        let content = fs::read_to_string(path)?;
        let parsed: AppConfig = toml::from_str(&content)?;
        Ok(parsed)
    }

    pub fn save(&self) -> Result<()> {
        let path = Self::default_path()?;
        self.save_to_path(&path)
    }

    pub fn save_to_path(&self, path: &Path) -> Result<()> {
        if let Some(parent) = path.parent() {
            fs::create_dir_all(parent)?;
        }
        let encoded = toml::to_string_pretty(self)?;
        fs::write(path, encoded)?;
        Ok(())
    }
}
