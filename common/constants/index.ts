export const ACCESS_TOKEN_TTL_MINUTES = 15;
export const REFRESH_TOKEN_TTL_HOURS = 168;

export const SUPPORTED_TOOLS = ["github", "slack", "jira", "notion", "linear"] as const;

export const AUDIT_STATUSES = {
  SUCCESS: "success",
  FAILED: "failed",
  PENDING_APPROVAL: "pending_approval"
} as const;

export const APPROVAL_WINDOW_MINUTES = 15;
