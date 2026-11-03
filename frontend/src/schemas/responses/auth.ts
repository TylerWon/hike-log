import * as z from "zod";

export const CheckAuthResponseSchema = z.object({
  authenticated: z.boolean(),
});
export type CheckAuthResponse = z.infer<typeof CheckAuthResponseSchema>;
