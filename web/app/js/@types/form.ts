import { APIResponse } from './api';

type FormValues = Record<string, string | number | boolean>;
type FormErrors<T extends FormValues> = {
  [K in keyof T]?: string;
};

interface FormOptions {
  data: FormValues;
  method?: 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  submitUrl: string;
  errors: FormErrors<FormValues>;
  validateFn(): FormErrors<FormValues>;
  onSuccess<K>(data: APIResponse<K>): undefined;
  onError(error: unknown): undefined;
}

export { FormOptions, FormValues, FormErrors };
