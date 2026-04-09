use nexus_cli::config::app_config::AppConfig;
use tempfile::tempdir;

#[test]
fn config_round_trip_is_lossless() {
    let temp_dir = tempdir().expect("tempdir");
    let config_path = temp_dir.path().join("config.toml");

    let config = AppConfig {
        gateway_url: "http://localhost:9000".to_string(),
        access_token: Some("test-token".to_string()),
    };

    config.save_to_path(&config_path).expect("save config");
    let loaded = AppConfig::load_from_path(&config_path).expect("load config");

    assert_eq!(loaded.gateway_url, "http://localhost:9000");
    assert_eq!(loaded.access_token.as_deref(), Some("test-token"));
}
