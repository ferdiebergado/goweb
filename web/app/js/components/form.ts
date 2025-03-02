import { APIResponse } from '../types';

type ValidateFn = () => FormErrors;
type HandleFn = <T>(data: T) => undefined;
type OnErrorFn = (error: unknown) => undefined;
type FinallyFn = () => undefined;
type Method = 'POST' | 'PUT' | 'PATCH' | 'DELETE';

interface FormData {
  [key: string]: string;
}

export interface FormErrors {
  [key: string]: string;
}

export interface FormOptions {
  data: FormData;
  method?: Method;
  submitUrl: string;
  validate: ValidateFn;
  handle: HandleFn;
  onError: OnErrorFn;
  onFinal: FinallyFn;
}

export const form = (options: FormOptions) => ({
  data: options.data || [],
  method: options.method ?? 'POST',
  submitUrl: options.submitUrl || '',
  validate: options.validate,
  handle: options.handle,
  errorFn: options.onError,
  finalFn: options.onFinal,
  errors: {} as FormErrors,
  isSubmitting: false,
  submitted: false,
  submissionError: false,

  validateForm(): boolean {
    this.errors = this.validate();
    return Object.keys(this.errors).length === 0;
  },

  async submitForm(): Promise<undefined> {
    if (!this.validateForm()) {
      console.log('invalid input', this.errors);

      return;
    }

    this.isSubmitting = true;
    this.submissionError = false;

    try {
      const res = await fetch(this.submitUrl, {
        method: this.method,
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(this.data),
      });

      if (!res.ok) {
        if (res.status in [400, 422]) {
          const data: APIResponse<undefined> = await res.json();
          this.errors = data.errors;
          return;
        }
        throw new Error('Network response was not ok');
      }

      const data = await res.json();

      console.log('Success:', data);
      if (this.handle) this.handle(data);

      this.submitted = true;
    } catch (error) {
      console.error(error);
      if (this.errorFn) this.errorFn(error);
      this.submissionError = true;
    } finally {
      if (this.finalFn) this.finalFn();
      this.isSubmitting = false;
    }

    console.log('Form Data:', this.data);
  },
});
