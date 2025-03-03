import { APIResponse } from './api';

type FormValues = Record<string, string | number | boolean>;
type FormErrors = Record<string, string | undefined>;

interface FormOptions<T extends FormValues, E extends FormErrors> {
  data: T;
  method?: 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  submitUrl: string;
  errors: E;
  validateFn(): FormErrors;
  onSuccess<K>(data: APIResponse<K>): undefined;
  onError(error: unknown): undefined;
}

export { FormValues, FormErrors, FormOptions };
