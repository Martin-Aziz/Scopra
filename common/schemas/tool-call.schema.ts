import { z } from "zod";
import { SUPPORTED_TOOLS } from "../constants";

export const toolCallSchema = z
  .object({
    agentId: z.string().uuid(),
    tool: z.enum(SUPPORTED_TOOLS),
    action: z.string().min(1),
    payload: z.record(z.unknown()),
    destructive: z.boolean()
  })
  .strict();

export type ToolCallSchema = z.infer<typeof toolCallSchema>;
