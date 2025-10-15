import { z } from 'zod';

/**
 * Creates a reactive form validation composable for client-side form management.
 * 
 * WHY THIS APPROACH:
 * - Provides real-time validation feedback for better UX
 * - Manages form state, errors, and touched fields automatically
 * - Uses Svelte 5 runes ($state, $derived) for optimal reactivity
 * - Separates UI concerns from server-side validation
 * 
 * USE CASE:
 * - Forms that need immediate user feedback (login, register, settings)
 * - Multi-step forms with complex validation requirements
 * - Forms where users need to see errors before submission
 * 
 * @param schema - Zod schema for validation rules
 * @returns Reactive form state management object
 * 
 * @example
 * ```typescript
 * const loginSchema = z.object({
 *   username: z.string().min(1, 'Required'),
 *   password: z.string().min(1, 'Required')
 * });
 * 
 * const form = createFormValidation(loginSchema);
 * // Now use form.formData, form.errors, form.validateField(), etc.
 * ```
 */
export function createFormValidation<T extends z.ZodType>(schema: T) {
	const formData = $state<Partial<z.infer<T>>>({} as Partial<z.infer<T>>);
	const errors = $state<Record<string, string>>({});
	const touched = $state<Record<string, boolean>>({});

	/**
 * Validates a single form field and updates reactive state.
 * 
 * WHY: Provides immediate feedback as users type, improving UX.
 * Marks field as "touched" to show errors only after user interaction.
 * 
 * @param field - Form field name to validate
 * @param value - Current field value (can be undefined from input binding)
 */
	function validateField(field: string, value: FormDataEntryValue | null | undefined) {
		touched[field] = true;
		const result = schema.safeParse({ ...formData, [field]: value || null });
		
		if (!result.success) {
			const fieldError = result.error.issues.find(i => i.path[0] === field);
			errors[field] = fieldError?.message || '';
		} else {
			delete errors[field];
		}
	}

	/**
	 * Validates the entire form against the schema.
	 * 
	 * WHY: Used before form submission to ensure all data is valid.
	 * Shows all validation errors at once and marks fields as touched.
	 * 
	 * @returns true if form is valid, false otherwise
	 */
	function validateForm() {
		const result = schema.safeParse(formData);
		if (!result.success) {
			result.error.issues.forEach(issue => {
				const field = issue.path[0] as string;
				errors[field] = issue.message;
				touched[field] = true;
			});
			return false;
		}
		// Clear all errors
		Object.keys(errors).forEach(key => {
			delete errors[key];
		});
		return true;
	}

	const isValid = $derived(Object.keys(errors).length === 0);

	return {
		formData,
		errors,
		touched,
		validateField,
		validateForm,
		isValid
	};
}
