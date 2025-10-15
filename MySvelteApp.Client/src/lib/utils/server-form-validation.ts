import { error } from '@sveltejs/kit';
import { z } from 'zod';

/**
 * Server-side form validation utilities for SvelteKit remote functions.
 * 
 * WHY THIS APPROACH:
 * - Provides data integrity and security validation on the server
 * - Handles experimental SvelteKit remote function formData type issues
 * - Separates server validation concerns from client-side UI state
 * - Ensures consistent validation logic across all form handlers
 * 
 * USE CASE:
 * - Authentication forms (login, register, password reset)
 * - Data submission forms that need server validation
 * - Any form where security and data integrity are critical
 */

/**
 * Helper to safely extract string value from FormData.
 * 
 * WHY: FormDataEntryValue can be string or File - we only want strings for validation.
 * 
 * @param value - FormData entry value that might be string or File
 * @returns Extracted string or empty string
 */
export const getString = (value: FormDataEntryValue | null) => (typeof value === 'string' ? value : '');

/**
 * Type for SvelteKit experimental remote function formData.get method.
 * 
 * WHY: Experimental remote functions have type definition issues.
 * This provides a proper type for the formData.get method.
 */
export type FormGetter = (key: string) => FormDataEntryValue | null;

/**
 * Type for experimental SvelteKit remote function formData object.
 * 
 * WHY: Experimental remote functions return formData with incomplete types.
 * This provides the correct interface for form data access.
 */
export type ExperimentalFormData = {
  get: (key: string) => FormDataEntryValue | null;
};

/**
 * Validates and extracts form data using a Zod schema.
 * 
 * WHY THIS APPROACH:
 * - Automatically extracts form data based on schema keys
 * - Handles experimental remote function type issues gracefully  
 * - Provides consistent error handling across all forms
 * - Returns validated, typed data for direct use
 * 
 * @param formData - The formData object from SvelteKit's experimental form function
 * @param schema - Zod schema for validation rules
 * @returns Validated and typed data ready for use
 * @throws 400 error with specific validation message if validation fails
 * 
 * @example
 * ```typescript
 * export const login = form(async (formData) => {
 *   const { username, password } = validateFormData(
 *     formData as unknown as ExperimentalFormData, 
 *     zLoginForm
 *   );
 *   // Now username and password are typed string values
 * });
 * ```
 */
export function validateFormData<T extends z.ZodType>(
  formData: ExperimentalFormData, 
  schema: T
): z.infer<T> {
  const get = formData.get;
  
  // Extract data based on schema keys
  const data: Record<string, FormDataEntryValue | null> = {};
  const keys = Object.keys((schema as unknown as z.ZodObject<Record<string, z.ZodType>>).shape);
  
  keys.forEach(key => {
    data[key] = getString(get(key));
  });
  
  // Validate with schema
  const parsed = schema.safeParse(data);
  
  if (!parsed.success) {
    const message = parsed.error.issues[0]?.message ?? 'Invalid form data';
    throw error(400, { message });
  }
  
  return parsed.data;
}
