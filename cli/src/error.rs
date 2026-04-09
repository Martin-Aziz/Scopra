use thiserror::Error;

#[derive(Debug, Error)]
pub enum CliError {
    #[error("io error: {0}")]
    Io(#[from] std::io::Error),

    #[error("http error: {0}")]
    Http(#[from] reqwest::Error),

    #[error("json error: {0}")]
    Json(#[from] serde_json::Error),

    #[error("toml de/serialize error: {0}")]
    Toml(#[from] toml::de::Error),

    #[error("toml serialize error: {0}")]
    TomlSerialize(#[from] toml::ser::Error),

    #[error("configuration error: {0}")]
    Config(String),

    #[error("api error: {0}")]
    Api(String),

    #[error("validation error: {0}")]
    Validation(String),
}

pub type Result<T> = std::result::Result<T, CliError>;
