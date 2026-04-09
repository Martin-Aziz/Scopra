# Technical Debt Register

## High Priority
1. Implement real OAuth 2.1 + PKCE flows for each supported connector and store encrypted provider tokens.
2. Introduce policy-as-code for resource-level authorization across tool calls.
3. Raise test coverage to >=90% and enforce coverage threshold in CI.

## Medium Priority
1. Convert Security Audit CI step from non-blocking to blocking once dependency remediation is complete.
2. Add distributed tracing and request replay tooling for audit investigations.
3. Add SSO/SAML integration and enterprise provisioning hooks.

## Low Priority
1. Replace connector stubs with resilient retries and provider-specific error mapping.
2. Introduce event streaming (Redis streams/Kafka) for high-volume audit fan-out.
