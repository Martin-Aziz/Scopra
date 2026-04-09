use crate::error::{CliError, Result};
use reqwest::{Client, RequestBuilder};
use serde::Deserialize;
use serde_json::{json, Value};

#[derive(Debug, Clone)]
pub struct ApiClient {
    base_url: String,
    client: Client,
    access_token: Option<String>,
}

#[derive(Debug, Deserialize)]
pub struct LoginResponse {
    #[serde(rename = "accessToken")]
    pub access_token: String,
}

impl ApiClient {
    pub fn new(base_url: String, access_token: Option<String>) -> Self {
        Self {
            base_url,
            client: Client::new(),
            access_token,
        }
    }

    pub async fn status(&self) -> Result<Value> {
        self.send_json(self.client.get(format!("{}/health", self.base_url)))
            .await
    }

    pub async fn register(&self, email: &str, password: &str, role: &str) -> Result<Value> {
        self.send_json(
            self.client
                .post(format!("{}/api/v1/auth/register", self.base_url))
                .json(&json!({"email": email, "password": password, "role": role})),
        )
        .await
    }

    pub async fn login(&self, email: &str, password: &str) -> Result<LoginResponse> {
        let response = self
            .client
            .post(format!("{}/api/v1/auth/login", self.base_url))
            .json(&json!({"email": email, "password": password}))
            .send()
            .await?;

        if !response.status().is_success() {
            let body = response
                .text()
                .await
                .unwrap_or_else(|_| "unable to read error body".to_string());
            return Err(CliError::Api(body));
        }

        Ok(response.json::<LoginResponse>().await?)
    }

    pub async fn connect_tool(&self, tool: &str, scopes: &[String]) -> Result<Value> {
        let request = self
            .client
            .post(format!(
                "{}/api/v1/connectors/{}/connect",
                self.base_url, tool
            ))
            .json(&json!({"scopes": scopes}));
        self.send_json(self.with_auth(request)?).await
    }

    pub async fn revoke_agent(&self, agent_id: &str) -> Result<Value> {
        let request = self.client.post(format!(
            "{}/api/v1/agents/{}/revoke",
            self.base_url, agent_id
        ));
        self.send_json(self.with_auth(request)?).await
    }

    pub async fn dashboard_summary(&self) -> Result<Value> {
        let request = self
            .client
            .get(format!("{}/api/v1/dashboard/summary", self.base_url));
        self.send_json(self.with_auth(request)?).await
    }

    fn with_auth(&self, request: RequestBuilder) -> Result<RequestBuilder> {
        let token = self.access_token.as_ref().ok_or_else(|| {
            CliError::Config("access token not found. Run login first.".to_string())
        })?;
        Ok(request.bearer_auth(token))
    }

    async fn send_json(&self, request: RequestBuilder) -> Result<Value> {
        let response = request.send().await?;
        if !response.status().is_success() {
            let body = response
                .text()
                .await
                .unwrap_or_else(|_| "unable to read error body".to_string());
            return Err(CliError::Api(body));
        }
        Ok(response.json::<Value>().await?)
    }
}
