var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
export default (options) => ({
    fields: options.fields || [],
    method: options.method || 'POST',
    submitUrl: options.submitUrl || '',
    formData: {},
    errors: {},
    isSubmitting: false,
    submitted: false,
    submissionError: false,
    init() {
        // Initialize formData with empty strings for each field
        this.fields.forEach((field) => {
            this.formData[field.name] = '';
        });
    },
    validateField(fieldName) {
        const field = this.fields.find((f) => f.name === fieldName);
        if (!field)
            return; // Field not found
        if (field.required && !this.formData[fieldName]) {
            this.errors[fieldName] = `${field.label} is required.`;
            return;
        }
        if (field.type === 'email' &&
            this.formData[fieldName] &&
            !this.isValidEmail(this.formData[fieldName])) {
            this.errors[fieldName] = 'Invalid email format.';
        }
        // Add more validation rules as needed (e.g., regex, length checks)
        // If we reach here, the field is valid, so remove any existing error:
        delete this.errors[fieldName];
    },
    isValidEmail(email) {
        // Basic email validation regex
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    },
    validateForm() {
        this.errors = {}; // Clear all errors
        this.fields.forEach((field) => {
            this.validateField(field.name);
        });
        return Object.keys(this.errors).length === 0; // Return true if no errors
    },
    submitForm() {
        return __awaiter(this, void 0, void 0, function* () {
            if (!this.validateForm()) {
                console.log('invalid', this.errors);
                return; // Don't submit if there are validation errors
            }
            this.isSubmitting = true;
            this.submissionError = false;
            try {
                const res = yield fetch(this.submitUrl, {
                    method: this.method,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.formData),
                });
                if (!res.ok) {
                    throw new Error('Network response was not ok');
                }
                const data = yield res.json();
                console.log('Success:', data);
                this.submitted = true;
            }
            catch (error) {
                console.error(error);
                this.submissionError = true;
            }
            finally {
                this.isSubmitting = false;
            }
            console.log('Form Data:', this.formData);
        });
    },
});
