interface FormField {
  name: string;
  label: string;
  type: 'text' | 'email' | 'textarea';
  required: boolean;
}

interface FormOptions {
  fields: FormField[];
  method: 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  submitUrl: string;
}

interface FormData {
  [key: string]: string;
}

interface FormErrors {
  [key: string]: string;
}

export default (options: FormOptions) => ({
  fields: options.fields || [],
  method: options.method || 'POST',
  submitUrl: options.submitUrl || '',
  formData: {} as FormData,
  errors: {} as FormErrors,
  isSubmitting: false,
  submitted: false,
  submissionError: false,

  init() {
    // Initialize formData with empty strings for each field
    this.fields.forEach((field) => {
      this.formData[field.name] = '';
    });
  },

  validateField(fieldName: string): void {
    const field = this.fields.find((f) => f.name === fieldName);
    if (!field) return; // Field not found

    if (field.required && !this.formData[fieldName]) {
      this.errors[fieldName] = `${field.label} is required.`;
      return;
    }

    if (
      field.type === 'email' &&
      this.formData[fieldName] &&
      !this.isValidEmail(this.formData[fieldName])
    ) {
      this.errors[fieldName] = 'Invalid email format.';
    }
    // Add more validation rules as needed (e.g., regex, length checks)

    // If we reach here, the field is valid, so remove any existing error:
    delete this.errors[fieldName];
  },

  isValidEmail(email: string): boolean {
    // Basic email validation regex
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  },

  validateForm(): boolean {
    this.errors = {}; // Clear all errors
    this.fields.forEach((field) => {
      this.validateField(field.name);
    });
    return Object.keys(this.errors).length === 0; // Return true if no errors
  },

  async submitForm(): Promise<void> {
    if (!this.validateForm()) {
      console.log('invalid', this.errors);

      return; // Don't submit if there are validation errors
    }

    this.isSubmitting = true;
    this.submissionError = false;

    try {
      const res = await fetch(this.submitUrl, {
        method: this.method,
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(this.formData),
      });

      if (!res.ok) {
        throw new Error('Network response was not ok');
      }

      const data = await res.json();

      console.log('Success:', data);
      this.submitted = true;
    } catch (error) {
      console.error(error);
      this.submissionError = true;
    } finally {
      this.isSubmitting = false;
    }

    console.log('Form Data:', this.formData);
  },
});
