import { APIResponse } from '../@types/api';
import { FormOptions } from '../@types/form';

export default function (opts: FormOptions) {
  const { data, method, submitUrl, errors, validateFn, onSuccess, onError } =
    opts;

  return {
    data,
    method: method ?? 'POST',
    submitUrl,
    isSubmitting: false,
    errors,
    validate() {
      this.errors = validateFn.call(this);
      return Object.keys(this.errors).length === 0;
    },
    async submit() {
      if (!this.validate()) return;

      this.isSubmitting = true;

      try {
        const response = await fetch(this.submitUrl, {
          method: this.method,
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(this.data),
        });

        if (!response.ok) {
          const { status } = response;
          if (status === 400 || status === 422) {
            const data: APIResponse<undefined> = await response.json();
            const { message, errors } = data;
            this.errors = errors;
            throw new Error(message);
          }
          throw new Error('Invalid credentials');
        }

        const data = await response.json();
        onSuccess(data);
      } catch (error) {
        console.error(error);
        onError(error);
      } finally {
        this.isSubmitting = false;
      }
    },
  };
}
