import { z } from "zod";

// The admin profile page's forms. Here rather than in the page so an app, a
// mobile client or a test validates a profile change the way the panel does.

export const PersonalInfoSchema = z.object({
  first_name: z.string().min(2, "First name must be at least 2 characters"),
  last_name: z.string().min(2, "Last name must be at least 2 characters"),
  email: z.string().email("Please enter a valid email"),
  // Only needed when the email changes.
  current_password: z.string().optional(),
});
export type PersonalInfoInput = z.infer<typeof PersonalInfoSchema>;

export const ProfessionalInfoSchema = z.object({
  job_title: z.string().optional().default(""),
  bio: z.string().optional().default(""),
});
export type ProfessionalInfoInput = z.infer<typeof ProfessionalInfoSchema>;

export const ChangePasswordSchema = z
  .object({
    current_password: z.string().min(1, "Enter your current password"),
    password: z.string().min(8, "Password must be at least 8 characters"),
    confirm_password: z.string().min(1, "Please confirm your password"),
  })
  .refine((d) => d.password === d.confirm_password, {
    message: "Passwords do not match",
    path: ["confirm_password"],
  });
export type ChangePasswordInput = z.infer<typeof ChangePasswordSchema>;
