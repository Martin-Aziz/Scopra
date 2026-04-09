mod client;
mod commands;
mod config;
mod error;
mod utils;

use clap::{Parser, Subcommand};
use client::api_client::ApiClient;
use commands::actions;
use config::app_config::AppConfig;
use error::Result;

#[derive(Debug, Parser)]
#[command(name = "nexus", version, about = "NEXUS-MCP CLI")]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Debug, Subcommand)]
enum Commands {
    Init {
        #[arg(long, default_value = "http://localhost:8080")]
        gateway_url: String,
    },
    Status,
    Register {
        #[arg(long)]
        email: String,
        #[arg(long)]
        password: String,
        #[arg(long, default_value = "user")]
        role: String,
    },
    Login {
        #[arg(long)]
        email: String,
        #[arg(long)]
        password: String,
    },
    Connect {
        tool: String,
        #[arg(long, value_delimiter = ',')]
        scopes: Vec<String>,
        #[arg(long, default_value_t = false)]
        headless: bool,
    },
    QuickConnect {
        tool: String,
        #[arg(long, value_delimiter = ',')]
        scopes: Vec<String>,
        #[arg(long, default_value = "nexus.yaml")]
        output: String,
    },
    Deploy {
        #[arg(long, default_value = "nexus.yaml")]
        manifest: String,
    },
    Revoke {
        agent_id: String,
    },
}

#[tokio::main]
async fn main() {
    if let Err(error) = run().await {
        eprintln!("nexus error: {error}");
        std::process::exit(1);
    }
}

async fn run() -> Result<()> {
    let cli = Cli::parse();
    let mut config = AppConfig::load_or_default()?;

    match cli.command {
        Commands::Init { gateway_url } => actions::run_init(&mut config, gateway_url),
        Commands::Status => {
            let client = ApiClient::new(config.gateway_url.clone(), config.access_token.clone());
            actions::run_status(&client).await
        }
        Commands::Register {
            email,
            password,
            role,
        } => {
            let client = ApiClient::new(config.gateway_url.clone(), config.access_token.clone());
            actions::run_register(&client, email, password, role).await
        }
        Commands::Login { email, password } => {
            let client = ApiClient::new(config.gateway_url.clone(), config.access_token.clone());
            actions::run_login(&client, &mut config, email, password).await
        }
        Commands::Connect {
            tool,
            scopes,
            headless,
        } => {
            let client = ApiClient::new(config.gateway_url.clone(), config.access_token.clone());
            actions::run_connect(&client, tool, scopes, headless).await
        }
        Commands::QuickConnect {
            tool,
            scopes,
            output,
        } => actions::run_quick_connect(tool, scopes, output),
        Commands::Deploy { manifest } => {
            let client = ApiClient::new(config.gateway_url.clone(), config.access_token.clone());
            actions::run_deploy(&client, manifest).await
        }
        Commands::Revoke { agent_id } => {
            let client = ApiClient::new(config.gateway_url.clone(), config.access_token.clone());
            actions::run_revoke(&client, agent_id).await
        }
    }
}
