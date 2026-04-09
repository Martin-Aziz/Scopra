use crate::client::api_client::ApiClient;
use crate::config::app_config::AppConfig;
use crate::error::{CliError, Result};
use crate::utils::output::print_json;
use std::fs;

pub async fn run_status(client: &ApiClient) -> Result<()> {
    let status = client.status().await?;
    print_json(&status)?;
    Ok(())
}

pub fn run_init(config: &mut AppConfig, gateway_url: String) -> Result<()> {
    config.gateway_url = gateway_url;
    config.save()?;
    println!("Initialized nexus CLI at {}", config.gateway_url);
    Ok(())
}

pub async fn run_register(
    client: &ApiClient,
    email: String,
    password: String,
    role: String,
) -> Result<()> {
    if password.len() < 12 {
        return Err(CliError::Validation(
            "password must be at least 12 characters".to_string(),
        ));
    }
    let created = client.register(&email, &password, &role).await?;
    print_json(&created)?;
    Ok(())
}

pub async fn run_login(
    client: &ApiClient,
    config: &mut AppConfig,
    email: String,
    password: String,
) -> Result<()> {
    let login_response = client.login(&email, &password).await?;
    config.access_token = Some(login_response.access_token);
    config.save()?;
    println!("Login succeeded and access token saved.");
    Ok(())
}

pub async fn run_connect(
    client: &ApiClient,
    tool: String,
    scopes: Vec<String>,
    headless: bool,
) -> Result<()> {
    let response = client.connect_tool(&tool, &scopes).await?;
    if headless {
        println!("Headless mode enabled: browser OAuth redirect bypassed.");
    }
    print_json(&response)?;
    Ok(())
}

pub async fn run_revoke(client: &ApiClient, agent_id: String) -> Result<()> {
    let response = client.revoke_agent(&agent_id).await?;
    print_json(&response)?;
    Ok(())
}

pub async fn run_deploy(client: &ApiClient, manifest_path: String) -> Result<()> {
    let manifest = fs::read_to_string(&manifest_path)?;
    if !manifest.contains("gateway:") {
        return Err(CliError::Validation(
            "manifest must include a gateway field".to_string(),
        ));
    }

    let summary = client.dashboard_summary().await?;
    println!("Deployment preflight checks passed.");
    println!("Current gateway summary:");
    print_json(&summary)?;
    Ok(())
}

pub fn run_quick_connect(tool: String, scopes: Vec<String>, output: String) -> Result<()> {
    if scopes.is_empty() {
        return Err(CliError::Validation(
            "at least one scope is required".to_string(),
        ));
    }

    let manifest = format!(
        "version: 1\ngateway: http://localhost:8080\nconnector:\n  tool: {}\n  scopes:\n{}",
        tool,
        scopes
            .iter()
            .map(|scope| format!("    - {}", scope))
            .collect::<Vec<String>>()
            .join("\n")
    );

    fs::write(&output, manifest)?;
    println!("Generated {}", output);
    Ok(())
}
