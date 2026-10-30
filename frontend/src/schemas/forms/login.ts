import * as z from "zod";

export const LoginFormDataSchema = z.object({
  password: z.string().min(1, "Required"),
});
export type LoginFormData = z.infer<typeof LoginFormDataSchema>;
