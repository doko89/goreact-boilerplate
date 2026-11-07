import { z } from "zod";
import { api } from "../../shared/lib/api";

export const userSchema = z.object({
  id: z.string(),
  username: z.string(),
});
export type User = z.infer<typeof userSchema>;

export const loginSchema = z.object({
  username: z.string().min(1, "username is required").max(64),
  password: z.string().min(1, "password is required").max(128),
});
export type LoginInput = z.infer<typeof loginSchema>;

export async function loginRequest(input: LoginInput): Promise<User> {
  const res: unknown = await api("/api/auth/login", { method: "POST", body: input });
  return userSchema.parse(res);
}

export async function logoutRequest(): Promise<void> {
  await api("/api/auth/logout", { method: "POST" });
}

export async function meRequest(): Promise<User> {
  const res: unknown = await api("/api/auth/me");
  return userSchema.parse(res);
}
