import { APIResponse } from '../@types/api';
import { FormErrors, FormOptions, FormValues } from '../@types/form';

export default function <T extends FormValues, E extends FormErrors>(
  opts: FormOptions<T, E>
) {
  const { data, method, submitUrl, errors, validateFn, onSuccess, onError } =
    opts;

  return {
    data,
    method: method ?? 'POST',
    submitUrl,
    isSubmitting: false,
    errors,
    validate() {
      this.errors = validateFn.call(this) as E;
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
          if (response.status in [400, 422]) {
            const data: APIResponse<undefined> = await response.json();
            const { message, errors } = data;
            console.error(message);
            this.errors = errors as E;
            return;
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
