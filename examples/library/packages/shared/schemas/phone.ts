import { z } from "zod";
import { isValidPhoneNumber } from "libphonenumber-js/min";

/**
 * A phone number in E.164 (+256772123456), checked against its country's
 * number lengths with libphonenumber's small ("min") metadata.
 *
 * Small on purpose: every page that imports this package's schemas carries it,
 * and the full metadata is twice the size. The admin's phone input loads the
 * full metadata on its own and the API applies the full rules, so a number
 * with a plausible length but an unassigned prefix is still refused.
 */
export const PhoneSchema = z
  .string()
  .trim()
  .regex(/^\+[1-9][0-9]{6,14}$/, "Enter the number with its country code, such as +256772123456")
  .refine((v) => isValidPhoneNumber(v), "This is not a valid phone number for its country");
