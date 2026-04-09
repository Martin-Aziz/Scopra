export type UserRole = "admin" | "user";

export type User = {
  id: string;
  email: string;
  role: UserRole;
  isActive: boolean;
  createdAt: string;
};

export type AgentStatus = "active" | "revoked";

export type Agent = {
  id: string;
  name: string;
  status: AgentStatus;
  createdAt: string;
  updatedAt: string;
};

export type ToolSlug = "github" | "slack" | "jira" | "notion" | "linear";

export type ToolCallRequest = {
  agentId: string;
  tool: ToolSlug;
  action: string;
  payload: Record<string, unknown>;
  destructive: boolean;
};

export type ToolCallResponse = {
  status: "success" | "pending";
  requiresApproval: boolean;
  approvalRequestId?: string;
  message: string;
  result?: Record<string, unknown>;
};

export type AuditEvent = {
  id: number;
  agentId: string;
  tool: string;
  action: string;
  status: string;
  correlationId: string;
  createdAt: string;
};
