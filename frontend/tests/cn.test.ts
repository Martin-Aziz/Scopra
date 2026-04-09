import { describe, expect, it } from "vitest";
import { cn } from "../src/utils/cn";

describe("cn", () => {
  it("merges class values and removes tailwind conflicts", () => {
    const className = cn("p-2", "p-4", "text-sm", false && "hidden");
    expect(className).toContain("p-4");
    expect(className).toContain("text-sm");
    expect(className).not.toContain("p-2");
  });
});
